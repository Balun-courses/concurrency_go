package main

import (
	"fmt"
	"sync"
)

var x int
var y int

var local_x int
var local_y int

func main() {
	index := 0
	for {
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			x = 1
			local_x = y
		}()

		go func() {
			defer wg.Done()

			y = 1
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
