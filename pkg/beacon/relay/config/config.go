package config

import "math/big"

// Chain contains the config data needed for the relay to operate.
type Chain struct {
	// GroupSize is the size of a group in the threshold relay.
	GroupSize int
	// Threshold is the minimum number of interacting group members needed to
	// produce a relay entry.
	Threshold int
	// TicketTimeout is the duration (in blocks) the staker has to submit
	// their tickets to satisfy the initial ticket submission (see group
	// selection, phase 2a)
	TicketTimeout int
	// MinimumStake is an on-chain value representing the minimum necessary
	// amount a client must lock up to submit a single ticket
	MinimumStake *big.Int
	// TokenSupply represents the total number of tokens that can exist in
	// the system
	TokenSupply *big.Int
	// NaturalThreshold is the value N virtual stakers' tickets would be
	// expected to fall below if the tokens were optimally staked, and the
	// tickets' values were evenly distributed in the domain of the
	// pseudorandom function
	NaturalThreshold *big.Int
}
