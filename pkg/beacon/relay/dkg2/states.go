package dkg2

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/net"
)

type keyGenerationState interface {
	// activeBlocks returns the number of blocks during which the current state
	// is active. Blocks are counted after the initiation process of the
	// current state has completed.
	activeBlocks() int

	// initiate performs all the required calculations and sends out all the
	// messages associated with the current state.
	initiate() error

	// receive is called each time a new message arrived. receive is expected to
	// be called for all broadcast channel messages, including the member's own
	// messages.
	receive(msg net.Message) error

	// nextState performs a state transition to the next state of the protocol.
	// If the current state is the last one, nextState returns `nil`.
	nextState() keyGenerationState

	// memberID returns the ID of member associated with the current state.
	memberID() gjkr.MemberID
}

func isMessageFromSelf(
	state keyGenerationState,
	message gjkr.ProtocolMessage,
) bool {
	if message.SenderID() == state.memberID() {
		return true
	}

	return false
}

func isSenderAccepted(
	filter gjkr.MessageFiltering,
	message gjkr.ProtocolMessage,
) bool {
	return filter.IsSenderAccepted(message.SenderID())
}

// initializationState is the starting state of key generation; it waits for
// activePeriod and then enters joinState. No messages are valid in this state.
type initializationState struct {
	channel net.BroadcastChannel
	member  *gjkr.LocalMember
}

func (is *initializationState) activeBlocks() int { return 3 }

func (is *initializationState) initiate() error {
	return nil
}

func (is *initializationState) receive(msg net.Message) error {
	return nil
}

func (is *initializationState) nextState() keyGenerationState {
	return &joinState{is.channel, is.member}
}

func (is *initializationState) memberID() gjkr.MemberID {
	return is.member.ID
}

// joinState is the state during which a member announces itself to the key
// generation broadcast channel to initiate the distributed protocol.
// `gjkr.JoinMessage`s are valid in this state.
type joinState struct {
	channel net.BroadcastChannel
	member  *gjkr.LocalMember
}

func (js *joinState) activeBlocks() int { return 3 }

func (js *joinState) initiate() error {
	return js.channel.Send(gjkr.NewJoinMessage(js.member.ID))
}

func (js *joinState) receive(msg net.Message) error {
	switch joinMsg := msg.Payload().(type) {
	case *gjkr.JoinMessage:
		js.member.AddToGroup(joinMsg.SenderID())
	}
	return nil
}

func (js *joinState) nextState() keyGenerationState {
	return &ephemeralKeyPairGenerationState{
		channel: js.channel,
		member:  js.member.InitializeEphemeralKeysGeneration(),
	}
}

func (js *joinState) memberID() gjkr.MemberID {
	return js.member.ID
}

// ephemeralKeyPairGenerationState is the state during which members broadcast
// public ephemeral keys generated for other members of the group.
// `gjkr.EphemeralPublicKeyMessage`s are valid in this state.
//
// State covers phase 1 of the protocol.
type ephemeralKeyPairGenerationState struct {
	channel net.BroadcastChannel
	member  *gjkr.EphemeralKeyPairGeneratingMember

	phaseMessages []*gjkr.EphemeralPublicKeyMessage
}

func (ekpgs *ephemeralKeyPairGenerationState) activeBlocks() int { return 3 }

func (ekpgs *ephemeralKeyPairGenerationState) initiate() error {
	message, err := ekpgs.member.GenerateEphemeralKeyPair()
	if err != nil {
		return err
	}

	if err := ekpgs.channel.Send(message); err != nil {
		return err
	}
	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *gjkr.EphemeralPublicKeyMessage:
		if !isMessageFromSelf(ekpgs, phaseMessage) &&
			isSenderAccepted(ekpgs.member, phaseMessage) {
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) nextState() keyGenerationState {
	return &symmetricKeyGenerationState{
		channel:               ekpgs.channel,
		member:                ekpgs.member.InitializeSymmetricKeyGeneration(),
		previousPhaseMessages: ekpgs.phaseMessages,
	}
}

func (ekpgs *ephemeralKeyPairGenerationState) memberID() gjkr.MemberID {
	return ekpgs.member.ID
}

// symmetricKeyGenerationState is the state during which members compute
// symmetric keys from the previously exchanged ephemeral public keys.
// No messages are valid in this state.
//
// State covers phase 2 of the protocol.
type symmetricKeyGenerationState struct {
	channel net.BroadcastChannel
	member  *gjkr.SymmetricKeyGeneratingMember

	previousPhaseMessages []*gjkr.EphemeralPublicKeyMessage
}

func (skgs *symmetricKeyGenerationState) activeBlocks() int { return 0 }

func (skgs *symmetricKeyGenerationState) initiate() error {
	skgs.member.MarkInactiveMembers(skgs.previousPhaseMessages)
	return skgs.member.GenerateSymmetricKeys(skgs.previousPhaseMessages)
}

func (skgs *symmetricKeyGenerationState) receive(msg net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) nextState() keyGenerationState {
	return &commitmentState{
		channel: skgs.channel,
		member:  skgs.member.InitializeCommitting(),
	}
}

func (skgs *symmetricKeyGenerationState) memberID() gjkr.MemberID {
	return skgs.member.ID
}

// commitmentState is the state during which members compute their individual
// shares and commitments to those shares. Two messages are valid in this state:
// - `gjkr.PeerSharesMessage`
// - `gjkr.MemberCommitmentsMessage`
//
// State covers phase 3 of the protocol.
type commitmentState struct {
	channel net.BroadcastChannel
	member  *gjkr.CommittingMember

	phaseSharesMessages      []*gjkr.PeerSharesMessage
	phaseCommitmentsMessages []*gjkr.MemberCommitmentsMessage
}

func (cs *commitmentState) activeBlocks() int { return 3 }

func (cs *commitmentState) initiate() error {
	sharesMsg, commitmentsMsg, err := cs.member.CalculateMembersSharesAndCommitments()
	if err != nil {
		return err
	}

	if err := cs.channel.Send(sharesMsg); err != nil {
		return err
	}

	if err := cs.channel.Send(commitmentsMsg); err != nil {
		return err
	}

	return nil
}

func (cs *commitmentState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *gjkr.PeerSharesMessage:
		if !isMessageFromSelf(cs, phaseMessage) &&
			isSenderAccepted(cs.member, phaseMessage) {
			cs.phaseSharesMessages = append(cs.phaseSharesMessages, phaseMessage)
		}

	case *gjkr.MemberCommitmentsMessage:
		if !isMessageFromSelf(cs, phaseMessage) {
			cs.phaseCommitmentsMessages = append(
				cs.phaseCommitmentsMessages,
				phaseMessage,
			)
		}
	}

	return nil
}

func (cs *commitmentState) nextState() keyGenerationState {
	return &commitmentsVerificationState{
		channel: cs.channel,
		member:  cs.member.InitializeCommitmentsVerification(),

		previousPhaseSharesMessages:      cs.phaseSharesMessages,
		previousPhaseCommitmentsMessages: cs.phaseCommitmentsMessages,
	}
}

func (cs *commitmentState) memberID() gjkr.MemberID {
	return cs.member.ID
}

// commitmentsVerificationState is the state during which members validate
// shares and commitments computed and published by other members in the
// previous phase. `gjkr.SecretShareAccusationMessage`s are valid in this state.
//
// State covers phase 4 of the protocol.
type commitmentsVerificationState struct {
	channel net.BroadcastChannel
	member  *gjkr.CommitmentsVerifyingMember

	previousPhaseSharesMessages      []*gjkr.PeerSharesMessage
	previousPhaseCommitmentsMessages []*gjkr.MemberCommitmentsMessage

	phaseAccusationsMessages []*gjkr.SecretSharesAccusationsMessage
}

func (cvs *commitmentsVerificationState) activeBlocks() int { return 3 }

func (cvs *commitmentsVerificationState) initiate() error {
	cvs.member.MarkInactiveMembers(
		cvs.previousPhaseSharesMessages,
		cvs.previousPhaseCommitmentsMessages,
	)
	accusationsMsg, err := cvs.member.VerifyReceivedSharesAndCommitmentsMessages(
		cvs.previousPhaseSharesMessages,
		cvs.previousPhaseCommitmentsMessages,
	)
	if err != nil {
		return err
	}

	if err := cvs.channel.Send(accusationsMsg); err != nil {
		return err
	}

	return nil
}

func (cvs *commitmentsVerificationState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *gjkr.SecretSharesAccusationsMessage:
		if !isMessageFromSelf(cvs, phaseMessage) &&
			isSenderAccepted(cvs.member, phaseMessage) {
			cvs.phaseAccusationsMessages = append(
				cvs.phaseAccusationsMessages,
				phaseMessage,
			)
		}
	}

	return nil
}

func (cvs *commitmentsVerificationState) nextState() keyGenerationState {
	return &sharesJustificationState{
		channel: cvs.channel,
		member:  cvs.member.InitializeSharesJustification(),

		previousPhaseAccusationsMessages: cvs.phaseAccusationsMessages,
	}
}

func (cvs *commitmentsVerificationState) memberID() gjkr.MemberID {
	return cvs.member.ID
}

// sharesJustificationState is the state during which members resolve
// accusations published by other group members in the previous state.
// No messages are valid in this state.
//
// State covers phase 5 of the protocol.
type sharesJustificationState struct {
	channel net.BroadcastChannel
	member  *gjkr.SharesJustifyingMember

	previousPhaseAccusationsMessages []*gjkr.SecretSharesAccusationsMessage
}

func (sjs *sharesJustificationState) activeBlocks() int { return 0 }

func (sjs *sharesJustificationState) initiate() error {
	disqualifiedMembers, err := sjs.member.ResolveSecretSharesAccusationsMessages(
		sjs.previousPhaseAccusationsMessages,
	)
	if err != nil {
		return err
	}

	// TODO: Handle member disqualification
	fmt.Printf("disqualified members = %v\n", disqualifiedMembers)

	return nil
}

func (sjs *sharesJustificationState) receive(msg net.Message) error {
	return nil
}

func (sjs *sharesJustificationState) nextState() keyGenerationState {
	return &qualificationState{
		channel: sjs.channel,
		member:  sjs.member.InitializeQualified(),
	}
}

func (sjs *sharesJustificationState) memberID() gjkr.MemberID {
	return sjs.member.ID
}

// qualificationState is the state during which group members combine all valid
// secret shares published by other group members in the previous states.
// No messages are valid in this state.
//
// State covers phase 6 of the protocol.
type qualificationState struct {
	channel net.BroadcastChannel
	member  *gjkr.QualifiedMember
}

func (qs *qualificationState) activeBlocks() int { return 0 }

func (qs *qualificationState) initiate() error {
	qs.member.CombineMemberShares()
	return nil
}

func (qs *qualificationState) receive(msg net.Message) error {
	return nil
}

func (qs *qualificationState) nextState() keyGenerationState {
	return &pointsShareState{
		channel: qs.channel,
		member:  qs.member.InitializeSharing(),
	}
}

func (qs *qualificationState) memberID() gjkr.MemberID {
	return qs.member.ID
}

// pointsShareState is the state during which group members calculate and
// publish their public key share points.
// `gjkr.MemberPublicKeySharePointsMessage`s are valid in this state.
//
// State covers phase 7 of the protocol.
type pointsShareState struct {
	channel net.BroadcastChannel
	member  *gjkr.SharingMember // TODO: SharingMember should be renamed to PointsSharingMember

	phaseMessages []*gjkr.MemberPublicKeySharePointsMessage
}

func (pss *pointsShareState) activeBlocks() int { return 3 }

func (pss *pointsShareState) initiate() error {
	message := pss.member.CalculatePublicKeySharePoints()
	if err := pss.channel.Send(message); err != nil {
		return err
	}

	return nil
}

func (pss *pointsShareState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *gjkr.MemberPublicKeySharePointsMessage:
		if !isMessageFromSelf(pss, phaseMessage) &&
			isSenderAccepted(pss.member, phaseMessage) {
			pss.phaseMessages = append(pss.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (pss *pointsShareState) nextState() keyGenerationState {
	return &pointsValidationState{
		channel: pss.channel,
		member:  pss.member,

		previousPhaseMessages: pss.phaseMessages,
	}
}

func (pss *pointsShareState) memberID() gjkr.MemberID {
	return pss.member.ID
}

// pointsValidationState is the state during which group members validate
// public key share points published by other group members in the previous
// state. `gjkr.PointsAccusationsMessage`s are valid in this state.
//
// State covers phase 8 of the protocol.
type pointsValidationState struct {
	channel net.BroadcastChannel
	member  *gjkr.SharingMember // TODO: split validation logic into PointsValidatingMember

	previousPhaseMessages []*gjkr.MemberPublicKeySharePointsMessage

	phaseMessages []*gjkr.PointsAccusationsMessage
}

func (pvs *pointsValidationState) activeBlocks() int { return 3 }

func (pvs *pointsValidationState) initiate() error {
	pvs.member.MarkInactiveMembers(pvs.previousPhaseMessages)
	accusationMsg, err := pvs.member.VerifyPublicKeySharePoints(
		pvs.previousPhaseMessages,
	)
	if err != nil {
		return err
	}

	if err := pvs.channel.Send(accusationMsg); err != nil {
		return err
	}

	return nil
}

func (pvs *pointsValidationState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *gjkr.PointsAccusationsMessage:
		if !isMessageFromSelf(pvs, phaseMessage) &&
			isSenderAccepted(pvs.member, phaseMessage) {
			pvs.phaseMessages = append(pvs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (pvs *pointsValidationState) nextState() keyGenerationState {
	return &pointsJustificationState{
		channel: pvs.channel,
		member:  pvs.member.InitializePointsJustification(),

		previousPhaseMessages: pvs.phaseMessages,
	}
}

func (pvs *pointsValidationState) memberID() gjkr.MemberID {
	return pvs.member.ID
}

// pointsJustificationState is the state during which group members resolve
// accusations published by other group members in the previous state.
// No messages are valid in this state.
//
// State covers phase 9 of the protocol.
type pointsJustificationState struct {
	channel net.BroadcastChannel
	member  *gjkr.PointsJustifyingMember

	previousPhaseMessages []*gjkr.PointsAccusationsMessage
}

func (pjs *pointsJustificationState) activeBlocks() int { return 0 }

func (pjs *pointsJustificationState) initiate() error {
	disqualifiedMembers, err := pjs.member.ResolvePublicKeySharePointsAccusationsMessages(
		pjs.previousPhaseMessages,
	)
	if err != nil {
		return err
	}

	// TODO: Handle member disqualification
	fmt.Printf("disqualified members = %v\n", disqualifiedMembers)

	return nil
}

func (pjs *pointsJustificationState) receive(msg net.Message) error {
	return nil
}

func (pjs *pointsJustificationState) nextState() keyGenerationState {
	return &keyRevealState{
		channel: pjs.channel,
		member:  pjs.member.InitializeRevealing(),
	}
}

func (pjs *pointsJustificationState) memberID() gjkr.MemberID {
	return pjs.member.ID
}

// keyRevealState is the state during which group members reveal ephemeral
// private keys used to create an ephemeral symmetric keys with disqualified
// members who share a group private key.
//
// State covers phase 10 of the protocol.
type keyRevealState struct {
	channel net.BroadcastChannel
	member  *gjkr.RevealingMember // TODO: Rename to KeyRevealingMember

	phaseMessages []*gjkr.DisqualifiedEphemeralKeysMessage
}

func (rs *keyRevealState) activeBlocks() int { return 1 }

func (rs *keyRevealState) initiate() error {
	revealMsg, err := rs.member.RevealDisqualifiedMembersKeys()
	if err != nil {
		return err
	}

	if err := rs.channel.Send(revealMsg); err != nil {
		return err
	}

	return nil
}

func (rs *keyRevealState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *gjkr.DisqualifiedEphemeralKeysMessage:
		if !isMessageFromSelf(rs, phaseMessage) &&
			isSenderAccepted(rs.member, phaseMessage) {
			rs.phaseMessages = append(rs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (rs *keyRevealState) nextState() keyGenerationState {
	return &reconstructionState{
		channel:               rs.channel,
		member:                rs.member.InitializeReconstruction(),
		previousPhaseMessages: rs.phaseMessages,
	}
}

func (rs *keyRevealState) memberID() gjkr.MemberID {
	return rs.member.ID
}

// reconstructionState is the state during which group members reconstruct
// individual keys of members disqualified in previous states. No messages are
// valid in this state.
//
// State covers phase 11 of the protocol.
type reconstructionState struct {
	channel net.BroadcastChannel
	member  *gjkr.ReconstructingMember

	previousPhaseMessages []*gjkr.DisqualifiedEphemeralKeysMessage
}

func (rs *reconstructionState) activeBlocks() int { return 0 }

func (rs *reconstructionState) initiate() error {
	rs.member.MarkInactiveMembers(rs.previousPhaseMessages)
	if err := rs.member.ReconstructDisqualifiedIndividualKeys(
		rs.previousPhaseMessages,
	); err != nil {
		return err
	}

	return nil
}

func (rs *reconstructionState) receive(msg net.Message) error {
	return nil
}

func (rs *reconstructionState) nextState() keyGenerationState {
	return &combinationState{
		channel: rs.channel,
		member:  rs.member.InitializeCombining(),
	}
}

func (rs *reconstructionState) memberID() gjkr.MemberID {
	return rs.member.ID
}

// combinationState is the state during which group members combine together all
// qualified key shares to form a group public key. No messages are valid in
// this state.
//
// State covers phase 12 of the protocol.
type combinationState struct {
	channel net.BroadcastChannel
	member  *gjkr.CombiningMember
}

func (cs *combinationState) activeBlocks() int { return 0 }

func (cs *combinationState) initiate() error {
	cs.member.CombineGroupPublicKey()
	return nil
}

func (cs *combinationState) receive(msg net.Message) error {
	return nil
}

func (cs *combinationState) nextState() keyGenerationState {
	return &finalizationState{
		channel: cs.channel,
		member:  cs.member.InitializeFinalization(),
	}
}

func (cs *combinationState) memberID() gjkr.MemberID {
	return cs.member.ID
}

// finalizationState is the last state of GJKR DKG protocol - in this state,
// distributed key generation is completed. No messages are valid in this state.
//
// State prepares a result to publish in phase 13 of the protocol but it does
// not execute that phase.
type finalizationState struct {
	channel net.BroadcastChannel
	member  *gjkr.FinalizingMember
}

func (fs *finalizationState) activeBlocks() int { return 0 }

func (fs *finalizationState) initiate() error {
	return nil
}

func (fs *finalizationState) receive(msg net.Message) error {
	return nil
}

func (fs *finalizationState) nextState() keyGenerationState {
	return nil
}

func (fs *finalizationState) memberID() gjkr.MemberID {
	return fs.member.ID
}

func (fs *finalizationState) result() *gjkr.Result {
	return fs.member.Result()
}

func (fs *finalizationState) thresholdSigner() *ThresholdSigner {
	return &ThresholdSigner{
		memberID:             fs.member.ID,
		groupPublicKey:       fs.member.GroupPublicKey(),
		groupPrivateKeyShare: fs.member.GroupPrivateKeyShare(),
	}
}
