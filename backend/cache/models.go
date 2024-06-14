package cache

import (
	"container/list"
	"sync"
	"time"
)

type CacheItem struct {
	value     interface{}
	timestamp time.Time
	ttl       time.Duration
}

type CacheEntry struct {
	key   string
	value *CacheItem
}

type LRUCache struct {
	capacity int
	ttl      time.Duration
	items    map[string]*list.Element
	order    *list.List
	mu       sync.Mutex
	stopCh   chan struct{}
}
