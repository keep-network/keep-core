package ethereum

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/operator"
)

// DefaultCurve is the default elliptic curve implementation used in the
// chain/ethereum package. Ethereum uses the secp256k1 curve and the specific
// implementation is provided by the go-ethereum package.
var DefaultCurve elliptic.Curve = crypto.S256()

// ChainPrivateKeyToOperatorKeyPair converts the Ethereum chain private key to
// a universal operator key pair. This conversion decouples the key from the
// chain-specific curve implementation while preserving the curve.
func ChainPrivateKeyToOperatorKeyPair(
	chainPrivateKey *ecdsa.PrivateKey,
) (*operator.PrivateKey, *operator.PublicKey, error) {
	chainPublicKey := &chainPrivateKey.PublicKey

	// Ethereum keys use secp256k1 curve underneath. If the given key doesn't
	// use that curve, something is wrong.
	if !isSameCurve(chainPublicKey.Curve, DefaultCurve) {
		return nil, nil, fmt.Errorf("ethereum chain key does not use secp256k1 curve")
	}

	publicKey := &operator.PublicKey{
		Curve: operator.Secp256k1,
		X:     chainPublicKey.X,
		Y:     chainPublicKey.Y,
	}

	privateKey := &operator.PrivateKey{
		PublicKey: *publicKey,
		D:         chainPrivateKey.D,
	}

	return privateKey, publicKey, nil
}

// operatorPublicKeyToChainPublicKey converts an operator public key to
// an Ethereum-specific chain public key that uses the go-ethereum-based
// secp256k1 curve implementation under the hood.
func operatorPublicKeyToChainPublicKey(
	operatorPublicKey *operator.PublicKey,
) (*ecdsa.PublicKey, error) {
	if operatorPublicKey.Curve != operator.Secp256k1 {
		return nil, fmt.Errorf("ethereum supports only secp256k1 operator keys")
	}

	return &ecdsa.PublicKey{
		Curve: DefaultCurve,
		X:     operatorPublicKey.X,
		Y:     operatorPublicKey.Y,
	}, nil
}

// operatorPublicKeyToChainAddress converts an operator public key to
// an Ethereum-specific chain address.
func operatorPublicKeyToChainAddress(
	operatorPublicKey *operator.PublicKey,
) (common.Address, error) {
	chainPublicKey, err := operatorPublicKeyToChainPublicKey(operatorPublicKey)
	if err != nil {
		return common.Address{}, fmt.Errorf(
			"cannot convert from operator to chain key: [%v]",
			err,
		)
	}

	return crypto.PubkeyToAddress(*chainPublicKey), nil
}

func isSameCurve(curve1, curve2 elliptic.Curve) bool {
	return reflect.DeepEqual(curve1.Params(), curve2.Params())
}
