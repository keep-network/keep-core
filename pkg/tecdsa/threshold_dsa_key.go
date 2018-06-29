package tecdsa

import (
	"math/big"

	"github.com/keep-network/paillier"
)

// ThresholdDsaKey represents DSA key for a fully initialized Signer.
//
// For (x, y = g^x) DSA key pair, y and E(x) are made public, where E is an
// additively homomorphic encryption scheme. This is an implicit (t, n) secret
// sharing of x since the decryption key of E is shared among the Signers.
type ThresholdDsaKey struct {
	x *paillier.Cypher
	y *CurvePoint
}

// dsaKeyShare represents a share of DSA key owned by LocalSigner before
// it's fully initialized into Signer.
//
// Each LocalSigner i generates a share of x and y DSA key, xi and yi.
// yi is broadcasted to other signers along with E(xi) where E is an additively
// homomorphic encryption scheme. It lets to compute
//
// E(x) = E(x1) + E(x2) + ... E(xn)
// y = y1 + y2 + ... + yn
//
// to create a ThresholdDsaKey.
type dsaKeyShare struct {
	xi *big.Int
	yi *CurvePoint
}
