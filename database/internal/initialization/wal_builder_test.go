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

	wal, err := CreateWAL(nil, zap.NewNop())
	require.NoError(t, err)
	require.Nil(t, wal)
}

func TestCreateWALWithEmptyConfigFields(t *testing.T) {
	t.Parallel()

	wal, err := CreateWAL(&configuration.WALConfig{}, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, wal)
}

func TestCreateWALWithIncorrectSegmentSize(t *testing.T) {
	t.Parallel()

	wal, err := CreateWAL(&configuration.WALConfig{MaxSegmentSize: "100PB"}, zap.NewNop())
	require.Error(t, err, "max segment size is incorrect")
	require.Nil(t, wal)
}

func TestCreateWAL(t *testing.T) {
	t.Parallel()

	cfg := &configuration.WALConfig{
		FlushingBatchLength:  200,
		FlushingBatchTimeout: 20 * time.Millisecond,
		MaxSegmentSize:       "20MB",
		DataDirectory:        "/data/wal",
	}

	wal, err := CreateWAL(cfg, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, wal)
}
