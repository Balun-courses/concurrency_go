package main

import "sync"

type RWMutex struct {
	notifier      *sync.Cond
	mutex         *sync.Mutex
	readersNumber int
	hasWriter     bool
}

func NewRWMutex() *RWMutex {
	var mutex sync.Mutex
	notifier := sync.NewCond(&mutex)

	return &RWMutex{
		mutex:    &mutex,
		notifier: notifier,
	}
}

func (m *RWMutex) Lock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for m.hasWriter {
		m.notifier.Wait()
	}

	m.hasWriter = true
	for m.readersNumber != 0 {
		m.notifier.Wait()
	}
}

func (m *RWMutex) Unlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.hasWriter = false
	m.notifier.Broadcast()
}

func (m *RWMutex) RLock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for m.hasWriter {
		m.notifier.Wait()
	}

	m.readersNumber++
}

func (m *RWMutex) RUnlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.readersNumber--
	if m.readersNumber == 0 {
		m.notifier.Broadcast()
	}
}
