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
		next := max(curr+delta, 0)
		if atomic.CompareAndSwapUint32(&wg.count, uint32(curr), uint32(next)) {
			if next == 0 {
				futex.WakeAll(&wg.count)
			}
			break
		}
	}
}

func (wg *WaitGroup) Done() {
	for {
		curr := atomic.LoadUint32(&wg.count)
		if curr == 0 {
			panic("Done")
		}

		next := curr - 1
		if atomic.CompareAndSwapUint32(&wg.count, curr, next) {
			if next == 0 {
				futex.WakeAll(&wg.count)
			}
			break
		}
	}

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
