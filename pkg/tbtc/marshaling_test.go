package tbtc

import (
	"crypto/ecdsa"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"math/big"
	"reflect"
	"testing"
)

func TestWalletMarshalling(t *testing.T) {
	marshaled := sampleWallet()

	unmarshaled := &wallet{}

	if err := pbutils.RoundTrip(marshaled, unmarshaled); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(marshaled, unmarshaled) {
		t.Fatal("unexpected content of unmarshaled wallet")
	}
}

func TestSignerMarshalling(t *testing.T) {
	marshaled := sampleSigner(t)

	unmarshaled := &signer{}

	if err := pbutils.RoundTrip(marshaled, unmarshaled); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(marshaled, unmarshaled) {
		t.Fatal("unexpected content of unmarshaled signer")
	}
}

func sampleWallet() *wallet {
	publicKeyX, _ := new(big.Int).SetString(
		"92492306118178589821640584737240636977398594678247616965910942704932180187323",
		10,
	)
	publicKeyY, _ := new(big.Int).SetString(
		"27954057508764275913470910100133573369328128015811591924683199269013496685879",
		10,
	)
	publicKey := &ecdsa.PublicKey{
		Curve: tecdsa.Curve,
		X:     publicKeyX,
		Y:     publicKeyY,
	}

	wallet := &wallet{
		publicKey: publicKey,
		signingGroupOperators: []chain.Address{
			"address-1",
			"address-2",
			"address-3",
			"address-3",
			"address-5",
		},
	}
	return wallet
}

func sampleSigner(t *testing.T) *signer {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	privateKeyShare := tecdsa.NewPrivateKeyShare(testData[0])

	wallet := sampleWallet()

	marshaled := &signer{
		wallet:                  *wallet,
		signingGroupMemberIndex: group.MemberIndex(1),
		privateKeyShare:         privateKeyShare,
	}
	return marshaled
}
