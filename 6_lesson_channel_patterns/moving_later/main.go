package main

import (
	"fmt"
	"math/rand"
	"time"
)

type DistributedDatabase struct{}

func (d *DistributedDatabase) Query(address string, key string) string {
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	return fmt.Sprintf("[%s]: value", address)
}

var database DistributedDatabase

func Query(addresses []string, query string) string {
	result := make(chan string)
	for _, address := range addresses {
		go func(address string) {
			select {
			case result <- database.Query(address, query):
			default:
				return
			}
		}(address)
	}

	return <-result
}

func main() {
	addresses := []string{
		"127.0.0.1",
		"127.0.0.2",
		"127.0.0.3",
	}

	value := Query(addresses, "GET key_1")
	fmt.Println(value)
}
