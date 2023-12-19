package initialization

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"spider/internal/configuration"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

var supportedLoggingLevels = map[string]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

const defaultEncoding = "json"
const defaultLevel = zapcore.InfoLevel
const defaultOutputPath = "output.log"

func CreateLogger(cfg *configuration.LoggingConfig) (*zap.Logger, error) {
	level := defaultLevel
	output := defaultOutputPath

	if cfg != nil {
		if cfg.Level != "" {
			var found bool
			if level, found = supportedLoggingLevels[cfg.Level]; !found {
				return nil, errors.New("logging level is incorrect")
			}
		}

		if cfg.Output != "" {
			output = cfg.Output
		}
	}

	loggerCfg := zap.Config{
		Encoding:    defaultEncoding,
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{output},
	}

	return loggerCfg.Build()
}
