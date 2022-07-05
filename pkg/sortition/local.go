package sortition

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/chain"
)

var errAuthorizationBelowMinimum = fmt.Errorf("authorization below the minimum")
var errOperatorAlreadyRegisteredInPool = fmt.Errorf("operator is already registered in the pool")

type localChain struct {
	operatorAddress chain.Address

	operatorToStakingProvider      map[chain.Address]chain.Address
	operatorToStakingProviderMutex sync.RWMutex

	// Stake for an operator, as set in the Sortition Pool contract.
	sortitionPool      map[chain.Address]*big.Int
	sortitionPoolMutex sync.RWMutex

	// Stake for an operator, as set in the Token Staking contract,
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

func (lc *localChain) OperatorToStakingProvider() (chain.Address, error) {
	lc.operatorToStakingProviderMutex.RLock()
	defer lc.operatorToStakingProviderMutex.RUnlock()

	stakingProvider, ok := lc.operatorToStakingProvider[lc.operatorAddress]
	if !ok {
		return "", ErrOperatorNotRegistered
	}

	return stakingProvider, nil
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

func (lc *localChain) JoinSortitionPool() error {
	lc.operatorToStakingProviderMutex.Lock()
	defer lc.operatorToStakingProviderMutex.Unlock()

	lc.sortitionPoolMutex.Lock()
	defer lc.sortitionPoolMutex.Unlock()

	stakingProvider, ok := lc.operatorToStakingProvider[lc.operatorAddress]
	if !ok {
		return ErrOperatorNotRegistered
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
