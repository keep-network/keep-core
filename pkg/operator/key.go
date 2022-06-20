package operator

import (
	"fmt"
	"math/big"
)

// Curve represents an elliptic curve name that denotes the curve used by
// the operator keys.
type Curve int

// Enum that lists all the elliptic curves that can be used with the operator
// keys.
const (
	Undefined Curve = iota
	Secp256k1
)

// Holds the bit size of the supported elliptic curves.
var curveBitSizes = map[Curve]int{
	Secp256k1: 256,
}

// PublicKey represents a public key that identifies the node operator.
// Actually this is an ECDSA public key that is not tied to any particular
// elliptic curve implementation and holds just the curve name instead.
// That trait makes the operator key flexible enough to be accepted by chain
// and network packages that can convert it to the chain and network key
// respectively. Chain and network keys should use their own implementations
// of the curve held by the operator key.
type PublicKey struct {
	Curve Curve
	X, Y  *big.Int
}

// PrivateKey represents an operator private key corresponding to an
// operator public key.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

// MarshalUncompressed marshals the given public key to the 65-byte uncompressed
// form. This implementation is based on the elliptic.Marshal function.
func MarshalUncompressed(publicKey *PublicKey) ([]byte, error) {
	curveBitSize, exists := curveBitSizes[publicKey.Curve]
	if !exists {
		return nil, fmt.Errorf("missing curve bit size information")
	}

	byteLength := (curveBitSize + 7) / 8

	uncompressed := make([]byte, 1+2*byteLength)
	uncompressed[0] = 4 // uncompressed point

	publicKey.X.FillBytes(uncompressed[1 : 1+byteLength])
	publicKey.Y.FillBytes(uncompressed[1+byteLength : 1+2*byteLength])

	return uncompressed, nil
}

// MarshalCompressed marshals the given public key to the 33-byte compressed
// form. This implementation is based on the elliptic.MarshalCompressed
// function.
func MarshalCompressed(publicKey *PublicKey) ([]byte, error) {
	curveBitSize, exists := curveBitSizes[publicKey.Curve]
	if !exists {
		return nil, fmt.Errorf("missing curve bit size information")
	}

	byteLength := (curveBitSize + 7) / 8
	compressed := make([]byte, 1+byteLength)
	compressed[0] = byte(publicKey.Y.Bit(0)) | 2
	publicKey.X.FillBytes(compressed[1:])

	return compressed, nil
}
