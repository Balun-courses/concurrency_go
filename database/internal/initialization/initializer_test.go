package initialization

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"spider/internal/configuration"
)

func TestInitializer(t *testing.T) {
	t.Parallel()

	initializer, err := NewInitializer(&configuration.Config{
		Network: &configuration.NetworkConfig{
			Address: "localhost:30003",
		},
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	err = initializer.StartDatabase(ctx)
	require.NoError(t, err)
}
