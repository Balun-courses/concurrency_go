package storage

import (
	"math"
	"sync/atomic"
)

type IDGenerator struct {
	counter atomic.Int64
}

func NewIDGenerator(previousID int64) *IDGenerator {
	generator := &IDGenerator{}
	generator.counter.Store(previousID)
	return generator
}

func (g *IDGenerator) Generate() int64 {
	g.counter.CompareAndSwap(math.MaxInt64, 0)
	return g.counter.Add(1)
}
