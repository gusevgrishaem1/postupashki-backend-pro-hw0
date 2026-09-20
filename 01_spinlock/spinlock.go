package spinlock

import (
	"runtime"
	"sync/atomic"
)

type Spinlock struct {
	locked atomic.Bool
}

func (s *Spinlock) Lock() {
	for !s.locked.CompareAndSwap(false, true) {
		runtime.Gosched()
	}
}

func (s *Spinlock) TryLock() bool {
	return s.locked.CompareAndSwap(false, true)
}

func (s *Spinlock) Unlock() {
	if !s.locked.CompareAndSwap(true, false) {
		panic("Unlock")
	}
}

type TTAS struct {
	locked atomic.Bool
}

func (s *TTAS) Lock() {
	for {
		if s.locked.Load() {
			runtime.Gosched()
			continue
		}
		if s.locked.CompareAndSwap(false, true) {
			break
		}
		runtime.Gosched()
	}
}

func (s *TTAS) TryLock() bool {
	return s.locked.CompareAndSwap(false, true)
}

func (s *TTAS) Unlock() {
	if !s.locked.CompareAndSwap(true, false) {
		panic("Unlock")
	}
}
