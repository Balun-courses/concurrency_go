package main

import (
	"errors"
	"sync"
	"time"
)

type Batcher struct {
	maxSize       int
	flushInterval time.Duration
	flushAction   func([]string)
	ticker        *time.Ticker

	mutex    sync.Mutex
	messages []string

	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewBatcher(action func([]string), size int, interval time.Duration) (*Batcher, error) {
	if action == nil {
		return nil, errors.New("invalid action")
	}

	if size <= 0 {
		return nil, errors.New("invalid size")
	}

	if interval <= 0 {
		return nil, errors.New("invalid interval")
	}

	return &Batcher{
		maxSize:       size,
		flushAction:   action,
		flushInterval: interval,
		closeCh:       make(chan struct{}),
		closeDoneCh:   make(chan struct{}),
	}, nil
}

// Append add message to batch
func (b *Batcher) Append(message string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	select {
	case <-b.closeCh:
		return errors.New("batcher is close")
	default:
	}

	b.messages = append(b.messages, message)
	if len(b.messages) == b.maxSize {
		b.flushLocked()
		b.ticker.Reset(b.flushInterval)
	}

	return nil
}

// Run start worker for periodic flushing
func (b *Batcher) Run() {
	if b.ticker != nil {
		return
	}

	b.ticker = time.NewTicker(b.flushInterval)

	go func() {
		defer close(b.closeDoneCh)

		for {
			select {
			case <-b.closeCh:
				b.flush()
				return
			default:
			}

			select {
			case <-b.closeCh:
				b.flush()
				return
			case <-b.ticker.C:
				b.flush()
			}
		}
	}()
}

func (b *Batcher) flush() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.flushLocked()
}

func (b *Batcher) flushLocked() {
	if len(b.messages) == 0 {
		return
	}

	messages := b.messages
	b.messages = nil
	go b.flushAction(messages)
}

// Close wait worker and flush buffer before closing
func (b *Batcher) Close() {
	select {
	case <-b.closeCh:
		return
	default:
	}

	b.mutex.Lock()
	close(b.closeCh)
	b.mutex.Unlock()

	<-b.closeDoneCh
	b.ticker.Stop()
}
