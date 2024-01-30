package main

import (
	"sync"
	"time"
)

type Batcher struct {
	size     int
	action   func([]string)
	mutex    sync.Mutex
	messages []string

	batchesCh   chan []string
	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewBatcher(action func([]string), size int) *Batcher {
	return &Batcher{
		size:        size,
		action:      action,
		batchesCh:   make(chan []string, 1),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}
}

func (b *Batcher) Append(message string) {
	select {
	case <-b.closeCh:
		return
	default:
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.messages = append(b.messages, message)
	if len(b.messages) == b.size {
		b.batchesCh <- b.messages
		b.messages = nil
	}
}

func (b *Batcher) Run(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
		close(b.closeDoneCh)
	}()

	for {
		select {
		case <-b.closeCh:
			b.flush()
		default:
		}

		select {
		case <-b.closeCh:
			b.flush()
			return
		case messages := <-b.batchesCh:
			b.action(messages)
			ticker.Reset(interval)
		case <-ticker.C:
			b.flush()
		}
	}
}

func (b *Batcher) flush() {
	b.mutex.Lock()
	messages := b.messages
	b.messages = nil
	b.mutex.Unlock()

	if len(messages) > 0 {
		b.action(messages)
	}
}

func (b *Batcher) Shutdown() {
	close(b.closeCh)
	<-b.closeDoneCh
}
