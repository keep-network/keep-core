package gjkr

import (
	"testing"
)

func TestGenerateResult(t *testing.T) {
	threshold := 4
	groupSize := 8

	members, err := initializePublishingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("%s", err)
	}

	var tests = map[string]struct {
		disqualifiedMemberIDs []MemberID
		inactiveMemberIDs     []MemberID
		expectedResult        *Result
	}{
		"no disqualified or inactive members - success": {
			expectedResult: &Result{
				Success:        true,
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   nil,
				Inactive:       nil,
			},
		},
		"one disqualified member - success": {
			disqualifiedMemberIDs: []MemberID{2},
			expectedResult: &Result{
				Success:        true,
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   []MemberID{2},
				Inactive:       nil,
			},
		},
		"two inactive members - success": {
			inactiveMemberIDs: []MemberID{3, 7},
			expectedResult: &Result{
				Success:        true,
				GroupPublicKey: members[0].groupPublicKey,
				Disqualified:   nil,
				Inactive:       []MemberID{3, 7},
			},
		},
		"more than half of threshold disqualified and inactive members - failure": {
			disqualifiedMemberIDs: []MemberID{2},
			inactiveMemberIDs:     []MemberID{3, 7},
			expectedResult: &Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   []MemberID{2},
				// inactive member ids; this value is nil in the case of a failure, as
				// only disqualified members are slashed
				Inactive: nil,
			},
		},
		"more than half of threshold inactive members - failure": {
			inactiveMemberIDs: []MemberID{3, 5, 7},
			expectedResult: &Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   nil,
				// inactive member ids; this value is nil in the case of a failure, as
				// only disqualified members are slashed
				Inactive: nil,
			},
		},
		"more than half of threshold disqualified members - failure": {
			disqualifiedMemberIDs: []MemberID{3, 5, 7},
			expectedResult: &Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   []MemberID{3, 5, 7},
				// inactive member ids; this value is nil in the case of a failure, as
				// only disqualified members are slashed
				Inactive: nil,
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

func initializePublishingMembersGroup(threshold, groupSize int) ([]*FinalizingMember, error) {
	combiningMembers, err := initializeCombiningMembersGroup(threshold, groupSize, nil)
	if err != nil {
		return nil, err
	}

	var publishingMembers []*FinalizingMember
	for _, cm := range combiningMembers {
		publishingMembers = append(publishingMembers, cm.InitializeFinalization())
	}
	return publishingMembers, nil
}
