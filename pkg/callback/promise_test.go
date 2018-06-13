package promise

import (
	"fmt"
	"testing"
)

func TestPromise2(t *testing.T) {
	done := make(chan bool)

	promise := NewPromise()

	promise.OnSuccess(func(in interface{}) {
		fmt.Printf("Got %v\n", in)
		done <- true
	})

	promise.OnFailure(func(err error) {
		fmt.Printf("Got error %v\n", err)
		done <- true
	})

	err := promise.Fulfill("dupa")
	if err != nil {
		t.Fatal(err)
	}

	<-done
}
