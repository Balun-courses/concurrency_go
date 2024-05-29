package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	interruptCh := make(chan os.Signal)
	signal.Notify(interruptCh, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-interruptCh:
				log.Print("Worker was stopped")
				return
			default:
			}

			select {
			case <-interruptCh:
				log.Print("Worker was stopped")
				return
			case <-ticker.C:
				log.Print("Do something")
			}
		}
	}()

	wg.Wait()
	log.Print("Application was stopped")
}
