package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Parallel()

	s := NewInMemoryStorage()
	s.Set(1, map[string]string{"key_0": "value_0"})
	s.Set(2, map[string]string{"key_1": "value_2"})
	s.Set(3, map[string]string{"key_1": "value_3"})
	s.Set(4, map[string]string{"key_1": "value_4"})

	assert.Equal(t, "", s.Get(1, "key_1"))
	assert.Equal(t, "value_2", s.Get(2, "key_1"))
	assert.Equal(t, "value_3", s.Get(3, "key_1"))
	assert.Equal(t, "value_4", s.Get(4, "key_1"))
	assert.Equal(t, "value_4", s.Get(5, "key_1"))

	assert.False(t, s.ExistsBetween(0, 1, map[string]string{"key_1": ""}))
	assert.True(t, s.ExistsBetween(2, 4, map[string]string{"key_1": ""}))
	assert.True(t, s.ExistsBetween(3, 5, map[string]string{"key_1": ""}))
}
