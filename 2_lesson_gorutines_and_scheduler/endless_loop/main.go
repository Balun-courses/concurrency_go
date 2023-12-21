package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)

	var i int
	go func() {
		for {
			i++
		}
	}()

	fmt.Println(i)
}
