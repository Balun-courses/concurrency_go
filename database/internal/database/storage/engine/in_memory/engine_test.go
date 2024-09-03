package in_memory

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"spider/internal/common"
)

func TestNewEngine(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		logger         *zap.Logger
		options        []EngineOption
		expectedErr    error
		expectedNilObj bool
	}{
		"create engine without logger": {
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create engine without options": {
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
		"create engine with partitions": {
			logger:      zap.NewNop(),
			options:     []EngineOption{WithPartitions(10)},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			engine, err := NewEngine(test.logger, test.options...)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, engine)
			} else {
				assert.NotNil(t, engine)
			}
		})
	}
}

func TestEngineSet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		engine *Engine
		key    string
		value  string
	}{
		"set with single partition": {
			engine: func() *Engine {
				engine, err := NewEngine(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
		"set with multiple partitions": {
			engine: func() *Engine {
				const partitionsNumber uint = 8
				engine, err := NewEngine(zap.NewNop(), WithPartitions(partitionsNumber))
				require.NoError(t, err)
				return engine
			}(),
			key:   "key",
			value: "value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const txID int64 = 1
			ctx := common.ContextWithTxID(context.Background(), txID)

			test.engine.Set(ctx, test.key, test.value)
			value, found := test.engine.Get(ctx, test.key)
			assert.True(t, found)
			assert.Equal(t, test.value, value)
		})
	}
}

func TestEngineDel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		engine *Engine
		key    string
	}{
		"del with single partition": {
			engine: func() *Engine {
				engine, err := NewEngine(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
		"del with multiple partitions": {
			engine: func() *Engine {
				const partitionsNumber uint = 8
				engine, err := NewEngine(zap.NewNop(), WithPartitions(partitionsNumber))
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const txID int64 = 1
			ctx := common.ContextWithTxID(context.Background(), txID)

			test.engine.Del(ctx, test.key)
			value, found := test.engine.Get(ctx, test.key)
			assert.False(t, found)
			assert.Empty(t, value)
		})
	}
}

func TestEngineGet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		engine *Engine
		key    string
	}{
		"get with single partition": {
			engine: func() *Engine {
				engine, err := NewEngine(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
		"get with multiple partitions": {
			engine: func() *Engine {
				const partitionsNumber uint = 8
				engine, err := NewEngine(zap.NewNop(), WithPartitions(partitionsNumber))
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const txID int64 = 1
			ctx := common.ContextWithTxID(context.Background(), txID)

			value, found := test.engine.Get(ctx, test.key)
			assert.False(t, found)
			assert.Empty(t, value)
		})
	}
}
