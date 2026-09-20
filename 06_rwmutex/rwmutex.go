package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	writer uint32 = 1 << 31
	reader        = writer - 1
)

type RWMutex struct {
	state uint32
}

func (rw *RWMutex) RLock() {
	for {
		curr := atomic.LoadUint32(&rw.state)

		if curr&writer != 0 {
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
			panic("rwmutex.RUnlock")
		}

		next := curr - 1
		if atomic.CompareAndSwapUint32(&rw.state, curr, next) {
			futex.Wake(&rw.state)
			break
		}
	}
}

func (rw *RWMutex) Lock() {
	for {
		curr := atomic.LoadUint32(&rw.state)

		if curr != 0 {
			futex.Wait(&rw.state, curr)
			continue
		}

		if atomic.CompareAndSwapUint32(&rw.state, curr, curr|writer) {
			break
		}
	}
}

func (rw *RWMutex) Unlock() {
	for {
		curr := atomic.LoadUint32(&rw.state)

		if curr&writer == 0 {
			panic("rwmutex.Unlock")
		}

		if atomic.CompareAndSwapUint32(&rw.state, curr, curr&^writer) {
			futex.Wake(&rw.state)
			break
		}
	}
}
