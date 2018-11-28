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
	memberIDs []int
	// IDs of group members who were disqualified during protocol execution.
	disqualifiedMemberIDs []int
	// IDs of group members who went inactive during protocol execution.
	inactiveMemberIDs []int
}

// MemberIDs returns IDs of all group members.
func (g *Group) MemberIDs() []int {
	return g.memberIDs
}

// RegisterMemberID adds a member to the list of group members.
func (g *Group) RegisterMemberID(id int) {
	g.memberIDs = append(g.memberIDs, id)
}

func (g *Group) DisqualifiedMembers() []int {
	// TODO: Implement
	return []int{}
}

func (g *Group) InactiveMembers() []int {
	// TODO: Implement
	return []int{}
}

func (g *Group) SetDisqualifiedMembers(dq []int) {
	g.disqualifiedMemberIDs = dq
}

func (g *Group) SetInactiveMembers(iq []int) {
	g.inactiveMemberIDs = iq
}
