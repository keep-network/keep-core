package membership

import (
	"github.com/btcsuite/btcd/btcec"
)

// Staker represents an on-chain identity and staked amount.
type Staker struct {
	PubKey btcec.PublicKey // Q_j

	// A staker's VirtualStakers is how many minimum-stake stakers a given
	// actual staker could form if they were to blitzpants their stake.
	VirtualStakers uint64
}

func NewStaker(pubKey btcec.PublicKey, weight uint64) *Staker {
	return &Staker{
		PubKey:         pubKey,
		VirtualStakers: weight,
	}
}
