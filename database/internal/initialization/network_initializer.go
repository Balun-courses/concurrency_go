package initialization

import (
	"errors"

	"go.uber.org/zap"

	"spider/internal/common"
	"spider/internal/configuration"
	"spider/internal/network"
)

const defaultServerAddress = ":3223"

func CreateNetwork(cfg *configuration.NetworkConfig, logger *zap.Logger) (*network.TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	address := defaultServerAddress
	var options []network.TCPServerOption

	if cfg != nil {
		if cfg.Address != "" {
			address = cfg.Address
		}

		if cfg.MaxConnections != 0 {
			options = append(options, network.WithServerMaxConnectionsNumber(uint(cfg.MaxConnections)))
		}

		if cfg.MaxMessageSize != "" {
			size, err := common.ParseSize(cfg.MaxMessageSize)
			if err != nil {
				return nil, errors.New("incorrect max message size")
			}

			options = append(options, network.WithServerBufferSize(uint(size)))
		}

		if cfg.IdleTimeout != 0 {
			options = append(options, network.WithServerIdleTimeout(cfg.IdleTimeout))
		}
	}

	return network.NewTCPServer(address, logger, options...)
}
