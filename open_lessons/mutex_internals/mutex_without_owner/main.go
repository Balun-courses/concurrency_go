package main

import "sync"

func main() {
	mutex := sync.Mutex{}
	mutex.Lock()

	/*
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			mutex.Unlock()
		}()

		wg.Wait()
	*/

	mutex.Lock()
	mutex.Unlock()
}
