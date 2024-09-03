package filesystem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSegmentUpperBound(t *testing.T) {
	t.Parallel()

	filename, err := SegmentUpperBound("test_data", "wal_0.log")
	require.NoError(t, err)
	require.Equal(t, "wal_1000.log", filename)

	filename, err = SegmentUpperBound("test_data", "wal_1000.log")
	require.NoError(t, err)
	require.Equal(t, "wal_2000.log", filename)

	filename, err = SegmentUpperBound("test_data", "wal_2000.log")
	require.NoError(t, err)
	require.Equal(t, "wal_3000.log", filename)

	filename, err = SegmentUpperBound("test_data", "wal_3000.log")
	require.NoError(t, err)
	require.Equal(t, "", filename)
}

func TestSegmentLast(t *testing.T) {
	t.Parallel()

	filename, err := SegmentLast("test_data")
	require.NoError(t, err)
	require.Equal(t, "wal_3000.log", filename)
}

/*
const testWALDirectory = "temp_test_data"

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
	maxSegmentSize := 100 << 10
	fsWriter := NewFSWriter(testWALDirectory, maxSegmentSize, zap.NewNop())

	batch := []Log{
		NewLog(1, compute.SetCommandID, []string{"key_1", "value_1"}),
		NewLog(2, compute.SetCommandID, []string{"key_2", "value_2"}),
		NewLog(3, compute.SetCommandID, []string{"key_3", "value_3"}),
	}

	now = func() time.Time {
		return time.Unix(1, 0)
	}

	fsWriter.WriteBatch(batch)
	for _, record := range batch {
		err := record.Result()
		require.NoError(t, err.Get())
	}

	stat, err := os.Stat(testWALDirectory + "/wal_1000.log")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())
}

func TestWALSegmentsRotation(t *testing.T) {
	maxSegmentSize := 10
	fsWriter := NewFSWriter(testWALDirectory, maxSegmentSize, zap.NewNop())

	batch := []Log{
		NewLog(4, compute.SetCommandID, []string{"key_4", "value_4"}),
		NewLog(5, compute.SetCommandID, []string{"key_5", "value_5"}),
		NewLog(6, compute.SetCommandID, []string{"key_6", "value_6"}),
	}

	now = func() time.Time {
		return time.Unix(2, 0)
	}

	fsWriter.WriteBatch(batch)
	for _, record := range batch {
		err := record.Result()
		require.NoError(t, err.Get())
	}

	batch = []Log{
		NewLog(7, compute.SetCommandID, []string{"key_7", "value_7"}),
		NewLog(8, compute.SetCommandID, []string{"key_8", "value_8"}),
		NewLog(9, compute.SetCommandID, []string{"key_9", "value_9"}),
	}

	now = func() time.Time {
		return time.Unix(3, 0)
	}

	fsWriter.WriteBatch(batch)
	for _, record := range batch {
		err := record.Result()
		require.NoError(t, err.Get())
	}

	stat, err := os.Stat(testWALDirectory + "/wal_2000.log")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())

	stat, err = os.Stat(testWALDirectory + "/wal_3000.log")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())
}
*/

/*func TestReadLogs(t *testing.T) {
	t.Parallel()

	reader := NewFSReader("test_data", zap.NewNop())

	logs, err := reader.ReadLogs()
	require.NoError(t, err)
	require.Equal(t, 9, len(logs))

	// from tests_data/wal_1000.log
	require.Equal(t, LogData{LSN: 1, CommandID: 1, Arguments: []string{"key_1", "value_1"}}, logs[0])
	require.Equal(t, LogData{LSN: 2, CommandID: 1, Arguments: []string{"key_2", "value_2"}}, logs[1])
	require.Equal(t, LogData{LSN: 3, CommandID: 1, Arguments: []string{"key_3", "value_3"}}, logs[2])

	// from tests_data/wal_2000.log
	require.Equal(t, LogData{LSN: 4, CommandID: 1, Arguments: []string{"key_4", "value_4"}}, logs[3])
	require.Equal(t, LogData{LSN: 5, CommandID: 1, Arguments: []string{"key_5", "value_5"}}, logs[4])
	require.Equal(t, LogData{LSN: 6, CommandID: 1, Arguments: []string{"key_6", "value_6"}}, logs[5])

	// from tests_data/wal_3000.log
	require.Equal(t, LogData{LSN: 7, CommandID: 1, Arguments: []string{"key_7", "value_7"}}, logs[6])
	require.Equal(t, LogData{LSN: 8, CommandID: 1, Arguments: []string{"key_8", "value_8"}}, logs[7])
	require.Equal(t, LogData{LSN: 9, CommandID: 1, Arguments: []string{"key_9", "value_9"}}, logs[8])
}*/
