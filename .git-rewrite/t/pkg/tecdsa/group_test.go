package tecdsa

import (
	"reflect"
	"testing"
)

func TestAddSignerID(t *testing.T) {
	group := signerGroup{}

	// Check empty signers group
	if len(group.signerIDs) != 0 {
		t.Fatal("signerIDs is not empty")
	}

	// Add Signer to the group
	expectedSigners := []string{"1001"}
	group.AddSignerID("1001")

	if !reflect.DeepEqual(expectedSigners, group.signerIDs) {
		t.Fatalf(
			"signer IDs list doesn't match expected\nExpected: %v\nActual: %v",
			expectedSigners,
			group.signerIDs,
		)
	}
}

func TestRemoveSignerID(t *testing.T) {
	expectedSigners := []string{"1001", "1002", "1003", "1004"}
	group := signerGroup{
		signerIDs: expectedSigners,
	}

	if !reflect.DeepEqual(group.signerIDs, expectedSigners) {
		t.Fatalf(
			"signer IDs list doesn't match expected\nExpected: %v\nActual: %v",
			expectedSigners,
			group.signerIDs,
		)
	}

	// Remove middle item
	expectedSigners = []string{"1001", "1002", "1004"}
	group.RemoveSignerID("1003")

	if !reflect.DeepEqual(group.signerIDs, expectedSigners) {
		t.Fatalf(
			"signer IDs list doesn't match expected\nExpected: %v\nActual: %v",
			expectedSigners,
			group.signerIDs,
		)
	}

	// Remove last item
	expectedSigners = []string{"1001", "1002"}
	group.RemoveSignerID("1004")

	if !reflect.DeepEqual(group.signerIDs, expectedSigners) {
		t.Fatalf(
			"signer IDs list doesn't match expected\nExpected: %v\nActual: %v",
			expectedSigners,
			group.signerIDs,
		)
	}

	// Remove first item
	expectedSigners = []string{"1002"}
	group.RemoveSignerID("1001")

	if !reflect.DeepEqual(group.signerIDs, expectedSigners) {
		t.Fatalf(
			"signer IDs list doesn't match expected\nExpected: %v\nActual: %v",
			expectedSigners,
			group.signerIDs,
		)
	}
}

func TestContains(t *testing.T) {
	expectedSigners := []string{"1001", "1002", "1003", "1004"}
	group := signerGroup{
		signerIDs: expectedSigners,
	}

	if !group.Contains("1003") {
		t.Fatal("signer with ID 1003 should be a member of the group")
	}

	if group.Contains("1009") {
		t.Fatal("signer with ID 1009 should not be a member of the group")
	}
}

func TestSignerCount(t *testing.T) {
	group := signerGroup{}
	if group.SignerCount() != 0 {
		t.Fatalf(
			"unexpected signer count\nExpected: %v\nActual: %v",
			0,
			group.SignerCount(),
		)
	}
	if group.PeerSignerCount() != 0 {
		t.Fatalf(
			"unexpected peer signer count\nExpected: %v\nActual:%v",
			0,
			group.PeerSignerCount(),
		)
	}

	group.signerIDs = []string{"1001", "1002", "1003", "1004"}
	if group.SignerCount() != 4 {
		t.Fatalf(
			"unexpected signer count\nExpected: %v\nActual: %v",
			4,
			group.SignerCount(),
		)
	}
	if group.PeerSignerCount() != 3 {
		t.Fatalf(
			"unexpected peer signer count\nExpected: %v\nActual:%v",
			3,
			group.PeerSignerCount(),
		)
	}
}

func TestIsSignerGroupComplete(t *testing.T) {
	signer := signerCore{
		signerGroup: &signerGroup{
			signerIDs: []string{"1001", "1002", "1003", "1004"},
		},
		publicParameters: &PublicParameters{},
	}

	var tests = map[string]struct {
		initialGroupSize int
		expectedResult   bool
	}{
		"group is complete": {
			initialGroupSize: 4,
			expectedResult:   true,
		},
		"group is not complete": {
			initialGroupSize: 5,
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
