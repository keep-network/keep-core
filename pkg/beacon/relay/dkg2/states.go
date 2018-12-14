package dkg2

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/net"
)

type keyGenerationState interface {
	activeBlocks() int

	initiate() error
	receive(msg net.Message) error
	nextState() (keyGenerationState, error)

	memberID() gjkr.MemberID
}

// initializationState is the starting state of key generation; it waits for
// activePeriod and then enters joinState. No messages are valid in this state.
type initializationState struct {
	channel net.BroadcastChannel
	member  *gjkr.EphemeralKeyPairGeneratingMember
}

func (is *initializationState) activeBlocks() int { return 1 }

func (is *initializationState) initiate() error {
	return nil
}

func (is *initializationState) receive(msg net.Message) error {
	return fmt.Errorf("unexpected message for initialization state: [%#v]", msg)
}

func (is *initializationState) nextState() (keyGenerationState, error) {
	return &joinState{is.channel, is.member}, nil
}

func (is *initializationState) memberID() gjkr.MemberID {
	return is.member.ID
}

// joinState is the state during which a member announces itself to the key
// generation broadcast channel to initiate the distributed protocol.
// `JoinMessage`s from other members are valid in this state.
type joinState struct {
	channel net.BroadcastChannel
	member  *gjkr.EphemeralKeyPairGeneratingMember
}

func (js *joinState) activeBlocks() int { return 1 }

func (js *joinState) initiate() error {
	return js.channel.Send(&gjkr.JoinMessage{SenderID: js.member.ID})
}

func (js *joinState) receive(msg net.Message) error {
	switch joinMsg := msg.Payload().(type) {
	case *gjkr.JoinMessage:
		if err := js.channel.RegisterIdentifier(
			msg.TransportSenderID(),
			joinMsg.SenderID,
		); err != nil {
			return err
		}

		js.member.AddToGroup(joinMsg.SenderID)
	}
	return nil
}

func (js *joinState) nextState() (keyGenerationState, error) {
	return &ephemeralKeyPairGeneratingState{
		channel: js.channel,
		member:  js.member,
	}, nil
}

func (js *joinState) memberID() gjkr.MemberID {
	return js.member.ID
}

// ephemeralKeyPairGeneratingState is the state during which members broadcast
// publish ephemeral keys generated for each other member in the group.
// `EphemeralPublicKeyMessage`s from other members are valid in this state.
type ephemeralKeyPairGeneratingState struct {
	channel net.BroadcastChannel
	member  *gjkr.EphemeralKeyPairGeneratingMember

	phaseMessages []*gjkr.EphemeralPublicKeyMessage
}

func (ekpgs *ephemeralKeyPairGeneratingState) activeBlocks() int { return 1 }

func (ekpgs *ephemeralKeyPairGeneratingState) initiate() error {
	message, err := ekpgs.member.GenerateEphemeralKeyPair()
	if err != nil {
		return fmt.Errorf("ephemeral key generation phase failed [%v]", err)
	}

	if err := ekpgs.channel.Send(message); err != nil {
		return fmt.Errorf("ephemeral key generation phase failed [%v]", err)
	}
	return nil
}

func (ekpgs *ephemeralKeyPairGeneratingState) receive(msg net.Message) error {
	switch publicKeyMessage := msg.Payload().(type) {
	case *gjkr.EphemeralPublicKeyMessage:
		if !messageFromSelf(ekpgs.memberID(), msg) {
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, publicKeyMessage)
		}

		return nil
	}

	return fmt.Errorf(
		"unexpected message for ephemeral key generation state: [%#v]",
		msg,
	)
}

func (ekpgs *ephemeralKeyPairGeneratingState) nextState() (keyGenerationState, error) {
	return &symmetricKeyGeneratingState{
		channel:               ekpgs.channel,
		member:                ekpgs.member.InitializeSymmetricKeyGeneration(),
		previousPhaseMessages: ekpgs.phaseMessages,
	}, nil
}

func (ekpgs *ephemeralKeyPairGeneratingState) memberID() gjkr.MemberID {
	return ekpgs.member.ID
}

type symmetricKeyGeneratingState struct {
	channel net.BroadcastChannel
	member  *gjkr.SymmetricKeyGeneratingMember

	previousPhaseMessages []*gjkr.EphemeralPublicKeyMessage
}

func (skgs *symmetricKeyGeneratingState) activeBlocks() int { return 1 }

func (skgs *symmetricKeyGeneratingState) initiate() error {
	return skgs.member.GenerateSymmetricKeys(skgs.previousPhaseMessages)
}

func (skgs *symmetricKeyGeneratingState) receive(msg net.Message) error {
	return fmt.Errorf(
		"unexpected message for symmetric key generation state: [%#v]",
		msg,
	)
}

func (skgs *symmetricKeyGeneratingState) nextState() (keyGenerationState, error) {
	return &committingState{
		channel: skgs.channel,
		member:  skgs.member.InitializeCommitting(),
	}, nil
}

func (skgs *symmetricKeyGeneratingState) memberID() gjkr.MemberID {
	return skgs.member.ID
}

type committingState struct {
	channel net.BroadcastChannel
	member  *gjkr.CommittingMember

	phaseSharesMessages      []*gjkr.PeerSharesMessage
	phaseCommitmentsMessages []*gjkr.MemberCommitmentsMessage
}

func (cs *committingState) activeBlocks() int { return 1 }

func (cs *committingState) initiate() error {
	sharesMsg, commitmentsMsg, err := cs.member.CalculateMembersSharesAndCommitments()
	if err != nil {
		return fmt.Errorf("committing phase failed [%v]", err)
	}

	if err := cs.channel.Send(sharesMsg); err != nil {
		return fmt.Errorf("committing phase failed [%v]", err)
	}

	if err := cs.channel.Send(commitmentsMsg); err != nil {
		return fmt.Errorf("committing phase failed [%v]", err)
	}

	return nil
}

func (cs *committingState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *gjkr.PeerSharesMessage:
		if !messageFromSelf(cs.memberID(), msg) {
			cs.phaseSharesMessages = append(cs.phaseSharesMessages, phaseMessage)
		}

		return nil

	case *gjkr.MemberCommitmentsMessage:
		if !messageFromSelf(cs.memberID(), msg) {
			cs.phaseCommitmentsMessages = append(
				cs.phaseCommitmentsMessages,
				phaseMessage,
			)
		}

		return nil
	}

	return fmt.Errorf("unexpected message for committing state: [%#v]", msg)
}

func (cs *committingState) nextState() (keyGenerationState, error) {
	return &commitmentsVerificationState{
		channel: cs.channel,
		member:  cs.member.InitializeCommitmentsVerification(),

		previousPhaseSharesMessages:      cs.phaseSharesMessages,
		previousPhaseCommitmentsMessages: cs.phaseCommitmentsMessages,
	}, nil
}

func (cs *committingState) memberID() gjkr.MemberID {
	return cs.member.ID
}

type commitmentsVerificationState struct {
	channel net.BroadcastChannel
	member  *gjkr.CommitmentsVerifyingMember

	previousPhaseSharesMessages      []*gjkr.PeerSharesMessage
	previousPhaseCommitmentsMessages []*gjkr.MemberCommitmentsMessage

	phaseAccusationsMessages []*gjkr.SecretSharesAccusationsMessage
}

func (cvs *commitmentsVerificationState) activeBlocks() int { return 1 }

func (cvs *commitmentsVerificationState) initiate() error {
	accusationsMsg, err := cvs.member.VerifyReceivedSharesAndCommitmentsMessages(
		cvs.previousPhaseSharesMessages,
		cvs.previousPhaseCommitmentsMessages,
	)
	if err != nil {
		return fmt.Errorf("commitments verification phase failed [%v]", err)
	}

	if err := cvs.channel.Send(accusationsMsg); err != nil {
		return fmt.Errorf("commitments verification phase failed [%v]", err)
	}

	return nil
}

func (cvs *commitmentsVerificationState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *gjkr.SecretSharesAccusationsMessage:
		if !messageFromSelf(cvs.memberID(), msg) {
			cvs.phaseAccusationsMessages = append(
				cvs.phaseAccusationsMessages,
				phaseMessage,
			)
		}

		return nil
	}

	return fmt.Errorf(
		"unexpected message for commitment verification state: [%#v]",
		msg,
	)
}

func (cvs *commitmentsVerificationState) nextState() (keyGenerationState, error) {
	return &sharesJustificationState{
		channel: cvs.channel,
		member:  cvs.member.InitializeSharesJustification(),

		previousPhaseAccusationsMessages: cvs.phaseAccusationsMessages,
	}, nil
}

func (cvs *commitmentsVerificationState) memberID() gjkr.MemberID {
	return cvs.member.ID
}

type sharesJustificationState struct {
	channel net.BroadcastChannel
	member  *gjkr.SharesJustifyingMember

	previousPhaseAccusationsMessages []*gjkr.SecretSharesAccusationsMessage
}

func (sjs *sharesJustificationState) activeBlocks() int { return 1 }

func (sjs *sharesJustificationState) initiate() error {
	disqualifiedMembers, err := sjs.member.ResolveSecretSharesAccusationsMessages(
		sjs.previousPhaseAccusationsMessages,
	)
	if err != nil {
		return fmt.Errorf("shares justification phase failed [%v]", err)
	}

	// TODO: Handle member disqualification
	fmt.Printf("disqualified members = %v\n", disqualifiedMembers)

	return nil
}

func (sjs *sharesJustificationState) receive(msg net.Message) error {
	return fmt.Errorf(
		"unexpected message for share justification phase: [%#v]",
		msg,
	)
}

func (sjs *sharesJustificationState) nextState() (keyGenerationState, error) {
	return nil, nil
}

func (sjs *sharesJustificationState) memberID() gjkr.MemberID {
	return sjs.member.ID
}

func messageFromSelf(selfMemberID gjkr.MemberID, message net.Message) bool {
	if senderID, ok := message.ProtocolSenderID().(gjkr.MemberID); ok {
		if senderID == selfMemberID {
			return true
		}
	}
	return false
}
