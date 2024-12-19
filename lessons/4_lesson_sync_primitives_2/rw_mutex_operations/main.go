package main

import "sync"

func RUnlockLockedMutex() {
	m := sync.RWMutex{}
	m.Lock()
	m.RUnlock()
}

func UnlockRLockedMutex() {
	m := sync.RWMutex{}
	m.RLock()
	m.Unlock()
}

func LockRLockedMutex() {
	m := sync.RWMutex{}
	m.Lock()
	m.RLock()
}

func main() {
}
