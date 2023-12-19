package tools

type FutureError = Future[error]

type Future[T any] struct {
	result <-chan T
}

func NewFuture[T any](result <-chan T) Future[T] {
	return Future[T]{
		result: result,
	}
}

func (f *Future[T]) Get() T {
	return <-f.result
}
