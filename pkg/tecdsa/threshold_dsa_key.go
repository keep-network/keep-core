package tecdsa

import "math/big"

// ThresholdDsaKey represents DSA key share for a fully initialized Signer
type ThresholdDsaKey struct {
	xi *big.Int
	y  *big.Int
}

type dsaKeyShare struct {
	xi *big.Int
	yi *CurvePoint
}
