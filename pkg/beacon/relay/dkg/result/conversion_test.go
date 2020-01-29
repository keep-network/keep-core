package result

import (
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

func TestConvertResult(t *testing.T) {
	groupSize := 5

	publicKey := new(bn256.G2).ScalarBaseMult(big.NewInt(2))
	marshalledPublicKey := publicKey.Marshal()

	var tests = map[string]struct {
		disqualifiedMemberIDs []group.MemberIndex
		inactiveMemberIDs     []group.MemberIndex
		gjkrResult            *gjkr.Result
		expectedResult        *relayChain.DKGResult
	}{
		"success: false, group public key: nil, DQ and IA: empty": {
			disqualifiedMemberIDs: []group.MemberIndex{},
			inactiveMemberIDs:     []group.MemberIndex{},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: nil,
				Group:          group.NewDkgGroup(3, 5),
			},
			expectedResult: &relayChain.DKGResult{
				GroupPublicKey: []byte{},
				Disqualified:   []byte{},
				Inactive:       []byte{},
			},
		},
		"success: true, group public key: provided, DQ and IA: provided": {
			disqualifiedMemberIDs: []group.MemberIndex{1, 3, 4},
			inactiveMemberIDs:     []group.MemberIndex{5},
			gjkrResult: &gjkr.Result{
				GroupPublicKey: publicKey,
				Group:          group.NewDkgGroup(3, 5),
			},
			expectedResult: &relayChain.DKGResult{
				GroupPublicKey: marshalledPublicKey,
				Disqualified:   []byte{0x01, 0x03, 0x04},
				Inactive:       []byte{0x05},
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

		convertedResult := convertResult(test.gjkrResult, groupSize)

		if !test.expectedResult.Equals(convertedResult) {
			t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedResult, convertedResult)
		}
	}
}
