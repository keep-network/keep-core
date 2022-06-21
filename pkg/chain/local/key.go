package local

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/keep-network/keep-core/pkg/operator"
)

// OperatorPublicKeyToChainPublicKey converts an operator public key to
// a local-chain-specific public key that uses the btcec secp256k1
// curve implementation under the hood.
func OperatorPublicKeyToChainPublicKey(
	operatorPublicKey *operator.PublicKey,
) (*ecdsa.PublicKey, error) {
	if operatorPublicKey.Curve != operator.Secp256k1 {
		return nil, fmt.Errorf("local chain supports only secp256k1 operator keys")
	}

	return &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     operatorPublicKey.X,
		Y:     operatorPublicKey.Y,
	}, nil
}

// OperatorPrivateKeyToChainKeyPair converts an operator private key to
// a local-chain-specific key pair that uses the btcec secp256k1
// curve implementation under the hood.
func OperatorPrivateKeyToChainKeyPair(
	operatorPrivateKey *operator.PrivateKey,
) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	operatorPublicKey := &operatorPrivateKey.PublicKey

	chainPublicKey, err := OperatorPublicKeyToChainPublicKey(
		operatorPublicKey,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot convert operator public key to chain public key: [%v]",
			err,
		)
	}

	chainPrivateKey := &ecdsa.PrivateKey{
		PublicKey: *chainPublicKey,
		D:         operatorPrivateKey.D,
	}

	return chainPrivateKey, chainPublicKey, nil
}
