package gjkr

import "fmt"

// Group is protocol's members group.
type Group struct {
	// The maximum number of group members who could be dishonest in order for the
	// generated key to be uncompromised.
	dishonestThreshold int
	// IDs of all members of the group. Contains local member's ID.
	// Initially empty, populated as each other member announces its presence.
	memberIDs []MemberID
	// IDs of all disqualified members of the group.
	disqualifiedMemberIDs []MemberID
	// IDs of all inactive members of the group.
	inactiveMemberIDs []MemberID
}

// MemberIDs returns IDs of all group members.
// TODO Add ActiveMemberIDs() method to return only active members of the group
// and use it across the protocol phases to get other members IDs.
func (g *Group) MemberIDs() []MemberID {
	return g.memberIDs
}

// RegisterMemberID adds a member to the list of group members.
func (g *Group) RegisterMemberID(memberID MemberID) error {
	if err := memberID.validate(); err != nil {
		return fmt.Errorf("cannot register member ID in the group [%v]", err)
	}

	for _, id := range g.memberIDs {
		if id == memberID {
			return nil // already there
		}
	}
	g.memberIDs = append(g.memberIDs, memberID)

	return nil
}

func (g *Group) eliminatedMembersCount() int {
	return len(g.disqualifiedMemberIDs) + len(g.inactiveMemberIDs)
}

// isThresholdSatisfied checks number of disqualified and inactive members in the
// group. If the number is less or equal half of dishonest threshold, returns true.
func (g *Group) isThresholdSatisfied() bool {
	return g.eliminatedMembersCount() <= g.dishonestThreshold/2
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
