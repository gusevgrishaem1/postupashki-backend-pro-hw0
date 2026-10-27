package spinlock

import (
	"runtime"
	"sync/atomic"
)

type lockState struct {
	locked atomic.Bool
}

func (s *lockState) TryLock() bool {
	return s.locked.CompareAndSwap(false, true)
}

func (s *lockState) Unlock() {
	if !s.locked.CompareAndSwap(true, false) {
		panic("spinlock: unlock of unlocked lock")
	}
}

type Spinlock struct {
	lockState
}

func (s *Spinlock) Lock() {
	for !s.TryLock() {
		runtime.Gosched()
	}
}

type TTAS struct {
	lockState
}

func (s *TTAS) Lock() {
	for {
		if s.locked.Load() {
			runtime.Gosched()
			continue
		}
		if s.TryLock() {
			break
		}
		runtime.Gosched()
	}
}
