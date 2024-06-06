package main

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"
)

type CacheRCU struct {
	data unsafe.Pointer
}

func NewCacheRCU(ctx context.Context) *CacheRCU {
	data := make(map[string]string)
	cache := &CacheRCU{
		data: unsafe.Pointer(&data),
	}

	go cache.synchronize(ctx)
	return cache
}

func (c *CacheRCU) synchronize(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var data map[string]string
			// data = ... - get from remote storage

			pointer := unsafe.Pointer(&data)
			atomic.StorePointer(&c.data, pointer)
		}
	}
}

func (c *CacheRCU) Get(key string) (string, bool) {
	pointer := atomic.LoadPointer(&c.data)
	data := *(*map[string]string)(pointer)

	value, found := data[key]
	return value, found
}
