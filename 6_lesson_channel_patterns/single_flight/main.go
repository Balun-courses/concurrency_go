package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type call struct {
	err   error
	value interface{}
	done  chan struct{}
}

type SingleFlight struct {
	mutex sync.Mutex
	calls map[string]*call
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight{
		calls: make(map[string]*call),
	}
}

func (s *SingleFlight) Do(key string, action func() (interface{}, error)) (interface{}, error) {
	s.mutex.Lock()
	if call, found := s.calls[key]; found {
		s.mutex.Unlock()
		return s.wait(call)
	}

	call := &call{
		done: make(chan struct{}),
	}

	s.calls[key] = call
	s.mutex.Unlock()

	go func() {
		defer func() {
			if v := recover(); v != nil {
				call.err = errors.New("fatal error from action")
			}

			s.mutex.Lock()
			close(call.done)
			delete(s.calls, key)
			s.mutex.Unlock()
		}()

		call.value, call.err = action()
	}()

	return s.wait(call)
}

func (s *SingleFlight) wait(call *call) (interface{}, error) {
	<-call.done
	return call.value, call.err
}

func main() {
	const inFlightRequests = 5
	wg := sync.WaitGroup{}
	wg.Add(inFlightRequests)

	singleFlight := NewSingleFlight()

	const key = "same_key"
	for i := 0; i < inFlightRequests; i++ {
		go func() {
			defer wg.Done()
			value, err := singleFlight.Do(key, func() (interface{}, error) {
				fmt.Println("single flight")
				time.Sleep(5 * time.Second)
				return "result", nil
			})

			fmt.Println(i, "=", value, err)
		}()
	}

	wg.Wait()

	singleFlight.Do(key, func() (interface{}, error) {
		fmt.Println("single flight")
		time.Sleep(5 * time.Second)
		return "result", nil
	})
}
