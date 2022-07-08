package result

import (
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/group"
)

func TestConvertResult(t *testing.T) {
	publicKey := new(bn256.G2).ScalarBaseMult(big.NewInt(2))
	marshalledPublicKey := publicKey.Marshal()

	var tests = map[string]struct {
		disqualifiedMemberIDs []group.MemberIndex
		inactiveMemberIDs     []group.MemberIndex
		gjkrResult            *gjkr.Result
		expectedResult        *beaconChain.DKGResult
	}{
		"group public key not provided, DQ and IA empty": {
			disqualifiedMemberIDs: []group.MemberIndex{},
			inactiveMemberIDs:     []group.MemberIndex{},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: nil,
				Group:          group.NewDkgGroup(32, 64),
			},
			expectedResult: &beaconChain.DKGResult{
				GroupPublicKey: []byte{},
				Misbehaved:     []byte{},
			},
		},
		"group public key provided, DQ and IA empty": {
			disqualifiedMemberIDs: []group.MemberIndex{},
			inactiveMemberIDs:     []group.MemberIndex{},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewDkgGroup(32, 64),
			},
			expectedResult: &beaconChain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Misbehaved:     []byte{},
			},
		},
		"group public key provided, both DQ and IA non-empty": {
			disqualifiedMemberIDs: []group.MemberIndex{1, 4, 3, 50},
			inactiveMemberIDs:     []group.MemberIndex{5, 3, 50},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewDkgGroup(32, 64),
			},
			expectedResult: &beaconChain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Misbehaved:     []byte{0x01, 0x03, 0x04, 0x05, 0x32},
			},
		},
		"group public key provided, DQ empty, IA non-empty": {
			disqualifiedMemberIDs: []group.MemberIndex{},
			inactiveMemberIDs:     []group.MemberIndex{5},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewDkgGroup(32, 64),
			},
			expectedResult: &beaconChain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Misbehaved:     []byte{0x05},
			},
		},
		"group public key provided, DQ non-empty, IA empty": {
			disqualifiedMemberIDs: []group.MemberIndex{60, 1, 5},
			inactiveMemberIDs:     []group.MemberIndex{},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewDkgGroup(32, 64),
			},
			expectedResult: &beaconChain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Misbehaved:     []byte{0x01, 0x05, 0x3C},
			},
		},
	}
	for _, test := range tests {
		for _, disqualifiedMember := range test.disqualifiedMemberIDs {
			test.gjkrResult.Group.MarkMemberAsDisqualified(disqualifiedMember)
		}

		for _, inactiveMember := range test.inactiveMemberIDs {
			test.gjkrResult.Group.MarkMemberAsInactive(inactiveMember)
		}

		convertedResult := convertGjkrResult(test.gjkrResult)

		if !test.expectedResult.Equals(convertedResult) {
			t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedResult, convertedResult)
		}
	}
}
