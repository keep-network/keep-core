package event

// Deduplicator decides whether the given event should be handled by the
// client or not.
//
// Event subscription may emit the same event two or more times. The same event
// can be emitted right after it's been emitted for the first time. The same
// event can also be emitted a long time after it's been emitted for the first
// time. It is deduplicator's responsibility to decide whether the given
// event is a duplicate and should be ignored or if it is not a duplicate and
// should be handled.
//
// Those events are supported:
// - group selection started
// - relay entry requested
type Deduplicator struct {
	groupSelectionDurationBlocks uint64

	groupSelectionTrack *groupSelectionTrack
}

// NewDeduplicator constructs a new Deduplicator instance.
func NewDeduplicator(groupSelectionDurationBlocks uint64) *Deduplicator {
	groupSelectionTrack := &groupSelectionTrack{}

	return &Deduplicator{
		groupSelectionDurationBlocks: groupSelectionDurationBlocks,
		groupSelectionTrack:          groupSelectionTrack,
	}
}

// NotifyGroupSelectionStarted notifies the client wants to start group
// selection upon receiving an event. It returns boolean indicating whether the
// client should proceed with the execution or ignore the event as a duplicate.
func (d *Deduplicator) NotifyGroupSelectionStarted(blockNumber uint64) bool {
	lastBlockNumber := d.groupSelectionTrack.get()

	if lastBlockNumber == 0 || blockNumber > lastBlockNumber+d.groupSelectionDurationBlocks {
		d.groupSelectionTrack.update(blockNumber)
		return true
	}

	return false
}
