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
	TicketInitialSubmissionTimeout uint64
	// TicketReactiveSubmissionTimeout is the duration (in blocks) the staker has to
	// submit any tickets that did not fall under the natural threshold. This
	// final chance to submit tickets is called reactive ticket submission
	// (defined in the group selection algorithm, 2b).
	TicketReactiveSubmissionTimeout uint64
	// TicketChallengeTimeout is the duration (in blocks) the staker has to
	// submit any challenges for tickets that fail any checks.
	TicketChallengeTimeout uint64
	// ResultPublicationBlockStep is the duration (in blocks) that has to pass
	// before group member with the given index is eligible to submit the
	// result.
	// Nth player becomes eligible to submit the result after
	// T_dkg + (N-1) * T_step
	// where T_dkg is time for phases 1-12 to complete and T_step is the result
	// publication block step.
	ResultPublicationBlockStep uint64
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
	// RelayEntryTimeout is a timeout in blocks on-chain for a relay
	// entry to be published by the selected group. Blocks are
	// counted from the moment relay request occur.
	RelayEntryTimeout uint64
}
