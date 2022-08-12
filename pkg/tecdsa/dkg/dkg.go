package dkg

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
	"github.com/keep-network/keep-core/pkg/subscription"
)

// ExecutorConfig carries the config for an Executor.
type ExecutorConfig struct {
	TssPreParamsPoolSize              int
	TssPreParamsPoolGenerationTimeout time.Duration
}

// Executor represents an ECDSA distributed key generation process executor.
type Executor struct {
	logger           log.StandardLogger
	tssPreParamsPool *tssPreParamsPool
}

// NewExecutor creates a new Executor instance.
func NewExecutor(
	logger log.StandardLogger,
	config *ExecutorConfig,
) *Executor {
	return &Executor{
		logger: logger,
		tssPreParamsPool: newTssPreParamsPool(
			logger,
			config.TssPreParamsPoolSize,
			config.TssPreParamsPoolGenerationTimeout,
		),
	}
}

// Execute runs the ECDSA distributed key generation protocol, given a
// broadcast channel to mediate with, a block counter used for time tracking,
// a member index to use in the group, dishonest threshold, and block height
// when DKG protocol should start.
func (e *Executor) Execute(
	seed *big.Int,
	startBlockNumber uint64,
	memberIndex group.MemberIndex,
	groupSize int,
	dishonestThreshold int,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
) (*Result, uint64, error) {
	e.logger.Debugf("[member:%v] initializing member", memberIndex)

	registerUnmarshallers(channel)

	member := newMember(
		e.logger,
		memberIndex,
		groupSize,
		dishonestThreshold,
		membershipValidator,
		seed.Text(16), // TODO: Should change on retry.,
		e.tssPreParamsPool.get(),
	)

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

// StartedDKGEvent represents a DKG start event.
type StartedDKGEvent struct {
	Seed        *big.Int
	BlockNumber uint64
}

// ResultSubmissionEvent represents a DKG result submission event. It is emitted
// after a submitted DKG result is positively validated on the chain. It contains
// the index of the member who submitted the result and a final public key of
// the group.
type ResultSubmissionEvent struct {
	MemberIndex         uint32
	GroupPublicKeyBytes []byte
	Misbehaved          []uint8

	BlockNumber uint64
}

type ResultSubmitter interface {
	// SubmitDKGResult sends DKG result to a chain, along with signatures over
	// result hash from group participants supporting the result.
	// Signatures over DKG result hash are collected in a map keyed by signer's
	// member index.
	SubmitDKGResult(
		participantIndex group.MemberIndex,
		result *Result,
		signatures map[group.MemberIndex][]byte,
	) error
	// IsGroupRegistered checks if group with the given public key is registered
	// on-chain.
	IsGroupRegistered(groupPublicKey []byte) (bool, error)
	// OnDKGResultSubmitted registers a callback that is invoked when an on-chain
	// notification of a new, valid submitted result is seen.
	OnDKGResultSubmitted(
		func(event *ResultSubmissionEvent),
	) subscription.EventSubscription
}

type ResultSigner interface {
	// Signing returns the chain's signer.
	Signing() chain.Signing
	// CalculateDKGResultHash calculates 256-bit hash of DKG result in standard
	// specific for the chain. Operation is performed off-chain.
	CalculateDKGResultHash(result *Result) (ResultHash, error)
}

type ResultHandler interface {
	ResultSigner
	ResultSubmitter
}

// Publish signs the DKG result for the given group member, collects signatures
// from other members and verifies them, and submits the DKG result.
func Publish(
	logger log.StandardLogger,
	memberIndex group.MemberIndex,
	membershipValidator *group.MembershipValidator,
	submissionConfig *SubmissionConfig,
	result *Result,
	channel net.BroadcastChannel,
	resultHandler ResultHandler,
	blockCounter chain.BlockCounter,
	startBlockHeight uint64,
) error {
	initialState := &resultSigningState{
		channel:       channel,
		resultHandler: resultHandler,
		blockCounter:  blockCounter,
		member: &signingMember{
			logger:              logger,
			index:               memberIndex,
			group:               result.Group,
			membershipValidator: membershipValidator,
			config:              submissionConfig,
		},
		result:                  result,
		signatureMessages:       make([]*dkgResultHashSignatureMessage, 0),
		signingStartBlockHeight: startBlockHeight,
	}

	stateMachine := state.NewMachine(logger, channel, blockCounter, initialState)

	lastState, _, err := stateMachine.Execute(startBlockHeight)
	if err != nil {
		return err
	}

	_, ok := lastState.(*resultSubmissionState)
	if !ok {
		return fmt.Errorf("execution ended on state %T", lastState)
	}

	return nil
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
