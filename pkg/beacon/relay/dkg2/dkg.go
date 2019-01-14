package dkg2

import (
	"fmt"
	"math/big"
	"time"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
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

// ExecuteDKG runs the full distributed key generation lifecycle.
func ExecuteDKG(
	requestID *big.Int,
	seed *big.Int,
	index int, // starts with 0
	groupSize int,
	threshold int,
	blockCounter chain.BlockCounter,
	relayChain relayChain.Interface,
	channel net.BroadcastChannel,
) error {
	// The staker index should begin with 1
	playerIndex := index + 1
	if playerIndex < 1 {
		return fmt.Errorf("player index must be >= 1")
	}

	gjkrResult, err := executeGJKR(playerIndex, blockCounter, channel, threshold, seed)
	if err != nil {
		return fmt.Errorf("GJKR execution failed [%v]", err)
	}

	err = executePublishing(
		requestID,
		playerIndex,
		relayChain,
		blockCounter,
		convertResult(gjkrResult, groupSize),
	)
	if err != nil {
		return fmt.Errorf("publishing failed [%v]", err)
	}

	return nil
}

// executeGJKR runs the GJKR distributed key generation  protocol, given a
// broadcast channel to mediate it, a block counter used for time tracking,
// a player index to use in the group, and a group size and threshold. If
// generation is successful, it returns a threshold group member who can
// participate in the group; if generation fails, it returns an error
// representing what went wrong.
func executeGJKR(
	playerIndex int,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	threshold int,
	seed *big.Int,
) (*gjkr.Result, error) {
	memberID := gjkr.MemberID(playerIndex)
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

	member, err := gjkr.NewMember(memberID, make([]gjkr.MemberID, 0), threshold, seed)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new member [%v]", err)
	}
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
				return finalState.result(), nil
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

// convertResult transforms GJKR protocol execution result to a chain specific
// DKG result form. It serializes a group public key to bytes and converts
// disqualified and inactive members lists to a boolean list where each entry
// corresponds to a member in the group and true/false value indicates status of
// the member.
func convertResult(gjkrResult *gjkr.Result, groupSize int) *relayChain.DKGResult {
	var serializedGroupPublicKey []byte
	if gjkrResult.GroupPublicKey != nil {
		serializedGroupPublicKey = gjkrResult.GroupPublicKey.Marshal()
	}

	// convertToBoolSlice converts slice containing members IDs to a slice of
	// group size length where true entry indicates the member was found on
	// passed members IDs slice. It assumes member IDs for a group starts iterating
	// from 1. E.g. for a group size of 3 with a passed members ID slice {2} the
	// resulting boolean slice will be {false, true, false}.
	convertToBoolSlice := func(memberIDsSlice []gjkr.MemberID) []bool {
		boolSlice := make([]bool, groupSize)
		for index := range boolSlice {
			for _, memberID := range memberIDsSlice {
				if memberID.Equals(index + 1) {
					boolSlice[index] = true
				}
			}
		}
		return boolSlice
	}

	return &relayChain.DKGResult{
		Success:        gjkrResult.Success,
		GroupPublicKey: serializedGroupPublicKey,
		Inactive:       convertToBoolSlice(gjkrResult.Inactive),
		Disqualified:   convertToBoolSlice(gjkrResult.Disqualified),
	}
}
