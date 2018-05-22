package dkg

import (
	"fmt"
	"math/rand"

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
// broadcast channel to mediate it and a group size and threshold. It returns a
// threshold group member who is participating in the group if the generation
// was successful, and an error representing what went wrong if not.
func ExecuteDKG(
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	groupSize int,
	threshold int,
) (*thresholdgroup.Member, error) {
	// Generate a nonzero memberID; loop until rand.Int31 returns something
	// other than 0, hopefully no more than once :)
	memberID := "0"
	for memberID = fmt.Sprintf("%v", rand.Int31()); memberID == "0"; {
	}
	fmt.Printf("[member:0x%010s] Initializing member.\n", memberID)
	localMember, err := thresholdgroup.NewMember(memberID, threshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize member: [%v]", err)
	}

	// Use an unbuffered channel to serialize message processing.
	recvChan := make(chan net.Message)
	channel.Recv(func(msg net.Message) error {
		recvChan <- msg
		return nil
	})

	var (
		currentState, pendingState keyGenerationState
		blockWaiter                <-chan int
	)

	stateTransition := func() error {
		fmt.Printf(
			"[member:%v, state:%T] Transitioning to state [%T]...\n",
			currentState.groupMember().MemberID(),
			currentState,
			pendingState,
		)
		err := pendingState.initiate()
		if err != nil {
			err :=
				fmt.Errorf(
					"failed to initialize state [%T]: [%v]",
					pendingState,
					err,
				)

			return err
		}

		currentState = pendingState
		pendingState = nil

		blockWaiter = blockCounter.BlockWaiter(currentState.activeBlocks())

		fmt.Printf(
			"[member:%v, state:%T] Transitioned to new state.\n",
			currentState.groupMember().MemberID(),
			currentState,
		)

		return nil
	}

	currentState = &initializationState{channel, localMember}
	pendingState = &initializationState{channel, localMember}
	stateTransition()
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
				err := fmt.Errorf(
					"[member:%v, state: %T] failed to receive message [%v]",
					currentState.groupMember().MemberID(),
					currentState,
					err,
				)
				return nil, err
			}

			nextState, err := currentState.nextState()
			if err != nil {
				err := fmt.Errorf(
					"[member:%v, state: %T] failed to move to next state [%v]",
					currentState.groupMember().MemberID(),
					currentState,
					err,
				)
				return nil, err
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

			err :=
				fmt.Errorf(
					"[member:%v, state: %T] failed to complete state inside active period [%v]",
					currentState.groupMember().MemberID(),
					currentState,
					currentState.activeBlocks(),
				)

			return nil, err
		}
	}
}
