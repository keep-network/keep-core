package gjkr

import (
	"reflect"
	"testing"
)

func TestDisqualifyMemberID(t *testing.T) {
	group := &Group{}

	id1 := MemberID(123)
	id2 := MemberID(321)

	if len(group.disqualifiedMemberIDs) != 0 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			0,
			len(group.disqualifiedMemberIDs),
		)
	}

	// Disqualify a member.
	group.DisqualifyMember(id1)
	if len(group.disqualifiedMemberIDs) != 1 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			1,
			len(group.disqualifiedMemberIDs),
		)
	}

	// Disqualify the same member for a second time.
	group.DisqualifyMember(id1)
	if len(group.disqualifiedMemberIDs) != 1 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			1,
			len(group.disqualifiedMemberIDs),
		)
	}

	// Disqualify a next member.
	group.DisqualifyMember(id2)
	if len(group.disqualifiedMemberIDs) != 2 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			2,
			len(group.disqualifiedMemberIDs),
		)
	}
}

func TestOperatingMembers(t *testing.T) {
	var tests = map[string]struct {
		initialMembers           []MemberID
		updateFunc               func(g *Group)
		expectedOperatingMembers []MemberID
	}{
		"all members remain operating": {
			initialMembers:           []MemberID{10, 12, 33, 11},
			expectedOperatingMembers: []MemberID{10, 12, 33, 11},
		},
		"one member disqualified": {
			initialMembers: []MemberID{99, 98, 12, 33, 44},
			updateFunc: func(g *Group) {
				g.DisqualifyMember(98)
			},
			expectedOperatingMembers: []MemberID{99, 12, 33, 44},
		},
		"one member inactive": {
			initialMembers: []MemberID{38, 19, 39, 22, 11},
			updateFunc: func(g *Group) {
				g.MarkMemberAsInactive(11)
			},
			expectedOperatingMembers: []MemberID{38, 19, 39, 22},
		},
		"one member disqualified and one member inactive": {
			initialMembers: []MemberID{19, 11, 31, 33},
			updateFunc: func(g *Group) {
				g.DisqualifyMember(19)
				g.MarkMemberAsInactive(33)
			},
			expectedOperatingMembers: []MemberID{11, 31},
		},
		"all but one inactive": {
			initialMembers: []MemberID{28, 19, 29},
			updateFunc: func(g *Group) {
				g.DisqualifyMember(19)
				g.DisqualifyMember(29)
			},
			expectedOperatingMembers: []MemberID{28},
		},
		"all but one disqualified": {
			initialMembers: []MemberID{92, 11, 20},
			updateFunc: func(g *Group) {
				g.DisqualifyMember(92)
				g.DisqualifyMember(11)
			},
			expectedOperatingMembers: []MemberID{20},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			group := &Group{}
			group.memberIDs = test.initialMembers

			if test.updateFunc != nil {
				test.updateFunc(group)
			}

			operatingMembers := group.OperatingMemberIDs()
			if !reflect.DeepEqual(
				test.expectedOperatingMembers,
				operatingMembers,
			) {
				t.Fatalf(
					"unexpected list of operating members\n[%v]\n[%v]",
					test.expectedOperatingMembers,
					operatingMembers,
				)
			}

		})
	}
}
