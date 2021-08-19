package event

import (
	"testing"
)

func TestGroupSelectionTrack_UpdateAndGet(t *testing.T) {
	gst := &groupSelectionTrack{}

	gst.update(1000)

	expectedValue := uint64(1000)
	actualValue := gst.get()

	if expectedValue != actualValue {
		t.Errorf(
			"\nexpected: %v\nactual:   %v",
			expectedValue,
			actualValue,
		)
	}
}

func TestRelayRequestTrack_Add(t *testing.T) {
	previousEntry1 := "0x12345"
	previousEntry2 := "0x67891"

	rrt := &relayRequestTrack{
		data: make(map[string]bool),
	}

	if !rrt.add(previousEntry1) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}

	if !rrt.add(previousEntry2) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}
}

func TestRelayRequestTrackAdd_Duplicate(t *testing.T) {
	previousEntry := "0x12345"

	rrt := &relayRequestTrack{
		data: make(map[string]bool),
	}

	if !rrt.add(previousEntry) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}

	if rrt.add(previousEntry) {
		t.Error("RelayEntryRequested event was emitted before; should not be added")
	}
}

func TestRelayRequestTrackRemove(t *testing.T) {
	previousEntry := "0x12345"

	rrt := &relayRequestTrack{
		data: make(map[string]bool),
	}

	if !rrt.add(previousEntry) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}

	rrt.remove(previousEntry)

	if !rrt.add(previousEntry) {
		t.Error("RelayEntryRequested event was removed; should be added successfully")
	}
}

func TestRelayRequestTrack_WhenEmpty(t *testing.T) {
	previousEntry := "0x12345"

	rrt := &relayRequestTrack{
		data: make(map[string]bool),
	}

	rrt.remove(previousEntry)

	if !rrt.add(previousEntry) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}
}
