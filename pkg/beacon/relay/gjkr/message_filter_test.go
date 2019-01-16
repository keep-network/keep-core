package gjkr

import (
	"reflect"
	"testing"
)

func TestExcludeInactiveMembers(t *testing.T) {
	var tests = map[string]struct {
		selfMemberID             MemberID
		groupMembers             []MemberID
		messageSenderIDs         []MemberID
		expectedOperatingMembers []MemberID
	}{
		"all other members active": {
			selfMemberID:             4,
			groupMembers:             []MemberID{3, 2, 4, 5, 1, 9},
			messageSenderIDs:         []MemberID{3, 2, 5, 9, 1},
			expectedOperatingMembers: []MemberID{3, 2, 4, 5, 1, 9},
		},
		"all other members inactive": {
			selfMemberID:             9,
			groupMembers:             []MemberID{9, 1, 2, 3},
			messageSenderIDs:         []MemberID{},
			expectedOperatingMembers: []MemberID{9},
		},
		"some members inactive": {
			selfMemberID:             3,
			groupMembers:             []MemberID{3, 4, 5, 1, 2, 8},
			messageSenderIDs:         []MemberID{1, 4, 2},
			expectedOperatingMembers: []MemberID{3, 4, 1, 2},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			group := &Group{
				memberIDs: test.groupMembers,
			}

			filter := &messageFilter{
				selfMemberID:       test.selfMemberID,
				group:              group,
				phaseActiveMembers: make([]MemberID, 0),
			}

			for _, member := range test.messageSenderIDs {
				filter.markMemberAsActive(member)
			}

			filter.flushInactiveMembers()

			actual := filter.group.OperatingMemberIDs()
			expected := test.expectedOperatingMembers

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf(
					"unexpected active members\n[%v]\n[%v]", actual, expected)
			}
		})
	}
}
