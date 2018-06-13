package promise

import (
	"fmt"
	"testing"
)

func TestPromise(t *testing.T) {
	Execute(func() (interface{}, error) {
		return "a", nil
	}).Then(func(in interface{}) (interface{}, error) {
		if s, ok := in.(string); !ok {
			return s + "b", nil
		}
		return nil, fmt.Errorf("Unexpected type %v", in)
	}).OnComplete(func(in interface{}) {
		fmt.Printf("Got %v\n", in)
	})
}
