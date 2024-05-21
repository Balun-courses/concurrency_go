package main

import "sync"

func withLock(mutex sync.Locker, action func()) {
	if action == nil {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	action()
}

func main() {
	mutex := sync.RWMutex{}
	withLock(&mutex, func() {
		// write lock
	})

	withLock(mutex.RLocker(), func() {
		// read lock
	})
}
