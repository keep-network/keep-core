package gjkr

// MessageFiltering interface defines method allowing to filter out messages
// from members that are not part of the group or were marked as IA or DQ.
type MessageFiltering interface {
	IsSenderAccepted(senderID MemberID) bool
}

// IsSenderAccepted returns true if the message from the given sender should be
// accepted for further processing. Otherwise, function returns false.
// Message from the given sender is allowed only if that member is a properly
// operating group member - it was not DQ or IA so far.
func (mc *memberCore) IsSenderAccepted(senderID MemberID) bool {
	return mc.group.isOperating(senderID)
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (em *SymmetricKeyGeneratingMember) MarkInactiveMembers(
	ephemeralPubKeyMessages []*EphemeralPublicKeyMessage,
) {
	filter := em.messageFilter()
	for _, message := range ephemeralPubKeyMessages {
		filter.markMemberAsActive(message.senderID)
	}

	filter.flushInactiveMembers()
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
				filter.markMemberAsActive(sharesMessage.senderID)
				break
			}
		}
	}

	filter.flushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (sm *SharingMember) MarkInactiveMembers(
	keySharePointsMessages []*MemberPublicKeySharePointsMessage,
) {
	filter := sm.messageFilter()
	for _, message := range keySharePointsMessages {
		filter.markMemberAsActive(message.senderID)
	}

	filter.flushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (rm *ReconstructingMember) MarkInactiveMembers(
	disqialifiedKeysMessages []*DisqualifiedEphemeralKeysMessage,
) {
	filter := rm.messageFilter()
	for _, message := range disqialifiedKeysMessages {
		filter.markMemberAsActive(message.senderID)
	}

	filter.flushInactiveMembers()
}

func (mc *memberCore) messageFilter() *inactiveMemberFilter {
	return &inactiveMemberFilter{
		selfMemberID:       mc.ID,
		group:              mc.group,
		phaseActiveMembers: make([]MemberID, 0),
	}
}

type inactiveMemberFilter struct {
	selfMemberID MemberID
	group        *Group

	phaseActiveMembers []MemberID
}

func (mf *inactiveMemberFilter) markMemberAsActive(memberID MemberID) {
	mf.phaseActiveMembers = append(mf.phaseActiveMembers, memberID)
}

func (mf *inactiveMemberFilter) flushInactiveMembers() {
	isActive := func(id MemberID) bool {
		if id == mf.selfMemberID {
			return true
		}

		for _, activeMemberID := range mf.phaseActiveMembers {
			if activeMemberID == id {
				return true
			}
		}

		return false
	}

	for _, operatingMemberID := range mf.group.OperatingMemberIDs() {
		if !isActive(operatingMemberID) {
			mf.group.MarkMemberAsInactive(operatingMemberID)
		}
	}
}
