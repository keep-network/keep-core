package util

import (
	"errors"
)

// AppendErrMsgs appends the ValidationError return to the error messages list.
func AppendErrMsgs(errMsgs []string, err error) []string {
	if err != nil {
		errMsgs = append(errMsgs, err.Error())
	}
	return errMsgs
}

// Err is typically passed to the return of a ValidationError function.
func Err(errMsgs []string) error {
	var err error
	if len(errMsgs) > 0 {
		err = errors.New(Join(errMsgs, "\n"))
	}
	return err
}
