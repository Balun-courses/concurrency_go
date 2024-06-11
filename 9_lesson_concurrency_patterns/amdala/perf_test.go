package amdala

import (
	"runtime"
	"sync"
	"testing"
)

// go test -bench=. bench_test.go
// a = 1 / ((1 - P) + (P / S))

func calculate(parallelFactor int) {
	value := 0
	for i := 0; i < 1000000; i++ {
		value += i
	}

	wg := sync.WaitGroup{}
	wg.Add(parallelFactor)
	for i := 0; i < parallelFactor; i++ {
		go func() {
			defer wg.Done()

			localValue := 0
			for j := 0; j < 1000000/parallelFactor; j++ {
				localValue += j
			}
		}()
	}

	wg.Wait()
}

func BenchmarkCalculation(b *testing.B) {
	parallelFactor := 1
	runtime.GOMAXPROCS(parallelFactor)
	for i := 0; i < b.N; i++ {
		calculate(parallelFactor)
	}
}
