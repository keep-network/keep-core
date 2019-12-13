package state

import (
	"fmt"
	"strconv"
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
func (m *Machine) Execute(startBlockHeight uint64) (State, uint64, error) {
	// Use an unbuffered channel to serialize message processing.
	recvChan := make(chan net.Message)
	handler := net.HandleMessageFunc{
		Type: fmt.Sprintf("dkg/%s", strconv.FormatInt(time.Now().UTC().UnixNano(), 16)),
		Handler: func(msg net.Message) error {
			recvChan <- msg
			return nil
		},
	}

	m.channel.Recv(handler)
	defer m.channel.UnregisterRecv(handler.Type)

	currentState := m.initialState

	logger.Infof(
		"[member:%v,channel:%s] waiting for block %v to start execution",
		currentState.MemberIndex(),
		m.channel.Name()[:5],
		startBlockHeight,
	)
	err := m.blockCounter.WaitForBlockHeight(startBlockHeight)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to wait for the execution start block")
	}

	lastStateEndBlockHeight := startBlockHeight

	blockWaiter, err := stateTransition(
		currentState,
		lastStateEndBlockHeight,
		m.blockCounter,
		m.channel.Name()[:5],
	)
	if err != nil {
		return nil, 0, err
	}

	for {
		select {
		case msg := <-recvChan:
			err := currentState.Receive(msg)
			if err != nil {
				logger.Errorf(
					"[member:%v,channel:%s, state: %T] failed to receive a message: [%v]",
					currentState.MemberIndex(),
					m.channel.Name()[:5],
					currentState,
					err,
				)
			}

		case lastStateEndBlockHeight := <-blockWaiter:
			nextState := currentState.Next()
			if nextState == nil {
				logger.Infof(
					"[member:%v,channel:%s,state:%T] reached final state at block: [%v]",
					currentState.MemberIndex(),
					m.channel.Name()[:5],
					currentState,
					lastStateEndBlockHeight,
				)
				return currentState, lastStateEndBlockHeight, nil
			}

			currentState = nextState

			blockWaiter, err = stateTransition(
				currentState,
				lastStateEndBlockHeight,
				m.blockCounter,
				m.channel.Name()[:5],
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
	lastStateEndBlockHeight uint64,
	blockCounter chain.BlockCounter,
	channelName string,
) (<-chan uint64, error) {
	logger.Infof(
		"[member:%v,channel:%s,state:%T] transitioning to a new state at block: [%v]",
		currentState.MemberIndex(),
		channelName,
		currentState,
		lastStateEndBlockHeight,
	)

	// We delay the initialization of the new state by one block to give all
	// other coopearating state machines a chance to enter the new state.
	// This is needed when, for example, during the initialization some
	// state-specific messages are sent.
	initiateDelay := lastStateEndBlockHeight + currentState.DelayBlocks()
	err := blockCounter.WaitForBlockHeight(initiateDelay)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to wait [%v] blocks entering state [%T]: [%v]",
			currentState.DelayBlocks(),
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
			"failed to initialize block height waiter at state [%T]: [%v]",
			currentState,
			err,
		)
	}

	logger.Infof(
		"[member:%v,channel:%s,state:%T] transitioned to new state",
		currentState.MemberIndex(),
		channelName,
		currentState,
	)

	return blockWaiter, nil
}
