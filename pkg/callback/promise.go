package promise

// Promise represents the eventual completion (or failure) of an a
// synchronous operation, and its resulting value
type Promise struct {
	successFn func(interface{})
	failureFn func(interface{})

	result  interface{}
	failure error
}

// OnSuccess registers an onComplete callback that is called when the Promise
// execution has completed successfuly. In case of a failure, onSuccess
// callback is not called at all. OnSuccess is a non-blocking operation.
func (p *Promise) OnSuccess(onSuccess func(interface{})) {
	p.successFn = onSuccess

	// TODO: Dirty hack, just for now; if result is ready, call the function
	if p.result != nil {
		onSuccess(p.result)
	}
}

// OnFailure registers an onFailure callback that is called when the Promise
// execution failed at any point. In case of a successful Promise execution
// onFailure callback is not called at all. OnFailure is a non-blocking operation.
func (p *Promise) OnFailure(onFailure func(interface{})) {
	p.failureFn = onFailure

	// TODO: Dirty hack, just for now; if we know we failed, call the function
	if p.failure != nil {
		onFailure(p.failure)
	}
}

// Then applies the given projection function to each value emitted by the source
// Promise and emits a new Promise with the updated value.
// Projection function returns either transformed value or an error if anything
// went wrong. In case of error, failed Promise is returned.
// Then is a non-blocking operation.
func (p *Promise) Then(project func(interface{}) (e interface{}, err error)) *Promise {
	if p.failure != nil {
		return p
	}

	nextResult, err := project(p.result)

	if err != nil {
		return &Promise{
			result:  nil,
			failure: err,
		}
	}

	return &Promise{
		result:  nextResult,
		failure: nil,
	}
}

// Execute creates a new Promise initialized with a value computed from
// the task function passed as an argument. If the initial value
// evaluation went wrong, failed Promise is returned. The Execute is
// a non-blocking function.
func Execute(task func() (r interface{}, err error)) *Promise {
	result, err := task()

	if err != nil {
		return &Promise{
			result:  nil,
			failure: err,
		}
	}

	return &Promise{
		result:  result,
		failure: nil,
	}
}
