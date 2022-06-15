package ethereum

import (
	"crypto/elliptic"
	"fmt"
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

func isSecp256k1(curve elliptic.Curve) bool {
	return reflect.DeepEqual(curve.Params(), crypto.S256().Params())
}
