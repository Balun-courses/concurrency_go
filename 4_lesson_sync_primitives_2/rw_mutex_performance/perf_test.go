package main

import (
	"sync"
	"testing"
)

// go test -bench=. perf_test.go

func BenchmarkMutexAdd(b *testing.B) {
	var number int32
	var mutex sync.Mutex
	for i := 0; i < b.N; i++ {
		mutex.Lock()
		number++
		mutex.Unlock()
	}
}

func BenchmarkRWMutexAdd(b *testing.B) {
	var number int32
	var mutex sync.RWMutex
	for i := 0; i < b.N; i++ {
		mutex.Lock()
		number++
		mutex.Unlock()
	}
}
