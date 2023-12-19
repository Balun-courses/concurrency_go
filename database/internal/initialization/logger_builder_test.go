package initialization

import (
	"github.com/stretchr/testify/require"
	"spider/internal/configuration"
	"testing"
)

func TestCreateLoggerWithoutConfig(t *testing.T) {
	t.Parallel()

	logger, err := CreateLogger(nil)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestCreateLoggerWithEmptyConfigFields(t *testing.T) {
	t.Parallel()

	logger, err := CreateLogger(&configuration.LoggingConfig{})
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestCreateLoggerWithIncorrectLevel(t *testing.T) {
	t.Parallel()

	logger, err := CreateLogger(&configuration.LoggingConfig{Level: "incorrect"})
	require.Error(t, err, "logging level is incorrect")
	require.Nil(t, logger)
}

func TestCreateLogger(t *testing.T) {
	t.Parallel()

	cfg := &configuration.LoggingConfig{
		Level:  DebugLevel,
		Output: "test_output.log",
	}

	logger, err := CreateLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)
}
