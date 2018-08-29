package tecdsa

import (
	"errors"
	"reflect"
	"testing"
)

func TestAddSignerID(t *testing.T) {
	group := signerGroup{}

	if len(group.signerIDs) != 0 {
		t.Fatal("signerIDs is not empty")
	}

	group.AddSignerID("1001")
	if len(group.signerIDs) != 1 || group.signerIDs[0] != "1001" {
		t.Fatalf("signerIDs should contain only one element with value %v, but is %v", "1001", group.signerIDs)
	}
}

func TestRemoveSignerID(t *testing.T) {
	group := signerGroup{
		signerIDs: []string{"1001", "1002", "1003", "1004"},
	}

	if !reflect.DeepEqual(group.signerIDs, []string{"1001", "1002", "1003", "1004"}) {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002 1003 1004]", group.signerIDs)
	}

	// Remove middle item
	group.RemoveSignerID("1003")

	if !reflect.DeepEqual(group.signerIDs, []string{"1001", "1002", "1004"}) {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002 1004]", group.signerIDs)
	}

	// Remove last item
	group.RemoveSignerID("1004")

	if !reflect.DeepEqual(group.signerIDs, []string{"1001", "1002"}) {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002]", group.signerIDs)
	}

	// Remove first item
	group.RemoveSignerID("1001")

	if !reflect.DeepEqual(group.signerIDs, []string{"1002"}) {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1002]", group.signerIDs)
	}
}

func TestIsActiveSigner(t *testing.T) {
	group := signerGroup{
		signerIDs: []string{"1001", "1002", "1003", "1004"},
	}

	if !group.IsActiveSigner("1003") {
		t.Fatal("signer with ID 1003 should be a member of the group")
	}

	if group.IsActiveSigner("1009") {
		t.Fatal("signer with ID 1009 should not be a member of the group")
	}
}

func TestSize(t *testing.T) {
	group := signerGroup{}
	if group.Size() != 0 {
		t.Fatalf("returned group size %v doesn't match expected %v", group.Size(), 0)
	}

	group.signerIDs = []string{"1001", "1002", "1003", "1004"}
	if group.Size() != 4 {
		t.Fatalf("returned group size %v doesn't match expected %v", group.Size(), 4)
	}
}

func TestIsSignerGroupComplete(t *testing.T) {
	signer := signerCore{
		signerGroup: &signerGroup{
			signerIDs: []string{"1001", "1002", "1003", "1004"},
		},
		groupParameters: &PublicParameters{},
	}

	var tests = map[string]struct {
		initialGroupSize int
		expectedResult   bool
		expectedError    error
	}{
		"positive validation": {
			initialGroupSize: 4,
			expectedResult:   true,
			expectedError:    nil,
		},
		"negative validation - group size doesn't match": {
			initialGroupSize: 3,
			expectedResult:   false,
			expectedError:    errors.New("current signers group size 4 doesn't match expected size 3"),
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			signer.signerGroup.GroupSize = test.initialGroupSize

			result, err := signer.signerGroup.IsSignerGroupComplete()

			if result != test.expectedResult {
				t.Fatalf(
					"unexpected result\nexpected: %v\nactual: %v",
					test.expectedResult,
					result,
				)
			}

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"unexpected error\nexpected: %v\nactual: %v",
					test.expectedError,
					err,
				)
			}
		})
	}
}
