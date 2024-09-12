package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type IRedisDb interface {
	Add(key, value string)
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, keys []string) ([]*string, error)
	Keys(ctx context.Context) ([]string, error)
}

// Questions?????
// - restrictions RAM
// - sync no more, than 1 minute
// - goroutine safe

type RedisCache struct {
	data sync.Map
	db   IRedisDb
}

func NewRedisCache(ctx context.Context, db IRedisDb) (*RedisCache, error) {
	cache := &RedisCache{db: db}
	go cache.synchronize(ctx)
	return cache, nil
}

func (r *RedisCache) synchronize(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.syncImplement(ctx)
		}
	}
}

func (r *RedisCache) syncImplement(ctx context.Context) {
	var data []*string
	var keys []string
	var err error

	if err := withRetry(ctx, 3, time.Millisecond*100, func() error {
		keys, err = r.db.Keys(ctx)
		return err
	}); err != nil || len(keys) == 0 {
		return
	}

	err = withRetry(ctx, 3, time.Millisecond*100, func() error {
		data, err = r.db.MGet(ctx, keys)
		return err
	})

	for index, key := range keys {
		if index < len(data) && data[index] != nil {
			r.data.Store(key, data[index])
		}
	}
}

func (r *RedisCache) Add(key, value string) {
	r.data.Store(key, value)
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if data, exists := r.data.Load(key); exists {
		return data.(string), nil
	} else {
		return "", errors.New("key not found")
	}
}

func (r *RedisCache) MGet(ctx context.Context, keys []string) ([]*string, error) {
	var data []*string

	for _, key := range keys {
		if value, exists := r.data.Load(key); exists {
			value := value.(string)
			data = append(data, &value)
		}
	}

	return data, nil
}

func (r *RedisCache) Keys(ctx context.Context) ([]string, error) {
	var keys []string

	r.data.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})

	return keys, nil
}

func withRetry(ctx context.Context, retries int, delay time.Duration, action func() error) error {
	var err error

	if action == nil || retries <= 0 || delay < 0 {
		return errors.New("invalid arguments")
	}

	for i := 0; i < retries; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err = action(); err == nil {
			return nil
		}

		time.Sleep(time.Duration(i+1) * delay)
	}

	return err
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redis := &RedisCache{}
	cache, err := NewRedisCache(ctx, redis.db)

	if err != nil {
		fmt.Println("Error creating RedisCache:", err)
		return
	}

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")
	cache.Add("key4", "value4")

	if data, err := cache.Get(ctx, "key1"); err == nil {
		fmt.Println("Get data from cache:", data)
	} else {
		fmt.Println("Error:", err)
	}

	keys := []string{"key1", "key2", "key3", "key4"}
	if data, err := cache.MGet(ctx, keys); err == nil {
		for key, value := range data {
			if value != nil {
				fmt.Printf("MGet %s: %s\n", keys[key], *value)
			} else {
				fmt.Printf("MGet %s: not found\n", keys[key])
			}
		}
	} else {
		fmt.Println("Error:", err)
	}

	if keys, err := cache.Keys(ctx); err == nil {
		fmt.Println("Keys in cache:", keys)
	} else {
		fmt.Println("Error:", err)
	}
}
