package ethereum

import (
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

type stakeMonitor struct {
	chain *Chain
}

func newStakeMonitor(chain *Chain) *stakeMonitor {
	return &stakeMonitor{chain: chain}
}

// TODO: Real implementation with v2 contracts.
func (sm *stakeMonitor) HasMinimumStake(
	operatorPublicKey *operator.PublicKey,
) (bool, error) {
	return true, nil
}

func (c *Chain) StakeMonitor() (chain.StakeMonitor, error) {
	return newStakeMonitor(c), nil
}
