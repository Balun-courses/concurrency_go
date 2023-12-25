package main

import "sync"

func done(wg sync.WaitGroup) {
	wg.Done()
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	done(wg)
	wg.Wait()
}
