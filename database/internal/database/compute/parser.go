package compute

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) (*Parser, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Parser{
		logger: logger,
	}, nil
}

func (p *Parser) ParseQuery(ctx context.Context, query string) ([]string, error) {
	machine := newStateMachine()
	tokens, err := machine.parse(query)
	if err != nil {
		return nil, err
	}

	txID := ctx.Value("tx").(int64)
	p.logger.Debug(
		"query parsed",
		zap.Int64("tx", txID),
		zap.Any("tokens", tokens),
	)

	return tokens, nil
}

func isWhiteSpace(symbol byte) bool {
	return symbol == '\t' || symbol == '\n' || symbol == ' '
}

func isLetter(symbol byte) bool {
	return (symbol >= 'a' && symbol <= 'z') ||
		(symbol >= 'A' && symbol <= 'Z') ||
		(symbol >= '0' && symbol <= '9') ||
		(symbol == '_')
}
