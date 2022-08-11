package dkg

import (
	"fmt"
	"github.com/ipfs/go-log"
	"math/big"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
	"github.com/keep-network/keep-core/pkg/subscription"
)

// TODO: Move structure and interfaces to appropriate files.

// DKGStartedEvent represents a DKG start event.
type DKGStartedEvent struct {
	Seed        *big.Int
	BlockNumber uint64
}

// DKGResultSubmissionEvent represents a DKG result submission event. It is emitted
// after a submitted DKG result is positively validated on the chain. It contains
// the index of the member who submitted the result and a final public key of
// the group.
type DKGResultSubmissionEvent struct {
	MemberIndex    uint32
	GroupPublicKey []byte
	Misbehaved     []uint8

	BlockNumber uint64
}

// GroupMemberIndex is an index of a wallet registry group member.
// Maximum value accepted by the chain is 255.
type GroupMemberIndex = uint8

type Chain interface {
	// SubmitDKGResult sends DKG result to a chain, along with signatures over
	// result hash from group participants supporting the result.
	// Signatures over DKG result hash are collected in a map keyed by signer's
	// member index.
	SubmitDKGResult(
		participantIndex GroupMemberIndex,
		result *Result,
		signatures map[GroupMemberIndex][]byte,
	) error
	// IsGroupRegistered checks if group with the given public key is registered
	// on-chain.
	IsGroupRegistered(groupPublicKey []byte) (bool, error)
	// OnDKGResultSubmitted registers a callback that is invoked when an on-chain
	// notification of a new, valid submitted result is seen.
	OnDKGResultSubmitted(
		func(event *DKGResultSubmissionEvent),
	) subscription.EventSubscription
	// Signing returns the chain's signer.
	Signing() chain.Signing
	// CalculateDKGResultHash calculates 256-bit hash of DKG result in standard
	// specific for the chain. Operation is performed off-chain.
	CalculateDKGResultHash(result *Result) (ResultHash, error)
}

// TODO: Description
func Publish(
	logger log.StandardLogger,
	memberIndex group.MemberIndex,
	dkgGroup *group.Group,
	membershipValidator *group.MembershipValidator,
	submissionConfig *SubmissionConfig,
	result *Result,
	channel net.BroadcastChannel,
	chain Chain,
	blockCounter chain.BlockCounter,
	startBlockHeight uint64,
) error {
	signingMember := NewSigningMember(
		logger,
		memberIndex,
		dkgGroup,
		membershipValidator,
		submissionConfig,
	)

	initialState := &resultSigningState{
		channel:                 channel,
		chain:                   chain,
		blockCounter:            blockCounter,
		member:                  signingMember,
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
