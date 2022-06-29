package sortition

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/operator"
)

// Start `waitUntilRegistered` when an operator is already registered.
func TestWaitUntilRegistered_AlreadyRegistered(t *testing.T) {
	err := log.SetLogLevel("*", "DEBUG")
	if err != nil {
		t.Fatalf("logger initialization failed: [%v]", err)
	}

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

	if localChain.operatorToStakingProviderAttempts != expectedOperatorToStakingProviderAttempts {
		t.Errorf(
			"unexpected OperatorToStakingProvider attempts\nexpected: [%d]\nactual:   [%d]",
			expectedOperatorToStakingProviderAttempts,
			localChain.operatorToStakingProviderAttempts,
		)
	}

	if actualStakingProvider != stakingProvider.Hex() {
		t.Errorf(
			"unexpected staking provider\nexpected: [%s]\nactual:   [%s]",
			stakingProvider,
			actualStakingProvider,
		)
	}

	// Wait for the `waitResolve` to resolve to verify if the `waitUntilRegistered`
	// function completes execution.
	resolveWait.Wait()
}

// Start `waitUntilRegistered` when an operator has not been registered yet.
// Register the operator when the waiting loop is running.
func TestWaitUntilRegistered_NotRegistered(t *testing.T) {
	err := log.SetLogLevel("*", "DEBUG")
	if err != nil {
		t.Fatalf("logger initialization failed: [%v]", err)
	}

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

	if localChain.operatorToStakingProviderAttempts != expectedOperatorToStakingProviderAttempts {
		t.Errorf(
			"unexpected OperatorToStakingProvider attempts\nexpected: [%d]\nactual:   [%d]",
			expectedOperatorToStakingProviderAttempts,
			localChain.operatorToStakingProviderAttempts,
		)
	}

	if actualStakingProvider != stakingProvider.Hex() {
		t.Errorf(
			"unexpected staking provider\nexpected: [%s]\nactual:   [%s]",
			stakingProvider,
			actualStakingProvider,
		)
	}

	// Wait for the `waitResolve` to resolve to verify if the `waitUntilRegistered`
	// function completes execution.
	resolveWait.Wait()
}

func TestJoinSortitionPoolWhenEligible(t *testing.T) {
	err := log.SetLogLevel("*", "DEBUG")
	if err != nil {
		t.Fatalf("logger initialization failed: [%v]", err)
	}

	operator, _ := generateAddress()
	stakingProvider, _ := generateAddress()
	blockCounter, _ := local.BlockCounter()
	localChain := connect(operator)

	localChain.registerOperator(stakingProvider, operator)

	// Overwrite the default delay to make things happen faster in the test.
	eligibilityRetryDelay = 1 * time.Second
	// Overwrite the default delay to make things happen faster in the test.
	poolLockRetryDelay = 1 * time.Second

	// Time delay after which the operator gets eligible in the test.
	eligibilityTestDelay := 2 * time.Second
	// Time delay after which the sortition pool gets unlocked in the test.
	lockTestDelay := 1 * time.Second

	// Each eligibility check happens every `eligibilityRetryDelay`.
	// We expect that the operator will join the pool after `eligibilityTestDelay`
	// and `lockTestDelay`.
	expectedEligibleStakeAttempts := 4

	// Lock the sortition pool.
	localChain.isPoolLocked = true

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
		t.Fatalf("operator already joined the pool")
	}

	go func() {
		time.Sleep(eligibilityTestDelay)
		localChain.eligibleStakes[stakingProvider.Hex()] = big.NewInt(100)

		time.Sleep(lockTestDelay)
		localChain.isPoolLocked = false
	}()

	// Wait for the `waitResolve` to resolve to verify if the `waitUntilRegistered`
	// function completes execution.
	resolveWait.Wait()

	if localChain.eligibleStakeAttempts != expectedEligibleStakeAttempts {
		t.Errorf(
			"unexpected EligibleStake attempts\nexpected: [%d]\nactual:   [%d]",
			expectedEligibleStakeAttempts,
			localChain.eligibleStakeAttempts,
		)
	}

	isInPool, _ = localChain.IsOperatorInPool()
	if !isInPool {
		t.Errorf("operator has not joined the pool")
	}

}

// Start `waitUntilJoined` when an operator already joined the sortition pool.
func TestWaitUntilJoined_AlreadyJoined(t *testing.T) {
	err := log.SetLogLevel("*", "DEBUG")
	if err != nil {
		t.Fatalf("logger initialization failed: [%v]", err)
	}

	operator, _ := generateAddress()
	blockCounter, _ := local.BlockCounter()
	localChain := connect(operator)

	joinedPoolChan := make(chan struct{})

	expectedIsOperatorInPoolAttempts := 1

	localChain.sortitionPool[operator.Hex()] = big.NewInt(200)

	resolveWait := &sync.WaitGroup{}
	resolveWait.Add(1)
	go func() {
		waitUntilJoined(context.Background(), blockCounter, localChain, joinedPoolChan)
		resolveWait.Done()
	}()

	<-joinedPoolChan
	if localChain.isOperatorInPoolAttempts != expectedIsOperatorInPoolAttempts {
		t.Errorf(
			"unexpected IsOperatorInPool attempts\nexpected: [%d]\nactual:   [%d]",
			expectedIsOperatorInPoolAttempts,
			localChain.isOperatorInPoolAttempts,
		)
	}

	isInPool, _ := localChain.IsOperatorInPool()
	if !isInPool {
		t.Errorf("operator has not joined the pool")
	}

	// Wait for the `waitResolve` to resolve to verify if the `waitUntilRegistered`
	// function completes execution.
	resolveWait.Wait()
}

// Start `waitUntilJoined` when an operator has not been registered yet.
// Register the operator when the waiting loop is running.
func TestWaitUntilJoined_NotJoined(t *testing.T) {
	err := log.SetLogLevel("*", "DEBUG")
	if err != nil {
		t.Fatalf("logger initialization failed: [%v]", err)
	}

	operator, _ := generateAddress()
	blockCounter, _ := local.BlockCounter()
	localChain := connect(operator)

	joinedPoolChan := make(chan struct{})

	// Time delay after which the operator joins the pool in the test.
	joinPoolTestDelay := 3 * time.Second
	// Each pool check happens every block (500 ms).
	// We expect that the operator will be registered after `joinPoolTestDelay`.
	expectedIsOperatorInPoolAttempts := 7

	resolveWait := &sync.WaitGroup{}
	resolveWait.Add(1)
	go func() {
		waitUntilJoined(context.Background(), blockCounter, localChain, joinedPoolChan)
		resolveWait.Done()
	}()

	go func() {
		time.Sleep(joinPoolTestDelay)
		localChain.sortitionPool[operator.Hex()] = big.NewInt(200)
	}()

	<-joinedPoolChan

	if localChain.isOperatorInPoolAttempts != expectedIsOperatorInPoolAttempts {
		t.Errorf(
			"unexpected OperatorToStakingProvider attempts\nexpected: [%d]\nactual:   [%d]",
			expectedIsOperatorInPoolAttempts,
			localChain.isOperatorInPoolAttempts,
		)
	}

	isInPool, _ := localChain.IsOperatorInPool()
	if !isInPool {
		t.Errorf("operator has not joined the pool")
	}

	// Wait for the `waitResolve` to resolve to verify if the `waitUntilRegistered`
	// function completes execution.
	resolveWait.Wait()
}

func generateAddress() (common.Address, error) {
	_, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		return common.Address{}, err
	}
	return operator.PubkeyToAddress(*publicKey), nil
}
