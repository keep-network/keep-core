package gjkr

import (
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

func TestGenerateResult(t *testing.T) {
	threshold := 4
	groupSize := 8

	var tests = map[string]struct {
		disqualifiedMemberIDs []group.MemberIndex
		inactiveMemberIDs     []group.MemberIndex
		expectedResult        func(groupPublicKey *bn256.G2) *Result
	}{
		"no disqualified or inactive members - success": {
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: groupPublicKey,
					Disqualified:   []group.MemberIndex{},
					Inactive:       []group.MemberIndex{},
				}
			},
		},
		"one disqualified member - success": {
			disqualifiedMemberIDs: []group.MemberIndex{2},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: groupPublicKey,
					Disqualified:   []group.MemberIndex{2},
					Inactive:       []group.MemberIndex{},
				}
			},
		},
		"two inactive members - success": {
			inactiveMemberIDs: []group.MemberIndex{3, 7},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: groupPublicKey,
					Disqualified:   []group.MemberIndex{},
					Inactive:       []group.MemberIndex{3, 7},
				}
			},
		},
		"more than half of threshold disqualified and inactive members - failure": {
			disqualifiedMemberIDs: []group.MemberIndex{2},
			inactiveMemberIDs:     []group.MemberIndex{3, 7},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: nil,
					Disqualified:   []group.MemberIndex{2},
					Inactive:       []group.MemberIndex{3, 7},
				}
			},
		},
		"more than half of threshold inactive members - failure": {
			inactiveMemberIDs: []group.MemberIndex{3, 5, 7},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: nil,
					Disqualified:   nil,
					Inactive:       []group.MemberIndex{3, 5, 7},
				}
			},
		},
		"more than half of threshold disqualified members - failure": {
			disqualifiedMemberIDs: []group.MemberIndex{3, 5, 7},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: nil,
					Disqualified:   []group.MemberIndex{3, 5, 7},
					Inactive:       []group.MemberIndex{},
				}
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializeFinalizingMembersGroup(threshold, groupSize)
			if err != nil {
				t.Fatal(err)
			}

			groupPubKey := members[0].groupPublicKey
			expectedResult := test.expectedResult(groupPubKey)

			for _, member := range members {
				for _, dq := range test.disqualifiedMemberIDs {
					member.group.MarkMemberAsDisqualified(dq)
				}
				for _, ia := range test.inactiveMemberIDs {
					member.group.MarkMemberAsInactive(ia)
				}

				resultToPublish := member.Result()

				if !expectedResult.Equals(resultToPublish) {
					t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedResult, resultToPublish)
				}
			}
		})
	}
}

func initializeFinalizingMembersGroup(threshold, groupSize int) ([]*FinalizingMember, error) {
	combiningMembers, err := initializeCombiningMembersGroup(threshold, groupSize)
	if err != nil {
		return nil, err
	}

	var finalizingMembers []*FinalizingMember
	for _, cm := range combiningMembers {
		finalizingMembers = append(finalizingMembers, cm.InitializeFinalization())
	}
	return finalizingMembers, nil
}
