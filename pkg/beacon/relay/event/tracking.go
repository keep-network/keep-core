package event

import (
	"sync"
)

// groupSelectionTrack is used to track a start for group selection after
// GroupSelectionStarted event is received. It is used to ensure that the
// process execution is not duplicated, i.e. when the client receives the same
// event multiple times. When event is received, its block number should be
// updated in this struct.
type groupSelectionTrack struct {
	blockNumber uint64
	mutex       sync.Mutex
}

// Update sets the new value of the last group selection start block number.
func (gst *groupSelectionTrack) update(newBlockNumber uint64) {
	gst.mutex.Lock()
	defer gst.mutex.Unlock()

	gst.blockNumber = newBlockNumber
}

// Get returns the last group selection start block number.
func (gst *groupSelectionTrack) get() uint64 {
	gst.mutex.Lock()
	defer gst.mutex.Unlock()

	return gst.blockNumber
}

// relayRequestTrack is used to track requests for new entries after RelayEntryRequested
// event is received. It is used to ensure that the process execution
// is not duplicated, i.e. when the client receives the same event multiple times.
// When event is received, it should be added in this struct. When a new entry
// generation process completes (no matter if it succeeded or failed), it should
// be removed from this struct.
type relayRequestTrack struct {
	data  map[string]bool // <previousEntry, bool>
	mutex sync.Mutex
}

// Add inserts a previous entry into a map that reflects whether a new relay entry
// has already been requested.
func (rrt *relayRequestTrack) add(previousEntry string) bool {
	rrt.mutex.Lock()
	defer rrt.mutex.Unlock()

	if rrt.data[previousEntry] {
		return false
	}

	rrt.data[previousEntry] = true

	return true
}

// Remove clears a map from a previous entry.
func (rrt *relayRequestTrack) remove(previousEntry string) {
	rrt.mutex.Lock()
	defer rrt.mutex.Unlock()

	delete(rrt.data, previousEntry)
}
