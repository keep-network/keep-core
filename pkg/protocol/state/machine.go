package state

import (
	"context"
	"fmt"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// For the entire time of state transition (delay + initiate), messages
// are not handled. We use a buffer to unblock producers and let
// them perform optional filtering/validation during that time.
// The size of that buffer should not be lower than the number of messages
// which can be delivered by the broadcast channel during the time the state
// is blocked on initiation.
// This version of the state machine requires a strict synchronization between
// participants, so this number is also the maximum number of messages that
// could be delivered in a single state.
const receiveBuffer = 128

// Machine is a state machine that executes states implementing the State
// interface.
type Machine struct {
	logger       log.StandardLogger
	channel      net.BroadcastChannel
	blockCounter chain.BlockCounter
	initialState State // first state from which execution starts
}

// NewMachine returns a new protocol state machine.
func NewMachine(
	logger log.StandardLogger,
	channel net.BroadcastChannel,
	blockCounter chain.BlockCounter,
	initialState State,
) *Machine {
	return &Machine{
		logger:       logger,
		channel:      channel,
		blockCounter: blockCounter,
		initialState: initialState,
	}
}

// Execute state machine starting with initial state up to finalization. It
// requires the broadcast channel to be pre-initialized.
func (m *Machine) Execute(startBlockHeight uint64) (State, uint64, error) {
	recvChan := make(chan net.Message, receiveBuffer)
	handler := func(msg net.Message) {
		recvChan <- msg
	}

	currentState := m.initialState
	ctx, cancelCtx := context.WithCancel(context.Background())
	m.channel.Recv(ctx, handler)

	m.logger.Infof(
		"[member:%v] waiting for block [%v] to start execution",
		currentState.MemberIndex(),
		startBlockHeight,
	)
	err := m.blockCounter.WaitForBlockHeight(startBlockHeight)
	if err != nil {
		cancelCtx()
		return nil, 0, fmt.Errorf("failed to wait for the execution start block")
	}

	lastStateEndBlockHeight := startBlockHeight

	blockWaiter, err := stateTransition(
		ctx,
		m.logger,
		currentState,
		lastStateEndBlockHeight,
		m.blockCounter,
	)
	if err != nil {
		cancelCtx()
		return nil, 0, err
	}

	for {
		select {
		case msg := <-recvChan:
			err := currentState.Receive(msg)
			if err != nil {
				m.logger.Errorf(
					"[member:%v,state:%T] failed to receive a message: [%v]",
					currentState.MemberIndex(),
					currentState,
					err,
				)
			}

		case lastStateEndBlockHeight := <-blockWaiter:
			cancelCtx()

			nextState, err := currentState.Next()
			if err != nil {
				return nil, 0, fmt.Errorf(
					"failed to complete state [%T]: [%w]",
					currentState,
					err,
				)
			}

			if nextState == nil {
				m.logger.Infof(
					"[member:%v,state:%T] reached final state at block: [%v]",
					currentState.MemberIndex(),
					currentState,
					lastStateEndBlockHeight,
				)
				return currentState, lastStateEndBlockHeight, nil
			}

			currentState = nextState
			ctx, cancelCtx = context.WithCancel(context.Background())
			m.channel.Recv(ctx, handler)

			blockWaiter, err = stateTransition(
				ctx,
				m.logger,
				currentState,
				lastStateEndBlockHeight,
				m.blockCounter,
			)
			if err != nil {
				cancelCtx()
				return nil, 0, err
			}
		}
	}
}

func stateTransition(
	ctx context.Context,
	logger log.StandardLogger,
	currentState State,
	lastStateEndBlockHeight uint64,
	blockCounter chain.BlockCounter,
) (<-chan uint64, error) {
	logger.Infof(
		"[member:%v,state:%T] transitioning to a new state at block: [%v]",
		currentState.MemberIndex(),
		currentState,
		lastStateEndBlockHeight,
	)

	// We delay the initialization of the new state by `initiateDelay` of blocks
	// to give all other participants a chance to enter the new state. This is
	// needed when state accepts only messages specific to that state.
	// In that case, if the message is sent too early, it is lost given that the
	// receiveBuffer has the retransmissions filtered out.
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

	err = currentState.Initiate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate new state [%w]", err)
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
		"[member:%v,state:%T] transitioned to new state",
		currentState.MemberIndex(),
		currentState,
	)

	return blockWaiter, nil
}
