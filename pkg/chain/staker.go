package chain

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// Staker represents a single staker in the system.
// DEPRECATED
// TODO: The "staker" should probably become "operator" to reflect random
//       beacon v2 structure.
type Staker interface {
	// Address returns staker's address
	Address() relaychain.StakerAddress
	// Stake returns the current stake of this staker according to the connected
	// chain state as a promise. If setup of the promise fails, an error is
	// returned.
	Stake() (*big.Int, error)
}
