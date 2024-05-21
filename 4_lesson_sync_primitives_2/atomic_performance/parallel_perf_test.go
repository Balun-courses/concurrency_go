package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// go test -bench=. parallel_perf_test.go

func BenchmarkParallelMutexAdd(b *testing.B) {
	var number int32
	var mutex sync.Mutex

	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				mutex.Lock()
				number++
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()
}

func BenchmarkParallelAtomicAdd(b *testing.B) {
	var number atomic.Int32

	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				number.Add(1)
			}
		}()
	}

	wg.Wait()
}
