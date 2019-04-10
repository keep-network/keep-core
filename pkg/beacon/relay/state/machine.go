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
func (m *Machine) Execute(startBlock uint64) (State, uint64, error) {
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

	fmt.Printf(
		"[member:%v] Waiting for block %v to start execution...\n",
		currentState.MemberIndex(),
		startBlock,
	)
	err := m.blockCounter.WaitForBlockHeight(startBlock)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to wait for the execution start block")
	}

	lastStateEndBlock := startBlock

	blockWaiter, err := stateTransition(
		currentState,
		lastStateEndBlock,
		m.blockCounter,
	)
	if err != nil {
		return nil, 0, err
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

		case endBlock := <-blockWaiter:
			lastStateEndBlock = endBlock
			nextState := currentState.Next()
			if nextState == nil {
				fmt.Printf(
					"[member:%v, state:%T] Final state reached at block [%v]\n",
					currentState.MemberIndex(),
					currentState,
					lastStateEndBlock,
				)
				return currentState, lastStateEndBlock, nil
			}

			currentState = nextState

			blockWaiter, err = stateTransition(
				currentState,
				lastStateEndBlock,
				m.blockCounter,
			)
			if err != nil {
				return nil, 0, err
			}

			continue
		}
	}
}

func stateTransition(
	currentState State,
	lastStateEndBlock uint64,
	blockCounter chain.BlockCounter,
) (<-chan uint64, error) {
	fmt.Printf(
		"[member:%v, state:%T] Transitioning to a new state at block [%v]...\n",
		currentState.MemberIndex(),
		currentState,
		lastStateEndBlock,
	)

	// We delay the initialization of the new state by one block to give all
	// other coopearating state machines a chance to enter the new state.
	// This is needed when, for example, during the initialization some
	// state-specific messages are sent.
	initiateDelay := lastStateEndBlock + 1
	err := blockCounter.WaitForBlockHeight(initiateDelay)
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

	blockWaiter, err := blockCounter.BlockHeightWaiter(
		initiateDelay + currentState.ActiveBlocks(),
	)
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
