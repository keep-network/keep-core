package gjkr

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRegisterMemberIDWithInvalidID(t *testing.T) {
	expectedError := fmt.Errorf("cannot register member ID in the group [member ID must be >= 1]")

	group := &Group{}
	err := group.RegisterMemberID(MemberID(0))

	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedError, err)
	}
}

func TestDisqualifyMemberID(t *testing.T) {
	group := &Group{
		memberIDs: []MemberID{93, 31, 32},
	}

	if group.isDisqualified(32) {
		t.Errorf("member 32 should not be disqualified yet")
	}

	group.DisqualifyMember(32)

	if !group.isDisqualified(32) {
		t.Errorf("member 32 should be disqualified")
	}
}

func TestMarkMemberAsInactive(t *testing.T) {
	group := &Group{
		memberIDs: []MemberID{18, 29, 19},
	}

	if group.isInactive(29) {
		t.Errorf("member 29 should not be marked as inactive yet")
	}

	group.MarkMemberAsInactive(29)

	if !group.isInactive(29) {
		t.Errorf("member 29 should be marked as inactive")
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
