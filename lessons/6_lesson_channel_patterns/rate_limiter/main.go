package main

import (
	"time"
)

type LeakyBucketLimiter struct {
	leakyBucketCh chan struct{}
}

func NewLeakyBucketLimiter(limit int, period time.Duration) *LeakyBucketLimiter {
	limiter := &LeakyBucketLimiter{
		leakyBucketCh: make(chan struct{}, limit),
	}

	leakInterval := period.Nanoseconds() / int64(limit)
	go limiter.startPeriodicLeak(time.Duration(leakInterval))
	return limiter
}

func (l *LeakyBucketLimiter) startPeriodicLeak(interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			select {
			case <-l.leakyBucketCh:
			default:
			}
		}
	}
}

func (l *LeakyBucketLimiter) Allow() bool {
	select {
	case l.leakyBucketCh <- struct{}{}:
		return true
	default:
		return false
	}
}
