package wal

import (
	"bytes"
	"encoding/gob"
)

type Log struct {
	LSN       int64
	CommandID int
	Arguments []string
}

func (l *Log) Encode(buffer *bytes.Buffer) error {
	encoder := gob.NewEncoder(buffer)
	return encoder.Encode(*l)
}

func (l *Log) Decode(buffer *bytes.Buffer) error {
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(l)
}
