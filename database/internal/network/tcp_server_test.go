package network

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTCPServer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverAddress := "localhost:22222"
	server, err := NewTCPServer(serverAddress, zap.NewNop())
	require.NoError(t, err)

	go func() {
		server.HandleQueries(ctx, func(ctx context.Context, data []byte) []byte {
			return []byte("hello-" + string(data))
		})
	}()

	time.Sleep(100 * time.Millisecond)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		connection, clientErr := net.Dial("tcp", serverAddress)
		require.NoError(t, clientErr)

		_, clientErr = connection.Write([]byte("client-1"))
		require.NoError(t, clientErr)

		buffer := make([]byte, 1024)
		size, clientErr := connection.Read(buffer)
		require.NoError(t, clientErr)

		clientErr = connection.Close()
		require.NoError(t, clientErr)

		assert.Equal(t, "hello-client-1", string(buffer[:size]))
	}()

	go func() {
		defer wg.Done()

		connection, clientErr := net.Dial("tcp", serverAddress)
		require.NoError(t, clientErr)

		_, clientErr = connection.Write([]byte("client-2"))
		require.NoError(t, clientErr)

		buffer := make([]byte, 1024)
		size, clientErr := connection.Read(buffer)
		require.NoError(t, clientErr)

		clientErr = connection.Close()
		require.NoError(t, clientErr)

		assert.Equal(t, "hello-client-2", string(buffer[:size]))
	}()

	wg.Wait()
	cancel()
}
