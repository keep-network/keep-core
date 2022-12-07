package group

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestFilterInactiveMembers(t *testing.T) {
	var tests = map[string]struct {
		selfMemberIndex          MemberIndex
		groupMembers             []MemberIndex
		messageSenderIndexes     []MemberIndex
		expectedOperatingMembers []MemberIndex
	}{
		"all other members active": {
			selfMemberIndex:          4,
			groupMembers:             []MemberIndex{3, 2, 4, 5, 1, 9},
			messageSenderIndexes:     []MemberIndex{3, 2, 5, 9, 1},
			expectedOperatingMembers: []MemberIndex{3, 2, 4, 5, 1, 9},
		},
		"all other members inactive": {
			selfMemberIndex:          9,
			groupMembers:             []MemberIndex{9, 1, 2, 3},
			messageSenderIndexes:     []MemberIndex{},
			expectedOperatingMembers: []MemberIndex{9},
		},
		"some members inactive": {
			selfMemberIndex:          3,
			groupMembers:             []MemberIndex{3, 4, 5, 1, 2, 8},
			messageSenderIndexes:     []MemberIndex{1, 4, 2},
			expectedOperatingMembers: []MemberIndex{3, 4, 1, 2},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			group := &Group{
				memberIndexes: test.groupMembers,
			}

			filter := &InactiveMemberFilter{
				logger:             &testutils.MockLogger{},
				selfMemberID:       test.selfMemberIndex,
				group:              group,
				phaseActiveMembers: make([]MemberIndex, 0),
			}

			for _, member := range test.messageSenderIndexes {
				filter.MarkMemberAsActive(member)
			}

			filter.FlushInactiveMembers()

			actual := filter.group.OperatingMemberIndexes()
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
