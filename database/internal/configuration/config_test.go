package configuration

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testCfgData = `
engine:
  type: "in_memory"
  partitions_number: 8
wal:
  flushing_batch_length: 100
  flushing_batch_timeout: "10ms"
  max_segment_size: "10MB"
  data_directory: "/data/spider/wal"
replication:
  replica_type: "slave"
  master_address: "127.0.0.1:3232"
  sync_interval: "1s"
network:
  address: "127.0.0.1:3223"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: 5m
logging:
  level: "info"
  output: "/log/output.log"
`

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfgData string

		expectedCfg Config
	}{
		"load empty config": {
			cfgData: ``,
		},
		"load config": {
			cfgData: testCfgData,
			expectedCfg: Config{
				Engine: &EngineConfig{
					Type:             "in_memory",
					PartitionsNumber: 8,
				},
				WAL: &WALConfig{
					FlushingBatchLength:  100,
					FlushingBatchTimeout: 10 * time.Millisecond,
					MaxSegmentSize:       "10MB",
					DataDirectory:        "/data/spider/wal",
				},
				Replication: &ReplicationConfig{
					ReplicaType:   "slave",
					MasterAddress: "127.0.0.1:3232",
					SyncInterval:  time.Second,
				},
				Network: &NetworkConfig{
					Address:        "127.0.0.1:3223",
					MaxConnections: 100,
					MaxMessageSize: "4KB",
					IdleTimeout:    time.Minute * 5,
				},
				Logging: &LoggingConfig{
					Level:  "info",
					Output: "/log/output.log",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reader := strings.NewReader(test.cfgData)
			cfg, err := Load(reader)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCfg, *cfg)
		})
	}
}
