package gjkr

import (
	"fmt"
	"math/big"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-gjkr")

// RegisterUnmarshallers initializes the given broadcast channel to be able to
// perform DKG protocol interactions by registering all the required protocol
// message unmarshallers.
// The channel needs to be fully initialized before Execute is called.
func RegisterUnmarshallers(channel net.BroadcastChannel) {
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
		return &MisbehavedEphemeralKeysMessage{}
	})
}

// Execute runs the GJKR distributed key generation  protocol, given a
// broadcast channel to mediate with, a block counter used for time tracking,
// a player index to use in the group, dishonest threshold, and block height
// when DKG protocol should start.
// If the generation is successful, it returns a threshold group member which
// can participate in the signing group; if the generation fails, it returns an
// error.
func Execute(
	memberIndex group.MemberIndex,
	groupSize int,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	dishonestThreshold int,
	seed *big.Int,
	startBlockHeight uint64,
) (*Result, uint64, error) {
	logger.Debugf("[member:%v] initializing member", memberIndex)

	member, err := NewMember(
		memberIndex,
		groupSize,
		dishonestThreshold,
		seed,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot create a new member: [%v]", err)
	}

	initialState := &ephemeralKeyPairGenerationState{
		channel: channel,
		member:  member.InitializeEphemeralKeysGeneration(),
	}

	stateMachine := state.NewMachine(channel, blockCounter, initialState)

	lastState, endBlockHeight, err := stateMachine.Execute(startBlockHeight)
	if err != nil {
		return nil, 0, err
	}

	finalizationState, ok := lastState.(*finalizationState)
	if !ok {
		return nil, 0, fmt.Errorf("execution ended on state: %T", lastState)
	}

	return finalizationState.result(), endBlockHeight, nil
}
