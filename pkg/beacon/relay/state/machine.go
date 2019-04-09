package state

import (
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// Machine is a state machine that executes over states implemented from State
// interface.
type Machine struct {
	channel      net.BroadcastChannel
	blockCounter chain.BlockCounter
	initialState State // first state from which execution starts
}

// NewMachine returns a new state machine. It requires a broadcast channel and
// an initialization function for the channel to be able to perform interactions.
func NewMachine(
	channel net.BroadcastChannel,
	blockCounter chain.BlockCounter,
	initialState State,
) *Machine {
	return &Machine{
		channel:      channel,
		blockCounter: blockCounter,
		initialState: initialState,
	}
}

// Execute state machine starting with initial state up to finalization. It
// requires the broadcast channel to be pre-initialized.
func (m *Machine) Execute() (State, error) {
	// Use an unbuffered channel to serialize message processing.
	recvChan := make(chan net.Message)
	handler := net.HandleMessageFunc{
		Type: fmt.Sprintf("dkg/%s", string(time.Now().UTC().UnixNano())),
		Handler: func(msg net.Message) error {
			recvChan <- msg
			return nil
		},
	}

	m.channel.Recv(handler)
	defer m.channel.UnregisterRecv(handler.Type)

	currentState := m.initialState

	blockWaiter, err := stateTransition(currentState, m.blockCounter)
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
			nextState := currentState.Next()
			if nextState == nil {
				fmt.Printf(
					"[member:%v, state:%T] Final state reached\n",
					currentState.MemberIndex(),
					currentState,
				)
				return currentState, nil
			}

			currentState = nextState

			blockWaiter, err = stateTransition(currentState, m.blockCounter)
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
) (<-chan uint64, error) {
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
