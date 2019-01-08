package tecdsa

import (
	"fmt"
	"math/big"
	"testing"
)

var ONE = big.NewInt(1)

func TestIt(t *testing.T) {
	fmt.Printf("%v\n", getMultiplicativeGroup(100))
}

func getMultiplicativeGroup(n int) []int {
	var multiplicativeGroup []int

	for i := 1; i < n; i++ {
		if isInMultiplicativeGroup(int64(i), int64(n)) {
			multiplicativeGroup = append(multiplicativeGroup, i)
		}
	}

	return multiplicativeGroup
}

func isInMultiplicativeGroup(x int64, n int64) bool {
	return ONE.Cmp(new(big.Int).GCD(
		nil,
		nil,
		big.NewInt(n),
		big.NewInt(x),
	)) == 0
}
