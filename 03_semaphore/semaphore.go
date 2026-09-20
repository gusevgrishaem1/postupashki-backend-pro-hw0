package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Semaphore struct {
	permits uint32
}

func New(n int) *Semaphore {
	s := new(Semaphore)
	atomic.StoreUint32(&s.permits, uint32(n))
	return s
}

func (s *Semaphore) Acquire() {
	for {
		curr := atomic.LoadUint32(&s.permits)
		if curr == 0 {
			futex.Wait(&s.permits, 0)
			continue
		}

		next := curr - 1
		if atomic.CompareAndSwapUint32(&s.permits, curr, next) {
			break
		}
	}
}

func (s *Semaphore) TryAcquire() bool {
	curr := atomic.LoadUint32(&s.permits)
	if curr == 0 {
		return false
	}

	next := curr - 1
	return atomic.CompareAndSwapUint32(&s.permits, curr, next)
}

func (s *Semaphore) Release() {
	for {
		curr := atomic.LoadUint32(&s.permits)

		next := curr + 1
		if atomic.CompareAndSwapUint32(&s.permits, curr, next) {
			futex.Wake(&s.permits)
			break
		}
	}
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}
