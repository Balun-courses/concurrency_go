package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	postgresStore sync.Map
	cacheStore    sync.Map
)

const timeout = time.Second

func getter(key string) (any, error) {
	// time.Sleep(1 * time.Microsecond)
	data, exists := postgresStore.Load(key)
	if !exists {
		return nil, errors.New("key not found")
	}
	return data, nil
}

func GetData(key string, getter func(key string) (any, error)) (any, error) {
	valCh := make(chan any)
	errCh := make(chan error)

	go func(key string) {
		data, err := getter(key)
		if err != nil {
			errCh <- err
			return
		}
		valCh <- data
	}(key)

	go func(key string) {
		data, exists := cacheStore.Load(key)
		if !exists {
			errCh <- errors.New("key not found")
			return
		}
		valCh <- data
	}(key)

	for {
		select {
		case value := <-valCh:
			return value, nil
		case err := <-errCh:
			return nil, err
		default:
		}

		select {
		case <-time.After(timeout):
			value := <-valCh
			return value, nil
		case err := <-errCh:
			return nil, err
		default:
		}
	}
}

func main() {
	postgresStore.Store("key", "postgres_value")
	cacheStore.Store("key", "cache_value")

	data, err := GetData("key", getter)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)
}