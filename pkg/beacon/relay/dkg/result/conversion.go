package result

import (
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

// convertResult transforms GJKR protocol execution result to a chain specific
// DKG result form. It serializes a group public key to bytes and converts
// disqualified and inactive members lists to a boolean list where each entry
// corresponds to a member in the group and true/false value indicates status of
// the member.
func convertResult(gjkrResult *gjkr.Result, groupSize int) *relayChain.DKGResult {
	groupPublicKey := make([]byte, 0)

	// We convert the point G2, to compress the point correctly
	// (ensuring we encode the parity bit).
	if gjkrResult.GroupPublicKey != nil {
		groupPublicKey = gjkrResult.GroupPublicKey.Marshal()
	}

	// convertToByteSlice converts slice containing member IDs to a slice of
	// bytes. E.g. for input = {1, 3, 20}, the output is {0x01, 0x03, 0x14}.
	// MemberIndex cannot be larger than 255.
	convertToByteSlice := func(memberIDsSlice []group.MemberIndex) []byte {
		bytes := make([]byte, len(memberIDsSlice))
		for i, memberID := range memberIDsSlice {
			bytes[i] = byte(memberID)
		}

		return bytes
	}

	return &relayChain.DKGResult{
		GroupPublicKey: groupPublicKey,
		Inactive:       convertToByteSlice(gjkrResult.Group.InactiveMemberIDs()),
		Disqualified:   convertToByteSlice(gjkrResult.Group.DisqualifiedMemberIDs()),
	}
}
