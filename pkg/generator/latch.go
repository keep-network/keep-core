package generator

import "sync"

// ProtocolLatch increases the internal counter every time protocol execution
// starts and decreases the counter every time protocol execution completes.
// The latch implements Protocol interface and can be registered in the
// Scheduler.
//
// The protocol code using the latch must guarantee that:
// 1. `Lock()` is always called before `Unlock()`
// 2. `Unlock()` is eventually called for every `Lock()`.
//
// Note that the Unlock() function may panic if the conditions are not met.
type ProtocolLatch struct {
	counter uint64
	mutex   sync.RWMutex
}

func NewProtocolLatch() *ProtocolLatch {
	return &ProtocolLatch{}
}

func (pl *ProtocolLatch) Lock() {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	pl.counter++
}

func (pl *ProtocolLatch) Unlock() {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	if pl.counter == 0 {
		panic("Lock was not called before Unlock")
	}

	pl.counter--
}

func (pl *ProtocolLatch) IsExecuting() bool {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()

	return pl.counter != 0
}
