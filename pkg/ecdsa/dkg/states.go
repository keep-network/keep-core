package dkg

import (
	"context"

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
		tssRoundThreeStateActiveBlocks
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
// commitments and Paillier public keys generated for other members of
// the group.
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
	message, err := tros.member.tssRoundOne(ctx)
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
		if !group.IsMessageFromSelf(tros.member.id, phaseMessage) &&
			group.IsSenderValid(tros.member, phaseMessage, msg.SenderPublicKey()) &&
			group.IsSenderAccepted(tros.member, phaseMessage) {
			tros.phaseMessages = append(tros.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (tros *tssRoundOneState) Next() state.State {
	//TODO implement me
	panic("implement me")
}

func (tros *tssRoundOneState) MemberIndex() group.MemberIndex {
	return tros.member.id
}

type tssRoundTwoState struct {
}

func (trts *tssRoundTwoState) DelayBlocks() uint64 {
	return tssRoundTwoStateDelayBlocks
}

func (trts *tssRoundTwoState) ActiveBlocks() uint64 {
	return tssRoundTwoStateActiveBlocks
}

func (trts *tssRoundTwoState) Initiate(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (trts *tssRoundTwoState) Receive(msg net.Message) error {
	//TODO implement me
	panic("implement me")
}

func (trts *tssRoundTwoState) Next() state.State {
	//TODO implement me
	panic("implement me")
}

func (trts *tssRoundTwoState) MemberIndex() group.MemberIndex {
	//TODO implement me
	panic("implement me")
}

type tssRoundThreeState struct {
}

func (trts *tssRoundThreeState) DelayBlocks() uint64 {
	return tssRoundThreeStateDelayBlocks
}

func (trts *tssRoundThreeState) ActiveBlocks() uint64 {
	return tssRoundThreeStateActiveBlocks
}

func (trts *tssRoundThreeState) Initiate(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (trts *tssRoundThreeState) Receive(msg net.Message) error {
	//TODO implement me
	panic("implement me")
}

func (trts *tssRoundThreeState) Next() state.State {
	//TODO implement me
	panic("implement me")
}

func (trts *tssRoundThreeState) MemberIndex() group.MemberIndex {
	//TODO implement me
	panic("implement me")
}

// finalizationState is the last state of the DKG protocol - in this state,
// distributed key generation is completed. No messages are valid in this state.
//
// State prepares a result to that is returned to the caller.
type finalizationState struct {
	channel net.BroadcastChannel
	member  *finalizingMember
}

func (fs *finalizationState) DelayBlocks() uint64 {
	return silentStateDelayBlocks
}

func (fs *finalizationState) ActiveBlocks() uint64 {
	return silentStateActiveBlocks
}

func (fs *finalizationState) Initiate(ctx context.Context) error {
	return nil
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
