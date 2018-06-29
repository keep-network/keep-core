package tecdsa

import "math/big"

// CurvePoint is a helper struct representing point on the Elliptic Curve.
// Neither Go nor Ethereum EC implementations do not define such a struct and
// all operations are done on coordinates. For T-ECDSA it's quite unwieldy
// and makes the code a bit clumsy. CurvePoint does not keep any logic
// and was created just to improve code readability.
type CurvePoint struct {
	x, y *big.Int
}
