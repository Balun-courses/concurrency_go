package tools

type Promise[T any] struct {
	result   chan T
	promised bool
}

func NewPromise[T any]() Promise[T] {
	return Promise[T]{
		result: make(chan T, 1),
	}
}

func (p *Promise[T]) Set(value T) {
	if p.promised == true {
		return
	}

	p.promised = true
	p.result <- value
	close(p.result)
}

func (p *Promise[T]) GetFuture() Future[T] {
	return NewFuture[T](p.result)
}
