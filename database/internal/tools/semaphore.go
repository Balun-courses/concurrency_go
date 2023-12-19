package tools

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(ticketsNumber int) Semaphore {
	return Semaphore{
		tickets: make(chan struct{}, ticketsNumber),
	}
}

func (s *Semaphore) Acquire() {
	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tickets
}

func (s *Semaphore) WithSemaphore(action func()) {
	if action == nil {
		return
	}

	s.Acquire()
	action()
	s.Release()
}
