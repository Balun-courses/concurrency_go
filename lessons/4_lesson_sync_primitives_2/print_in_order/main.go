package main

import "sync/atomic"

type Foo struct {
	firstJobDone  atomic.Bool
	secondJobDone atomic.Bool
}

func (f *Foo) first(printFirst func()) {
	printFirst()
	f.firstJobDone.Store(true)
}

func (f *Foo) second(printSecond func()) {
	for !f.firstJobDone.Load() {
		// active waiting...
	}

	printSecond()
	f.secondJobDone.Store(true)
}

func (f *Foo) third(printThird func()) {
	for !f.secondJobDone.Load() {
		// active waiting...
	}

	printThird()
}
