package event

import (
	"sync"
)

// GroupSelectionTrack is used to track a start for group selection after GroupSelectionStarted
// event is received. It is used to ensure that the process execution
// is not duplicated, i.e. when the client receives the same event multiple times.
// When event is received, it should be added in this struct. When a group
// selection process completes (no matter if it succeeded or failed), it should
// be removed from this struct.
type GroupSelectionTrack struct {
	Data  map[string]bool // <entry, bool>
	Mutex *sync.Mutex
}

// Add inserts a new entry into a map that reflects whether a new group selection
// has already started using an entry as a seed.
func (gst *GroupSelectionTrack) Add(entry string) bool {
	gst.Mutex.Lock()
	defer gst.Mutex.Unlock()

	if gst.Data[entry] {
		return false
	}

	gst.Data[entry] = true

	return true
}

// Remove clears a map from an entry that was used to start a group selection.
func (gst *GroupSelectionTrack) Remove(entry string) {
	gst.Mutex.Lock()
	defer gst.Mutex.Unlock()

	delete(gst.Data, entry)
}

// RelayRequestTrack is used to track requests for new entries after RelayEntryRequested
// event is received. It is used to ensure that the process execution
// is not duplicated, i.e. when the client receives the same event multiple times.
// When event is received, it should be added in this struct. When a new entry
// generation process completes (no matter if it succeeded or failed), it should
// be removed from this struct.
type RelayRequestTrack struct {
	Data  map[string]bool // <previousEntry, bool>
	Mutex *sync.Mutex
}

// Add inserts a previous entry into a map that reflects whether a new relay entry
// has already been requested.
func (rrt *RelayRequestTrack) Add(previousEntry string) bool {
	rrt.Mutex.Lock()
	defer rrt.Mutex.Unlock()

	if rrt.Data[previousEntry] {
		return false
	}

	rrt.Data[previousEntry] = true

	return true
}

// Remove clears a map from a previous entry.
func (rrt *RelayRequestTrack) Remove(previousEntry string) {
	rrt.Mutex.Lock()
	defer rrt.Mutex.Unlock()

	delete(rrt.Data, previousEntry)
}
