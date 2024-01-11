package main

import (
	"fmt"
	"sync"
)

var ch chan int

func producer() {
	for i := 0; i < 5; i++ {
		ch <- i
	}

	close(ch)
}

func consumer() {
	for {
		value, ok := <-ch
		if !ok {
			break
		}

		fmt.Println(value)
	}
}

func main() {
	ch = make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		producer()
	}()

	go func() {
		defer wg.Done()
		consumer()
	}()

	wg.Wait()
}
