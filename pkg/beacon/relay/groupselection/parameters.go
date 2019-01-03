package groupselection

import "math/big"

// MinimumStake is an on-chain value representing the minimum necessary amount a
// client must lock up to submit a single ticket
var MinimumStake = big.NewInt(1)

var tokensTotal = big.NewInt(1000000000) // 1 Billion tokens
var spaceTickets = big.NewInt((2 ^ 256) - 1)

func quotientTokensTotalOverMinStake() *big.Int {
	return tokensTotal.Quo(tokensTotal, MinimumStake)
}

func virtualStakersDistribution(n *big.Int) *big.Int {
	return big.NewInt(1).Mul(n, spaceTickets)
}

// NaturalThreshold takes the number of virtual stakers a client wishes to submit
// tickets for and computes the natural threshold value defined as:
// floor(N * spaceTickets / (tokensTotal / minimumStake))
func NaturalThreshold(n *big.Int) *big.Int {
	return virtualStakersDistribution(n).Quo(
		spaceTickets, quotientTokensTotalOverMinStake(),
	)
}
