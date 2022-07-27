package gjkr

import "github.com/keep-network/keep-core/pkg/protocol/group"

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (em *SymmetricKeyGeneratingMember) MarkInactiveMembers(
	ephemeralPubKeyMessages []*EphemeralPublicKeyMessage,
) {
	filter := em.inactiveMemberFilter()
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
	filter := cvm.inactiveMemberFilter()
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
	filter := cvm.inactiveMemberFilter()
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
	filter := sm.inactiveMemberFilter()
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
	filter := cvm.inactiveMemberFilter()
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
	filter := rm.inactiveMemberFilter()
	for _, message := range messages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// inactiveMemberFilter returns a new instance of the inactive member filter.
func (mc *memberCore) inactiveMemberFilter() *group.InactiveMemberFilter {
	return group.NewInactiveMemberFilter(mc.logger, mc.ID, mc.group)
}

// shouldAcceptMessage indicates whether the given member should accept
// a message from the given sender.
func (mc *memberCore) shouldAcceptMessage(
	senderID group.MemberIndex,
	senderPublicKey []byte,
) bool {
	isMessageFromSelf := senderID == mc.ID
	isSenderValid := mc.membershipValidator.IsValidMembership(
		senderID,
		senderPublicKey,
	)
	isSenderAccepted := mc.group.IsOperating(senderID)

	return !isMessageFromSelf && isSenderValid && isSenderAccepted
}
