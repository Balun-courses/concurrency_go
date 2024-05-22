package main

import (
	"fmt"
	"sync"
)

func notifier(signals chan<- struct{}) {
	signals <- struct{}{}
}

func subscriber(signals <-chan struct{}) {
	<-signals
	fmt.Println("signaled")
}

func main() {
	signals := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		notifier(signals)
	}()

	go func() {
		defer wg.Done()
		subscriber(signals)
	}()

	wg.Wait()
}
