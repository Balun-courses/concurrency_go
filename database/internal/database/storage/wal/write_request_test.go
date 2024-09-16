package wal

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"spider/internal/database/compute"
)

func TestNewWriteRequest(t *testing.T) {
	t.Parallel()

	lsn := int64(100)
	commandID := compute.GetCommandID
	argumnets := []string{"key"}

	request := NewWriteRequest(lsn, compute.GetCommandID, []string{"key"})
	assert.Equal(t, lsn, request.log.LSN)
	assert.Equal(t, commandID, request.log.CommandID)
	assert.True(t, reflect.DeepEqual(argumnets, request.log.Arguments))
}

func TestWriteRequestWithError(t *testing.T) {
	t.Parallel()

	request := NewWriteRequest(100, compute.GetCommandID, []string{"key"})
	future := request.FutureResponse()

	go func() {
		request.SetResponse(errors.New("error"))
	}()

	err := future.Get()
	assert.Error(t, err, "error")
}

func TestWriteRequest(t *testing.T) {
	t.Parallel()

	request := NewWriteRequest(100, compute.GetCommandID, []string{"key"})
	future := request.FutureResponse()

	go func() {
		request.SetResponse(nil)
	}()

	err := future.Get()
	assert.NoError(t, err)
}
