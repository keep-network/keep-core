package sortition

import (
	"errors"
	"math/big"
)

// ErrOperatorNotRegistered is an error that is returned by the sortition
// Handle implementation when an operator has not been registered for a staking
// provider.
var ErrOperatorNotRegistered = errors.New("operator not registered")

// Handle for interaction with the sortition pool contracts.
type Handle interface {
	OperatorToStakingProvider() (string, error)
	EligibleStake(stakingProvider string) (*big.Int, error)
	IsPoolLocked() (bool, error)
	IsOperatorInPool() (bool, error)
	JoinSortitionPool() error
}
