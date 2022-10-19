// Package state contains a generic state machine implementation that is
// meant to be used with interactive protocols which require a synchronization
// mechanism between protocol members. The synchronization is based on
// a fixed number of active and delay blocks. Even if the given participant
// received all the necessary information to continue the protocol, the state
// machine waits with proceeding to the next step for the fixed duration of
// blocks.
package state

import (
	"context"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// State is and interface that should be implemented by protocol states.
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
	// function is scoped to the lifetime of the current state but is cancelled
	// upon state's end only after Initiate returns. Because of that,
	// Initiate should not be blocking for a longer period of time and
	// should always return before state's end.
	Initiate(ctx context.Context) error

	// Receive is called each time a new message arrived. Receive is expected to
	// be called for all broadcast channel messages, including the member's own
	// messages.
	Receive(msg net.Message) error

	// Next performs a state transition to the next state of the protocol.
	// If the current state is the last one, nextState returns `nil`.
	Next() (State, error)

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
