// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package config_test

import (
	"testing"
	"os"
	"github.com/pkg/errors"
	"github.com/keep-network/keep-core/go/libs/config"
)

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected to see a panic")
		}
	}()
	filePanicFcn()
}

func TestPanic2(t *testing.T) {
	assertPanic(t, filePanicFcn)
}

func TestPanic3(t *testing.T) {
	assertPanic(t, zeroPanicFcn)
}


// -------------
//    Helpers
// -------------

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected to see a panic")
		}
	}()
	f()
}

func filePanicFcn() {
	_, err := os.Open("doesnot-exist.txt")
	if err != nil {
		config.HandlePanic(errors.Wrap(err, "unable to read file"))
	}
}

func divFcn(d int) error {
	if d == 0 {
		return errors.New("divide by 0 attempted")
	}
	return nil
}

func zeroPanicFcn() {
	err := divFcn(0)
	if err != nil {
		config.HandlePanic(err)
	}
}
