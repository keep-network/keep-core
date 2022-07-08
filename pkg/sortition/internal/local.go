package sortition

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/chain"
)

var errOperatorUnknown = fmt.Errorf("operator not registered for the staking provider")
var errAuthorizationBelowMinimum = fmt.Errorf("authorization below the minimum")
var errOperatorAlreadyRegisteredInPool = fmt.Errorf("operator is already registered in the pool")

type localChain struct {
	operatorAddress chain.Address

	operatorToStakingProvider      map[chain.Address]chain.Address
	operatorToStakingProviderMutex sync.RWMutex

	// Weight of an operator, as set in the Sortition Pool contract.
	sortitionPool      map[chain.Address]*big.Int
	sortitionPoolMutex sync.RWMutex

	// Stake for a staking provider, as set in the Token Staking contract,
	// minus the pending authorization decrease.
	eligibleStake      map[chain.Address]*big.Int
	eligibleStakeMutex sync.RWMutex

	isPoolLocked bool
}

func connectLocal(operatorAddress chain.Address) *localChain {
	return &localChain{
		operatorAddress:           operatorAddress,
		operatorToStakingProvider: make(map[chain.Address]chain.Address),
		sortitionPool:             make(map[chain.Address]*big.Int),
		eligibleStake:             make(map[chain.Address]*big.Int),
	}
}

// This is a test util function to setup the chain
func (lc *localChain) registerOperator(stakingProvider chain.Address, operator chain.Address) {
	lc.operatorToStakingProviderMutex.Lock()
	defer lc.operatorToStakingProviderMutex.Unlock()

	lc.operatorToStakingProvider[lc.operatorAddress] = stakingProvider
}

// This is a test util function to setup the chain
func (lc *localChain) setEligibleStake(stakingProvider chain.Address, stake *big.Int) {
	lc.eligibleStakeMutex.Lock()
	defer lc.eligibleStakeMutex.Unlock()

	lc.eligibleStake[stakingProvider] = stake
}

func (lc *localChain) OperatorToStakingProvider() (chain.Address, bool, error) {
	lc.operatorToStakingProviderMutex.RLock()
	defer lc.operatorToStakingProviderMutex.RUnlock()

	stakingProvider, ok := lc.operatorToStakingProvider[lc.operatorAddress]
	return stakingProvider, ok, nil
}

func (lc *localChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	lc.eligibleStakeMutex.RLock()
	defer lc.eligibleStakeMutex.RUnlock()

	eligibleStake, ok := lc.eligibleStake[stakingProvider]
	if !ok {
		return big.NewInt(0), nil
	}

	return eligibleStake, nil
}

func (lc *localChain) IsPoolLocked() (bool, error) {
	return lc.isPoolLocked, nil
}

func (lc *localChain) IsOperatorInPool() (bool, error) {
	lc.sortitionPoolMutex.RLock()
	defer lc.sortitionPoolMutex.RUnlock()

	_, ok := lc.sortitionPool[lc.operatorAddress]

	return ok, nil
}

func (lc *localChain) IsOperatorUpToDate() (bool, error) {
	lc.sortitionPoolMutex.RLock()
	defer lc.sortitionPoolMutex.RUnlock()

	lc.eligibleStakeMutex.RLock()
	defer lc.eligibleStakeMutex.RUnlock()

	stakingProvider, isRegistered := lc.operatorToStakingProvider[lc.operatorAddress]
	if !isRegistered {
		return false, errOperatorUnknown
	}

	eligibleStake, hasStake := lc.eligibleStake[stakingProvider]
	weight, isInPool := lc.sortitionPool[lc.operatorAddress]

	if !isInPool {
		return !hasStake || eligibleStake.Cmp(big.NewInt(0)) == 0, nil
	} else {
		return weight.Cmp(eligibleStake) == 0, nil
	}
}

func (lc *localChain) JoinSortitionPool() error {
	lc.operatorToStakingProviderMutex.Lock()
	defer lc.operatorToStakingProviderMutex.Unlock()

	lc.sortitionPoolMutex.Lock()
	defer lc.sortitionPoolMutex.Unlock()

	stakingProvider, ok := lc.operatorToStakingProvider[lc.operatorAddress]
	if !ok {
		return errOperatorUnknown
	}

	eligibleStake, ok := lc.eligibleStake[stakingProvider]
	if !ok || eligibleStake.Cmp(big.NewInt(0)) == 0 {
		return errAuthorizationBelowMinimum
	}

	_, ok = lc.sortitionPool[lc.operatorAddress]
	if ok {
		return errOperatorAlreadyRegisteredInPool
	}

	lc.sortitionPool[lc.operatorAddress] = eligibleStake

	return nil
}

func (lc *localChain) UpdateOperatorStatus() error {
	lc.eligibleStakeMutex.RLock()
	defer lc.eligibleStakeMutex.RUnlock()

	lc.sortitionPoolMutex.Lock()
	defer lc.sortitionPoolMutex.Unlock()

	stakingProvider, isRegistered := lc.operatorToStakingProvider[lc.operatorAddress]
	if !isRegistered {
		return errOperatorUnknown
	}

	eligibleStake := lc.eligibleStake[stakingProvider]
	lc.sortitionPool[lc.operatorAddress] = eligibleStake

	return nil
}
