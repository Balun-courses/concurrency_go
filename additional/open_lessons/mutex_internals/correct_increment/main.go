package main

import (
	"fmt"
	"sync"
)

const goroutinesNumber = 1000

func main() {
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(goroutinesNumber)

	value := 0
	for i := 0; i < goroutinesNumber; i++ {
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
