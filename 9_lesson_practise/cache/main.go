package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

type RedisDatabase interface {
	Get(context.Context, string) (string, error)
	MGet(context.Context, []string) ([]*string, error)
	Keys(context.Context) ([]string, error)
}

type RedisDatabaseWithCache struct {
	mutex sync.RWMutex
	cache map[string]string

	mutexOp  sync.Mutex
	activeOp map[string]chan struct{}

	database RedisDatabase
}

func NewRedisDatabaseWithCache(ctx context.Context, database RedisDatabase) *RedisDatabaseWithCache {
	db := &RedisDatabaseWithCache{
		cache:    make(map[string]string),
		database: database,
	}

	go func() {
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
				db.synchronize(ctx)
			}
		}
	}()

	return db
}

func (c *RedisDatabaseWithCache) synchronize(ctx context.Context) {
	var keys []string
	err := withRetry(ctx, 3, 50, func() error {
		var err error
		keys, err = c.database.Keys(ctx)
		return err
	})

	if err != nil {
		return
	}

	cache := make(map[string]string, len(keys))
	values, err := c.database.MGet(ctx, keys)
	if err != nil {
		return
	}

	for idx, value := range values {
		if value != nil {
			key := keys[idx]
			cache[key] = cache[*value]
		}
	}

	c.mutex.Lock()
	c.cache = cache
	c.mutex.Unlock()
}

func (c *RedisDatabaseWithCache) GetValue(ctx context.Context, key string) (string, error) {
	var found bool
	var value string
	withLock(c.mutex.RLocker(), func() {
		value, found = c.cache[key]
	})

	if found {
		return value, nil
	}

	var isActive bool
	var barrier chan struct{}
	withLock(&c.mutexOp, func() {
		if barrier, isActive = c.activeOp[key]; !isActive {
			barrier = make(chan struct{})
			c.activeOp[key] = barrier
		}
	})

	if isActive {
		<-barrier
		withLock(c.mutex.RLocker(), func() {
			value, found = c.cache[key]
		})

		return value, nil
	}

	defer func() {
		withLock(&c.mutexOp, func() {
			close(barrier)
			delete(c.activeOp, key)
		})
	}()

	err := withRetry(ctx, 3, 50, func() error {
		var err error
		value, err = c.database.Get(ctx, key)
		return err
	})

	if err != nil {
		return "", err
	}

	withLock(&c.mutex, func() {
		c.cache[key] = value
	})

	return value, nil
}

func (c *RedisDatabaseWithCache) Keys(ctx context.Context) ([]string, error) {
	var keys []string
	return keys, withRetry(ctx, 3, 50, func() error {
		var err error
		keys, err = c.database.Keys(ctx)
		return err
	})
}

func withRetry(ctx context.Context, retriesNumber int, initialDelayMs time.Duration, action func() error) error {
	if action == nil {
		return errors.New("incorrect action")
	}

	var err error
	for retry := 1; retry <= retriesNumber; retry++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err = action(); err == nil {
			break
		}

		time.Sleep(time.Millisecond * initialDelayMs * time.Duration(retry))
	}

	return err
}

func withLock(locker sync.Locker, action func()) {
	if action == nil {
		return
	}

	locker.Lock()
	action()
	locker.Unlock()
}
