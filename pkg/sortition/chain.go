package sortition

import (
	"errors"
	"math/big"

	"github.com/keep-network/keep-core/pkg/chain"
)

// ErrOperatorNotRegistered is an error that is returned by the sortition
// Handle implementation when an operator has not been registered for a staking
// provider. Each chain may have its own representation of a zero-address, and
// this error provides a common representation of an unknown operator.
var ErrOperatorNotRegistered = errors.New("operator not registered")

// Chain handle for interaction with the sortition pool contracts.
type Chain interface {
	OperatorToStakingProvider() (chain.Address, error)
	EligibleStake(stakingProvider chain.Address) (*big.Int, error)
	IsPoolLocked() (bool, error)
	IsOperatorInPool() (bool, error)
	JoinSortitionPool() error
}
