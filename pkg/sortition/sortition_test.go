package sortition

import (
	"context"
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

	time.Sleep(registrationTestDelay)
	localChain.registerOperator(stakingProvider, operator)

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

func generateAddress() (common.Address, error) {
	_, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		return common.Address{}, err
	}
	return operator.PubkeyToAddress(*publicKey), nil
}
