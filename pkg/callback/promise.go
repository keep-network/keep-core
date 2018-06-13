package promise

import (
	"fmt"
)

// Promise represents the eventual completion or failure of an
// ansynchronous operation and its resulting value. Promise can
// be either fulfilled or failed and it can happen only one time.
type Promise struct {
	successFn func(interface{})
	failureFn func(error)

	isCompleted bool
}

// NewPromise creates a new, uncompleted Promise instance with
// no success or failure callback configured. You need to install
// those callbacks before you fail or fulfill the Promise.
func NewPromise() *Promise {
	return &Promise{
		isCompleted: false,
	}
}

// OnSuccess registers an onSuccess callback that is called when the Promise
// has been fulfilled. In case of a failed Promise, onSucess callback is not
// called at all. OnSuccess is a non-blocking operation. Only one onSuccess
// callback can be registered for a Promise.
func (p *Promise) OnSuccess(onSuccess func(interface{})) *Promise {
	p.successFn = onSuccess
	return p
}

// OnFailure registers an onFailure callback that is called when the Promise
// execution failed. In case of a fulfilled Promise, onFailure callback is not
// called at all. OnFailure is a non-blocking operation. Only one onFailure
// callback can be registered for a Promise.
func (p *Promise) OnFailure(onFailure func(error)) *Promise {
	p.failureFn = onFailure
	return p
}

// Fulfill can happen only once for a Promise and it results in calling
// the OnSuccess callback. If Promise has been already completed by either
// fulfilling or failing, this function reports an error. If no OnSuccess
// callback has been registered for a Promise, error is reported as well.
func (p *Promise) Fulfill(value interface{}) error {
	if p.isCompleted {
		return fmt.Errorf("promise already completed")
	}

	if p.successFn == nil {
		return fmt.Errorf("success callback not registered")
	}

	p.isCompleted = true
	go func() {
		p.successFn(value)
	}()
	return nil
}

// Fail can happen only once for a Promise and it results in calling
// the OnFailure callback. If Promise has been already completed by either
// fulfilling or failing, this function reports an error. If no OnFailure
// callback has been registered for a Promise, error is reported as well.
func (p *Promise) Fail(err error) error {
	if p.isCompleted {
		return fmt.Errorf("promise already completed")
	}

	if p.failureFn == nil {
		return fmt.Errorf("failure callback not registered")
	}

	p.isCompleted = true
	go func() {
		p.failureFn(err)
	}()
	return nil
}
