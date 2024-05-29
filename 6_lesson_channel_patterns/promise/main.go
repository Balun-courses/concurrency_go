package main

import (
	"fmt"
	"time"
)

type Promise struct {
	waitCh chan struct{}
	value  interface{}
	err    error
}

func NewPromise(task func() (interface{}, error)) *Promise {
	if task == nil {
		return nil
	}

	promise := &Promise{
		waitCh: make(chan struct{}),
	}

	go func() {
		defer close(promise.waitCh)
		promise.value, promise.err = task()
	}()

	return promise
}

func (p *Promise) Then(successCb func(interface{}), errCb func(error)) {
	<-p.waitCh
	if p.err == nil {
		successCb(p.value)
	} else {
		errCb(p.err)
	}
}

func main() {
	callback := func() (interface{}, error) {
		time.Sleep(time.Second)
		return "ok", nil
	}

	promise := NewPromise(callback)
	promise.Then(
		func(value interface{}) {
			fmt.Println("success", value)
		},
		func(err error) {
			fmt.Println("error", err.Error())
		},
	)
}
