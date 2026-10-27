package once

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	open = iota
	processing
	ready
)

type Once struct {
	state uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.state) == ready {
		return
	}

	if atomic.CompareAndSwapUint32(&o.state, open, processing) {
		defer func() {
			atomic.StoreUint32(&o.state, ready)
			futex.WakeAll(&o.state)
		}()
		f()
		return
	}

	for {
		if atomic.LoadUint32(&o.state) == ready {
			return
		}
		futex.Wait(&o.state, processing)
	}
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.state) == ready
}
