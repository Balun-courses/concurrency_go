package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	fmt.Println("CPU:", runtime.NumCPU())

	runtime.GOMAXPROCS(16)

	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	fmt.Println("CPU:", runtime.NumCPU())
}
