package tecdsa

import "math/big"

// CurvePoint is a point on the Elliptic Curve
type CurvePoint struct {
	x, y *big.Int
}
