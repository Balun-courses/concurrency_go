package once_implementation

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	mutex sync.Mutex
	state uintptr
}

func NewOnce() *Once {
	return &Once{}
}

func (o *Once) Do(action func()) {
	if action == nil {
		return
	}

	if atomic.LoadUintptr(&o.state) == 0 {
		o.doOnce(action)
	}
}

func (o *Once) doOnce(action func()) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.state == 0 {
		action()
		atomic.StoreUintptr(&o.state, 1)
	}
}
