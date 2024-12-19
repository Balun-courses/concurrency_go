package main

import (
	"fmt"
	"sync"
)

// Need to show solution

var value int

func inc() {
	mutex := sync.Mutex{}

	mutex.Lock()
	value++
	mutex.Unlock()
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			inc()
		}()
	}

	wg.Wait()

	fmt.Println(value)
}
