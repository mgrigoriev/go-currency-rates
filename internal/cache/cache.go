package cache

import (
	"sync"
)

type Cache struct {
	rates map[string]float64
	mu    sync.RWMutex
}

func New() *Cache {
	return &Cache{
		rates: make(map[string]float64),
	}
}

func (c *Cache) Set(key string, value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rates[key] = value
}

func (c *Cache) Get(key string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.rates[key]

	return value, ok
}
