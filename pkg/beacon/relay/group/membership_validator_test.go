package group

import (
	"math/big"
	"testing"

	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestIsInGroup(t *testing.T) {
	chain := local.Connect(3, 3, big.NewInt(100))
	signing := chain.Signing()

	publicKey1 := generatePublicKey(t)
	publicKey2 := generatePublicKey(t)
	publicKey3 := generatePublicKey(t)

	address1 := signing.PublicKeyToAddress(*publicKey1)
	address2 := signing.PublicKeyToAddress(*publicKey2)

	validator := NewStakersMembershipValidator(
		[]relaychain.StakerAddress{address1, address2, address2},
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
	chain := local.Connect(3, 3, big.NewInt(100))
	signing := chain.Signing()

	publicKey1 := generatePublicKeyBytes(t)
	publicKey2 := generatePublicKeyBytes(t)
	publicKey3 := generatePublicKeyBytes(t)

	address1 := signing.PublicKeyBytesToAddress(publicKey1)
	address2 := signing.PublicKeyBytesToAddress(publicKey2)

	validator := NewStakersMembershipValidator(
		[]relaychain.StakerAddress{address2, address1, address2},
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

func generatePublicKey(t *testing.T) *ecdsa.PublicKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	return &key.PublicKey
}

func generatePublicKeyBytes(t *testing.T) []byte {
	publicKey := generatePublicKey(t)
	return elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
}
