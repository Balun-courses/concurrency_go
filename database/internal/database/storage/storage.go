package storage

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"spider/internal/database/compute"
	"spider/internal/database/storage/wal"
	"spider/internal/tools"
)

type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

type WAL interface {
	Start()
	Recover() ([]wal.LogData, error)
	Set(context.Context, string, string) tools.FutureError
	Del(context.Context, string) tools.FutureError
	Shutdown()
}

type Storage struct {
	engine Engine
	wal    WAL
	logger *zap.Logger
}

func NewStorage(engine Engine, wal WAL, logger *zap.Logger) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	storage := &Storage{
		engine: engine,
		wal:    wal,
		logger: logger,
	}

	if wal != nil {
		records, err := wal.Recover()
		if err != nil {
			logger.Error("failed to recover database from WAL")
		}

		storage.recover(records)
		wal.Start()
	}

	return storage, nil
}

func (s *Storage) Set(ctx context.Context, key, value string) error {
	if s.wal != nil {
		future := s.wal.Set(ctx, key, value)
		if err := future.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	value, _ := s.engine.Get(ctx, key)
	return value, nil
}

func (s *Storage) Del(ctx context.Context, key string) error {
	if s.wal != nil {
		future := s.wal.Del(ctx, key)
		if err := future.Get(); err != nil {
			return err
		}
	}

	s.engine.Del(ctx, key)
	return nil
}

func (s *Storage) recover(logs []wal.LogData) {
	for _, log := range logs {
		switch log.CommandID {
		case compute.SetCommandID:
			s.engine.Set(context.Background(), log.Arguments[0], log.Arguments[1])
		case compute.DelCommandID:
			s.engine.Del(context.Background(), log.Arguments[0])
		}
	}
}
