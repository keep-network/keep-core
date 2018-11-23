package gjkr

import (
	"math/big"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	threshold := 3
	groupSize := 5

	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize, nil)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	for _, member := range committingMembers {
		err := member.CalculateMembersSharesAndCommitments()
		if err != nil {
			t.Fatalf("shares and commitments calculation failed [%s]", err)
		}
	}

	receivedPeerSharesAndCommitments := make(map[int][]*SharesAndCommitments)

	for _, member := range committingMembers {
		for _, peer := range committingMembers {
			if member.ID != peer.ID {
				receivedPeerSharesAndCommitments[member.ID] = append(
					receivedPeerSharesAndCommitments[member.ID],
					&SharesAndCommitments{
						peerID:      peer.ID,
						shareS:      peer.evaluatedSecretSharesS[member.ID],
						shareT:      peer.evaluatedSecretSharesT[member.ID],
						commitments: peer.commitments,
					})
			}
		}
	}

	for _, member := range committingMembers {
		for _, sharesAndCommitments := range receivedPeerSharesAndCommitments[member.ID] {
			member.VerifyReceivedSharesAndCommitments(
				sharesAndCommitments.peerID,
				sharesAndCommitments.shareS, sharesAndCommitments.shareT,
				sharesAndCommitments.commitments,
			)
		}

		if len(member.accusedMembersIDs) > 0 {
			t.Fatalf("\nexpected: 0 accusations\nactual:   %d\n",
				member.accusedMembersIDs,
			)
		}

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

	publicKeySharePointsMessages := make([]*MemberPublicKeySharePointsMessage, groupSize)
	for i, member := range sharingMembers {
		publicKeySharePointsMessages[i] = member.CalculatePublicKeySharePoints()
	}

	for _, member := range sharingMembers {
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

	var combiningMembers []*CombiningMember
	for _, sm := range sharingMembers {
		// TODO: Handle transition from SharingMember to ReconstructingMember in Next() function
		combiningMembers = append(combiningMembers, &CombiningMember{
			ReconstructingMember: &ReconstructingMember{
				PointsJustifyingMember: &PointsJustifyingMember{
					SharingMember: sm,
				},
			},
		})
	}

	for _, member := range combiningMembers {
		member.CombineGroupPublicKey()
	}

}

type SharesAndCommitments struct {
	peerID         int
	shareS, shareT *big.Int
	commitments    []*big.Int
}
