package cache

import "time"

func (cache *LRUCache) startCleaner() {
	ticker := time.NewTicker(cache.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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

func (cache *LRUCache) Set(key string, value interface{}) {
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
		ttl:       cache.ttl,
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

func (cache *LRUCache) Delete(key string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if element, found := cache.items[key]; found {
		cache.order.Remove(element)
		delete(cache.items, key)
	}
}

func (cache *LRUCache) Stop() {
	close(cache.stopCh)
}
