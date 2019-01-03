package groupselection

import "math/big"

var MinimumStake = big.NewInt(1)

var tokensTotal = big.NewInt(1000000000) // 1 Billion tokens
var spaceTickets = big.NewInt((2 ^ 256) - 1)

func quotientTokensTotalOverMinStake() *big.Int {
	return tokensTotal.Quo(tokensTotal, MinimumStake)
}

func virtualStakersDistribution(n *big.Int) *big.Int {
	return big.NewInt(1).Mul(n, spaceTickets)
}

func NaturalThreshold(n *big.Int) *big.Int {
	return virtualStakersDistribution(n).Quo(
		spaceTickets, quotientTokensTotalOverMinStake(),
	)
}
