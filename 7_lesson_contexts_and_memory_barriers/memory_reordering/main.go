package main

import (
	"runtime"
	"time"
)

func main() {
	var a []int // nil
	var b bool  // false

	go func() {
		a = make([]int, 3)
		b = true
	}()

	for !b {
		time.Sleep(time.Second)
		runtime.Gosched()
	}

	a[0], a[1], a[2] = 0, 1, 2 // might panic
}
