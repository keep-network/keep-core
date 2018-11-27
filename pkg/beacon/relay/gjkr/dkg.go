package gjkr

import (
	crand "crypto/rand"
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/chain"
)

// DKG contains the configuration data needed for the DKG protocol execution.
type DKG struct {
	// P, Q are big primes, such that `p = 2q + 1`
	P, Q *big.Int

	chain *Chain
}

// Chain contains handle to interact with blockchain along with parameters specific
// for block height tracking.
type Chain struct {
	handle chain.Handle

	// Block height when the protocol execution started. Value needs to be set
	// at the begining of the protocol execution.
	initialBlockHeight int // t_init
	// Predefined expected duration of the protocol execution. Relates to DKG
	// Phase 13.
	expectedProtocolDuration int // t_dkg
	// Predefined step for each publishing window. The value is used to determine
	// eligible publishing member. Relates to DKG Phase 13.
	blockStep int // t_step

	TMax      int // M_Max
	TConflict int // T_conflict

	TNow     int // T_Now - Current block height at the current time
	TFirst   int // T_First -
	PerGroup map[string]ValidationGroupState
}

// RandomQ generates a random `big.Int` in range (0, q).
func (d *DKG) RandomQ() (*big.Int, error) {
	for {
		x, err := crand.Int(crand.Reader, d.Q)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number [%s]", err)
		}
		if x.Sign() > 0 {
			return x, nil
		}
	}
}

// ChainHandle returns blockchain handle that provides access to chain interactions.
func (d *DKG) ChainHandle() chain.Handle {
	return d.chain.handle
}

// CurrentBlock returns current block height on a chain.
func (d *Chain) CurrentBlock() (int, error) {
	blockCounter, err := d.handle.BlockCounter()
	if err != nil {
		return 0, err
	}
	return blockCounter.CurrentBlock()
}
