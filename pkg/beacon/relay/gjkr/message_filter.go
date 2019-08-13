package gjkr

import "github.com/keep-network/keep-core/pkg/beacon/relay/group"

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (em *SymmetricKeyGeneratingMember) MarkInactiveMembers(
	ephemeralPubKeyMessages []*EphemeralPublicKeyMessage,
) {
	filter := em.messageIAFilter()
	for _, message := range ephemeralPubKeyMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (cvm *CommitmentsVerifyingMember) MarkInactiveMembers(
	sharesMessages []*PeerSharesMessage,
	commitmentsMessages []*MemberCommitmentsMessage,
) {
	filter := cvm.messageIAFilter()
	for _, sharesMessage := range sharesMessages {
		for _, commitmentsMessage := range commitmentsMessages {
			if sharesMessage.senderID == commitmentsMessage.senderID {
				filter.MarkMemberAsActive(sharesMessage.senderID)
				break
			}
		}
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (cvm *SharesJustifyingMember) MarkInactiveMembers(
	sharesAccusationsMessages []*SecretSharesAccusationsMessage,
) {
	filter := cvm.messageIAFilter()
	for _, message := range sharesAccusationsMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (sm *SharingMember) MarkInactiveMembers(
	keySharePointsMessages []*MemberPublicKeySharePointsMessage,
) {
	filter := sm.messageIAFilter()
	for _, message := range keySharePointsMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (cvm *PointsJustifyingMember) MarkInactiveMembers(
	pointsAccusationsMessages []*PointsAccusationsMessage,
) {
	filter := cvm.messageIAFilter()
	for _, message := range pointsAccusationsMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (rm *ReconstructingMember) MarkInactiveMembers(
	disqialifiedKeysMessages []*DisqualifiedEphemeralKeysMessage,
) {
	filter := rm.messageIAFilter()
	for _, message := range disqialifiedKeysMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

func (mc *memberCore) messageIAFilter() *group.InactiveMemberFilter {
	return group.NewInactiveMemberFilter(mc.ID, mc.group)
}

func (mc *memberCore) messageDQFilter() *group.DisqualifiedMemberFilter {
	return group.NewDisqualifiedMemberFilter(mc.ID, mc.group)
}

func (mc *memberCore) IsSenderAccepted(senderID group.MemberIndex) bool {
	return mc.group.IsOperating(senderID)
}
