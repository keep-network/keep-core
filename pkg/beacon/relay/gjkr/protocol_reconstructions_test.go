package gjkr

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestReconstructIndividualPrivateKeys(t *testing.T) {
	threshold := 2
	groupSize := 5

	disqualifiedMembersIDs := []MemberID{3, 5}

	group, err := initializeReconstructingMembersGroup(threshold, groupSize)
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

	disqualifiedMembersIDs := []int{4, 5} // m

	reconstructedIndividualPrivateKeys := make( // z_m
		map[MemberID]*big.Int,
		len(disqualifiedMembersIDs),
	)
	reconstructedIndividualPrivateKeys[4] = big.NewInt(14) // z_4
	reconstructedIndividualPrivateKeys[5] = big.NewInt(15) // z_5

	expectedIndividualPublicKeys := make( // y_m = g^{z_m}
		map[MemberID]*bn256.G1,
		len(disqualifiedMembersIDs),
	)
	expectedIndividualPublicKeys[4] = new(bn256.G1).ScalarBaseMult(
		reconstructedIndividualPrivateKeys[4],
	)
	expectedIndividualPublicKeys[5] = new(bn256.G1).ScalarBaseMult(
		reconstructedIndividualPrivateKeys[5],
	)

	members, err := initializeReconstructingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	for _, member := range members {
		// Simulate phase where individual private keys are reconstructed.
		member.reconstructedIndividualPrivateKeys = reconstructedIndividualPrivateKeys
	}

	for _, reconstructingMember := range members {
		reconstructingMember.ReconstructIndividualPublicKeys()

		for disqualifiedMemberID, expectedIndividualPublicKey := range expectedIndividualPublicKeys {
			actualPublicKey := reconstructingMember.reconstructedIndividualPublicKeys[disqualifiedMemberID]
			if actualPublicKey.String() != expectedIndividualPublicKey.String() {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedIndividualPublicKey,
					actualPublicKey,
				)
			}
		}
	}
}

func TestCombineGroupPublicKey(t *testing.T) {
	threshold := 2
	groupSize := 3

	expectedGroupPublicKey := new(bn256.G1).ScalarBaseMult(
		big.NewInt(243), // 10 + 20 + 30 + 91 + 92
	)
	members, err := initializeCombiningMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}
	member := members[0]

	// Member's public coefficients. Zeroth coefficient is member's individual
	// public key.
	member.publicKeySharePoints = []*bn256.G1{
		new(bn256.G1).ScalarBaseMult(big.NewInt(10)),
		new(bn256.G1).ScalarBaseMult(big.NewInt(11)),
		new(bn256.G1).ScalarBaseMult(big.NewInt(12)),
	}

	// Public coefficients received from peer members. Each peer member's zeroth
	// coefficient is their individual public key.
	member.receivedValidPeerPublicKeySharePoints[2] = []*bn256.G1{
		new(bn256.G1).ScalarBaseMult(big.NewInt(20)),
		new(bn256.G1).ScalarBaseMult(big.NewInt(21)),
		new(bn256.G1).ScalarBaseMult(big.NewInt(22)),
	}
	member.receivedValidPeerPublicKeySharePoints[3] = []*bn256.G1{
		new(bn256.G1).ScalarBaseMult(big.NewInt(30)),
		new(bn256.G1).ScalarBaseMult(big.NewInt(31)),
		new(bn256.G1).ScalarBaseMult(big.NewInt(32)),
	}

	// Reconstructed individual public keys for disqualified members.
	member.reconstructedIndividualPublicKeys[4] = new(bn256.G1).ScalarBaseMult(
		big.NewInt(91),
	)
	member.reconstructedIndividualPublicKeys[5] = new(bn256.G1).ScalarBaseMult(
		big.NewInt(92),
	)

	// Combine individual public keys of group members to get group public key.
	member.CombineGroupPublicKey()

	if member.groupPublicKey.String() != expectedGroupPublicKey.String() {
		t.Fatalf(
			"incorrect group public key for member %d\nexpected: %v\nactual:   %v\n",
			member.ID,
			expectedGroupPublicKey,
			member.groupPublicKey,
		)
	}
}

func initializeReconstructingMembersGroup(threshold, groupSize int) (
	[]*ReconstructingMember, error) {
	// TODO When whole protocol is implemented check if SharingMember type is really
	// the one expected here (should be the member from Phase 10)
	pointsJustifyingMembers, err := initializePointsJustifyingMemberGroup(
		threshold,
		groupSize,
	)
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

func initializeCombiningMembersGroup(
	threshold,
	groupSize int,
) ([]*CombiningMember, error) {
	// TODO When whole protocol is implemented check if SharingMember type is really
	// the one expected here (should be the member from Phase 10)
	reconstructingMembers, err := initializeReconstructingMembersGroup(
		threshold,
		groupSize,
	)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var combiningMembers []*CombiningMember
	for _, rm := range reconstructingMembers {
		combiningMembers = append(combiningMembers, rm.InitializeCombining())
	}

	return combiningMembers, nil
}
