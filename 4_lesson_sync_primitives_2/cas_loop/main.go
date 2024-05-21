package main

import (
	"sync/atomic"
)

func IncrementAndGet(pointer *int32) int32 {
	for {
		currentValue := atomic.LoadInt32(pointer)
		nextValue := currentValue + 1
		if atomic.CompareAndSwapInt32(pointer, currentValue, nextValue) {
			return nextValue
		}
	}
}
