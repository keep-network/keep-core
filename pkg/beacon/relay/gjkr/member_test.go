package gjkr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

func TestNewMemberWithInvalidID(t *testing.T) {
	expectedError := fmt.Errorf("could not create a new member [member index must be >= 1]")

	_, err := NewMember(group.MemberIndex(0), 13, 13, nil)

	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedError, err)
	}
}

func TestMemberIDValidate(t *testing.T) {
	var tests = map[string]struct {
		id            group.MemberIndex
		expectedError error
	}{
		"id = 0": {
			id:            group.MemberIndex(0),
			expectedError: fmt.Errorf("member index must be >= 1"),
		},
		"id = 1": {
			id:            1,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		err := test.id.Validate()

		if !reflect.DeepEqual(err, test.expectedError) {
			t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedError, err)
		}
	}
}
