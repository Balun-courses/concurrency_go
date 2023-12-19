package initialization

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"spider/internal/configuration"
	"spider/internal/database"
	"spider/internal/database/compute"
	"spider/internal/database/storage"
	"spider/internal/network"
)

type Initializer struct {
	engine storage.Engine
	server *network.TCPServer
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

	dbEngine, err := CreateEngine(cfg.Engine, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize engine: %w", err)
	}

	tcpServer, err := CreateNetwork(cfg.Network, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network: %w", err)
	}

	return &Initializer{
		engine: dbEngine,
		server: tcpServer,
		logger: logger,
	}, nil
}

func (i *Initializer) StartDatabase(ctx context.Context) error {
	compute, err := i.createComputeLayer()
	if err != nil {
		return err
	}

	storage, err := i.createStorageLayer()
	if err != nil {
		return err
	}

	database, err := database.NewDatabase(compute, storage, i.logger)
	if err != nil {
		i.logger.Error("failed to start application", zap.Error(err))
		return err
	}

	i.server.HandleQueries(ctx, func(ctx context.Context, query []byte) []byte {
		response := database.HandleQuery(ctx, string(query))
		return []byte(response)
	})

	return nil
}

func (i *Initializer) createComputeLayer() (*compute.Compute, error) {
	queryParser, err := compute.NewParser(i.logger)
	if err != nil {
		i.logger.Error("failed to start application", zap.Error(err))
		return nil, err
	}

	queryAnalyzer, err := compute.NewAnalyzer(i.logger)
	if err != nil {
		i.logger.Error("failed to start application", zap.Error(err))
		return nil, err
	}

	compute, err := compute.NewCompute(queryParser, queryAnalyzer, i.logger)
	if err != nil {
		i.logger.Error("failed to start application", zap.Error(err))
		return nil, err
	}

	return compute, nil
}

func (i *Initializer) createStorageLayer() (*storage.Storage, error) {
	storage, err := storage.NewStorage(i.engine, nil, i.logger) // TODO
	if err != nil {
		i.logger.Error("failed to start application", zap.Error(err))
		return nil, err
	}

	return storage, nil
}
