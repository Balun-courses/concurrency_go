package initialization

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"spider/internal/configuration"
)

func TestCreateReplica(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		replicaCfg *configuration.ReplicationConfig
		walCfg     *configuration.WALConfig
		logger     *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create replica without logger": {
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create replica without replica config": {
			logger:         zap.NewNop(),
			expectedNilObj: true,
		},
		"create replica without wal config": {
			replicaCfg:     &configuration.ReplicationConfig{},
			logger:         zap.NewNop(),
			expectedErr:    errors.New("replication without wal"),
			expectedNilObj: true,
		},
		"create replica with empty config fields": {
			replicaCfg:     &configuration.ReplicationConfig{},
			walCfg:         &configuration.WALConfig{},
			logger:         zap.NewNop(),
			expectedErr:    errors.New("replica type is incorrect"),
			expectedNilObj: true,
		},
		"create replica with incorrect type": {
			replicaCfg:     &configuration.ReplicationConfig{ReplicaType: "incorrect"},
			walCfg:         &configuration.WALConfig{},
			logger:         zap.NewNop(),
			expectedErr:    errors.New("replica type is incorrect"),
			expectedNilObj: true,
		},
		"create replica with incorrect master address": {
			replicaCfg:     &configuration.ReplicationConfig{ReplicaType: masterType},
			walCfg:         &configuration.WALConfig{},
			logger:         zap.NewNop(),
			expectedErr:    errors.New("master address is incorrect"),
			expectedNilObj: true,
		},
		"create master with config fields": {
			replicaCfg: &configuration.ReplicationConfig{
				ReplicaType:       masterType,
				MasterAddress:     "localhost:9090",
				SyncInterval:      time.Second * 10,
				MaxReplicasNumber: 10,
			},
			walCfg: &configuration.WALConfig{
				DataDirectory:  "./temp",
				MaxSegmentSize: "16MB",
			},
			logger: zap.NewNop(),
		},
		"create slave with config fields": {
			replicaCfg: &configuration.ReplicationConfig{
				ReplicaType:   slaveType,
				MasterAddress: "localhost:9090",
				SyncInterval:  time.Second * 10,
			},
			walCfg: &configuration.WALConfig{
				DataDirectory:  "./temp",
				MaxSegmentSize: "16MB",
			},
			logger: zap.NewNop(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			replica, err := CreateReplica(test.replicaCfg, test.walCfg, test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, replica)
			} else {
				assert.NotNil(t, replica)
			}
		})
	}
}
