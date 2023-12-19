package compute

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

const (
	setQueryArgumentsNumber = 2
	getQueryArgumentsNumber = 1
	delQueryArgumentsNumber = 1
)

var (
	errInvalidSymbol    = errors.New("invalid symbol")
	errInvalidCommand   = errors.New("invalid command")
	errInvalidArguments = errors.New("invalid arguments")
)

type Analyzer struct {
	handlers []func(context.Context, Query) error
	logger   *zap.Logger
}

func NewAnalyzer(logger *zap.Logger) (*Analyzer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	analyser := &Analyzer{
		logger: logger,
	}

	analyser.handlers = []func(context.Context, Query) error{
		SetCommandID: analyser.analyzeSetQuery,
		GetCommandID: analyser.analyzeGetQuery,
		DelCommandID: analyser.analyzeDelQuery,
	}

	return analyser, nil
}

func (a *Analyzer) AnalyzeQuery(ctx context.Context, tokens []string) (Query, error) {
	if len(tokens) == 0 {
		txID := ctx.Value("tx").(int64)
		a.logger.Debug("invalid query", zap.Int64("tx", txID))
		return Query{}, errInvalidCommand
	}

	command := tokens[0]
	commandID := CommandNameToCommandID(command)
	if commandID == UnknownCommandID {
		txID := ctx.Value("tx").(int64)
		a.logger.Debug(
			"invalid command",
			zap.Int64("tx", txID),
			zap.String("command", command),
		)
		return Query{}, errInvalidCommand
	}

	query := NewQuery(commandID, tokens[1:])
	handler := a.handlers[commandID]
	if err := handler(ctx, query); err != nil {
		return Query{}, err
	}

	txID := ctx.Value("tx").(int64)
	a.logger.Debug(
		"query analyzed",
		zap.Int64("tx", txID),
		zap.Any("query", query),
	)

	return query, nil
}

func (a *Analyzer) analyzeSetQuery(ctx context.Context, query Query) error {
	if len(query.Arguments()) != setQueryArgumentsNumber {
		txID := ctx.Value("tx").(int64)
		a.logger.Debug(
			"invalid arguments for set query",
			zap.Int64("tx", txID),
			zap.Any("args", query.Arguments()),
		)
		return errInvalidArguments
	}

	return nil
}

func (a *Analyzer) analyzeGetQuery(ctx context.Context, query Query) error {
	if len(query.Arguments()) != getQueryArgumentsNumber {
		txID := ctx.Value("tx").(int64)
		a.logger.Debug(
			"invalid arguments for get query",
			zap.Int64("tx", txID),
			zap.Any("args", query.Arguments()),
		)
		return errInvalidArguments
	}

	return nil
}

func (a *Analyzer) analyzeDelQuery(ctx context.Context, query Query) error {
	if len(query.Arguments()) != delQueryArgumentsNumber {
		txID := ctx.Value("tx").(int64)
		a.logger.Debug(
			"invalid arguments for del query",
			zap.Int64("tx", txID),
			zap.Any("args", query.Arguments()),
		)
		return errInvalidArguments
	}

	return nil
}
