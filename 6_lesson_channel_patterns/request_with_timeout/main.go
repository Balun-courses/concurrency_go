package main

import (
	"errors"
	"time"
)

func request(chan string)

func requestWithTimeout(timeout time.Duration) (string, error) {
	result := make(chan string)
	go request(result)

	select {
	case data := <-result:
		return data, nil
	case <-time.After(timeout):
		return "", errors.New("timeout")
	}
}
