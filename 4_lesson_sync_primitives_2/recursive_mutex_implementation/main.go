package main

import (
	"bytes"
	"errors"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
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
	mutex sync.Mutex
	count atomic.Int32
	owner atomic.Int32
}

func NewRecursiveMutex() *RecursiveMutex {
	return &RecursiveMutex{}
}

func (m *RecursiveMutex) Lock() {
	id, err := goid()
	if err != nil {
		panic("recursive_mutex: " + err.Error())
	}

	if m.owner.Load() == int32(id) {
		m.count.Add(1)
	} else {
		m.mutex.Lock()
		m.owner.Store(int32(id))
		m.count.Store(1)
	}
}

func (m *RecursiveMutex) Unlock() {
	id, err := goid()
	if err != nil || m.owner.Load() != int32(id) {
		panic("recursive_mutex: " + err.Error())
	}

	m.count.Add(-1)
	if m.count.Load() == 0 {
		m.owner.Store(0)
		m.mutex.Unlock()
	}
}
