package state

import (
	"context"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-relay-state")

// State is and interface against which relay states should be implemented.
type State interface {
	// DelayBlocks returns the number of blocks for which the current state
	// initialization is delayed. We delay the initialization to give all other
	// cooperating state machines (e.g. those running on other clients)
	// a chance to enter the new state.
	//
	// State machine moves to the new state immediately after the active blocks
	// time of the previous state has been reached but waits with the
	// initialization of the new state before delay blocks timeout is reached.
	DelayBlocks() uint64

	// ActiveBlocks returns the number of blocks during which the current state
	// is active. Blocks are counted after the initiation process of the
	// current state has completed.
	ActiveBlocks() uint64

	// Initiate performs all the required calculations and sends out all the
	// messages associated with the current state. The context passed to this
	// function is scoped to the lifetime of the current state.
	Initiate(ctx context.Context) error

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

// SilentStateDelayBlocks is a delay in blocks for a state that do not
// exchange any network messages as a part of its execution.
//
// There is no delay before such state enters initialization.
const SilentStateDelayBlocks = 0

// SilentStateActiveBlocks is a number of blocks for which a state that do not
// exchange any network messages as a part of its execution should be active.
//
// Such state is active only for a time needed to perform required computations
// and not any longer.
const SilentStateActiveBlocks = 0
