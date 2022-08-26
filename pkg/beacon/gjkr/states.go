package gjkr

import (
	"context"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

type keyGenerationState = state.State

const (
	silentStateDelayBlocks  = 0
	silentStateActiveBlocks = 0

	ephemeralKeyPairStateDelayBlocks  = 1
	ephemeralKeyPairStateActiveBlocks = 5

	commitmentStateDelayBlocks  = 1
	commitmentStateActiveBlocks = 5

	commitmentVerificationStateDelayBlocks  = 1
	commitmentVerificationStateActiveBlocks = 10

	pointsShareStateDelayBlocks  = 1
	pointsShareStateActiveBlocks = 5

	pointsValidationStateDelayBlocks  = 1
	pointsValidationStateActiveBlocks = 10

	keyRevealStateDelayBlocks  = 1
	keyRevealStateActiveBlocks = 5

	combinationStateDelayBlocks  = 0
	combinationStateActiveBlocks = 20
)

// ProtocolBlocks returns the total number of blocks it takes to execute
// all the required work defined by the GJKR protocol.
func ProtocolBlocks() uint64 {
	return ephemeralKeyPairStateDelayBlocks +
		ephemeralKeyPairStateActiveBlocks +
		commitmentStateDelayBlocks +
		commitmentStateActiveBlocks +
		commitmentVerificationStateDelayBlocks +
		commitmentVerificationStateActiveBlocks +
		pointsShareStateDelayBlocks +
		pointsShareStateActiveBlocks +
		pointsValidationStateDelayBlocks +
		pointsValidationStateActiveBlocks +
		keyRevealStateDelayBlocks +
		keyRevealStateActiveBlocks +
		combinationStateDelayBlocks +
		combinationStateActiveBlocks
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

func (ekpgs *ephemeralKeyPairGenerationState) DelayBlocks() uint64 {
	return ephemeralKeyPairStateDelayBlocks
}

func (ekpgs *ephemeralKeyPairGenerationState) ActiveBlocks() uint64 {
	return ephemeralKeyPairStateActiveBlocks
}

func (ekpgs *ephemeralKeyPairGenerationState) Initiate(ctx context.Context) error {
	message, err := ekpgs.member.GenerateEphemeralKeyPair()
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
	case *EphemeralPublicKeyMessage:
		if ekpgs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) {
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Next() (state.State, error) {
	return &symmetricKeyGenerationState{
		channel:               ekpgs.channel,
		member:                ekpgs.member.InitializeSymmetricKeyGeneration(),
		previousPhaseMessages: ekpgs.phaseMessages,
	}, nil
}

func (ekpgs *ephemeralKeyPairGenerationState) MemberIndex() group.MemberIndex {
	return ekpgs.member.ID
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

func (skgs *symmetricKeyGenerationState) DelayBlocks() uint64 {
	return silentStateDelayBlocks
}

func (skgs *symmetricKeyGenerationState) ActiveBlocks() uint64 {
	return silentStateActiveBlocks
}

func (skgs *symmetricKeyGenerationState) Initiate(ctx context.Context) error {
	skgs.member.MarkInactiveMembers(skgs.previousPhaseMessages)
	return skgs.member.GenerateSymmetricKeys(skgs.previousPhaseMessages)
}

func (skgs *symmetricKeyGenerationState) Receive(msg net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) Next() (state.State, error) {
	return &commitmentState{
		channel: skgs.channel,
		member:  skgs.member.InitializeCommitting(),
	}, nil
}

func (skgs *symmetricKeyGenerationState) MemberIndex() group.MemberIndex {
	return skgs.member.ID
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

func (cs *commitmentState) DelayBlocks() uint64 {
	return commitmentStateDelayBlocks
}

func (cs *commitmentState) ActiveBlocks() uint64 {
	return commitmentStateActiveBlocks
}

func (cs *commitmentState) Initiate(ctx context.Context) error {
	sharesMsg, commitmentsMsg, err := cs.member.CalculateMembersSharesAndCommitments()
	if err != nil {
		return err
	}

	if err := cs.channel.Send(ctx, sharesMsg); err != nil {
		return err
	}

	if err := cs.channel.Send(ctx, commitmentsMsg); err != nil {
		return err
	}

	return nil
}

func (cs *commitmentState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *PeerSharesMessage:
		if cs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) {
			cs.phaseSharesMessages = append(cs.phaseSharesMessages, phaseMessage)
		}

	case *MemberCommitmentsMessage:
		if cs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) {
			cs.phaseCommitmentsMessages = append(
				cs.phaseCommitmentsMessages,
				phaseMessage,
			)
		}
	}

	return nil
}

func (cs *commitmentState) Next() (state.State, error) {
	return &commitmentsVerificationState{
		channel: cs.channel,
		member:  cs.member.InitializeCommitmentsVerification(),

		previousPhaseSharesMessages:      cs.phaseSharesMessages,
		previousPhaseCommitmentsMessages: cs.phaseCommitmentsMessages,
	}, nil
}

func (cs *commitmentState) MemberIndex() group.MemberIndex {
	return cs.member.ID
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

func (cvs *commitmentsVerificationState) DelayBlocks() uint64 {
	return commitmentVerificationStateDelayBlocks
}

func (cvs *commitmentsVerificationState) ActiveBlocks() uint64 {
	return commitmentVerificationStateActiveBlocks
}

func (cvs *commitmentsVerificationState) Initiate(ctx context.Context) error {
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

	if err := cvs.channel.Send(ctx, accusationsMsg); err != nil {
		return err
	}

	return nil
}

func (cvs *commitmentsVerificationState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *SecretSharesAccusationsMessage:
		if cvs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) {
			cvs.phaseAccusationsMessages = append(
				cvs.phaseAccusationsMessages,
				phaseMessage,
			)
		}
	}

	return nil
}

func (cvs *commitmentsVerificationState) Next() (state.State, error) {
	return &sharesJustificationState{
		channel: cvs.channel,
		member:  cvs.member.InitializeSharesJustification(),

		previousPhaseAccusationsMessages: cvs.phaseAccusationsMessages,
	}, nil
}

func (cvs *commitmentsVerificationState) MemberIndex() group.MemberIndex {
	return cvs.member.ID
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

func (sjs *sharesJustificationState) DelayBlocks() uint64 {
	return silentStateDelayBlocks
}

func (sjs *sharesJustificationState) ActiveBlocks() uint64 {
	return silentStateActiveBlocks
}

func (sjs *sharesJustificationState) Initiate(ctx context.Context) error {
	sjs.member.MarkInactiveMembers(sjs.previousPhaseAccusationsMessages)

	err := sjs.member.ResolveSecretSharesAccusationsMessages(
		sjs.previousPhaseAccusationsMessages,
	)
	if err != nil {
		return err
	}

	return nil
}

func (sjs *sharesJustificationState) Receive(msg net.Message) error {
	return nil
}

func (sjs *sharesJustificationState) Next() (state.State, error) {
	return &qualificationState{
		channel: sjs.channel,
		member:  sjs.member.InitializeQualified(),
	}, nil
}

func (sjs *sharesJustificationState) MemberIndex() group.MemberIndex {
	return sjs.member.ID
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

func (qs *qualificationState) DelayBlocks() uint64 {
	return silentStateDelayBlocks
}

func (qs *qualificationState) ActiveBlocks() uint64 {
	return silentStateActiveBlocks
}

func (qs *qualificationState) Initiate(ctx context.Context) error {
	qs.member.CombineMemberShares()
	return nil
}

func (qs *qualificationState) Receive(msg net.Message) error {
	return nil
}

func (qs *qualificationState) Next() (state.State, error) {
	return &pointsShareState{
		channel: qs.channel,
		member:  qs.member.InitializeSharing(),
	}, nil
}

func (qs *qualificationState) MemberIndex() group.MemberIndex {
	return qs.member.ID
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

func (pss *pointsShareState) DelayBlocks() uint64 {
	return pointsShareStateDelayBlocks
}

func (pss *pointsShareState) ActiveBlocks() uint64 {
	return pointsShareStateActiveBlocks
}

func (pss *pointsShareState) Initiate(ctx context.Context) error {
	message := pss.member.CalculatePublicKeySharePoints()
	if err := pss.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (pss *pointsShareState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *MemberPublicKeySharePointsMessage:
		if pss.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) {
			pss.phaseMessages = append(pss.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (pss *pointsShareState) Next() (state.State, error) {
	return &pointsValidationState{
		channel: pss.channel,
		member:  pss.member,

		previousPhaseMessages: pss.phaseMessages,
	}, nil
}

func (pss *pointsShareState) MemberIndex() group.MemberIndex {
	return pss.member.ID
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

func (pvs *pointsValidationState) DelayBlocks() uint64 {
	return pointsValidationStateDelayBlocks
}

func (pvs *pointsValidationState) ActiveBlocks() uint64 {
	return pointsValidationStateActiveBlocks
}

func (pvs *pointsValidationState) Initiate(ctx context.Context) error {
	pvs.member.MarkInactiveMembers(pvs.previousPhaseMessages)
	accusationMsg, err := pvs.member.VerifyPublicKeySharePoints(
		pvs.previousPhaseMessages,
	)
	if err != nil {
		return err
	}

	if err := pvs.channel.Send(ctx, accusationMsg); err != nil {
		return err
	}

	return nil
}

func (pvs *pointsValidationState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *PointsAccusationsMessage:
		if pvs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) {
			pvs.phaseMessages = append(pvs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (pvs *pointsValidationState) Next() (state.State, error) {
	return &pointsJustificationState{
		channel: pvs.channel,
		member:  pvs.member.InitializePointsJustification(),

		previousPhaseMessages: pvs.phaseMessages,
	}, nil
}

func (pvs *pointsValidationState) MemberIndex() group.MemberIndex {
	return pvs.member.ID
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

func (pjs *pointsJustificationState) DelayBlocks() uint64 {
	return silentStateDelayBlocks
}

func (pjs *pointsJustificationState) ActiveBlocks() uint64 {
	return silentStateActiveBlocks
}

func (pjs *pointsJustificationState) Initiate(ctx context.Context) error {
	pjs.member.MarkInactiveMembers(pjs.previousPhaseMessages)

	err := pjs.member.ResolvePublicKeySharePointsAccusationsMessages(
		pjs.previousPhaseMessages,
	)
	if err != nil {
		return err
	}

	return nil
}

func (pjs *pointsJustificationState) Receive(msg net.Message) error {
	return nil
}

func (pjs *pointsJustificationState) Next() (state.State, error) {
	return &keyRevealState{
		channel: pjs.channel,
		member:  pjs.member.InitializeRevealing(),
	}, nil
}

func (pjs *pointsJustificationState) MemberIndex() group.MemberIndex {
	return pjs.member.ID
}

// keyRevealState is the state during which group members reveal ephemeral
// private keys used to create an ephemeral symmetric keys with disqualified
// members who share a group private key.
//
// State covers phase 10 of the protocol.
type keyRevealState struct {
	channel net.BroadcastChannel
	member  *RevealingMember // TODO: Rename to KeyRevealingMember

	phaseMessages []*MisbehavedEphemeralKeysMessage
}

func (rs *keyRevealState) DelayBlocks() uint64 {
	return keyRevealStateDelayBlocks
}

func (rs *keyRevealState) ActiveBlocks() uint64 {
	return keyRevealStateActiveBlocks
}

func (rs *keyRevealState) Initiate(ctx context.Context) error {
	revealMsg, err := rs.member.RevealMisbehavedMembersKeys()
	if err != nil {
		return err
	}

	if err := rs.channel.Send(ctx, revealMsg); err != nil {
		return err
	}

	return nil
}

func (rs *keyRevealState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *MisbehavedEphemeralKeysMessage:
		if rs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) {
			rs.phaseMessages = append(rs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (rs *keyRevealState) Next() (state.State, error) {
	return &reconstructionState{
		channel:               rs.channel,
		member:                rs.member.InitializeReconstruction(),
		previousPhaseMessages: rs.phaseMessages,
	}, nil
}

func (rs *keyRevealState) MemberIndex() group.MemberIndex {
	return rs.member.ID
}

// reconstructionState is the state during which group members reconstruct
// individual keys of members disqualified in previous states. No messages are
// valid in this state.
//
// State covers phase 11 of the protocol.
type reconstructionState struct {
	channel net.BroadcastChannel
	member  *ReconstructingMember

	previousPhaseMessages []*MisbehavedEphemeralKeysMessage
}

func (rs *reconstructionState) DelayBlocks() uint64 {
	return silentStateDelayBlocks
}

func (rs *reconstructionState) ActiveBlocks() uint64 {
	return silentStateActiveBlocks
}

func (rs *reconstructionState) Initiate(ctx context.Context) error {
	rs.member.MarkInactiveMembers(rs.previousPhaseMessages)
	if err := rs.member.ReconstructMisbehavedIndividualKeys(
		rs.previousPhaseMessages,
	); err != nil {
		return err
	}

	return nil
}

func (rs *reconstructionState) Receive(msg net.Message) error {
	return nil
}

func (rs *reconstructionState) Next() (state.State, error) {
	return &combinationState{
		channel: rs.channel,
		member:  rs.member.InitializeCombining(),
	}, nil
}

func (rs *reconstructionState) MemberIndex() group.MemberIndex {
	return rs.member.ID
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

func (cs *combinationState) DelayBlocks() uint64 {
	return combinationStateDelayBlocks
}

func (cs *combinationState) ActiveBlocks() uint64 {
	return combinationStateActiveBlocks
}

func (cs *combinationState) Initiate(ctx context.Context) error {
	cs.member.ComputeGroupPublicKeyShares()
	cs.member.CombineGroupPublicKey()
	return nil
}

func (cs *combinationState) Receive(msg net.Message) error {
	return nil
}

func (cs *combinationState) Next() (state.State, error) {
	return &finalizationState{
		channel: cs.channel,
		member:  cs.member.InitializeFinalization(),
	}, nil
}

func (cs *combinationState) MemberIndex() group.MemberIndex {
	return cs.member.ID
}

// finalizationState is the last state of GJKR DKG protocol - in this state,
// distributed key generation is completed. No messages are valid in this state.
//
// State prepares a result to publish in phase 13 of the protocol but it does
// not execute that phase.
type finalizationState struct {
	channel net.BroadcastChannel
	member  *FinalizingMember
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

func (fs *finalizationState) Next() (state.State, error) {
	return nil, nil
}

func (fs *finalizationState) MemberIndex() group.MemberIndex {
	return fs.member.ID
}

func (fs *finalizationState) result() *Result {
	return fs.member.Result()
}
