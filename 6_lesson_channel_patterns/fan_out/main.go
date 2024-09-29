package main

import (
	"fmt"
	"sync"
	"time"
)

func fanOut(channelMain chan int, amount int) []chan int {
	if amount <= 0 {
		amount = 1
	}

	channelPull := make([]chan int, amount)
	for i := 0; i < amount; i++ {
		channelPull[i] = make(chan int)
	}

	go func() {
		defer func() {
			for _, closed := range channelPull {
				close(closed)
			}
		}()

		index := 0
		for value := range channelMain {
			channelPull[index] <- value
			index = (index + 1) % amount
		}
	}()

	return channelPull
}

func main() {
	wg := sync.WaitGroup{}
	amount := 2
	wg.Add(amount)

	channelMain := make(chan int)
	channelPull := fanOut(channelMain, amount)

	go func() {
		defer close(channelMain)
		for i := 0; i < 10; i++ {
			channelMain <- i
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for index, channel := range channelPull {
		go func(channel chan int, index int) {
			defer wg.Done()
			for value := range channel {
				fmt.Printf("ch%d: %d\n", index+1, value)
			}
		}(channel, index)
	}

	wg.Wait()
}
