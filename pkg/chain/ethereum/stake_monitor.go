package ethereum

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

// TODO: Remove all

// StakeProviderResolver is a handle for interaction with contracts able
// to resolve staking provider based on the given operator address.
type StakeProviderResolver interface {
	// OperatorToStakingProvider returns the staking provider address for the
	// given operator. If the staking provider has not been registered for the
	// operator, the returned address is empty and the boolean flag is set to
	// false. If the staking provider has been registered, the address is not
	// empty and the boolean flag indicates true.
	OperatorToStakingProvider(operator chain.Address) (chain.Address, bool, error)
}

type stakeMonitor struct {
	stakeProviderResolver StakeProviderResolver

	chain *Chain
}

func newStakeMonitor(stakeProviderResolver StakeProviderResolver, chain *Chain) *stakeMonitor {
	return &stakeMonitor{
		stakeProviderResolver: stakeProviderResolver,
		chain:                 chain,
	}
}

// HasMinimumStake checks if the specified account has enough active stake
// to become network operator and that the operator contract the client is
// working with has been authorized for potential slashing.
//
// Having the required minimum of active stake makes the operator eligible
// to join the network. If the active stake is not currently undelegating,
// operator is also eligible for work selection.
func (sm *stakeMonitor) HasMinimumStake(
	operatorPublicKey *operator.PublicKey,
) (bool, error) {
	operatorAddress, err := operatorPublicKeyToChainAddress(operatorPublicKey)
	if err != nil {
		return false, fmt.Errorf(
			"cannot convert from operator key to chain address: [%v]",
			err,
		)
	}

	stakingProvider, isRegistered, err :=
		sm.stakeProviderResolver.OperatorToStakingProvider(
			chain.Address(operatorAddress.Hex()),
		)
	if err != nil {
		return false, fmt.Errorf("could not resolve staking provider: [%v]", err)
	}

	if !isRegistered {
		return false, fmt.Errorf(
			"staking provider not registered for the given operator address[%s]",
			operatorPublicKey.String(),
		)
	}

	// TODO: Should OperatorToStakingProvider be resolved in both
	// `WalletRegistry` and `RandomBeacon` at the same time?

	// Check if the staking provider has an owner. This check ensures that there
	// is/was a stake delegation for the given staking provider.
	_, _, _, hasStakeDelegation, err := sm.chain.RolesOf(stakingProvider)
	if err != nil {
		return false, err
	}

	if !hasStakeDelegation {
		return false, fmt.Errorf("staking provider has no staking delegation")
	}

	return true, nil
}
