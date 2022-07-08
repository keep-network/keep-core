package ethereum

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

type StakeProviderRetriever interface {
	OperatorToStakingProvider(address chain.Address) (chain.Address, error)
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

	// Verify staking provider for the given operator.
	// TODO: Use address obtained from `operatorPublicKeyToChainAddress`
	walletRegistryStakingProvider, err :=
		sm.stakeProviderRetriever.OperatorToStakingProvider(chain.Address(""))
	if err != nil {
		return false, err
	}
	if walletRegistryStakingProvider == "" {
		return false, nil
	}

	// TODO: OperatorToStakingProvider should probably be resolved in both
	// `WalletRegistry` and `RandomBeacon`

	// Ensure the staking provider has an owner
	owner, _, _, err := sm.chain.RolesOf(walletRegistryStakingProvider)
	if err != nil {
		return false, err
	}

	if owner == "" {
		return false, nil
	}

	// TODO: Continue with the implementation
	return true, nil
}
