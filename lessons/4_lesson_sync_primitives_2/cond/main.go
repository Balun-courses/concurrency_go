package main

import (
	"log"
	"sync"
	"time"
)

func subscribe(name string, data map[string]string, c *sync.Cond) {
	c.L.Lock()

	for len(data) == 0 {
		c.Wait()
	}

	log.Printf("[%s] %s\n", name, data["key"])

	c.L.Unlock()
}

func publish(name string, data map[string]string, c *sync.Cond) {
	time.Sleep(time.Second)

	c.L.Lock()
	data["key"] = "value"
	c.L.Unlock()

	log.Printf("[%s] data publisher\n", name)
	c.Broadcast()
}

func main() {
	data := map[string]string{}
	cond := sync.NewCond(&sync.Mutex{})

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		subscribe("subscriber_1", data, cond)
	}()

	go func() {
		defer wg.Done()
		subscribe("subscriber_2", data, cond)
	}()

	go func() {
		defer wg.Done()
		publish("publisher", data, cond)
	}()

	wg.Wait()
}
