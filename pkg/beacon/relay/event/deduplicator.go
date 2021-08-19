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
	relayRequestTrack   *relayRequestTrack
}

// NewDeduplicator constructs a new Deduplicator instance.
func NewDeduplicator(
	groupSelectionDurationBlocks uint64,
) *Deduplicator {
	return &Deduplicator{
		groupSelectionDurationBlocks: groupSelectionDurationBlocks,
		groupSelectionTrack:          &groupSelectionTrack{},
		relayRequestTrack:            &relayRequestTrack{data: make(map[string]bool)},
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

// NotifyRelayEntryStarted notifies the client wants to start relay entry
// generation upon receiving an event. It returns boolean indicating whether the
// client should proceed with the execution or ignore the event as a duplicate.
//
// In case the client proceeds with the relay entry generation, it should call
// NotifyRelayEntryCompleted once the protocol completes, no matter if it
// failed or succeeded.
func (d *Deduplicator) NotifyRelayEntryStarted(previousEntry string) bool {
	return d.relayRequestTrack.add(previousEntry)
}

// NotifyRelayEntryCompleted should be called once client completed relay
// entry generation protocol, no matter if it succeeded or not.
func (d *Deduplicator) NotifyRelayEntryCompleted(previousEntry string) {
	d.relayRequestTrack.remove(previousEntry)
}
