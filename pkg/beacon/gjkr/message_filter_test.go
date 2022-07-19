package gjkr

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/protocol/group"
)

func TestFilterSymmetricKeyGeneratingMembers(t *testing.T) {
	member := (&LocalMember{
		memberCore: &memberCore{
			ID:    13,
			group: group.NewDkgGroup(8, 15),
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
			ID:    93,
			group: group.NewDkgGroup(49, 96),
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
			ID:    24,
			group: group.NewDkgGroup(13, 24),
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
			ID:    44,
			group: group.NewDkgGroup(23, 44),
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

	messages := []*MisbehavedEphemeralKeysMessage{
		{senderID: 41},
	}

	member.MarkInactiveMembers(messages)

	assertAcceptsFrom(member, 44, t) // should accept from self
	assertAcceptsFrom(member, 41, t)
	assertNotAcceptFrom(member, 42, t)
	assertNotAcceptFrom(member, 43, t)
}

func assertAcceptsFrom(member group.MessageFiltering, senderID group.MemberIndex, t *testing.T) {
	if !member.IsSenderAccepted(senderID) {
		t.Errorf("member should accept messages from [%v]", senderID)
	}
}

func assertNotAcceptFrom(member group.MessageFiltering, senderID group.MemberIndex, t *testing.T) {
	if member.IsSenderAccepted(senderID) {
		t.Errorf("member should not accept messages from [%v]", senderID)
	}
}
