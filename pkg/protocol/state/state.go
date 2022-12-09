// Package state contains generic state machine implementations.
//
// state.SyncMachine is meant to be used with interactive protocols when
// participants are expected to synchronize based on the number of blocks being
// mined at the time the protocol is executing. Even if the given participant
// received all the necessary information to continue the protocol, the state
// machine waits with proceeding to the next step for the fixed duration of
// blocks. This approach is the most optimal when the protocol may finish
// successfully even if some members expected to participate in the execution
// are inactive.
//
// state.AsyncMachine is meant to be used with interactive protocols when
// participants are expected to synchronize based on the messages being sent and
// some participants may be slower than others. Each protocol participant, to
// finish the execution, must wait for the slowest participant. This approach is
// the most optimal when the protocol may finish successfully only if all
// members expected to participate in the execution are actively participating.
package state

import (
	"context"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// SilentStateDelayBlocks is a delay in blocks for a SyncState that do not
// exchange any network messages as a part of its execution.
//
// There is no delay before such state enters initialization.
const SilentStateDelayBlocks = 0

// SilentStateActiveBlocks is a number of blocks for which a SyncState that do
// not exchange any network messages as a part of its execution should be
// active.
//
// Such state is active only for a time needed to perform required computations
// and not any longer.
const SilentStateActiveBlocks = 0

// SyncState is and interface that should be implemented by protocol states
// of sync.Machine.
type SyncState interface {
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
	Next() (SyncState, error)

	// MemberIndex returns the index of member associated with the current state.
	MemberIndex() group.MemberIndex
}

// AsyncState is and interface that should be implemented by protocol states of
// async.Machine.
//
// The synchronization of the execution is based on an signal from each
// participant that they are ready to proceed to the next step when, for
// example, they received all the necessary information from all the other
// participants.
//
// This approach allows for faster execution of protocols but has strict
// requirements regarding the implementation of states.
//
// Requirement 1: Context lifetime and retransmissions
//
// The context passed to `state.AsyncMachine` must be active as long as the
// result is not published to the chain or until a fixed time for the protocol
// execution has not passed.
//
// The context is used for the retransmission of messages and all protocol
// participants must have a chance to receive messages from other participants.
//
// Consider the following example: there are two participants of the protocol:
// A and B, and they are executing the final state of the protocol. If the
// context was canceled right after the completion of work by the participant,
// without a confirmation that the result was published on-chain, we could run
// into a situation when A received a message from B and exited immediately,
// without giving B a chance to receive a message:
//
// A: |-S------R-|
// B:        |-S----------(...)
//
// |- denotes when the last state starts and -| denotes when it ends. S denotes
// when the given participant sent its message and R denotes when the given
// participant received the message from the other. Since B initiated the last
// state later than A (it could execute some time-consuming computations in the
// last state), A already sent its message and now B needs to wait for the
// retransmissions. B sends its message immediately upon the initiation and this
// message is received by A. Since A exits and cancels the context immediately
// after receiving a message from B, it is no longer retransmitting its message
// and B hangs forever.
//
// There are two solutions. One is that both A and B observe the chain and they
// keep the context active until the result is published. Another is that each
// member waits for some time before completing the protocol and canceling the
// context for retransmissions or that the entire protocol has a fixed maximum
// execution time until the timeout.
//
// Requirement 2: Store all received messages
//
// Since the state machine does not require strict synchronization between
// participants, it is not guaranteed at which moment of the execution the rest
// of the group is. If the current member is at the first state of the
// execution, and all other members advanced to further states, the current
// member should accept and store messages from further states "for the future"
// instead of rejecting them. This can be achieved using `BaseAsyncState`
// structure.
type AsyncState interface {

	// CanTransition indicates if the state has all the information needed and
	// is ready to advance to the next one.
	// It is guaranteed that CanTransition() is not called for the current state
	// until Initiate() completes the work.
	CanTransition() bool

	// Initiate performs all the required calculations and sends out all the
	// messages associated with the current state. The context passed to this
	// function is scoped to the lifetime of the entire state machine and is
	// cancelled when te state machine completed. Use this context for message
	// retransmission to ensure all late protocol participants can catch up.
	Initiate(ctx context.Context) error

	// Receive is called each time a new message arrived. Receive is expected to
	// be called for all broadcast channel messages, including the member's own
	// messages.
	Receive(msg net.Message) error

	// Next performs a state transition to the next state of the protocol.
	// If the current state is the last one, nextState returns `nil`.
	Next() (AsyncState, error)

	// MemberIndex returns the index of member associated with the current state.
	MemberIndex() group.MemberIndex
}

// BaseAsyncState allows to store all received messages even if they are not
// used for the currently executed state. This is important to allow the state
// machine to eventually synchronize with other participants if the current
// state machine is late compared to the rest of the group.
type BaseAsyncState struct {
	messages      map[string][]net.Message
	messagesMutex sync.RWMutex
}

func NewBaseAsyncState() *BaseAsyncState {
	return &BaseAsyncState{
		messages: make(map[string][]net.Message),
	}
}

// ReceiveToHistory stores the received message by its Type() in the internal
// history of received messages.
//
// This function is NOT performing any validation of the received net.Message,
// especially if the public key matches member index or if the given operator is
// allowed to publish messages in the broadcast channel. This function should
// be wrapped with a set of validations.
func (bas *BaseAsyncState) ReceiveToHistory(msg net.Message) {
	bas.messagesMutex.Lock()
	defer bas.messagesMutex.Unlock()

	messageType := msg.Type()
	bas.messages[messageType] = append(bas.messages[messageType], msg)
}

// GetAllReceivedMessages returns all messages of the given Type() received
// so far. If no messages of the given Type() were received, nil slice is
// returned.
func (bas *BaseAsyncState) GetAllReceivedMessages(messageType string) []net.Message {
	bas.messagesMutex.RLock()
	defer bas.messagesMutex.RUnlock()

	return bas.messages[messageType]
}

// ExtractMessagesPayloads gets all messages of type messageType from the
// state's history, extracts their payloads and tries to cast them to type T.
// If the payload cannot be cast to T, it is ignored. It is the caller's
// responsibility to ensure that the type T actually corresponds to the
// type of payloads extracted from messages of type messageType.
func ExtractMessagesPayloads[T any](
	state *BaseAsyncState,
	messageType string,
) []T {
	payloads := make([]T, 0)
	for _, msg := range state.GetAllReceivedMessages(messageType) {
		payload, ok := msg.Payload().(T)
		if !ok {
			continue
		}

		payloads = append(payloads, payload)
	}

	return payloads
}

// DeduplicateMessagesPayloads returns a deduplicated copy of payloads.
// The deduplication logic determines duplicates using the provided keyFn that
// is used to produce a key for the given payload. Payloads using the same key
// are considered equal. Only the first payload with the given key is put to
// the resulting slice.
func DeduplicateMessagesPayloads[T any](payloads []T, keyFn func(T) string) []T {
	keys := make(map[string]bool)
	result := make([]T, 0)

	for _, payload := range payloads {
		key := keyFn(payload)
		if _, exists := keys[key]; !exists {
			keys[key] = true
			result = append(result, payload)
		}
	}

	return result
}
