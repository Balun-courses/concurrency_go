package main

import "sync"

// go run -race main.go

var buffer chan int

func main() {
	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			buffer = make(chan int)
		}()
	}

	wg.Wait()
}
