package tbtc

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"reflect"
	"testing"
)

func TestFinalSigningGroup(t *testing.T) {
	beaconConfig := &ChainConfig{
		GroupSize:       5,
		HonestThreshold: 3,
	}

	selectedOperators := []chain.Address{
		"0xAA",
		"0xBB",
		"0xCC",
		"0xDD",
		"0xEE",
	}

	var tests = map[string]struct {
		selectedOperators                  []chain.Address
		operatingMembersIndexes            []group.MemberIndex
		expectedSigningGroupOperators      []chain.Address
		expectedSigningGroupMembersIndexes map[group.MemberIndex]group.MemberIndex
		expectedError                      error
	}{
		"selected operators count not equal to the group size": {
			selectedOperators:       selectedOperators[:4],
			operatingMembersIndexes: []group.MemberIndex{1, 2, 3, 4, 5},
			expectedError:           fmt.Errorf("invalid input parameters"),
		},
		"all selected operators are operating": {
			selectedOperators:                  selectedOperators,
			operatingMembersIndexes:            []group.MemberIndex{5, 4, 3, 2, 1},
			expectedSigningGroupOperators:      selectedOperators,
			expectedSigningGroupMembersIndexes: map[group.MemberIndex]group.MemberIndex{1: 1, 2: 2, 3: 3, 4: 4, 5: 5},
		},
		"honest majority of selected operators are operating": {
			selectedOperators:                  selectedOperators,
			operatingMembersIndexes:            []group.MemberIndex{5, 1, 3},
			expectedSigningGroupOperators:      []chain.Address{"0xAA", "0xCC", "0xEE"},
			expectedSigningGroupMembersIndexes: map[group.MemberIndex]group.MemberIndex{1: 1, 3: 2, 5: 3},
		},
		"less than honest majority of selected operators are operating": {
			selectedOperators:       selectedOperators,
			operatingMembersIndexes: []group.MemberIndex{5, 1},
			expectedError:           fmt.Errorf("invalid input parameters"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualGroupOperators, actualSigningGroupMembersIndexes, err :=
				finalSigningGroup(
					test.selectedOperators,
					test.operatingMembersIndexes,
					beaconConfig,
				)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Errorf(
					"unexpected error\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if !reflect.DeepEqual(
				test.expectedSigningGroupOperators,
				actualGroupOperators,
			) {
				t.Errorf(
					"unexpected group operators\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedSigningGroupOperators,
					actualGroupOperators,
				)
			}

			if !reflect.DeepEqual(
				test.expectedSigningGroupMembersIndexes,
				actualSigningGroupMembersIndexes,
			) {
				t.Errorf(
					"unexpected group members indexes\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedSigningGroupMembersIndexes,
					actualSigningGroupMembersIndexes,
				)
			}
		})
	}
}
