package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}
