package main

import (
	"log"
	"runtime"
)

var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
	if done {
		log.Println(len(a)) // always 12 once printed
	}
}

func main() {
	go setup()

	for !done {
		runtime.Gosched()
	}

	log.Println(a) // expected to print: hello, world
}
