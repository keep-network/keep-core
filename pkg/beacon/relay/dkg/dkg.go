package dkg

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
)

// Init initializes a given broadcast channel to be able to perform distributed
// key generation interactions.
func Init(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &JoinMessage{} })
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &MemberCommitmentsMessage{} })
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &MemberShareMessage{} })
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &AccusationsMessage{} })
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &JustificationsMessage{} })
}

// ExecuteDKG runs the full distributed key generation lifecycle, given a
// broadcast channel to mediate it, a player index to use in the group, and a
// group size and threshold. If generation is successful, it returns a threshold
// group member who can participate in the group; if generation fails, it
// returns an error representing what went wrong.
func ExecuteDKG(
	playerIndex int,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	groupSize int,
	threshold int,
) (*thresholdgroup.Member, error) {
	if playerIndex < 0 {
		return nil, fmt.Errorf(
			"playerIndex must be >= 0, got [%v]",
			playerIndex,
		)
	}
	memberID := fmt.Sprintf("%v", playerIndex+1)

	fmt.Printf("[member:0x%010s] Initializing member.\n", memberID)

	var (
		currentState, pendingState keyGenerationState
		blockWaiter                <-chan int
	)

	localMember, err := thresholdgroup.NewMember(memberID, threshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf(
			"in state [%T], failed to initialize block wait: [%v]",
			currentState,
			err,
		)
	}

	// Use an unbuffered channel to serialize message processing.
	recvChan := make(chan net.Message)
	channel.Recv(func(msg net.Message) error {
		recvChan <- msg
		return nil
	})

	stateTransition := func() error {
		fmt.Printf(
			"[member:%v, state:%T] Transitioning to state [%T]...\n",
			currentState.groupMember().MemberID(),
			currentState,
			pendingState,
		)
		err := pendingState.initiate()
		if err != nil {
			return fmt.Errorf(
				"failed to initialize state [%T]: [%v]",
				pendingState,
				err,
			)
		}

		currentState = pendingState
		pendingState = nil

		blockWaiter, err = blockCounter.BlockWaiter(currentState.activeBlocks())
		if err != nil {
			return fmt.Errorf(
				"failed to initialize blockCounter.BlockWaiter state [%T]: [%v]",
				currentState,
				err,
			)
		}

		fmt.Printf(
			"[member:%v, state:%T] Transitioned to new state.\n",
			currentState.groupMember().MemberID(),
			currentState,
		)

		return nil
	}

	currentState = &initializationState{channel, localMember}
	pendingState = &initializationState{channel, localMember}
	if err := stateTransition(); err != nil {
		return nil, err
	}
	pendingState, err = currentState.nextState()
	if err != nil {
		return nil, fmt.Errorf(
			"[member:%v] failed to start distributed key generation [%v]",
			currentState.groupMember().MemberID(),
			err,
		)
	}

	for {
		select {
		case msg := <-recvChan:
			fmt.Printf(
				"[member:%v, state:%T] Processing message.\n",
				currentState.groupMember().MemberID(),
				currentState,
			)

			err := currentState.receive(msg)
			if err != nil {
				return nil, fmt.Errorf(
					"[member:%v, state: %T] failed to receive message [%v]",
					currentState.groupMember().MemberID(),
					currentState,
					err,
				)
			}

			nextState, err := currentState.nextState()
			if err != nil {
				return nil, fmt.Errorf(
					"[member:%v, state: %T] failed to move to next state [%v]",
					currentState.groupMember().MemberID(),
					currentState,
					err,
				)
			}

			if nextState != currentState {
				pendingState = nextState

				fmt.Printf(
					"[member:%v, state:%T] Waiting for active period to elapse...\n",
					currentState.groupMember().MemberID(),
					currentState,
				)
			}

		case <-blockWaiter:
			if pendingState != nil {
				err := stateTransition()
				if err != nil {
					return nil, err
				}

				continue
			} else if finalized, ok := currentState.groupMember().(*thresholdgroup.Member); ok {
				return finalized, nil
			}

			return nil, fmt.Errorf(
				"[member:%v, state: %T] failed to complete state inside active period [%v]",
				currentState.groupMember().MemberID(),
				currentState,
				currentState.activeBlocks(),
			)
		}
	}
}
