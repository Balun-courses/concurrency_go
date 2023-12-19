package wal

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"log"
	"os"
	"spider/pkg/processing"
	"testing"
)

const testWALDirectory = "./test_wal"

func TestMain(m *testing.M) {
	if err := os.Mkdir(testWALDirectory, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := os.RemoveAll(testWALDirectory); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestBatchWritingToWALSegment(t *testing.T) {
	t.Parallel()

	maxSegmentSize := 100 << 10
	fsWriter := NewFSWriter(testWALDirectory, maxSegmentSize, zap.NewNop())

	batch := []LogRecord{
		NewLogRecord(1, processing.SetCommandID, []string{"key_1", "value_1"}),
		NewLogRecord(2, processing.SetCommandID, []string{"key_2", "value_2"}),
		NewLogRecord(3, processing.SetCommandID, []string{"key_2", "value_2"}),
	}

	fsWriter.WriteBatch(batch)
	for _, record := range batch {
		err := record.Result()
		require.NoError(t, err.Get())
	}

	stat, err := os.Stat(testWALDirectory + "/0.wal")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())
}

func TestWALSegmentsRotation(t *testing.T) {
	t.Parallel()

	maxSegmentSize := 10
	fsWriter := NewFSWriter(testWALDirectory, maxSegmentSize, zap.NewNop())
	fsWriter.lastLSN = 4

	batch := []LogRecord{
		NewLogRecord(5, processing.SetCommandID, []string{"key_1", "value_1"}),
		NewLogRecord(6, processing.SetCommandID, []string{"key_2", "value_2"}),
		NewLogRecord(7, processing.SetCommandID, []string{"key_2", "value_2"}),
	}

	fsWriter.WriteBatch(batch)
	for _, record := range batch {
		err := record.Result()
		require.NoError(t, err.Get())
	}

	stat, err := os.Stat(testWALDirectory + "/4.wal")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())

	stat, err = os.Stat(testWALDirectory + "/5.wal")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())

	stat, err = os.Stat(testWALDirectory + "/6.wal")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())
}
