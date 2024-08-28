package initialization

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"spider/internal/configuration"
)

func TestCreateLogger(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg            *configuration.LoggingConfig
		expectedErr    error
		expectedNilObj bool
	}{
		"create logger without config": {
			expectedErr: nil,
		},
		"create logger with empty config fields": {
			cfg:         &configuration.LoggingConfig{},
			expectedErr: nil,
		},
		"create logger with config fields": {
			cfg: &configuration.LoggingConfig{
				Level:  debugLevel,
				Output: "test.log",
			},
			expectedErr: nil,
		},
		"create logger with incorrect level": {
			cfg:            &configuration.LoggingConfig{Level: "incorrect"},
			expectedErr:    errors.New("logging level is incorrect"),
			expectedNilObj: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := CreateLogger(test.cfg)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, logger)
			} else {
				assert.NotNil(t, logger)
			}
		})
	}
}
