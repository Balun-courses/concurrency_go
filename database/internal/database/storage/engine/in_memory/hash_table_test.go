package in_memory

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateStorage(t *testing.T) {
	t.Parallel()

	table := NewHashTable()
	require.NotNil(t, table.data)
}

func TestSet(t *testing.T) {
	t.Parallel()

	table := &HashTable{
		data: map[string]string{
			"key_1": "value_1",
			"key_2": "value_2",
		},
	}

	t.Run("test set not existing key", func(t *testing.T) {
		table.Set("key_3", "new_value")

		value, found := table.data["key_3"]
		require.Equal(t, "new_value", value)
		require.True(t, found)
	})

	t.Run("test set existing key", func(t *testing.T) {
		table.Set("key_1", "new_value")

		value, found := table.data["key_1"]
		require.Equal(t, "new_value", value)
		require.True(t, found)
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	table := &HashTable{
		data: map[string]string{
			"key_1": "value_1",
			"key_2": "value_2",
		},
	}

	t.Run("test get not existing key", func(t *testing.T) {
		value, found := table.Get("key_3")
		require.Equal(t, "", value)
		require.False(t, found)
	})

	t.Run("test get existing key", func(t *testing.T) {
		value, found := table.Get("key_1")
		require.Equal(t, "value_1", value)
		require.True(t, found)
	})
}

func TestDel(t *testing.T) {
	t.Parallel()

	table := &HashTable{
		data: map[string]string{
			"key_1": "value_1",
			"key_2": "value_2",
		},
	}

	t.Run("test del not existing key", func(t *testing.T) {
		table.Del("key_3")

		_, found := table.data["key_3"]
		require.False(t, found)
		_, found = table.data["key_2"]
		require.True(t, found)
		_, found = table.data["key_1"]
		require.True(t, found)
	})

	t.Run("test del existing key", func(t *testing.T) {
		table.Del("key_1")

		_, found := table.data["key_1"]
		require.False(t, found)
		_, found = table.data["key_2"]
		require.True(t, found)
	})
}
