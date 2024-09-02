package concurrency

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(ticketsNumber int) Semaphore {
	return Semaphore{
		tickets: make(chan struct{}, ticketsNumber),
	}
}

func (s *Semaphore) Acquire() {
	if s == nil || s.tickets == nil {
		return
	}

	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	if s == nil || s.tickets == nil {
		return
	}

	<-s.tickets
}

func (s *Semaphore) WithAcquire(action func()) {
	if s == nil || action == nil {
		return
	}

	s.Acquire()
	action()
	s.Release()
}
