package main

import (
	"log"
	"runtime"
	"sync/atomic"
)

var a string
var done atomic.Bool

func setup() {
	a = "hello, world"
	done.Store(true)
	if done.Load() {
		log.Println(len(a)) // always 12 once printed
	}
}

func main() {
	go setup()

	for !done.Load() {
		runtime.Gosched()
	}

	log.Println(a) // hello, world
}
