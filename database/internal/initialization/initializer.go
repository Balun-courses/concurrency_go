package initialization

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"spider/internal/configuration"
	"spider/internal/database"
	"spider/internal/database/compute"
	"spider/internal/database/storage"
	"spider/internal/database/storage/engine/in_memory"
	"spider/internal/database/storage/replication"
	"spider/internal/database/storage/wal"
	"spider/internal/network"
)

type Initializer struct {
	wal    *wal.WAL
	engine *in_memory.Engine
	server *network.TCPServer
	slave  *replication.Slave
	master *replication.Master
	logger *zap.Logger
}

func NewInitializer(cfg *configuration.Config) (*Initializer, error) {
	if cfg == nil {
		return nil, errors.New("failed to initialize: config is invalid")
	}

	logger, err := CreateLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	wal, err := CreateWAL(cfg.WAL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize wal: %w", err)
	}

	engine, err := CreateEngine(cfg.Engine, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize engine: %w", err)
	}

	server, err := CreateNetwork(cfg.Network, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network: %w", err)
	}

	replica, err := CreateReplica(cfg.Replication, cfg.WAL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize replication: %w", err)
	}

	initializer := &Initializer{
		wal:    wal,
		engine: engine,
		server: server,
		logger: logger,
	}

	switch v := replica.(type) {
	case *replication.Slave:
		initializer.slave = v
	case *replication.Master:
		initializer.master = v
	}

	return initializer, nil
}

func (i *Initializer) StartDatabase(ctx context.Context) error {
	compute, err := compute.NewCompute(i.logger)
	if err != nil {
		return err
	}

	var options []storage.StorageOption
	if i.wal != nil {
		options = append(options, storage.WithWAL(i.wal))
	}

	if i.master != nil {
		options = append(options, storage.WithReplication(i.master))
	} else if i.slave != nil {
		options = append(options, storage.WithReplication(i.slave))
		options = append(options, storage.WithReplicationStream(i.slave.ReplicationStream()))
	}

	// TODO: need to start WAL and replication
	storage, err := storage.NewStorage(i.engine, i.logger, options...)
	if err != nil {
		return err
	}

	database, err := database.NewDatabase(compute, storage, i.logger)
	if err != nil {
		return err
	}

	group, groupCtx := errgroup.WithContext(ctx)
	if i.master != nil {
		group.Go(func() error {
			return i.master.HandleSynchronizations(groupCtx)
		})
	}

	group.Go(func() error {
		return i.server.HandleQueries(groupCtx, func(ctx context.Context, query []byte) []byte {
			response := database.HandleQuery(ctx, string(query))
			return []byte(response)
		})
	})

	return group.Wait()
}
