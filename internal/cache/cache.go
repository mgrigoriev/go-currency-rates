package cache

import (
	"sync"
)

type Cache struct {
	rates sync.Map
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Set(key string, value interface{}) {
	c.rates.Store(key, value)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.rates.Load(key)
}
