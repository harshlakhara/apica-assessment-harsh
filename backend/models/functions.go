package cache

import (
	"container/list"
	"time"
)

func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		ttl:      ttl,
		items:    make(map[string]*list.Element),
		order:    list.New(),
		stopCh:   make(chan struct{}),
	}

	go cache.startCleaner()
	return cache
}
