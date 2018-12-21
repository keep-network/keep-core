package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

// TODO Add test with many messages from accusers and many accused in the message.
func TestResolveSecretSharesAccusations(t *testing.T) {
	threshold := 3
	groupSize := 5

	currentMemberID := MemberID(2) // i

	var tests = map[string]struct {
		accuserID               MemberID // j
		accusedID               MemberID // m
		modifyShareS            func(shareS *big.Int) *big.Int
		modifyShareT            func(shareT *big.Int) *big.Int
		modifyCommitments       func(commitments []*big.Int) []*big.Int
		modifyAccusedPrivateKey func(symmetricKey *ephemeral.PrivateKey) *ephemeral.PrivateKey
		expectedResult          []MemberID
		expectedError           error
	}{
		"false accusation - accuser is punished": {
			accuserID:      3,
			accusedID:      4,
			expectedResult: []MemberID{3},
		},
		"current member as an accuser - accusation skipped": {
			accuserID:      currentMemberID,
			accusedID:      3,
			expectedResult: []MemberID{},
		},
		"current member as an accused - accusation skipped": {
			accuserID:      3,
			accusedID:      currentMemberID,
			expectedResult: []MemberID{},
		},
		"incorrect shareS - accused member is punished": {
			accuserID: 3,
			accusedID: 4,
			modifyShareS: func(shareS *big.Int) *big.Int {
				return new(big.Int).Sub(shareS, big.NewInt(1))
			},
			expectedResult: []MemberID{4},
		},
		"incorrect shareT - accused member is punished": {
			accuserID: 3,
			accusedID: 4,
			modifyShareT: func(shareT *big.Int) *big.Int {
				return new(big.Int).Sub(shareT, big.NewInt(13))
			},
			expectedResult: []MemberID{4},
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
			expectedResult: []MemberID{4},
		},
		"incorrect accused private key - error returned": {
			accuserID: 3,
			accusedID: 4,
			modifyAccusedPrivateKey: func(symmetricKey *ephemeral.PrivateKey) *ephemeral.PrivateKey {
				return &ephemeral.PrivateKey{D: big.NewInt(12)}
			},
			// TODO Should we disqualify accuser/accused member here?
			expectedResult: nil,
			expectedError:  fmt.Errorf("could not decrypt shares [cannot decrypt share S [could not decrypt S share [symmetric key decryption failed]]]"),
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializeSharesJustifyingMemberGroup(threshold, groupSize, nil)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			member := findSharesJustifyingMemberByID(members, currentMemberID)

			accuser := findSharesJustifyingMemberByID(members, test.accuserID)
			modifiedShareS := accuser.receivedValidSharesS[test.accusedID]
			modifiedShareT := accuser.receivedValidSharesT[test.accusedID]

			if test.modifyShareS != nil {
				modifiedShareS = test.modifyShareS(modifiedShareS)
			}
			if test.modifyShareT != nil {
				modifiedShareT = test.modifyShareT(modifiedShareT)
			}
			if test.modifyCommitments != nil {
				member.receivedValidPeerCommitments[test.accusedID] =
					test.modifyCommitments(member.receivedValidPeerCommitments[test.accusedID])
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
			shares := make(map[MemberID]*peerShares)
			shares[test.accuserID] = &peerShares{encryptedShareS, encryptedShareT}
			member.evidenceLog.PutPeerSharesMessage(
				&PeerSharesMessage{
					senderID: test.accusedID,
					shares:   shares,
				},
			)

			// Generate SecretSharesAccusationsMessages
			accusedMembersKeys := make(map[MemberID]*ephemeral.PrivateKey)
			accusedMembersKeys[test.accusedID] = accuser.ephemeralKeyPairs[test.accusedID].PrivateKey
			if test.modifyAccusedPrivateKey != nil {
				accusedMembersKeys[test.accusedID] = test.modifyAccusedPrivateKey(accusedMembersKeys[test.accusedID])
			}
			var messages []*SecretSharesAccusationsMessage
			messages = append(messages, &SecretSharesAccusationsMessage{
				senderID:           test.accuserID,
				accusedMembersKeys: accusedMembersKeys,
			})

			result, err := member.ResolveSecretSharesAccusationsMessages(
				messages,
			)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf("\nexpected: %s\nactual:   %s\n", test.expectedError, err)
			}

			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Fatalf("\nexpected: %d\nactual:   %d\n", test.expectedResult, result)
			}
		})
	}
}

func TestRecoverSymmetricKey(t *testing.T) {
	member1ID := MemberID(1)
	member1KeyPair, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	member2ID := MemberID(2)
	member2KeyPair, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	expectedSymmetricKey := member1KeyPair.PrivateKey.Ecdh(member2KeyPair.PublicKey)

	messageBuffer := newDkgEvidenceLog()

	ephemeralPublicKeys1 := make(map[MemberID]*ephemeral.PublicKey)
	ephemeralPublicKeys1[member2ID] = member1KeyPair.PublicKey
	messageBuffer.PutEphemeralMessage(&EphemeralPublicKeyMessage{
		member1ID,
		ephemeralPublicKeys1,
	})

	ephemeralPublicKeys2 := make(map[MemberID]*ephemeral.PublicKey)
	ephemeralPublicKeys2[member1ID] = member2KeyPair.PublicKey
	messageBuffer.PutEphemeralMessage(&EphemeralPublicKeyMessage{
		member2ID,
		ephemeralPublicKeys2,
	})

	recoveredSymmetricKey, err := recoverSymmetricKey(
		messageBuffer,
		member1ID,                 // sender
		member2ID,                 // receiver
		member2KeyPair.PrivateKey, // receiver's private key
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedSymmetricKey, recoveredSymmetricKey) {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", expectedSymmetricKey, recoveredSymmetricKey)
	}
}

func TestRecoverShares(t *testing.T) {
	member1ID := MemberID(1)
	member1KeyPair, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	member2ID := MemberID(2)
	member2KeyPair, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	symmetricKey := member1KeyPair.PrivateKey.Ecdh(member2KeyPair.PublicKey)

	shareS := big.NewInt(21)
	shareT := big.NewInt(22)
	encryptedShareS, err := symmetricKey.Encrypt(shareS.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	encryptedShareT, err := symmetricKey.Encrypt(shareT.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	messageBuffer := newDkgEvidenceLog()
	shares := make(map[MemberID]*peerShares)
	shares[member2ID] = &peerShares{encryptedShareS, encryptedShareT}
	messageBuffer.PutPeerSharesMessage(&PeerSharesMessage{member1ID, shares})

	var tests = map[string]struct {
		symmetricKey   ephemeral.SymmetricKey
		expectedShareS *big.Int
		expectedShareT *big.Int
		expectedError  error
	}{
		"shares successfully recovered": {
			expectedShareS: shareS,
			expectedShareT: shareT,
			symmetricKey:   symmetricKey,
		},
		"shares recovery failed - incorrect symmetric key": {
			symmetricKey:  &ephemeral.SymmetricEcdhKey{},
			expectedError: fmt.Errorf("cannot decrypt share S [could not decrypt S share [symmetric key decryption failed]]"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			recoveredShareS, recoveredShareT, err := recoverShares(
				messageBuffer,
				member1ID,
				member2ID,
				test.symmetricKey,
			)
			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf("\nexpected: %s\nactual:   %s\n", test.expectedError, err)
			}
			if test.expectedShareS != nil {
				if test.expectedShareS.Cmp(recoveredShareS) != 0 {
					t.Fatalf("\nexpected: %s\nactual:   %s\n",
						test.expectedShareS,
						recoveredShareS,
					)
				}
			}
			if test.expectedShareT != nil {
				if test.expectedShareT.Cmp(recoveredShareT) != 0 {
					t.Fatalf("\nexpected: %s\nactual:   %s\n",
						test.expectedShareT,
						recoveredShareT,
					)
				}
			}
		})
	}
}

// TODO Add test with many messages from accusers and many accused in the message.
func TestResolvePublicKeySharePointsAccusationsMessages(t *testing.T) {
	threshold := 3
	groupSize := 5

	currentMemberID := MemberID(2) // i

	var tests = map[string]struct {
		accuserID                  MemberID // j
		accusedID                  MemberID // m
		modifyShareS               func(shareS *big.Int) *big.Int
		modifyPublicKeySharePoints func(coefficients []*big.Int) []*big.Int
		expectedResult             []MemberID
	}{
		"false accusation - sender is punished": {
			accuserID:      3,
			accusedID:      4,
			expectedResult: []MemberID{3},
		},
		"current member as a sender - accusation skipped": {
			accuserID:      currentMemberID,
			accusedID:      3,
			expectedResult: []MemberID{},
		},
		"current member as an accused - accusation skipped": {
			accuserID:      3,
			accusedID:      currentMemberID,
			expectedResult: []MemberID{},
		},
		"incorrect shareS - accused member is punished": {
			accuserID: 3,
			accusedID: 4,
			modifyShareS: func(shareS *big.Int) *big.Int {
				return new(big.Int).Sub(shareS, big.NewInt(1))
			},
			expectedResult: []MemberID{4},
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
			expectedResult: []MemberID{4},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializePointsJustifyingMemberGroup(threshold, groupSize, nil)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			member := findCoefficientsJustifyingMemberByID(members, currentMemberID)

			accuser := findCoefficientsJustifyingMemberByID(members, test.accuserID)
			modifiedShareS := accuser.receivedValidSharesS[test.accusedID]
			if test.modifyShareS != nil {
				modifiedShareS = test.modifyShareS(modifiedShareS)
			}
			if test.modifyPublicKeySharePoints != nil {
				member.receivedValidPeerPublicKeySharePoints[test.accusedID] =
					test.modifyPublicKeySharePoints(member.receivedValidPeerPublicKeySharePoints[test.accusedID])
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
			shares := make(map[MemberID]*peerShares)
			shares[test.accuserID] = &peerShares{encryptedShareS, encryptedShareT}
			member.evidenceLog.PutPeerSharesMessage(
				&PeerSharesMessage{test.accusedID, shares},
			)

			// Generate PointsAccusationMessages
			accusedMembersKeys := make(map[MemberID]*ephemeral.PrivateKey)
			accusedMembersKeys[test.accusedID] = accuser.ephemeralKeyPairs[test.accusedID].PrivateKey
			// if test.modifyAccusedPrivateKey != nil {
			// 	accusedMembersKeys[test.accusedID] = test.modifyAccusedPrivateKey(accusedMembersKeys[test.accusedID])
			// }
			var messages []*PointsAccusationsMessage
			messages = append(messages, &PointsAccusationsMessage{
				senderID:           test.accuserID,
				accusedMembersKeys: accusedMembersKeys,
			})

			result, err := member.ResolvePublicKeySharePointsAccusationsMessages(
				messages,
			)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Fatalf("\nexpected: %d\nactual:   %d\n", test.expectedResult, result)
			}
		})
	}
}

func findSharesJustifyingMemberByID(members []*SharesJustifyingMember, id MemberID) *SharesJustifyingMember {
	for _, m := range members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

func findCoefficientsJustifyingMemberByID(
	members []*PointsJustifyingMember,
	id MemberID,
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
	for _, cvm := range commitmentsVerifyingMembers {
		sharesJustifyingMembers = append(sharesJustifyingMembers,
			cvm.InitializeSharesJustification())
	}

	// Maps which will keep coefficients and commitments of all group members,
	// with members IDs as keys.
	groupCoefficientsA := make(map[MemberID][]*big.Int, groupSize)
	groupCoefficientsB := make(map[MemberID][]*big.Int, groupSize)
	groupCommitments := make(map[MemberID][]*big.Int, groupSize)

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

			commitments[k] = m.protocolConfig.vss.CalculateCommitment(
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
