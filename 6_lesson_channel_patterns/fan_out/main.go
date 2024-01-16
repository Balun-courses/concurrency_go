package main

import (
	"fmt"
	"time"
)

func SplitChannel(inputCh <-chan int, n int) []chan int {
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

	for {
		select {
		case value, ok := <-channels[0]:
			if !ok {
				return
			}

			fmt.Println(value)
		case value, ok := <-channels[1]:
			if !ok {
				return
			}

			fmt.Println(value)
		}
	}
}
