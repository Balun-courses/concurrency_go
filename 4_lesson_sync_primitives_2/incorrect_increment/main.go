package main

import (
	"fmt"
	"sync"
)

// Need to show solution

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1000)

	var value int
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			value++
		}()
	}

	wg.Wait()

	fmt.Println(value)
}
