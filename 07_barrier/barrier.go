package barrier

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Barrier struct {
	need    uint32
	arrived uint32
	round   uint32
}

func New(n int) *Barrier {
	b := new(Barrier)
	b.need = uint32(n)
	return b
}

func (b *Barrier) Wait() {
	for {
		curr := atomic.LoadUint32(&b.arrived)

		next := curr + 1
		if atomic.CompareAndSwapUint32(&b.arrived, curr, next) {
			if next == b.need {
				atomic.AddUint32(&b.round, 1)
				atomic.StoreUint32(&b.arrived, 0)
				futex.WakeAll(&b.round)
				return
			}

			futex.Wait(&b.round, b.round)
			break
		}
	}
}
