package util

import (
	"fmt"
	"reflect"
	"testing"
)

// Assert fails when condition is false
func Assert(tb testing.TB, condition bool, message string, got ...interface{}) {
	tb.Helper()
	if !condition {
		tb.Fatalf(redMsg(message), got...)
	}
}

// Ok fails when err is not nil
func Ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf(redMsg("unexpected error: %v"), err)
	}
}

// NotOk fails when err is nil
func NotOk(tb testing.TB, err error, msgFormat string, msgArgs ...interface{}) {
	tb.Helper()

	if err == nil {
		if len(msgArgs) != 0 {
			tb.Fatalf(redMsg("expected error where: "+msgFormat+", got none"), msgArgs...)
		} else {
			tb.Fatalf(redMsg("expected error, got none"))
		}
	}
}

// Equals fails when expected is not equal to got
func Equals(tb testing.TB, expected, got interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(expected, got) {
		tb.Fatalf(redMsg("expected: %#v\n\ngot: %#v"), expected, got)
	}
}

func redMsg(msg string) string {
	return fmt.Sprintf("\033[31m%s\033[39m\n", msg)
}
