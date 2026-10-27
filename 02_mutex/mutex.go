package mutex

import (
	"primitives/internal/futex"
	"runtime"
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

	for range 4 {
		if atomic.LoadUint32(&m.state) == free && m.TryLock() {
			return
		}
		runtime.Gosched()
	}

	for atomic.SwapUint32(&m.state, contended) != free {
		futex.Wait(&m.state, contended)
	}
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapUint32(&m.state, free, held)
}

func (m *Mutex) Unlock() {
	old := atomic.SwapUint32(&m.state, free)
	if old == free {
		panic("mutex: unlock of unlocked mutex")
	}
	if old == contended {
		futex.Wake(&m.state)
	}
}
