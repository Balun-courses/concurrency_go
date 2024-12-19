package main

import (
	"sync"
	"time"
)

type TimedMutex struct {
	mutex sync.Mutex
}

func NewTimedMutex() *TimedMutex {
	return &TimedMutex{}
}

func (m *TimedMutex) Lock() {
	m.mutex.Lock()
}

func (m *TimedMutex) TryLock() bool {
	return m.mutex.TryLock()
}

func (m *TimedMutex) TryLockFor(duration time.Duration) bool {
	const periodsNumber = 10
	period := duration / periodsNumber
	for duration > 0 {
		if m.TryLock() {
			return true
		}

		time.Sleep(period)
		duration -= period
	}

	return false
}

func (m *TimedMutex) Unlock() {
	m.mutex.Unlock()
}
