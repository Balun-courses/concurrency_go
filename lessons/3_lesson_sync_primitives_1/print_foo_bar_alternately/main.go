package main

import "sync"

type FooBar struct {
	number   int
	fooMutex sync.Mutex
	barMutex sync.Mutex
}

func NewFooBar(number int) *FooBar {
	fb := &FooBar{number: number}
	fb.barMutex.Lock()
	return fb
}

func (fb *FooBar) Foo(printFoo func()) {
	for i := 0; i < fb.number; i++ {
		fb.fooMutex.Lock()
		printFoo()
		fb.barMutex.Unlock()
	}
}

func (fb *FooBar) Bar(printBar func()) {
	for i := 0; i < fb.number; i++ {
		fb.barMutex.Lock()
		printBar()
		fb.fooMutex.Unlock()
	}
}
