package signing

import (
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"math/big"
)

// Execute runs the tECDSA signing protocol, given a message to sign,
// broadcast channel to mediate with, a block counter used for time tracking,
// a member index to use in the group, private key share, dishonest threshold,
// and block height when signing protocol should start.
//
// This function also supports signing execution with a subset of the signing
// group by passing a non-empty excludedMembers slice holding the members that
// should be excluded.
func Execute(
	logger log.StandardLogger,
	message *big.Int,
	sessionID string,
	startBlockNumber uint64,
	memberIndex group.MemberIndex,
	privateKeyShare *tecdsa.PrivateKeyShare,
	groupSize int,
	dishonestThreshold int,
	excludedMembers []group.MemberIndex,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
) (*Result, error) {
	logger.Debugf("[member:%v] initializing member", memberIndex)

	registerUnmarshallers(channel)

	member := newMember(
		logger,
		memberIndex,
		groupSize,
		dishonestThreshold,
		membershipValidator,
		sessionID,
	)

	// Mark excluded members as disqualified in order to not exchange messages
	// with them.
	for _, excludedMember := range excludedMembers {
		if excludedMember != member.id {
			member.group.MarkMemberAsDisqualified(excludedMember)
		}
	}

	initialState := &ephemeralKeyPairGenerationState{
		channel: channel,
		member:  member.initializeEphemeralKeysGeneration(),
	}

	stateMachine := state.NewMachine(logger, channel, blockCounter, initialState)

	lastState, _, err := stateMachine.Execute(startBlockNumber)
	if err != nil {
		return nil, err
	}

	_, ok := lastState.(*symmetricKeyGenerationState)
	if !ok {
		return nil, fmt.Errorf("execution ended on state: %T", lastState)
	}

	// TODO: Temporary result. Eventually, take the result from the
	//       finalization state.
	return &Result{
		R:          new(big.Int).Add(message, big.NewInt(1)),
		S:          new(big.Int).Add(message, big.NewInt(2)),
		RecoveryID: 1,
	}, nil
}

// registerUnmarshallers initializes the given broadcast channel to be able to
// perform signing protocol interactions by registering all the required
// protocol message unmarshallers.
func registerUnmarshallers(channel net.BroadcastChannel) {
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &ephemeralPublicKeyMessage{}
	})
}