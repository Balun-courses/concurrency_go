package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Ticker struct {
	C        chan struct{}
	interval int64
	closed   atomic.Bool
}

func NewTicker(interval time.Duration) *Ticker {
	ticker := &Ticker{
		C:        make(chan struct{}),
		interval: int64(interval),
	}

	go func() {
		for !ticker.closed.Load() {
			duration := atomic.LoadInt64(&ticker.interval)
			time.Sleep(time.Duration(duration))
			ticker.C <- struct{}{}
		}
	}()

	return ticker
}

func (t *Ticker) Stop() {
	t.closed.Store(true)
}

func (t *Ticker) Reset(interval time.Duration) {
	atomic.StoreInt64(&t.interval, int64(interval))
}

func main() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		fmt.Println("tick")
	}
}
