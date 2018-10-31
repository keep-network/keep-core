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
	var commitmentMessages []*MemberCommitmentsMessage
	for _, member := range committingMembers {
		sharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
		if err != nil {
			t.Fatalf("shares and commitments calculation failed [%s]", err)
		}
		sharesMessages = append(sharesMessages, sharesMessage...)
		commitmentMessages = append(commitmentMessages, commitmentsMessage)
	}

	for i := range committingMembers {
		committingMember := committingMembers[i]

		accusedMessage, err := committingMember.VerifyReceivedSharesAndCommitmentsMessages(
			filterPeerSharesMessage(sharesMessages, committingMember.ID),
			filterMemberCommitmentsMessages(commitmentMessages, committingMember.ID),
		)
		if err != nil {
			t.Fatalf("shares and commitments verification failed [%s]", err)
		}

		if len(accusedMessage.accusedIDs) > 0 {
			t.Fatalf("found accused members but was not expecting to")
		}
	}
}
