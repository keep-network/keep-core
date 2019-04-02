package gjkr

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

func newDKGResult(
	groupSize int,
	disqualifiedMembers []group.MemberIndex,
	inactiveMembers []group.MemberIndex,
) *Result {
	group := initalizeGroup(groupSize, disqualifiedMembers, inactiveMembers)

	return &Result{Group: group}
}
