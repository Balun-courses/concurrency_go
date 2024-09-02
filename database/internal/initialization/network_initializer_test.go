package initialization

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"spider/internal/configuration"
)

func TestCreateNetwork(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg    *configuration.NetworkConfig
		logger *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create network without logger": {
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create network without config": {
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
		"create network with empty config fields": {
			logger: zap.NewNop(),
			cfg: &configuration.NetworkConfig{
				Address: "localhost:20002",
			},
			expectedErr: nil,
		},
		"create network with config fields": {
			logger: zap.NewNop(),
			cfg: &configuration.NetworkConfig{
				Address:        "localhost:10001",
				MaxConnections: 100,
				MaxMessageSize: "2KB",
				IdleTimeout:    time.Second,
			},
			expectedErr: nil,
		},
		"create network with incorrect size": {
			logger: zap.NewNop(),
			cfg: &configuration.NetworkConfig{
				MaxMessageSize: "2incorrect",
			},
			expectedErr:    errors.New("incorrect max message size"),
			expectedNilObj: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			network, err := CreateNetwork(test.cfg, test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, network)
			} else {
				assert.NotNil(t, network)
			}
		})
	}
}
