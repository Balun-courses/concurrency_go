package main

import (
	"fmt"
	"time"
)

func printOwn() {
	fmt.Println("goroutine")
}

func syncPrint() {
	for i := 0; i < 100; i++ {
		printOwn()
	}
}

func asyncPrint() {
	numCPU := 4 // runtime.NumCPU()
	for cpu := 0; cpu < numCPU; cpu++ {
		go func() {
			for i := 0; i < 100/numCPU; i++ {
				printOwn()
			}
		}()
	}

	time.Sleep(2 * time.Second)
}

func main() {
	asyncPrint()
}
