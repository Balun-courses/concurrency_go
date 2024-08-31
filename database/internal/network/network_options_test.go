package network

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithClientIdleTimeout(t *testing.T) {
	t.Parallel()

	idleTimeout := time.Second
	option := WithClientIdleTimeout(time.Second)

	var client TCPClient
	option(&client)

	assert.Equal(t, idleTimeout, client.idleTimeout)
}

func TestWithClientBufferSize(t *testing.T) {
	t.Parallel()

	var bufferSize uint = 10 << 10
	option := WithClientBufferSize(bufferSize)

	var client TCPClient
	option(&client)

	assert.Equal(t, bufferSize, uint(client.bufferSize))
}

func TestWithServerIdleTimeout(t *testing.T) {
	t.Parallel()

	idleTimeout := time.Second
	option := WithServerIdleTimeout(time.Second)

	var server TCPServer
	option(&server)

	assert.Equal(t, idleTimeout, server.idleTimeout)
}

func TestWithServerBufferSize(t *testing.T) {
	t.Parallel()

	var bufferSize uint = 10 << 10
	option := WithServerBufferSize(bufferSize)

	var server TCPServer
	option(&server)

	assert.Equal(t, bufferSize, uint(server.bufferSize))
}

func TestWithServerMaxConnectionsNumber(t *testing.T) {
	t.Parallel()

	var maxConnections uint = 100
	option := WithServerMaxConnectionsNumber(maxConnections)

	var server TCPServer
	option(&server)

	assert.Equal(t, maxConnections, uint(server.maxConnections))
}
