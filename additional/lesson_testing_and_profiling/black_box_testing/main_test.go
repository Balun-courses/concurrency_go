package main

import (
	"sync"
	"testing"
)

func ParallelFunction() int {
	var wg sync.WaitGroup
	var result int
	numWorkers := 4

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			defer wg.Done()
			result += id
		}(i)
	}

	wg.Wait()
	return result
}

func TestParallelFunction(t *testing.T) {
	expected := 6
	result := ParallelFunction()

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}
