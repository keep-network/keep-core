package promise

import (
	"fmt"
	"sync"
)

// Promise represents the eventual completion or failure of an
// ansynchronous operation and its resulting value. Promise can
// be either fulfilled or failed and it can happen only one time.
// All Promise operations are thread-safe, guarded by a mutex.
type Promise struct {
	mutex     sync.Mutex
	successFn func(interface{})
	failureFn func(error)

	isComplete bool
}

// NewPromise creates a new, uncompleted Promise instance with
// no success or failure callback configured.
func NewPromise() *Promise {
	return &Promise{
		isComplete: false,
	}
}

// OnSuccess registers a function to be called when the Promise
// has been fulfilled. In case of a failed Promise, function is not
// called at all. OnSuccess is a non-blocking operation. Only one on success
// function can be registered for a Promise.
func (p *Promise) OnSuccess(onSuccess func(interface{})) *Promise {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.successFn = onSuccess
	return p
}

// OnFailure registers a function to be called when the Promise
// execution failed. In case of a fulfilled Promise, function is not
// called at all. OnFailure is a non-blocking operation. Only one on failure
// function can be registered for a Promise.
func (p *Promise) OnFailure(onFailure func(error)) *Promise {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.failureFn = onFailure
	return p
}

// Fulfill can happen only once for a Promise and it results in calling
// the OnSuccess callback, if registered. If Promise has been already
// completed by either fulfilling or failing, this function reports
// an error.
func (p *Promise) Fulfill(value interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isComplete {
		return fmt.Errorf("promise already completed")
	}

	p.isComplete = true
	if p.successFn != nil {
		go func() {
			p.successFn(value)
		}()
	}

	return nil
}

// Fail can happen only once for a Promise and it results in calling
// the OnFailure callback, if registered. If Promise has been already
// completed by either fulfilling or failing, this function reports
// an error.
func (p *Promise) Fail(err error) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isComplete {
		return fmt.Errorf("promise already completed")
	}

	p.isComplete = true
	if p.failureFn != nil {
		go func() {
			p.failureFn(err)
		}()
	}

	return nil
}
