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

var queryArgumentsNumber = map[int]int{
	SetCommandID: setQueryArgumentsNumber,
	GetCommandID: getQueryArgumentsNumber,
	DelCommandID: delQueryArgumentsNumber,
}

var (
	errInvalidSymbol    = errors.New("invalid symbol")
	errInvalidCommand   = errors.New("invalid command")
	errInvalidArguments = errors.New("invalid arguments")
)

type Analyzer struct {
	logger *zap.Logger
}

func NewAnalyzer(logger *zap.Logger) (*Analyzer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Analyzer{
		logger: logger,
	}, nil
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
	argumentsNumber := queryArgumentsNumber[commandID]
	if len(query.Arguments()) != argumentsNumber {
		txID := ctx.Value("tx").(int64)
		a.logger.Debug(
			"invalid arguments for query",
			zap.Int64("tx", txID),
			zap.Any("args", query.Arguments()),
		)
		return Query{}, errInvalidArguments
	}

	txID := ctx.Value("tx").(int64)
	a.logger.Debug(
		"query analyzed",
		zap.Int64("tx", txID),
		zap.Any("query", query),
	)

	return query, nil
}
