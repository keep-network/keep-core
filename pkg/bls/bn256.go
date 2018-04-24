package bls

import (
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare/bn256"
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

func G1HashToPoint(m []byte) (*big.Int, *big.Int) {

	one, three := big.NewInt(1), big.NewInt(3)

	h := sha256.Sum256(m)

	x := mod(new(big.Int).SetBytes(h[:]), bn256.P)

	for {
		x3 = product(x, x, x)
		y := modSqrt(sum(x3, three), p)
		if y != nil {
			return x, y
		}

		x.Add(x, one)
	}
}
