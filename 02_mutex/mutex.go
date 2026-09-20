package mutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	free = iota
	held
	contended
)

type Mutex struct {
	state uint32
}

func (m *Mutex) Lock() {
	if atomic.CompareAndSwapUint32(&m.state, free, held) {
		return
	}

	for {
		if atomic.CompareAndSwapUint32(&m.state, free, contended) {
			break
		}

		if atomic.LoadUint32(&m.state) == contended {
			futex.Wait(&m.state, contended)
			continue
		}

		if atomic.CompareAndSwapUint32(&m.state, held, contended) {
			futex.Wait(&m.state, contended)
		}
	}
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapUint32(&m.state, free, held)
}

func (m *Mutex) Unlock() {
	if atomic.CompareAndSwapUint32(&m.state, held, free) {
		return
	}

	if atomic.CompareAndSwapUint32(&m.state, contended, free) {
		futex.Wake(&m.state)
		return
	}

	panic("Unlock")
}
