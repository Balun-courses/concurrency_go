package main

import "sync"

// Need to show solution

type Counters struct {
	mu sync.Mutex
	m  map[string]int
}

func (c *Counters) Load(key string) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, found := c.m[key]
	return value, found
}

func (c *Counters) Store(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[key] = value
}
