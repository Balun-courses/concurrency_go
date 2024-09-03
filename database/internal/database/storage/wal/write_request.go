package wal

import "spider/internal/concurrency"

type WriteRequest struct {
	log     Log
	promise concurrency.PromiseError
}

func NewWriteRequest(lsn int64, commandID int, args []string) WriteRequest {
	return WriteRequest{
		log: Log{
			LSN:       lsn,
			CommandID: commandID,
			Arguments: args,
		},
		promise: concurrency.NewPromise[error](),
	}
}

func (l *WriteRequest) Log() Log {
	return l.log
}

func (l *WriteRequest) SetResponse(err error) {
	l.promise.Set(err)
}

func (l *WriteRequest) FutureResponse() concurrency.FutureError {
	return l.promise.GetFuture()
}
