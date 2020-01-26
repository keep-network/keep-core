package groupselection

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

	address1 := signing.PublicKeyToAddress(publicKey1)
	address2 := signing.PublicKeyToAddress(publicKey2)

	validator := NewMembershipValidator(
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

func TestIsSelectedAtIndex(t *testing.T) {
	chain := local.Connect(3, 3, big.NewInt(100))
	signing := chain.Signing()

	publicKey1 := generatePublicKey(t)
	publicKey2 := generatePublicKey(t)
	publicKey3 := generatePublicKey(t)

	address1 := signing.PublicKeyToAddress(publicKey1)
	address2 := signing.PublicKeyToAddress(publicKey2)

	validator := NewMembershipValidator(
		[]relaychain.StakerAddress{address2, address1, address2},
		signing,
	)

	if !validator.IsSelectedAtIndex(0, publicKey2) {
		t.Errorf("staker with public key 2 has been selected at index [0]")
	}
	if !validator.IsSelectedAtIndex(1, publicKey1) {
		t.Errorf("staker with public key 1 has been selected at index [1]")
	}
	if !validator.IsSelectedAtIndex(2, publicKey2) {
		t.Errorf("staker with public key 2 has been selected at index [2]")
	}
	if validator.IsSelectedAtIndex(1, publicKey2) {
		t.Errorf("staker with public key 2 has not been selected at index [2]")
	}
	if validator.IsSelectedAtIndex(1, publicKey3) {
		t.Errorf("staker with public key 3 has not been selected")
	}
}

func generatePublicKey(t *testing.T) ecdsa.PublicKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	return key.PublicKey
}
