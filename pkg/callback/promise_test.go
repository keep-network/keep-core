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

	err := promise.Fulfill(expectedResult)

	if err != nil {
		t.Fatal(err)
	}

	result := <-done
	if result != expectedResult {
		t.Errorf("Unexpected value passed to callback\nExpected: %v\nActual:%v\n", expectedResult, result)
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
				t.Fatalf("Errors don't match\nExpected: %v\nActual: %v\n", test.expectedError, error)
			}

		})
	}
	<-done
}
