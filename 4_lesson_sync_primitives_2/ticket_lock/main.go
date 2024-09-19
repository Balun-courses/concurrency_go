package main

import (
	"runtime"
	"sync/atomic"
)

type TicketLock struct {
	ownerTicket    atomic.Int64
	nextFreeTicket atomic.Int64
}

func NewTicketLock() *TicketLock {
	return &TicketLock{}
}

func (t *TicketLock) Lock() {
	ticket := t.nextFreeTicket.Add(1)
	for t.ownerTicket.Load() != ticket-1 {
		runtime.Gosched()
	}
}

func (t *TicketLock) Unlock() {
	t.ownerTicket.Add(1)
}
