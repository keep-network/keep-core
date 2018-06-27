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
func Ok(tb testing.TB, errWrap ErrWrap) {
	tb.Helper()
	if errWrap.Err != nil {
		tb.Fatalf(redMsg("unexpected error: %v"), errWrap)
	}
}

// NotOk fails when err is nil
// util.NotOk(t, err, "input=%v", test.input)
func NotOk(tb testing.TB, errWrap ErrWrap, errType int, msgFormat string, msgArgs ...interface{}) {
	tb.Helper()

	if errWrap.Err == nil {
		if len(msgArgs) != 0 {
			tb.Fatalf(redMsg("expected error where: "+msgFormat+", got none"), msgArgs...)
		}
		tb.Fatalf(redMsg("expected error, got none"))
	} else if errWrap.ErrNo != errType {
		tb.Fatalf(redMsg(fmt.Sprintf("expected error (%d) got (%d): %v", errType, errWrap.ErrorNumber(), errWrap.Err.Error())))
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
