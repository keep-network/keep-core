package tecdsa

import (
	"crypto/elliptic"
	"reflect"
	"testing"

	"github.com/bnb-chain/tss-lib/crypto"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestPrivateKeyShareMarshalling(t *testing.T) {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	privateKeyShare := NewPrivateKeyShare(testData[0])

	unmarshaled := &PrivateKeyShare{}

	if err := pbutils.RoundTrip(privateKeyShare, unmarshaled); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(privateKeyShare, unmarshaled) {
		t.Fatal("unexpected content of unmarshaled private key share")
	}
}

func TestPrivateKeyShareMarshalling_NonTECDSAKey(t *testing.T) {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	privateKeyShare := NewPrivateKeyShare(testData[0])

	p256 := elliptic.P256()

	ecPoint, err := crypto.NewECPoint(p256, p256.Params().Gx, p256.Params().Gy)
	if err != nil {
		t.Fatal(err)
	}

	// Use a non-secp256k1 based key to cause the expected failure.
	privateKeyShare.data.ECDSAPub = ecPoint

	_, err = privateKeyShare.Marshal()

	testutils.AssertErrorsSame(t, ErrIncompatiblePublicKey, err)
}
