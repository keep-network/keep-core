// Package async code is auto generated.
package async

import (
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon/event"
)

// EventRequestPromise represents an eventual completion of an ansynchronous
// operation and its resulting value. Promise can be either fulfilled or failed
// nd it can happen only one time. All Promise operations are thread-safe.
//
// To create a promise use: `&EventRequestPromise{}`
type EventRequestPromise struct {
	mutex      sync.Mutex
	successFn  func(*event.Request)
	failureFn  func(error)
	completeFn func(*event.Request, error)

	isComplete bool
	value      *event.Request
	err        error
}

// OnSuccess registers a function to be called when the Promise
// has been fulfilled. In case of a failed Promise, function is not
// called at all. OnSuccess is a non-blocking operation. Only one on success
// function can be registered for a Promise. If the Promise has been already
// fulfilled, the function is called immediatelly.
func (p *EventRequestPromise) OnSuccess(onSuccess func(*event.Request)) *EventRequestPromise {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.successFn = onSuccess

	if p.isComplete && p.err == nil {
		p.callSuccessFn()
	}

	return p
}

// OnFailure registers a function to be called when the Promise
// execution failed. In case of a fulfilled Promise, function is not
// called at all. OnFailure is a non-blocking operation. Only one on failure
// function can be registered for a Promise. If the Promise has already failed,
// the function is called immediatelly.
func (p *EventRequestPromise) OnFailure(onFailure func(error)) *EventRequestPromise {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.failureFn = onFailure

	if p.isComplete && p.err != nil {
		p.callFailureFn()
	}

	return p
}

// OnComplete registers a function to be called when the Promise
// execution completed no matter if it succeded or failed.
// In case of a successful execution, error passed to the callback
// function is nil. In case of a failed execution, there is no
// value evaluated so the value parameter is nil. OnComplete is
// a non-blocking operation. Only one on complete function can be
// registered for a Promise. If the Promise has already completed,
// the function is called immediatelly.
func (p *EventRequestPromise) OnComplete(onComplete func(*event.Request, error)) *EventRequestPromise {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.completeFn = onComplete

	if p.isComplete {
		p.callCompleteFn()
	}

	return p
}

// Fulfill can happen only once for a Promise and it results in calling
// the OnSuccess callback, if registered. If Promise has been already
// completed by either fulfilling or failing, this function reports
// an error.
func (p *EventRequestPromise) Fulfill(value *event.Request) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isComplete {
		return fmt.Errorf("promise already completed")
	}

	p.isComplete = true
	p.value = value

	p.callSuccessFn()
	p.callCompleteFn()

	return nil
}

// Fail can happen only once for a Promise and it results in calling
// the OnFailure callback, if registered. If Promise has been already
// completed by either fulfilling or failing, this function reports
// an error. Also, this function reports an error if `err` parameter
// is `nil`.
func (p *EventRequestPromise) Fail(err error) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if err == nil {
		return fmt.Errorf("error cannot be nil")
	}

	if p.isComplete {
		return fmt.Errorf("promise already completed")
	}

	p.isComplete = true
	p.err = err

	p.callFailureFn()
	p.callCompleteFn()

	return nil
}

func (p *EventRequestPromise) callCompleteFn() {
	if p.completeFn != nil {
		go func() {
			p.completeFn(p.value, p.err)
		}()
	}
}

func (p *EventRequestPromise) callSuccessFn() {
	if p.successFn != nil {
		go func() {
			p.successFn(p.value)
		}()
	}
}

func (p *EventRequestPromise) callFailureFn() {
	if p.failureFn != nil {
		go func() {
			p.failureFn(p.err)
		}()
	}
}
