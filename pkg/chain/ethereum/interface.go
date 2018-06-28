package ethereum

type Watcher interface {
	RegisterSuccessCallback(success SuccessFunc) (string, error)
	RegisterFailureCallback(name string, fail func(error) error) error
}

type SuccessFunc interface {
	Type() string
}
