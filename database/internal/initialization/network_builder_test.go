package initialization

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"spider/internal/configuration"
)

func TestCreateNetworkWithoutConfig(t *testing.T) {
	server, err := CreateNetwork(nil, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestCreateNetworkWithEmptyConfigFields(t *testing.T) {
	server, err := CreateNetwork(&configuration.NetworkConfig{}, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestCreateLoggerWithIncorrectMessageSize(t *testing.T) {
	t.Parallel()

	server, err := CreateNetwork(&configuration.NetworkConfig{MaxMessageSize: "10PB"}, zap.NewNop())
	require.Error(t, err)
	require.Nil(t, server)
}

func TestCreateNetwork(t *testing.T) {
	t.Parallel()

	cfg := &configuration.NetworkConfig{
		Address: "localhost:9898",
	}

	server, err := CreateNetwork(cfg, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, server)
}
