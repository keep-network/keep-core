package zkp

import (
	"crypto/sha256"
	"math/big"
)

func sum256(data ...[]byte) [sha256.Size]byte {
	accumulator := make([]byte, 0)
	for _, d := range data {
		accumulator = append(accumulator, d...)
	}
	return sha256.Sum256(accumulator)
}

// Evaluates discrete result of a^b mod c. Since we operate on integers,
// negative exponent is interpreted as the multiplicative inverse
// a^b modulo c.
func discreteExp(a, b, c *big.Int) *big.Int {
	if b.Cmp(big.NewInt(0)) == -1 { // b < 0 ?
		ret := new(big.Int).Exp(a, new(big.Int).Abs(b), c)
		return new(big.Int).ModInverse(ret, c)
	}
	return new(big.Int).Exp(a, b, c)
}

// Returns true if number ∈ [start, end).
// Returns false otherwise.
func isInRange(number, start, end *big.Int) bool {
	return number.Cmp(start) != -1 && number.Cmp(end) == -1
}
