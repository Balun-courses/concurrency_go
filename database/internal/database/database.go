package database

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"spider/internal/database/compute"
)

type computeLayer interface {
	HandleQuery(context.Context, string) (compute.Query, error)
}

type storageLayer interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type Database struct {
	computeLayer computeLayer
	storageLayer storageLayer
	idGenerator  *IDGenerator
	logger       *zap.Logger
}

func NewDatabase(computeLayer computeLayer, storageLayer storageLayer, logger *zap.Logger) (*Database, error) {
	if computeLayer == nil {
		return nil, errors.New("compute is invalid")
	}

	if storageLayer == nil {
		return nil, errors.New("storage is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Database{
		computeLayer: computeLayer,
		storageLayer: storageLayer,
		idGenerator:  NewIDGenerator(),
		logger:       logger,
	}, nil
}

func (d *Database) HandleQuery(ctx context.Context, queryStr string) string {
	txID := d.idGenerator.Generate()
	ctx = context.WithValue(ctx, "tx", txID)

	d.logger.Debug(
		"handling query",
		zap.Int64("tx", txID),
		zap.String("query", queryStr),
	)

	query, err := d.computeLayer.HandleQuery(ctx, queryStr)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	switch query.CommandID() {
	case compute.SetCommandID:
		return d.handleSetQuery(ctx, query)
	case compute.GetCommandID:
		return d.handleGetQuery(ctx, query)
	case compute.DelCommandID:
		return d.handleDelQuery(ctx, query)
	}

	d.logger.Error("compute layer is incorrect", zap.Int64("tx", txID))
	return "[error] internal configuration error"
}

func (d *Database) handleSetQuery(ctx context.Context, query compute.Query) string {
	arguments := query.Arguments()
	if err := d.storageLayer.Set(ctx, arguments[0], arguments[1]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[ok]"
}

func (d *Database) handleGetQuery(ctx context.Context, query compute.Query) string {
	arguments := query.Arguments()
	value, err := d.storageLayer.Get(ctx, arguments[0])
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return fmt.Sprintf("[ok] %s", value)
}

func (d *Database) handleDelQuery(ctx context.Context, query compute.Query) string {
	arguments := query.Arguments()
	if err := d.storageLayer.Del(ctx, arguments[0]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[ok]"
}
