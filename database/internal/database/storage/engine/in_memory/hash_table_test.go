package in_memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHashTable(t *testing.T) {
	t.Parallel()

	table := NewHashTable()
	require.NotNil(t, table)
	assert.NotNil(t, table.data)
}

func TestHashTableSet(t *testing.T) {
	t.Parallel()

	table := &HashTable{
		data: map[string]string{
			"key_1": "value_1",
			"key_2": "value_2",
		},
	}

	tests := map[string]struct {
		key           string
		value         string
		expectedValue string
	}{
		"set not existing key": {
			key:           "key_3",
			value:         "new_value",
			expectedValue: "new_value",
		},
		"set existing key": {
			key:           "key_1",
			value:         "new_value",
			expectedValue: "new_value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			table.Set(test.key, test.value)
			value, found := table.Get(test.key)
			assert.Equal(t, test.expectedValue, value)
			assert.True(t, found)
		})
	}
}

func TestHashTableGet(t *testing.T) {
	t.Parallel()

	table := &HashTable{
		data: map[string]string{
			"key_1": "value_1",
			"key_2": "value_2",
		},
	}

	tests := map[string]struct {
		key           string
		expectedValue string
		found         bool
	}{
		"get not existing key": {
			key:           "key_3",
			found:         false,
			expectedValue: "",
		},
		"get existing key": {
			key:           "key_1",
			found:         true,
			expectedValue: "value_1",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value, found := table.Get(test.key)
			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.found, found)
		})
	}
}

func TestHashTableDel(t *testing.T) {
	t.Parallel()

	table := &HashTable{
		data: map[string]string{
			"key_1": "value_1",
			"key_2": "value_2",
		},
	}

	tests := map[string]struct {
		key string
	}{
		"del not existing key": {
			key: "key_3",
		},
		"del existing key": {
			key: "key_1",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			table.Del(test.key)
			value, found := table.Get(test.key)
			assert.Equal(t, "", value)
			assert.False(t, found)
		})
	}
}
