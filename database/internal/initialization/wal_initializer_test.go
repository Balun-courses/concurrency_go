package initialization

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"spider/internal/configuration"
)

func TestCreateWAL(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg    *configuration.WALConfig
		logger *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create wal without logger": {
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create wal without config": {
			logger:         zap.NewNop(),
			expectedNilObj: true,
		},
		"create wal with empty config fields": {
			cfg:    &configuration.WALConfig{},
			logger: zap.NewNop(),
		},
		"create wal with incorrect size": {
			cfg: &configuration.WALConfig{
				MaxSegmentSize: "incorrect",
			},
			logger:         zap.NewNop(),
			expectedErr:    errors.New("max segment size is incorrect"),
			expectedNilObj: true,
		},
		"create wal with config fields": {
			cfg: &configuration.WALConfig{
				FlushingBatchLength:  100,
				FlushingBatchTimeout: time.Millisecond,
				MaxSegmentSize:       "10MB",
				DataDirectory:        "./temp",
			},
			logger: zap.NewNop(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wal, err := CreateWAL(test.cfg, test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, wal)
			} else {
				assert.NotNil(t, wal)
			}
		})
	}
}
