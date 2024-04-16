package inactivity

import (
	"bytes"
	"context"
	"strconv"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

// claimSigningState is the state during which group members sign their
// preferred inactivity claim (by hashing their inactivity, and then signing the
// result), and share this over the broadcast channel.
type claimSigningState struct {
	*state.BaseAsyncState

	channel        net.BroadcastChannel
	claimSigner    ClaimSigner
	claimSubmitter ClaimSubmitter

	member *signingMember

	claim *Claim
}

func (css *claimSigningState) Initiate(ctx context.Context) error {
	message, err := css.member.signClaim(css.claim, css.claimSigner)
	if err != nil {
		return err
	}

	if err := css.channel.Send(
		ctx,
		message,
		net.BackoffRetransmissionStrategy,
	); err != nil {
		return err
	}

	return nil
}

func (css *claimSigningState) Receive(netMessage net.Message) error {
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
	// produce a signature over the inactivity claim hash. If the keys don't
	// match, it means that an incorrect key was used to sign inactivity claim
	// hash and the message should be rejected.
	isValidKeyUsed := func(signatureMessage *claimSignatureMessage) bool {
		return bytes.Equal(signatureMessage.publicKey, netMessage.SenderPublicKey())
	}

	// As there is only one message type exchanged during result publication,
	// we can simplify the code and cast directly to the concrete type
	// `*resultSignatureMessage` instead of casting to the generic `message`.
	if signatureMessage, ok := netMessage.Payload().(*claimSignatureMessage); ok {
		if css.member.shouldAcceptMessage(
			signatureMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && isValidKeyUsed(
			signatureMessage,
		) && css.member.sessionID == signatureMessage.sessionID {
			css.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (css *claimSigningState) CanTransition() bool {
	// Although there is no hard requirement to expect signature messages
	// from all participants, it makes sense to do so because this is an
	// additional participant availability check that allows to maximize
	// the final count of active participants. Moreover, this check does not
	// bound the signing state to a fixed duration and one can move to the
	// next state as soon as possible.
	messagingDone := len(receivedMessages[*claimSignatureMessage](css.BaseAsyncState)) ==
		len(css.member.group.OperatingMemberIndexes())-1

	// TODO: Modify the above code so that only 51 members are needed. Since it
	//       is executed after a failed heartbeat, we cannot expect all the
	//       members to sign the claim. In the future consider taking the number
	//       of active signers from the heartbeat procedure.

	return messagingDone
}

func (css *claimSigningState) Next() (state.AsyncState, error) {
	return &signaturesVerificationState{
		BaseAsyncState:  css.BaseAsyncState,
		channel:         css.channel,
		claimSigner:     css.claimSigner,
		claimSubmitter:  css.claimSubmitter,
		member:          css.member,
		claim:           css.claim,
		validSignatures: make(map[group.MemberIndex][]byte),
	}, nil
}

func (css *claimSigningState) MemberIndex() group.MemberIndex {
	return css.member.memberIndex
}

type signaturesVerificationState struct {
	*state.BaseAsyncState

	channel        net.BroadcastChannel
	claimSigner    ClaimSigner
	claimSubmitter ClaimSubmitter

	member *signingMember

	claim *Claim

	validSignatures map[group.MemberIndex][]byte
}

func (svs *signaturesVerificationState) Initiate(ctx context.Context) error {
	svs.validSignatures = svs.member.verifyInactivityClaimSignatures(
		receivedMessages[*claimSignatureMessage](svs.BaseAsyncState),
		svs.claimSigner,
	)
	return nil
}

func (svs *signaturesVerificationState) Receive(msg net.Message) error {
	return nil
}

func (svs *signaturesVerificationState) CanTransition() bool {
	return true
}

func (svs *signaturesVerificationState) Next() (state.AsyncState, error) {
	return &claimSubmissionState{
		BaseAsyncState: svs.BaseAsyncState,
		channel:        svs.channel,
		claimSubmitter: svs.claimSubmitter,
		member:         svs.member.initializeSubmittingMember(),
		claim:          svs.claim,
		signatures:     svs.validSignatures,
	}, nil
}

func (svs *signaturesVerificationState) MemberIndex() group.MemberIndex {
	return svs.member.memberIndex
}

type claimSubmissionState struct {
	*state.BaseAsyncState

	channel        net.BroadcastChannel
	claimSubmitter ClaimSubmitter

	member *submittingMember

	claim      *Claim
	signatures map[group.MemberIndex][]byte
}

func (css *claimSubmissionState) Initiate(ctx context.Context) error {
	return css.member.submitClaim(
		ctx,
		css.claim,
		css.signatures,
		css.claimSubmitter,
	)
}

func (css *claimSubmissionState) Receive(msg net.Message) error {
	return nil
}

func (css *claimSubmissionState) CanTransition() bool {
	return true
}

func (css *claimSubmissionState) Next() (state.AsyncState, error) {
	// returning nil represents this is the final state
	return nil, nil
}

func (css *claimSubmissionState) MemberIndex() group.MemberIndex {
	return css.member.memberIndex
}

// receivedMessages returns all messages of type T that have been received
// and validated so far. Returned messages are deduplicated so there is a
// guarantee that only one message of the given type is returned for the
// given sender.
func receivedMessages[T message](base *state.BaseAsyncState) []T {
	var messageTemplate T

	payloads := state.ExtractMessagesPayloads[T](base, messageTemplate.Type())

	return state.DeduplicateMessagesPayloads(
		payloads,
		func(message T) string {
			return strconv.Itoa(int(message.SenderID()))
		},
	)
}
