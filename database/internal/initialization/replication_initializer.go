package initialization

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"spider/internal/configuration"
	"spider/internal/database/storage/replication"
	"spider/internal/network"
	"spider/internal/size"
)

const (
	masterType = "master"
	slaveType  = "slave"
)

const (
	defaultReplicationSyncInterval = time.Second
	defaultMaxReplicasNumber       = 5
)

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

	if replicationCfg.MasterAddress == "" {
		return nil, errors.New("master address is incorrect")
	}

	maxMessageSize := defaultMaxSegmentSize
	masterAddress := replicationCfg.MasterAddress
	syncInterval := defaultReplicationSyncInterval
	walDirectory := defaultWALDataDirectory

	if replicationCfg.SyncInterval != 0 {
		syncInterval = replicationCfg.SyncInterval
	}

	if walCfg.DataDirectory != "" {
		walDirectory = walCfg.DataDirectory
	}

	if walCfg.MaxSegmentSize != "" {
		size, _ := size.ParseSize(walCfg.MaxSegmentSize)
		maxMessageSize = size
	}

	idleTimeout := syncInterval * 3
	if replicationCfg.ReplicaType == masterType {
		maxReplicasNumber := defaultMaxReplicasNumber
		if replicationCfg.MaxReplicasNumber != 0 {
			maxReplicasNumber = replicationCfg.MaxReplicasNumber
		}

		var options []network.TCPServerOption
		options = append(options, network.WithServerIdleTimeout(idleTimeout))
		options = append(options, network.WithServerBufferSize(uint(maxMessageSize)))
		options = append(options, network.WithServerMaxConnectionsNumber(uint(maxReplicasNumber)))
		server, err := network.NewTCPServer(masterAddress, logger, options...)
		if err != nil {
			return nil, err
		}

		return replication.NewMaster(server, walDirectory, logger)
	} else {
		var options []network.TCPClientOption
		options = append(options, network.WithClientIdleTimeout(idleTimeout))
		options = append(options, network.WithClientBufferSize(uint(maxMessageSize)))
		client, err := network.NewTCPClient(masterAddress, options...)
		if err != nil {
			return nil, err
		}

		return replication.NewSlave(client, walDirectory, syncInterval, logger)
	}
}
