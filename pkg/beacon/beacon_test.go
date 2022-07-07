package beacon

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/subscription"
)

const (
	currentRelayRequestConfirmationRetries = 10
	currentRelayRequestConfirmationDelay   = 10 * time.Millisecond
)

func TestConfirmRelayRequestOnFirstAttempt(t *testing.T) {
	expectedRequestStartBlock := 1888
	onConfirmedExecuted := false
	onConfirmed := func() {
		onConfirmedExecuted = true
	}

	chain := newMockRelayChain(func(attempt int) (int, error) {
		return expectedRequestStartBlock, nil
	})

	confirmCurrentRelayRequest(
		uint64(expectedRequestStartBlock),
		chain,
		onConfirmed,
		currentRelayRequestConfirmationRetries,
		currentRelayRequestConfirmationDelay,
	)

	if !onConfirmedExecuted {
		t.Errorf("onConfirmed function not executed")
	}
	if chain.currentRequestStartBlockExecutionCount != 1 {
		t.Errorf(
			"unexpected number of CurrentRequestStartBlock executions\nexpected: [%v]\nactual:   [%v]",
			1,
			chain.currentRequestStartBlockExecutionCount,
		)
	}
}

func TestConfirmRelayRequestOnLastAttempt(t *testing.T) {
	expectedRequestStartBlock := 1888
	onConfirmedExecuted := false
	onConfirmed := func() {
		onConfirmedExecuted = true
	}

	chain := newMockRelayChain(func(attempt int) (int, error) {
		if attempt != currentRelayRequestConfirmationRetries {
			return 0, nil
		}
		return expectedRequestStartBlock, nil
	})

	confirmCurrentRelayRequest(
		uint64(expectedRequestStartBlock),
		chain,
		onConfirmed,
		currentRelayRequestConfirmationRetries,
		currentRelayRequestConfirmationDelay,
	)

	if !onConfirmedExecuted {
		t.Errorf("onConfirmed function not executed")
	}
	if chain.currentRequestStartBlockExecutionCount != currentRelayRequestConfirmationRetries {
		t.Errorf(
			"unexpected number of CurrentRequestStartBlock executions\nexpected: [%v]\nactual:   [%v]",
			currentRelayRequestConfirmationRetries,
			chain.currentRequestStartBlockExecutionCount,
		)
	}
}

func TestConfirmRelayRequestOnLastAttemptBecauseOfErrors(t *testing.T) {
	expectedRequestStartBlock := 1888
	onConfirmedExecuted := false
	onConfirmed := func() {
		onConfirmedExecuted = true
	}

	chain := newMockRelayChain(func(attempt int) (int, error) {
		if attempt != currentRelayRequestConfirmationRetries {
			return 0, fmt.Errorf("VERY BAD ERROR")
		}
		return expectedRequestStartBlock, nil
	})

	confirmCurrentRelayRequest(
		uint64(expectedRequestStartBlock),
		chain,
		onConfirmed,
		currentRelayRequestConfirmationRetries,
		currentRelayRequestConfirmationDelay,
	)

	if !onConfirmedExecuted {
		t.Errorf("onConfirmed function not executed")
	}
	if chain.currentRequestStartBlockExecutionCount != currentRelayRequestConfirmationRetries {
		t.Errorf(
			"unexpected number of CurrentRequestStartBlock executions\nexpected: [%v]\nactual:   [%v]",
			currentRelayRequestConfirmationRetries,
			chain.currentRequestStartBlockExecutionCount,
		)
	}
}

func TestConfirmRelayRequestWithMoreRecentEntryPending(t *testing.T) {
	expectedRequestStartBlock := 1888
	onConfirmedExecuted := false
	onConfirmed := func() {
		onConfirmedExecuted = true
	}

	chain := newMockRelayChain(func(attempt int) (int, error) {
		return expectedRequestStartBlock + 1, nil
	})

	confirmCurrentRelayRequest(
		uint64(expectedRequestStartBlock),
		chain,
		onConfirmed,
		currentRelayRequestConfirmationRetries,
		currentRelayRequestConfirmationDelay,
	)

	if onConfirmedExecuted {
		t.Errorf("onConfirmed function should not be executed")
	}
	if chain.currentRequestStartBlockExecutionCount != 1 {
		t.Errorf(
			"unexpected number of CurrentRequestStartBlock executions\nexpected: [%v]\nactual:   [%v]",
			1,
			chain.currentRequestStartBlockExecutionCount,
		)
	}
}

func TestConfirmRelayRequestWithAllAttemptsFailed(t *testing.T) {
	expectedRequestStartBlock := 1888
	onConfirmedExecuted := false
	onConfirmed := func() {
		onConfirmedExecuted = true
	}

	chain := newMockRelayChain(func(attempt int) (int, error) {
		return 0, nil
	})

	confirmCurrentRelayRequest(
		uint64(expectedRequestStartBlock),
		chain,
		onConfirmed,
		currentRelayRequestConfirmationRetries,
		currentRelayRequestConfirmationDelay,
	)

	if onConfirmedExecuted {
		t.Errorf("onConfirmed function should not be executed")
	}
	if chain.currentRequestStartBlockExecutionCount != currentRelayRequestConfirmationRetries {
		t.Errorf(
			"unexpected number of CurrentRequestStartBlock executions\nexpected: [%v]\nactual:   [%v]",
			currentRelayRequestConfirmationRetries,
			chain.currentRequestStartBlockExecutionCount,
		)
	}
}

func TestConfirmRelayRequestWithAllAttemptsFailedBecauseOfError(t *testing.T) {
	expectedRequestStartBlock := 1888
	onConfirmedExecuted := false
	onConfirmed := func() {
		onConfirmedExecuted = true
	}

	chain := newMockRelayChain(func(attempt int) (int, error) {
		return 0, fmt.Errorf("VERY BAD ERROR")
	})

	confirmCurrentRelayRequest(
		uint64(expectedRequestStartBlock),
		chain,
		onConfirmed,
		currentRelayRequestConfirmationRetries,
		currentRelayRequestConfirmationDelay,
	)

	if onConfirmedExecuted {
		t.Errorf("onConfirmed function should not be executed")
	}
	if chain.currentRequestStartBlockExecutionCount != currentRelayRequestConfirmationRetries {
		t.Errorf(
			"unexpected number of CurrentRequestStartBlock executions\nexpected: [%v]\nactual:   [%v]",
			currentRelayRequestConfirmationRetries,
			chain.currentRequestStartBlockExecutionCount,
		)
	}
}

func newMockRelayChain(
	currentRequestStartBlockFn func(int) (int, error),
) *mockRelayChain {
	return &mockRelayChain{
		currentRequestStartBlockFn:             currentRequestStartBlockFn,
		currentRequestStartBlockExecutionCount: 0,
	}
}

type mockRelayChain struct {
	currentRequestStartBlockExecutionCount int
	currentRequestStartBlockFn             func(int) (int, error)
}

func (mrc *mockRelayChain) SubmitRelayEntry(
	entry []byte,
) *async.EventEntrySubmittedPromise {
	panic("not implemented")
}

func (mrc *mockRelayChain) OnRelayEntrySubmitted(
	func(entry *event.EntrySubmitted),
) subscription.EventSubscription {
	panic("not implemented")
}

func (mrc *mockRelayChain) OnRelayEntryRequested(
	func(request *event.Request),
) subscription.EventSubscription {
	panic("not implemented")
}

func (mrc *mockRelayChain) ReportRelayEntryTimeout() error {
	panic("not implemented")
}

func (mrc *mockRelayChain) IsEntryInProgress() (bool, error) {
	panic("not implemented")
}

func (mrc *mockRelayChain) CurrentRequestStartBlock() (*big.Int, error) {
	mrc.currentRequestStartBlockExecutionCount++
	startBlock, err := mrc.currentRequestStartBlockFn(mrc.currentRequestStartBlockExecutionCount)
	return big.NewInt(int64(startBlock)), err
}

func (mrc *mockRelayChain) CurrentRequestPreviousEntry() ([]byte, error) {
	panic("not implemented")
}

func (mrc *mockRelayChain) CurrentRequestGroupPublicKey() ([]byte, error) {
	panic("not implemented")
}
