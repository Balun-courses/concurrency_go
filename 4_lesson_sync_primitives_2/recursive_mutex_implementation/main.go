package main

import (
	"bytes"
	"errors"
	"runtime"
	"strconv"
	"sync"
)

// This is terrible, slow, and should never be used.
func goid() (int, error) {
	buf := make([]byte, 32)
	n := runtime.Stack(buf, false)
	buf = buf[:n]
	// goroutine 1 [running]: ...

	buf, ok := bytes.CutPrefix(buf, []byte("goroutine "))
	if !ok {
		return 0, errors.New("bad stack")
	}

	i := bytes.IndexByte(buf, ' ')
	if i < 0 {
		return 0, errors.New("bad stack")
	}

	return strconv.Atoi(string(buf[:i]))
}

type RecursiveMutex struct {
	mutex    sync.Mutex
	notifier sync.Cond
	count    int
	owner    int
	locked   bool
}

func NewRecursiveMutex() *RecursiveMutex {
	return &RecursiveMutex{}
}

func (m *RecursiveMutex) Lock() {
	id, err := goid()
	if err != nil {
		panic("recursive_mutex: " + err.Error())
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for m.locked && id != m.owner {
		m.notifier.Wait()
	}

	m.count++
	m.locked = true
}

func (m *RecursiveMutex) Unlock() {
	id, err := goid()
	if err != nil {
		panic("recursive_mutex: " + err.Error())
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.locked || id != m.owner {
		panic("incorrect usage")
	}

	m.count--
	if m.count == 0 {
		m.locked = false
		m.notifier.Broadcast()
	}
}
