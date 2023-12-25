package main

import (
	"fmt"
	"sync"
	"time"
)

// Starvation may be with processor, memory, file descriptors,
// connections with database and so on...

func main() {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	const runtime = 1 * time.Second

	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			mutex.Lock()
			time.Sleep(3 * time.Nanosecond)
			mutex.Unlock()
			count++
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			mutex.Lock()
			time.Sleep(1 * time.Nanosecond)
			mutex.Unlock()

			mutex.Lock()
			time.Sleep(1 * time.Nanosecond)
			mutex.Unlock()

			mutex.Lock()
			time.Sleep(1 * time.Nanosecond)
			mutex.Unlock()

			count++
		}

		fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()

	wg.Wait()
}
