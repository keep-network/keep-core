package tbtc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
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

func TestSignerMarshalling_NonTECDSAKey(t *testing.T) {
	signer := sampleSigner(t)

	p256 := elliptic.P256()

	// Use a non-secp256k1 based key to cause the expected failure.
	signer.wallet.publicKey = &ecdsa.PublicKey{
		Curve: p256,
		X:     p256.Params().Gx,
		Y:     p256.Params().Gy,
	}

	_, err := signer.Marshal()

	testutils.AssertErrorsSame(t, errIncompatiblePublicKey, err)
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
		wallet: wallet{
			publicKey:             walletPublicKey,
			signingGroupOperators: walletSigningGroupOperators,
		},
		signingGroupMemberIndex: group.MemberIndex(1),
		privateKeyShare:         privateKeyShare,
	}
	return marshaled
}
