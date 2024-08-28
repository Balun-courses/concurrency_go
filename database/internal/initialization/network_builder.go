package initialization

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"spider/internal/configuration"
	"spider/internal/network"
	"spider/internal/tools"
)

const defaultServerAddress = "localhost:3223"
const defaultMaxConnectionNumber = 100
const defaultMaxMessageSize = 2 << 10
const defaultIdleTimeout = time.Minute * 5

func CreateNetwork(cfg *configuration.NetworkConfig, logger *zap.Logger) (*network.TCPServer, error) {
	address := defaultServerAddress
	maxConnectionsNumber := defaultMaxConnectionNumber
	maxMessageSize := defaultMaxMessageSize
	idleTimeout := defaultIdleTimeout

	if cfg != nil {
		if cfg.Address != "" {
			address = cfg.Address
		}

		if cfg.MaxConnections != 0 {
			maxConnectionsNumber = cfg.MaxConnections
		}

		if cfg.MaxMessageSize != "" {
			size, err := tools.ParseSize(cfg.MaxMessageSize)
			if err != nil {
				return nil, errors.New("incorrect max message size")
			}

			maxMessageSize = size
		}

		if cfg.IdleTimeout != 0 {
			idleTimeout = cfg.IdleTimeout
		}
	}

	return network.NewTCPServer(address, maxConnectionsNumber, maxMessageSize, idleTimeout, logger)
}
