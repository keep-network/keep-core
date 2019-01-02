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
	// IDs of all disqualified members of the group.
	disqualifiedMemberIDs []MemberID
}

// MemberIDs returns IDs of all group members.
// TODO Add ActiveMemberIDs() method to return only active members of the group
// and use it across the protocol phases to get other members IDs.
func (g *Group) MemberIDs() []MemberID {
	return g.memberIDs
}

// RegisterMemberID adds a member to the list of group members.
func (g *Group) RegisterMemberID(id MemberID) {
	g.memberIDs = append(g.memberIDs, id)
}

// DisqualifyMemberID adds a member to the list of disqualified members.
func (g *Group) DisqualifyMemberID(id MemberID) {
	for _, currentID := range g.disqualifiedMemberIDs {
		if currentID == id {
			return
		}
	}

	g.disqualifiedMemberIDs = append(g.disqualifiedMemberIDs, id)
}
