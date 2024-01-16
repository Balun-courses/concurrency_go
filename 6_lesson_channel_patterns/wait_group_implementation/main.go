package main

import "fmt"

type WaitGroup struct {
	ch   chan struct{}
	size int
}

func NewWaitGroup(size int) WaitGroup {
	return WaitGroup{
		ch:   make(chan struct{}, size),
		size: size,
	}
}

func (wg *WaitGroup) Done() {
	wg.ch <- struct{}{}
}

func (wg *WaitGroup) Wait() {
	for i := 0; i < wg.size; i++ {
		<-wg.ch
	}
}

func main() {
	size := 5
	wg := NewWaitGroup(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			fmt.Println("message")
		}()
	}

	wg.Wait()
}
