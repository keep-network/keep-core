package operator

import (
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
