package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Need to show solution

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1000)

	var value atomic.Int32
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			value.Add(1)
		}()
	}

	wg.Wait()

	fmt.Println(value)
}
