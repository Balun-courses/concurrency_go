package main

import (
	"context"
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
	// time.Sleep(2 * time.Second)
	data, exists := postgresStore.Load(key)
	if !exists {
		return nil, errors.New("key not found")
	}

	return data, nil
}

func GetData(key string, getter func(key string) (any, error)) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dataPostgres, err := getter(key)
	if err != nil {
		return nil, err
	}

	dataCache, exists := cacheStore.Load(key)
	if !exists {
		return nil, errors.New("key not found")
	}

	select {
	case <-ctx.Done():
		return dataCache, nil
	default:
		return dataPostgres, nil
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
