package util

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"testing"
)

// Ok fails when err is not nil.
func Ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Errorf("expected: no errors\nactual: %v", err)
	}
}

// NotOk fails when err is nil or when the error is an runtime error due to file not found.
func NotOk(tb testing.TB, err error, msgFormat string, msgArgs ...interface{}) {
	tb.Helper()

	if err == nil {
		if len(msgArgs) != 0 {
			tb.Errorf("expected an error where: "+msgFormat+", actual: no errors", msgArgs...)
		} else {
			tb.Errorf("expected an error, actual: no errors")
		}
	} else {
		// Report unexpected error if caused by failure to read a file
		if perr, ok := err.(*os.PathError); ok {
			switch perr.Err.(syscall.Errno) {
			case syscall.ENOENT:
				var errMsg string
				if err, ok := err.(*os.PathError); ok {
					errMsg = fmt.Sprintf("file at path (%s) failed to open", err.Path)
				}
				tb.Errorf(errMsg)
			default:
				tb.Errorf(fmt.Sprintf("actual, unexpected error: %v", err))
			}
		}
	}
}

// Equals fails when expected value is not equal to the actual value.
func Equals(tb testing.TB, expected, actual interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(expected, actual) {
		tb.Errorf("expected: %#v\nactual: %#v", expected, actual)
	}
}
