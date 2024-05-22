package main

import (
	"fmt"
	"sync"
)

func producer(ch chan int) {
	defer close(ch)
	for i := 0; i < 5; i++ {
		ch <- i
	}
}

func consumer(ch chan int) {
	for value := range ch { // syntax sugar
		fmt.Println(value)
	}
}

func main() {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		producer(ch)
	}()

	go func() {
		defer wg.Done()
		consumer(ch)
	}()

	wg.Wait()
}
