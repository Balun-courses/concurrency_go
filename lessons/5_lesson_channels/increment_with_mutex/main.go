package main

import (
	"fmt"
	"sync"
)

func main() {
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(2)

	value := 0
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			mutex.Lock()
			value++
			mutex.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println(value)
}
