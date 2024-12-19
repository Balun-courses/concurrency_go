package main

import (
	"fmt"
	"sync"
)

func notifier(signals chan int) {
	close(signals)
}

func subscriber(signals chan int) {
	<-signals
	fmt.Println("signaled")
}

func main() {
	signals := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		notifier(signals)
	}()

	go func() {
		defer wg.Done()
		subscriber(signals)
	}()

	go func() {
		defer wg.Done()
		subscriber(signals)
	}()

	wg.Wait()
}
