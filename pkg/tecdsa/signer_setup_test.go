package tecdsa

import (
	"fmt"
	"testing"
)

func setupGroup(group []*LocalSigner) error {
	var err error

	// Initialize master public key for multi-trapdoor commitment scheme.
	// Each signer generates a master public key share which is a point in
	// G2 abstract cyclic group of bn256 curve. The share is broadcasted in
	// MasterPublicKeyShareMessage.
	// The shares are combined by adding the points which results in a point
	// which is a master public key.
	masterPublicKeyShareMessages := make([]*MasterPublicKeyShareMessage, len(group))
	for i, signer := range group {
		masterPublicKeyShareMessages[i], err = signer.GenerateMasterPublicKeyShare()
		if err != nil {
			return err
		}
	}

	masterPublicKey, err := group[0].CombineMasterPublicKeyShares(masterPublicKeyShareMessages)
	if err != nil {
		return err
	}

	for _, signer := range group {
		signer.commitmentMasterPublicKey = masterPublicKey
	}

	return nil
}

func TestRegisterSignerID(t *testing.T) {
	signer := NewLocalSigner(nil, nil, nil)

	if len(signer.signerIDs) != 0 {
		t.Fatal("signerIDs is not empty")
	}

	signer.RegisterSignerID("1001")
	if len(signer.signerIDs) != 1 || signer.signerIDs[0] != "1001" {
		t.Fatalf("signerIDs should contain only one element with value %v, but is %v", "1001", signer.signerIDs)
	}
}

func TestRemoveSignerID(t *testing.T) {
	signer := NewLocalSigner(nil, nil, nil)

	signer.signerIDs = []string{"1001", "1002", "1003", "1004"}

	if len(signer.signerIDs) != 4 || fmt.Sprint(signer.signerIDs) != "[1001 1002 1003 1004]" {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002 1003 1004]", signer.signerIDs)
	}

	// Remove middle item
	signer.RemoveSignerID("1003")

	if len(signer.signerIDs) != 3 || fmt.Sprint(signer.signerIDs) != "[1001 1002 1004]" {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002 1004]", signer.signerIDs)
	}

	// Remove last item
	signer.RemoveSignerID("1004")

	if len(signer.signerIDs) != 2 || fmt.Sprint(signer.signerIDs) != "[1001 1002]" {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1001 1002]", signer.signerIDs)
	}

	// Remove first item
	signer.RemoveSignerID("1001")

	if len(signer.signerIDs) != 1 || fmt.Sprint(signer.signerIDs) != "[1002]" {
		t.Fatalf("signer IDs list doesn't match expected\nExpected: %v\nActual: %v", "[1002]", signer.signerIDs)
	}
}
