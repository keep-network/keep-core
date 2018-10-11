package dkg

import "math/big"

// Group is protocol's members group.
type Group struct {
	// The number of members in the complete group.
	groupSize int
	// The maximum number of group members who could be dishonest in order for the
	// generated key to be uncompromised.
	dishonestThreshold int
	// IDs of all members of the group. Contains local member's ID.
	// Initially empty, populated as each other member announces its presence.
	memberIDs []*big.Int
}

// MemberIDs returns IDs of all group members.
func (g *Group) MemberIDs() []*big.Int {
	return g.memberIDs
}
