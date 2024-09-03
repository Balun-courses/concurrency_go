package wal

import (
	"bytes"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"spider/internal/database/compute"
)

// mockgen -source=logs_writer.go -destination=logs_writer_mock.go -package=wal

func TestNewLogsWriter(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		segment segment
		logger  *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create logs writer without segment": {
			expectedErr:    errors.New("segment is invalid"),
			expectedNilObj: true,
		},
		"create logs writer without logger": {
			segment:        NewMockwalSegment(ctrl),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create logs writer": {
			segment: NewMockwalSegment(ctrl),
			logger:  zap.NewNop(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			writer, err := NewLogsWriter(test.segment, test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, writer)
			} else {
				assert.NotNil(t, writer)
			}
		})
	}
}

func TestWriteWithErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("write error")
	requests := []WriteRequest{
		NewWriteRequest(100, compute.SetCommandID, []string{"key", "value"}),
		NewWriteRequest(200, compute.GetCommandID, []string{"key"}),
		NewWriteRequest(300, compute.DelCommandID, []string{"key"}),
	}

	var buffer bytes.Buffer
	for idx := range requests {
		log := requests[idx].log
		err := log.Encode(&buffer)
		require.NoError(t, err)
	}

	ctrl := gomock.NewController(t)
	segment := NewMockwalSegment(ctrl)
	segment.EXPECT().
		Write(buffer.Bytes()).
		Return(expectedErr)

	writer, err := NewLogsWriter(segment, zap.NewNop())
	require.NoError(t, err)
	writer.Write(requests)

	for _, request := range requests {
		futureResponse := request.FutureResponse()
		assert.Equal(t, expectedErr, futureResponse.Get())
	}
}

func TestWrite(t *testing.T) {
	t.Parallel()

	requests := []WriteRequest{
		NewWriteRequest(100, compute.SetCommandID, []string{"key", "value"}),
		NewWriteRequest(200, compute.GetCommandID, []string{"key"}),
		NewWriteRequest(300, compute.DelCommandID, []string{"key"}),
	}

	var buffer bytes.Buffer
	for idx := range requests {
		log := requests[idx].log
		err := log.Encode(&buffer)
		require.NoError(t, err)
	}

	ctrl := gomock.NewController(t)
	segment := NewMockwalSegment(ctrl)
	segment.EXPECT().
		Write(buffer.Bytes()).
		Return(nil)

	writer, err := NewLogsWriter(segment, zap.NewNop())
	require.NoError(t, err)
	writer.Write(requests)

	for _, request := range requests {
		futureResponse := request.FutureResponse()
		assert.Nil(t, futureResponse.Get())
	}
}
