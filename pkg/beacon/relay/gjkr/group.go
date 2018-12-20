package gjkr

// Group is protocol's members group.
type Group struct {
	// The number of members in the complete group.
	groupSize int
	// The maximum number of group members who could be dishonest in order for the
	// generated key to be uncompromised.
	dishonestThreshold int
	// IDs of all members of the group. Contains local member's ID.
	// Initially empty, populated as each other member announces its presence.
	memberIDs []MemberID
	// IDs of group members who were disqualified during protocol execution.
	disqualifiedMemberIDs []MemberID
	// IDs of group members who went inactive during protocol execution.
	inactiveMemberIDs []MemberID
}

// MemberIDs returns IDs of all group members.
func (g *Group) MemberIDs() []MemberID {
	return g.memberIDs
}

// RegisterMemberID adds a member to the list of group members.
func (g *Group) RegisterMemberID(id MemberID) {
	g.memberIDs = append(g.memberIDs, id)
}

// DisqualifiedMembers returns members disqualified during protocol execution.
func (g *Group) DisqualifiedMembers() []MemberID {
	return g.disqualifiedMemberIDs
}

// InactiveMembers returns members inactive during protocol execution.
func (g *Group) InactiveMembers() []MemberID {
	return g.inactiveMemberIDs
}

func (g *Group) eliminatedMembersCount() int {
	return len(g.disqualifiedMemberIDs) + len(g.inactiveMemberIDs)
}

// isThresholdSatisfied checks number of disqualified and inactive members in the
// group. If the number is less or equal half of dishonest threshold returns true.
func (g *Group) isThresholdSatisfied() bool {
	return g.eliminatedMembersCount() <= g.dishonestThreshold/2
}
