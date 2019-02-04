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

func TestMarkMemberAsDisqualified(t *testing.T) {
	var tests = map[string]struct {
		initialMembers              []MemberID
		updateFunc                  func(g *Group)
		expectedDisqualifiedMembers []MemberID
		expectedInactiveMembers     []MemberID
	}{
		"mark member as disqualified": {
			initialMembers: []MemberID{19, 11, 31, 33},
			updateFunc: func(g *Group) {
				g.MarkMemberAsDisqualified(19)
			},
			expectedDisqualifiedMembers: []MemberID{19},
			expectedInactiveMembers:     []MemberID{},
		},
		"mark member as disqualified twice": {
			initialMembers: []MemberID{19, 11, 31, 33},
			updateFunc: func(g *Group) {
				g.MarkMemberAsDisqualified(11)
				g.MarkMemberAsDisqualified(11)
			},
			expectedDisqualifiedMembers: []MemberID{11},
			expectedInactiveMembers:     []MemberID{},
		},
		"mark member from out of the group as disqualified": {
			initialMembers: []MemberID{19, 11, 31, 33},
			updateFunc: func(g *Group) {
				g.MarkMemberAsDisqualified(88)
			},
			expectedDisqualifiedMembers: []MemberID{},
			expectedInactiveMembers:     []MemberID{},
		},
		"mark all members as disqualified": {
			initialMembers: []MemberID{11, 12, 13},
			updateFunc: func(g *Group) {
				g.MarkMemberAsDisqualified(11)
				g.MarkMemberAsDisqualified(13)
				g.MarkMemberAsDisqualified(12)
			},
			expectedDisqualifiedMembers: []MemberID{11, 13, 12},
			expectedInactiveMembers:     []MemberID{},
		},
		"mark member as inactive": {
			initialMembers: []MemberID{19, 11, 31, 33},
			updateFunc: func(g *Group) {
				g.MarkMemberAsInactive(31)
			},
			expectedDisqualifiedMembers: []MemberID{},
			expectedInactiveMembers:     []MemberID{31},
		},
		"mark member as inactive twice": {
			initialMembers: []MemberID{19, 11, 31, 33},
			updateFunc: func(g *Group) {
				g.MarkMemberAsInactive(33)
				g.MarkMemberAsInactive(33)
			},
			expectedDisqualifiedMembers: []MemberID{},
			expectedInactiveMembers:     []MemberID{33},
		},
		"mark member from out of the group as inactive": {
			initialMembers: []MemberID{19, 11, 31, 33},
			updateFunc: func(g *Group) {
				g.MarkMemberAsInactive(99)
			},
			expectedDisqualifiedMembers: []MemberID{},
			expectedInactiveMembers:     []MemberID{},
		},
		"mark all members as inactive": {
			initialMembers: []MemberID{19, 18, 17, 16},
			updateFunc: func(g *Group) {
				g.MarkMemberAsInactive(17)
				g.MarkMemberAsInactive(19)
				g.MarkMemberAsInactive(16)
				g.MarkMemberAsInactive(18)
			},
			expectedDisqualifiedMembers: []MemberID{},
			expectedInactiveMembers:     []MemberID{17, 19, 16, 18},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			group := &Group{
				memberIDs:             test.initialMembers,
				disqualifiedMemberIDs: []MemberID{},
				inactiveMemberIDs:     []MemberID{},
			}

			if test.updateFunc != nil {
				test.updateFunc(group)
			}

			if !reflect.DeepEqual(
				test.expectedDisqualifiedMembers,
				group.disqualifiedMemberIDs,
			) {
				t.Fatalf(
					"unexpected list of disqualified members\nexpected: %v\nactual:   %v\n",
					test.expectedDisqualifiedMembers,
					group.disqualifiedMemberIDs,
				)
			}

			if !reflect.DeepEqual(
				test.expectedInactiveMembers,
				group.inactiveMemberIDs,
			) {
				t.Fatalf(
					"unexpected list of inactive members\nexpected: %v\nactual:   %v\n",
					test.expectedInactiveMembers,
					group.inactiveMemberIDs,
				)
			}
		})
	}
}

func TestIsDisqualified(t *testing.T) {
	group := &Group{
		memberIDs: []MemberID{19, 11, 31, 33},
	}

	if group.isDisqualified(19) {
		t.Errorf("member should not be disqualified at this point")
	}

	group.MarkMemberAsDisqualified(19)

	if !group.isDisqualified(19) {
		t.Errorf("member should be disqualified at this point")
	}
}

func TestIsInactive(t *testing.T) {
	group := &Group{
		memberIDs: []MemberID{19, 11, 31, 33},
	}

	if group.isInactive(31) {
		t.Errorf("member should ne be inactive at this point")
	}

	group.MarkMemberAsInactive(31)

	if !group.isInactive(31) {
		t.Errorf("member should be inactive at this point")
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
				g.MarkMemberAsDisqualified(98)
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
				g.MarkMemberAsDisqualified(19)
				g.MarkMemberAsInactive(33)
			},
			expectedOperatingMembers: []MemberID{11, 31},
		},
		"all but one inactive": {
			initialMembers: []MemberID{28, 19, 29},
			updateFunc: func(g *Group) {
				g.MarkMemberAsDisqualified(19)
				g.MarkMemberAsDisqualified(29)
			},
			expectedOperatingMembers: []MemberID{28},
		},
		"all but one disqualified": {
			initialMembers: []MemberID{92, 11, 20},
			updateFunc: func(g *Group) {
				g.MarkMemberAsDisqualified(92)
				g.MarkMemberAsDisqualified(11)
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
					"unexpected list of operating members\nexpected: %v\nactual:   %v\n",
					test.expectedOperatingMembers,
					operatingMembers,
				)
			}

		})
	}
}
