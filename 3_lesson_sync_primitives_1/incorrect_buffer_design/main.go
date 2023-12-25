package main

import "sync"

type Buffer struct {
	mtx  sync.Mutex
	data []int
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) Add(value int) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.data = append(b.data, value)
}

func (b *Buffer) Data() []int {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return b.data
}
