package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	deadline time.Time
	value    int
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

func (c *Cache) Set(key string, val int, ttl time.Duration) {
	c.mu.Lock()
	c.items[key] = CacheItem{
		deadline: time.Now().Add(ttl),
		value:    val,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) *int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Now().Before(c.items[key].deadline) {
		v := c.items[key]
		return &v.value
	}

	return nil
}

func (c *Cache) Del(key string) {
	delete(c.items, key)
}
