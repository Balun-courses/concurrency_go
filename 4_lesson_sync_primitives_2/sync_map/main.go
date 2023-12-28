package main

import "sync"

type Counters struct {
	m sync.Map
}

func (c *Counters) Load(key string) (int, bool) {
	value, found := c.m.Load(key)
	return value.(int), found
}

func (c *Counters) Store(key string, value int) {
	c.m.Store(key, value)
}
