package state

import (
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/chain"
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
	MemberIndex() member.Index

	// IsFinalState returns true when final state is reached.
	IsFinalState() bool
}

// Execute state machine starting with initial state up to finalization.
// It requires a broadcast channel and an initialization function for the channel
// to be able to perform interactions.
func Execute(
	channel net.BroadcastChannel,
	channelInitialization func(channel net.BroadcastChannel),
	initialState State,
	blockCounter chain.BlockCounter,
) (State, error) {
	// Use an unbuffered channel to serialize message processing.
	recvChan := make(chan net.Message)
	handler := net.HandleMessageFunc{
		Type: fmt.Sprintf("dkg/%s", string(time.Now().UTC().UnixNano())),
		Handler: func(msg net.Message) error {
			recvChan <- msg
			return nil
		},
	}

	channelInitialization(channel)

	channel.Recv(handler)
	defer channel.UnregisterRecv(handler.Type)

	var (
		currentState State
	)

	currentState = initialState

	blockWaiter, err := stateTransition(currentState, blockCounter)
	if err != nil {
		return nil, err
	}

	for {
		select {
		case msg := <-recvChan:
			fmt.Printf(
				"[member:%v, state:%T] Processing message\n",
				currentState.MemberIndex(),
				currentState,
			)

			err := currentState.Receive(msg)
			if err != nil {
				fmt.Printf(
					"[member:%v, state: %T] Failed to receive a message [%v]\n",
					currentState.MemberIndex(),
					currentState,
					err,
				)
			}

		case <-blockWaiter:
			if currentState.IsFinalState() {
				return currentState, nil
			}

			currentState = currentState.Next()
			blockWaiter, err = stateTransition(currentState, blockCounter)
			if err != nil {
				return nil, err
			}

			continue
		}
	}
}

func stateTransition(
	currentState State,
	blockCounter chain.BlockCounter,
) (<-chan int, error) {
	fmt.Printf(
		"[member:%v, state:%T] Transitioning to a new state...\n",
		currentState.MemberIndex(),
		currentState,
	)

	err := blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to wait 1 block entering state [%T]: [%v]",
			currentState,
			err,
		)
	}

	err = currentState.Initiate()
	if err != nil {
		return nil, fmt.Errorf("failed to initiate new state [%v]", err)
	}

	blockWaiter, err := blockCounter.BlockWaiter(currentState.ActiveBlocks())
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize blockCounter.BlockWaiter state [%T]: [%v]",
			currentState,
			err,
		)
	}

	fmt.Printf(
		"[member:%v, state:%T] Transitioned to new state\n",
		currentState.MemberIndex(),
		currentState,
	)

	return blockWaiter, nil
}
