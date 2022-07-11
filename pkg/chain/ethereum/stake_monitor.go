package ethereum

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

type StakeProviderRetriever interface {
	OperatorToStakingProvider(address chain.Address) (chain.Address, bool, error)
}

type stakeMonitor struct {
	stakeProviderRetriever StakeProviderRetriever

	chain *Chain
}

func newStakeMonitor(stakeProviderRetriever StakeProviderRetriever, chain *Chain) *stakeMonitor {
	return &stakeMonitor{
		stakeProviderRetriever: stakeProviderRetriever,
		chain:                  chain,
	}
}

func (sm *stakeMonitor) HasMinimumStake(
	operatorPublicKey *operator.PublicKey,
) (bool, error) {
	_, err := operatorPublicKeyToChainAddress(operatorPublicKey)
	if err != nil {
		return false, fmt.Errorf(
			"cannot convert from operator key to chain address: [%v]",
			err,
		)
	}

	stakingProvider, isRegistered, err :=
		sm.stakeProviderRetriever.OperatorToStakingProvider(chain.Address(""))
	if err != nil {
		return false, fmt.Errorf("could not resolve staking provider: [%v]", err)
	}

	if !isRegistered {
		return false, fmt.Errorf(
			"operator not registered for the staking provider")
	}

	// TODO: Should OperatorToStakingProvider be resolved in both
	// `WalletRegistry` and `RandomBeacon` at the same time?

	// Ensure the staking provider has an owner
	_, _, _, stakeHashOwner, err := sm.chain.RolesOf(stakingProvider)
	if err != nil {
		return false, err
	}

	if !stakeHashOwner {
		return false, fmt.Errorf("staking provider has no owner set")
	}

	// TODO: Continue with the implementation
	return true, nil
}
