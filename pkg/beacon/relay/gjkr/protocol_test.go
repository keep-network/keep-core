package gjkr

import (
	"math/big"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	threshold := 3
	groupSize := 5

	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	var sharesMessages []*PeerSharesMessage
	var commitmentsMessages []*MemberCommitmentsMessage
	for _, member := range committingMembers {
		sharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
		if err != nil {
			t.Fatalf("shares and commitments calculation failed [%s]", err)
		}
		sharesMessages = append(sharesMessages, sharesMessage...)
		commitmentsMessages = append(commitmentsMessages, commitmentsMessage)
	}

	for _, member := range committingMembers {
		accusedSecretSharesMessage, err := member.VerifyReceivedSharesAndCommitmentsMessages(
			filterPeerSharesMessage(sharesMessages, member.ID),
			filterMemberCommitmentsMessages(commitmentsMessages, member.ID),
		)
		if err != nil {
			t.Fatalf("shares and commitments verification failed [%s]", err)
		}

		if len(accusedSecretSharesMessage.accusedIDs) > 0 {
			t.Fatalf("\nexpected: 0 accusations\nactual:   %d\n",
				accusedSecretSharesMessage.accusedIDs,
			)
		}
	}

	var qualifiedMembers []*QualifiedMember
	// TODO: Handle transition from CommittingMember to SharingMember in Next() function
	for _, cm := range committingMembers {
		qualifiedMembers = append(qualifiedMembers, &QualifiedMember{
			SharesJustifyingMember: &SharesJustifyingMember{
				CommittingMember: cm,
			},
		})
	}

	for _, member := range qualifiedMembers {
		member.CombineMemberShares()
	}

	var sharingMembers []*SharingMember
	// TODO: Handle transition from CommittingMember to SharingMember in Next() function
	for _, qm := range qualifiedMembers {
		sharingMembers = append(sharingMembers, &SharingMember{
			QualifiedMember:                       qm,
			receivedValidPeerPublicKeySharePoints: make(map[int][]*big.Int, groupSize-1),
		})
	}

	for _, member := range sharingMembers {
		if len(member.receivedValidSharesS) != groupSize-1 {
			t.Fatalf("\nexpected: %d received shares S\nactual:   %d\n",
				groupSize-1,
				len(member.receivedValidSharesS),
			)
		}
		if len(member.receivedValidSharesT) != groupSize-1 {
			t.Fatalf("\nexpected: %d received shares T\nactual:   %d\n",
				groupSize-1,
				len(member.receivedValidSharesT),
			)
		}
		member.CombineMemberShares()
	}

	publicKeySharePointsMessages := make([]*MemberPublicKeySharePointsMessage, groupSize)
	for i, member := range sharingMembers {
		publicKeySharePointsMessages[i] = member.CalculatePublicKeySharePoints()

		member.receivedGroupPublicKeyShares = make(map[int]*big.Int, groupSize-1)
	}

	var reconstructingMembers []*ReconstructingMember
	for _, sm := range sharingMembers {
		reconstructingMembers = append(reconstructingMembers, &ReconstructingMember{
			SharingMember: sm,
		})
	}

	for i := range reconstructingMembers {
		reconstructingMembers[i].CombineGroupPublicKey()
	}

	for i := range reconstructingMembers {
		member := reconstructingMembers[i]

		accusedPointsMessage, err := member.VerifyPublicKeySharePoints(
			filterMemberPublicKeySharePointsMessages(publicKeySharePointsMessages, member.ID),
		)
		if err != nil {
			t.Fatalf("public coefficients verification failed [%s]", err)
		}
		if len(accusedPointsMessage.accusedIDs) > 0 {
			t.Fatalf("\nexpected: 0 accusations\nactual:   %d\n",
				accusedPointsMessage.accusedIDs,
			)
		}
	}
}
