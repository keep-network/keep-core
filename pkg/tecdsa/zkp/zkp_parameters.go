package zkp

import (
	"crypto/elliptic"
	"math/big"
)

// PublicParameters for Gennaro's ZKPs
type PublicParameters struct {
	N      *big.Int
	NTilde *big.Int
	h1     *big.Int
	h2     *big.Int

	q     *big.Int
	curve elliptic.Curve
}
