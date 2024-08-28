package storage

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"spider/internal/database/compute"
	"spider/internal/database/storage/wal"
	"spider/internal/tools"
)

var (
	ErrorNotFound  = errors.New("not found")
	ErrorMutableTX = errors.New("mutable transaction on slave")
)

type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

type WAL interface {
	Recover() ([]wal.LogData, error)
	Set(context.Context, string, string) tools.FutureError
	Del(context.Context, string) tools.FutureError
	Shutdown()
}

type Replica interface {
	IsMaster() bool
	Shutdown()
}

type Storage struct {
	engine    Engine
	replica   Replica
	wal       WAL
	stream    <-chan []wal.LogData
	generator *IDGenerator
	logger    *zap.Logger
}

func NewStorage(engine Engine, logger *zap.Logger, options ...StorageOption) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	storage := &Storage{
		engine:    engine,
		logger:    logger,
		generator: NewIDGenerator(0), // TODO: need to update after recovering
	}

	for _, option := range options {
		option(storage)
	}

	if storage.wal != nil {
		logs, err := storage.wal.Recover()
		if err != nil {
			logger.Error("failed to recover data from WAL", zap.Error(err))
		} else {
			storage.applyData(logs)
		}
	}

	if storage.stream != nil {
		go func() {
			for logs := range storage.stream {
				storage.applyData(logs)
			}
		}()
	}

	return storage, nil
}

func (s *Storage) Set(ctx context.Context, key, value string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return ErrorMutableTX
	} else if ctx.Err() != nil {
		return ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = context.WithValue(ctx, "tx", txID)

	if s.wal != nil {
		future := s.wal.Set(ctx, key, value)
		if err := future.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)
	return nil
}

func (s *Storage) Del(ctx context.Context, key string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return ErrorMutableTX
	} else if ctx.Err() != nil {
		return ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = context.WithValue(ctx, "tx", txID)

	if s.wal != nil {
		future := s.wal.Del(ctx, key)
		if err := future.Get(); err != nil {
			return err
		}
	}

	s.engine.Del(ctx, key)
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = context.WithValue(ctx, "tx", txID)

	value, found := s.engine.Get(ctx, key)
	if !found {
		return "", ErrorNotFound
	}

	return value, nil
}

func (s *Storage) applyData(logs []wal.LogData) {
	for _, log := range logs {
		ctx := context.WithValue(context.Background(), "tx", log.LSN)
		switch log.CommandID {
		case compute.SetCommandID:
			s.engine.Set(ctx, log.Arguments[0], log.Arguments[1])
		case compute.DelCommandID:
			s.engine.Del(ctx, log.Arguments[0])
		}
	}
}
