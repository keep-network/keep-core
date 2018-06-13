package promise

import (
	"fmt"
	"testing"
)

func TestPromise(t *testing.T) {
	Execute(func() (interface{}, error) {
		return "a", nil
	}).Then(func(in interface{}) (interface{}, error) {
		return in.(string) + "b", nil
	}).OnSuccess(func(in interface{}) {
		fmt.Printf("Got %v\n", in)
	})
}
