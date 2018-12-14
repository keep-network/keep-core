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

// TODO: rename to isMessageFromSelf
func messageFromSelf(selfMemberID gjkr.MemberID, message net.Message) bool {
	if senderID, ok := message.ProtocolSenderID().(gjkr.MemberID); ok {
		if senderID == selfMemberID {
			return true
		}
	}
	return false
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
	return fmt.Errorf("unexpected message for initialization phase: [%#v]", msg)
}

func (is *initializationState) nextState() (keyGenerationState, error) {
	return &joinState{is.channel, is.member}, nil
}

func (is *initializationState) memberID() gjkr.MemberID {
	return is.member.ID
}

// joinState is the state during which a member announces itself to the key
// generation broadcast channel to initiate the distributed protocol.
// `gjkr.JoinMessage`s are valid in this state.
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
// public ephemeral keys generated for other members of the group.
// `gjkr.EphemeralPublicKeyMessage`s are valid in this state.
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
		"unexpected message for ephemeral key generation phase: [%#v]",
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

// symmetricKeyGeneratingState is the state during which members compute
// symmetric keys from the previously exchanged ephemeral public keys.
// No messages are valid in this state.
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
		"unexpected message for symmetric key generation phase: [%#v]",
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

// committingState is the state during which members compute their individual
// shares and commitments to those shares. Two messages are valid in this state:
// - `gjkr.PeerSharesMessage`
// - `gjkr.MemberCommitmentsMessage`
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

	return fmt.Errorf("unexpected message for committing phase: [%#v]", msg)
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

// commitmentsVerificationState is the state during which members validate
// shares and commitments computed and published by other members in the
// previous phase. `gjkr.SecretShareAccusationMessage`s are valid in this state.
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
		"unexpected message for commitment verification phase: [%#v]",
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

// sharesJustificationState is the state during which members resolve
// accusations published by other group members in the previous state.
// No messages are valid in this state.
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
	return &qualifiedState{
		channel: sjs.channel,
		member:  sjs.member.InitializeQualified(),
	}, nil
}

func (sjs *sharesJustificationState) memberID() gjkr.MemberID {
	return sjs.member.ID
}

// qualifiedState is the state during which group members combine all valid
// secret shares published by other group members in the previous states.
// No messages are valid in this state.
type qualifiedState struct {
	channel net.BroadcastChannel
	member  *gjkr.QualifiedMember
}

func (qs *qualifiedState) activeBlocks() int { return 1 }

func (qs *qualifiedState) initiate() error {
	qs.member.CombineMemberShares()
	return nil
}

func (qs *qualifiedState) receive(msg net.Message) error {
	return fmt.Errorf("unexpected message for qualified phase: [%#v]", msg)
}

func (qs *qualifiedState) nextState() (keyGenerationState, error) {
	return &pointsSharingState{
		channel: qs.channel,
		member:  qs.member.InitializeSharing(),
	}, nil
}

func (qs *qualifiedState) memberID() gjkr.MemberID {
	return qs.member.ID
}

// pointsSharingState is the state during which group members calculate and
// publish their public key share points.
// `gjkr.MemberPublicKeySharePointsMessage`s are valid in this state.
type pointsSharingState struct {
	channel net.BroadcastChannel
	member  *gjkr.SharingMember // TODO: SharingMember should be renamed to PointsSharingMember

	phaseMessages []*gjkr.MemberPublicKeySharePointsMessage
}

func (pss *pointsSharingState) activeBlocks() int { return 1 }

func (pss *pointsSharingState) initiate() error {
	message := pss.member.CalculatePublicKeySharePoints()
	if err := pss.channel.Send(message); err != nil {
		return fmt.Errorf("points sharing phase failed [%v]", err)
	}

	return nil
}

func (pss *pointsSharingState) receive(msg net.Message) error {
	switch pointsMessage := msg.Payload().(type) {
	case *gjkr.MemberPublicKeySharePointsMessage:
		if !messageFromSelf(pss.memberID(), msg) {
			pss.phaseMessages = append(pss.phaseMessages, pointsMessage)
		}
		return nil
	}

	return fmt.Errorf(
		"unexpected message for points sharing phase: [%#v]",
		msg,
	)
}

func (pss *pointsSharingState) nextState() (keyGenerationState, error) {
	return &pointsValidationState{
		channel: pss.channel,
		member:  pss.member,

		previousPhaseMessages: pss.phaseMessages,
	}, nil
}

func (pss *pointsSharingState) memberID() gjkr.MemberID {
	return pss.member.ID
}

// pointsValidationState is the state during which group members validate
// public key share points published by other group members in the previous
// state. `gjkr.PointsAccusationsMessage`s are valid in this state.
type pointsValidationState struct {
	channel net.BroadcastChannel
	member  *gjkr.SharingMember // TODO: split validation logic into PointsValidatingMember

	previousPhaseMessages []*gjkr.MemberPublicKeySharePointsMessage

	phaseMessages []*gjkr.PointsAccusationsMessage
}

func (pvs *pointsValidationState) activeBlocks() int { return 1 }

func (pvs *pointsValidationState) initiate() error {
	accusationMsg, err := pvs.member.VerifyPublicKeySharePoints(
		pvs.previousPhaseMessages,
	)
	if err != nil {
		return fmt.Errorf("points validation phase failed [%v]", err)
	}

	if err := pvs.channel.Send(accusationMsg); err != nil {
		return fmt.Errorf("points validation phase failed [%v]", err)
	}

	return nil
}

func (pvs *pointsValidationState) receive(msg net.Message) error {
	switch pointsAccusationMessage := msg.Payload().(type) {
	case *gjkr.PointsAccusationsMessage:
		if !messageFromSelf(pvs.memberID(), msg) {
			pvs.phaseMessages = append(pvs.phaseMessages, pointsAccusationMessage)
		}

		return nil
	}

	return fmt.Errorf(
		"unexpected message for points validation phase: [%#v]",
		msg,
	)
}

func (pvs *pointsValidationState) nextState() (keyGenerationState, error) {
	return &pointsJustificationState{
		channel: pvs.channel,
		member:  pvs.member.InitializePointsJustification(),

		previousPhaseMessages: pvs.phaseMessages,
	}, nil
}

func (pvs *pointsValidationState) memberID() gjkr.MemberID {
	return pvs.member.ID
}

// pointsJustificationState is the state during which group members resolve
// accusations published by other group members in the previous state.
// No messages are valid in this state.
type pointsJustificationState struct {
	channel net.BroadcastChannel
	member  *gjkr.PointsJustifyingMember

	previousPhaseMessages []*gjkr.PointsAccusationsMessage
}

func (pjs *pointsJustificationState) activeBlocks() int { return 1 }

func (pjs *pointsJustificationState) initiate() error {
	disqualifiedMembers, err := pjs.member.ResolvePublicKeySharePointsAccusationsMessages(
		pjs.previousPhaseMessages,
	)
	if err != nil {
		return fmt.Errorf("points justification phase failed [%v]", err)
	}

	// TODO: Handle member disqualification
	fmt.Printf("disqualified members = %v\n", disqualifiedMembers)

	return nil
}

func (pjs *pointsJustificationState) receive(msg net.Message) error {
	return fmt.Errorf(
		"unexpected message for points justification phase: [%#v]",
		msg,
	)
}

func (pjs *pointsJustificationState) nextState() (keyGenerationState, error) {
	return &reconstructionState{
		channel: pjs.channel,
		member:  pjs.member.InitializeReconstruction(),
	}, nil
}

func (pjs *pointsJustificationState) memberID() gjkr.MemberID {
	return pjs.member.ID
}

// reconstructionState is the state during which group members reconstruct
// individual keys of members disqualified in previous states. No messages are
// valid in this state.
type reconstructionState struct {
	channel net.BroadcastChannel
	member  *gjkr.ReconstructingMember
}

func (rp *reconstructionState) activeBlocks() int { return 1 }

func (rp *reconstructionState) initiate() error {
	// TODO: implement once member disqualification will be ready
	return nil
}

func (rp *reconstructionState) receive(msg net.Message) error {
	return fmt.Errorf("unexpected message for reconstruction phase: [%#v]", msg)
}

func (rp *reconstructionState) nextState() (keyGenerationState, error) {
	return &combiningState{
		channel: rp.channel,
		member:  rp.member.InitializeCombining(),
	}, nil
}

func (rp *reconstructionState) memberID() gjkr.MemberID {
	return rp.member.ID
}

// combiningState is the final state of GJKR protocol during which group
// members combine together all qualified key shares to form a group public key.
// No messages are valid in this state.
type combiningState struct {
	channel net.BroadcastChannel
	member  *gjkr.CombiningMember
}

func (cs *combiningState) activeBlocks() int { return 1 }

func (cs *combiningState) initiate() error {
	cs.member.CombineGroupPublicKey()
	return nil
}

func (cs *combiningState) receive(msg net.Message) error {
	return fmt.Errorf("unexpected message for combining phase: [%#v]", msg)
}

func (cs *combiningState) nextState() (keyGenerationState, error) {
	return nil, nil
}

func (cs *combiningState) memberID() gjkr.MemberID {
	return cs.member.ID
}
