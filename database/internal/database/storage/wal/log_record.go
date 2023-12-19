package wal

import "spider/internal/tools"

type Data struct {
	LSN       int64
	CommandID int
	Arguments []string
}

type LogRecord struct {
	data         Data
	writePromise tools.Promise[error]
}

func NewLogRecord(lsn int64, commandID int, args []string) LogRecord {
	return LogRecord{
		data: Data{
			LSN:       lsn,
			CommandID: commandID,
			Arguments: args,
		},
		writePromise: tools.NewPromise[error](),
	}
}

func (r *LogRecord) Data() Data {
	return r.data
}

func (r *LogRecord) LSN() int64 {
	return r.data.LSN
}

func (r *LogRecord) CommandID() int {
	return r.data.CommandID
}

func (r *LogRecord) Arguments() []string {
	return r.data.Arguments
}

func (r *LogRecord) SetResult(err error) {
	r.writePromise.Set(err)
}

func (r *LogRecord) Result() tools.Future[error] {
	return r.writePromise.GetFuture()
}
