package promise

import (
	"fmt"
	"testing"
)

func TestPromise2(t *testing.T) {
	done := make(chan bool)

	promise := newPromise()

	promise.onSuccess(func(in interface{}) {
		fmt.Printf("Got %v\n", in)
		done <- true
	})

	promise.onFailure(func(err error) {
		fmt.Printf("Got error %v\n", err)
		done <- true
	})

	err := promise.fulfill("dupa")
	if err != nil {
		t.Fatal(err)
	}

	<-done
}
