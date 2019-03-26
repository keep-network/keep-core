package gjkr

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// Execute runs the GJKR distributed key generation  protocol, given a
// broadcast channel to mediate it, a block counter used for time tracking,
// a player index to use in the group, and a group size and threshold. If
// generation is successful, it returns a threshold group member who can
// participate in the group; if generation fails, it returns an error
// representing what went wrong.
func Execute(
	memberIndex member.Index,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	threshold int,
	seed *big.Int,
) (*Result, error) {
	fmt.Printf("[member:0x%010v] Initializing member\n", memberIndex)

	member, err := NewMember(
		memberIndex,
		make([]member.Index, 0),
		threshold,
		seed,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new member [%v]", err)
	}

	initialState := &initializationState{
		channel: channel,
		member:  member,
	}

	stateMachine := state.NewMachine(channel, blockCounter, initialState)

	lastState, err := stateMachine.Execute(channelInitialization)

	finalizationState, ok := lastState.(*finalizationState)
	if !ok {
		return nil, fmt.Errorf("execution ended on state %T", lastState)
	}

	return finalizationState.result(),
		nil
}

// channelInitialization initializes a given broadcast channel to be able to
// perform distributed key generation interactions.
func channelInitialization(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &JoinMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &EphemeralPublicKeyMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &MemberCommitmentsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &PeerSharesMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &SecretSharesAccusationsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &MemberPublicKeySharePointsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &PointsAccusationsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &DisqualifiedEphemeralKeysMessage{}
	})
}
