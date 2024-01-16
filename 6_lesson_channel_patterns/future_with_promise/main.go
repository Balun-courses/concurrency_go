package main

import (
	"fmt"
	"time"
)

type Future struct {
	result <-chan interface{}
}

func NewFuture(result <-chan interface{}) Future {
	return Future{
		result: result,
	}
}

func (f *Future) Get() interface{} {
	return <-f.result
}

type Promise struct {
	result   chan interface{}
	promised bool
}

func NewPromise() Promise {
	return Promise{
		result: make(chan interface{}, 1),
	}
}

// Set don't use with any goroutines
func (p *Promise) Set(value interface{}) {
	if p.promised {
		return
	}

	p.promised = true
	p.result <- value
	close(p.result)
}

func (p *Promise) GetFuture() Future {
	return NewFuture(p.result)
}

func main() {
	promise := NewPromise()
	go func() {
		time.Sleep(time.Second)
		promise.Set("Test")
	}()

	future := promise.GetFuture()
	value := future.Get()
	fmt.Println(value)
}
