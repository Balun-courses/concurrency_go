package main

import (
	"fmt"
	"sync"
)

// Need to show solution

func WaitToClose(lhs, rhs chan struct{}) {
	lhsClosed, rhsClosed := false, false
	for !lhsClosed || !rhsClosed {
		select {
		case _, ok := <-lhs:
			fmt.Println("lhs", ok)
			if !ok {
				lhsClosed = true
			}
		case _, ok := <-rhs:
			fmt.Println("rhs", ok)
			if !ok {
				rhsClosed = true
			}
		}
	}
}

func main() {
	lhs := make(chan struct{}, 1)
	rhs := make(chan struct{}, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		WaitToClose(lhs, rhs)
	}()

	lhs <- struct{}{}
	rhs <- struct{}{}

	close(lhs)
	close(rhs)

	wg.Wait()
}
