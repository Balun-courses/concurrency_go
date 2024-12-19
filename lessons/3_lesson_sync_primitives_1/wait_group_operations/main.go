package main

import "sync"

func makeNegativeCounter() {
	wg := sync.WaitGroup{}
	wg.Add(-10)
}

func waitZeroCounter() {
	wg := sync.WaitGroup{}
	wg.Wait()
}

func main() {
}
