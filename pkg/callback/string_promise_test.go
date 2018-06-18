package callback

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestStringPromiseOnSuccessFulfill(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	done := make(chan interface{})

	expectedResult := "batman"

	promise := StringPromise{}

	promise.OnSuccess(func(in string) {
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

	expectedResult := "conan the barbarian"

	promise := StringPromise{}

	// first fulfill, then install callback
	err := promise.Fulfill(expectedResult)
	if err != nil {
		t.Fatal(err)
	}

	promise.OnSuccess(func(in string) {
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

	expectedResult := "robin"

	promise := StringPromise{}

	promise.OnComplete(func(in string, err error) {
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

	expectedResult := "conan the conqueror"

	promise := StringPromise{}

	// first fulfill, then install callback
	err := promise.Fulfill(expectedResult)
	if err != nil {
		t.Fatal(err)
	}

	promise.OnComplete(func(in string, err error) {
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

	promise := StringPromise{}

	promise.OnFailure(func(err error) {
		done <- err
	})

	promise.OnSuccess(func(in string) {
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

	promise := StringPromise{}

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

	promise := StringPromise{}

	promise.OnComplete(func(in string, err error) {
		if in != "" {
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

	promise := StringPromise{}

	// first fail, then install callback
	err := promise.Fail(expectedError)
	if err != nil {
		t.Fatal(err)
	}

	promise.OnComplete(func(in string, err error) {
		if in != "" {
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

func TestPromiseFulfill(t *testing.T) {
	promise := StringPromise{}

	if promise.isComplete {
		t.Error("Promise is completed")
	}

	err := promise.Fulfill("")

	if err != nil {
		t.Errorf("Fulfill returned an error: %v", err)
	}

	if !promise.isComplete {
		t.Error("Promise is not completed")
	}
}

func TestPromiseFail(t *testing.T) {
	promise := StringPromise{}

	if promise.isComplete {
		t.Error("Promise is completed")
	}

	err := promise.Fail(nil)
	if err != nil {
		t.Errorf("Fail returned an error: %v", err)
	}

	if !promise.isComplete {
		t.Error("Promise is not completed")
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
				promise := StringPromise{}
				promise.OnSuccess(func(in string) { done <- true })
				promise.Fulfill("")
				return promise.Fulfill("")
			},
			expectedError: fmt.Errorf("promise already completed"),
		},
		"Fail with result `promise already completed`": {
			function: func() error {
				promise := StringPromise{}
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
