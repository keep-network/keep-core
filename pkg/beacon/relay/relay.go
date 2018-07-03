package relay

import (
	"encoding/binary"
	"math/big"

	"sync"
	"sync/atomic"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

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
	lastSeenEntryLock sync.RWMutex
	lastSeenEntry     event.Entry
}

// IsNextGroup returns true if the next group expected to generate a threshold
// signature is the same as the group the NodeState belongs to.
func (state *NodeState) IsNextGroup() bool {
	group := atomic.LoadUint32(&state.group)
	groupCount := atomic.LoadUint32(&state.groupCount)

	state.lastSeenEntryLock.RLock()
	defer state.lastSeenEntryLock.RUnlock()
	return binary.BigEndian.Uint32(state.lastSeenEntry.Value[:])%groupCount == group
}

func (ns *NodeState) AddGroup() {
	// increment internal group counter for modding
	atomic.AddUint32(&ns.groupCount, 1)
	// add group pubkey associated with group index
}

// EmptyState returns an empty NodeState with no group, zero group count, and
// a nil last seen entry.
func EmptyState() *NodeState {
	return &NodeState{
		groupCount: 0,
		group:      0,
		GroupID:    0,
		lastSeenEntry: event.Entry{
			RequestID:     &big.Int{},
			Value:         [32]byte{},
			GroupID:       &big.Int{},
			PreviousEntry: &big.Int{},
		},
	}
}
