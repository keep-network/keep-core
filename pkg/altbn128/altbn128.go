package altbn128

import (
	"crypto/sha256"
	"errors"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

type G1Point struct {
	*bn256.G1
}

type G2Point struct {
	*bn256.G2
}

// Quadratic extension field element as seen in bn256/gfp2.go
type gfP2 struct {
	x, y *big.Int
}

// Twist curve B constant from bn256/twist.go. This is used to calculate y using y² = x³ + twistB.
var twistB = &gfP2{
	bigFromBase10("19485874751759354771024239261021720505790618469301721065564631296452457478373"),
	bigFromBase10("266929791119991161246907387137283842545076965332900288569378510910307636690"),
}

// Root of the point where x and y are equal. This is used to calculate square root of y in gfP2.
var hexRoot = &gfP2{
	bigFromBase10("21573744529824266246521972077326577680729363968861965890554801909984373949499"),
	bigFromBase10("16854739155576650954933913186877292401521110422362946064090026408937773542853"),
}

// bigFromBase10 returns a big number from it's string representation.
func bigFromBase10(s string) *big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	return n
}

func sum(ints ...*big.Int) *big.Int {
	acc := big.NewInt(0)
	for _, num := range ints {
		acc.Add(acc, num)
	}
	return acc
}

func product(ints ...*big.Int) *big.Int {
	acc := big.NewInt(1)
	for _, num := range ints {
		acc.Mul(acc, num)
	}
	return acc
}

func mod(i, m *big.Int) *big.Int {
	return new(big.Int).Mod(i, m)
}

// modSqrt returns square root of x mod p if such a square root exists. The
// modulus p must be an odd prime. If x is not a square mod p, function returns
// nil.
func modSqrt(x, p *big.Int) *big.Int {
	return new(big.Int).ModSqrt(x, p)
}

// yFromX calculates and returns only one of the two possible Ys, by
// solving the curve equation for X, the two Ys can be distinguished by
// their parity.
func yFromX(x *big.Int) *big.Int {
	return modSqrt(sum(product(x, x, x), big.NewInt(3)), bn256.P)
}

// G1FromInts returns G1 point based on the provided x and y.
func G1FromInts(x *big.Int, y *big.Int) (*bn256.G1, error) {
	if len(x.Bytes()) > 32 || len(y.Bytes()) > 32 {
		return nil, errors.New("points on G1 are limited to 256-bit coordinates")
	}

	paddedX, _ := byteutils.LeftPadTo32Bytes(x.Bytes())
	paddedY, _ := byteutils.LeftPadTo32Bytes(y.Bytes())
	m := append(paddedX, paddedY...)

	g1 := new(bn256.G1)

	_, err := g1.Unmarshal(m)

	return g1, err
}

// G2FromInts returns G2 point based on the provided x and y in Fp^2.
func G2FromInts(x *gfP2, y *gfP2) (*bn256.G2, error) {

	if len(x.x.Bytes()) > 32 || len(x.y.Bytes()) > 32 || len(y.x.Bytes()) > 32 || len(y.y.Bytes()) > 32 {
		return nil, errors.New("points on G2 are limited to two 256-bit coordinates")
	}

	paddedXX, _ := byteutils.LeftPadTo32Bytes(x.x.Bytes())
	paddedXY, _ := byteutils.LeftPadTo32Bytes(x.y.Bytes())
	paddedX := append(paddedXY, paddedXX...)

	paddedYX, _ := byteutils.LeftPadTo32Bytes(y.x.Bytes())
	paddedYY, _ := byteutils.LeftPadTo32Bytes(y.y.Bytes())
	paddedY := append(paddedYY, paddedYX...)

	m := append(paddedX, paddedY...)

	g2 := new(bn256.G2)

	_, err := g2.Unmarshal(m)

	return g2, err
}

// G1HashToPoint hashes the provided byte slice, maps it into a G1
// and returns it as a G1 point.
func G1HashToPoint(m []byte) *bn256.G1 {

	one := big.NewInt(1)

	h := sha256.Sum256(m)

	x := mod(new(big.Int).SetBytes(h[:]), bn256.P)

	for {
		y := yFromX(x)
		if y != nil {
			g1, _ := G1FromInts(x, y)
			return g1
		}

		x.Add(x, one)
	}
}

// yParity calculates whether the provided Y coordinate is an even or odd
// number. Returns 0x01 if Y is an even number and 0x00 if it's odd.
func yParity(y *big.Int) byte {
	arr := y.Bytes()
	return arr[len(arr)-1] & 1
}

// Compress compresses point by using X value and the parity bit of Y
// encoded into the first byte.
func (g G1Point) Compress() []byte {

	rt := make([]byte, 32)

	marshalled := g.Marshal()

	for i := 0; i < 32; i++ {
		rt[i] = marshalled[i]
	}

	y := new(big.Int).SetBytes(marshalled[32:])

	// Encode the parity (even/oddness) of y as a 0 or 1 in the topmost bit of
	// the compressed point.
	mask := yParity(y) << 7
	rt[0] |= mask

	return rt
}

// Compress compresses point by using X value and the parity bit of Y
// encoded into the first byte.
func (g G2Point) Compress() []byte {

	// X of G2 point is a 64 bytes value.
	rt := make([]byte, 64)

	marshalled := g.Marshal()

	for i := 0; i < 64; i++ {
		rt[i] = marshalled[i]
	}

	y := new(big.Int).SetBytes(marshalled[64:96])

	// Encode the parity (even/oddness) of y as a 0 or 1 in the topmost bit of
	// the compressed point.
	mask := yParity(y) << 7
	rt[0] |= mask

	return rt
}

// DecompressToG1 decompresses byte slice into G1 point by extracting Y parity
// bit from the first byte, extracting X value and calculating original Y
// value based on the extracted Y parity. The parity bit is encoded in the
// top byte as 0x01 (even) or 0x00 (odd).
func DecompressToG1(m []byte) (*bn256.G1, error) {

	// Get the original X.
	x := new(big.Int).SetBytes(append([]byte{m[0] & 0x7F}, m[1:]...))

	// Get one of the two possible Y.
	y := yFromX(x)

	if y == nil {
		return nil, errors.New("failed to decompress G1")
	}

	// Compare calculated Y parity with the original Y parity in the top bit of
	// the compressed point. If it doesn't match, we know `Y1 + Y2 = P`, so we
	// recover the correct Y using bn256.P.
	if m[0]&0x80>>7 != yParity(y) {
		y = new(big.Int).Add(bn256.P, new(big.Int).Neg(y))
	}

	return G1FromInts(x, y)
}

// DecompressToG2 decompresses byte slice into G2 point by extracting Y parity
// bit from the first byte, extracting X value and calculating original Y
// value based on the extracted Y parity. The parity bit is encoded in the
// top byte as 0x01 (even) or 0x00 (odd).
func DecompressToG2(m []byte) (*bn256.G2, error) {

	// Get the X.
	x := new(gfP2)
	x.x = new(big.Int).SetBytes(m[32:64])
	// Strip Y parity bit when recovering the upper bytes.
	x.y = new(big.Int).SetBytes(append([]byte{m[0] & 0x7F}, m[1:32]...))

	// Get one of the two possible Y on curve y² = x³ + twistB.
	y2 := new(gfP2).pow(x, big.NewInt(3))
	y2.add(y2, twistB)
	y := sqrtGfP2(y2)

	// Compare calculated Y parity with the original Y parity in the top bit of
	// the compressed point. If it doesn't match, we know `Y1 + Y2 = P`, so we
	// recover the correct Y using bn256.P.
	if m[0]&0x80>>7 != yParity(y.y) {
		y.x = new(big.Int).Add(bn256.P, new(big.Int).Neg(y.x))
		y.y = new(big.Int).Add(bn256.P, new(big.Int).Neg(y.y))
	}

	return G2FromInts(x, y)
}

// multiply returns multiplication of two gfP2 elements.
func (e *gfP2) multiply(a, b *gfP2) *gfP2 {
	xx := mod(new(big.Int).Mul(a.x, b.x), bn256.P)
	xy := mod(new(big.Int).Mul(a.x, b.y), bn256.P)
	yx := mod(new(big.Int).Mul(a.y, b.x), bn256.P)
	yy := mod(new(big.Int).Mul(a.y, b.y), bn256.P)
	e.x = mod(new(big.Int).Sub(xx, yy), bn256.P)
	e.y = mod(new(big.Int).Add(xy, yx), bn256.P)
	return e
}

// add returns addition of two gfP2 elements.
func (e *gfP2) add(a, b *gfP2) *gfP2 {
	e.x = mod(new(big.Int).Add(a.x, b.x), bn256.P)
	e.y = mod(new(big.Int).Add(a.y, b.y), bn256.P)
	return e
}

// x2y compares if y^2 equals x.
func x2y(x, y *gfP2) bool {
	y = new(gfP2).pow(y, big.NewInt(2))
	return y.x.Cmp(x.x) == 0 && y.y.Cmp(x.y) == 0
}

// sqrtGfP2 returns square root of a gfP2 element.
func sqrtGfP2(x *gfP2) *gfP2 {

	// (bn256.p^2 + 15) // 32)
	var exp = bigFromBase10("14971724250519463826312126413021210649976634891596900701138993820439690427699319920245032869357433499099632259837909383182382988566862092145199781964622")

	y := new(gfP2).pow(x, exp)

	// Multiply y by hexRoot constant to find correct y.
	for !x2y(x, y) {
		y.multiply(y, hexRoot)
	}
	return y
}

// pow returns gfP2 element to the power of the provided exponent.
func (e *gfP2) pow(base *gfP2, exp *big.Int) *gfP2 {

	e.x = big.NewInt(1)
	e.y = big.NewInt(0)

	// Reduce exp with right shift operator (divide by 2) gradually to 0
	// while computing e when exp is an odd number.
	for exp.Cmp(big.NewInt(0)) == 1 {

		if yParity(exp) == 1 {
			e.multiply(e, base)
		}

		exp = new(big.Int).Rsh(exp, 1)
		base = new(gfP2).multiply(base, base)
	}
	return e
}
