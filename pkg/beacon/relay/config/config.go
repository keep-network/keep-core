package config

import "math/big"

// Chain contains the config data needed for the relay to operate.
type Chain struct {
	// GroupSize is the size of a group in the threshold relay.
	GroupSize int
	// Threshold is the minimum number of interacting group members needed to
	// produce a relay entry.
	Threshold int
	// TicketInitialSubmissionTimeout is the duration (in blocks) the staker has to submit
	// tickets that fall under the natural threshold to satisfy the initial
	// ticket timeout (see group selection, phase 2a).
	TicketInitialSubmissionTimeout int
	// TicketReactiveSubmissionTimeout is the duration (in blocks) the staker has to
	// submit any tickets that did not fall under the natural threshold. This
	// final chance to submit tickets is called reactive ticket submission
	// (defined in the group selection algorithm, 2b).
	TicketReactiveSubmissionTimeout int
	// TicketChallengeTimeout is the duration (in blocks) the staker has to
	// submit any challenges for tickets that fail any checks.
	TicketChallengeTimeout int
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

// Honest threshold is the sufficient amount of valid signature shares required
// to reconstruct group BLS signature after threshold signing.
func (c *Chain) HonestThreshold() int {
	return c.Threshold + 1
}
