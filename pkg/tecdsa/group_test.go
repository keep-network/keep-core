package tecdsa

import (
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
	if group.SignerCount() != 0 {
		t.Fatalf("returned group size %v doesn't match expected %v", group.SignerCount(), 0)
	}

	group.signerIDs = []string{"1001", "1002", "1003", "1004"}
	if group.SignerCount() != 4 {
		t.Fatalf("returned group size %v doesn't match expected %v", group.SignerCount(), 4)
	}
}

func TestIsSignerGroupComplete(t *testing.T) {
	signer := signerCore{
		signerGroup: &signerGroup{
			signerIDs: []string{"1001", "1002", "1003", "1004"},
		},
		signatureParameters: &PublicSignatureParameters{},
	}

	var tests = map[string]struct {
		initialGroupSize int
		expectedResult   bool
	}{
		"positive validation": {
			initialGroupSize: 4,
			expectedResult:   true,
		},
		"negative validation - group size doesn't match": {
			initialGroupSize: 3,
			expectedResult:   false,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			signer.signerGroup.InitialGroupSize = test.initialGroupSize

			if signer.signerGroup.IsSignerGroupComplete() != test.expectedResult {
				t.Fatalf(
					"unexpected result\nexpected: %v\nactual: %v",
					test.expectedResult,
					signer.signerGroup.IsSignerGroupComplete(),
				)
			}
		})
	}
}
