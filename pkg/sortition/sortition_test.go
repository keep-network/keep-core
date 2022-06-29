package sortition

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/operator"
)

// Start `waitUntilRegistered` when an operator is already registered.
func TestWaitUntilRegistered_AlreadyRegistered(t *testing.T) {
	operator, _ := generateAddress()
	stakingProvider, _ := generateAddress()
	blockCounter, _ := local.BlockCounter()
	localChain := connect(operator)

	operatorRegisteredChan := make(chan string)

	expectedOperatorToStakingProviderAttempts := 1

	localChain.registerOperator(stakingProvider, operator)

	resolveWait := &sync.WaitGroup{}
	resolveWait.Add(1)
	go func() {
		waitUntilRegistered(context.Background(), blockCounter, localChain, operatorRegisteredChan)
		resolveWait.Done()
	}()

	actualStakingProvider := <-operatorRegisteredChan
	if actualStakingProvider != stakingProvider.Hex() {
		t.Errorf(
			"unexpected staking provider\nexpected: [%s]\nactual:   [%s]",
			stakingProvider,
			actualStakingProvider,
		)
	}

	if localChain.operatorToStakingProviderAttempts != expectedOperatorToStakingProviderAttempts {
		t.Errorf(
			"unexpected OperatorToStakingProvider attempts\nexpected: [%d]\nactual:   [%d]",
			expectedOperatorToStakingProviderAttempts,
			localChain.operatorToStakingProviderAttempts,
		)
	}

	// Wait for the `waitResolve` to resolve to verify if the `waitUntilRegistered`
	// function completes execution.
	resolveWait.Wait()

}

// Start `waitUntilRegistered` when an operator has not been registered yet.
// Register the operator when the waiting loop is running.
func TestWaitUntilRegistered_NotRegistered(t *testing.T) {
	operator, _ := generateAddress()
	stakingProvider, _ := generateAddress()
	blockCounter, _ := local.BlockCounter()
	localChain := connect(operator)

	operatorRegisteredChan := make(chan string)

	// Overwrite the default delay to make things happen faster in the test.
	operatorRegistrationRetryDelay = 1 * time.Second
	// Time delay after which the operator gets registered in the test.
	registrationTestDelay := 3 * time.Second
	// Each registration check happens every `operatorRegistrationRetryDelay`.
	// We expect that the operator will be registered after `registrationTestDelay`.
	expectedOperatorToStakingProviderAttempts := 4

	resolveWait := &sync.WaitGroup{}
	resolveWait.Add(1)
	go func() {
		waitUntilRegistered(context.Background(), blockCounter, localChain, operatorRegisteredChan)
		resolveWait.Done()
	}()

	go func() {
		time.Sleep(registrationTestDelay)
		localChain.registerOperator(stakingProvider, operator)
	}()

	actualStakingProvider := <-operatorRegisteredChan
	if actualStakingProvider != stakingProvider.Hex() {
		t.Errorf(
			"unexpected staking provider\nexpected: [%s]\nactual:   [%s]",
			stakingProvider,
			actualStakingProvider,
		)
	}

	if localChain.operatorToStakingProviderAttempts != expectedOperatorToStakingProviderAttempts {
		t.Errorf(
			"unexpected OperatorToStakingProvider attempts\nexpected: [%d]\nactual:   [%d]",
			expectedOperatorToStakingProviderAttempts,
			localChain.operatorToStakingProviderAttempts,
		)
	}

	// Wait for the `waitResolve` to resolve to verify if the `waitUntilRegistered`
	// function completes execution.
	resolveWait.Wait()
}

func TestJoinSortitionPoolWhenEligible(t *testing.T) {
	operator, _ := generateAddress()
	stakingProvider, _ := generateAddress()
	blockCounter, _ := local.BlockCounter()
	localChain := connect(operator)

	localChain.registerOperator(stakingProvider, operator)

	// Overwrite the default delay to make things happen faster in the test.
	eligibilityRetryDelay = 1 * time.Second
	// Time delay after which the operator gets eligible in the test.
	eligibilityTestDelay := 2 * time.Second
	// Each eligibility check happens every `eligibilityRetryDelay`.
	// We expect that the operator will join the pool after `eligibilityTestDelay`.
	expectedEligibleStakeAttempts := 3

	resolveWait := &sync.WaitGroup{}
	resolveWait.Add(1)
	go func() {
		joinSortitionPoolWhenEligible(
			context.Background(),
			stakingProvider.Hex(),
			blockCounter,
			localChain,
		)
		resolveWait.Done()
	}()

	isInPool, _ := localChain.IsOperatorInPool()
	if isInPool {
		t.Fatalf("operator is already registered in the pool")
	}

	go func() {
		time.Sleep(eligibilityTestDelay)
		localChain.eligibleStakes[stakingProvider.Hex()] = big.NewInt(100)
	}()

	// Wait for the `waitResolve` to resolve to verify if the `waitUntilRegistered`
	// function completes execution.
	resolveWait.Wait()

	isInPool, _ = localChain.IsOperatorInPool()
	if !isInPool {
		t.Errorf("operator was not registered in the pool")
	}

	if localChain.eligibleStakeAttempts != expectedEligibleStakeAttempts {
		t.Errorf(
			"unexpected EligibleStake attempts\nexpected: [%d]\nactual:   [%d]",
			expectedEligibleStakeAttempts,
			localChain.eligibleStakeAttempts,
		)
	}
}

func generateAddress() (common.Address, error) {
	_, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		return common.Address{}, err
	}
	return operator.PubkeyToAddress(*publicKey), nil
}
