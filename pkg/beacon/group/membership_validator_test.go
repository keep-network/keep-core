package group

import (
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/operator"
	"math/big"
	"testing"
)

func TestIsInGroup(t *testing.T) {
	localChain := local_v1.Connect(3, 3, big.NewInt(100))
	signing := localChain.Signing()

	publicKey1 := generatePublicKey(t)
	publicKey2 := generatePublicKey(t)
	publicKey3 := generatePublicKey(t)

	address1, err := signing.PublicKeyToAddress(publicKey1)
	if err != nil {
		t.Fatal(err)
	}

	address2, err := signing.PublicKeyToAddress(publicKey2)
	if err != nil {
		t.Fatal(err)
	}

	validator := NewStakersMembershipValidator(
		[]chain.Address{address1, address2, address2},
		signing,
	)

	if !validator.IsInGroup(publicKey1) {
		t.Errorf("staker with public key 1 has been selected")
	}
	if !validator.IsInGroup(publicKey2) {
		t.Errorf("staker with public key 2 has been selected")
	}
	if validator.IsInGroup(publicKey3) {
		t.Errorf("staker with public key 2 has not been selected")
	}
}

func TestIsValidMembership(t *testing.T) {
	localChain := local_v1.Connect(3, 3, big.NewInt(100))
	signing := localChain.Signing()

	publicKey1 := generatePublicKeyBytes(t)
	publicKey2 := generatePublicKeyBytes(t)
	publicKey3 := generatePublicKeyBytes(t)

	address1 := signing.PublicKeyBytesToAddress(publicKey1)
	address2 := signing.PublicKeyBytesToAddress(publicKey2)

	validator := NewStakersMembershipValidator(
		[]chain.Address{address2, address1, address2},
		signing,
	)

	if !validator.IsValidMembership(1, publicKey2) {
		t.Errorf("staker with public key 2 has been selected at index [0]")
	}
	if !validator.IsValidMembership(2, publicKey1) {
		t.Errorf("staker with public key 1 has been selected at index [1]")
	}
	if !validator.IsValidMembership(3, publicKey2) {
		t.Errorf("staker with public key 2 has been selected at index [2]")
	}
	if validator.IsValidMembership(4, publicKey2) {
		t.Errorf("staker with public key 2 has not been selected at index [2]")
	}
	if validator.IsValidMembership(5, publicKey3) {
		t.Errorf("staker with public key 3 has not been selected")
	}
}

func generatePublicKey(t *testing.T) *operator.PublicKey {
	_, operatorPublicKey, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	return operatorPublicKey
}

func generatePublicKeyBytes(t *testing.T) []byte {
	operatorPublicKey := generatePublicKey(t)

	operatorPublicKeyBytes := operator.MarshalUncompressed(operatorPublicKey)

	return operatorPublicKeyBytes
}
