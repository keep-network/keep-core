package config

import (
	"math/big"
)

// Chain contains the config data needed for the relay to operate.
type Chain struct {
	// GroupSize is the size of a group in the threshold relay.
	GroupSize int
	// Threshold is the minimum number of interacting group members needed to
	// produce a relay entry.
	Threshold int
}

// DKG contains the configuration data needed for the DKG protocol execution.
type DKG struct {
	// P, Q are big primes, such that `p = 2q + 1`
	P, Q *big.Int
}
