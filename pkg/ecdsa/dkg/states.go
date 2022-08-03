package dkg

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
	tbtcchain "github.com/keep-network/keep-core/pkg/tbtc/chain"
)

// represents a given state in the state machine for signing DKG results
type signingState = state.State

const (
	silentStateDelayBlocks  = 0
	silentStateActiveBlocks = 0

	ephemeralKeyPairStateDelayBlocks  = 1
	ephemeralKeyPairStateActiveBlocks = 5

	tssRoundOneStateDelayBlocks  = 1
	tssRoundOneStateActiveBlocks = 5

	tssRoundTwoStateDelayBlocks  = 1
	tssRoundTwoStateActiveBlocks = 5

	tssRoundThreeStateDelayBlocks  = 1
	tssRoundThreeStateActiveBlocks = 5

	finalizationStateDelayBlocks  = 1
	finalizationStateActiveBlocks = 2

	resultSigningStateDelayBlocks  = 1
	resultSigningStateActiveBlocks = 5
)

// ProtocolBlocks returns the total number of blocks it takes to execute
// all the required work defined by the DKG protocol.
func ProtocolBlocks() uint64 {
	return ephemeralKeyPairStateDelayBlocks +
		ephemeralKeyPairStateActiveBlocks +
		tssRoundOneStateDelayBlocks +
		tssRoundOneStateActiveBlocks +
		tssRoundTwoStateDelayBlocks +
		tssRoundTwoStateActiveBlocks +
		tssRoundThreeStateDelayBlocks +
		tssRoundThreeStateActiveBlocks +
		finalizationStateDelayBlocks +
		finalizationStateActiveBlocks +
		resultSigningStateDelayBlocks +
		resultSigningStateActiveBlocks
}

// ephemeralKeyPairGenerationState is the state during which members broadcast
// public ephemeral keys generated for other members of the group.
// `ephemeralPublicKeyMessage`s are valid in this state.
type ephemeralKeyPairGenerationState struct {
	channel net.BroadcastChannel
	member  *ephemeralKeyPairGeneratingMember

	phaseMessages []*ephemeralPublicKeyMessage
}

func (ekpgs *ephemeralKeyPairGenerationState) DelayBlocks() uint64 {
	return ephemeralKeyPairStateDelayBlocks
}

func (ekpgs *ephemeralKeyPairGenerationState) ActiveBlocks() uint64 {
	return ephemeralKeyPairStateActiveBlocks
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

func (ekpgs *ephemeralKeyPairGenerationState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *ephemeralPublicKeyMessage:
		if ekpgs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) {
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Next() state.State {
	return &symmetricKeyGenerationState{
		channel:               ekpgs.channel,
		member:                ekpgs.member.initializeSymmetricKeyGeneration(),
		previousPhaseMessages: ekpgs.phaseMessages,
	}
}

func (ekpgs *ephemeralKeyPairGenerationState) MemberIndex() group.MemberIndex {
	return ekpgs.member.id
}

// symmetricKeyGenerationState is the state during which members compute
// symmetric keys from the previously exchanged ephemeral public keys.
// No messages are valid in this state.
type symmetricKeyGenerationState struct {
	channel net.BroadcastChannel
	member  *symmetricKeyGeneratingMember

	previousPhaseMessages []*ephemeralPublicKeyMessage
}

func (skgs *symmetricKeyGenerationState) DelayBlocks() uint64 {
	return silentStateDelayBlocks
}

func (skgs *symmetricKeyGenerationState) ActiveBlocks() uint64 {
	return silentStateActiveBlocks
}

func (skgs *symmetricKeyGenerationState) Initiate(ctx context.Context) error {
	skgs.member.MarkInactiveMembers(skgs.previousPhaseMessages)

	if len(skgs.member.group.OperatingMemberIDs()) != skgs.member.group.GroupSize() {
		return fmt.Errorf("inactive members detected")
	}

	return skgs.member.generateSymmetricKeys(skgs.previousPhaseMessages)
}

func (skgs *symmetricKeyGenerationState) Receive(msg net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) Next() state.State {
	return &tssRoundOneState{
		channel: skgs.channel,
		member:  skgs.member.initializeTssRoundOne(),
	}
}

func (skgs *symmetricKeyGenerationState) MemberIndex() group.MemberIndex {
	return skgs.member.id
}

// tssRoundOneState is the state during which members broadcast TSS
// commitments and the Paillier public key.
// `tssRoundOneMessage`s are valid in this state.
type tssRoundOneState struct {
	channel net.BroadcastChannel
	member  *tssRoundOneMember

	phaseMessages []*tssRoundOneMessage
}

func (tros *tssRoundOneState) DelayBlocks() uint64 {
	return tssRoundOneStateDelayBlocks
}

func (tros *tssRoundOneState) ActiveBlocks() uint64 {
	return tssRoundOneStateActiveBlocks
}

func (tros *tssRoundOneState) Initiate(ctx context.Context) error {
	// The ctx instance passed as Initiate argument is scoped to the lifetime
	// of the current state. However, the Initiate method is blocking and the
	// ctx instance is cancelled properly only after Initiate returns. Because
	// of that, we cannot use ctx as round timeout signal as Initiate would
	// hang forever if something goes wrong. To avoid such a resource leak,
	// we set a round timeout based on state block duration and an average
	// block time. The exact duration doesn't need to be super-accurate because
	// if the timeout is hit, the execution will fail anyway. We just want
	// to give enough time for round computation and make sure the round
	// terminates regardless of the result.
	stateBlocks := tros.DelayBlocks() + tros.ActiveBlocks()
	stateDuration := 15 * time.Duration(stateBlocks)
	roundCtx, roundCtxCancel := context.WithTimeout(ctx, stateDuration)
	defer roundCtxCancel()

	message, err := tros.member.tssRoundOne(roundCtx)
	if err != nil {
		return err
	}

	if err := tros.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (tros *tssRoundOneState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundOneMessage:
		if tros.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && tros.member.sessionID == phaseMessage.sessionID {
			tros.phaseMessages = append(tros.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (tros *tssRoundOneState) Next() state.State {
	return &tssRoundTwoState{
		channel:               tros.channel,
		member:                tros.member.initializeTssRoundTwo(),
		previousPhaseMessages: tros.phaseMessages,
	}
}

func (tros *tssRoundOneState) MemberIndex() group.MemberIndex {
	return tros.member.id
}

// tssRoundOneState is the state during which members broadcast TSS
// shares and de-commitments.
// `tssRoundTwoMessage`s are valid in this state.
type tssRoundTwoState struct {
	channel net.BroadcastChannel
	member  *tssRoundTwoMember

	previousPhaseMessages []*tssRoundOneMessage

	phaseMessages []*tssRoundTwoMessage
}

func (trts *tssRoundTwoState) DelayBlocks() uint64 {
	return tssRoundTwoStateDelayBlocks
}

func (trts *tssRoundTwoState) ActiveBlocks() uint64 {
	return tssRoundTwoStateActiveBlocks
}

func (trts *tssRoundTwoState) Initiate(ctx context.Context) error {
	trts.member.MarkInactiveMembers(trts.previousPhaseMessages)

	if len(trts.member.group.OperatingMemberIDs()) != trts.member.group.GroupSize() {
		return fmt.Errorf("inactive members detected")
	}

	// The ctx instance passed as Initiate argument is scoped to the lifetime
	// of the current state. However, the Initiate method is blocking and the
	// ctx instance is cancelled properly only after Initiate returns. Because
	// of that, we cannot use ctx as round timeout signal as Initiate would
	// hang forever if something goes wrong. To avoid such a resource leak,
	// we set a round timeout based on state block duration and an average
	// block time. The exact duration doesn't need to be super-accurate because
	// if the timeout is hit, the execution will fail anyway. We just want
	// to give enough time for round computation and make sure the round
	// terminates regardless of the result.
	stateBlocks := trts.DelayBlocks() + trts.ActiveBlocks()
	stateDuration := 15 * time.Duration(stateBlocks)
	roundCtx, roundCtxCancel := context.WithTimeout(ctx, stateDuration)
	defer roundCtxCancel()

	message, err := trts.member.tssRoundTwo(roundCtx, trts.previousPhaseMessages)
	if err != nil {
		return err
	}

	if err := trts.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trts *tssRoundTwoState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundTwoMessage:
		if trts.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trts.member.sessionID == phaseMessage.sessionID {
			trts.phaseMessages = append(trts.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trts *tssRoundTwoState) Next() state.State {
	return &tssRoundThreeState{
		channel:               trts.channel,
		member:                trts.member.initializeTssRoundThree(),
		previousPhaseMessages: trts.phaseMessages,
	}
}

func (trts *tssRoundTwoState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// tssRoundOneState is the state during which members broadcast the TSS Paillier
// proof.
// `tssRoundThreeMessage`s are valid in this state.
type tssRoundThreeState struct {
	channel net.BroadcastChannel
	member  *tssRoundThreeMember

	previousPhaseMessages []*tssRoundTwoMessage

	phaseMessages []*tssRoundThreeMessage
}

func (trts *tssRoundThreeState) DelayBlocks() uint64 {
	return tssRoundThreeStateDelayBlocks
}

func (trts *tssRoundThreeState) ActiveBlocks() uint64 {
	return tssRoundThreeStateActiveBlocks
}

func (trts *tssRoundThreeState) Initiate(ctx context.Context) error {
	trts.member.MarkInactiveMembers(trts.previousPhaseMessages)

	if len(trts.member.group.OperatingMemberIDs()) != trts.member.group.GroupSize() {
		return fmt.Errorf("inactive members detected")
	}

	// The ctx instance passed as Initiate argument is scoped to the lifetime
	// of the current state. However, the Initiate method is blocking and the
	// ctx instance is cancelled properly only after Initiate returns. Because
	// of that, we cannot use ctx as round timeout signal as Initiate would
	// hang forever if something goes wrong. To avoid such a resource leak,
	// we set a round timeout based on state block duration and an average
	// block time. The exact duration doesn't need to be super-accurate because
	// if the timeout is hit, the execution will fail anyway. We just want
	// to give enough time for round computation and make sure the round
	// terminates regardless of the result.
	stateBlocks := trts.DelayBlocks() + trts.ActiveBlocks()
	stateDuration := 15 * time.Duration(stateBlocks)
	roundCtx, roundCtxCancel := context.WithTimeout(ctx, stateDuration)
	defer roundCtxCancel()

	message, err := trts.member.tssRoundThree(roundCtx, trts.previousPhaseMessages)
	if err != nil {
		return err
	}

	if err := trts.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trts *tssRoundThreeState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundThreeMessage:
		if trts.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trts.member.sessionID == phaseMessage.sessionID {
			trts.phaseMessages = append(trts.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trts *tssRoundThreeState) Next() state.State {
	return &finalizationState{
		channel:               trts.channel,
		member:                trts.member.initializeFinalization(),
		previousPhaseMessages: trts.phaseMessages,
	}
}

func (trts *tssRoundThreeState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// finalizationState is the last state of the DKG protocol - in this state,
// distributed key generation is completed. No messages are valid in this state.
//
// State prepares a result to that is returned to the caller.
type finalizationState struct {
	channel net.BroadcastChannel
	member  *finalizingMember

	previousPhaseMessages []*tssRoundThreeMessage
}

func (fs *finalizationState) DelayBlocks() uint64 {
	return finalizationStateDelayBlocks
}

func (fs *finalizationState) ActiveBlocks() uint64 {
	return finalizationStateActiveBlocks
}

func (fs *finalizationState) Initiate(ctx context.Context) error {
	fs.member.MarkInactiveMembers(fs.previousPhaseMessages)

	if len(fs.member.group.OperatingMemberIDs()) != fs.member.group.GroupSize() {
		return fmt.Errorf("inactive members detected")
	}

	// The ctx instance passed as Initiate argument is scoped to the lifetime
	// of the current state. However, the Initiate method is blocking and the
	// ctx instance is cancelled properly only after Initiate returns. Because
	// of that, we cannot use ctx as round timeout signal as Initiate would
	// hang forever if something goes wrong. To avoid such a resource leak,
	// we set a round timeout based on state block duration and an average
	// block time. The exact duration doesn't need to be super-accurate because
	// if the timeout is hit, the execution will fail anyway. We just want
	// to give enough time for round computation and make sure the round
	// terminates regardless of the result.
	stateBlocks := fs.DelayBlocks() + fs.ActiveBlocks()
	stateDuration := 15 * time.Duration(stateBlocks)
	roundCtx, roundCtxCancel := context.WithTimeout(ctx, stateDuration)
	defer roundCtxCancel()

	return fs.member.tssFinalize(roundCtx, fs.previousPhaseMessages)
}

func (fs *finalizationState) Receive(msg net.Message) error {
	return nil
}

func (fs *finalizationState) Next() state.State {
	return nil
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
	channel      net.BroadcastChannel
	tbtcChain    tbtcchain.Chain
	blockCounter chain.BlockCounter

	member *SigningMember

	result *tbtcchain.DKGResult

	signatureMessages []*dkgResultHashSignatureMessage

	signingStartBlockHeight uint64
}

func (rss *resultSigningState) DelayBlocks() uint64 {
	return resultSigningStateDelayBlocks
}

func (rss *resultSigningState) ActiveBlocks() uint64 {
	return resultSigningStateActiveBlocks
}

func (rss *resultSigningState) Initiate(ctx context.Context) error {
	message, err := rss.member.SignDKGResult(rss.result, rss.tbtcChain)
	if err != nil {
		return err
	}
	if err := rss.channel.Send(ctx, message); err != nil {
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
	isValidKeyUsed := func(phaseMessage *dkgResultHashSignatureMessage) bool {
		return bytes.Compare(phaseMessage.publicKey, msg.SenderPublicKey()) == 0
	}

	switch signedMessage := msg.Payload().(type) {
	case *dkgResultHashSignatureMessage:
		if rss.member.shouldAcceptMessage(
			signedMessage.SenderID(),
			msg.SenderPublicKey(),
		) && isValidKeyUsed(signedMessage) {
			rss.signatureMessages = append(rss.signatureMessages, signedMessage)
		}
	}

	return nil
}

func (rss *resultSigningState) Next() signingState {
	return &signaturesVerificationState{
		channel:           rss.channel,
		tbtcChain:         rss.tbtcChain,
		blockCounter:      rss.blockCounter,
		member:            rss.member,
		result:            rss.result,
		signatureMessages: rss.signatureMessages,
		validSignatures:   make(map[group.MemberIndex][]byte),
		verificationStartBlockHeight: rss.signingStartBlockHeight +
			rss.DelayBlocks() +
			rss.ActiveBlocks(),
	}
}

func (rss *resultSigningState) MemberIndex() group.MemberIndex {
	return rss.member.index
}

// signaturesVerificationState is the state during which group members verify
// all validSignatures that valid submitters sent over the broadcast channel in
// the previous state. Valid validSignatures are added to the state.
type signaturesVerificationState struct {
	channel      net.BroadcastChannel
	tbtcChain    tbtcchain.Chain
	blockCounter chain.BlockCounter

	member *SigningMember

	result *tbtcchain.DKGResult

	signatureMessages []*dkgResultHashSignatureMessage
	validSignatures   map[group.MemberIndex][]byte

	verificationStartBlockHeight uint64
}

func (svs *signaturesVerificationState) DelayBlocks() uint64 {
	return state.SilentStateDelayBlocks
}

func (svs *signaturesVerificationState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
}

func (svs *signaturesVerificationState) Initiate(ctx context.Context) error {
	signatures, err := svs.member.VerifyDKGResultSignatures(
		svs.signatureMessages,
		svs.tbtcChain.Signing(),
	)
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
		tbtcChain:    svs.tbtcChain,
		blockCounter: svs.blockCounter,
		member:       NewSubmittingMember(svs.member.logger, svs.member.index),
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
type resultSubmissionState struct {
	channel      net.BroadcastChannel
	tbtcChain    tbtcchain.Chain
	blockCounter chain.BlockCounter

	member *SubmittingMember

	result     *tbtcchain.DKGResult
	signatures map[group.MemberIndex][]byte

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
		rss.result,
		rss.signatures,
		rss.tbtcChain,
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
