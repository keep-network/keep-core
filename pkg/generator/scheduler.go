package generator

import (
	"context"
	"sync"
)

type state int

const (
	working state = iota
	stopped
)

// Scheduler allows managing computationally heavy operations: stopping and
// resuming them. The client needs to generate parameters for cryptographic
// algorithms and generating these parameters requires a lot of CPU cycles.
// The generation process may starve other processes in the client when it comes
// to access to the CPU. The scheduler allows starting the parameter generation
// when no other processes such as key generation or signing are executed on the
// client. This way, the client that would normally be idle, can spend CPU
// cycles on computationally heavy operations and stop these operations when CPU
// cycles are needed elsewhere.
type Scheduler struct {
	state   state
	workers []func(context.Context)
	stops   []context.CancelFunc

	mutex sync.Mutex
}

// Compute takes the worker function and starts the computations in a separate
// goroutine if the scheduler status is "working". Otherwise, when the scheduler
// status is "stopped", the worker function is scheduled for execution later.
// The function accepts the context and is required to stop the execution if
// the context is done. The function will be called in a loop until the
// scheduler is stopped.
func (s *Scheduler) compute(workerFn func(context.Context)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.workers = append(s.workers, workerFn)

	if s.state == working {
		s.startWorker(workerFn)
	}
}

// Stop asks all worker functions to stop their work. The context passed to
// the function is cancelled and no further calls to the worker function are
// done until the scheduler work is resumed.
func (s *Scheduler) stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.state == stopped {
		return
	}

	logger.Info("stopping computations\n")
	s.state = stopped

	for _, stop := range s.stops {
		stop()
	}
	s.stops = nil
}

// Resume resumes the work of all worker functions, each in a separate
// goroutine.
func (s *Scheduler) resume() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.state == working {
		return
	}

	logger.Info("resuming computations\n")
	s.state = working

	for _, worker := range s.workers {
		s.startWorker(worker)
	}
}

func (s *Scheduler) startWorker(miningFn func(context.Context)) {
	ctx, cancelFn := context.WithCancel(context.Background())
	s.stops = append(s.stops, cancelFn)

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
