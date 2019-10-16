package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

// TODO Add test with many messages from accusers and many accused in the message.
func TestResolveSecretSharesAccusations(t *testing.T) {
	dishonestThreshold := 2
	groupSize := 5

	currentMemberID := group.MemberIndex(2) // i

	var tests = map[string]struct {
		accuserID               group.MemberIndex // j
		accusedID               group.MemberIndex // m
		modifyEvidenceLog       func(evidenceLog evidenceLog) evidenceLog
		modifyShareS            func(shareS *big.Int) *big.Int
		modifyShareT            func(shareT *big.Int) *big.Int
		modifyCommitments       func(commitments []*bn256.G1) []*bn256.G1
		modifyAccusedPrivateKey func(symmetricKey *ephemeral.PrivateKey) *ephemeral.PrivateKey
		expectedResult          []group.MemberIndex
		expectedError           error
	}{
		"false accusation - accuser is disqualified": {
			accuserID:      3,
			accusedID:      4,
			expectedResult: []group.MemberIndex{3},
		},
		"current member as an accused - accuser is disqualified": {
			accuserID:      3,
			accusedID:      currentMemberID,
			expectedResult: []group.MemberIndex{3},
		},
		"incorrect shareS - accused member is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyShareS: func(shareS *big.Int) *big.Int {
				return new(big.Int).Sub(shareS, big.NewInt(1))
			},
			expectedResult: []group.MemberIndex{4},
		},
		"incorrect shareT - accused member is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyShareT: func(shareT *big.Int) *big.Int {
				return new(big.Int).Sub(shareT, big.NewInt(13))
			},
			expectedResult: []group.MemberIndex{4},
		},
		"incorrect commitments - accused member is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyCommitments: func(commitments []*bn256.G1) []*bn256.G1 {
				newCommitments := make([]*bn256.G1, len(commitments))
				for i := range newCommitments {
					newCommitments[i] = new(bn256.G1).ScalarBaseMult(
						big.NewInt(int64(i)),
					)
				}
				return newCommitments
			},
			expectedResult: []group.MemberIndex{4},
		},
		"no commitments - accused member is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyCommitments: func(commitments []*bn256.G1) []*bn256.G1 {
				return []*bn256.G1{}
			},
			expectedResult: []group.MemberIndex{4},
		},
		"incorrect accused private key - accuser is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyAccusedPrivateKey: func(symmetricKey *ephemeral.PrivateKey) *ephemeral.PrivateKey {
				return &ephemeral.PrivateKey{D: big.NewInt(12)}
			},
			expectedResult: []group.MemberIndex{3},
		},
		"inactive member as an accused (no EphemeralPublicKeyMessage sent) - " +
			"accuser is disqualified": {
			accuserID: 3,
			accusedID: 5,
			modifyEvidenceLog: func(evidenceLog evidenceLog) evidenceLog {
				if dkgEvidenceLog, ok := evidenceLog.(*dkgEvidenceLog); ok {
					dkgEvidenceLog.pubKeyMessageLog.removeMessage(
						group.MemberIndex(5),
					)
					return dkgEvidenceLog
				}
				return evidenceLog
			},
			expectedResult: []group.MemberIndex{3},
		},
		"inactive member as an accused (no PeerSharesMessage sent) - " +
			"accuser is disqualified": {
			accuserID: 3,
			accusedID: 5,
			modifyEvidenceLog: func(evidenceLog evidenceLog) evidenceLog {
				if dkgEvidenceLog, ok := evidenceLog.(*dkgEvidenceLog); ok {
					dkgEvidenceLog.peerSharesMessageLog.removeMessage(
						group.MemberIndex(5),
					)
					return dkgEvidenceLog
				}
				return evidenceLog
			},
			expectedResult: []group.MemberIndex{3},
		},
		"shares could not be decrypted - accused member is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyEvidenceLog: func(evidenceLog evidenceLog) evidenceLog {
				if dkgEvidenceLog, ok := evidenceLog.(*dkgEvidenceLog); ok {
					message := dkgEvidenceLog.peerSharesMessage(group.MemberIndex(4))
					message.shares[group.MemberIndex(3)] = &peerShares{
						[]byte{0x00},
						[]byte{0x00},
					}
					return dkgEvidenceLog
				}
				return evidenceLog
			},
			expectedResult: []group.MemberIndex{4},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializeSharesJustifyingMemberGroup(
				dishonestThreshold,
				groupSize,
			)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			justifyingMember := findSharesJustifyingMemberByID(members, currentMemberID)

			accuser := findSharesJustifyingMemberByID(members, test.accuserID)
			modifiedShareS := accuser.receivedQualifiedSharesS[test.accusedID]
			modifiedShareT := accuser.receivedQualifiedSharesT[test.accusedID]

			if test.modifyShareS != nil {
				modifiedShareS = test.modifyShareS(modifiedShareS)
			}
			if test.modifyShareT != nil {
				modifiedShareT = test.modifyShareT(modifiedShareT)
			}
			if test.modifyCommitments != nil {
				justifyingMember.receivedPeerCommitments[test.accusedID] =
					test.modifyCommitments(
						justifyingMember.receivedPeerCommitments[test.accusedID],
					)
			}

			// Simulate received PeerSharesMessage send by accused member.
			symmetricKey := accuser.symmetricKeys[test.accusedID]
			encryptedShareS, err := symmetricKey.Encrypt(modifiedShareS.Bytes())
			if err != nil {
				t.Fatalf("unexpected error: [%v]", err)
			}
			encryptedShareT, err := symmetricKey.Encrypt(modifiedShareT.Bytes())
			if err != nil {
				t.Fatalf("unexpected error: [%v]", err)
			}
			shares := make(map[group.MemberIndex]*peerShares)
			shares[test.accuserID] = &peerShares{encryptedShareS, encryptedShareT}
			justifyingMember.evidenceLog.PutPeerSharesMessage(
				&PeerSharesMessage{
					senderID: test.accusedID,
					shares:   shares,
				},
			)

			if test.modifyEvidenceLog != nil {
				justifyingMember.evidenceLog = test.modifyEvidenceLog(
					justifyingMember.evidenceLog,
				)
			}

			// Generate SecretSharesAccusationsMessages
			accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			accusedMembersKeys[test.accusedID] = accuser.ephemeralKeyPairs[test.accusedID].PrivateKey
			if test.modifyAccusedPrivateKey != nil {
				accusedMembersKeys[test.accusedID] = test.modifyAccusedPrivateKey(accusedMembersKeys[test.accusedID])
			}
			var messages []*SecretSharesAccusationsMessage
			messages = append(messages, &SecretSharesAccusationsMessage{
				senderID:           test.accuserID,
				accusedMembersKeys: accusedMembersKeys,
			})

			err = justifyingMember.ResolveSecretSharesAccusationsMessages(
				messages,
			)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf("\nexpected: %s\nactual:   %s\n", test.expectedError, err)
			}

			result := justifyingMember.group.DisqualifiedMemberIDs()
			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Fatalf("\nexpected: %d\nactual:   %d\n", test.expectedResult, result)
			}
		})
	}
}

// TODO Add test with many messages from accusers and many accused in the message.
func TestResolvePublicKeySharePointsAccusationsMessages(t *testing.T) {
	dishonestThreshold := 2
	groupSize := 5

	currentMemberID := group.MemberIndex(2) // i

	var tests = map[string]struct {
		accuserID                  group.MemberIndex // j
		accusedID                  group.MemberIndex // m
		modifyEvidenceLog          func(evidenceLog evidenceLog) evidenceLog
		modifyShareS               func(shareS *big.Int) *big.Int
		modifyPublicKeySharePoints func(points []*bn256.G2) []*bn256.G2
		modifyAccusedPrivateKey    func(symmetricKey *ephemeral.PrivateKey) *ephemeral.PrivateKey
		expectedResult             []group.MemberIndex
	}{
		"false accusation - accuser is disqualified": {
			accuserID:      3,
			accusedID:      4,
			expectedResult: []group.MemberIndex{3},
		},
		"current member as an accused - accuser is disqualified": {
			accuserID:      3,
			accusedID:      currentMemberID,
			expectedResult: []group.MemberIndex{3},
		},
		"incorrect shareS - accused member is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyShareS: func(shareS *big.Int) *big.Int {
				return new(big.Int).Sub(shareS, big.NewInt(1))
			},
			expectedResult: []group.MemberIndex{4},
		},
		"incorrect commitments - accused member is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyPublicKeySharePoints: func(points []*bn256.G2) []*bn256.G2 {
				newPoints := make([]*bn256.G2, len(points))
				for i := range newPoints {
					newPoints[i] = new(bn256.G2).ScalarBaseMult(
						big.NewInt(int64(i)),
					)
				}
				return newPoints
			},
			expectedResult: []group.MemberIndex{4},
		},
		"no commitments - accused member is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyPublicKeySharePoints: func(points []*bn256.G2) []*bn256.G2 {
				return []*bn256.G2{}
			},
			expectedResult: []group.MemberIndex{4},
		},
		"incorrect accused private key - accuser is disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyAccusedPrivateKey: func(symmetricKey *ephemeral.PrivateKey) *ephemeral.PrivateKey {
				return &ephemeral.PrivateKey{D: big.NewInt(12)}
			},
			expectedResult: []group.MemberIndex{3},
		},
		"inactive member as an accused (no EphemeralPublicKeyMessage sent) - " +
			"accuser is disqualified": {
			accuserID: 3,
			accusedID: 5,
			modifyEvidenceLog: func(evidenceLog evidenceLog) evidenceLog {
				if dkgEvidenceLog, ok := evidenceLog.(*dkgEvidenceLog); ok {
					dkgEvidenceLog.pubKeyMessageLog.removeMessage(
						group.MemberIndex(5),
					)
					return dkgEvidenceLog
				}
				return evidenceLog
			},
			expectedResult: []group.MemberIndex{3},
		},
		"inactive member as an accused (no PeerSharesMessage sent) - " +
			"accuser is disqualified": {
			accuserID: 3,
			accusedID: 5,
			modifyEvidenceLog: func(evidenceLog evidenceLog) evidenceLog {
				if dkgEvidenceLog, ok := evidenceLog.(*dkgEvidenceLog); ok {
					dkgEvidenceLog.peerSharesMessageLog.removeMessage(
						group.MemberIndex(5),
					)
					return dkgEvidenceLog
				}
				return evidenceLog
			},
			expectedResult: []group.MemberIndex{3},
		},
		"shares could not be decrypted - both are disqualified": {
			accuserID: 3,
			accusedID: 4,
			modifyEvidenceLog: func(evidenceLog evidenceLog) evidenceLog {
				if dkgEvidenceLog, ok := evidenceLog.(*dkgEvidenceLog); ok {
					message := dkgEvidenceLog.peerSharesMessage(group.MemberIndex(4))
					message.shares[group.MemberIndex(3)] = &peerShares{
						[]byte{0x00},
						[]byte{0x00},
					}
					return dkgEvidenceLog
				}
				return evidenceLog
			},
			expectedResult: []group.MemberIndex{3, 4},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializePointsJustifyingMemberGroup(
				dishonestThreshold,
				groupSize,
			)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			justifyingMember := findCoefficientsJustifyingMemberByID(members, currentMemberID)

			accuser := findCoefficientsJustifyingMemberByID(members, test.accuserID)
			modifiedShareS := accuser.receivedQualifiedSharesS[test.accusedID]
			if test.modifyShareS != nil {
				modifiedShareS = test.modifyShareS(modifiedShareS)
			}
			if test.modifyPublicKeySharePoints != nil {
				justifyingMember.receivedValidPeerPublicKeySharePoints[test.accusedID] =
					test.modifyPublicKeySharePoints(
						justifyingMember.receivedValidPeerPublicKeySharePoints[test.accusedID],
					)
			}

			// Simulate received PeerSharesMessage send by accused member.
			symmetricKey := accuser.symmetricKeys[test.accusedID]
			encryptedShareS, err := symmetricKey.Encrypt(modifiedShareS.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			encryptedShareT, err := symmetricKey.Encrypt(big.NewInt(13).Bytes())
			if err != nil {
				t.Fatal(err)
			}
			shares := make(map[group.MemberIndex]*peerShares)
			shares[test.accuserID] = &peerShares{encryptedShareS, encryptedShareT}
			justifyingMember.evidenceLog.PutPeerSharesMessage(
				&PeerSharesMessage{test.accusedID, shares},
			)

			if test.modifyEvidenceLog != nil {
				justifyingMember.evidenceLog = test.modifyEvidenceLog(
					justifyingMember.evidenceLog,
				)
			}

			// Generate PointsAccusationMessages
			accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			accusedMembersKeys[test.accusedID] = accuser.ephemeralKeyPairs[test.accusedID].PrivateKey
			if test.modifyAccusedPrivateKey != nil {
				accusedMembersKeys[test.accusedID] = test.modifyAccusedPrivateKey(accusedMembersKeys[test.accusedID])
			}
			var messages []*PointsAccusationsMessage
			messages = append(messages, &PointsAccusationsMessage{
				senderID:           test.accuserID,
				accusedMembersKeys: accusedMembersKeys,
			})

			err = justifyingMember.ResolvePublicKeySharePointsAccusationsMessages(
				messages,
			)
			if err != nil {
				t.Fatal(err)
			}

			result := justifyingMember.group.DisqualifiedMemberIDs()
			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Fatalf("\nexpected: %d\nactual:   %d\n", test.expectedResult, result)
			}
		})
	}
}

func findSharesJustifyingMemberByID(
	members []*SharesJustifyingMember,
	id group.MemberIndex,
) *SharesJustifyingMember {
	for _, m := range members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

func findCoefficientsJustifyingMemberByID(
	members []*PointsJustifyingMember,
	id group.MemberIndex) *PointsJustifyingMember {
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
func initializeSharesJustifyingMemberGroup(dishonestThreshold, groupSize int) (
	[]*SharesJustifyingMember,
	error,
) {
	commitmentsVerifyingMembers, err :=
		initializeCommitmentsVerifiyingMembersGroup(dishonestThreshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var sharesJustifyingMembers []*SharesJustifyingMember
	for _, cvm := range commitmentsVerifyingMembers {
		sharesJustifyingMembers = append(sharesJustifyingMembers,
			cvm.InitializeSharesJustification())
	}

	// Maps which will keep coefficients and commitments of all group members,
	// with members IDs as keys.
	groupCoefficientsA := make(map[group.MemberIndex][]*big.Int, groupSize)
	groupCoefficientsB := make(map[group.MemberIndex][]*big.Int, groupSize)
	groupCommitments := make(map[group.MemberIndex][]*bn256.G1, groupSize)

	for _, m := range sharesJustifyingMembers {
		memberCoefficientsA, err := generatePolynomial(dishonestThreshold)
		if err != nil {
			return nil, fmt.Errorf("polynomial generation failed [%s]", err)
		}
		memberCoefficientsB, err := generatePolynomial(dishonestThreshold)
		if err != nil {
			return nil, fmt.Errorf("polynomial generation failed [%s]", err)
		}

		// polynomial is of degree dishonestThreshold so it has
		// dishonestThreshold+1 coefficients including a constant coefficient
		commitments := make([]*bn256.G1, dishonestThreshold+1)
		for k := range memberCoefficientsA {

			commitments[k] = m.calculateCommitment(
				memberCoefficientsA[k],
				memberCoefficientsB[k],
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
				p.receivedQualifiedSharesS[m.ID] = m.evaluateMemberShare(p.ID, groupCoefficientsA[m.ID])
				p.receivedQualifiedSharesT[m.ID] = m.evaluateMemberShare(p.ID, groupCoefficientsB[m.ID])
				p.receivedPeerCommitments[m.ID] = groupCommitments[m.ID]
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
	dishonestThreshold, groupSize int,
) ([]*PointsJustifyingMember, error) {
	sharingMembers, err := initializeSharingMembersGroup(dishonestThreshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var pointsJustifyingMembers []*PointsJustifyingMember
	for _, sm := range sharingMembers {
		pointsJustifyingMembers = append(pointsJustifyingMembers,
			sm.InitializePointsJustification())
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
