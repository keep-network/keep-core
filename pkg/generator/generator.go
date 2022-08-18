package generator

import (
	"time"

	"github.com/ipfs/go-log"
)

var logger = log.Logger("keep-generator")

// Protocol defines the interface that allows the StartScheduler function to
// determine if the protocol is executing or not. This interface should be
// implemented by all important protocols of the client, such as distributed
// key generation or signing.
type Protocol interface {
	isExecuting() bool
}

const checkTick = 1 * time.Second

// StartScheduler creates a new instance of a Scheduler that is responsible
// for managing long-running, computationally-expensive operations.
// The scheduler stops and resumes operations based on the state of provided
// protocols. If at least one of the protocols is currently executing, the
// scheduler stops all computations. Computations are automatically resumed once
// none of the protocols is executing.
func StartScheduler(protocols []Protocol) *Scheduler {
	scheduler := &Scheduler{}

	go func() {
		for {
			atLeastOneProtocolExecuting := false

			for _, protocol := range protocols {
				if protocol.isExecuting() {
					atLeastOneProtocolExecuting = true
					break
				}
			}

			if atLeastOneProtocolExecuting {
				scheduler.stop()
			} else {
				scheduler.resume()
			}

			time.Sleep(checkTick)
		}
	}()

	return scheduler
}
