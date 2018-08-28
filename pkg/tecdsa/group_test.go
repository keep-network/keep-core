package tecdsa

import (
	"fmt"
	"testing"
)

func TestRegisterSignerID(t *testing.T) {
	group := signerGroup{}

	if len(group.signerIDs) != 0 {
		t.Fatal("signerIDs is not empty")
	}

	group.RegisterSignerID("1001")
	if len(group.signerIDs) != 1 || group.signerIDs[0] != "1001" {
		t.Fatalf("signerIDs should contain only one element with value %v, but is %v", "1001", group.signerIDs)
	}
}

func TestRemoveSignerID(t *testing.T) {
	group := signerGroup{
		signerIDs: []string{"1001", "1002", "1003", "1004"},
	}

	if len(group.signerIDs) != 4 || fmt.Sprint(group.signerIDs) != "[1001 1002 1003 1004]" {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002 1003 1004]", group.signerIDs)
	}

	// Remove middle item
	group.RemoveSignerID("1003")

	if len(group.signerIDs) != 3 || fmt.Sprint(group.signerIDs) != "[1001 1002 1004]" {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002 1004]", group.signerIDs)
	}

	// Remove last item
	group.RemoveSignerID("1004")

	if len(group.signerIDs) != 2 || fmt.Sprint(group.signerIDs) != "[1001 1002]" {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002]", group.signerIDs)
	}

	// Remove first item
	group.RemoveSignerID("1001")

	if len(group.signerIDs) != 1 || fmt.Sprint(group.signerIDs) != "[1002]" {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1002]", group.signerIDs)
	}
}

func TestIsActiveSigner(t *testing.T) {
	group := signerGroup{
		signerIDs: []string{"1001", "1002", "1003", "1004"},
	}

	if !group.IsActiveSigner("1003") {
		t.Fatal("signer with ID 1003 should be active")
	}

	if group.IsActiveSigner("1009") {
		t.Fatal("signer with ID 1009 should not be active")
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
