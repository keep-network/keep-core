package result

import (
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

// represents a given state in the state machine for signing dkg results
type signingState = state.State

// resultSigningState is the state during which group members sign their preferred
// dkg result (by hashing their dkg result, and then signing the result), and
// share this over the broadcast channel.
//
// State is part of phase 13 of the protocol.
type resultSigningState struct {
	channel net.BroadcastChannel
	handle  chain.Handle

	member *SigningMember

	requestID *big.Int
	result    *relayChain.DKGResult

	signatureMessages []*DKGResultHashSignatureMessage
}

func (rss *resultSigningState) ActiveBlocks() int { return 3 }

func (rss *resultSigningState) Initiate() error {
	message, err := rss.member.SignDKGResult(rss.result, rss.handle)
	if err != nil {
		return err
	}
	if err := rss.channel.Send(message); err != nil {
		return err
	}
	return nil
}

func (rss *resultSigningState) Receive(msg net.Message) error {
	switch signedMessage := msg.Payload().(type) {
	case *DKGResultHashSignatureMessage:
		// ignore messages from ourselves
		if !group.IsMessageFromSelf(rss.member.index, signedMessage) &&
			group.IsSenderAccepted(rss.member, signedMessage) {
			rss.signatureMessages = append(rss.signatureMessages, signedMessage)
		}
	}

	return nil
}

func (rss *resultSigningState) Next() signingState {
	// set up the verification state, phase 13 part 2
	return &signaturesVerificationState{
		channel:           rss.channel,
		handle:            rss.handle,
		member:            rss.member,
		requestID:         rss.requestID,
		result:            rss.result,
		signatureMessages: rss.signatureMessages,
		validSignatures:   make(map[group.MemberIndex]operator.Signature),
	}

}

func (rss *resultSigningState) MemberIndex() group.MemberIndex {
	return rss.member.index
}

// signaturesVerificationState is the state during which group members verify all validSignatures
// that valid submitters sent over the broadcast channel in the previous state.
// Valid validSignatures are added to the state.
//
// State is part of phase 13 of the protocol.
type signaturesVerificationState struct {
	channel net.BroadcastChannel
	handle  chain.Handle

	member *SigningMember

	requestID *big.Int
	result    *relayChain.DKGResult

	signatureMessages []*DKGResultHashSignatureMessage
	validSignatures   map[group.MemberIndex]operator.Signature
}

func (svs *signaturesVerificationState) ActiveBlocks() int { return 0 }

func (svs *signaturesVerificationState) Initiate() error {
	signatures, err := svs.member.VerifyDKGResultSignatures(svs.signatureMessages)
	if err != nil {
		return err
	}

	svs.validSignatures = signatures
	return nil
}

func (svs *signaturesVerificationState) Receive(msg net.Message) error {
	return nil
}

func (svs *signaturesVerificationState) Next() signingState {
	return &resultSubmissionState{
		channel:    svs.channel,
		handle:     svs.handle,
		member:     NewSubmittingMember(svs.member.index),
		requestID:  svs.requestID,
		result:     svs.result,
		signatures: svs.validSignatures,
	}

}

func (svs *signaturesVerificationState) MemberIndex() group.MemberIndex {
	return svs.member.index
}

// resultSubmissionState is the state during which group members submit the dkg
// result to the chain.
//
// State covers phase 14 of the protocol.
type resultSubmissionState struct {
	channel net.BroadcastChannel
	handle  chain.Handle

	member *SubmittingMember

	requestID  *big.Int
	result     *relayChain.DKGResult
	signatures map[group.MemberIndex]operator.Signature
}

func (rss *resultSubmissionState) ActiveBlocks() int { return 3 }

func (rss *resultSubmissionState) Initiate() error {
	return rss.member.SubmitDKGResult(
		rss.requestID,
		rss.result,
		rss.signatures,
		rss.handle,
	)
}

func (rss *resultSubmissionState) Receive(msg net.Message) error {
	return nil
}

func (rss *resultSubmissionState) Next() signingState {
	return nil
}

func (rss *resultSubmissionState) MemberIndex() group.MemberIndex {
	return rss.member.index
}
