package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(ticketsNumber int) Semaphore {
	return Semaphore{
		tickets: make(chan struct{}, ticketsNumber),
	}
}

func (s *Semaphore) Acquire() {
	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tickets
}

func (s *Semaphore) WithSemaphore(action func()) {
	if action == nil {
		return
	}

	s.Acquire()
	action()
	s.Release()
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(6)

	semaphore := NewSemaphore(5)
	for i := 0; i < 6; i++ {
		semaphore.Acquire()
		go func() {
			defer func() {
				wg.Done()
				semaphore.Release()
			}()

			fmt.Println("working...")
			time.Sleep(time.Second * 2)
			fmt.Println("exiting...")
		}()
	}

	wg.Wait()
}
