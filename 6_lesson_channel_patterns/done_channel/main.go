package main

import (
	"fmt"
)

func doWork(closeCh chan struct{}) <-chan struct{} {
	closeDoneCh := make(chan struct{})

	go func() {
		defer close(closeDoneCh)

		for {
			select {
			case <-closeCh:
				return
			default:
				// ... do some work
			}
		}
	}()

	return closeDoneCh
}

func main() {
	closeCh := make(chan struct{})
	closeDoneCh := doWork(closeCh)

	close(closeCh)
	<-closeDoneCh

	fmt.Println("terminated")
}
