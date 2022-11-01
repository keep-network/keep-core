package dkg

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

// ephemeralKeyPairGenerationState is the state during which members broadcast
// public ephemeral keys generated for other members of the group.
// `ephemeralPublicKeyMessage`s are valid in this state.
type ephemeralKeyPairGenerationState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *ephemeralKeyPairGeneratingMember
}

func (ekpgs *ephemeralKeyPairGenerationState) Initiate(ctx context.Context) error {
	message, err := ekpgs.member.generateEphemeralKeyPair()
	if err != nil {
		return err
	}

	if err := ekpgs.channel.Send(ctx, message); err != nil {
		return err
	}
	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if ekpgs.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && ekpgs.member.sessionID == protocolMessage.SessionID() {
			ekpgs.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) CanTransition() bool {
	messagingDone := len(receivedMessages[*ephemeralPublicKeyMessage](ekpgs.BaseAsyncState)) ==
		len(ekpgs.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (ekpgs *ephemeralKeyPairGenerationState) Next() (state.AsyncState, error) {
	return &symmetricKeyGenerationState{
		BaseAsyncState: ekpgs.BaseAsyncState,
		channel:        ekpgs.channel,
		member:         ekpgs.member.initializeSymmetricKeyGeneration(),
	}, nil
}

func (ekpgs *ephemeralKeyPairGenerationState) MemberIndex() group.MemberIndex {
	return ekpgs.member.id
}

// symmetricKeyGenerationState is the state during which members compute
// symmetric keys from the previously exchanged ephemeral public keys.
// No messages are valid in this state.
type symmetricKeyGenerationState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *symmetricKeyGeneratingMember
}

func (skgs *symmetricKeyGenerationState) Initiate(ctx context.Context) error {
	return skgs.member.generateSymmetricKeys(
		receivedMessages[*ephemeralPublicKeyMessage](skgs.BaseAsyncState),
	)
}

func (skgs *symmetricKeyGenerationState) Receive(msg net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) CanTransition() bool {
	return true
}

func (skgs *symmetricKeyGenerationState) Next() (state.AsyncState, error) {
	member, err := skgs.member.initializeTssRoundOne()
	if err != nil {
		return nil, fmt.Errorf(
			"cannot initialize TSS round one member: [%w]",
			err,
		)
	}

	return &tssRoundOneState{
		BaseAsyncState: skgs.BaseAsyncState,
		channel:        skgs.channel,
		member:         member,
	}, nil
}

func (skgs *symmetricKeyGenerationState) MemberIndex() group.MemberIndex {
	return skgs.member.id
}

// tssRoundOneState is the state during which members broadcast TSS
// commitments and the Paillier public key.
// `tssRoundOneMessage`s are valid in this state.
type tssRoundOneState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundOneMember
}

func (tros *tssRoundOneState) Initiate(ctx context.Context) error {
	message, err := tros.member.tssRoundOne(ctx)
	if err != nil {
		return err
	}

	if err := tros.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (tros *tssRoundOneState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if tros.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && tros.member.sessionID == protocolMessage.SessionID() {
			tros.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (tros *tssRoundOneState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundOneMessage](tros.BaseAsyncState)) ==
		len(tros.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (tros *tssRoundOneState) Next() (state.AsyncState, error) {
	return &tssRoundTwoState{
		BaseAsyncState: tros.BaseAsyncState,
		channel:        tros.channel,
		member:         tros.member.initializeTssRoundTwo(),
	}, nil
}

func (tros *tssRoundOneState) MemberIndex() group.MemberIndex {
	return tros.member.id
}

// tssRoundTwoState is the state during which members broadcast TSS
// shares and de-commitments.
// `tssRoundTwoMessage`s are valid in this state.
type tssRoundTwoState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundTwoMember
}

func (trts *tssRoundTwoState) Initiate(ctx context.Context) error {
	message, err := trts.member.tssRoundTwo(
		ctx,
		receivedMessages[*tssRoundOneMessage](trts.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trts.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trts *tssRoundTwoState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trts.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trts.member.sessionID == protocolMessage.SessionID() {
			trts.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trts *tssRoundTwoState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundTwoMessage](trts.BaseAsyncState)) ==
		len(trts.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trts *tssRoundTwoState) Next() (state.AsyncState, error) {
	return &tssRoundThreeState{
		BaseAsyncState: trts.BaseAsyncState,
		channel:        trts.channel,
		member:         trts.member.initializeTssRoundThree(),
	}, nil
}

func (trts *tssRoundTwoState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// tssRoundThreeState is the state during which members broadcast the TSS Paillier
// proof.
// `tssRoundThreeMessage`s are valid in this state.
type tssRoundThreeState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundThreeMember
}

func (trts *tssRoundThreeState) Initiate(ctx context.Context) error {
	message, err := trts.member.tssRoundThree(
		ctx,
		receivedMessages[*tssRoundTwoMessage](trts.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trts.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trts *tssRoundThreeState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trts.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trts.member.sessionID == protocolMessage.SessionID() {
			trts.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trts *tssRoundThreeState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundThreeMessage](trts.BaseAsyncState)) ==
		len(trts.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trts *tssRoundThreeState) Next() (state.AsyncState, error) {
	return &finalizationState{
		BaseAsyncState: trts.BaseAsyncState,
		channel:        trts.channel,
		member:         trts.member.initializeFinalization(),
	}, nil
}

func (trts *tssRoundThreeState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// finalizationState is the state during which members finalize the TSS process
// and prepare the distributed key generation result.
// `tssFinalizationMessage`s are valid in this state.
type finalizationState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *finalizingMember
}

func (fs *finalizationState) Initiate(ctx context.Context) error {
	message, err := fs.member.tssFinalize(
		ctx,
		receivedMessages[*tssRoundThreeMessage](fs.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := fs.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (fs *finalizationState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if fs.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && fs.member.sessionID == protocolMessage.SessionID() {
			fs.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (fs *finalizationState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssFinalizationMessage](fs.BaseAsyncState)) ==
		len(fs.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (fs *finalizationState) Next() (state.AsyncState, error) {
	return nil, nil
}

func (fs *finalizationState) MemberIndex() group.MemberIndex {
	return fs.member.id
}

func (fs *finalizationState) result() *Result {
	return fs.member.Result()
}

// resultSigningState is the state during which group members sign their
// preferred DKG result (by hashing their DKG result, and then signing the
// result), and share this over the broadcast channel.
type resultSigningState struct {
	*state.BaseAsyncState

	channel         net.BroadcastChannel
	resultSigner    ResultSigner
	resultSubmitter ResultSubmitter

	member *signingMember

	result *Result
}

func (rss *resultSigningState) Initiate(ctx context.Context) error {
	message, err := rss.member.signDKGResult(rss.result, rss.resultSigner)
	if err != nil {
		return err
	}

	if err := rss.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (rss *resultSigningState) Receive(netMessage net.Message) error {
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
	isValidKeyUsed := func(signatureMessage *resultSignatureMessage) bool {
		return bytes.Equal(signatureMessage.publicKey, netMessage.SenderPublicKey())
	}

	// As there is only one message type exchanged during result publication,
	// we can simplify the code and cast directly to the concrete type
	// `*resultSignatureMessage` instead of casting to the generic `message`.
	if signatureMessage, ok := netMessage.Payload().(*resultSignatureMessage); ok {
		if rss.member.shouldAcceptMessage(
			signatureMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && isValidKeyUsed(
			signatureMessage,
		) && rss.member.sessionID == signatureMessage.sessionID {
			rss.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (rss *resultSigningState) CanTransition() bool {
	// Although there is no hard requirement to expect signature messages
	// from all participants, it makes sense to do so because this is an
	// additional participant availability check that allows to maximize
	// the final count of active participants. Moreover, this check does not
	// bound the signing state to a fixed duration and one can move to the
	// next state as soon as possible.
	messagingDone := len(receivedMessages[*resultSignatureMessage](rss.BaseAsyncState)) ==
		len(rss.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (rss *resultSigningState) Next() (state.AsyncState, error) {
	return &signaturesVerificationState{
		BaseAsyncState:  rss.BaseAsyncState,
		channel:         rss.channel,
		resultSigner:    rss.resultSigner,
		resultSubmitter: rss.resultSubmitter,
		member:          rss.member,
		result:          rss.result,
		validSignatures: make(map[group.MemberIndex][]byte),
	}, nil
}

func (rss *resultSigningState) MemberIndex() group.MemberIndex {
	return rss.member.memberIndex
}

// signaturesVerificationState is the state during which group members verify
// all validSignatures that valid submitters sent over the broadcast channel in
// the previous state. Valid validSignatures are added to the state.
type signaturesVerificationState struct {
	*state.BaseAsyncState

	channel         net.BroadcastChannel
	resultSigner    ResultSigner
	resultSubmitter ResultSubmitter

	member *signingMember

	result *Result

	validSignatures map[group.MemberIndex][]byte
}

func (svs *signaturesVerificationState) Initiate(ctx context.Context) error {
	svs.validSignatures = svs.member.verifyDKGResultSignatures(
		receivedMessages[*resultSignatureMessage](svs.BaseAsyncState),
		svs.resultSigner,
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
	return &resultSubmissionState{
		BaseAsyncState:  svs.BaseAsyncState,
		channel:         svs.channel,
		resultSubmitter: svs.resultSubmitter,
		member:          svs.member.initializeSubmittingMember(),
		result:          svs.result,
		signatures:      svs.validSignatures,
	}, nil
}

func (svs *signaturesVerificationState) MemberIndex() group.MemberIndex {
	return svs.member.memberIndex
}

// resultSubmissionState is the state during which group members submit the dkg
// result to the chain. This state concludes the DKG protocol.
type resultSubmissionState struct {
	*state.BaseAsyncState

	channel         net.BroadcastChannel
	resultSubmitter ResultSubmitter

	member *submittingMember

	result     *Result
	signatures map[group.MemberIndex][]byte
}

func (rss *resultSubmissionState) Initiate(ctx context.Context) error {
	return rss.member.submitDKGResult(
		ctx,
		rss.result,
		rss.signatures,
		rss.resultSubmitter,
	)
}

func (rss *resultSubmissionState) Receive(msg net.Message) error {
	return nil
}

func (rss *resultSubmissionState) CanTransition() bool {
	return true
}

func (rss *resultSubmissionState) Next() (state.AsyncState, error) {
	// returning nil represents this is the final state
	return nil, nil
}

func (rss *resultSubmissionState) MemberIndex() group.MemberIndex {
	return rss.member.memberIndex
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
