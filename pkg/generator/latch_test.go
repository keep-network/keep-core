package generator

import (
	"testing"
)

func TestIsExecuting(t *testing.T) {
	latch := NewProtocolLatch()

	if latch.IsExecuting() {
		t.Errorf("protocol is not executing now")
	}

	latch.Lock()

	if !latch.IsExecuting() {
		t.Errorf("protocol is executing now")
	}

	latch.Unlock()

	if latch.IsExecuting() {
		t.Errorf("protocol is not executing now")
	}

	latch.Lock()
	latch.Lock()
	//lint:ignore SA2001 empty critical section for test purposes
	latch.Unlock()

	if !latch.IsExecuting() {
		t.Errorf("protocol is executing now")
	}

	latch.Unlock()

	if latch.IsExecuting() {
		t.Errorf("protocol is not executing now")
	}
}

// TestUnlockPanic ensures the Unlock() function panics if Lock() was called
// before.
func TestUnlockPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Unlock should panic")
		}
	}()

	latch := NewProtocolLatch()
	latch.Unlock()
}
