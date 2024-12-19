package main

import (
	"sync"
	"time"
)

func withLock(mutex sync.Locker, action func()) {
	if mutex == nil || action == nil {
		return
	}

	mutex.Lock()
	action()
	mutex.Unlock()
}

func main() {
	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	wg := sync.WaitGroup{}
	wg.Add(1)

	parallelTx2 := func() {
		defer wg.Done()
		tx2 := s.StartTransaction()
		tx2.Set("key_2", "value_2")
		time.Sleep(time.Millisecond * 200)
		tx2.Set("key_1", "value_1")
		tx2.Commit()
	}

	tx1 := s.StartTransaction()
	tx1.Get("key_1")
	go parallelTx2()
	time.Sleep(time.Millisecond * 100)
	tx1.Get("key_2")
	tx1.Commit()

	wg.Wait()
}
