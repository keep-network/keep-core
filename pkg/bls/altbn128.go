package bls

import (
	"crypto/sha256"
	"errors"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

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

func modSqrt(i, m *big.Int) *big.Int {
	return new(big.Int).ModSqrt(i, m)
}

// Returns only one of the two possible Y (even/odd) for provided X
func yFromX(x *big.Int) *big.Int {
	return modSqrt(sum(product(x, x, x), big.NewInt(3)), bn256.P)
}

func G1FromInts(x *big.Int, y *big.Int) (*bn256.G1, error) {
	if len(x.Bytes()) > 32 || len(y.Bytes()) > 32 {
		return nil, errors.New("Points on G1 are limited to 256-bit coordinates.")
	}

	paddedX, _ := byteutils.LeftPadTo32Bytes(x.Bytes())
	paddedY, _ := byteutils.LeftPadTo32Bytes(y.Bytes())
	m := append(paddedX, paddedY...)

	g1 := new(bn256.G1)

	_, err := g1.Unmarshal(m)

	return g1, err
}

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

// Calculates whether the provided Y coordinate is an even or odd number.
// Returns 0x01 if Y is an even number and 0x00 if it's odd.
func ySign(y *big.Int) byte {
	arr := y.Bytes()
	return arr[len(arr)-1] & 1
}

// Compress point by using X value and the sign of Y (even/odd) encoded
// into the first byte. This is achieved by using bitwise `OR` operator.
// Original byte and the sign can be later extracted with bitswise `AND`.
func Compress(g *bn256.G1) []byte {

	rt := make([]byte, 32)

	marshalled := g.Marshal()

	for i := 31; i >= 0; i-- {
		rt[i] = marshalled[i]
	}

	y := new(big.Int).SetBytes(marshalled[32:])

	// Prepare bytes mask with (even/odd) sign
	mask := ySign(y) << 7

	// Use `OR` operator to save the sign
	rt[0] |= mask

	return rt
}

// Decompress byte slice into G1 point by extracting Y sign from the first
// byte, extracting X value and calculating original Y value based on the
// extracted Y sign. Extracting sign and original byte is achieved by
// using `AND` bitwise operator:
// ... & 0x7F => to get the original byte
// ... & 0x80 => to get the sign
func Decompress(m []byte) (*bn256.G1, error) {

	// Get the original X
	x := new(big.Int).SetBytes(append([]byte{m[0] & 0x7F}, m[1:]...))

	// Get one of the two possible Y
	y := yFromX(x)

	if y == nil {
		return nil, errors.New("Failed to decompress G1.")
	}

	// Compare calculated Y sign with the original Y sign and if it doesn't match
	// get the right Y by extracting the calculated one from the bn256.P
	if m[0] & 0x80 >> 7 != ySign(y) {
		y = new(big.Int).Add(bn256.P, new(big.Int).Neg(y))
	}

	return G1FromInts(x, y)
}
