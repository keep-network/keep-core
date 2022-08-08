package dkg

import (
	"context"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

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
		finalizationStateActiveBlocks
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
	stateDuration := 15 * time.Second * time.Duration(stateBlocks)
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
	stateDuration := 15 * time.Second * time.Duration(stateBlocks)
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
	stateDuration := 15 * time.Second * time.Duration(stateBlocks)
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
	stateDuration := 15 * time.Second * time.Duration(stateBlocks)
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
