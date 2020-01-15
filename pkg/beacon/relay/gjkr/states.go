package gjkr

import (
	"context"

	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/net"
)

type keyGenerationState = state.State

// joinState is the state during which a member announces itself to the key
// generation broadcast channel to initiate the distributed protocol.
// `JoinMessage`s are valid in this state.
type joinState struct {
	channel net.BroadcastChannel
	member  *LocalMember
}

func (js *joinState) DelayBlocks() uint64 {
	return state.DefaultMessagingStateDelayBlocks
}

func (js *joinState) ActiveBlocks() uint64 {
	return state.DefaultMessagingStateActiveBlocks
}

func (js *joinState) Initiate(ctx context.Context) error {
	return nil
}

func (js *joinState) Receive(msg net.Message) error {
	return nil
}

func (js *joinState) Next() keyGenerationState {
	return &ephemeralKeyPairGenerationState{
		channel: js.channel,
		member:  js.member.InitializeEphemeralKeysGeneration(),
	}
}

func (js *joinState) MemberIndex() group.MemberIndex {
	return js.member.ID
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
	return state.DefaultMessagingStateDelayBlocks
}

func (ekpgs *ephemeralKeyPairGenerationState) ActiveBlocks() uint64 {
	return state.DefaultMessagingStateActiveBlocks
}

func (ekpgs *ephemeralKeyPairGenerationState) Initiate(ctx context.Context) error {
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
		if !group.IsMessageFromSelf(ekpgs.member.ID, phaseMessage) &&
			group.IsSenderAccepted(ekpgs.member, phaseMessage) {
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Next() keyGenerationState {
	return &symmetricKeyGenerationState{
		channel:               ekpgs.channel,
		member:                ekpgs.member.InitializeSymmetricKeyGeneration(),
		previousPhaseMessages: ekpgs.phaseMessages,
	}
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
	return state.SilentStateDelayBlocks
}

func (skgs *symmetricKeyGenerationState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
}

func (skgs *symmetricKeyGenerationState) Initiate(ctx context.Context) error {
	skgs.member.MarkInactiveMembers(skgs.previousPhaseMessages)
	return skgs.member.GenerateSymmetricKeys(skgs.previousPhaseMessages)
}

func (skgs *symmetricKeyGenerationState) Receive(msg net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) Next() keyGenerationState {
	return &commitmentState{
		channel: skgs.channel,
		member:  skgs.member.InitializeCommitting(),
	}
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
	return state.DefaultMessagingStateDelayBlocks
}

func (cs *commitmentState) ActiveBlocks() uint64 {
	return state.DefaultMessagingStateActiveBlocks
}

func (cs *commitmentState) Initiate(ctx context.Context) error {
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
		if !group.IsMessageFromSelf(cs.member.ID, phaseMessage) &&
			group.IsSenderAccepted(cs.member, phaseMessage) {
			cs.phaseSharesMessages = append(cs.phaseSharesMessages, phaseMessage)
		}

	case *MemberCommitmentsMessage:
		if !group.IsMessageFromSelf(cs.member.ID, phaseMessage) &&
			group.IsSenderAccepted(cs.member, phaseMessage) {
			cs.phaseCommitmentsMessages = append(
				cs.phaseCommitmentsMessages,
				phaseMessage,
			)
		}
	}

	return nil
}

func (cs *commitmentState) Next() keyGenerationState {
	return &commitmentsVerificationState{
		channel: cs.channel,
		member:  cs.member.InitializeCommitmentsVerification(),

		previousPhaseSharesMessages:      cs.phaseSharesMessages,
		previousPhaseCommitmentsMessages: cs.phaseCommitmentsMessages,
	}
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
	return state.DefaultMessagingStateDelayBlocks
}

func (cvs *commitmentsVerificationState) ActiveBlocks() uint64 {
	return state.DefaultMessagingStateActiveBlocks
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

	if err := cvs.channel.Send(accusationsMsg); err != nil {
		return err
	}

	return nil
}

func (cvs *commitmentsVerificationState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *SecretSharesAccusationsMessage:
		if !group.IsMessageFromSelf(cvs.member.ID, phaseMessage) &&
			group.IsSenderAccepted(cvs.member, phaseMessage) {
			cvs.phaseAccusationsMessages = append(
				cvs.phaseAccusationsMessages,
				phaseMessage,
			)
		}
	}

	return nil
}

func (cvs *commitmentsVerificationState) Next() keyGenerationState {
	return &sharesJustificationState{
		channel: cvs.channel,
		member:  cvs.member.InitializeSharesJustification(),

		previousPhaseAccusationsMessages: cvs.phaseAccusationsMessages,
	}
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
	return state.SilentStateDelayBlocks
}

func (sjs *sharesJustificationState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
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

func (sjs *sharesJustificationState) Next() keyGenerationState {
	return &qualificationState{
		channel: sjs.channel,
		member:  sjs.member.InitializeQualified(),
	}
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
	return state.SilentStateDelayBlocks
}

func (qs *qualificationState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
}

func (qs *qualificationState) Initiate(ctx context.Context) error {
	qs.member.CombineMemberShares()
	return nil
}

func (qs *qualificationState) Receive(msg net.Message) error {
	return nil
}

func (qs *qualificationState) Next() keyGenerationState {
	return &pointsShareState{
		channel: qs.channel,
		member:  qs.member.InitializeSharing(),
	}
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
	return state.DefaultMessagingStateDelayBlocks
}

func (pss *pointsShareState) ActiveBlocks() uint64 {
	return state.DefaultMessagingStateActiveBlocks
}

func (pss *pointsShareState) Initiate(ctx context.Context) error {
	message := pss.member.CalculatePublicKeySharePoints()
	if err := pss.channel.Send(message); err != nil {
		return err
	}

	return nil
}

func (pss *pointsShareState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *MemberPublicKeySharePointsMessage:
		if !group.IsMessageFromSelf(pss.member.ID, phaseMessage) &&
			group.IsSenderAccepted(pss.member, phaseMessage) {
			pss.phaseMessages = append(pss.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (pss *pointsShareState) Next() keyGenerationState {
	return &pointsValidationState{
		channel: pss.channel,
		member:  pss.member,

		previousPhaseMessages: pss.phaseMessages,
	}
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
	return state.DefaultMessagingStateDelayBlocks
}

func (pvs *pointsValidationState) ActiveBlocks() uint64 {
	return state.DefaultMessagingStateActiveBlocks
}

func (pvs *pointsValidationState) Initiate(ctx context.Context) error {
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
		if !group.IsMessageFromSelf(pvs.member.ID, phaseMessage) &&
			group.IsSenderAccepted(pvs.member, phaseMessage) {
			pvs.phaseMessages = append(pvs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (pvs *pointsValidationState) Next() keyGenerationState {
	return &pointsJustificationState{
		channel: pvs.channel,
		member:  pvs.member.InitializePointsJustification(),

		previousPhaseMessages: pvs.phaseMessages,
	}
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
	return state.SilentStateDelayBlocks
}

func (pjs *pointsJustificationState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
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

func (pjs *pointsJustificationState) Next() keyGenerationState {
	return &keyRevealState{
		channel: pjs.channel,
		member:  pjs.member.InitializeRevealing(),
	}
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
	return state.DefaultMessagingStateDelayBlocks
}

func (rs *keyRevealState) ActiveBlocks() uint64 {
	return state.DefaultMessagingStateActiveBlocks
}

func (rs *keyRevealState) Initiate(ctx context.Context) error {
	revealMsg, err := rs.member.RevealMisbehavedMembersKeys()
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
	case *MisbehavedEphemeralKeysMessage:
		if !group.IsMessageFromSelf(rs.member.ID, phaseMessage) &&
			group.IsSenderAccepted(rs.member, phaseMessage) {
			rs.phaseMessages = append(rs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (rs *keyRevealState) Next() keyGenerationState {
	return &reconstructionState{
		channel:               rs.channel,
		member:                rs.member.InitializeReconstruction(),
		previousPhaseMessages: rs.phaseMessages,
	}
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
	return state.SilentStateDelayBlocks
}

func (rs *reconstructionState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
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

func (rs *reconstructionState) Next() keyGenerationState {
	return &combinationState{
		channel: rs.channel,
		member:  rs.member.InitializeCombining(),
	}
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
	return state.SilentStateDelayBlocks
}

func (cs *combinationState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
}

func (cs *combinationState) Initiate(ctx context.Context) error {
	cs.member.CombineGroupPublicKey()
	return nil
}

func (cs *combinationState) Receive(msg net.Message) error {
	return nil
}

func (cs *combinationState) Next() keyGenerationState {
	return &finalizationState{
		channel: cs.channel,
		member:  cs.member.InitializeFinalization(),
	}
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
	return state.SilentStateDelayBlocks
}

func (fs *finalizationState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
}

func (fs *finalizationState) Initiate(ctx context.Context) error {
	return nil
}

func (fs *finalizationState) Receive(msg net.Message) error {
	return nil
}

func (fs *finalizationState) Next() keyGenerationState {
	return nil
}

func (fs *finalizationState) MemberIndex() group.MemberIndex {
	return fs.member.ID
}

func (fs *finalizationState) result() *Result {
	return fs.member.Result()
}
