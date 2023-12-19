package initialization

import (
	"go.uber.org/zap"
	"spider/internal/configuration"
	"spider/internal/network"
	"time"
)

const defaultServerAddress = "localhost:3223"
const defaultMaxConnectionNumber = 100
const defaultMaxMessageSize = 2048
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

		if cfg.MaxMessageSize != 0 {
			maxMessageSize = cfg.MaxMessageSize
		}

		if cfg.IdleTimeout != 0 {
			idleTimeout = cfg.IdleTimeout
		}
	}

	return network.NewTCPServer(address, maxConnectionsNumber, maxMessageSize, idleTimeout, logger)
}
