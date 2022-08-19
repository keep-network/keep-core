package miner

import (
	"context"
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
//
// There are two requirements for goroutines using the miner:
// 1. Always call Stop() before Resume(),
// 2. Never forget to call Resume().
type Miner struct {
	state   state
	latch   sync.WaitGroup
	workers []func(context.Context)
	stops   []context.CancelFunc

	mutex sync.Mutex
}

// Mine takes the worker function and starts the computations in a separate
// goroutine if the miner status is "working". Otherwise, when the miner status
// is "stopped", the worker function is scheduled for execution later.
// The function accepts the context and is required to stop the execution if
// the context is done. The function will be called in a loop until the
// miner is stopped.
func (m *Miner) Mine(miningFn func(context.Context)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.workers = append(m.workers, miningFn)

	if m.state == working {
		m.startWorker(miningFn)
	}
}

// Stop asks all worker functions to stop their work. The context passed to
// the function is cancelled and no further calls to the function are done
// until the miner work is resumed.
// Stop can be called multiple times. Each call increases a count on the
// internal latch by one. To resume the work, Resume() function needs to be as
// many times as Stop(). This way, several goroutines can stop the computations
// and they all need to agree on resuming them.
//
// Requirement:
// Each goroutine calling Stop() must call Resume() once it is done.
func (m *Miner) Stop() {
	m.latch.Add(1)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.state == stopped {
		return
	}

	logger.Info("stopping computations\n")
	m.state = stopped

	for _, stop := range m.stops {
		stop()
	}
	m.stops = nil
}

// Resume resumes the work of all worker functions, each in a separate
// goroutine. If Stop() has been called multiple times, Resume() needs to be
// called the same number of times. Resume() blocks until the work can be
// resumed. The function panics if all requirements are not met.
//
// Requirements:
// Each goroutine calling Resume() must have had Stop() called before.
// Each goroutine must call Stop() and Resume() only one time.
func (m *Miner) Resume() {
	m.latch.Done()
	m.latch.Wait()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.state == working {
		return
	}

	logger.Info("resuming computations\n")
	m.state = working

	for _, worker := range m.workers {
		m.startWorker(worker)
	}
}

func (m *Miner) startWorker(miningFn func(context.Context)) {
	ctx, cancelFn := context.WithCancel(context.Background())
	m.stops = append(m.stops, cancelFn)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				miningFn(ctx)
			}
		}
	}()
}
