package main

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type Context struct {
	done   chan struct{}
	closed int32
}

func WithTimeout(parent Context, duration time.Duration) (*Context, func()) {
	if atomic.LoadInt32(&parent.closed) == 1 {
		return nil, nil // don't use nil
	}

	ctx := &Context{
		done: make(chan struct{}),
	}

	cancel := func() {
		if atomic.CompareAndSwapInt32(&ctx.closed, 0, 1) {
			close(ctx.done)
		}
	}

	go func() {
		timer := time.NewTimer(duration)
		defer timer.Stop()

		select {
		case <-parent.Done():
		case <-timer.C:
		}

		cancel()
	}()

	return ctx, cancel
}

func (c *Context) Done() <-chan struct{} {
	return c.done
}

func (c *Context) Err() error {
	select {
	case <-c.done:
		return errors.New("context deadline exceeded")
	default:
		return nil
	}
}

func (c *Context) Deadline() (time.Time, bool) {
	// not implemented
	return time.Time{}, false
}

func (c *Context) Value(any) any {
	// not implemented
	return nil
}

func main() {
	ctx, cancel := WithTimeout(Context{}, time.Second)
	defer cancel()

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Println("finished")
	case <-ctx.Done():
		fmt.Println("canceled")
	}
}
