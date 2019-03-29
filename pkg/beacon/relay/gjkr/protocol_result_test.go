package gjkr

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
)

func TestGenerateResult(t *testing.T) {
	threshold := 4
	groupSize := 8

	members, err := initializeFinalizingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		disqualifiedMemberIDs []member.Index
		inactiveMemberIDs     []member.Index
		expectedResult        *Result
	}{
		"no disqualified or inactive members - success": {
			expectedResult: &Result{
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   []member.Index{},
				Inactive:       []member.Index{},
			},
		},
		"one disqualified member - success": {
			disqualifiedMemberIDs: []member.Index{2},
			expectedResult: &Result{
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   []member.Index{2},
				Inactive:       []member.Index{},
			},
		},
		"two inactive members - success": {
			inactiveMemberIDs: []member.Index{3, 7},
			expectedResult: &Result{
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   []member.Index{},
				Inactive:       []member.Index{3, 7},
			},
		},
		"more than half of threshold disqualified and inactive members - failure": {
			disqualifiedMemberIDs: []member.Index{2},
			inactiveMemberIDs:     []member.Index{3, 7},
			expectedResult: &Result{
				GroupPublicKey: nil,
				Disqualified:   []member.Index{2},
				Inactive:       []member.Index{3, 7},
			},
		},
		"more than half of threshold inactive members - failure": {
			inactiveMemberIDs: []member.Index{3, 5, 7},
			expectedResult: &Result{
				GroupPublicKey: nil,
				Disqualified:   nil,
				Inactive:       []member.Index{3, 5, 7},
			},
		},
		"more than half of threshold disqualified members - failure": {
			disqualifiedMemberIDs: []member.Index{3, 5, 7},
			expectedResult: &Result{
				GroupPublicKey: nil,
				Disqualified:   []member.Index{3, 5, 7},
				Inactive:       []member.Index{},
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			for _, member := range members {
				member.group.disqualifiedMemberIDs = test.disqualifiedMemberIDs
				member.group.inactiveMemberIDs = test.inactiveMemberIDs

				resultToPublish := member.Result()

				if !test.expectedResult.Equals(resultToPublish) {
					t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, resultToPublish)
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
