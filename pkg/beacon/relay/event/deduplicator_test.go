package event

import (
	"math/big"
	"testing"
)

func TestStartGroupSelection_NoPriorGroupSelections(t *testing.T) {
	chain := &testChain{
		numberOfCreatedGroupsValue:    big.NewInt(0),
		getGroupRegistrationTimeValue: map[uint64]*big.Int{},
	}

	// In case of first group selection for that node, the last group selection
	// block number held by the deduplicator is zero.
	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	canGenerate, err := deduplicator.NotifyGroupSelectionStarted(5)
	if err != nil {
		t.Fatal(err)
	}

	if !canGenerate {
		t.Fatal("should be allowed to start group selection")
	}
}

func TestStartGroupSelection_LastGroupSelectionTimedOut(t *testing.T) {
	chain := &testChain{
		numberOfCreatedGroupsValue:    big.NewInt(0),
		getGroupRegistrationTimeValue: map[uint64]*big.Int{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	// Simulate the last group selection occured at block 100
	deduplicator.groupSelectionTrack.update(100)

	// Group selection will be possible at block 100 + 200 + 1 = 301
	canGenerate, err := deduplicator.NotifyGroupSelectionStarted(301)
	if err != nil {
		t.Fatal(err)
	}

	if !canGenerate {
		t.Fatal("should be allowed to start group selection")
	}
}

func TestStartGroupSelection_LastGroupSelectionNotTimedOut(t *testing.T) {
	chain := &testChain{
		numberOfCreatedGroupsValue:    big.NewInt(0),
		getGroupRegistrationTimeValue: map[uint64]*big.Int{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	// Simulate the last group selection occured at block 100
	deduplicator.groupSelectionTrack.update(100)

	// Group selection will be possible at block 100 + 200 + 1 = 301
	canGenerate, err := deduplicator.NotifyGroupSelectionStarted(300)
	if err != nil {
		t.Fatal(err)
	}

	if canGenerate {
		t.Fatal("should not be allowed to start group selection")
	}
}

func TestStartGroupSelection_LastGroupSelectionNotTimedOutButSuccessful(
	t *testing.T,
) {
	chain := &testChain{
		numberOfCreatedGroupsValue: big.NewInt(10),
		getGroupRegistrationTimeValue: map[uint64]*big.Int{
			9: big.NewInt(200),
		},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	// Simulate the last group selection occured at block 100
	deduplicator.groupSelectionTrack.update(100)

	// Group selection will be possible at block 100 + 200 + 1 = 301 but
	// the last group was registered at T=200 and the new group selection
	// start block has T=300
	canGenerate, err := deduplicator.NotifyGroupSelectionStarted(300)
	if err != nil {
		t.Fatal(err)
	}

	if !canGenerate {
		t.Fatal("should be allowed to start group selection")
	}
}

func TestStartRelayEntry(t *testing.T) {
	chain := &testChain{
		numberOfCreatedGroupsValue:    big.NewInt(0),
		getGroupRegistrationTimeValue: map[uint64]*big.Int{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	canGenerate := deduplicator.NotifyRelayEntryStarted("entry")

	if !canGenerate {
		t.Fatal("should be allowed to start relay entry")
	}
}

func TestStartRelayEntry_AlreadyStarted(t *testing.T) {
	chain := &testChain{
		numberOfCreatedGroupsValue:    big.NewInt(0),
		getGroupRegistrationTimeValue: map[uint64]*big.Int{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	deduplicator.NotifyRelayEntryStarted("entry")
	canGenerate := deduplicator.NotifyRelayEntryStarted("entry")

	if canGenerate {
		t.Fatal("should be not allowed to start relay entry")
	}
}

func TestStartRelayEntry_AlreadyStartedButCompleted(t *testing.T) {
	chain := &testChain{
		numberOfCreatedGroupsValue:    big.NewInt(0),
		getGroupRegistrationTimeValue: map[uint64]*big.Int{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	deduplicator.NotifyRelayEntryStarted("entry")
	deduplicator.NotifyRelayEntryCompleted("entry")
	canGenerate := deduplicator.NotifyRelayEntryStarted("entry")

	if !canGenerate {
		t.Fatal("should be allowed to start relay entry")
	}
}

type testChain struct {
	numberOfCreatedGroupsValue    *big.Int
	getGroupRegistrationTimeValue map[uint64]*big.Int // <groupIndex, time>
}

func (tc *testChain) BlockTimestamp(blockNumber uint64) (uint64, error) {
	return blockNumber, nil
}

func (tc *testChain) GetNumberOfCreatedGroups() (*big.Int, error) {
	return tc.numberOfCreatedGroupsValue, nil
}

func (tc *testChain) GetGroupRegistrationTime(
	groupIndex *big.Int,
) (*big.Int, error) {
	return tc.getGroupRegistrationTimeValue[groupIndex.Uint64()], nil
}
