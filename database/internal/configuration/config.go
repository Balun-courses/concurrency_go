package configuration

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Engine      *EngineConfig      `yaml:"engine"`
	WAL         *WALConfig         `yaml:"wal"`
	Replication *ReplicationConfig `yaml:"replication"`
	Network     *NetworkConfig     `yaml:"network"`
	Logging     *LoggingConfig     `yaml:"logging"`
}

type EngineConfig struct {
	Type string `yaml:"type"`
}

type WALConfig struct {
	FlushingBatchLength  int           `yaml:"flushing_batch_length"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string        `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

type ReplicationConfig struct {
	ReplicaType   string        `yaml:"replica_type"`
	MasterAddress string        `yaml:"master_address"`
	SyncInterval  time.Duration `yaml:"sync_interval"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func Load(filename string) (*Config, error) {
	if filename == "" {
		return &Config{}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
