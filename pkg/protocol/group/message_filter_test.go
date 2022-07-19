package group

import (
	"reflect"
	"testing"
)

func TestFilterInactiveMembers(t *testing.T) {
	var tests = map[string]struct {
		selfMemberID             MemberIndex
		groupMembers             []MemberIndex
		messageSenderIDs         []MemberIndex
		expectedOperatingMembers []MemberIndex
	}{
		"all other members active": {
			selfMemberID:             4,
			groupMembers:             []MemberIndex{3, 2, 4, 5, 1, 9},
			messageSenderIDs:         []MemberIndex{3, 2, 5, 9, 1},
			expectedOperatingMembers: []MemberIndex{3, 2, 4, 5, 1, 9},
		},
		"all other members inactive": {
			selfMemberID:             9,
			groupMembers:             []MemberIndex{9, 1, 2, 3},
			messageSenderIDs:         []MemberIndex{},
			expectedOperatingMembers: []MemberIndex{9},
		},
		"some members inactive": {
			selfMemberID:             3,
			groupMembers:             []MemberIndex{3, 4, 5, 1, 2, 8},
			messageSenderIDs:         []MemberIndex{1, 4, 2},
			expectedOperatingMembers: []MemberIndex{3, 4, 1, 2},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			group := &Group{
				memberIDs: test.groupMembers,
			}

			filter := &InactiveMemberFilter{
				selfMemberID:       test.selfMemberID,
				group:              group,
				phaseActiveMembers: make([]MemberIndex, 0),
			}

			for _, member := range test.messageSenderIDs {
				filter.MarkMemberAsActive(member)
			}

			filter.FlushInactiveMembers()

			actual := filter.group.OperatingMemberIDs()
			expected := test.expectedOperatingMembers

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf(
					"unexpected active members\nexpected: %v\nactual:   %v\n",
					expected,
					actual,
				)
			}
		})
	}
}
