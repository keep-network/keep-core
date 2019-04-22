package chain

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// Staker represents a single staker in the system. Stakers provide access to a
// unique identifier, a stake value, and a way to attach an event listener for
// stake changes.
type Staker interface {
	// ID returns a unique identifier for the represented stake. Two
	// representations of the same stake should return the same ID.
	ID() relaychain.StakerAddress
	// Stake returns the current stake of this staker according to the connected
	// chain state as a promise. If setup of the promise fails, an error is
	// returned.
	Stake() (*big.Int, error)
	// OnStakeChanged registers a callback to be invoked when the given staker's
	// stake changes.
	OnStakeChanged(func(newStake *big.Int))
}
