package gjkr

import (
	"reflect"
	"testing"
)

func TestFilterInactiveMembers(t *testing.T) {
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

			filter := &inactiveMemberFilter{
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
					"unexpected active members\nexpected: %v\nactual:   %v\n",
					expected,
					actual,
				)
			}
		})
	}
}

func TestFilterSymmetricKeyGeneratingMembers(t *testing.T) {
	member := (&LocalMember{
		memberCore: &memberCore{
			ID: 13,
			group: &Group{
				memberIDs: []MemberID{11, 12, 13, 14, 15},
			},
		},
	}).InitializeEphemeralKeysGeneration().
		InitializeSymmetricKeyGeneration()

	messages := []*EphemeralPublicKeyMessage{
		&EphemeralPublicKeyMessage{senderID: 11},
		&EphemeralPublicKeyMessage{senderID: 14},
	}

	member.MarkInactiveMembers(messages)

	assertAcceptsFrom(member, 13, t) // should accept from self
	assertAcceptsFrom(member, 11, t)
	assertAcceptsFrom(member, 14, t)
	assertNotAcceptFrom(member, 12, t)
	assertNotAcceptFrom(member, 15, t)
}

func TestFilterCommitmentsVefiryingMembers(t *testing.T) {
	member := (&LocalMember{
		memberCore: &memberCore{
			ID: 93,
			group: &Group{
				memberIDs: []MemberID{91, 92, 93, 94, 95, 96},
			},
		},
	}).InitializeEphemeralKeysGeneration().
		InitializeSymmetricKeyGeneration().
		InitializeCommitting().
		InitializeCommitmentsVerification()

	sharesMessages := []*PeerSharesMessage{
		&PeerSharesMessage{senderID: 91},
		&PeerSharesMessage{senderID: 92},
		&PeerSharesMessage{senderID: 94},
	}

	commitmentsMessages := []*MemberCommitmentsMessage{
		&MemberCommitmentsMessage{senderID: 92},
		&MemberCommitmentsMessage{senderID: 94},
		&MemberCommitmentsMessage{senderID: 95},
	}

	member.MarkInactiveMembers(sharesMessages, commitmentsMessages)

	// should accept from self
	assertAcceptsFrom(member, 93, t)

	// 92 and 94 sent both shares message and commitments message
	assertAcceptsFrom(member, 92, t)
	assertAcceptsFrom(member, 94, t)

	// 95 did not send shares message
	assertNotAcceptFrom(member, 95, t)

	// 91 did not send commitments message
	assertNotAcceptFrom(member, 91, t)

	// 96 did not send shares message nor commitments message
	assertNotAcceptFrom(member, 96, t)
}

func TestFilterSharingMembers(t *testing.T) {
	member := (&LocalMember{
		memberCore: &memberCore{
			ID: 24,
			group: &Group{
				memberIDs: []MemberID{21, 22, 23, 24},
			},
		},
	}).InitializeEphemeralKeysGeneration().
		InitializeSymmetricKeyGeneration().
		InitializeCommitting().
		InitializeCommitmentsVerification().
		InitializeSharesJustification().
		InitializeQualified().
		InitializeSharing()

	messages := []*MemberPublicKeySharePointsMessage{
		&MemberPublicKeySharePointsMessage{senderID: 21},
		&MemberPublicKeySharePointsMessage{senderID: 23},
	}

	member.MarkInactiveMembers(messages)

	assertAcceptsFrom(member, 24, t) // should accept from self
	assertAcceptsFrom(member, 21, t)
	assertAcceptsFrom(member, 23, t)
	assertNotAcceptFrom(member, 22, t)
}

func TestFilterReconstructingMember(t *testing.T) {
	member := (&LocalMember{
		memberCore: &memberCore{
			ID: 44,
			group: &Group{
				memberIDs: []MemberID{41, 42, 43, 44},
			},
		},
	}).InitializeEphemeralKeysGeneration().
		InitializeSymmetricKeyGeneration().
		InitializeCommitting().
		InitializeCommitmentsVerification().
		InitializeSharesJustification().
		InitializeQualified().
		InitializeSharing().
		InitializePointsJustification().
		InitializeRevealing().
		InitializeReconstruction()

	messages := []*DisqualifiedEphemeralKeysMessage{
		&DisqualifiedEphemeralKeysMessage{senderID: 41},
	}

	member.MarkInactiveMembers(messages)

	assertAcceptsFrom(member, 44, t) // should accept from self
	assertAcceptsFrom(member, 41, t)
	assertNotAcceptFrom(member, 42, t)
	assertNotAcceptFrom(member, 43, t)
}

func assertAcceptsFrom(member MessageFiltering, senderID MemberID, t *testing.T) {
	if !member.IsSenderAccepted(senderID) {
		t.Errorf("member should accept messages from [%v]", senderID)
	}
}

func assertNotAcceptFrom(member MessageFiltering, senderID MemberID, t *testing.T) {
	if member.IsSenderAccepted(senderID) {
		t.Errorf("member should not accept messages from [%v]", senderID)
	}
}
