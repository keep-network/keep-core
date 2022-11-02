package state

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/net"
)

// asyncReceiveBuffer is a buffer for messages received from the broadcast
// channel used when the state machine's current state is temporarily too slow
// to handle them. The asynchronous state machine tries to read from this buffer
// all the time.
const asyncReceiveBuffer = 512

// The time interval with which the CanTransition of the AsyncState condition
// is checked.
const transitionCheckInterval = 100 * time.Millisecond

// AsyncMachine is a state machine that executes states implementing the
// AsyncState interface.
//
// AsyncMachine execution state is based on an signal from each participant that
// they are ready to proceed to the next step when, for example, they received
// all the necessary information from all the other participants.
//
// AsyncMachine is meant to be used with interactive protocols when participants
// are expected to synchronize based on the messages being sent and some
// participants may be slower than others. Each protocol participant, to finish
// the execution, must wait for the slowest participant. This approach is
// the most optimal when the protocol may finish successfully only if all
// members expected to participate in the execution are actively participating.
//
// This approach allows for faster execution of protocols but has
// strict requirements regarding the implementation of states. Please see
// AsyncState documentation for more information.
type AsyncMachine struct {
	logger       log.StandardLogger
	ctx          context.Context
	channel      net.BroadcastChannel
	initialState AsyncState // first state from which execution starts
}

// NewAsyncMachine returns a new protocol asynchronous state machine
//
// The context passed to `faststate.NewSyncMachine` must be active as long as the
// result is not published to the chain or until a fixed time for the protocol
// execution has not passed.
//
// The context is used for the retransmission of messages and all protocol
// participants must have a chance to receive messages from other participants.
func NewAsyncMachine(
	logger log.StandardLogger,
	ctx context.Context,
	channel net.BroadcastChannel,
	initialState AsyncState,
) *AsyncMachine {
	return &AsyncMachine{
		logger:       logger,
		ctx:          ctx,
		channel:      channel,
		initialState: initialState,
	}
}

// Execute state machine starting with initial state up to finalization. It
// requires the broadcast channel to be pre-initialized.
func (am *AsyncMachine) Execute() (AsyncState, error) {
	recvChan := make(chan net.Message, asyncReceiveBuffer)
	handler := func(msg net.Message) {
		recvChan <- msg
	}
	am.channel.Recv(am.ctx, handler)

	currentState := am.initialState

	onStateDone := asyncStateTransition(
		am.ctx,
		am.logger,
		currentState,
	)

	for {
		select {
		case msg := <-recvChan:
			err := currentState.Receive(msg)
			if err != nil {
				am.logger.Errorf(
					"[member:%v,state:%T] failed to receive a message: [%v]",
					currentState.MemberIndex(),
					currentState,
					err,
				)
			}

		case err := <-onStateDone:
			if err != nil {
				return nil, fmt.Errorf(
					"failed to initiate state [%T]: [%w]",
					currentState,
					err,
				)
			}

			nextState, err := currentState.Next()
			if err != nil {
				return nil, fmt.Errorf(
					"failed to complete state [%T]: [%w]",
					currentState,
					err,
				)
			}

			if nextState == nil {
				am.logger.Infof(
					"[member:%v,state:%T] reached final state",
					currentState.MemberIndex(),
					currentState,
				)
				return currentState, nil
			}

			currentState = nextState
			onStateDone = asyncStateTransition(
				am.ctx,
				am.logger,
				currentState,
			)

		case <-am.ctx.Done():
			return nil, am.ctx.Err()
		}
	}
}

// asyncStateTransition kicks of the state initiation calling Initiate()
// function and returns the channel that gets closed when the initiation is
// done. In case the initiation failed, the channel receives an error from
// the Initiate() call. This was, the message Recv loop in the state machine
// is not blocked for the time of initiation.
func asyncStateTransition(
	ctx context.Context,
	logger log.StandardLogger,
	currentState AsyncState,
) <-chan error {
	logger.Infof(
		"[member:%v,state:%T] transitioning to a new state",
		currentState.MemberIndex(),
		currentState,
	)

	onDone := make(chan error)

	go func() {
		err := currentState.Initiate(ctx)
		if err != nil {
			onDone <- err
			return
		}

		ticker := time.NewTicker(transitionCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if currentState.CanTransition() {
					close(onDone)
					return
				}
			}
		}
	}()

	logger.Infof(
		"[member:%v,state:%T] transitioned to new state",
		currentState.MemberIndex(),
		currentState,
	)

	return onDone
}
