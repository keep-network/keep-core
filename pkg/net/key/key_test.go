package key

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/operator"
	btcec_v2 "github.com/btcsuite/btcd/btcec/v2"
)

// `geth` uses `go-ethereum` library to generate key with secp256k1 curve.
// `libp2p` does not recognize this curve and when it comes to creating peer's
// ID or deserializing the key, operation fails with unrecognized curve error.
//
// To overcome this limitation we rewrite ECDSA key referencing secp256k1 curve
// from `go-ethereum` library into a new key instance supported by `libP2P` and
// referencing secp256k1 curve from `btcsuite` used by `libp2p` under the hood.
// This happens in `toLibp2pKey` function.
//
// As long as all curve parameters are the same, this operation is valid.
// This test ensures that secp256k1 from `go-ethereum` and secp256k1 from
// `btcsuite` are the same. If this test starts to fails, we'll need to revisit
// how the key is ported from one instance to another in `toLibp2pKey` function.
func TestSameCurveAsEthereum(t *testing.T) {
	privateKey, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	_, libp2pKey := OperatorKeyToNetworkKey(privateKey, publicKey)
	ethereumCurve := privateKey.Curve.Params()

	libp2pCurve := NetworkKeyToECDSAKey(libp2pKey).Curve.Params()

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
	staticPrivateKey, staticPublicKey, err := operator.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	libp2pPrivKey, libp2pPubKey := OperatorKeyToNetworkKey(
		staticPrivateKey, staticPublicKey,
	)

	ecdsaPrivateKey := (*btcec_v2.PrivateKey)(libp2pPrivKey).ToECDSA()

	if staticPrivateKey.D.Cmp(ecdsaPrivateKey.D) != 0 {
		t.Errorf(
			"unexpected D\nexpected: %v\nactual: %v",
			staticPrivateKey.D,
			ecdsaPrivateKey.D,
		)
	}

	ecdsaPublicKey := NetworkKeyToECDSAKey(libp2pPubKey)

	if staticPublicKey.X.Cmp(ecdsaPublicKey.X) != 0 {
		t.Errorf(
			"unexpected X\nexpected: %v\nactual: %v",
			staticPublicKey.X,
			ecdsaPublicKey.X,
		)
	}

	if staticPrivateKey.PublicKey.Y.Cmp(ecdsaPublicKey.Y) != 0 {
		t.Errorf(
			"unexpected Y\nexpected: %v\nactual: %v",
			staticPrivateKey.PublicKey.Y,
			ecdsaPublicKey.Y,
		)
	}
}

func TestSameKeySerializationMethod(t *testing.T) {
	privateOperatorKey, publicOperatorKey, err := operator.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	_, publicNetworkKey := OperatorKeyToNetworkKey(
		privateOperatorKey, publicOperatorKey,
	)

	operatorKeyBytes := operator.Marshal(publicOperatorKey)
	networkKeyBytes := Marshal(publicNetworkKey)

	testutils.AssertBytesEqual(t, operatorKeyBytes, networkKeyBytes)
}

func TestNetworkPubKeyToAddress(t *testing.T) {
	staticPrivateKey, staticPublicKey, err := operator.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	_, libp2pPublicKey := OperatorKeyToNetworkKey(
		staticPrivateKey, staticPublicKey,
	)

	ecdsaPublicKey := NetworkKeyToECDSAKey(libp2pPublicKey)

	ethAddress := crypto.PubkeyToAddress(*ecdsaPublicKey).String()

	libp2pAddress := NetworkPubKeyToChainAddress(libp2pPublicKey)

	if ethAddress != libp2pAddress {
		t.Errorf(
			"unexpected address\nexpected: %v\nactual: %v",
			ethAddress,
			libp2pAddress,
		)
	}
}

func TestECDSAKeyToNetworkKey(t *testing.T) {
	_, ecdsaPublicKey, err := operator.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	addressFromECDSAKey := crypto.PubkeyToAddress(*ecdsaPublicKey).String()

	networkPublicKey := ECDSAKeyToNetworkKey(ecdsaPublicKey)
	addressFromNetworkKey := NetworkPubKeyToChainAddress(networkPublicKey)

	if addressFromECDSAKey != addressFromNetworkKey {
		t.Errorf(
			"unexpected address\nexpected: %v\nactual: %v",
			addressFromECDSAKey,
			addressFromNetworkKey,
		)
	}
}
