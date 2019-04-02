package gjkr

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

func initalizeGroup(
	groupSize int,
	disqualifiedMembers []group.MemberIndex,
	inactiveMembers []group.MemberIndex,
) *group.Group {
	dkgGroup := group.NewEmptyDkgGroup(groupSize/2 + 1)
	for i := 1; i <= groupSize; i++ {
		dkgGroup.RegisterMemberID(group.MemberIndex(i))
	}

	for _, disqualified := range disqualifiedMembers {
		dkgGroup.MarkMemberAsDisqualified(disqualified)
	}
	for _, inactive := range inactiveMembers {
		dkgGroup.MarkMemberAsInactive(inactive)
	}
	return dkgGroup
}

func newDKGResult(
	groupSize int,
	disqualifiedMembers []group.MemberIndex,
	inactiveMembers []group.MemberIndex,
) *Result {
	group := initalizeGroup(groupSize, disqualifiedMembers, inactiveMembers)

	return &Result{Group: group}
}
