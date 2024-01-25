package main

import (
	"sync"
	"time"
)

type Batcher struct {
	f    func([]string)
	m    sync.Mutex
	msgs []string

	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewBatcher(f func([]string)) *Batcher {
	b := &Batcher{
		f:           f,
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}

	go b.run()
	return b
}

func (b *Batcher) Append(s string) {
	b.m.Lock()
	defer b.m.Unlock()
	b.msgs = append(b.msgs, s)
}

func (b *Batcher) run() {
	for {
		select {
		case <-b.closeCh:
			b.closeDoneCh <- struct{}{}
			return
		case <-time.After(time.Second):
			b.call()
		}
	}
}

func (b *Batcher) call() {
	b.m.Lock()
	msgs := b.msgs
	b.msgs = nil
	b.m.Unlock()

	if len(msgs) > 0 {
		b.f(msgs)
	}
}

func (b *Batcher) Cancel() {
	select {
	case b.closeCh <- struct{}{}:
		<-b.closeDoneCh
		b.call()
	default:
	}
}
