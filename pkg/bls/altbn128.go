package bls

import (
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"math/big"
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

func yFromX(x *big.Int) *big.Int {
    return modSqrt(sum(product(x, x, x), big.NewInt(3)), p)
}

func G1FromInts(x *big.Int, y *big.Int) (*bn256.G1) {
    // TODO error out if ints are over 32 bytes

	m := append(x.Bytes(), y.Bytes())

	return new(bn256.G1).Unmarshal(m)
}

func G1HashToPoint(m []byte) *bn256.G1 {

	one, three := big.NewInt(1), big.NewInt(3)

	h := sha256.Sum256(m)

	x := mod(new(big.Int).SetBytes(h[:]), bn256.P)

	for {
		y := yFromX(x)
		if y != nil {
            return G1FromInts(x, y)
        }

		x.Add(x, one)
	}
}

func ySign(y *big.Int) byte {
    arr := y.Bytes()
    return arr[len(arr-1)] & 1
}

func (g *bn256.G1) Compress() []byte {

    rt := make([]byte, 32)

    // x := g.p.x.Bytes()
    // y := g.p.y.Bytes()

    // for i := len(x)-1; i >= 0; i-- {
    //     rt[i] = x[i]
    // }

	marshalled := g.Marshal()

    for i := 31; i >= 0; i-- {
        rt[i] = marshalled[i]
    }

    mask := ySign(y) << 7

    rt[0] |= mask

    return rt
}

func Decompress(m []byte) *bn256.G1 {

    x := new(big.Int).SetBytes(append([]byte{m[0] & 011111111}, m[1:]))
    y := yFromX(x)

    if ySign(m[0]) != ySign(y){
        y = new(big.Int).minus(bn256.P, y)
    }

    return G1FromInts(x, y)
}
