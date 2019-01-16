package gjkr

func (mc *memberCore) messageFilter() *messageFilter {
	return &messageFilter{
		selfMemberID:       mc.ID,
		group:              mc.group,
		phaseActiveMembers: make([]MemberID, 0),
	}
}

type messageFilter struct {
	selfMemberID MemberID
	group        *Group

	phaseActiveMembers []MemberID
}

func (mf *messageFilter) markMemberAsActive(memberID MemberID) {
	mf.phaseActiveMembers = append(mf.phaseActiveMembers, memberID)
}

func (mf *messageFilter) flushInactiveMembers() {
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
