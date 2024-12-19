package main

import "sync"

type Counters struct {
	mu sync.RWMutex
	m  map[string]int
}

func (c *Counters) Load(key string) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	return value, found
}

func (c *Counters) Store(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[key] = value
}
