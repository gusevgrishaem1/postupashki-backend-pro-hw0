package waitgroup

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type WaitGroup struct {
	count uint32
}

func (wg *WaitGroup) Add(delta int) {
	for {
		curr := int(atomic.LoadUint32(&wg.count))
		next := curr + delta
		if next < 0 {
			panic("waitgroup: negative counter")
		}
		if atomic.CompareAndSwapUint32(&wg.count, uint32(curr), uint32(next)) {
			if next == 0 {
				futex.WakeAll(&wg.count)
			}
			break
		}
	}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	for {
		curr := atomic.LoadUint32(&wg.count)
		if curr == 0 {
			break
		}
		futex.Wait(&wg.count, curr)
	}
}
