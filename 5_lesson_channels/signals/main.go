package main

import "sync"

func notifier(signals chan<- struct{}) {
	signals <- struct{}{}
}

func subscriber(signals <-chan struct{}) {
	<-signals
}

func main() {
	signals := make(chan struct{})
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

	wg.Wait()
}
