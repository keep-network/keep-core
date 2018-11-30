package gjkr

import (
	"fmt"
	"math/big"
	"testing"
)

func TestReconstructIndividualPrivateKeys(t *testing.T) {
	threshold := 2
	groupSize := 5

	disqualifiedMembersIDs := []MemberID{3, 5}

	group, err := initializeReconstructingMembersGroup(threshold, groupSize, nil)
	if err != nil {
		t.Fatal(err)
	}

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

func contains(slice []MemberID, value MemberID) bool {
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

	reconstructedIndividualPrivateKeys := make(map[MemberID]*big.Int, len(disqualifiedMembersIDs)) // z_m
	reconstructedIndividualPrivateKeys[4] = big.NewInt(14)                                         // z_4
	reconstructedIndividualPrivateKeys[5] = big.NewInt(15)                                         // z_5

	expectedIndividualPublicKeys := make(map[MemberID]*big.Int, len(disqualifiedMembersIDs)) // y_m = g^{z_m} mod p
	expectedIndividualPublicKeys[4] = big.NewInt(43)                                         // 7^14 mod 179
	expectedIndividualPublicKeys[5] = big.NewInt(122)                                        // 7^15 mod 179

	members, err := initializeReconstructingMembersGroup(threshold, groupSize, dkg)
	if err != nil {
		t.Fatal(err)
	}

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

	expectedGroupPublicKey := big.NewInt(1620) // 10*20*30*91*92 mod 1620

	members, err := initializeCombiningMembersGroup(threshold, groupSize, dkg)
	if err != nil {
		t.Fatal(err)
	}
	member := members[0]

	// Member's public coefficients. Zeroth coefficient is member's individual
	// public key.
	member.publicKeySharePoints = []*big.Int{big.NewInt(10), big.NewInt(11), big.NewInt(12)}

	// Public coefficients received from peer members. Each peer member's zeroth
	// coefficient is their individual public key.
	member.receivedValidPeerPublicKeySharePoints[2] = []*big.Int{big.NewInt(20), big.NewInt(21), big.NewInt(22)}
	member.receivedValidPeerPublicKeySharePoints[3] = []*big.Int{big.NewInt(30), big.NewInt(31), big.NewInt(32)}

	// Reconstructed individual public keys for disqualified members.
	member.reconstructedIndividualPublicKeys[4] = big.NewInt(91)
	member.reconstructedIndividualPublicKeys[5] = big.NewInt(92)

	// Combine individual public keys of group members to get group public key.
	member.CombineGroupPublicKey()

	if member.groupPublicKey.Cmp(expectedGroupPublicKey) != 0 {
		t.Fatalf(
			"incorrect group public key for member %d\nexpected: %v\nactual:   %v\n",
			member.ID,
			expectedGroupPublicKey,
			member.groupPublicKey,
		)
	}
}

func initializeReconstructingMembersGroup(threshold, groupSize int, dkg *DKG) (
	[]*ReconstructingMember, error) {
	// TODO When whole protocol is implemented check if SharingMember type is really
	// the one expected here (should be the member from Phase 10)
	pointsJustifyingMembers, err := initializePointsJustifyingMemberGroup(threshold, groupSize, dkg)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var reconstructingMembers []*ReconstructingMember
	for _, pjm := range pointsJustifyingMembers {
		reconstructingMembers = append(reconstructingMembers,
			pjm.InitializeReconstruction())
	}

	return reconstructingMembers, nil
}

// disqualifyMembers disqualifies specific members for a test run. It collects
// shares calculated by disqualified members for their peers and reveals them.
func disqualifyMembers(
	members []*ReconstructingMember,
	disqualifiedMembersIDs []MemberID,
) []*DisqualifiedShares {
	allDisqualifiedShares := make([]*DisqualifiedShares, len(disqualifiedMembersIDs))
	for i, disqualifiedMemberID := range disqualifiedMembersIDs {
		sharesReceivedFromDisqualifiedMember := make(map[MemberID]*big.Int,
			len(members)-len(disqualifiedMembersIDs))
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

func initializeCombiningMembersGroup(threshold, groupSize int, dkg *DKG) ([]*CombiningMember, error) {
	// TODO When whole protocol is implemented check if SharingMember type is really
	// the one expected here (should be the member from Phase 10)
	reconstructingMembers, err := initializeReconstructingMembersGroup(threshold, groupSize, dkg)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var combiningMembers []*CombiningMember
	for _, rm := range reconstructingMembers {
		combiningMembers = append(combiningMembers, rm.InitializeCombining())
	}

	return combiningMembers, nil
}
