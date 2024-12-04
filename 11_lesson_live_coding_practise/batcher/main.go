package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Batcher struct {
	items   []string
	size    int
	timeout time.Duration
	mutex   sync.Mutex
	cond    *sync.Cond
	counter int
}

func NewBatcher(size int, timeout time.Duration) (*Batcher, error) {
	if size <= 0 {
		return nil, errors.New("invalid argument")
	}

	bt := &Batcher{
		items:   make([]string, 0, size),
		size:    size,
		timeout: timeout,
	}

	bt.cond = sync.NewCond(&bt.mutex)
	go bt.runBatcher()
	return bt, nil
}

func (b *Batcher) runBatcher() {
	ticker := time.NewTicker(b.timeout)
	defer ticker.Stop()

	for {
		b.mutex.Lock()
		for len(b.items) < b.size {
			b.cond.Wait()
		}

		batch := b.items[:b.size]
		b.items = b.items[b.size:]

		<-ticker.C
		b.counter++
		fmt.Printf("Batch %d: %s\n", b.counter, batch)
		b.mutex.Unlock()
	}
}

func (b *Batcher) append(item string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.items = append(b.items, item)
	if len(b.items) >= b.size {
		b.cond.Signal()
	}
}

func (b *Batcher) flush() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.items) > 0 {
		batch := b.items
		b.items = nil
		b.counter++
		fmt.Printf("Batch %d: %s\n", b.counter, batch)
	}
}

func main() {
	batcher, err := NewBatcher(4, 2*time.Second)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	for i := 1; i <= 10; i++ {
		batcher.append(fmt.Sprintf("Item %d", i))
		time.Sleep(100 * time.Millisecond)
	}

	batcher.flush()
}
