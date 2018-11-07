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

	chain chain.Handle

	// Blockchain block heigh when the protocol execution started.
	// TODO: Move it to chain.BlockCounter ?
	initialBlockHeight int // t_init
	expectedDuration   int // t_dkg
	blockStep          int // t_step
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
