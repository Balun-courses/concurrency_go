package main

import (
	"fmt"
	"sync"
)

type Barrier struct {
	count    int
	size     int
	mutex    sync.Mutex
	beforeCh chan int
	afterCh  chan int
}

func NewBarrier(size int) *Barrier {
	return &Barrier{
		size:     size,
		beforeCh: make(chan int, size),
		afterCh:  make(chan int, size),
	}
}

func (b *Barrier) Before() {
	b.mutex.Lock()

	b.count++
	if b.count == b.size {
		for i := 0; i < b.size; i++ {
			b.beforeCh <- 1
		}
	}

	b.mutex.Unlock()
	<-b.beforeCh
}

func (b *Barrier) After() {
	b.mutex.Lock()

	b.count--
	if b.count == 0 {
		for i := 0; i < b.size; i++ {
			b.afterCh <- 1
		}
	}

	b.mutex.Unlock()
	<-b.afterCh
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	bootstrap := func() {
		fmt.Println("bootstrap")
	}

	work := func() {
		fmt.Println("work")
	}

	barrier := NewBarrier(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				// wait for all workers to finish previous loop
				barrier.Before()
				bootstrap()
				// wait for other workers to bootstrap
				barrier.After()
				work()
			}
		}()
	}

	wg.Wait()
}
