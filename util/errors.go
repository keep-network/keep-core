package util

import "fmt"

// ErrWrap wraps an error and its error number
type ErrWrap struct {
	ErrNo int
	Err   error
}

// Error is a behavior only available for the ErrWrap type
func (ew *ErrWrap) Error() string {
	return fmt.Sprintf("[%d] %s", ew.ErrNo, ew.Err.Error())
}

// ErrorNumber returns the error number associated with the wrapped error
func (ew ErrWrap) ErrorNumber() int {
	return ew.ErrNo
}
