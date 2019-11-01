package config

import "math/big"

// Chain contains the config data needed for the relay to operate.
type Chain struct {
	// GroupSize is the size of a group in the threshold relay.
	GroupSize int
	// HonestThreshold is the minimum number of active participants behaving
	// according to the protocol needed to generate a new relay entry.
	HonestThreshold int
	// TicketSubmissionTimeout is the duration (in blocks) the staker has to
	// submit any tickets to candidate to a new group.
	TicketSubmissionTimeout uint64
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
	// RelayEntryTimeout is a timeout in blocks on-chain for a relay
	// entry to be published by the selected group. Blocks are
	// counted from the moment relay request occur.
	RelayEntryTimeout uint64
}

// DishonestThreshold is the maximum number of misbehaving participants for
// which it is still possible to generate a new relay entry.
// Misbehaviour is any misconduct to the protocol, including inactivity.
func (c *Chain) DishonestThreshold() int {
	return c.GroupSize - c.HonestThreshold
}
