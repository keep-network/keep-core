package sortition

import (
	"errors"
	"math/big"
)

// ErrOperatorNotRegistered ...
var ErrOperatorNotRegistered = errors.New("operator not registered")

// Handle for interaction with the Random Beacon module contracts.
type Handle interface {
	OperatorToStakingProvider() (string, error)
	EligibleStake(stakingProvider string) (*big.Int, error)
	IsOperatorInPool() (bool, error)
	JoinSortitionPool() error
}
