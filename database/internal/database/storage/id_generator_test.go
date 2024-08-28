package storage

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDGeneratorCreatingWithoutPreviousID(t *testing.T) {
	t.Parallel()

	generator := NewIDGenerator(0)
	assert.Equal(t, int64(0), generator.counter.Load())
}

func TestIDGeneratorCreatingWithPreviousID(t *testing.T) {
	t.Parallel()

	generator := NewIDGenerator(1000)
	assert.Equal(t, int64(1000), generator.counter.Load())
}

func TestGenerateID(t *testing.T) {
	t.Parallel()

	generator := NewIDGenerator(0)

	goroutinesNumber := 999
	wg := sync.WaitGroup{}
	wg.Add(goroutinesNumber)

	for i := 0; i < goroutinesNumber; i++ {
		go func() {
			defer wg.Done()
			_ = generator.Generate()
		}()
	}

	wg.Wait()

	nextID := generator.Generate()
	expectedID := goroutinesNumber + 1
	assert.Equal(t, int64(expectedID), nextID)
}

func TestGenerateIDOverflow(t *testing.T) {
	t.Parallel()

	generator := NewIDGenerator(math.MaxInt64)

	nextID := generator.Generate()
	assert.Equal(t, int64(1), nextID)
}
