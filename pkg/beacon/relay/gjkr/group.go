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
}

// MemberIDs returns IDs of all group members.
func (g *Group) MemberIDs() []MemberID {
	return g.memberIDs
}

// RegisterMemberID adds a member to the list of group members.
func (g *Group) RegisterMemberID(memberID MemberID) {
	for _, id := range g.memberIDs {
		if id == memberID {
			return // already there
		}
	}
	g.memberIDs = append(g.memberIDs, memberID)
}
