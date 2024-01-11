package configuration

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLoadNonExistentFile(t *testing.T) {
	t.Parallel()

	cfg, err := Load("test_data/non_existent_config.yaml")
	require.Error(t, err)
	require.Nil(t, cfg)
}

func TestLoadWithEmptyFilename(t *testing.T) {
	t.Parallel()

	cfg, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestLoadEmptyConfig(t *testing.T) {
	t.Parallel()

	cfg, err := Load("test_data/empty_config.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cfg, err := Load("test_data/config.yaml")
	require.NoError(t, err)

	require.Equal(t, "in_memory", cfg.Engine.Type)

	require.Equal(t, 100, cfg.WAL.FlushingBatchLength)
	require.Equal(t, time.Millisecond*10, cfg.WAL.FlushingBatchTimeout)
	require.Equal(t, "10MB", cfg.WAL.MaxSegmentSize)
	require.Equal(t, "/data/spider/wal", cfg.WAL.DataDirectory)

	require.Equal(t, "slave", cfg.Replication.ReplicaType)
	require.Equal(t, "127.0.0.1:3232", cfg.Replication.MasterAddress)
	require.Equal(t, time.Second, cfg.Replication.SyncInterval)

	require.Equal(t, "127.0.0.1:3223", cfg.Network.Address)
	require.Equal(t, 100, cfg.Network.MaxConnections)
	require.Equal(t, "4KB", cfg.Network.MaxMessageSize)
	require.Equal(t, time.Minute*5, cfg.Network.IdleTimeout)

	require.Equal(t, "info", cfg.Logging.Level)
	require.Equal(t, "/log/output.log", cfg.Logging.Output)
}
