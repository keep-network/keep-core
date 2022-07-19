package local

import (
	"crypto/elliptic"
	"github.com/btcsuite/btcd/btcec"
)

// DefaultCurve is the default elliptic curve implementation used in the
// net/local package. Local network uses the secp256k1 curve and the specific
// implementation is provided by the btcec package.
var DefaultCurve elliptic.Curve = btcec.S256()
