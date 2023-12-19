package tools

import "sync"

func WithLock(mutex sync.Locker, action func()) {
	if action == nil {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	action()
}
