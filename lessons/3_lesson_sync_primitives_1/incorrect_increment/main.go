package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1000)

	value := 0
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			value++
		}()
	}

	wg.Wait()

	fmt.Println(value)
}
