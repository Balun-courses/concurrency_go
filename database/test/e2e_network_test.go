package main

import (
	"net"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2ENetwork(t *testing.T) {
	buffer := make([]byte, 1024)
	const serverAddress = "localhost:3223"

	cmd := exec.Command("../spider-server")
	require.NoError(t, cmd.Start())

	time.Sleep(time.Second)

	connection, clientErr := net.Dial("tcp", serverAddress)
	require.NoError(t, clientErr)

	_, clientErr = connection.Write([]byte("GET key1"))
	require.NoError(t, clientErr)

	size, clientErr := connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "[not found]", string(buffer[:size]))

	_, clientErr = connection.Write([]byte("SET key1 value1"))
	require.NoError(t, clientErr)

	size, clientErr = connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "[ok]", string(buffer[:size]))

	_, clientErr = connection.Write([]byte("GET key1"))
	require.NoError(t, clientErr)

	size, clientErr = connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "[ok] value1", string(buffer[:size]))

	_, clientErr = connection.Write([]byte("DEL key1"))
	require.NoError(t, clientErr)

	size, clientErr = connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "[ok]", string(buffer[:size]))

	_, clientErr = connection.Write([]byte("GET key1"))
	require.NoError(t, clientErr)

	size, clientErr = connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "[not found]", string(buffer[:size]))

	require.NoError(t, connection.Close())
	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))

	time.Sleep(time.Second)
}
