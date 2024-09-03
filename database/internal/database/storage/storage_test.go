package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"spider/internal/concurrency"
	"spider/internal/database/compute"
	"spider/internal/database/storage/wal"
)

// mockgen -source=storage.go -destination=storage_mock.go -package=storage

func TestNewStorage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	writeAheadLog := NewMockWAL(ctrl)
	writeAheadLog.EXPECT().
		Recover().
		Return(nil, nil)

	tests := map[string]struct {
		engine  Engine
		logger  *zap.Logger
		options []StorageOption

		expectedErr    error
		expectedNilObj bool
	}{
		"create storage without engine": {
			expectedErr:    errors.New("engine is invalid"),
			expectedNilObj: true,
		},
		"create storage without logger": {
			engine:         NewMockEngine(ctrl),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create engine without options": {
			engine:      NewMockEngine(ctrl),
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
		"create engine with replication stream": {
			engine:      NewMockEngine(ctrl),
			logger:      zap.NewNop(),
			options:     []StorageOption{WithReplicationStream(make(<-chan []wal.Log))},
			expectedErr: nil,
		},
		"create engine with wal": {
			engine:      NewMockEngine(ctrl),
			logger:      zap.NewNop(),
			options:     []StorageOption{WithWAL(writeAheadLog)},
			expectedErr: nil,
		},
		"create engine with replica": {
			engine:      NewMockEngine(ctrl),
			logger:      zap.NewNop(),
			options:     []StorageOption{WithReplication(NewMockReplica(ctrl))},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			storage, err := NewStorage(test.engine, test.logger, test.options...)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, storage)
			} else {
				assert.NotNil(t, storage)
			}
		})
	}
}

func TestStorageSet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		engine  func() Engine
		replica func() Replica
		wal     func() WAL

		expectedErr error
	}{
		"set with slave replica": {
			engine: func() Engine { return NewMockEngine(ctrl) },
			replica: func() Replica {
				replica := NewMockReplica(ctrl)
				replica.EXPECT().
					IsMaster().
					Return(false)
				return replica
			},
			wal: func() WAL {
				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				return wal
			},
			expectedErr: ErrorMutableTX,
		},
		"set without wal": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Set(gomock.Any(), "key", "value")
				return engine
			},
			replica: func() Replica { return nil },
			wal:     func() WAL { return nil },
		},
		"set with error from wal": {
			engine:  func() Engine { return NewMockEngine(ctrl) },
			replica: func() Replica { return nil },
			wal: func() WAL {
				result := make(chan error, 1)
				result <- errors.New("wal error")
				future := concurrency.NewFuture(result)

				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				wal.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(future)
				return wal
			},
			expectedErr: errors.New("wal error"),
		},
		"set with wal": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Set(gomock.Any(), "key", "value")
				return engine
			},
			replica: func() Replica { return nil },
			wal: func() WAL {
				result := make(chan error, 1)
				result <- nil
				future := concurrency.NewFuture(result)

				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				wal.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(future)
				return wal
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			options := []StorageOption{
				WithWAL(test.wal()),
				WithReplication(test.replica()),
			}

			storage, err := NewStorage(test.engine(), zap.NewNop(), options...)
			require.NoError(t, err)

			err = storage.Set(context.Background(), "key", "value")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestStorageDel(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		engine  func() Engine
		replica func() Replica
		wal     func() WAL

		expectedErr error
	}{
		"del with slave replica": {
			engine: func() Engine { return NewMockEngine(ctrl) },
			replica: func() Replica {
				replica := NewMockReplica(ctrl)
				replica.EXPECT().
					IsMaster().
					Return(false)
				return replica
			},
			wal: func() WAL {
				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				return wal
			},
			expectedErr: ErrorMutableTX,
		},
		"del without wal": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Del(gomock.Any(), "key")
				return engine
			},
			replica: func() Replica { return nil },
			wal:     func() WAL { return nil },
		},
		"del with error from wal": {
			engine:  func() Engine { return NewMockEngine(ctrl) },
			replica: func() Replica { return nil },
			wal: func() WAL {
				result := make(chan error, 1)
				result <- errors.New("wal error")
				future := concurrency.NewFuture(result)

				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				wal.EXPECT().
					Del(gomock.Any(), "key").
					Return(future)
				return wal
			},
			expectedErr: errors.New("wal error"),
		},
		"del with wal": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Del(gomock.Any(), "key")
				return engine
			},
			replica: func() Replica { return nil },
			wal: func() WAL {
				result := make(chan error, 1)
				result <- nil
				future := concurrency.NewFuture(result)

				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				wal.EXPECT().
					Del(gomock.Any(), "key").
					Return(future)
				return wal
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			options := []StorageOption{
				WithWAL(test.wal()),
				WithReplication(test.replica()),
			}

			storage, err := NewStorage(test.engine(), zap.NewNop(), options...)
			require.NoError(t, err)

			err = storage.Del(context.Background(), "key")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestStorageGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		engine func() Engine

		expectedValue string
		expectedErr   error
	}{
		"get with unexesiting element": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Get(gomock.Any(), "key").
					Return("", false)
				return engine
			},
			expectedErr: ErrorNotFound,
		},
		"get with exesiting element": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Get(gomock.Any(), "key").
					Return("value", true)
				return engine
			},
			expectedValue: "value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			storage, err := NewStorage(test.engine(), zap.NewNop())
			require.NoError(t, err)

			value, err := storage.Get(context.Background(), "key")
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func TestStorageWithReplicationStream(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)
	engine.EXPECT().
		Set(gomock.Any(), "key_1", "value_1")
	engine.EXPECT().
		Del(gomock.Any(), "key_2")

	replicationStream := make(chan []wal.Log)
	_, err := NewStorage(engine, zap.NewNop(), WithReplicationStream(replicationStream))
	require.NoError(t, err)

	replicationStream <- []wal.Log{
		{
			LSN:       1,
			CommandID: compute.SetCommandID,
			Arguments: []string{"key_1", "value_1"},
		},
		{
			LSN:       2,
			CommandID: compute.DelCommandID,
			Arguments: []string{"key_2"},
		},
	}

	close(replicationStream)
	time.Sleep(100 * time.Millisecond) // TODO: need to fix time waiting
}
