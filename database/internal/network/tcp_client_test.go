package network

import (
	"errors"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTCPClient(t *testing.T) {
	t.Parallel()

	const serverResponse = "hello client"
	const serverAdress = "localhost:11111"
	listener, err := net.Listen("tcp", serverAdress)
	require.NoError(t, err)

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}

			_, err = connection.Read(make([]byte, 2048))
			require.NoError(t, err)

			_, err = connection.Write([]byte(serverResponse))
			require.NoError(t, err)
		}
	}()

	tests := map[string]struct {
		request string
		client  func() *TCPClient

		expectedResponse string
		expectedErr      error
	}{
		"client with incorrect server address": {
			request: "hello server",
			client: func() *TCPClient {
				client, err := NewTCPClient("localhost:1010")
				require.ErrorIs(t, err, syscall.ECONNREFUSED)
				return client
			},
			expectedResponse: serverResponse,
		},
		"client with small max message size": {
			request: "hello server",
			client: func() *TCPClient {
				client, err := NewTCPClient(serverAdress, WithClientBufferSize(5))
				require.NoError(t, err)
				return client
			},
			expectedErr: errors.New("small buffer size"),
		},
		"client with idle timeout": {
			request: "hello server",
			client: func() *TCPClient {
				client, err := NewTCPClient(serverAdress, WithClientIdleTimeout(100*time.Millisecond))
				require.NoError(t, err)
				return client
			},
			expectedResponse: serverResponse,
		},
		"client without options": {
			request: "hello server",
			client: func() *TCPClient {
				client, err := NewTCPClient(serverAdress)
				require.NoError(t, err)
				return client
			},
			expectedResponse: serverResponse,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := test.client()
			if client == nil {
				return
			}

			response, err := client.Send([]byte(test.request))
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResponse, string(response))
			client.Close()
		})
	}
}
