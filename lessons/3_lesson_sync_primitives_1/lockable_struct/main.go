package main

import "sync"

type Lockable[T any] struct {
	sync.Mutex
	Data T
}

func main() {
	var l1 Lockable[int32]
	l1.Lock()
	l1.Data = 100
	l1.Unlock()

	var l2 Lockable[string]
	l2.Lock()
	l2.Data = "test"
	l2.Unlock()
}
