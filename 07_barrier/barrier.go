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
	round := atomic.LoadUint32(&b.round)

	for {
		curr := atomic.LoadUint32(&b.arrived)

		next := curr + 1
		if atomic.CompareAndSwapUint32(&b.arrived, curr, next) {
			if next == b.need {
				atomic.StoreUint32(&b.arrived, 0)
				atomic.AddUint32(&b.round, 1)
				futex.WakeAll(&b.round)
				break
			}

			for atomic.LoadUint32(&b.round) == round {
				futex.Wait(&b.round, round)
			}

			break
		}
	}
}
