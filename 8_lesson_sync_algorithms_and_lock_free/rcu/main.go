package main

import (
	"context"
	"sync"
	"time"
)

// Need to show solution

type Cache struct {
	mutex sync.RWMutex
	data  map[string]string
}

func NewCache(ctx context.Context) *Cache {
	cache := &Cache{
		data: make(map[string]string),
	}

	cache.synchronize(ctx)
	return cache
}

func (c *Cache) synchronize(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var temp map[string]string
			// data = ... - get from remote storage

			c.mutex.Lock()
			c.data = temp
			c.mutex.Unlock()
		}
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, found := c.data[key]
	return value, found
}
