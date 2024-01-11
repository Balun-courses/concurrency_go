package initialization

import (
	"context"
	"github.com/stretchr/testify/require"
	"spider/internal/configuration"
	"testing"
	"time"
)

func TestFailedInitializerCreation(t *testing.T) {
	t.Parallel()

	initializer, err := NewInitializer(&configuration.Config{Logging: &configuration.LoggingConfig{Level: "incorrect"}})
	require.Error(t, err)
	require.Nil(t, initializer)

	initializer, err = NewInitializer(&configuration.Config{WAL: &configuration.WALConfig{MaxSegmentSize: "100PB"}})
	require.Error(t, err)
	require.Nil(t, initializer)

	initializer, err = NewInitializer(&configuration.Config{Engine: &configuration.EngineConfig{Type: "incorrect"}})
	require.Error(t, err)
	require.Nil(t, initializer)

	initializer, err = NewInitializer(&configuration.Config{Network: &configuration.NetworkConfig{MaxMessageSize: "10PB"}})
	require.Error(t, err)
	require.Nil(t, initializer)

	initializer, err = NewInitializer(&configuration.Config{Replication: &configuration.ReplicationConfig{ReplicaType: "non-master"}})
	require.Error(t, err)
	require.Nil(t, initializer)
}

func TestInitializer(t *testing.T) {
	t.Parallel()

	initializer, err := NewInitializer(&configuration.Config{})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	err = initializer.StartDatabase(ctx)
	require.NoError(t, err)
}
