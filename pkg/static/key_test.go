package static

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pborman/uuid"
)

// `geth` uses `go-ethereum` library to generate key with secp256k1 curve.

// As long as all curve parameters are the same, this operation is valid.
// This test ensures that secp256k1 from `go-ethereum` and secp256k1 from
// `btcsuite` are the same. If this test starts to fails, we'll need to revisit
// how the key is ported from one instance to another in `toLibp2pKey` function.
func TestSameCurveAsEthereum(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	staticPrivateKey, _ := EthereumKeyToStaticKey(ethereumKey)

	ethereumCurve := ethereumKey.PrivateKey.Curve.Params()
	staticKeyCurve := staticPrivateKey.Curve.Params()

	if ethereumCurve.P.Cmp(staticKeyCurve.P) != 0 {
		t.Errorf(
			"unexpected P\nexpected: %v\nactual: %v",
			ethereumCurve.P,
			staticKeyCurve.P,
		)
	}

	if ethereumCurve.N.Cmp(staticKeyCurve.N) != 0 {
		t.Errorf(
			"unexpected N\nexpected: %v\nactual: %v",
			ethereumCurve.N,
			staticKeyCurve.N,
		)
	}

	if ethereumCurve.B.Cmp(staticKeyCurve.B) != 0 {
		t.Errorf(
			"unexpected B\nexpected: %v\nactual: %v",
			ethereumCurve.B,
			staticKeyCurve.B,
		)
	}

	if ethereumCurve.Gx.Cmp(staticKeyCurve.Gx) != 0 {
		t.Errorf(
			"unexpected Gx\nexpected: %v\nactual: %v",
			ethereumCurve.Gx,
			staticKeyCurve.Gx,
		)
	}

	if ethereumCurve.Gy.Cmp(staticKeyCurve.Gy) != 0 {
		t.Errorf(
			"unexpected Gy\nexpected: %v\nactual: %v",
			ethereumCurve.Gy,
			staticKeyCurve.Gy,
		)
	}

	if ethereumCurve.BitSize != staticKeyCurve.BitSize {
		t.Errorf(
			"unexpected BitSize\nexpected: %v\nactual: %v",
			ethereumCurve.BitSize,
			staticKeyCurve.BitSize,
		)
	}
}

func TestSameKeyAsEthereum(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	staticPrivKey, staticPubKey := EthereumKeyToStaticKey(
		ethereumKey,
	)

	if ethereumKey.PrivateKey.D.Cmp(staticPrivKey.D) != 0 {
		t.Errorf(
			"unexpected D\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.D,
			staticPrivKey.D,
		)
	}

	if ethereumKey.PrivateKey.PublicKey.X.Cmp(staticPubKey.X) != 0 {
		t.Errorf(
			"unexpected X\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.PublicKey.X,
			staticPubKey.X,
		)
	}

	if ethereumKey.PrivateKey.PublicKey.Y.Cmp(staticPubKey.Y) != 0 {
		t.Errorf(
			"unexpected Y\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.PublicKey.Y,
			staticPubKey.Y,
		)
	}
}

func TestNetworkPubKeyToAddress(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	_, staticPublicKey := EthereumKeyToStaticKey(ethereumKey)

	ethAddress := crypto.PubkeyToAddress(ethereumKey.PrivateKey.PublicKey).String()

	staticKeyAddress := PubKeyToEthAddress(staticPublicKey)

	if ethAddress != staticKeyAddress {
		t.Errorf(
			"unexpected address\nexpected: %v\nactual: %v",
			ethAddress,
			staticKeyAddress,
		)
	}
}

func generateEthereumKey() (*keystore.Key, error) {
	ethCurve := secp256k1.S256()

	ethereumKey, err := ecdsa.GenerateKey(ethCurve, rand.Reader)
	if err != nil {
		return nil, err
	}

	id := uuid.NewRandom()

	return &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(ethereumKey.PublicKey),
		PrivateKey: ethereumKey,
	}, nil
}
