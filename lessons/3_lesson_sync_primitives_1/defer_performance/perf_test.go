package main

import (
	"sync"
	"testing"
)

// go test -bench=. perf_test.go

func BenchmarkMutex(b *testing.B) {
	var counter int64
	var mutex sync.Mutex
	for j := 0; j < b.N; j++ {
		func() {
			mutex.Lock()
			counter++
			mutex.Unlock()
		}()
	}
}

func BenchmarkMutexWithDefer(b *testing.B) {
	var counter int64
	var mutex sync.Mutex
	for j := 0; j < b.N; j++ {
		func() {
			mutex.Lock()
			defer mutex.Unlock()
			counter++
		}()
	}
}
