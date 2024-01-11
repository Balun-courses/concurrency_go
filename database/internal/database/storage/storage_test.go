package storage

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"spider/internal/tools"
	"testing"
)

// mockgen -source=storage.go -destination=storage_mock.go -package=storage

func TestNewStorage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)

	storage, err := NewStorage(nil, nil, nil, nil)
	require.Error(t, err, "engine is invalid")
	require.Nil(t, storage)

	storage, err = NewStorage(engine, nil, nil, nil)
	require.Error(t, err, "logger is invalid")
	require.Nil(t, storage)

	storage, err = NewStorage(engine, nil, nil, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestSetWithCanceledContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel()

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)

	storage, err := NewStorage(engine, nil, nil, zap.NewNop())
	require.NoError(t, err)

	err = storage.Set(ctxWithCancel, "key", "value")
	require.Error(t, err, context.Canceled)
}

func TestSetWithWALError(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	result := make(chan error, 1)
	result <- errors.New("wal error")

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)

	wal := NewMockWAL(ctrl)
	wal.EXPECT().
		Set(ctx, "key", "value").
		Return(tools.NewFuture(result))

	storage, err := NewStorage(engine, wal, nil, zap.NewNop())
	require.NoError(t, err)

	err = storage.Set(ctx, "key", "value")
	require.Error(t, err, "wal error")
}

func TestSuccessfulSet(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	result := make(chan error, 1)
	result <- nil

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)
	engine.EXPECT().
		Set(ctx, "key", "value")

	wal := NewMockWAL(ctrl)
	wal.EXPECT().
		Set(ctx, "key", "value").
		Return(tools.NewFuture(result))

	storage, err := NewStorage(engine, wal, nil, zap.NewNop())
	require.NoError(t, err)

	err = storage.Set(ctx, "key", "value")
	require.NoError(t, err)
}

func TestGetWithCanceledContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel()

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)

	storage, err := NewStorage(engine, nil, nil, zap.NewNop())
	require.NoError(t, err)

	value, err := storage.Get(ctxWithCancel, "key")
	require.Error(t, err, context.Canceled)
	require.Equal(t, "", value)
}

func TestSuccessfulGet(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	result := make(chan error, 1)
	result <- nil

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)
	engine.EXPECT().
		Get(ctx, "key").Return("value", true)

	storage, err := NewStorage(engine, nil, nil, zap.NewNop())
	require.NoError(t, err)

	value, err := storage.Get(ctx, "key")
	require.NoError(t, err)
	require.Equal(t, "value", value)
}

func TestDelWithCanceledContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel()

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)

	storage, err := NewStorage(engine, nil, nil, zap.NewNop())
	require.NoError(t, err)

	err = storage.Del(ctxWithCancel, "key")
	require.Error(t, err, context.Canceled)
}

func TestDelWithWALError(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	result := make(chan error, 1)
	result <- errors.New("wal error")

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)

	wal := NewMockWAL(ctrl)
	wal.EXPECT().
		Del(ctx, "key").
		Return(tools.NewFuture(result))

	storage, err := NewStorage(engine, wal, nil, zap.NewNop())
	require.NoError(t, err)

	err = storage.Del(ctx, "key")
	require.Error(t, err, "wal error")
}

func TestSuccessfulDel(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	result := make(chan error, 1)
	result <- nil

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)
	engine.EXPECT().
		Del(ctx, "key")

	wal := NewMockWAL(ctrl)
	wal.EXPECT().
		Del(ctx, "key").
		Return(tools.NewFuture(result))

	storage, err := NewStorage(engine, wal, nil, zap.NewNop())
	require.NoError(t, err)

	err = storage.Del(ctx, "key")
	require.NoError(t, err)
}
