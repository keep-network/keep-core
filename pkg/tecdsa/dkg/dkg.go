package dkg

import (
	"fmt"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

// Executor represents an ECDSA distributed key generation process executor.
type Executor struct {
	logger           log.StandardLogger
	tssPreParamsPool *tssPreParamsPool
}

// NewExecutor creates a new Executor instance.
func NewExecutor(
	logger log.StandardLogger,
	scheduler *generator.Scheduler,
	preParamsPoolSize int,
	preParamsGenerationTimeout time.Duration,
	preParamsGenerationDelay time.Duration,
	preParamsGenerationConcurrency int,
) *Executor {
	return &Executor{
		logger: logger,
		tssPreParamsPool: newTssPreParamsPool(
			logger,
			scheduler,
			preParamsPoolSize,
			preParamsGenerationTimeout,
			preParamsGenerationDelay,
			preParamsGenerationConcurrency,
		),
	}
}

// Execute runs the tECDSA distributed key generation protocol, given a
// broadcast channel to mediate with, a block counter used for time tracking,
// a member index to use in the group, dishonest threshold, and block height
// when DKG protocol should start.
//
// This function also supports DKG execution with a subset of the selected
// group by passing a non-empty excludedMembers slice holding the members that
// should be excluded.
func (e *Executor) Execute(
	sessionID string,
	startBlockNumber uint64,
	memberIndex group.MemberIndex,
	groupSize int,
	dishonestThreshold int,
	excludedMembers []group.MemberIndex,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
) (*Result, uint64, error) {
	e.logger.Debugf("[member:%v] initializing member", memberIndex)

	registerUnmarshallers(channel)

	preParams, err := e.tssPreParamsPool.Get()
	if err != nil {
		return nil, 0, fmt.Errorf("failed fetching pre-params: [%v]", err)
	}

	member := newMember(
		e.logger,
		memberIndex,
		groupSize,
		dishonestThreshold,
		membershipValidator,
		sessionID,
		preParams.data,
	)

	// Mark excluded members as disqualified in order to not exchange messages
	// with them and have them recorded as misbehaving in the final result.
	for _, excludedMember := range excludedMembers {
		if excludedMember != member.id {
			member.group.MarkMemberAsDisqualified(excludedMember)
		}
	}

	initialState := &ephemeralKeyPairGenerationState{
		channel: channel,
		member:  member.initializeEphemeralKeysGeneration(),
	}

	stateMachine := state.NewMachine(e.logger, channel, blockCounter, initialState)

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
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &tssRoundOneMessage{}
	})
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &tssRoundTwoMessage{}
	})
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &tssRoundThreeMessage{}
	})
}
