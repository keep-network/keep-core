package dkg

import (
	"sort"

	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// Result of distributed key generation protocol.
type Result struct {
	// Group represents the group state, including members, disqualified,
	// and inactive members.
	Group *group.Group
	// TODO: Temporary result. Add real items.
	GroupPublicKey []byte
}

// TODO: Consider removing and using just one DKG result type.
func convertToChainResult(result *Result) *DKGResult {
	convertToMisbehaved := func(
		inactive []group.MemberIndex,
		disqualified []group.MemberIndex,
	) []byte {
		// merge IA and DQ into 'misbehaved' set
		misbehaving := make(map[group.MemberIndex]bool)
		for _, ia := range inactive {
			misbehaving[ia] = true
		}
		for _, dq := range disqualified {
			misbehaving[dq] = true
		}

		// convert 'misbehaved' set into sorted list
		var sorted []group.MemberIndex
		for m := range misbehaving {
			sorted = append(sorted, m)
		}
		sort.Slice(sorted[:], func(i, j int) bool {
			return sorted[i] < sorted[j]
		})

		// convert sorted list of member indexes into bytes
		bytes := make([]byte, len(sorted))
		for i, m := range sorted {
			bytes[i] = byte(m)
		}

		return bytes
	}

	return &DKGResult{
		// TODO: Check if GroupPublicKey needs further conversion
		GroupPublicKey: result.GroupPublicKey,
		Misbehaved: convertToMisbehaved(
			result.Group.InactiveMemberIDs(),
			result.Group.DisqualifiedMemberIDs(),
		),
	}
}
