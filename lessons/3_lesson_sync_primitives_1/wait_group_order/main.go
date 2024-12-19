package main

import (
	"log"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("test")
		}()

		wg.Wait()
	}
}
