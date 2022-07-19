package result

import (
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"sort"

	"github.com/keep-network/keep-core/pkg/beacon/gjkr"
	"github.com/keep-network/keep-core/pkg/group"
)

// convertGjkrResult transforms GJKR protocol execution result to a chain
// specific DKG result form. It serializes a group public key to bytes and
// converts disqualified and inactive members lists to one list of misbehaving
// participants where each byte represents misbehaving member index.
func convertGjkrResult(gjkrResult *gjkr.Result) *beaconchain.DKGResult {
	groupPublicKey := make([]byte, 0)

	// We convert the point G2, to compress the point correctly
	// (ensuring we encode the parity bit).
	if bytes, err := gjkrResult.GroupPublicKeyBytes(); err == nil {
		groupPublicKey = bytes
	}

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

		// contert sorted list of member indexes into bytes
		bytes := make([]byte, len(sorted))
		for i, m := range sorted {
			bytes[i] = byte(m)
		}

		return bytes
	}

	return &beaconchain.DKGResult{
		GroupPublicKey: groupPublicKey,
		Misbehaved: convertToMisbehaved(
			gjkrResult.Group.InactiveMemberIDs(),
			gjkrResult.Group.DisqualifiedMemberIDs(),
		),
	}
}
