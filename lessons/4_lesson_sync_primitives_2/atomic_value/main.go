package main

import (
	"sync/atomic"
	"time"
)

func loadConfig() map[string]string {
	return make(map[string]string)
}

func requests() chan int {
	return make(chan int)
}

func main() {
	var config atomic.Value
	config.Store(loadConfig())

	go func() {
		for {
			time.Sleep(10 * time.Second)
			config.Store(loadConfig())
		}
	}()

	for r := range requests() {
		c := config.Load().(map[string]string)
		_, _ = r, c
	}
}
