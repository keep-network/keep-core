package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestResolvePublicCoefficientsAccusations(t *testing.T) {
	threshold := 3
	groupSize := 5

	members, err := initializeSharingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}
	member := members[1]

	var tests = map[string]struct {
		senderID                 int
		accusedID                int
		modifyShareS             func(shareS *big.Int) *big.Int
		modifyPublicCoefficients func(coefficients []*big.Int) []*big.Int
		expectedResult           int
		expectedError            error
	}{
		"false accusation - sender is punished": {
			senderID:       3,
			accusedID:      4,
			expectedResult: 3,
		},
		"current member as a sender - error returned": {
			senderID:       2,
			accusedID:      3,
			expectedResult: 0,
			expectedError:  fmt.Errorf("current member cannot be a part of a dispute"),
		},
		"current member as an accused - error returned": {
			senderID:       3,
			accusedID:      2,
			expectedResult: 0,
			expectedError:  fmt.Errorf("current member cannot be a part of a dispute"),
		},
		"incorrect shareS - accused member is punished": {
			senderID:  3,
			accusedID: 4,
			modifyShareS: func(shareS *big.Int) *big.Int {
				return new(big.Int).Sub(shareS, big.NewInt(1))
			},
			expectedResult: 4,
		},
		"incorrect commitments - accused member is punished": {
			senderID:  3,
			accusedID: 4,
			modifyPublicCoefficients: func(coefficients []*big.Int) []*big.Int {
				newCoefficients := make([]*big.Int, len(coefficients))
				for i := range newCoefficients {
					newCoefficients[i] = big.NewInt(int64(990 + i))
				}
				return newCoefficients
			},
			expectedResult: 4,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			setupPublicCoefficients(members)
			sender := findSharingMemberByID(members, test.senderID)
			revealedShareS := sender.receivedSharesS[test.accusedID]
			if test.modifyShareS != nil {
				revealedShareS = test.modifyShareS(revealedShareS)
			}
			if test.modifyPublicCoefficients != nil {
				member.receivedPublicCoefficients[test.accusedID] = test.modifyPublicCoefficients(member.receivedPublicCoefficients[test.accusedID])
			}
			result, err := member.ResolvePublicCoefficientsAccusations(
				test.senderID,
				test.accusedID,
				revealedShareS,
			)
			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf("\nexpected: %s\nactual:   %s\n", test.expectedError, err)
			}
			if result != test.expectedResult {
				t.Fatalf("\nexpected: %d\nactual:   %d\n", test.expectedResult, result)
			}
		})
	}
}

// setupPublicCoefficients simulates public coefficients calculation and sharing
// between members. It expects secret coefficients to be already stored in
// secretCoefficients field for each group member. At the end it stores
// values for each member just like they would be received from peers.
func setupPublicCoefficients(members []*SharingMember) {
	// Calculate public coefficients for each group member.
	for _, m := range members {
		m.CalculatePublicCoefficients()
	}
	// Simulate phase where members store received public coefficients from peers.
	for _, m := range members {
		for _, p := range members {
			if m.ID != p.ID {
				m.receivedPublicCoefficients[p.ID] = p.publicCoefficients
			}
		}
	}
}

func findSharingMemberByID(members []*SharingMember, id int) *SharingMember {
	for _, m := range members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

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
			QualifiedMember:            qm,
			receivedPublicCoefficients: make(map[int][]*big.Int, groupSize-1),
		})
	}

	sharingMember := sharingMembers[0]
	if len(sharingMember.receivedSharesS) != groupSize-1 {
		t.Fatalf("\nexpected: %d received shares\nactual:   %d\n",
			groupSize-1,
			len(sharingMember.receivedSharesS),
		)
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
