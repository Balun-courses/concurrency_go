package initialization

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"spider/internal/configuration"
	"testing"
)

func TestCreateEngineWithoutConfig(t *testing.T) {
	t.Parallel()

	engine, err := CreateEngine(nil, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, engine)
}

func TestCreateEngineWithEmptyConfigFields(t *testing.T) {
	t.Parallel()

	engine, err := CreateEngine(&configuration.EngineConfig{}, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, engine)
}

func TestCreateEngineWithIncorrectType(t *testing.T) {
	t.Parallel()

	engine, err := CreateEngine(&configuration.EngineConfig{Type: "incorrect"}, zap.NewNop())
	require.Error(t, err, "engine type is incorrect")
	require.Nil(t, engine)
}

func TestCreateEngine(t *testing.T) {
	t.Parallel()

	cfg := &configuration.EngineConfig{
		Type: "in_memory",
	}

	engine, err := CreateEngine(cfg, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, engine)
}
