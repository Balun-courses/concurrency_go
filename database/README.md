# spider

### Introduction
Spider - is in-memory key-value database with support for asynchronous replication (physic) and WAL (Write-Ahead Logging). Implented only like reference of the final project for course Concurrency in Go.

### Server launch

For local start of the server, you can use the following instructions:

```bash
make run-server
```

### Server configuration

By default you don't need to use any parameters for start database, but implemented support of different configation paramaters with YAML format, for example:

```yaml
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
```

For launch with configuration you need to use the following instructions:

```bash
make run-server CONFIG_FILE_NAME=some_address
```

### CLI launch

For local start of the CLI, you can use the following instructions:

```bash
make run-cli
```

