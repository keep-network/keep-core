package gjkr

import (
	"math/big"
	"testing"
)

func TestReconstructIndividualPrivateKeys(t *testing.T) {
	threshold := 2
	groupSize := 5

	disqualifiedMembersIDs := []int{3, 5}

	group := initializeReconstructingMembersGroup(threshold, groupSize, nil)

	disqualifiedMember1 := group[2] // for ID = 3
	disqualifiedMember2 := group[4] // for ID = 5

	// polynomial's zeroth coefficient is member's individual private key
	expectedIndividualPrivateKey1 := disqualifiedMember1.secretCoefficients[0]
	expectedIndividualPrivateKey2 := disqualifiedMember2.secretCoefficients[0]

	allDisqualifiedShares := disqualifyMembers(group, disqualifiedMembersIDs)

	for _, m := range group {
		if !contains(disqualifiedMembersIDs, m.ID) {
			m.ReconstructIndividualPrivateKeys(allDisqualifiedShares)

			if m.reconstructedIndividualPrivateKeys[disqualifiedMember1.ID].Cmp(expectedIndividualPrivateKey1) != 0 {
				t.Fatalf("invalid reconstructed private key 1\nexpected: %s\nactual:   %s\n",
					expectedIndividualPrivateKey1,
					m.reconstructedIndividualPrivateKeys[disqualifiedMember1.ID],
				)
			}

			if m.reconstructedIndividualPrivateKeys[disqualifiedMember2.ID].Cmp(expectedIndividualPrivateKey2) != 0 {
				t.Fatalf("invalid reconstructed private key 2\nexpected: %s\nactual:   %s\n",
					expectedIndividualPrivateKey2,
					m.reconstructedIndividualPrivateKeys[disqualifiedMember2.ID],
				)
			}
		}
	}
}

func contains(slice []int, value int) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}

func TestCalculateReconstructedIndividualPublicKeys(t *testing.T) {
	groupSize := 3
	threshold := 2
	dkg := &DKG{P: big.NewInt(179), Q: big.NewInt(89)}
	g := big.NewInt(7) // `g` value for public key calculation `y_m = g^{z_m} mod p`

	disqualifiedMembersIDs := []int{4, 5} // m

	reconstructedIndividualPrivateKeys := make(map[int]*big.Int, len(disqualifiedMembersIDs)) // z_m
	reconstructedIndividualPrivateKeys[4] = big.NewInt(14)                                    // z_4
	reconstructedIndividualPrivateKeys[5] = big.NewInt(15)                                    // z_5

	expectedIndividualPublicKeys := make(map[int]*big.Int, len(disqualifiedMembersIDs)) // y_m = g^{z_m} mod p
	expectedIndividualPublicKeys[4] = big.NewInt(43)                                    // 7^14 mod 179
	expectedIndividualPublicKeys[5] = big.NewInt(122)                                   // 7^15 mod 179

	members := initializeReconstructingMembersGroup(threshold, groupSize, dkg)
	for _, member := range members {
		member.vss.G = g // set fixed `g` value
		// Simulate phase where individual private keys are reconstructed.
		member.reconstructedIndividualPrivateKeys = reconstructedIndividualPrivateKeys
	}

	for _, reconstructingMember := range members {
		reconstructingMember.ReconstructIndividualPublicKeys()

		for disqualifiedMemberID, expectedIndividualPublicKey := range expectedIndividualPublicKeys {
			if reconstructingMember.reconstructedIndividualPublicKeys[disqualifiedMemberID].
				Cmp(expectedIndividualPublicKey) != 0 {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedIndividualPublicKey,
					reconstructingMember.reconstructedIndividualPublicKeys[disqualifiedMemberID],
				)
			}
		}
	}
}

func TestCombineGroupPublicKey(t *testing.T) {
	threshold := 2
	groupSize := 3

	dkg := &DKG{P: big.NewInt(1907), Q: big.NewInt(953)}

	var tests = map[string]struct {
		disqualifiedMembers    int // number of disqualified members
		expectedGroupPublicKey *big.Int
	}{
		"no disqualified members - no reconstructed individual public key": {
			expectedGroupPublicKey: big.NewInt(279), // 10*20*30 mod 1907 = 279
		},
		"2 disqualified members - 2 reconstructed individual public keys": {
			disqualifiedMembers:    2,
			expectedGroupPublicKey: big.NewInt(1620), // 10*20*30*91*92 mod 1620
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members := initializeReconstructingMembersGroup(threshold, groupSize, dkg)

			// Generate member's public coefficients.
			for _, member := range members {
				member.publicKeySharePoints = make([]*big.Int, threshold+1)
				for k := range member.publicKeySharePoints {
					member.publicKeySharePoints[k] = big.NewInt(int64(member.ID*10 + k))
				}
			}

			// Configure public coefficients received from peer members.
			for _, member := range members {
				member.receivedValidPeerPublicKeySharePoints = make(map[int][]*big.Int, groupSize-1)
				for _, peer := range members {
					if member.ID != peer.ID {
						member.receivedValidPeerPublicKeySharePoints[peer.ID] =
							peer.publicKeySharePoints
					}
				}
			}

			// Configure reconstructed individual public key of disqualified members.
			//
			// Create as many reconstructed public keys as specified by disqualifiedMembers.
			// Reconstructed public keys will have an integer value starting
			// from 91 (91, 92, 93, ...).
			for _, member := range members {
				member.reconstructedIndividualPublicKeys = make(map[int]*big.Int, test.disqualifiedMembers)
				for m := 1; m <= test.disqualifiedMembers; m++ {
					disqualifiedMemberID := groupSize + m
					member.reconstructedIndividualPublicKeys[disqualifiedMemberID] =
						big.NewInt(int64(90 + m))
				}
			}

			for _, member := range members {
				member.CombineGroupPublicKey()

				if member.groupPublicKey.Cmp(test.expectedGroupPublicKey) != 0 {
					t.Fatalf(
						"incorrect group public key for member %d\nexpected: %v\nactual:   %v\n",
						member.ID,
						test.expectedGroupPublicKey,
						member.groupPublicKey,
					)
				}
			}
		})
	}
}

func initializeReconstructingMembersGroup(threshold, groupSize int, dkg *DKG) []*ReconstructingMember {
	// TODO When whole protocol is implemented check if SharingMember type is really
	// the one expected here (should be the member from Phase 10)
	sharingMembers, _ := initializeSharingMembersGroup(threshold, groupSize, dkg)

	var reconstructingMembers []*ReconstructingMember
	// TODO Should be handled by the `.Next()`` function
	for _, sm := range sharingMembers {
		reconstructingMembers = append(reconstructingMembers,
			&ReconstructingMember{
				SharingMember: sm,
			},
		)
	}

	return reconstructingMembers
}

// disqualifyMembers disqualifies specific members for a test run. It collects
// shares calculated by disqualified members for their peers and reveals them.
func disqualifyMembers(
	members []*ReconstructingMember,
	disqualifiedMembersIDs []int,
) []*DisqualifiedShares {
	allDisqualifiedShares := make([]*DisqualifiedShares, len(disqualifiedMembersIDs))
	for i, disqualifiedMemberID := range disqualifiedMembersIDs {
		sharesReceivedFromDisqualifiedMember := make(map[int]*big.Int, len(members)-len(disqualifiedMembersIDs))
		// for each group member
		for _, m := range members {
			// if the member has not been disqualified
			if !contains(disqualifiedMembersIDs, m.ID) {
				// collect all shares which this member received from disqualified
				// member and store them in sharesReceivedFromDisqualifiedMember
				for peerID, receivedShare := range m.receivedValidSharesS {
					if peerID == disqualifiedMemberID {
						sharesReceivedFromDisqualifiedMember[m.ID] = receivedShare
						break
					}
				}
			}
		}
		allDisqualifiedShares[i] = &DisqualifiedShares{
			disqualifiedMemberID: disqualifiedMemberID,
			peerSharesS:          sharesReceivedFromDisqualifiedMember,
		}
	}

	return allDisqualifiedShares
}
