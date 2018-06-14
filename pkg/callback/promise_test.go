package promise

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPromiseOnSuccessFulfill(t *testing.T) {
	done := make(chan interface{})

	expectedResult := "batman"

	promise := NewPromise()

	promise.OnSuccess(func(in interface{}) {
		done <- in
	})

	promise.OnFailure(func(err error) {
		t.Fatal("`OnFailure` was called for `Fulfill`")
	})

	err := promise.Fulfill(expectedResult)

	if err != nil {
		t.Fatal(err)
	}

	result := <-done
	if result != expectedResult {
		t.Errorf(
			"Unexpected value passed to callback\nExpected: %v\nActual:%v\n",
			expectedResult,
			result)
	}
}

func TestPromiseOnCompleteFulfill(t *testing.T) {
	done := make(chan interface{})

	expectedResult := "robin"

	promise := NewPromise()

	promise.OnComplete(func(in interface{}, err error) {
		if err != nil {
			t.Fatal("Error should be nil")
		}

		done <- in
	})

	err := promise.Fulfill(expectedResult)

	if err != nil {
		t.Fatal(err)
	}

	result := <-done
	if result != expectedResult {
		t.Errorf(
			"Unexpected value passed to callback\nExpected: %v\nActual:%v\n",
			expectedResult,
			result,
		)
	}
}

func TestPromiseOnFailureFail(t *testing.T) {
	done := make(chan interface{})

	expectedResult := fmt.Errorf("it's not working")

	promise := NewPromise()

	promise.OnFailure(func(err error) {
		done <- err
	})

	promise.OnSuccess(func(in interface{}) {
		t.Fatal("`OnSuccess` was called for `Fail`")
	})

	err := promise.Fail(expectedResult)

	if err != nil {
		t.Fatal(err)
	}

	result := <-done
	if result != expectedResult {
		t.Errorf(
			"Unexpected value passed to callback\nExpected: %v\nActual:%v\n",
			expectedResult,
			result)
	}
}

func TestPromiseOnCompleteFail(t *testing.T) {
	done := make(chan interface{})

	expectedFailure := fmt.Errorf("catwoman")

	promise := NewPromise()

	promise.OnComplete(func(in interface{}, err error) {
		if in != nil {
			t.Fatal("Evaluated value should be nil")
		}

		done <- err
	})

	err := promise.Fail(expectedFailure)

	if err != nil {
		t.Fatal(err)
	}

	result := <-done
	if result != expectedFailure {
		t.Errorf(
			"Unexpected failure passed to callback\nExpected: %v\nActual:%v\n",
			expectedFailure,
			result,
		)
	}
}

func TestPromiseFulfill(t *testing.T) {
	promise := NewPromise()

	if promise.isComplete {
		t.Error("Promise is completed")
	}

	err := promise.Fulfill(nil)

	if err != nil {
		t.Errorf("Fulfill returned an error: %v", err)
	}

	if !promise.isComplete {
		t.Error("Promise is not completed")
	}
}

func TestPromiseFail(t *testing.T) {
	promise := NewPromise()

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
	done := make(chan bool)

	var tests = map[string]struct {
		function      func() error
		expectedError error
	}{
		"Fulfill with result `promise already completed`": {
			function: func() error {
				promise := NewPromise().OnSuccess(func(in interface{}) { done <- true })
				promise.Fulfill(nil)
				return promise.Fulfill(nil)
			},
			expectedError: fmt.Errorf("promise already completed"),
		},
		"Fail with result `promise already completed`": {
			function: func() error {
				promise := NewPromise().OnFailure(func(error) { done <- true })
				promise.Fail(nil)
				return promise.Fail(nil)
			},
			expectedError: fmt.Errorf("promise already completed"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			error := test.function()
			if !reflect.DeepEqual(test.expectedError, error) {
				t.Errorf(
					"Errors don't match\nExpected: %v\nActual: %v\n",
					test.expectedError,
					error)
			}

		})
	}
	<-done
}
