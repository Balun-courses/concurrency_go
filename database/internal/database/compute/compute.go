package compute

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

type parser interface {
	ParseQuery(context.Context, string) ([]string, error)
}

type analyzer interface {
	AnalyzeQuery(context.Context, []string) (Query, error)
}

type Compute struct {
	parser   parser
	analyzer analyzer
	logger   *zap.Logger
}

func NewCompute(parser parser, analyzer analyzer, logger *zap.Logger) (*Compute, error) {
	if parser == nil {
		return nil, errors.New("query parser is invalid")
	}

	if parser == nil {
		return nil, errors.New("query analyzer is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Compute{
		parser:   parser,
		analyzer: analyzer,
		logger:   logger,
	}, nil
}

func (d *Compute) HandleQuery(ctx context.Context, queryStr string) (Query, error) {
	tokens, err := d.parser.ParseQuery(ctx, queryStr)
	if err != nil {
		return Query{}, err
	}

	query, err := d.analyzer.AnalyzeQuery(ctx, tokens)
	if err != nil {
		return Query{}, err
	}

	return query, nil
}
