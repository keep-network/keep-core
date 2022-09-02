package signing

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
)

// ProtocolBlocks returns the total number of blocks it takes to execute
// all the required work defined by the signing protocol.
func ProtocolBlocks() uint64 {
	return ephemeralKeyPairStateDelayBlocks +
		ephemeralKeyPairStateActiveBlocks
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
		) && ekpgs.member.sessionID == phaseMessage.sessionID {
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Next() (state.State, error) {
	return &symmetricKeyGenerationState{
		channel:               ekpgs.channel,
		member:                ekpgs.member.initializeSymmetricKeyGeneration(),
		previousPhaseMessages: ekpgs.phaseMessages,
	}, nil
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

	if len(skgs.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(skgs.member.group.InactiveMemberIDs())
	}

	return skgs.member.generateSymmetricKeys(skgs.previousPhaseMessages)
}

func (skgs *symmetricKeyGenerationState) Receive(msg net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) Next() (state.State, error) {
	return nil, nil
}

func (skgs *symmetricKeyGenerationState) MemberIndex() group.MemberIndex {
	return skgs.member.id
}