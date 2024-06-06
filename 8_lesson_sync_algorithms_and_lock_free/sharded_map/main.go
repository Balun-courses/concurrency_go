package main

import (
	"hash/fnv"
	"sync"
)

type Map struct {
	mutex sync.Mutex
	data  map[string]string
}

func NewMap() *Map {
	return &Map{
		data: make(map[string]string),
	}
}

func (m *Map) Get(key string) (string, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	value, found := m.data[key]
	return value, found
}

func (m *Map) Set(key, value string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = value
}

func (m *Map) Delete(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, key)
}

type ShardedMap struct {
	shards       []Map
	shardsNumber int
}

func NewShardedMap(shardsNumber int) *ShardedMap {
	return &ShardedMap{
		shards:       make([]Map, shardsNumber),
		shardsNumber: shardsNumber,
	}
}

func (sm *ShardedMap) Get(key string) (string, bool) {
	idx := sm.shardIdx(key)
	return sm.shards[idx].Get(key)
}

func (sm *ShardedMap) Set(key, value string) {
	idx := sm.shardIdx(key)
	sm.shards[idx].Set(key, value)
}

func (sm *ShardedMap) Delete(key string) {
	idx := sm.shardIdx(key)
	sm.shards[idx].Delete(key)
}

func (sm *ShardedMap) shardIdx(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))

	hash := int(h.Sum32())
	return hash % sm.shardsNumber
}
