package main

import (
	"sync"
)

// Need to show solution

type Cache struct {
	mutex sync.Mutex
	data  map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = value
}

func (c *Cache) Get(key string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Size() > 0 {
		return c.data[key]
	}

	return ""
}

func (c *Cache) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return len(c.data)
}
