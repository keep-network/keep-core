package state

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/net"
)

// State is and interface against which relay states should be implemented.
type State interface {
	// ActiveBlocks returns the number of blocks during which the current state
	// is active. Blocks are counted after the initiation process of the
	// current state has completed.
	ActiveBlocks() int

	// Initiate performs all the required calculations and sends out all the
	// messages associated with the current state.
	Initiate() error

	// Receive is called each time a new message arrived. Receive is expected to
	// be called for all broadcast channel messages, including the member's own
	// messages.
	Receive(msg net.Message) error

	// NextState performs a state transition to the next state of the protocol.
	// If the current state is the last one, nextState returns `nil`.
	Next() State

	// MemberIndex returns the index of member associated with the current state.
	MemberIndex() group.MemberIndex
}
