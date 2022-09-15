package group

import "github.com/ipfs/go-log/v2"

// InactiveMemberFilter is a proxy facilitates filtering out inactive members
// in the given phase and registering their final list in the Group.
type InactiveMemberFilter struct {
	logger log.StandardLogger

	selfMemberID MemberIndex
	group        *Group

	phaseActiveMembers []MemberIndex
}

// NewInactiveMemberFilter creates a new instance of InactiveMemberFilter.
// It accepts member index of the current member (the one which will be
// filtering out other group members for inactivity) and the reference to Group
// to which all those members belong.
func NewInactiveMemberFilter(
	logger log.StandardLogger,
	selfMemberIndex MemberIndex,
	group *Group,
) *InactiveMemberFilter {
	return &InactiveMemberFilter{
		logger:             logger,
		selfMemberID:       selfMemberIndex,
		group:              group,
		phaseActiveMembers: make([]MemberIndex, 0),
	}
}

// MarkMemberAsActive marks member with the given index as active in the given
// phase.
func (mf *InactiveMemberFilter) MarkMemberAsActive(memberID MemberIndex) {
	mf.phaseActiveMembers = append(mf.phaseActiveMembers, memberID)
}

// FlushInactiveMembers takes all members who were not previously marked as
// active and flushes them to the group as inactive members.
func (mf *InactiveMemberFilter) FlushInactiveMembers() {
	isActive := func(id MemberIndex) bool {
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
			mf.logger.Warnf(
				"[member:%v] marking member [%v] as inactive",
				mf.selfMemberID,
				operatingMemberID,
			)
			mf.group.MarkMemberAsInactive(operatingMemberID)
		}
	}
}
