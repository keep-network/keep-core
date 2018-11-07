package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestResolveSecretSharesAccusations(t *testing.T) {
	threshold := 3
	groupSize := 5

	members, err := initializeSharesJustifyingMemberGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}
	member := members[1]

	var tests = map[string]struct {
		senderID          int
		accusedID         int
		modifyShareS      func(shareS *big.Int) *big.Int
		modifyShareT      func(shareT *big.Int) *big.Int
		modifyCommitments func(commitments []*big.Int) []*big.Int
		expectedResult    int
		expectedError     error
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
		"incorrect shareT - accused member is punished": {
			senderID:  3,
			accusedID: 4,
			modifyShareT: func(shareT *big.Int) *big.Int {
				return new(big.Int).Sub(shareT, big.NewInt(13))
			},
			expectedResult: 4,
		},
		"incorrect commitments - accused member is punished": {
			senderID:  3,
			accusedID: 4,
			modifyCommitments: func(commitments []*big.Int) []*big.Int {
				newCommitments := make([]*big.Int, len(commitments))
				for i := range newCommitments {
					newCommitments[i] = big.NewInt(int64(990 + i))
				}
				return newCommitments
			},
			expectedResult: 4,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			err := setupSharesAndCommitments(members, threshold)
			if err != nil {
				t.Fatalf("unexpected error [%s]", err)
			}

			sender := findSharesJustifyingMemberByID(members, test.senderID)
			revealedShareS := sender.receivedValidSharesS[test.accusedID]
			revealedShareT := sender.receivedValidSharesT[test.accusedID]

			if test.modifyShareS != nil {
				revealedShareS = test.modifyShareS(revealedShareS)
			}

			if test.modifyShareT != nil {
				revealedShareT = test.modifyShareT(revealedShareT)
			}

			if test.modifyCommitments != nil {
				member.receivedValidPeerCommitments[test.accusedID] =
					test.modifyCommitments(member.receivedValidPeerCommitments[test.accusedID])
			}

			result, err := member.ResolveSecretSharesAccusations(
				test.senderID,
				test.accusedID,
				revealedShareS,
				revealedShareT,
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

func TestResolvePublicCoefficientsAccusations(t *testing.T) {
	threshold := 3
	groupSize := 5

	members, err := initializeCoefficientsJustifyingMemberGroup(threshold, groupSize)
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
			sender := findCoefficientsJustifyingMemberByID(members, test.senderID)
			revealedShareS := sender.receivedValidSharesS[test.accusedID]
			if test.modifyShareS != nil {
				revealedShareS = test.modifyShareS(revealedShareS)
			}
			if test.modifyPublicCoefficients != nil {
				member.receivedValidPeerCoefficients[test.accusedID] = test.modifyPublicCoefficients(member.receivedValidPeerCoefficients[test.accusedID])
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

// setupSharesAndCommitments simulates shares calculation and commitments sharing
// betwen members. It generates coefficients for each group member, calculates
// commitments and shares for each peer member individually. At the end it stores
// values for each member just like they would be received from peers.
func setupSharesAndCommitments(members []*SharesJustifyingMember, threshold int) error {
	groupSize := len(members)

	// Maps which will keep coefficients and commitments of all group members,
	// with members IDs as keys.
	groupCoefficientsA := make(map[int][]*big.Int, groupSize)
	groupCoefficientsB := make(map[int][]*big.Int, groupSize)
	groupCommitments := make(map[int][]*big.Int, groupSize)

	// Generate threshold+1 coefficients and commitments for each group member.
	for _, m := range members {
		memberCoefficientsA, err := generatePolynomial(threshold, m.protocolConfig)
		if err != nil {
			return fmt.Errorf("polynomial generation failed [%s]", err)
		}
		memberCoefficientsB, err := generatePolynomial(threshold, m.protocolConfig)
		if err != nil {
			return fmt.Errorf("polynomial generation failed [%s]", err)
		}

		commitments := make([]*big.Int, threshold+1)
		for k := range memberCoefficientsA {

			commitments[k] = m.vss.CalculateCommitment(
				memberCoefficientsA[k],
				memberCoefficientsB[k],
				m.protocolConfig.P,
			)
		}
		// Store generated values in maps.
		groupCoefficientsA[m.ID] = memberCoefficientsA
		groupCoefficientsB[m.ID] = memberCoefficientsB
		groupCommitments[m.ID] = commitments
	}
	// Simulate phase where members are calculating shares individually for each
	// peer member and store received shares and commitments from peers.
	for _, m := range members {
		for _, p := range members {
			if m.ID != p.ID {
				p.receivedValidSharesS[m.ID] = evaluateMemberShare(p.ID, groupCoefficientsA[m.ID])
				p.receivedValidSharesT[m.ID] = evaluateMemberShare(p.ID, groupCoefficientsB[m.ID])

				p.receivedValidPeerCommitments[m.ID] = groupCommitments[m.ID]
			}
		}
	}
	return nil
}

// setupPublicCoefficients simulates public coefficients calculation and sharing
// between members. It expects secret coefficients to be already stored in
// secretCoefficients field for each group member. At the end it stores
// values for each member just like they would be received from peers.
func setupPublicCoefficients(members []*CoefficientsJustifyingMember) {
	// Calculate public coefficients for each group member.
	for _, m := range members {
		m.CalculatePublicCoefficients()
	}
	// Simulate phase where members store received public coefficients from peers.
	for _, m := range members {
		for _, p := range members {
			if m.ID != p.ID {
				m.receivedValidPeerCoefficients[p.ID] = p.publicCoefficients
			}
		}
	}
}

func findSharesJustifyingMemberByID(members []*SharesJustifyingMember, id int) *SharesJustifyingMember {
	for _, m := range members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

func findCoefficientsJustifyingMemberByID(
	members []*CoefficientsJustifyingMember,
	id int,
) *CoefficientsJustifyingMember {
	for _, m := range members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

func initializeSharesJustifyingMemberGroup(threshold, groupSize int) ([]*SharesJustifyingMember, error) {
	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var sharesJustifyingMember []*SharesJustifyingMember
	for _, jm := range committingMembers {
		sharesJustifyingMember = append(sharesJustifyingMember, &SharesJustifyingMember{
			CommittingMember: jm,
		})
	}

	return sharesJustifyingMember, nil
}

func initializeCoefficientsJustifyingMemberGroup(threshold, groupSize int) ([]*CoefficientsJustifyingMember, error) {
	sharingMembers, err := initializeSharingMembersGroup(threshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var coefficientsJustifyingMember []*CoefficientsJustifyingMember
	for _, sm := range sharingMembers {
		coefficientsJustifyingMember = append(coefficientsJustifyingMember,
			&CoefficientsJustifyingMember{
				SharingMember: sm,
			})
	}

	return coefficientsJustifyingMember, nil
}
