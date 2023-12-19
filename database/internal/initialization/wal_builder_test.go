package initialization

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"spider/internal/configuration"
	"testing"
	"time"
)

func TestCreateWALWithoutConfig(t *testing.T) {
	t.Parallel()

	logger, err := CreateWAL(nil, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestCreateWALWithEmptyConfigFields(t *testing.T) {
	t.Parallel()

	logger, err := CreateWAL(&configuration.WALConfig{}, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestCreateWALWithIncorrectSegmentSize(t *testing.T) {
	t.Parallel()

	logger, err := CreateWAL(&configuration.WALConfig{MaxSegmentSize: "100INCORRECT"}, zap.NewNop())
	require.Error(t, err, "max segment size is incorrect")
	require.Nil(t, logger)
}

func TestCreateWAL(t *testing.T) {
	t.Parallel()

	cfg := &configuration.WALConfig{
		FlushingBatchSize:    200,
		FlushingBatchTimeout: 20 * time.Millisecond,
		MaxSegmentSize:       "20MB",
		DataDirectory:        "/data/wal",
	}

	logger, err := CreateWAL(cfg, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, logger)
}
