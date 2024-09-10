package main

import (
	"fmt"
	"sync"
	"time"
)

func fanIn(channels ...chan int) chan int {
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	merged := make(chan int)
	for _, channel := range channels {
		go func(channel chan int) {
			defer wg.Done()
			for value := range channel {
				merged <- value
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		defer func() {
			close(ch1)
			close(ch2)
			close(ch3)
		}()

		for i := 0; i < 100; i += 3 {
			ch1 <- i
			ch2 <- i + 1
			ch3 <- i + 2
		}
		time.Sleep(100 * time.Millisecond)
	}()

	for value := range fanIn(ch1, ch2, ch3) {
		fmt.Println(value)
	}
}
