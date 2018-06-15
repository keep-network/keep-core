package relay

import (
	"encoding/binary"
	"time"

	callback "github.com/keep-network/keep-core/pkg/callback"
)

// ChainInterface represents the interface that the relay expects to interact
// with the anchoring blockchain on.
type ChainInterface interface {
	// SubmitGroupPublicKey submits a 96-byte BLS public key to the blockchain,
	// associated with a string groupID. An error is generally only returned in
	// case of connectivity issues; on-chain errors are reported through event
	// callbacks.
	SubmitGroupPublicKey(groupID string, key [96]byte) *callback.Promise
}

// NodeState represents the current state of a relay node.
type NodeState struct {
	// groupCount is the total number of groups in the relay.
	groupCount uint32
	// group is the id of the relay group this node belongs to. 0 if none.
	// Necessarily less than groupCount.
	group uint32
	// groupId is the id of this node within its relay group. 0 if none.
	GroupID uint16
	// lastSeenEntry is the last relay entry this node is aware of.
	lastSeenEntry Entry
}

// Entry represents one entry in the threshold relay.
type Entry struct {
	Value     [8]byte
	Timestamp time.Time
}

// IsNextGroup returns true if the next group expected to generate a threshold
// signature is the same as the group the NodeState belongs to.
func (state NodeState) IsNextGroup() bool {
	return binary.BigEndian.Uint32(state.lastSeenEntry.Value[:])%state.groupCount == state.group
}

// EmptyState returns an empty NodeState with no group, zero group count, and
// a nil last seen entry.
func EmptyState() NodeState {
	return NodeState{groupCount: 0, group: 0, GroupID: 0, lastSeenEntry: Entry{[8]byte{}, time.Unix(0, 0)}}
}
