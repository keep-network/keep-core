package event

import (
	"encoding/hex"
	"testing"
)

func TestStartGroupSelection_NoPriorGroupSelections(t *testing.T) {
	chain := &testChain{
		currentRequestPreviousEntryValue: []byte{},
	}

	// In case of first group selection for that node, the last group selection
	// block number held by the deduplicator is zero.
	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	canGenerate := deduplicator.NotifyGroupSelectionStarted(5)

	if !canGenerate {
		t.Fatal("should be allowed to start group selection")
	}
}

func TestStartGroupSelection_MinGroupSelectionDurationPassed(t *testing.T) {
	chain := &testChain{
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	// Simulate the last group selection occured at block 100
	deduplicator.currentGroupSelectionStartBlock = 100

	// Group selection will be possible at block 100 + 200 + 1 = 301
	canGenerate := deduplicator.NotifyGroupSelectionStarted(301)

	if !canGenerate {
		t.Fatal("should be allowed to start group selection")
	}
}

func TestStartGroupSelection_MinGroupSelectionDurationNotPassed(t *testing.T) {
	chain := &testChain{
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	// Simulate the last group selection occured at block 100
	deduplicator.currentGroupSelectionStartBlock = 100

	// Group selection will be possible at block 100 + 200 + 1 = 301
	canGenerate := deduplicator.NotifyGroupSelectionStarted(300)

	if canGenerate {
		t.Fatal("should not be allowed to start group selection")
	}
}

func TestStartRelayEntry_NoPriorRelayEntries(t *testing.T) {
	chain := &testChain{
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	canGenerate, err := deduplicator.NotifyRelayEntryStarted(
		5,
		"entry",
	)
	if err != nil {
		t.Fatal(err)
	}

	if !canGenerate {
		t.Fatal("should be allowed to start relay entry")
	}
}

func TestStartRelayEntry_SmallerStartBlock(t *testing.T) {
	chain := &testChain{
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	deduplicator.currentRequestStartBlock = 100

	canGenerate, err := deduplicator.NotifyRelayEntryStarted(
		5,
		"entry",
	)
	if err != nil {
		t.Fatal(err)
	}

	if canGenerate {
		t.Fatal("should not be allowed to start relay entry")
	}
}

func TestStartRelayEntry_BiggerStartBlock_DifferentPreviousEntry(t *testing.T) {
	chain := &testChain{
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	deduplicator.currentRequestStartBlock = 100
	deduplicator.currentRequestPreviousEntry = "01"

	canGenerate, err := deduplicator.NotifyRelayEntryStarted(
		101,
		"02",
	)
	if err != nil {
		t.Fatal(err)
	}

	if !canGenerate {
		t.Fatal("should be allowed to start relay entry")
	}
}

func TestStartRelayEntry_BiggerStartBlock_SamePreviousEntryConfirmedOnChain(t *testing.T) {
	bytes, err := hex.DecodeString("01")
	if err != nil {
		t.Fatal(err)
	}

	chain := &testChain{
		currentRequestPreviousEntryValue: bytes,
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	deduplicator.currentRequestStartBlock = 100
	deduplicator.currentRequestPreviousEntry = "01"

	canGenerate, err := deduplicator.NotifyRelayEntryStarted(
		101,
		"01",
	)
	if err != nil {
		t.Fatal(err)
	}

	if !canGenerate {
		t.Fatal("should be allowed to start relay entry")
	}
}

func TestStartRelayEntry_BiggerStartBlock_SamePreviousEntryNotConfirmedOnChain(t *testing.T) {
	bytes, err := hex.DecodeString("02")
	if err != nil {
		t.Fatal(err)
	}

	chain := &testChain{
		currentRequestPreviousEntryValue: bytes,
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	deduplicator.currentRequestStartBlock = 100
	deduplicator.currentRequestPreviousEntry = "01"

	canGenerate, err := deduplicator.NotifyRelayEntryStarted(
		101,
		"01",
	)
	if err != nil {
		t.Fatal(err)
	}

	if canGenerate {
		t.Fatal("should not be allowed to start relay entry")
	}
}

type testChain struct {
	currentRequestPreviousEntryValue []byte
}

func (tc *testChain) CurrentRequestPreviousEntry() ([]byte, error) {
	return tc.currentRequestPreviousEntryValue, nil
}
