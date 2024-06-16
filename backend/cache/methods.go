package cache

import (
	"lru-cache/socket"
	"time"
)

func (cache *LRUCache) startCleaner() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			socket.Broadcast(cache.Snapshot())
			cache.removeExpiredItems()
		case <-cache.stopCh:
			return
		}
	}
}

func (cache *LRUCache) removeExpiredItems() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	now := time.Now()
	for e := cache.order.Back(); e != nil; e = e.Prev() {
		entry := e.Value.(*CacheEntry)
		if now.Sub(entry.value.timestamp) > entry.value.ttl {
			cache.order.Remove(e)
			delete(cache.items, entry.key)
		} else {
			break
		}
	}
}

func (cache *LRUCache) Set(key string, value interface{}, ttl time.Duration) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if element, found := cache.items[key]; found {
		cache.order.MoveToFront(element)
		element.Value.(*CacheEntry).value.value = value
		element.Value.(*CacheEntry).value.timestamp = time.Now()
		return
	}

	if cache.order.Len() >= cache.capacity {
		oldest := cache.order.Back()
		if oldest != nil {
			cache.order.Remove(oldest)
			delete(cache.items, oldest.Value.(*CacheEntry).key)
		}
	}

	newItem := &CacheItem{
		value:     value,
		timestamp: time.Now(),
		ttl:       ttl,
	}

	entry := &CacheEntry{
		key:   key,
		value: newItem,
	}

	element := cache.order.PushFront(entry)
	cache.items[key] = element
}

func (cache *LRUCache) Get(key string) (interface{}, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if element, found := cache.items[key]; found {
		item := element.Value.(*CacheEntry).value
		if time.Since(item.timestamp) > item.ttl {
			cache.order.Remove(element)
			delete(cache.items, key)
			return nil, false
		}
		cache.order.MoveToFront(element)
		item.timestamp = time.Now()
		return item.value, true
	}
	return nil, false
}

func (cache *LRUCache) Delete(key string) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	element, found := cache.items[key]
	if found {
		cache.order.Remove(element)
		delete(cache.items, key)
	}
	return found
}

func (cache *LRUCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for k := range cache.items {
		delete(cache.items, k)
	}

	for e := cache.order.Front(); e != nil; {
		next := e.Next()
		cache.order.Remove(e)
		e = next
	}
}

func (cache *LRUCache) Snapshot() []interface{} {
	res := make([]interface{}, 0, len(cache.items))
	for e := cache.order.Front(); e != nil; e = e.Next() {
		value := make(map[string]interface{})
		value["key"] = e.Value.(*CacheEntry).key
		value["value"] = e.Value.(*CacheEntry).value.value
		value["ttl"] = e.Value.(*CacheEntry).value.ttl.Seconds()
		res = append(res, value)
	}
	return res
}

func (cache *LRUCache) Stop() {
	close(cache.stopCh)
}
