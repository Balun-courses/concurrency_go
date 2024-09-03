package wal

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spider/internal/database/compute"
)

func TestSingleSerialization(t *testing.T) {
	t.Parallel()

	expectedLog := Log{
		LSN:       100,
		CommandID: compute.SetCommandID,
		Arguments: []string{"key", "value"},
	}

	var writeBuffer bytes.Buffer
	err := expectedLog.Encode(&writeBuffer)
	require.NoError(t, err)

	data := writeBuffer.Bytes()
	readBuffer := bytes.NewBuffer(data)

	var log Log
	err = log.Decode(readBuffer)
	require.NoError(t, err)

	assert.True(t, reflect.DeepEqual(expectedLog, log))
}

func TestMultipleSerialization(t *testing.T) {
	t.Parallel()

	expectedLogs := []Log{
		{
			LSN:       200,
			CommandID: compute.SetCommandID,
			Arguments: []string{"key", "value"},
		},
		{
			LSN:       200,
			CommandID: compute.GetCommandID,
			Arguments: []string{"key"},
		},
		{
			LSN:       300,
			CommandID: compute.DelCommandID,
			Arguments: []string{"key"},
		},
	}

	var writeBuffer bytes.Buffer
	for idx := range expectedLogs {
		expectedLog := expectedLogs[idx]
		err := expectedLog.Encode(&writeBuffer)
		require.NoError(t, err)
	}

	data := writeBuffer.Bytes()
	readBuffer := bytes.NewBuffer(data)

	var logs []Log
	for readBuffer.Len() > 0 {
		var log Log
		err := log.Decode(readBuffer)
		require.NoError(t, err)
		logs = append(logs, log)
	}

	assert.True(t, reflect.DeepEqual(expectedLogs, logs))
}
