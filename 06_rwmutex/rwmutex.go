package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	writer        uint32 = 1 << 31
	writerWaiting uint32 = 1 << 30
	reader               = writerWaiting - 1
)

type RWMutex struct {
	state uint32
}

func (rw *RWMutex) RLock() {
	for {
		curr := atomic.LoadUint32(&rw.state)

		if curr&(writer|writerWaiting) != 0 {
			futex.Wait(&rw.state, curr)
			continue
		}

		next := curr + 1
		if atomic.CompareAndSwapUint32(&rw.state, curr, next) {
			break
		}
	}
}

func (rw *RWMutex) RUnlock() {
	for {
		curr := atomic.LoadUint32(&rw.state)

		if curr&reader == 0 {
			panic("rwmutex: runlock of unlocked read lock")
		}

		next := curr - 1
		if atomic.CompareAndSwapUint32(&rw.state, curr, next) {
			if next&reader == 0 {
				futex.WakeAll(&rw.state)
			}
			break
		}
	}
}

func (rw *RWMutex) Lock() {
	if atomic.CompareAndSwapUint32(&rw.state, 0, writer) {
		return
	}
	for {
		curr := atomic.LoadUint32(&rw.state)

		if curr&(writer|writerWaiting) != 0 {
			futex.Wait(&rw.state, curr)
			continue
		}

		if atomic.CompareAndSwapUint32(&rw.state, curr, curr|writerWaiting) {
			break
		}
	}

	for {
		curr := atomic.LoadUint32(&rw.state)
		if curr == writerWaiting {
			if atomic.CompareAndSwapUint32(&rw.state, writerWaiting, writer) {
				return
			}
			continue
		}
		futex.Wait(&rw.state, curr)
	}
}

func (rw *RWMutex) Unlock() {
	for {
		curr := atomic.LoadUint32(&rw.state)

		if curr&writer == 0 {
			panic("rwmutex: unlock of unlocked write lock")
		}

		if atomic.CompareAndSwapUint32(&rw.state, curr, curr&^writer) {
			futex.WakeAll(&rw.state)
			break
		}
	}
}
