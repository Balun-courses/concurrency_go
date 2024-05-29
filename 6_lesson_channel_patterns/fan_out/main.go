package main

import (
	"fmt"
	"sync"
	"time"
)

func SplitChannel(inputCh <-chan int, n int) []chan int {
	if n <= 0 {
		n = 1
	}

	outputCh := make([]chan int, n)
	for i := 0; i < n; i++ {
		outputCh[i] = make(chan int)
	}

	go func() {
		idx := 0
		for value := range inputCh {
			outputCh[idx] <- value
			idx = (idx + 1) % n
		}

		for _, ch := range outputCh {
			close(ch)
		}
	}()

	return outputCh
}

func main() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
	}()

	channels := SplitChannel(ch, 2)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for value := range channels[0] {
			fmt.Println("ch1: ", value)
		}
	}()

	go func() {
		defer wg.Done()
		for value := range channels[1] {
			fmt.Println("ch2: ", value)
		}
	}()

	wg.Wait()
}
