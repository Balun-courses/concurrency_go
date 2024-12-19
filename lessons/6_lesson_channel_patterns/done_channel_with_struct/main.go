package main

import (
	"fmt"
	"time"
)

type Worker struct {
	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewWorker() Worker {
	worker := Worker{
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer func() {
			ticker.Stop()
			close(worker.closeDoneCh)
		}()

		for {
			select {
			case <-worker.closeCh:
				fmt.Println("Worker was stopped")

				return
			default:
			}

			select {
			case <-worker.closeCh:
				fmt.Println("Worker was stopped")
				return
			case <-ticker.C:
				fmt.Println("Do something")
			}
		}
	}()

	return worker
}

func (w *Worker) Shutdown() {
	close(w.closeCh)
	<-w.closeDoneCh
}

func main() {
	worker := NewWorker()
	time.Sleep(5 * time.Second)
	worker.Shutdown()
}
