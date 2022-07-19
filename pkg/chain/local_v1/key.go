package local_v1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/keep-network/keep-core/pkg/operator"
)

// DefaultCurve is the default elliptic curve implementation used in the
// chain/local package. The local chain uses the secp256k1 curve and the
// specific implementation is provided by the btcec package. Normally,
// the go-ethereum implementation could be used as well but that implementation
// has an unset elliptic curve name (elliptic.Curve.Params().Name
// returns an empty string) so that curve cannot be easily used with
// operator.GenerateKeyPair which is widely used in the local chain
// environment. This is why the btcec implementation is used.
var DefaultCurve elliptic.Curve = btcec.S256()

// operatorPublicKeyToChainPublicKey converts an operator public key to
// a local-chain-specific public key that uses the btcec secp256k1
// curve implementation under the hood.
func operatorPublicKeyToChainPublicKey(
	operatorPublicKey *operator.PublicKey,
) (*ecdsa.PublicKey, error) {
	if operatorPublicKey.Curve != operator.Secp256k1 {
		return nil, fmt.Errorf("local chain supports only secp256k1 operator keys")
	}

	return &ecdsa.PublicKey{
		Curve: DefaultCurve,
		X:     operatorPublicKey.X,
		Y:     operatorPublicKey.Y,
	}, nil
}

// operatorPrivateKeyToChainKeyPair converts an operator private key to
// a local-chain-specific key pair that uses the btcec secp256k1
// curve implementation under the hood.
func operatorPrivateKeyToChainKeyPair(
	operatorPrivateKey *operator.PrivateKey,
) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	operatorPublicKey := &operatorPrivateKey.PublicKey

	chainPublicKey, err := operatorPublicKeyToChainPublicKey(
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
