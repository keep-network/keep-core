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

	currentMemberID := 2 // i

	var tests = map[string]struct {
		accuserID         int // j
		accusedID         int // m
		modifyShareS      func(shareS *big.Int) *big.Int
		modifyShareT      func(shareT *big.Int) *big.Int
		modifyCommitments func(commitments []*big.Int) []*big.Int
		expectedResult    int
		expectedError     error
	}{
		"false accusation - accuser is punished": {
			accuserID:      3,
			accusedID:      4,
			expectedResult: 3,
		},
		"current member as an accuser - error returned": {
			accuserID:      currentMemberID,
			accusedID:      3,
			expectedResult: 0,
			expectedError:  fmt.Errorf("current member cannot be a part of a dispute"),
		},
		"current member as an accused - error returned": {
			accuserID:      3,
			accusedID:      currentMemberID,
			expectedResult: 0,
			expectedError:  fmt.Errorf("current member cannot be a part of a dispute"),
		},
		"incorrect shareS - accused member is punished": {
			accuserID: 3,
			accusedID: 4,
			modifyShareS: func(shareS *big.Int) *big.Int {
				return new(big.Int).Sub(shareS, big.NewInt(1))
			},
			expectedResult: 4,
		},
		"incorrect shareT - accused member is punished": {
			accuserID: 3,
			accusedID: 4,
			modifyShareT: func(shareT *big.Int) *big.Int {
				return new(big.Int).Sub(shareT, big.NewInt(13))
			},
			expectedResult: 4,
		},
		"incorrect commitments - accused member is punished": {
			accuserID: 3,
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
			members, err := initializeSharesJustifyingMemberGroup(threshold, groupSize, nil)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			member := findSharesJustifyingMemberByID(members, currentMemberID)

			// Simulate shares reveal by accuser `j`
			accuser := findSharesJustifyingMemberByID(members, test.accuserID)
			revealedShareS := accuser.receivedValidSharesS[test.accusedID]
			revealedShareT := accuser.receivedValidSharesT[test.accusedID]

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
				test.accuserID,
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

func TestResolvePublicKeySharePointsAccusations(t *testing.T) {
	threshold := 3
	groupSize := 5

	currentMemberID := 2 // i

	var tests = map[string]struct {
		accuserID                  int // j
		accusedID                  int // m
		modifyShareS               func(shareS *big.Int) *big.Int
		modifyPublicKeySharePoints func(coefficients []*big.Int) []*big.Int
		expectedResult             int
		expectedError              error
	}{
		"false accusation - sender is punished": {
			accuserID:      3,
			accusedID:      4,
			expectedResult: 3,
		},
		"current member as a sender - error returned": {
			accuserID:      currentMemberID,
			accusedID:      3,
			expectedResult: 0,
			expectedError:  fmt.Errorf("current member cannot be a part of a dispute"),
		},
		"current member as an accused - error returned": {
			accuserID:      3,
			accusedID:      currentMemberID,
			expectedResult: 0,
			expectedError:  fmt.Errorf("current member cannot be a part of a dispute"),
		},
		"incorrect shareS - accused member is punished": {
			accuserID: 3,
			accusedID: 4,
			modifyShareS: func(shareS *big.Int) *big.Int {
				return new(big.Int).Sub(shareS, big.NewInt(1))
			},
			expectedResult: 4,
		},
		"incorrect commitments - accused member is punished": {
			accuserID: 3,
			accusedID: 4,
			modifyPublicKeySharePoints: func(points []*big.Int) []*big.Int {
				newPoints := make([]*big.Int, len(points))
				for i := range newPoints {
					newPoints[i] = big.NewInt(int64(990 + i))
				}
				return newPoints
			},
			expectedResult: 4,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializePointsJustifyingMemberGroup(threshold, groupSize, nil)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			member := findCoefficientsJustifyingMemberByID(members, currentMemberID)

			sender := findCoefficientsJustifyingMemberByID(members, test.accuserID)
			revealedShareS := sender.receivedValidSharesS[test.accusedID]
			if test.modifyShareS != nil {
				revealedShareS = test.modifyShareS(revealedShareS)
			}
			if test.modifyPublicKeySharePoints != nil {
				member.receivedValidPeerPublicKeySharePoints[test.accusedID] =
					test.modifyPublicKeySharePoints(member.receivedValidPeerPublicKeySharePoints[test.accusedID])
			}
			result, err := member.ResolvePublicKeySharePointsAccusations(
				test.accuserID,
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

func findSharesJustifyingMemberByID(members []*SharesJustifyingMember, id int) *SharesJustifyingMember {
	for _, m := range members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

func findCoefficientsJustifyingMemberByID(
	members []*PointsJustifyingMember,
	id int,
) *PointsJustifyingMember {
	for _, m := range members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

// InitializeSharesJustifyingMemberGroup generates a group of members and simulates
// shares calculation and commitments sharing betwen members (Phases 3 and 4).
// It generates coefficients for each group member, calculates commitments and
// shares for each peer member individually. At the end it stores values for each
// member just like they would be received from peers.
func initializeSharesJustifyingMemberGroup(threshold, groupSize int, dkg *DKG) ([]*SharesJustifyingMember, error) {
	commitmentsVerifyingMembers, err := initializeCommitmentsVerifiyingMembersGroup(threshold, groupSize, dkg)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var sharesJustifyingMembers []*SharesJustifyingMember
	// TODO: Handle transition from CommittingMember to SharesJustifyingMember in Next() function
	for _, cvm := range commitmentsVerifyingMembers {
		sharesJustifyingMembers = append(sharesJustifyingMembers, &SharesJustifyingMember{
			CommitmentsVerifyingMember: cvm,
		})
	}

	// Maps which will keep coefficients and commitments of all group members,
	// with members IDs as keys.
	groupCoefficientsA := make(map[int][]*big.Int, groupSize)
	groupCoefficientsB := make(map[int][]*big.Int, groupSize)
	groupCommitments := make(map[int][]*big.Int, groupSize)

	// Generate threshold+1 coefficients and commitments for each group member.
	for _, m := range sharesJustifyingMembers {
		memberCoefficientsA, err := generatePolynomial(threshold, m.protocolConfig)
		if err != nil {
			return nil, fmt.Errorf("polynomial generation failed [%s]", err)
		}
		memberCoefficientsB, err := generatePolynomial(threshold, m.protocolConfig)
		if err != nil {
			return nil, fmt.Errorf("polynomial generation failed [%s]", err)
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
	for _, m := range sharesJustifyingMembers {
		for _, p := range sharesJustifyingMembers {
			if m.ID != p.ID {
				p.receivedValidSharesS[m.ID] = m.evaluateMemberShare(p.ID, groupCoefficientsA[m.ID])
				p.receivedValidSharesT[m.ID] = m.evaluateMemberShare(p.ID, groupCoefficientsB[m.ID])
				p.receivedValidPeerCommitments[m.ID] = groupCommitments[m.ID]
			}
		}
	}

	return sharesJustifyingMembers, nil
}

// initializePointsJustifyingMemberGroup generates a group of members and
// simulates public coefficients calculation and sharing between members
// (Phase 7 and 8). It expects secret coefficients to be already stored in
// secretCoefficients field for each group member. At the end it stores
// values for each member just like they would be received from peers.
func initializePointsJustifyingMemberGroup(
	threshold, groupSize int,
	dkg *DKG,
) ([]*PointsJustifyingMember, error) {
	sharingMembers, err := initializeSharingMembersGroup(threshold, groupSize, dkg)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var pointsJustifyingMembers []*PointsJustifyingMember
	// TODO: Handle transition from SharingMember to PointsJustifyingMember in Next() function
	for _, sm := range sharingMembers {
		pointsJustifyingMembers = append(pointsJustifyingMembers,
			&PointsJustifyingMember{
				SharingMember: sm,
			})
	}

	// Calculate public key share points for each group member (Phase 7).
	for _, m := range pointsJustifyingMembers {
		m.CalculatePublicKeySharePoints()
	}
	// Simulate phase where members store received public key share points from
	// peers (Phase 8).
	for _, m := range pointsJustifyingMembers {
		for _, p := range pointsJustifyingMembers {
			if m.ID != p.ID {
				m.receivedValidPeerPublicKeySharePoints[p.ID] = p.publicKeySharePoints
			}
		}
	}

	return pointsJustifyingMembers, nil
}
