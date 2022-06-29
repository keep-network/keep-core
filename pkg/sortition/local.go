package sortition

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type localChain struct {
	operatorAddress common.Address
	eligibleStakes  map[string]*big.Int

	operatorToStakingProvider      map[string]common.Address
	operatorToStakingProviderMutex sync.RWMutex

	sortitionPool      map[string]*big.Int
	sortitionPoolMutex sync.RWMutex

	operatorToStakingProviderAttempts int
	eligibleStakeAttempts             int
}

func connect(operatorAddress common.Address) *localChain {
	return &localChain{
		operatorAddress:           operatorAddress,
		eligibleStakes:            make(map[string]*big.Int),
		operatorToStakingProvider: make(map[string]common.Address),
		sortitionPool:             make(map[string]*big.Int),
	}
}

// This is a test util function to setup the chain
func (lc *localChain) registerOperator(stakingProvider common.Address, operator common.Address) {
	lc.operatorToStakingProvider[lc.operatorAddress.Hex()] = stakingProvider
}

func (lc *localChain) OperatorToStakingProvider() (string, error) {
	lc.operatorToStakingProviderAttempts++

	lc.operatorToStakingProviderMutex.RLock()
	defer lc.operatorToStakingProviderMutex.RUnlock()

	stakingProvider, ok := lc.operatorToStakingProvider[lc.operatorAddress.Hex()]
	if !ok {
		return "", ErrOperatorNotRegistered
	}

	return stakingProvider.Hex(), nil
}

func (lc *localChain) EligibleStake(stakingProvider string) (*big.Int, error) {
	lc.eligibleStakeAttempts++

	eligibleStake, ok := lc.eligibleStakes[stakingProvider]
	if !ok {
		return big.NewInt(0), nil
	}

	return eligibleStake, nil
}

func (lc *localChain) IsOperatorInPool() (bool, error) {
	lc.sortitionPoolMutex.RLock()
	defer lc.sortitionPoolMutex.RUnlock()

	_, ok := lc.sortitionPool[lc.operatorAddress.Hex()]

	return ok, nil
}

func (lc *localChain) JoinSortitionPool() error {
	lc.operatorToStakingProviderMutex.Lock()
	defer lc.operatorToStakingProviderMutex.Unlock()
	lc.sortitionPoolMutex.Lock()
	defer lc.sortitionPoolMutex.Unlock()

	stakingProvider, ok := lc.operatorToStakingProvider[lc.operatorAddress.Hex()]
	if !ok {
		return fmt.Errorf("Unknown operator")
	}

	eligibleStake, ok := lc.eligibleStakes[stakingProvider.Hex()]
	if !ok || eligibleStake.Cmp(big.NewInt(0)) == 0 {
		return fmt.Errorf("Authorization below the minimum")
	}

	_, ok = lc.sortitionPool[lc.operatorAddress.Hex()]
	if ok {
		return fmt.Errorf("Operator is already registered in the pool")
	}

	lc.sortitionPool[lc.operatorAddress.Hex()] = eligibleStake

	return nil
}
