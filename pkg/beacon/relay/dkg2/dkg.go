package dkg2

import (
	"fmt"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// Init initializes a given broadcast channel to be able to perform distributed
// key generation interactions.
func Init(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.JoinMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.EphemeralPublicKeyMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.MemberCommitmentsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.PeerSharesMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.SecretSharesAccusationsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.MemberPublicKeySharePointsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.PointsAccusationsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.DisqualifiedEphemeralKeysMessage{}
	})
}

// ExecuteDKG runs the full distributed key generation lifecycle, given a
// broadcast channel to mediate it, a block counter used for time tracking,
// a player index to use in the group, and a group size and threshold. If
// generation is successful, it returns a threshold group member who can
// participate in the group; if generation fails, it returns an error
// representing what went wrong.
func ExecuteDKG(
	playerIndex int,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	groupSize int,
	threshold int,
	seed *big.Int,
) ([]byte, error) {
	if playerIndex < 0 {
		return nil, fmt.Errorf("playerIndex must be >= 0, got: %v", playerIndex)
	}
	memberID := gjkr.MemberID(playerIndex + 1)

	fmt.Printf("[member:0x%010v] Initializing member\n", memberID)

	// Use an unbuffered channel to serialize message processing.
	recvChan := make(chan net.Message)
	handler := net.HandleMessageFunc{
		Type: fmt.Sprintf("dkg/%s", string(time.Now().UTC().UnixNano())),
		Handler: func(msg net.Message) error {
			recvChan <- msg
			return nil
		},
	}

	channel.Recv(handler)
	defer channel.UnregisterRecv(handler.Type)

	var (
		currentState keyGenerationState
		blockWaiter  <-chan int
	)

	member := gjkr.NewMember(memberID, make([]gjkr.MemberID, 0), threshold, seed)
	currentState = &initializationState{channel, member}

	if err := stateTransition(
		currentState,
		blockCounter,
		blockWaiter,
	); err != nil {
		return nil, err
	}

	for {
		select {
		case msg := <-recvChan:
			fmt.Printf(
				"[member:%v, state:%T] Processing message\n",
				currentState.memberID(),
				currentState,
			)

			err := currentState.receive(msg)
			if err != nil {
				fmt.Printf(
					"[member:%v, state: %T] Failed to receive a message [%v]",
					currentState.memberID(),
					currentState,
					err,
				)
			}

		case <-blockWaiter:
			if finalState, ok := currentState.(*finalizationState); ok {
				return finalState.result().GroupPublicKey.Marshal(), nil
			}

			currentState = currentState.nextState()
			if err := stateTransition(
				currentState,
				blockCounter,
				blockWaiter,
			); err != nil {
				return nil, err
			}

			continue
		}
	}
}

func stateTransition(
	currentState keyGenerationState,
	blockCounter chain.BlockCounter,
	blockWaiter <-chan int,
) error {
	fmt.Printf(
		"[member:%v, state:%T] Transitioning to a new state...\n",
		currentState.memberID(),
		currentState,
	)

	err := blockCounter.WaitForBlocks(1)
	if err != nil {
		return fmt.Errorf(
			"failed to wait 1 block entering state [%T]: [%v]",
			currentState,
			err,
		)
	}

	err = currentState.initiate()
	if err != nil {
		return fmt.Errorf("failed to initiate new state [%v]", err)
	}

	blockWaiter, err = blockCounter.BlockWaiter(currentState.activeBlocks())
	if err != nil {
		return fmt.Errorf(
			"failed to initialize blockCounter.BlockWaiter state [%T]: [%v]",
			currentState,
			err,
		)
	}

	fmt.Printf(
		"[member:%v, state:%T] Transitioned to new state\n",
		currentState.memberID(),
		currentState,
	)

	return nil
}
