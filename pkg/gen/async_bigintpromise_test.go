package main

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/gen/async"
)

func TestBigIntPromiseOnSuccessFulfill(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan *big.Int)

	expectedResult := big.NewInt(8)

	promise := &async.BigIntPromise{}

	promise.OnSuccess(func(in *big.Int) {
		done <- in
	})

	promise.OnFailure(func(err error) {
		t.Fatal("`OnFailure` was called for `Fulfill`")
	})

	err := promise.Fulfill(expectedResult)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case result := <-done:
		if result != expectedResult {
			t.Errorf(
				"Unexpected value passed to callback\nExpected: %v\nActual:%v\n",
				expectedResult,
				result,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestPromiseOnSuccessAlreadyFulfilled(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan interface{})

	expectedResult := big.NewInt(18)

	promise := &async.BigIntPromise{}

	// first fulfill, then install callback
	err := promise.Fulfill(expectedResult)
	if err != nil {
		t.Fatal(err)
	}

	promise.OnSuccess(func(in *big.Int) {
		done <- in
	})

	select {
	case result := <-done:
		if result != expectedResult {
			t.Errorf(
				"Unexpected value passed to callback\nExpected: %v\nActual:%v\n",
				expectedResult,
				result,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestPromiseOnCompleteFulfill(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan interface{})

	expectedResult := big.NewInt(128)

	promise := &async.BigIntPromise{}

	promise.OnComplete(func(in *big.Int, err error) {
		if err != nil {
			t.Fatal("Error should be nil")
		}

		done <- in
	})

	err := promise.Fulfill(expectedResult)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case result := <-done:
		if result != expectedResult {
			t.Errorf(
				"Unexpected value passed to callback\nExpected: %v\nActual:%v\n",
				expectedResult,
				result,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestPromiseOnCompleteAlreadyFulfilled(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan interface{})

	expectedResult := big.NewInt(1238)

	promise := &async.BigIntPromise{}

	// first fulfill, then install callback
	err := promise.Fulfill(expectedResult)
	if err != nil {
		t.Fatal(err)
	}

	promise.OnComplete(func(in *big.Int, err error) {
		if err != nil {
			t.Fatal("Error should be nil")
		}

		done <- in
	})

	select {
	case actualResult := <-done:
		if expectedResult != actualResult {
			t.Errorf(
				"Unexpected value passed to callback\nExpected: %v\nActual %v\n",
				expectedResult,
				actualResult,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestPromiseOnFailureFail(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan interface{})

	expectedResult := fmt.Errorf("it's not working")

	promise := &async.BigIntPromise{}

	promise.OnFailure(func(err error) {
		done <- err
	})

	promise.OnSuccess(func(in *big.Int) {
		t.Fatal("`OnSuccess` was called for `Fail`")
	})

	err := promise.Fail(expectedResult)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case result := <-done:
		if result != expectedResult {
			t.Errorf(
				"Unexpected value passed to callback\nExpected: %v\nActual:%v\n",
				expectedResult,
				result,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestPromiseOnFailureAlreadyFailed(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan interface{})

	expectedError := fmt.Errorf("i just can't")

	promise := &async.BigIntPromise{}

	// first fail, then install callback
	err := promise.Fail(expectedError)
	if err != nil {
		t.Fatal(err)
	}

	promise.OnFailure(func(err error) {
		done <- err
	})

	select {
	case actualError := <-done:
		if !reflect.DeepEqual(expectedError, actualError) {
			t.Errorf(
				"Unexpected error passed to callback\nExpected: %v\nActual %v\n",
				expectedError,
				actualError,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestPromiseOnCompleteFail(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan interface{})

	expectedFailure := fmt.Errorf("catwoman")

	promise := &async.BigIntPromise{}

	promise.OnComplete(func(in *big.Int, err error) {
		if in != nil {
			t.Fatal("Evaluated value should be nil")
		}

		done <- err
	})

	err := promise.Fail(expectedFailure)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case result := <-done:
		if result != expectedFailure {
			t.Errorf(
				"Unexpected failure passed to callback\nExpected: %v\nActual:%v\n",
				expectedFailure,
				result,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestPromiseOnCompleteAlreadyFailed(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan interface{})

	expectedError := fmt.Errorf("nope nope nope")

	promise := &async.BigIntPromise{}

	// first fail, then install callback
	err := promise.Fail(expectedError)
	if err != nil {
		t.Fatal(err)
	}

	promise.OnComplete(func(in *big.Int, err error) {
		if in != nil {
			t.Fatal("Promise's value should be nil")
		}

		done <- err
	})

	select {
	case actualError := <-done:
		if !reflect.DeepEqual(expectedError, actualError) {
			t.Errorf(
				"Unexpected error passed to callback\nExpected: %v\nActual %v\n",
				expectedError,
				actualError,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestPromiseAlreadyCompleted(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan bool)

	var tests = map[string]struct {
		function      func() error
		expectedError error
	}{
		"Fulfill with result `promise already completed`": {
			function: func() error {
				promise := &async.BigIntPromise{}
				promise.OnSuccess(func(in *big.Int) { done <- true })
				promise.Fulfill(nil)
				return promise.Fulfill(nil)
			},
			expectedError: fmt.Errorf("promise already completed"),
		},
		"Fail with result `promise already completed`": {
			function: func() error {
				promise := &async.BigIntPromise{}
				promise.OnFailure(func(error) { done <- true })
				promise.Fail(nil)
				return promise.Fail(nil)
			},
			expectedError: fmt.Errorf("promise already completed"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			error := test.function()

			select {
			case <-done:
				if !reflect.DeepEqual(test.expectedError, error) {
					t.Errorf(
						"Errors don't match\nExpected: %v\nActual: %v\n",
						test.expectedError,
						error,
					)
				}
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
		})
	}
}

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}
