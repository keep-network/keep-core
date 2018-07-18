package curve

import "math/big"

// Point is a helper struct representing point on the Elliptic Curve.
// Neither Go nor Ethereum EC implementations define such a struct and
// all operations are done on coordinates. For T-ECDSA it's quite unwieldy
// and makes the code a bit clumsy. Point does not keep any curve operation
// logic and was created just to improve code readability.
type Point struct {
	X, Y *big.Int
}

func NewPoint(x, y *big.Int) *Point {
	return &Point{x, y}
}

func (cp *Point) Bytes() []byte {
	xb := cp.X.Bytes()
	yb := cp.Y.Bytes()
	return append(xb, yb...)
}
