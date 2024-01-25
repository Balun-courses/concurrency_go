package main

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"
)

type RedisDatabase interface {
	GetValue(string) (string, error)
	GetAllKeys() ([]string, error)
}

type RedisDatabaseWithCache struct {
	cache    *map[string]string
	database RedisDatabase
}

func NewRedisDatabaseWithCache(ctx context.Context) *RedisDatabaseWithCache {
	database := &RedisDatabaseWithCache{}

	go func() {
		for {
			var err error
			var keys []string
			for retry := 1; retry <= 3; retry++ {
				if keys, err = database.GetAllKeys(); err != nil {
					time.Sleep(time.Millisecond * 50 * time.Duration(retry))
				} else {
					break
				}
			}

			if err != nil {
				continue
			}

			cache := make(map[string]string)
			for _, key := range keys {
				value, _ := database.GetValue(key)
				cache[key] = value
			}

			new := unsafe.Pointer(&cache)
			old := unsafe.Pointer(database.cache)
			atomic.SwapPointer(&old, new)
		}
	}()

	return database
}

func (c *RedisDatabaseWithCache) GetValue(key string) (string, error) {
	value, found := (*c.cache)[key]
	if found {
		return value, nil
	}

	value, _ = c.database.GetValue(key)
	(*c.cache)[key] = value
	return value, nil
}

func (c *RedisDatabaseWithCache) GetAllKeys() ([]string, error) {

}
