package in_memory

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

// mockgen -source=engine.go -destination=engine_mock.go -package=in_memory

func TestNewEngine(t *testing.T) {
	t.Parallel()

	tableBuilder := func() hashTable {
		ctrl := gomock.NewController(t)
		return NewMockhashTable(ctrl)
	}

	engine, err := NewEngine(nil, -1, nil)
	require.Error(t, err, "hash table builder is invalid")
	require.Nil(t, engine)

	engine, err = NewEngine(tableBuilder, -1, nil)
	require.Error(t, err, "partitions number is invalid")
	require.Nil(t, engine)

	engine, err = NewEngine(tableBuilder, 1, nil)
	require.Error(t, err, "logger is invalid")
	require.Nil(t, engine)

	engine, err = NewEngine(tableBuilder, 10, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, engine)
}

func TestSetQuery(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	tableBuilder := func() hashTable {
		ctrl := gomock.NewController(t)
		table := NewMockhashTable(ctrl)
		table.EXPECT().Set("key_1", "value_1")
		return table
	}

	engine, err := NewEngine(tableBuilder, 1, zap.NewNop())
	require.NoError(t, err)

	engine.Set(ctx, "key_1", "value_1")
}

func TestGetQuery(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	tableBuilder := func() hashTable {
		ctrl := gomock.NewController(t)
		table := NewMockhashTable(ctrl)
		table.EXPECT().Get("key_1").Return("value_1", true)
		table.EXPECT().Get("key_2").Return("", false)
		return table
	}

	engine, err := NewEngine(tableBuilder, 1, zap.NewNop())
	require.NoError(t, err)

	value, found := engine.Get(ctx, "key_1")
	require.Equal(t, "value_1", value)
	require.True(t, found)

	value, found = engine.Get(ctx, "key_2")
	require.Equal(t, "", value)
	require.False(t, found)
}

func TestDelQuery(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	tableBuilder := func() hashTable {
		ctrl := gomock.NewController(t)
		table := NewMockhashTable(ctrl)
		table.EXPECT().Del("key_1")
		return table
	}

	engine, err := NewEngine(tableBuilder, 1, zap.NewNop())
	require.NoError(t, err)

	engine.Del(ctx, "key_1")
}
