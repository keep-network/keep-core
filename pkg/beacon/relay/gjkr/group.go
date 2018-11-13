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
}

// MemberIDs returns IDs of all group members.
func (g *Group) MemberIDs() []int {
	return g.memberIDs
}

// RegisterMemberID adds a member to the list of group members.
func (g *Group) RegisterMemberID(id int) {
	g.memberIDs = append(g.memberIDs, id)
}

// DisqualifiedMembers returns members disqualified during protocol execution.
func (g *Group) DisqualifiedMembers() []int {
	// TODO: Implement
	return []int{}
}

// InactiveMembers returns members inactive during protocol execution.
func (g *Group) InactiveMembers() []int {
	// TODO: Implement
	return []int{}
}
