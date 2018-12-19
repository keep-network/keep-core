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
) ([]byte, error) {
	if playerIndex < 0 {
		return nil, fmt.Errorf("playerIndex must be >= 0, got: %v", playerIndex)
	}
	memberID := gjkr.MemberID(playerIndex + 1)

	var (
		currentState keyGenerationState
		blockWaiter  <-chan int
	)

	stateTransition := func() error {
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

	fmt.Printf("[member:0x%010v] Initializing member\n", memberID)
	dkg := preconfiguredDKG()
	member := gjkr.NewMember(memberID, make([]gjkr.MemberID, 0), threshold, dkg)

	currentState = &initializationState{channel, member}
	if err := stateTransition(); err != nil {
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
			if finalState, ok := currentState.(*finalState); ok {
				return finalState.groupPublicKey().Bytes(), nil
			}

			currentState = currentState.nextState()
			if err := stateTransition(); err != nil {
				return nil, err
			}

			continue
		}
	}
}

// We use preconfigured cryptographic parameters which entails the protocol is
// not secure! Once we switch the protocol to use elliptic curves, we'll remote
// this method and the threat will be gone.
func preconfiguredDKG() *gjkr.DKG {
	p := new(big.Int)
	p.SetString("95334665770371710735165175185539898748507557931415093801643066444636813028167", 10)

	q := new(big.Int)
	q.SetString("47667332885185855367582587592769949374253778965707546900821533222318406514083", 10)

	g := new(big.Int)
	g.SetString("26356443122116168287367397770636591434494873554012367991347050804959071206175", 10)

	h := new(big.Int)
	h.SetString("44701896773967475854551946130428071586025426046371976839556825967620150305387", 10)

	return gjkr.NewDKG(p, q, g, h)
}
