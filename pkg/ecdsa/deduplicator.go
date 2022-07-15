package ecdsa

import "math/big"

// deduplicator decides whether the given event should be handled by the
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
// - DKG started
type deduplicator struct {
}

func newDeduplicator() *deduplicator {
	return &deduplicator{}
}

// notifyDKGStarted notifies the client wants to start the distributed key
// generation upon receiving an event. It returns boolean indicating whether the
// client should proceed with the execution or ignore the event as a duplicate.
func (d *deduplicator) notifyDKGStarted(
	newDKGSeed *big.Int,
) bool {
	// TODO: Implementation.
	return true
}
