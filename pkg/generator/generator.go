package generator

import (
	"time"

	"github.com/ipfs/go-log"
)

var logger = log.Logger("keep-generator")

const checkTick = 1 * time.Second

// StartScheduler creates a new instance of a Scheduler that is responsible
// for managing long-running, computationally-expensive operations.
// The scheduler stops and resumes operations based on the state of registered
// protocols. If at least one of the protocols is currently executing, the
// scheduler stops all computations. Computations are automatically resumed once
// none of the protocols is executing.
func StartScheduler() *Scheduler {
	scheduler := &Scheduler{}

	go func() {
		for {
			scheduler.checkProtocols()
			time.Sleep(checkTick)
		}
	}()

	return scheduler
}
