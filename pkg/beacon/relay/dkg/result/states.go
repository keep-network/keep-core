package result

import (
	"bytes"
	"context"
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
	channel      net.BroadcastChannel
	relayChain   relayChain.Interface
	blockCounter chain.BlockCounter

	member *SigningMember

	requestID *big.Int
	result    *relayChain.DKGResult

	signatureMessages []*DKGResultHashSignatureMessage

	signingStartBlockHeight uint64
}

func (rss *resultSigningState) DelayBlocks() uint64 {
	return state.MessagingStateDelayBlocks
}

func (rss *resultSigningState) ActiveBlocks() uint64 {
	return state.MessagingStateActiveBlocks
}

func (rss *resultSigningState) Initiate(context.Context) error {
	message, err := rss.member.SignDKGResult(rss.result, rss.relayChain)
	if err != nil {
		return err
	}
	if err := rss.channel.Send(message); err != nil {
		return err
	}
	return nil
}

func (rss *resultSigningState) Receive(msg net.Message) error {
	// The network layer determines the message sender's public key based on
	// the network client's pinned identity. The sender can not use any other
	// public key than the one it is identified with in the network.
	// Furthermore, the sender must possess the associated private key - each
	// network message is signed with it.
	//
	// The network layer rejects any message with an incorrect signature or
	// altered public key. By this point, we've conducted enough checks to
	// be very certain that the sender' public key presented in the network
	// net.Message is the correct one.
	//
	// In this final step, we compare the pinned network key with one used to
	// produce a signature over the DKG result hash. If the keys don't match,
	// it means that an incorrect key was used to sign DKG result hash and
	// the message should be rejected.
	isValidKeyUsed := func(phaseMessage *DKGResultHashSignatureMessage) bool {
		return bytes.Compare(
			operator.Marshal(phaseMessage.publicKey),
			msg.SenderPublicKey(),
		) == 0
	}

	switch signedMessage := msg.Payload().(type) {
	case *DKGResultHashSignatureMessage:
		if !group.IsMessageFromSelf(rss.member.index, signedMessage) &&
			group.IsSenderAccepted(rss.member, signedMessage) &&
			isValidKeyUsed(signedMessage) {
			rss.signatureMessages = append(rss.signatureMessages, signedMessage)
		}
	}

	return nil
}

func (rss *resultSigningState) Next() signingState {
	// set up the verification state, phase 13 part 2
	return &signaturesVerificationState{
		channel:           rss.channel,
		relayChain:        rss.relayChain,
		blockCounter:      rss.blockCounter,
		member:            rss.member,
		requestID:         rss.requestID,
		result:            rss.result,
		signatureMessages: rss.signatureMessages,
		validSignatures:   make(map[group.MemberIndex]operator.Signature),
		verificationStartBlockHeight: rss.signingStartBlockHeight +
			rss.DelayBlocks() +
			rss.ActiveBlocks(),
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
	channel      net.BroadcastChannel
	relayChain   relayChain.Interface
	blockCounter chain.BlockCounter

	member *SigningMember

	requestID *big.Int
	result    *relayChain.DKGResult

	signatureMessages []*DKGResultHashSignatureMessage
	validSignatures   map[group.MemberIndex]operator.Signature

	verificationStartBlockHeight uint64
}

func (svs *signaturesVerificationState) DelayBlocks() uint64 {
	return state.SilentStateDelayBlocks
}

func (svs *signaturesVerificationState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
}

func (svs *signaturesVerificationState) Initiate(context.Context) error {
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
		channel:      svs.channel,
		relayChain:   svs.relayChain,
		blockCounter: svs.blockCounter,
		member:       NewSubmittingMember(svs.member.index),
		requestID:    svs.requestID,
		result:       svs.result,
		signatures:   svs.validSignatures,
		submissionStartBlockHeight: svs.verificationStartBlockHeight +
			svs.DelayBlocks() +
			svs.ActiveBlocks(),
	}

}

func (svs *signaturesVerificationState) MemberIndex() group.MemberIndex {
	return svs.member.index
}

// resultSubmissionState is the state during which group members submit the dkg
// result to the chain. This state concludes the DKG protocol.
//
// State covers, the final phase, phase 14 of the protocol.
type resultSubmissionState struct {
	channel      net.BroadcastChannel
	relayChain   relayChain.Interface
	blockCounter chain.BlockCounter

	member *SubmittingMember

	requestID  *big.Int
	result     *relayChain.DKGResult
	signatures map[group.MemberIndex]operator.Signature

	submissionStartBlockHeight uint64
}

func (rss *resultSubmissionState) DelayBlocks() uint64 {
	return state.SilentStateDelayBlocks
}

func (rss *resultSubmissionState) ActiveBlocks() uint64 {
	// We do not exchange any messages in this phase. We publish result to the
	// chain but it is an action blocking all group members for the same time
	// - members exit when the first valid result is accepted by the chain.
	// How long it takes depends on the block step and group size.
	return state.SilentStateActiveBlocks
}

func (rss *resultSubmissionState) Initiate(ctx context.Context) error {
	return rss.member.SubmitDKGResult(
		ctx,
		rss.requestID,
		rss.result,
		rss.signatures,
		rss.relayChain,
		rss.blockCounter,
		rss.submissionStartBlockHeight,
	)
}

func (rss *resultSubmissionState) Receive(msg net.Message) error {
	return nil
}

func (rss *resultSubmissionState) Next() signingState {
	// returning nil represents this is the final state
	return nil
}

func (rss *resultSubmissionState) MemberIndex() group.MemberIndex {
	return rss.member.index
}
