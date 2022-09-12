package signing

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// InactiveMembersError is raised when inactive members were detected during
// the execution of the signing protocol. A member is considered inactive when
// a required network message from him is not received within the expected
// time window.
type InactiveMembersError struct {
	InactiveMembersIndexes []group.MemberIndex
}

func newInactiveMembersError(
	inactiveMembersIndexes []group.MemberIndex,
) *InactiveMembersError {
	return &InactiveMembersError{inactiveMembersIndexes}
}

func (ime *InactiveMembersError) Error() string {
	return fmt.Sprintf(
		"inactive members: [%v]",
		ime.InactiveMembersIndexes,
	)
}
