package relay

import "time"

// NodeState represents the current state of a relay node.
type NodeState struct {
	group         int
	groupCount    int
	lastSeenEntry Entry
}

// Entry represents one entry in the threshold relay.
type Entry struct {
	value     int
	timestamp time.Time
}

// IsNextGroup returns true if the next group expected to generate a threshold
// signature is the same as the group the NodeState belongs to.
func (state NodeState) IsNextGroup() bool {
	return state.lastSeenEntry.value%state.groupCount == state.group
}

// EmptyState returns an empty NodeState with no group, zero group count, and
// a nil last seen entry.
func EmptyState() NodeState {
	return NodeState{group: 0, groupCount: 0, lastSeenEntry: Entry{0, time.Unix(0, 0)}}
}
