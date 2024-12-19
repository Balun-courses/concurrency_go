package main

import (
	"fmt"
	"math/rand"
	"time"
)

type DistributedDatabase struct{}

func (d *DistributedDatabase) Query(address string, key string) string {
	time.Sleep(time.Second * time.Duration(rand.Intn(3)))
	return fmt.Sprintf("[%s]: value", address)
}

var database DistributedDatabase

func DistributedQuery(addresses []string, query string) string {
	responseCh := make(chan string, 1) // buffered necessary
	for _, address := range addresses {
		go func(address string) {
			select {
			case responseCh <- database.Query(address, query):
			default:
				return
			}
		}(address)
	}

	return <-responseCh
}

func main() {
	addresses := []string{
		"127.0.0.1",
		"127.0.0.2",
		"127.0.0.3",
	}

	value := DistributedQuery(addresses, "GET key_1")
	fmt.Println(value)
}
