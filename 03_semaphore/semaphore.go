package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Semaphore struct {
	permits uint32
	waiters int32
}

func New(n int) *Semaphore {
	if n < 0 {
		panic("semaphore: negative permit count")
	}
	return &Semaphore{permits: uint32(n)}
}

func (s *Semaphore) Acquire() {
	for !s.TryAcquire() {
		atomic.AddInt32(&s.waiters, 1)
		futex.Wait(&s.permits, 0)
		atomic.AddInt32(&s.waiters, -1)
	}
}

func (s *Semaphore) TryAcquire() bool {
	for {
		curr := atomic.LoadUint32(&s.permits)
		if curr == 0 {
			return false
		}
		if atomic.CompareAndSwapUint32(&s.permits, curr, curr-1) {
			return true
		}
	}
}

func (s *Semaphore) Release() {
	atomic.AddUint32(&s.permits, 1)
	if atomic.LoadInt32(&s.waiters) != 0 {
		futex.Wake(&s.permits)
	}
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}
