package gjkr

import "github.com/keep-network/keep-core/pkg/protocol/group"

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (em *SymmetricKeyGeneratingMember) MarkInactiveMembers(
	ephemeralPubKeyMessages []*EphemeralPublicKeyMessage,
) {
	filter := em.messageFilter()
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
	filter := cvm.messageFilter()
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
	filter := cvm.messageFilter()
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
	filter := sm.messageFilter()
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
	filter := cvm.messageFilter()
	for _, message := range pointsAccusationsMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (rm *ReconstructingMember) MarkInactiveMembers(
	messages []*MisbehavedEphemeralKeysMessage,
) {
	filter := rm.messageFilter()
	for _, message := range messages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

func (mc *memberCore) messageFilter() *group.InactiveMemberFilter {
	return group.NewInactiveMemberFilter(mc.ID, mc.group)
}

func (mc *memberCore) IsSenderAccepted(senderID group.MemberIndex) bool {
	return mc.group.IsOperating(senderID)
}

func (mc *memberCore) IsSenderValid(
	senderID group.MemberIndex,
	senderPublicKey []byte,
) bool {
	return mc.membershipValidator.IsValidMembership(senderID, senderPublicKey)
}
