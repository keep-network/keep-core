package state

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// For the entire time of state transition (delay + initiate), messages
// are not handled. We use a buffer to unblock producers and let
// them perform optional filtering/validation during that time.
// The size of that buffer should be equal to the biggest possible
// message count which can be delivered by the broadcast channel
// in the same moment.
const receiveBuffer = 128

// Machine is a state machine that executes states implementing the State
// interface.
type Machine struct {
	logger       log.StandardLogger
	channel      net.BroadcastChannel
	blockCounter chain.BlockCounter
	initialState State // first state from which execution starts
}

// NewMachine returns a new protocol state machine. It requires a broadcast
// channel and an initialization function for the channel to be able to
// perform interactions between protocol parties.
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
		ctx,
		m.logger,
		currentState,
		lastStateEndBlockHeight,
		m.blockCounter,
		m.channel.Name()[:5],
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
					"[member:%v,channel:%s, state: %T] failed to receive a message: [%v]",
					currentState.MemberIndex(),
					m.channel.Name()[:5],
					currentState,
					err,
				)
			}

		case lastStateEndBlockHeight := <-blockWaiter:
			cancelCtx()
			nextState := currentState.Next()
			if nextState == nil {
				m.logger.Infof(
					"[member:%v,channel:%s,state:%T] reached final state at block: [%v]",
					currentState.MemberIndex(),
					m.channel.Name()[:5],
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
				m.channel.Name()[:5],
			)
			if err != nil {
				cancelCtx()
				return nil, 0, err
			}

			continue
		}
	}
}

func stateTransition(
	ctx context.Context,
	logger log.StandardLogger,
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
		"[member:%v,channel:%s,state:%T] transitioned to new state",
		currentState.MemberIndex(),
		channelName,
		currentState,
	)

	return blockWaiter, nil
}
