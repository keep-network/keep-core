package event

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func TestStartGroupSelection_NoPriorGroupSelections(t *testing.T) {
	chain := &testChain{
		currentRequestStartBlockValue:    nil,
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
		currentRequestStartBlockValue:    nil,
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	// Simulate the last group selection occured at block 100
	deduplicator.NotifyGroupSelectionStarted(100)

	// Group selection will be possible at block 100 + 200 + 1 = 301
	canGenerate := deduplicator.NotifyGroupSelectionStarted(301)

	if !canGenerate {
		t.Fatal("should be allowed to start group selection")
	}
}

func TestStartGroupSelection_MinGroupSelectionDurationNotPassed(t *testing.T) {
	chain := &testChain{
		currentRequestStartBlockValue:    nil,
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	// Simulate the last group selection occured at block 100
	deduplicator.NotifyGroupSelectionStarted(100)

	// Group selection will be possible at block 100 + 200 + 1 = 301
	canGenerate := deduplicator.NotifyGroupSelectionStarted(300)

	if canGenerate {
		t.Fatal("should not be allowed to start group selection")
	}
}

func TestStartRelayEntry_NoPriorRelayEntries(t *testing.T) {
	chain := &testChain{
		currentRequestStartBlockValue:    nil,
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

func TestStartRelayEntry_LowerStartBlock(t *testing.T) {
	chain := &testChain{
		currentRequestStartBlockValue:    nil,
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	_, err := deduplicator.NotifyRelayEntryStarted(100, "01")
	if err != nil {
		t.Fatal(err)
	}

	canGenerate, err := deduplicator.NotifyRelayEntryStarted(
		5,
		"02",
	)
	if err != nil {
		t.Fatal(err)
	}

	if canGenerate {
		t.Fatal("should not be allowed to start relay entry")
	}
}

func TestStartRelayEntry_HigherStartBlock_DifferentPreviousEntry(t *testing.T) {
	chain := &testChain{
		currentRequestStartBlockValue:    nil,
		currentRequestPreviousEntryValue: []byte{},
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	_, err := deduplicator.NotifyRelayEntryStarted(100, "01")
	if err != nil {
		t.Fatal(err)
	}

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

func TestStartRelayEntry_HigherStartBlock_SamePreviousEntry_ConfirmedOnChain(t *testing.T) {
	bytes, err := hex.DecodeString("01")
	if err != nil {
		t.Fatal(err)
	}

	chain := &testChain{
		currentRequestStartBlockValue:    big.NewInt(101),
		currentRequestPreviousEntryValue: bytes,
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	_, err = deduplicator.NotifyRelayEntryStarted(100, "01")
	if err != nil {
		t.Fatal(err)
	}

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

func TestStartRelayEntry_HigherStartBlock_SamePreviousEntry_PreviousEntryNotConfirmedOnChain(t *testing.T) {
	bytes, err := hex.DecodeString("02")
	if err != nil {
		t.Fatal(err)
	}

	chain := &testChain{
		currentRequestStartBlockValue:    big.NewInt(101),
		currentRequestPreviousEntryValue: bytes,
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	_, err = deduplicator.NotifyRelayEntryStarted(100, "01")
	if err != nil {
		t.Fatal(err)
	}

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

func TestStartRelayEntry_HigherStartBlock_SamePreviousEntry_StartBlockNotConfirmedOnChain(t *testing.T) {
	bytes, err := hex.DecodeString("01")
	if err != nil {
		t.Fatal(err)
	}

	chain := &testChain{
		currentRequestStartBlockValue:    big.NewInt(100),
		currentRequestPreviousEntryValue: bytes,
	}

	deduplicator := NewDeduplicator(
		chain,
		200,
	)

	_, err = deduplicator.NotifyRelayEntryStarted(100, "01")
	if err != nil {
		t.Fatal(err)
	}

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
	currentRequestStartBlockValue    *big.Int
	currentRequestPreviousEntryValue []byte
}

func (tc *testChain) CurrentRequestStartBlock() (*big.Int, error) {
	return tc.currentRequestStartBlockValue, nil
}

func (tc *testChain) CurrentRequestPreviousEntry() ([]byte, error) {
	return tc.currentRequestPreviousEntryValue, nil
}
