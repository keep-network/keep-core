package state

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/net"
)

// For the entire time of initiating the state transition, messages are not
// handled. We use a buffer to unblock producers and let them perform optional
// filtering/validation during that time. The size of that buffer should not be
// lower than the number of messages which can be delivered by the broadcast
// channel during the time the state is blocked on initiation.
// The async version of the state machine does not require a strict
// synchronization between participants, so this number is also the maximum
// number of messages that could be delivered when the current machine is
// blocked on initiation and the rest of participants advance with work.
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

	onStateDone, err := asyncStateTransition(
		am.ctx,
		am.logger,
		currentState,
	)
	if err != nil {
		return nil, err
	}

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

		case <-onStateDone:
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
			onStateDone, err = asyncStateTransition(
				am.ctx,
				am.logger,
				currentState,
			)
			if err != nil {
				return nil, err
			}
		}
	}
}

func asyncStateTransition(
	ctx context.Context,
	logger log.StandardLogger,
	currentState AsyncState,
) (<-chan interface{}, error) {
	logger.Infof(
		"[member:%v,state:%T] transitioning to a new state",
		currentState.MemberIndex(),
		currentState,
	)

	err := currentState.Initiate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate new state [%w]", err)
	}

	onDone := make(chan interface{})
	ticker := time.NewTicker(transitionCheckInterval)

	go func() {
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

	return onDone, nil
}
