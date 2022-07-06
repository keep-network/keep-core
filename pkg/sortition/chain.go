package sortition

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/chain"
)

// Chain handle for interaction with the sortition pool contracts.
type Chain interface {
	OperatorToStakingProvider() (chain.Address, bool, error)
	EligibleStake(stakingProvider chain.Address) (*big.Int, error)
	IsPoolLocked() (bool, error)
	IsOperatorInPool() (bool, error)
	JoinSortitionPool() error
}
