package promise

import (
	"fmt"
)

// Promise represents the eventual completion (or failure) of an a
// synchronous operation, and its resulting value
type Promise struct {
	successFn func(interface{})
	failureFn func(error)

	isCompleted bool
}

func newPromise() *Promise {
	return &Promise{
		isCompleted: false,
	}
}

// onSuccess registers an onComplete callback that is called when the Promise
// execution has completed successfuly. In case of a failure, onSuccess
// callback is not called at all. onSuccess is a non-blocking operation.
func (p *Promise) onSuccess(onSuccess func(interface{})) *Promise {
	p.successFn = onSuccess
	return p
}

// onFailure registers an onFailure callback that is called when the Promise
// execution failed at any point. In case of a successful Promise execution
// onFailure callback is not called at all. onFailure is a non-blocking operation.
func (p *Promise) onFailure(onFailure func(error)) *Promise {
	p.failureFn = onFailure
	return p
}

func (p *Promise) fulfill(value interface{}) error {
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

func (p *Promise) fail(err error) error {
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
