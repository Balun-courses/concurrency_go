package initialization

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"spider/internal/configuration"
	"testing"
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

func TestCreateNetwork(t *testing.T) {
	t.Parallel()

	cfg := &configuration.NetworkConfig{
		Address: "localhost:9898",
	}

	server, err := CreateNetwork(cfg, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, server)
}
