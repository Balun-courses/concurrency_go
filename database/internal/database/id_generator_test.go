package database

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateIDGenerator(t *testing.T) {
	generator := NewIDGenerator()
	require.Equal(t, int64(0), generator.counter.Load())
}

func TestGenerateID(t *testing.T) {
	generator := NewIDGenerator()

	id := generator.Generate()
	require.Equal(t, int64(1), id)

	id = generator.Generate()
	require.Equal(t, int64(2), id)
}
