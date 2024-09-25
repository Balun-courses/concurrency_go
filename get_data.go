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

const (
	timeout = time.Second
)

func getter(key string) (any, error) {
	data, exists := postgresStore.Load(key)
	if !exists {
		return nil, errors.New("key not found")
	}

	return data, nil
}

func GetData(key string, getter func(key string) (any, error)) (any, error) {
	ch := make(chan struct {
		val any
		err error
	})

	go func() {
		defer close(ch)
		dataPostgres, err := getter(key)
		if err != nil {
			ch <- struct {
				val any
				err error
			}{err: err}
			return
		}
		ch <- struct {
			val any
			err error
		}{val: dataPostgres}
	}()

	select {
	case dataPostgres := <-ch:
		return dataPostgres.val, nil
	case err := <-ch:
		return nil, err.err
	case <-time.After(timeout):
		dataCache, exists := cacheStore.Load(key)
		if !exists {
			return nil, errors.New("key not found")
		}
		return dataCache, nil
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