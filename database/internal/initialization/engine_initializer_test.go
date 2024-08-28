package initialization

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"spider/internal/configuration"
)

func TestCreateEngine(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg            *configuration.EngineConfig
		logger         *zap.Logger
		expectedErr    error
		expectedNilObj bool
	}{
		"create engine without logger": {
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create engine without config": {
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
		"create engine with empty config fields": {
			cfg:         &configuration.EngineConfig{},
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
		"create engine with config fields": {
			cfg: &configuration.EngineConfig{
				Type:             "in_memory",
				PartitionsNumber: 10,
			},
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
		"create engine with incorrect type": {
			cfg:            &configuration.EngineConfig{Type: "incorrect"},
			logger:         zap.NewNop(),
			expectedErr:    errors.New("engine type is incorrect"),
			expectedNilObj: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			engine, err := CreateEngine(test.cfg, test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, engine)
			} else {
				assert.NotNil(t, engine)
			}
		})
	}
}
