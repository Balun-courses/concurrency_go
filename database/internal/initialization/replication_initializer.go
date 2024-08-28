package initialization

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"spider/internal/configuration"
	"spider/internal/database/storage/replication"
	"spider/internal/network"
)

const masterType = "master"
const slaveType = "slave"

const defaultReplicationMasterAddress = "localhost:3232"
const defaultReplicationSyncInterval = time.Second

func CreateReplica(
	replicationCfg *configuration.ReplicationConfig,
	walCfg *configuration.WALConfig,
	logger *zap.Logger,
) (interface{}, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	} else if replicationCfg == nil {
		return nil, nil
	} else if walCfg == nil && replicationCfg != nil {
		return nil, errors.New("replication without wal")
	}

	supportedTypes := map[string]struct{}{
		masterType: {},
		slaveType:  {},
	}

	if _, found := supportedTypes[replicationCfg.ReplicaType]; !found {
		return nil, errors.New("replica type is incorrect")
	}

	masterAddress := defaultReplicationMasterAddress
	syncInterval := defaultReplicationSyncInterval
	walDirectory := defaultWALDataDirectory // TODO

	if replicationCfg.MasterAddress != "" {
		masterAddress = replicationCfg.MasterAddress
	}

	if replicationCfg.SyncInterval != 0 {
		syncInterval = replicationCfg.SyncInterval
	}

	if walCfg.DataDirectory != "" {
		walDirectory = walCfg.DataDirectory
	}

	const maxReplicasNumber = 5
	const maxMessageSize = 16 << 20
	idleTimeout := syncInterval * 3

	if replicationCfg.ReplicaType == masterType {
		server, err := network.NewTCPServer(masterAddress, maxReplicasNumber, maxMessageSize, idleTimeout, logger)
		if err != nil {
			return nil, err
		}

		return replication.NewMaster(server, walDirectory, logger)
	} else {
		client, err := network.NewTCPClient(masterAddress, maxMessageSize, idleTimeout)
		if err != nil {
			return nil, err
		}

		return replication.NewSlave(client, syncInterval, logger)
	}
}
