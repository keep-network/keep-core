package tbtc

import (
	"crypto/ecdsa"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"reflect"
	"testing"
)

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

func sampleSigner(t *testing.T) *signer {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	share := testData[0]

	walletPublicKey := &ecdsa.PublicKey{
		Curve: share.ECDSAPub.Curve(),
		X:     share.ECDSAPub.X(),
		Y:     share.ECDSAPub.Y(),
	}

	walletSigningGroupOperators := []chain.Address{
		"address-1",
		"address-2",
		"address-3",
		"address-3",
		"address-5",
	}

	privateKeyShare := tecdsa.NewPrivateKeyShare(share)

	marshaled := &signer{
		wallet:                  wallet{
			publicKey:             walletPublicKey,
			signingGroupOperators: walletSigningGroupOperators,
		},
		signingGroupMemberIndex: group.MemberIndex(1),
		privateKeyShare:         privateKeyShare,
	}
	return marshaled
}
