package result

import (
	"math/big"
	"testing"

	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/gjkr"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

func TestConvertResult(t *testing.T) {
	publicKey := new(bn256.G2).ScalarBaseMult(big.NewInt(2))
	marshalledPublicKey := publicKey.Marshal()

	var tests = map[string]struct {
		disqualifiedMemberIndexes []group.MemberIndex
		inactiveMemberIndexes     []group.MemberIndex
		gjkrResult                *gjkr.Result
		expectedResult            *beaconchain.DKGResult
	}{
		"group public key not provided, DQ and IA empty": {
			disqualifiedMemberIndexes: []group.MemberIndex{},
			inactiveMemberIndexes:     []group.MemberIndex{},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: nil,
				Group:          group.NewGroup(32, 64),
			},
			expectedResult: &beaconchain.DKGResult{
				GroupPublicKey: []byte{},
				Misbehaved:     []byte{},
			},
		},
		"group public key provided, DQ and IA empty": {
			disqualifiedMemberIndexes: []group.MemberIndex{},
			inactiveMemberIndexes:     []group.MemberIndex{},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewGroup(32, 64),
			},
			expectedResult: &beaconchain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Misbehaved:     []byte{},
			},
		},
		"group public key provided, both DQ and IA non-empty": {
			disqualifiedMemberIndexes: []group.MemberIndex{1, 4, 3, 50},
			inactiveMemberIndexes:     []group.MemberIndex{5, 3, 50},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewGroup(32, 64),
			},
			expectedResult: &beaconchain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Misbehaved:     []byte{0x01, 0x03, 0x04, 0x05, 0x32},
			},
		},
		"group public key provided, DQ empty, IA non-empty": {
			disqualifiedMemberIndexes: []group.MemberIndex{},
			inactiveMemberIndexes:     []group.MemberIndex{5},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewGroup(32, 64),
			},
			expectedResult: &beaconchain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Misbehaved:     []byte{0x05},
			},
		},
		"group public key provided, DQ non-empty, IA empty": {
			disqualifiedMemberIndexes: []group.MemberIndex{60, 1, 5},
			inactiveMemberIndexes:     []group.MemberIndex{},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewGroup(32, 64),
			},
			expectedResult: &beaconchain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Misbehaved:     []byte{0x01, 0x05, 0x3C},
			},
		},
	}
	for _, test := range tests {
		for _, disqualifiedMember := range test.disqualifiedMemberIndexes {
			test.gjkrResult.Group.MarkMemberAsDisqualified(disqualifiedMember)
		}

		for _, inactiveMember := range test.inactiveMemberIndexes {
			test.gjkrResult.Group.MarkMemberAsInactive(inactiveMember)
		}

		convertedResult := convertGjkrResult(test.gjkrResult)

		if !test.expectedResult.Equals(convertedResult) {
			t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedResult, convertedResult)
		}
	}
}
