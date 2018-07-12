package tecdsa

import "math/big"

// CurvePoint is a helper struct representing point on the Elliptic Curve.
// Neither Go nor Ethereum EC implementations do not define such a struct and
// all operations are done on coordinates. For T-ECDSA it's quite unwieldy
// and makes the code a bit clumsy. CurvePoint does not keep any curve operation
// logic and was created just to improve code readability.
type CurvePoint struct {
	X, Y *big.Int
}

func NewCurvePoint(x, y *big.Int) *CurvePoint {
	return &CurvePoint{x, y}
}

func (cp *CurvePoint) Bytes() []byte {
	xb := cp.X.Bytes()
	yb := cp.Y.Bytes()
	return append(xb, yb...)
}
