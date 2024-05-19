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

// Questions?????
// - restrictions RAM
// - sync no more, than 1 minute
// - goroutine safe

type RedisDatabaseWithCache struct {
	mutex sync.RWMutex
	data  map[string]string

	operationsMutex  sync.Mutex
	activeOperations map[string]chan struct{}

	database RedisDatabase
}

func NewRedisDatabaseWithCache(ctx context.Context, database RedisDatabase) (*RedisDatabaseWithCache, error) {
	if database == nil {
		return nil, errors.New("invalid arguments")
	}

	db := &RedisDatabaseWithCache{
		activeOperations: make(map[string]chan struct{}),
		data:             make(map[string]string),
		database:         database,
	}

	go db.synchronize(ctx)
	return db, nil
}

func (c *RedisDatabaseWithCache) synchronize(ctx context.Context) {
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
			c.syncImpl(ctx)
		}
	}
}

func (c *RedisDatabaseWithCache) syncImpl(ctx context.Context) {
	var err error
	var keys []string
	err = withRetry(ctx, 3, time.Millisecond*100, func() error {
		keys, err = c.database.Keys(ctx)
		return err
	})

	if err != nil || len(keys) == 0 {
		return
	}

	var values []*string
	err = withRetry(ctx, 3, time.Millisecond*100, func() error {
		values, err = c.database.MGet(ctx, keys)
		return err
	})

	if err != nil {
		return
	}

	data := make(map[string]string, len(keys))
	for idx, key := range keys {
		value := values[idx]
		if value != nil {
			data[key] = *value
		}
	}

	c.mutex.Lock()
	c.data = data
	c.mutex.Unlock()
}

// Get need to proxy only one method from interface
func (c *RedisDatabaseWithCache) Get(ctx context.Context, key string) (string, error) {
	var found bool
	var value string
	withLock(c.mutex.RLocker(), func() {
		value, found = c.data[key]
	})

	if found {
		return value, nil
	}

	var barrier chan struct{}
	var isActiveOperation bool
	withLock(&c.operationsMutex, func() {
		if barrier, isActiveOperation = c.activeOperations[key]; !isActiveOperation {
			barrier = make(chan struct{})
			c.activeOperations[key] = barrier
		}
	})

	if isActiveOperation {
		<-barrier
		withLock(c.mutex.RLocker(), func() {
			value, found = c.data[key]
		})

		return value, nil
	}

	defer func() {
		close(barrier)
		withLock(&c.operationsMutex, func() {
			delete(c.activeOperations, key)
		})
	}()

	var err error
	err = withRetry(ctx, 3, time.Millisecond*100, func() error {
		value, err = c.database.Get(ctx, key)
		return err
	})

	if err != nil {
		return "", err
	}

	withLock(&c.mutex, func() {
		c.data[key] = value
	})

	return value, nil
}

func withLock(locker sync.Locker, action func()) {
	if action == nil {
		return
	}

	locker.Lock()
	action()
	locker.Unlock()
}

func withRetry(ctx context.Context, retriesNumber int, delay time.Duration, action func() error) error {
	if action == nil || retriesNumber <= 0 || delay < 0 {
		return errors.New("invalid arguments")
	}

	var err error
	for retry := 1; retry <= retriesNumber; retry++ {
		if ctx.Err() != nil {
			return err
		}

		if err = action(); err == nil {
			return nil
		}

		time.Sleep(time.Duration(retry) * delay)
	}

	return err
}
