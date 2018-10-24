package libp2p

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pborman/uuid"
)

func TestSameCurveAsEthereum(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	libp2pKey := toLibp2pKey(ethereumKey)

	ethereumCurve := ethereumKey.PrivateKey.Curve.Params()
	libp2pCurve := libp2pKey.Curve.Params()

	if ethereumCurve.P.Cmp(libp2pCurve.P) != 0 {
		t.Errorf(
			"unexpected P\nexpected: %v\nactual: %v",
			ethereumCurve.P,
			libp2pCurve.P,
		)
	}

	if ethereumCurve.N.Cmp(libp2pCurve.N) != 0 {
		t.Errorf(
			"unexpected N\nexpected: %v\nactual: %v",
			ethereumCurve.N,
			libp2pCurve.N,
		)
	}

	if ethereumCurve.B.Cmp(libp2pCurve.B) != 0 {
		t.Errorf(
			"unexpected B\nexpected: %v\nactual: %v",
			ethereumCurve.B,
			libp2pCurve.B,
		)
	}

	if ethereumCurve.Gx.Cmp(libp2pCurve.Gx) != 0 {
		t.Errorf(
			"unexpected Gx\nexpected: %v\nactual: %v",
			ethereumCurve.Gx,
			libp2pCurve.Gx,
		)
	}

	if ethereumCurve.Gy.Cmp(libp2pCurve.Gy) != 0 {
		t.Errorf(
			"unexpected Gy\nexpected: %v\nactual: %v",
			ethereumCurve.Gy,
			libp2pCurve.Gy,
		)
	}

	if ethereumCurve.BitSize != libp2pCurve.BitSize {
		t.Errorf(
			"unexpected BitSize\nexpected: %v\nactual: %v",
			ethereumCurve.BitSize,
			libp2pCurve.BitSize,
		)
	}
}

func TestSameKeyAsEthereum(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	libp2pKey := toLibp2pKey(ethereumKey)

	if ethereumKey.PrivateKey.D.Cmp(libp2pKey.D) != 0 {
		t.Errorf(
			"unexpected D\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.D,
			libp2pKey.D,
		)
	}

	if ethereumKey.PrivateKey.PublicKey.X.Cmp(libp2pKey.PublicKey.X) != 0 {
		t.Errorf(
			"unexpected X\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.PublicKey.X,
			libp2pKey.PublicKey.X,
		)
	}

	if ethereumKey.PrivateKey.PublicKey.Y.Cmp(libp2pKey.PublicKey.Y) != 0 {
		t.Errorf(
			"unexpected Y\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.PublicKey.Y,
			libp2pKey.PublicKey.Y,
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
