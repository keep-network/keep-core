package gjkr

import (
	"math/big"
	"reflect"
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
		disqualifiedMemberIDs []int
		inactiveMemberIDs     []int
		expectedResult        *Result
	}{
		"no disqualified or inactive members - success": {
			expectedResult: &Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   nil,
				Inactive:       nil,
			},
		},
		"one disqualified member - success": {
			disqualifiedMemberIDs: []int{2},
			expectedResult: &Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   []int{2},
				Inactive:       nil,
			},
		},
		"two inactive members - success": {
			inactiveMemberIDs: []int{3, 7},
			expectedResult: &Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   nil,
				Inactive:       []int{3, 7},
			},
		},
		"more than half of threshold disqualified and inactive members - failure": {
			disqualifiedMemberIDs: []int{2},
			inactiveMemberIDs:     []int{3, 7},
			expectedResult: &Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   []int{2},
				// inactive member ids; this value is nil in the case of a failure, as
				// only disqualified members are slashed
				Inactive: nil,
			},
		},
		"more than half of threshold inactive members - failure": {
			inactiveMemberIDs: []int{3, 5, 7},
			expectedResult: &Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   nil,
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

				if !reflect.DeepEqual(test.expectedResult, resultToPublish) {
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
