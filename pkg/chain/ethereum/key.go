package ethereum

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/operator"
)

// ChainKeyToOperatorKeyPair converts the Ethereum chain key to a universal
// operator key pair. This conversion decouples the key from the chain-specific
// curve implementation while preserving the curve name.
func ChainKeyToOperatorKeyPair(
	chainKey *keystore.Key,
) (*operator.PrivateKey, *operator.PublicKey, error) {
	chainPrivateKey := chainKey.PrivateKey
	chainPublicKey := &chainPrivateKey.PublicKey

	// Ethereum keys use secp256k1 curve underneath. If the given key doesn't
	// use that curve, something wrong.
	if !isSecp256k1(chainPublicKey.Curve) {
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

// OperatorPublicKeyToChainPublicKey converts an operator public key to
// an Ethereum-specific chain public key that uses the go-ethereum-based
// secp256k1 curve implementation under the hood.
func OperatorPublicKeyToChainPublicKey(
	operatorPublicKey *operator.PublicKey,
) (*ecdsa.PublicKey, error) {
	if operatorPublicKey.Curve != operator.Secp256k1 {
		return nil, fmt.Errorf("ethereum supports only secp256k1 operator keys")
	}

	return &ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     operatorPublicKey.X,
		Y:     operatorPublicKey.Y,
	}, nil
}

// OperatorPublicKeyToChainAddress converts an operator public key to
// an Ethereum-specific chain address.
func OperatorPublicKeyToChainAddress(
	operatorPublicKey *operator.PublicKey,
) (common.Address, error) {
	chainPublicKey, err := OperatorPublicKeyToChainPublicKey(operatorPublicKey)
	if err != nil {
		return common.Address{}, fmt.Errorf(
			"cannot convert from operator to chain key: [%v]",
			err,
		)
	}

	return crypto.PubkeyToAddress(*chainPublicKey), nil
}

func isSecp256k1(curve elliptic.Curve) bool {
	return reflect.DeepEqual(curve.Params(), crypto.S256().Params())
}
