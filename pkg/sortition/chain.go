package sortition

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/chain"
)

// Chain handle for interaction with the sortition pool contracts.
type Chain interface {
	// OperatorToStakingProvider returns the staking provider address for the
	// operator. If the staking provider has not been registered for the
	// operator, the returned address is empty and the boolean flag is set to
	// false. If the staking provider has been registered, the address is not
	// empty and the boolean flag indicates true.
	OperatorToStakingProvider() (chain.Address, bool, error)

	// EligibleStake returns the current value of the staking provider's
	// eligible stake. Eligible stake is defined as the currently authorized
	// stake minus the pending authorization decrease. Eligible stake
	// is what is used for operator's weight in the sortition pool.
	// If the authorized stake minus the pending authorization decrease
	// is below the minimum authorization, eligible stake is 0.
	EligibleStake(stakingProvider chain.Address) (*big.Int, error)

	// IsPoolLocked returns true if the sortition pool is locked and no state
	// changes are allowed.
	IsPoolLocked() (bool, error)

	// IsOperatorInPool returns true if the operator is registered in
	// the sortition pool.
	IsOperatorInPool() (bool, error)

	// IsOperatorUpToDate checks if the operator's authorized stake is in sync
	// with operator's weight in the sortition pool.
	// If the operator's authorized stake is not in sync with sortition pool
	// weight, function returns false.
	// If the operator is not in the sortition pool and their authorized stake
	// is non-zero, function returns false.
	IsOperatorUpToDate() (bool, error)

	// JoinSortitionPool executes a transaction to have the operator join the
	// sortition pool.
	JoinSortitionPool() error

	// UpdateOperatorStatus executes a transaction to update the operator's
	// state in the sortition pool.
	UpdateOperatorStatus() error

	// IsEligibleForRewards checks whether the operator is eligible for rewards
	// or not.
	IsEligibleForRewards() (bool, error)

	// Checks whether the operator is able to restore their eligibility for
	// rewards right away.
	CanRestoreRewardEligibility() (bool, error)

	// Restores reward eligibility for the operator.
	RestoreRewardEligibility() error

	// Returns true if the chaosnet phase is active, false otherwise.
	IsChaosnetActive() (bool, error)

	// Returns true if operator is a beta operator, false otherwise.
	// Chaosnet status does not matter.
	IsBetaOperator() (bool, error)

	// GetOperatorID returns the operator ID for the given operator address.
	GetOperatorID(operatorAddress chain.Address) (chain.OperatorID, error)
}
