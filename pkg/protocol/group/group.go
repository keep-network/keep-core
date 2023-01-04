// Package group contains an implementation of a generic protocol group
// and auxiliary tools that help during group-related operations.
package group

// MemberIndex is an index of a member in a group. The maximum member index
// value is 255.
type MemberIndex = uint8

// MaxMemberIndex denotes the maximum value of the MemberIndex type. That type
// is represented as uint8 so the maximum member index is 255.
const MaxMemberIndex = 255

// Group is protocol's members group.
type Group struct {
	// The maximum number of misbehaving participants for which it is still
	// possible to generate a signature.
	dishonestThreshold int
	// Indexes of all disqualified members of the group.
	disqualifiedMemberIndexes []MemberIndex
	// Indexes of all inactive members of the group.
	inactiveMemberIndexes []MemberIndex
	// All member indexes in this group.
	memberIndexes []MemberIndex
}

// NewGroup creates a new Group with the provided dishonest threshold and empty
// IA and DQ members list.
func NewGroup(dishonestThreshold int, size int) *Group {
	memberIndexes := make([]MemberIndex, size)
	for i := 0; i < size; i++ {
		memberIndexes[i] = MemberIndex(i + 1)
	}

	return &Group{
		dishonestThreshold:        dishonestThreshold,
		disqualifiedMemberIndexes: []MemberIndex{},
		inactiveMemberIndexes:     []MemberIndex{},
		memberIndexes:             memberIndexes,
	}
}

// MemberIndexes returns indexes of all group members, as initially selected to
// the group. Returned list contains indexes of all members, including those
// marked as inactive or disqualified.
func (g *Group) MemberIndexes() []MemberIndex {
	return g.memberIndexes
}

// GroupSize returns the full size of the group, including IA and DQ members.
func (g *Group) GroupSize() int {
	return len(g.memberIndexes)
}

// DishonestThreshold returns value of the dishonest members threshold as set
// for the group.
func (g *Group) DishonestThreshold() int {
	return g.dishonestThreshold
}

// HonestThreshold returns value of the honest members threshold as set
// for the group.
func (g *Group) HonestThreshold() int {
	return g.GroupSize() - g.DishonestThreshold()
}

// DisqualifiedMemberIndexes returns indexes of all group members that have been
// disqualified during the protocol execution.
func (g *Group) DisqualifiedMemberIndexes() []MemberIndex {
	return g.disqualifiedMemberIndexes
}

// InactiveMemberIndexes returns indexes of all group members that have been
// marked as inactive during the protocol execution.
func (g *Group) InactiveMemberIndexes() []MemberIndex {
	return g.inactiveMemberIndexes
}

// OperatingMemberIndexes returns indexes of all group members that are active
// and have not been disqualified. All those members are properly operating in
// the group at the moment of calling this method.
func (g *Group) OperatingMemberIndexes() []MemberIndex {
	operatingMembers := make([]MemberIndex, 0)
	for _, member := range g.MemberIndexes() {
		if g.IsOperating(member) {
			operatingMembers = append(operatingMembers, member)
		}
	}

	return operatingMembers
}

// MarkMemberAsDisqualified adds the member with the given index to the list of
// disqualified members. If the member is not a part of the group, is already
// disqualified or marked as inactive, method does nothing.
func (g *Group) MarkMemberAsDisqualified(memberIndex MemberIndex) {
	if g.IsOperating(memberIndex) {
		g.disqualifiedMemberIndexes = append(g.disqualifiedMemberIndexes, memberIndex)
	}
}

// MarkMemberAsInactive adds the member with the given index to the list of
// inactive members. If the member is not a part of the group, is already
// disqualified or marked as inactive, method does nothing.
func (g *Group) MarkMemberAsInactive(memberIndex MemberIndex) {
	if g.IsOperating(memberIndex) {
		g.inactiveMemberIndexes = append(g.inactiveMemberIndexes, memberIndex)
	}
}

// IsOperating returns true if member with the given index has not been marked
// as IA or DQ in the group.
func (g *Group) IsOperating(memberIndex MemberIndex) bool {
	return g.isInGroup(memberIndex) &&
		!g.isInactive(memberIndex) &&
		!g.isDisqualified(memberIndex)
}

func (g *Group) isInGroup(memberIndex MemberIndex) bool {
	for _, groupMember := range g.MemberIndexes() {
		if groupMember == memberIndex {
			return true
		}
	}

	return false
}

func (g *Group) isInactive(memberIndex MemberIndex) bool {
	for _, inactiveMemberIndex := range g.inactiveMemberIndexes {
		if memberIndex == inactiveMemberIndex {
			return true
		}
	}

	return false
}

func (g *Group) isDisqualified(memberIndex MemberIndex) bool {
	for _, disqualifiedMemberIndex := range g.disqualifiedMemberIndexes {
		if memberIndex == disqualifiedMemberIndex {
			return true
		}
	}

	return false
}
