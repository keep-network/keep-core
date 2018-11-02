package gjkr

import "testing"

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

	for i := range committingMembers {
		committingMember := committingMembers[i]

		accusedSecretSharesMessage, err := committingMember.VerifyReceivedSharesAndCommitmentsMessages(
			filterPeerSharesMessage(sharesMessages, committingMember.ID),
			filterMemberCommitmentsMessages(commitmentsMessages, committingMember.ID),
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

	var sharingMembers []*SharingMember
	for _, cm := range committingMembers {
		sharingMembers = append(sharingMembers, &SharingMember{
			CommittingMember: cm,
		})
	}

	sharingMember := sharingMembers[0]
	if len(sharingMember.receivedSharesS) != groupSize-1 {
		t.Fatalf("\nexpected: %d received shares\nactual:   %d\n",
			groupSize-1,
			len(sharingMember.receivedSharesS),
		)
	}

	for _, member := range sharingMembers {
		member.CombineMemberShares()
	}

	publicCoefficientsMessages := make([]*MemberPublicCoefficientsMessage, groupSize)
	for i, member := range sharingMembers {
		publicCoefficientsMessages[i] = member.CalculatePublicCoefficients()
	}

	for i := range sharingMembers {
		member := sharingMembers[i]

		accusedCoefficientsMessage, err := member.VerifyPublicCoefficients(
			filterMemberPublicCoefficientsMessages(publicCoefficientsMessages, member.ID),
		)
		if err != nil {
			t.Fatalf("public coefficients verification failed [%s]", err)
		}
		if len(accusedCoefficientsMessage.accusedIDs) > 0 {
			t.Fatalf("\nexpected: 0 accusations\nactual:   %d\n",
				accusedCoefficientsMessage.accusedIDs,
			)
		}
	}
}
