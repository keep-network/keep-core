package zkp

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
	"math/big"
)

func randomFromMultiplicativeGroup(
	random io.Reader,
	n *big.Int,
) (*big.Int, error) {
	for {
		r, err := rand.Int(random, n)

		if err != nil {
			return nil, err
		}

		if new(big.Int).GCD(nil, nil, r, n).Cmp(big.NewInt(1)) == 0 {
			return r, nil
		}
	}
}

func sum256(data ...[]byte) [sha256.Size]byte {
	accumulator := make([]byte, 0)
	for _, d := range data {
		accumulator = append(accumulator, d...)
	}
	return sha256.Sum256(accumulator)
}

func discreteExp(a, b, c *big.Int) *big.Int {
	if b.Cmp(big.NewInt(0)) == -1 { // b < 0 ?
		ret := new(big.Int).Exp(a, new(big.Int).Neg(b), c)
		return new(big.Int).ModInverse(ret, c)
	}
	return new(big.Int).Exp(a, b, c)
}
