package main

import "sync"

func waitWithoutLock() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.Wait()
}

func waitAfterSignal() {
	cond := sync.NewCond(&sync.Mutex{})

	cond.L.Lock()
	cond.Wait()
	cond.L.Unlock()
}

func main() {
}
