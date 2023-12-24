package network

import (
	"context"
	"github.com/stretchr/testify/require"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestTCPClient(t *testing.T) {
	t.Parallel()

	request := "hello server"
	response := "hello client"

	listener, err := net.Listen("tcp", ":10001")
	require.NoError(t, err)

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}

		buffer := make([]byte, 2048)
		count, err := connection.Read(buffer)
		require.NoError(t, err)
		require.True(t, reflect.DeepEqual([]byte(request), buffer[:count]))

		_, err = connection.Write([]byte(response))
		require.NoError(t, err)

		defer func() {
			err = connection.Close()
			require.NoError(t, err)
			err = listener.Close()
			require.NoError(t, err)
		}()
	}()

	time.Sleep(100 * time.Millisecond)

	client, err := NewTCPClient("127.0.0.1:10001", 2048, time.Minute)
	require.NoError(t, err)

	buffer, err := client.Send([]byte(request))
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual([]byte(response), buffer))
}

func TestTCPIdleClientConnection(t *testing.T) {
	t.Parallel()

	request := "hello server"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener, err := net.Listen("tcp", ":10002")
	require.NoError(t, err)

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}

		buffer := make([]byte, 2048)
		count, err := connection.Read(buffer)
		require.NoError(t, err)
		require.True(t, reflect.DeepEqual([]byte(request), buffer[:count]))

		<-ctx.Done()
		defer func() {
			err = connection.Close()
			require.NoError(t, err)
			err = listener.Close()
			require.NoError(t, err)
		}()
	}()

	time.Sleep(100 * time.Millisecond)

	client, err := NewTCPClient("127.0.0.1:10002", 2048, time.Millisecond*50)
	require.NoError(t, err)

	_, err = client.Send([]byte(request))
	require.Error(t, err)
}
