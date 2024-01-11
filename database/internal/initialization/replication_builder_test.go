package initialization

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"spider/internal/configuration"
	"testing"
	"time"
)

func TestCreateReplicaWithoutConfig(t *testing.T) {
	t.Parallel()

	replica, err := CreateReplica(nil, nil, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, replica)
}

func TestCreateReplicaWithEmptyConfigFields(t *testing.T) {
	t.Parallel()

	replica, err := CreateReplica(&configuration.ReplicationConfig{MasterAddress: "localhost:4441"}, nil, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, replica)
}

func TestCreateReplicaWithIncorrectType(t *testing.T) {
	t.Parallel()

	replica, err := CreateReplica(&configuration.ReplicationConfig{ReplicaType: "non-master"}, nil, zap.NewNop())
	require.Error(t, err)
	require.Nil(t, replica)
}

func TestCreateReplica(t *testing.T) {
	t.Parallel()

	cfg := &configuration.ReplicationConfig{
		ReplicaType:   "master",
		MasterAddress: "localhost:4442",
		SyncInterval:  time.Minute,
	}

	replica, err := CreateReplica(cfg, nil, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, replica)
}
