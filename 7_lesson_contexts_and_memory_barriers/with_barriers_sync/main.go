package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var x int
var y int

var local_x int
var local_y int

var barrier atomic.Bool

func main() {
	index := 0
	for {
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			x = 1
			barrier.Store(barrier.Load())
			local_x = y
		}()

		go func() {
			defer wg.Done()

			y = 1
			barrier.Store(barrier.Load())
			local_y = x
		}()

		wg.Wait()

		if local_x == 0 && local_y == 0 {
			fmt.Println("broken CPU, iteration =", index)
			return
		} else {
			fmt.Println("iteration =", index)
		}

		index++
		x, y = 0, 0
		local_x, local_y = 0, 0
	}
}
