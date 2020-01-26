package group

import "crypto/ecdsa"

// MembershipValidator operates on a group selection result and lets to
// validate one's membership based on the provided public key.
//
// Validator is used to filter out messages from parties not selected to
// the group. It is also used to confirm the position in the group of
// a party that was selected. This is used to validate messages sent by that
// party to all other group members.
type MembershipValidator interface {
	// IsInGroup returns true if party with the given public key has been selected
	// to the group. Otherwise, function returns false.
	IsInGroup(publicKey *ecdsa.PublicKey) bool

	// IsSelectedAtIndex returns true if party with the given public key has been
	// selected to the group at the given position. If the position does not match
	// function returns false. The same happens when the party was not selected
	// to the group.
	IsSelectedAtIndex(index int, publicKey *ecdsa.PublicKey) bool
}

// Group is protocol's members group.
type Group struct {
	// The maximum number of misbehaving participants for which it is still
	// possible to generate a signature.
	dishonestThreshold int
	// IDs of all disqualified members of the group.
	disqualifiedMemberIDs []MemberIndex
	// IDs of all inactive members of the group.
	inactiveMemberIDs []MemberIndex
	// All member IDs in this group.
	memberIDs []MemberIndex
}

// NewDkgGroup creates a new Group with the provided dishonest threshold, member
// identifiers, and empty IA and DQ members list.
func NewDkgGroup(dishonestThreshold int, size int) *Group {
	memberIDs := make([]MemberIndex, size)
	for i := 0; i < size; i++ {
		memberIDs[i] = MemberIndex(i + 1)
	}

	return &Group{
		dishonestThreshold:    dishonestThreshold,
		disqualifiedMemberIDs: []MemberIndex{},
		inactiveMemberIDs:     []MemberIndex{},
		memberIDs:             memberIDs,
	}
}

// MemberIDs returns IDs of all group members, as initially selected to the
// group. Returned list contains IDs of all members, including those marked as
// inactive or disqualified.
func (g *Group) MemberIDs() []MemberIndex {
	return g.memberIDs
}

// GroupSize returns the full size of the group, including IA and DQ members.
func (g *Group) GroupSize() int {
	return len(g.memberIDs)
}

// DishonestThreshold returns value of the dishonest members threshold as set
// for the group.
func (g *Group) DishonestThreshold() int {
	return g.dishonestThreshold
}

// DisqualifiedMemberIDs returns indexes of all group members that have been
// disqualified during DKG protocol execution.
func (g *Group) DisqualifiedMemberIDs() []MemberIndex {
	return g.disqualifiedMemberIDs
}

// InactiveMemberIDs returns indexes of all group members that have been marked
// as inactive during DKG protocol execution.
func (g *Group) InactiveMemberIDs() []MemberIndex {
	return g.inactiveMemberIDs
}

// OperatingMemberIDs returns IDs of all group members that are active and have
// not been disqualified. All those members are properly operating in the group
// at the moment of calling this method.
func (g *Group) OperatingMemberIDs() []MemberIndex {
	operatingMembers := make([]MemberIndex, 0)
	for _, member := range g.MemberIDs() {
		if g.IsOperating(member) {
			operatingMembers = append(operatingMembers, member)
		}
	}

	return operatingMembers
}

// MarkMemberAsDisqualified adds the member with the given ID to the list of
// disqualified members. If the member is not a part of the group, is already
// disqualified or marked as inactive, method does nothing.
func (g *Group) MarkMemberAsDisqualified(memberID MemberIndex) {
	if g.IsOperating(memberID) {
		g.disqualifiedMemberIDs = append(g.disqualifiedMemberIDs, memberID)
	}
}

// MarkMemberAsInactive adds the member with the given ID to the list of
// inactive members. If the member is not a part of the group, is already
// disqualified or marked as inactive, method does nothing.
func (g *Group) MarkMemberAsInactive(memberID MemberIndex) {
	if g.IsOperating(memberID) {
		g.inactiveMemberIDs = append(g.inactiveMemberIDs, memberID)
	}
}

// IsOperating returns true if member with the given index has not been marked
// as IA or DQ in the group.
func (g *Group) IsOperating(memberID MemberIndex) bool {
	return g.isInGroup(memberID) &&
		!g.isInactive(memberID) &&
		!g.isDisqualified(memberID)
}

func (g *Group) isInGroup(memberID MemberIndex) bool {
	for _, groupMember := range g.MemberIDs() {
		if groupMember == memberID {
			return true
		}
	}

	return false
}

func (g *Group) isInactive(memberID MemberIndex) bool {
	for _, inactiveMemberID := range g.inactiveMemberIDs {
		if memberID == inactiveMemberID {
			return true
		}
	}

	return false
}

func (g *Group) isDisqualified(memberID MemberIndex) bool {
	for _, disqualifiedMemberID := range g.disqualifiedMemberIDs {
		if memberID == disqualifiedMemberID {
			return true
		}
	}

	return false
}

func (g *Group) eliminatedMembersCount() int {
	return len(g.disqualifiedMemberIDs) + len(g.inactiveMemberIDs)
}

// isThresholdSatisfied checks number of disqualified and inactive members in
// the group. If the number is less or equal half of dishonest threshold,
// returns true.
func (g *Group) isThresholdSatisfied() bool {
	return g.eliminatedMembersCount() <= g.dishonestThreshold/2
}
