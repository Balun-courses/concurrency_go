package wal

import (
	"context"
	"errors"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spider/internal/common"
)

// mockgen -source=wal.go -destination=wal_mock.go -package=wal

func TestNewWAL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		reader logsReader
		writer logsWriter

		expectedErr    error
		expectedNilObj bool
	}{
		"create wal without writer": {
			expectedErr:    errors.New("writer is invalid"),
			expectedNilObj: true,
		},
		"create wal without reader": {
			writer:         NewMocklogsWriter(ctrl),
			expectedErr:    errors.New("reader is invalid"),
			expectedNilObj: true,
		},
		"create wal": {
			reader: NewMocklogsReader(ctrl),
			writer: NewMocklogsWriter(ctrl),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wal, err := NewWAL(test.writer, test.reader, time.Millisecond, 100)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, wal)
			} else {
				assert.NotNil(t, wal)
			}
		})
	}
}

func TestWALFlushByTimeout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	logsReader := NewMocklogsReader(ctrl)
	logsWriter := NewMocklogsWriter(ctrl)
	logsWriter.EXPECT().
		Write(gomock.Any()).
		Do(func(requests []WriteRequest) {
			for _, request := range requests {
				request.SetResponse(nil)
			}
		})

	const timeout = 50 * time.Millisecond
	wal, err := NewWAL(logsWriter, logsReader, timeout, 1000)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wal.Start(ctx)

	future1 := wal.Set(common.ContextWithTxID(context.Background(), 10), "key1", "value1")
	future2 := wal.Set(common.ContextWithTxID(context.Background(), 20), "key2", "value2")

	time.Sleep(timeout)
	assert.NoError(t, future1.Get())
	assert.NoError(t, future2.Get())
}

func TestWALFlushBySize(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	logsReader := NewMocklogsReader(ctrl)
	logsWriter := NewMocklogsWriter(ctrl)
	logsWriter.EXPECT().
		Write(gomock.Any()).
		Do(func(requests []WriteRequest) {
			for _, request := range requests {
				request.SetResponse(nil)
			}
		})

	wal, err := NewWAL(logsWriter, logsReader, time.Minute, 2)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wal.Start(ctx)

	future1 := wal.Set(common.ContextWithTxID(context.Background(), 10), "key1", "value1")
	future2 := wal.Set(common.ContextWithTxID(context.Background(), 20), "key2", "value2")

	time.Sleep(100 * time.Millisecond)
	assert.NoError(t, future1.Get())
	assert.NoError(t, future2.Get())
}
