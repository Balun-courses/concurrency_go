package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// go test -bench=. perf_test.go

type MutexCounter struct {
	value int32
	mutex sync.Mutex
}

func (c *MutexCounter) Increment(int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value++
}

func (c *MutexCounter) Get() int32 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.value
}

type AtomicCounter struct {
	value atomic.Int32
}

func (c *AtomicCounter) Increment(int) {
	c.value.Add(1)
}

func (c *AtomicCounter) Get() int32 {
	return c.value.Load()
}

type ShardedAtomicCounter struct {
	shards [10]AtomicCounter
}

func (c *ShardedAtomicCounter) Increment(idx int) {
	c.shards[idx].value.Add(1)
}

func (c *ShardedAtomicCounter) Get() int32 {
	var value int32
	for idx := 0; idx < 10; idx++ {
		value += c.shards[idx].Get()
	}

	return value
}

func BenchmarkAtomicCounter(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	counter := MutexCounter{}
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				counter.Increment(idx)
			}
		}(i)
	}

	wg.Wait()
}
