package event

import (
	"sync"
	"testing"
)

func TestGroupSelectionTrack_Add(t *testing.T) {
	entry1 := "0x12345"
	entry2 := "0x67891"

	gst := &GroupSelectionTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	if !gst.Add(entry1) {
		t.Error("GroupSelectionStarted event wasn't emitted before; should be added successfully")
	}

	if !gst.Add(entry2) {
		t.Error("GroupSelectionStarted event wasn't emitted before; should be added successfully")
	}
}

func TestGroupSelectionTrackAdd_Duplicate(t *testing.T) {
	entry := "0x12345"

	gst := &GroupSelectionTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	if !gst.Add(entry) {
		t.Error("GroupSelectionStarted event wasn't emitted before; should be added successfully")
	}

	if gst.Add(entry) {
		t.Error("GroupSelectionStarted event was emitted before; should not be added")
	}
}

func TestGroupSelectionTrackRemove(t *testing.T) {
	entry := "0x12345"

	gst := &GroupSelectionTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	if !gst.Add(entry) {
		t.Error("GroupSelectionStarted event wasn't emitted before; should be added successfully")
	}

	gst.Remove(entry)

	if !gst.Add(entry) {
		t.Error("GroupSelectionStarted event was removed; should be added successfully")
	}
}

func TestGroupSelectionTrack_WhenEmpty(t *testing.T) {
	entry := "0x12345"

	gst := &GroupSelectionTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	gst.Remove(entry)

	if !gst.Add(entry) {
		t.Error("GroupSelectionStarted event wasn't emitted before; should be added successfully")
	}
}

func TestRelayRequestTrack_Add(t *testing.T) {
	previousEntry1 := "0x12345"
	previousEntry2 := "0x67891"

	rrt := &RelayRequestTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	if !rrt.Add(previousEntry1) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}

	if !rrt.Add(previousEntry2) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}
}

func TestRelayRequestTrackAdd_Duplicate(t *testing.T) {
	previousEntry := "0x12345"

	rrt := &RelayRequestTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	if !rrt.Add(previousEntry) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}

	if rrt.Add(previousEntry) {
		t.Error("RelayEntryRequested event was emitted before; should not be added")
	}
}

func TestRelayRequestTrackRemove(t *testing.T) {
	previousEntry := "0x12345"

	rrt := &RelayRequestTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	if !rrt.Add(previousEntry) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}

	rrt.Remove(previousEntry)

	if !rrt.Add(previousEntry) {
		t.Error("RelayEntryRequested event was removed; should be added successfully")
	}
}

func TestRelayRequestTrack_WhenEmpty(t *testing.T) {
	previousEntry := "0x12345"

	rrt := &RelayRequestTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	rrt.Remove(previousEntry)

	if !rrt.Add(previousEntry) {
		t.Error("RelayEntryRequested event wasn't emitted before; should be added successfully")
	}
}
