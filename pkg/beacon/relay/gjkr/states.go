package gjkr

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/states"
	"github.com/keep-network/keep-core/pkg/net"
)

type keyGenerationState states.DKGState

func isMessageFromSelf(
	state keyGenerationState,
	message ProtocolMessage,
) bool {
	if message.SenderID() == MemberID(state.MemberID()) {
		return true
	}

	return false
}

func isSenderAccepted(
	filter MessageFiltering,
	message ProtocolMessage,
) bool {
	return filter.IsSenderAccepted(message.SenderID())
}

// Init initializes a given broadcast channel to be able to perform distributed
// key generation interactions.
func Init(channel net.BroadcastChannel) func(channel net.BroadcastChannel) {
	return func(channel net.BroadcastChannel) {
		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &JoinMessage{}
		})
		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &EphemeralPublicKeyMessage{}
		})
		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &MemberCommitmentsMessage{}
		})
		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &PeerSharesMessage{}
		})
		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &SecretSharesAccusationsMessage{}
		})
		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &MemberPublicKeySharePointsMessage{}
		})
		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &PointsAccusationsMessage{}
		})
		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &DisqualifiedEphemeralKeysMessage{}
		})
	}
}

func InitializationState(
	channel net.BroadcastChannel,
	member *LocalMember,
) *initializationState {
	return &initializationState{channel, member}
}

// initializationState is the starting state of key generation; it waits for
// activePeriod and then enters joinState. No messages are valid in this state.
type initializationState struct {
	channel net.BroadcastChannel
	member  *LocalMember
}

func (is *initializationState) ActiveBlocks() int { return 3 }

func (is *initializationState) Initiate() error {
	return nil
}

func (is *initializationState) Receive(msg net.Message) error {
	return nil
}

func (is *initializationState) NextState() states.DKGState {
	return &joinState{is.channel, is.member}
}

func (is *initializationState) MemberID() states.MemberIndex {
	return states.MemberIndex(is.member.ID)
}

func (is *initializationState) IsFinalState() bool {
	return false
}

// joinState is the state during which a member announces itself to the key
// generation broadcast channel to initiate the distributed protocol.
// `JoinMessage`s are valid in this state.
type joinState struct {
	channel net.BroadcastChannel
	member  *LocalMember
}

func (js *joinState) ActiveBlocks() int { return 3 }

func (js *joinState) Initiate() error {
	return js.channel.Send(NewJoinMessage(js.member.ID))
}

func (js *joinState) Receive(msg net.Message) error {
	switch joinMsg := msg.Payload().(type) {
	case *JoinMessage:
		js.member.AddToGroup(joinMsg.SenderID())
	}
	return nil
}

func (js *joinState) NextState() states.DKGState {
	return &ephemeralKeyPairGenerationState{
		channel: js.channel,
		member:  js.member.InitializeEphemeralKeysGeneration(),
	}
}

func (js *joinState) MemberID() states.MemberIndex {
	return states.MemberIndex(js.member.ID)
}

func (js *joinState) IsFinalState() bool {
	return false
}

// ephemeralKeyPairGenerationState is the state during which members broadcast
// public ephemeral keys generated for other members of the group.
// `EphemeralPublicKeyMessage`s are valid in this state.
//
// State covers phase 1 of the protocol.
type ephemeralKeyPairGenerationState struct {
	channel net.BroadcastChannel
	member  *EphemeralKeyPairGeneratingMember

	phaseMessages []*EphemeralPublicKeyMessage
}

func (ekpgs *ephemeralKeyPairGenerationState) ActiveBlocks() int { return 3 }

func (ekpgs *ephemeralKeyPairGenerationState) Initiate() error {
	message, err := ekpgs.member.GenerateEphemeralKeyPair()
	if err != nil {
		return err
	}

	if err := ekpgs.channel.Send(message); err != nil {
		return err
	}
	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *EphemeralPublicKeyMessage:
		if !isMessageFromSelf(ekpgs, phaseMessage) &&
			isSenderAccepted(ekpgs.member, phaseMessage) {
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) NextState() states.DKGState {
	return &symmetricKeyGenerationState{
		channel:               ekpgs.channel,
		member:                ekpgs.member.InitializeSymmetricKeyGeneration(),
		previousPhaseMessages: ekpgs.phaseMessages,
	}
}

func (ekpgs *ephemeralKeyPairGenerationState) MemberID() states.MemberIndex {
	return states.MemberIndex(ekpgs.member.ID)
}

func (ekpgs *ephemeralKeyPairGenerationState) IsFinalState() bool {
	return false
}

// symmetricKeyGenerationState is the state during which members compute
// symmetric keys from the previously exchanged ephemeral public keys.
// No messages are valid in this state.
//
// State covers phase 2 of the protocol.
type symmetricKeyGenerationState struct {
	channel net.BroadcastChannel
	member  *SymmetricKeyGeneratingMember

	previousPhaseMessages []*EphemeralPublicKeyMessage
}

func (skgs *symmetricKeyGenerationState) ActiveBlocks() int { return 0 }

func (skgs *symmetricKeyGenerationState) Initiate() error {
	skgs.member.MarkInactiveMembers(skgs.previousPhaseMessages)
	return skgs.member.GenerateSymmetricKeys(skgs.previousPhaseMessages)
}

func (skgs *symmetricKeyGenerationState) Receive(msg net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) NextState() states.DKGState {
	return &commitmentState{
		channel: skgs.channel,
		member:  skgs.member.InitializeCommitting(),
	}
}

func (skgs *symmetricKeyGenerationState) MemberID() states.MemberIndex {
	return states.MemberIndex(skgs.member.ID)
}

func (skgs *symmetricKeyGenerationState) IsFinalState() bool {
	return false
}

// commitmentState is the state during which members compute their individual
// shares and commitments to those shares. Two messages are valid in this state:
// - `PeerSharesMessage`
// - `MemberCommitmentsMessage`
//
// State covers phase 3 of the protocol.
type commitmentState struct {
	channel net.BroadcastChannel
	member  *CommittingMember

	phaseSharesMessages      []*PeerSharesMessage
	phaseCommitmentsMessages []*MemberCommitmentsMessage
}

func (cs *commitmentState) ActiveBlocks() int { return 3 }

func (cs *commitmentState) Initiate() error {
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

func (cs *commitmentState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *PeerSharesMessage:
		if !isMessageFromSelf(cs, phaseMessage) &&
			isSenderAccepted(cs.member, phaseMessage) {
			cs.phaseSharesMessages = append(cs.phaseSharesMessages, phaseMessage)
		}

	case *MemberCommitmentsMessage:
		if !isMessageFromSelf(cs, phaseMessage) {
			cs.phaseCommitmentsMessages = append(
				cs.phaseCommitmentsMessages,
				phaseMessage,
			)
		}
	}

	return nil
}

func (cs *commitmentState) NextState() states.DKGState {
	return &commitmentsVerificationState{
		channel: cs.channel,
		member:  cs.member.InitializeCommitmentsVerification(),

		previousPhaseSharesMessages:      cs.phaseSharesMessages,
		previousPhaseCommitmentsMessages: cs.phaseCommitmentsMessages,
	}
}

func (cs *commitmentState) MemberID() states.MemberIndex {
	return states.MemberIndex(cs.member.ID)
}

func (cs *commitmentState) IsFinalState() bool {
	return false
}

// commitmentsVerificationState is the state during which members validate
// shares and commitments computed and published by other members in the
// previous phase. `SecretShareAccusationMessage`s are valid in this state.
//
// State covers phase 4 of the protocol.
type commitmentsVerificationState struct {
	channel net.BroadcastChannel
	member  *CommitmentsVerifyingMember

	previousPhaseSharesMessages      []*PeerSharesMessage
	previousPhaseCommitmentsMessages []*MemberCommitmentsMessage

	phaseAccusationsMessages []*SecretSharesAccusationsMessage
}

func (cvs *commitmentsVerificationState) ActiveBlocks() int { return 3 }

func (cvs *commitmentsVerificationState) Initiate() error {
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

func (cvs *commitmentsVerificationState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *SecretSharesAccusationsMessage:
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

func (cvs *commitmentsVerificationState) NextState() states.DKGState {
	return &sharesJustificationState{
		channel: cvs.channel,
		member:  cvs.member.InitializeSharesJustification(),

		previousPhaseAccusationsMessages: cvs.phaseAccusationsMessages,
	}
}

func (cvs *commitmentsVerificationState) MemberID() states.MemberIndex {
	return states.MemberIndex(cvs.member.ID)
}

func (cvs *commitmentsVerificationState) IsFinalState() bool {
	return false
}

// sharesJustificationState is the state during which members resolve
// accusations published by other group members in the previous state.
// No messages are valid in this state.
//
// State covers phase 5 of the protocol.
type sharesJustificationState struct {
	channel net.BroadcastChannel
	member  *SharesJustifyingMember

	previousPhaseAccusationsMessages []*SecretSharesAccusationsMessage
}

func (sjs *sharesJustificationState) ActiveBlocks() int { return 0 }

func (sjs *sharesJustificationState) Initiate() error {
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

func (sjs *sharesJustificationState) Receive(msg net.Message) error {
	return nil
}

func (sjs *sharesJustificationState) NextState() states.DKGState {
	return &qualificationState{
		channel: sjs.channel,
		member:  sjs.member.InitializeQualified(),
	}
}

func (sjs *sharesJustificationState) MemberID() states.MemberIndex {
	return states.MemberIndex(sjs.member.ID)
}

func (sjs *sharesJustificationState) IsFinalState() bool {
	return false
}

// qualificationState is the state during which group members combine all valid
// secret shares published by other group members in the previous states.
// No messages are valid in this state.
//
// State covers phase 6 of the protocol.
type qualificationState struct {
	channel net.BroadcastChannel
	member  *QualifiedMember
}

func (qs *qualificationState) ActiveBlocks() int { return 0 }

func (qs *qualificationState) Initiate() error {
	qs.member.CombineMemberShares()
	return nil
}

func (qs *qualificationState) Receive(msg net.Message) error {
	return nil
}

func (qs *qualificationState) NextState() states.DKGState {
	return &pointsShareState{
		channel: qs.channel,
		member:  qs.member.InitializeSharing(),
	}
}

func (qs *qualificationState) MemberID() states.MemberIndex {
	return states.MemberIndex(qs.member.ID)
}

func (qs *qualificationState) IsFinalState() bool {
	return false
}

// pointsShareState is the state during which group members calculate and
// publish their public key share points.
// `MemberPublicKeySharePointsMessage`s are valid in this state.
//
// State covers phase 7 of the protocol.
type pointsShareState struct {
	channel net.BroadcastChannel
	member  *SharingMember // TODO: SharingMember should be renamed to PointsSharingMember

	phaseMessages []*MemberPublicKeySharePointsMessage
}

func (pss *pointsShareState) ActiveBlocks() int { return 3 }

func (pss *pointsShareState) Initiate() error {
	message := pss.member.CalculatePublicKeySharePoints()
	if err := pss.channel.Send(message); err != nil {
		return err
	}

	return nil
}

func (pss *pointsShareState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *MemberPublicKeySharePointsMessage:
		if !isMessageFromSelf(pss, phaseMessage) &&
			isSenderAccepted(pss.member, phaseMessage) {
			pss.phaseMessages = append(pss.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (pss *pointsShareState) NextState() states.DKGState {
	return &pointsValidationState{
		channel: pss.channel,
		member:  pss.member,

		previousPhaseMessages: pss.phaseMessages,
	}
}

func (pss *pointsShareState) MemberID() states.MemberIndex {
	return states.MemberIndex(pss.member.ID)
}

func (pss *pointsShareState) IsFinalState() bool {
	return false
}

// pointsValidationState is the state during which group members validate
// public key share points published by other group members in the previous
// state. `PointsAccusationsMessage`s are valid in this state.
//
// State covers phase 8 of the protocol.
type pointsValidationState struct {
	channel net.BroadcastChannel
	member  *SharingMember // TODO: split validation logic into PointsValidatingMember

	previousPhaseMessages []*MemberPublicKeySharePointsMessage

	phaseMessages []*PointsAccusationsMessage
}

func (pvs *pointsValidationState) ActiveBlocks() int { return 3 }

func (pvs *pointsValidationState) Initiate() error {
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

func (pvs *pointsValidationState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *PointsAccusationsMessage:
		if !isMessageFromSelf(pvs, phaseMessage) &&
			isSenderAccepted(pvs.member, phaseMessage) {
			pvs.phaseMessages = append(pvs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (pvs *pointsValidationState) NextState() states.DKGState {
	return &pointsJustificationState{
		channel: pvs.channel,
		member:  pvs.member.InitializePointsJustification(),

		previousPhaseMessages: pvs.phaseMessages,
	}
}

func (pvs *pointsValidationState) MemberID() states.MemberIndex {
	return states.MemberIndex(pvs.member.ID)
}

func (pvs *pointsValidationState) IsFinalState() bool {
	return false
}

// pointsJustificationState is the state during which group members resolve
// accusations published by other group members in the previous state.
// No messages are valid in this state.
//
// State covers phase 9 of the protocol.
type pointsJustificationState struct {
	channel net.BroadcastChannel
	member  *PointsJustifyingMember

	previousPhaseMessages []*PointsAccusationsMessage
}

func (pjs *pointsJustificationState) ActiveBlocks() int { return 0 }

func (pjs *pointsJustificationState) Initiate() error {
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

func (pjs *pointsJustificationState) Receive(msg net.Message) error {
	return nil
}

func (pjs *pointsJustificationState) NextState() states.DKGState {
	return &keyRevealState{
		channel: pjs.channel,
		member:  pjs.member.InitializeRevealing(),
	}
}

func (pjs *pointsJustificationState) MemberID() states.MemberIndex {
	return states.MemberIndex(pjs.member.ID)
}

func (pjs *pointsJustificationState) IsFinalState() bool {
	return false
}

// keyRevealState is the state during which group members reveal ephemeral
// private keys used to create an ephemeral symmetric keys with disqualified
// members who share a group private key.
//
// State covers phase 10 of the protocol.
type keyRevealState struct {
	channel net.BroadcastChannel
	member  *RevealingMember // TODO: Rename to KeyRevealingMember

	phaseMessages []*DisqualifiedEphemeralKeysMessage
}

func (rs *keyRevealState) ActiveBlocks() int { return 1 }

func (rs *keyRevealState) Initiate() error {
	revealMsg, err := rs.member.RevealDisqualifiedMembersKeys()
	if err != nil {
		return err
	}

	if err := rs.channel.Send(revealMsg); err != nil {
		return err
	}

	return nil
}

func (rs *keyRevealState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *DisqualifiedEphemeralKeysMessage:
		if !isMessageFromSelf(rs, phaseMessage) &&
			isSenderAccepted(rs.member, phaseMessage) {
			rs.phaseMessages = append(rs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (rs *keyRevealState) NextState() states.DKGState {
	return &reconstructionState{
		channel:               rs.channel,
		member:                rs.member.InitializeReconstruction(),
		previousPhaseMessages: rs.phaseMessages,
	}
}

func (rs *keyRevealState) MemberID() states.MemberIndex {
	return states.MemberIndex(rs.member.ID)
}

func (rs *keyRevealState) IsFinalState() bool {
	return false
}

// reconstructionState is the state during which group members reconstruct
// individual keys of members disqualified in previous states. No messages are
// valid in this state.
//
// State covers phase 11 of the protocol.
type reconstructionState struct {
	channel net.BroadcastChannel
	member  *ReconstructingMember

	previousPhaseMessages []*DisqualifiedEphemeralKeysMessage
}

func (rs *reconstructionState) ActiveBlocks() int { return 0 }

func (rs *reconstructionState) Initiate() error {
	rs.member.MarkInactiveMembers(rs.previousPhaseMessages)
	if err := rs.member.ReconstructDisqualifiedIndividualKeys(
		rs.previousPhaseMessages,
	); err != nil {
		return err
	}

	return nil
}

func (rs *reconstructionState) Receive(msg net.Message) error {
	return nil
}

func (rs *reconstructionState) NextState() states.DKGState {
	return &combinationState{
		channel: rs.channel,
		member:  rs.member.InitializeCombining(),
	}
}

func (rs *reconstructionState) MemberID() states.MemberIndex {
	return states.MemberIndex(rs.member.ID)
}

func (rs *reconstructionState) IsFinalState() bool {
	return false
}

// combinationState is the state during which group members combine together all
// qualified key shares to form a group public key. No messages are valid in
// this state.
//
// State covers phase 12 of the protocol.
type combinationState struct {
	channel net.BroadcastChannel
	member  *CombiningMember
}

func (cs *combinationState) ActiveBlocks() int { return 0 }

func (cs *combinationState) Initiate() error {
	cs.member.CombineGroupPublicKey()
	return nil
}

func (cs *combinationState) Receive(msg net.Message) error {
	return nil
}

func (cs *combinationState) NextState() states.DKGState {
	return &FinalizationState{
		channel: cs.channel,
		member:  cs.member.InitializeFinalization(),
	}
}

func (cs *combinationState) MemberID() states.MemberIndex {
	return states.MemberIndex(cs.member.ID)
}

func (cs *combinationState) IsFinalState() bool {
	return false
}

// finalizationState is the last state of GJKR DKG protocol - in this state,
// distributed key generation is completed. No messages are valid in this state.
//
// State prepares a result to publish in phase 13 of the protocol but it does
// not execute that phase.
type FinalizationState struct {
	channel net.BroadcastChannel
	member  *FinalizingMember
}

func (fs *FinalizationState) ActiveBlocks() int { return 0 }

func (fs *FinalizationState) Initiate() error {
	return nil
}

func (fs *FinalizationState) Receive(msg net.Message) error {
	return nil
}

func (fs *FinalizationState) NextState() states.DKGState {
	return nil
}

func (fs *FinalizationState) MemberID() states.MemberIndex {
	return states.MemberIndex(fs.member.ID)
}

func (fs *FinalizationState) IsFinalState() bool {
	return true
}

func (fs *FinalizationState) Result() *Result {
	return fs.member.Result()
}

func (fs *FinalizationState) GroupPrivateKeyShare() *big.Int {
	return fs.member.GroupPrivateKeyShare()
}
