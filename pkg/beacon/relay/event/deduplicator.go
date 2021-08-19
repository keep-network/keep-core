package event

import (
	"fmt"
	"math/big"
)

// Local chain interface to avoid import cycles.
type chain interface {
	BlockTimestamp(blockNumber uint64) (uint64, error)
	GetNumberOfCreatedGroups() (*big.Int, error)
	GetGroupRegistrationTime(groupIndex *big.Int) (*big.Int, error)
}

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
	chain                        chain
	groupSelectionDurationBlocks uint64

	groupSelectionTrack *groupSelectionTrack
	relayRequestTrack   *relayRequestTrack
}

// NewDeduplicator constructs a new Deduplicator instance.
func NewDeduplicator(
	chain chain,
	groupSelectionDurationBlocks uint64,
) *Deduplicator {
	return &Deduplicator{
		chain:                        chain,
		groupSelectionDurationBlocks: groupSelectionDurationBlocks,
		groupSelectionTrack:          &groupSelectionTrack{},
		relayRequestTrack:            &relayRequestTrack{data: make(map[string]bool)},
	}
}

// NotifyGroupSelectionStarted notifies the client wants to start group
// selection upon receiving an event. It returns boolean indicating whether the
// client should proceed with the execution or ignore the event as a duplicate.
func (d *Deduplicator) NotifyGroupSelectionStarted(
	newGroupSelectionStartBlock uint64,
) (bool, error) {
	lastGroupSelectionStartBlock := d.groupSelectionTrack.get()

	lastGroupSelectionTimeoutBlock := lastGroupSelectionStartBlock +
		d.groupSelectionDurationBlocks

	if lastGroupSelectionStartBlock == 0 ||
		newGroupSelectionStartBlock > lastGroupSelectionTimeoutBlock {
		return d.groupSelectionTrack.update(newGroupSelectionStartBlock), nil
	}

	// If the notification comes before last group selection timeout elapsed,
	// there is still a chance a new group selection has been triggered shortly
	// after the last one terminated with success. The client must handle
	// this case.

	groupsNumber, err := d.chain.GetNumberOfCreatedGroups()
	if err != nil {
		return false, fmt.Errorf(
			"could not get number of created groups: [%v]",
			err,
		)
	}

	if groupsNumber != nil && groupsNumber.Uint64() > 0 {
		lastGroupIndex := new(big.Int).Sub(groupsNumber, big.NewInt(1))
		lastGroupTime, err := d.chain.GetGroupRegistrationTime(
			lastGroupIndex,
		)
		if err != nil {
			return false, fmt.Errorf(
				"could not get group registration time: [%v]",
				err,
			)
		}

		lastGroupSelectionStartTime, err := d.chain.BlockTimestamp(
			lastGroupSelectionStartBlock,
		)
		if err != nil {
			return false, fmt.Errorf(
				"could not get last group selection start time: [%v]",
				err,
			)
		}

		if lastGroupTime.Uint64() > lastGroupSelectionStartTime {
			return d.groupSelectionTrack.update(newGroupSelectionStartBlock), nil
		}
	}

	return false, nil
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
