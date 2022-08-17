package miner

import (
	"sync"

	"github.com/ipfs/go-log"
)

var logger = log.Logger("keep-miner")

type state int

const (
	working state = iota
	stopped
)

// Miner allows managing computationally heavy operations: stopping and resuming
// them. The client needs to generate parameters for cryptographic algorithms
// and generating these parameters requires a lot of CPU cycles. The generation
// process may starve other processes in the client when it comes to access to
// the CPU. The miner allows starting the parameter generation when no other
// processes such as key generation or signing are executed on the client. This
// way, the client that would normally be idle, can spend CPU cycles on
// computationally heavy operations and stop these operations when CPU cycles
// are needed elsewhere.
type Miner struct {
	state   state
	workers []func()
	stops   []chan interface{}

	mutex sync.Mutex
}

// Mine takes the worker function and starts the computations in a separate
// goroutine if the miner status is "working". Otherwise, when the miner status
// is "stopped", the worker function is scheduled for execution later.
func (m *Miner) Mine(miningFn func()) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.workers = append(m.workers, miningFn)

	if m.state == working {
		m.startWorker(miningFn)
	}
}

// Stop asks all worker functions to stop their work. Note that there is no
// guarantee the worker function will stop immediately but the miner will not
// call the worker function again until the miner's work is resumed.
func (m *Miner) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	logger.Info("stopping computations")
	m.state = stopped

	for _, stop := range m.stops {
		stop <- struct{}{}
	}
	m.stops = nil
}

// Resume resumes the work of all worker functions, each in a separate
// goroutine.
func (m *Miner) Resume() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	logger.Info("resuming computations")
	m.state = working

	for _, worker := range m.workers {
		m.startWorker(worker)
	}
}

func (m *Miner) startWorker(miningFn func()) {
	stopChan := make(chan interface{})
	m.stops = append(m.stops, stopChan)

	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				miningFn()
			}
		}
	}()
}
