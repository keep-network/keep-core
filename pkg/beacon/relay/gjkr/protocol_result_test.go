package gjkr

import (
	"testing"
)

func TestGenerateResult(t *testing.T) {
	threshold := 4
	groupSize := 8

	members, err := initializeFinalizingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		disqualifiedMemberIDs []MemberID
		inactiveMemberIDs     []MemberID
		expectedResult        *Result
	}{
		"no disqualified or inactive members - success": {
			expectedResult: &Result{
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   []MemberID{},
				Inactive:       []MemberID{},
			},
		},
		"one disqualified member - success": {
			disqualifiedMemberIDs: []MemberID{2},
			expectedResult: &Result{
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   []MemberID{2},
				Inactive:       []MemberID{},
			},
		},
		"two inactive members - success": {
			inactiveMemberIDs: []MemberID{3, 7},
			expectedResult: &Result{
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   []MemberID{},
				Inactive:       []MemberID{3, 7},
			},
		},
		"more than half of threshold disqualified and inactive members - failure": {
			disqualifiedMemberIDs: []MemberID{2},
			inactiveMemberIDs:     []MemberID{3, 7},
			expectedResult: &Result{
				GroupPublicKey: nil,
				Disqualified:   []MemberID{2},
				Inactive:       []MemberID{3, 7},
			},
		},
		"more than half of threshold inactive members - failure": {
			inactiveMemberIDs: []MemberID{3, 5, 7},
			expectedResult: &Result{
				GroupPublicKey: nil,
				Disqualified:   nil,
				Inactive:       []MemberID{3, 5, 7},
			},
		},
		"more than half of threshold disqualified members - failure": {
			disqualifiedMemberIDs: []MemberID{3, 5, 7},
			expectedResult: &Result{
				GroupPublicKey: nil,
				Disqualified:   []MemberID{3, 5, 7},
				Inactive:       []MemberID{},
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
