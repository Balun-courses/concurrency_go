package main

import "sync"

type Map struct {
	once   sync.Once
	values map[interface{}]interface{}
}

func NewMap() *Map {
	return &Map{}
}

func (m *Map) Add(key, value interface{}) {
	m.init()
	m.values[key] = value
}

func (m *Map) init() {
	m.once.Do(func() {
		m.values = make(map[interface{}]interface{})
	})
}
