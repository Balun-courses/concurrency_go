package wal

import "spider/internal/concurrency"

type LogData struct {
	LSN       int64
	CommandID int
	Arguments []string
}

type Log struct {
	data         LogData
	writePromise concurrency.Promise[error]
}

func NewLog(lsn int64, commandID int, args []string) Log {
	return Log{
		data: LogData{
			LSN:       lsn,
			CommandID: commandID,
			Arguments: args,
		},
		writePromise: concurrency.NewPromise[error](),
	}
}

func (l *Log) Data() LogData {
	return l.data
}

func (l *Log) LSN() int64 {
	return l.data.LSN
}

func (l *Log) CommandID() int {
	return l.data.CommandID
}

func (l *Log) Arguments() []string {
	return l.data.Arguments
}

func (l *Log) SetResult(err error) {
	l.writePromise.Set(err)
}

func (l *Log) Result() concurrency.Future[error] {
	return l.writePromise.GetFuture()
}
