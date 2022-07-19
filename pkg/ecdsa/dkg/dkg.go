package dkg

import (
	"fmt"
	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

var logger = log.Logger("keep-ecdsa-dkg")

// Execute runs the ECDSA distributed key generation protocol, given a
// broadcast channel to mediate with, a block counter used for time tracking,
// a member index to use in the group, dishonest threshold, and block height
// when DKG protocol should start.
func Execute(
	startBlockNumber uint64,
	memberIndex group.MemberIndex,
	groupSize int,
	dishonestThreshold int,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	membershipValidator group.MembershipValidator,
) (*Result, uint64, error) {
	logger.Debugf("[member:%v] initializing member", memberIndex)

	registerUnmarshallers(channel)

	member := newMember(
		memberIndex,
		groupSize,
		dishonestThreshold,
		membershipValidator,
	)

	initialState := &ephemeralKeyPairGenerationState{
		channel: channel,
		member:  member.initializeEphemeralKeysGeneration(),
	}

	stateMachine := state.NewMachine(channel, blockCounter, initialState)

	lastState, endBlockNumber, err := stateMachine.Execute(startBlockNumber)
	if err != nil {
		return nil, 0, err
	}

	finalizationState, ok := lastState.(*finalizationState)
	if !ok {
		return nil, 0, fmt.Errorf("execution ended on state: %T", lastState)
	}

	return finalizationState.result(), endBlockNumber, nil
}

// registerUnmarshallers initializes the given broadcast channel to be able to
// perform DKG protocol interactions by registering all the required protocol
// message unmarshallers.
func registerUnmarshallers(channel net.BroadcastChannel) {
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &ephemeralPublicKeyMessage{}
	})
}
