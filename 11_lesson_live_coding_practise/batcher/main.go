package main

import (
	"errors"
	"sync"
	"time"
)

type Batcher struct {
	mutex    sync.Mutex
	messages []string

	size       int
	action     func([]string)
	messagesCh chan []string

	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewBatcher(action func([]string), size int) (*Batcher, error) {
	if action == nil || size <= 0 {
		return nil, errors.New("invalid arguments")
	}

	return &Batcher{
		action:      action,
		size:        size,
		messagesCh:  make(chan []string, 1),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}, nil
}

// Append add message to batch
func (b *Batcher) Append(message string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.messages = append(b.messages, message)
	if b.size == len(b.messages) {
		b.messagesCh <- b.messages
		b.messages = nil
	}
}

// Run start worker for periodic flushing
func (b *Batcher) Run(interval time.Duration) {
	go func() {
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
			case messages := <-b.messagesCh:
				b.flushMessages(messages)
				ticker.Reset(interval)
			case <-ticker.C:
				b.flush()
			}
		}
	}()
}

func (b *Batcher) flush() {
	b.mutex.Lock()
	messages := b.messages
	b.messages = nil
	b.mutex.Unlock()

	b.flushMessages(messages)
}

func (b *Batcher) flushMessages(messages []string) {
	if len(messages) != 0 {
		b.action(messages)
	}
}

// Shutdown wait worker and flush buffer before closing
func (b *Batcher) Shutdown() {
	close(b.closeCh)
	<-b.closeDoneCh
}
