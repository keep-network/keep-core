package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestRevealDisqualifiedMembersKeys(t *testing.T) {
	threshold := 2
	groupSize := 8

	members, err := initializeRevealingMembersGroup(threshold, groupSize, nil)
	if err != nil {
		t.Fatal(err)
	}
	member := members[0]

	disqualifiedSharingMember1 := MemberID(2)
	disqualifiedSharingMember2 := MemberID(3)
	disqualifiedNotSharingMember := MemberID(6)
	member.group.DisqualifyMemberID(disqualifiedSharingMember1)
	member.group.DisqualifyMemberID(disqualifiedSharingMember2)
	member.group.DisqualifyMemberID(disqualifiedNotSharingMember)

	// Simulate a case where member is disqualified in Phase 5.
	member.receivedValidSharesS[disqualifiedNotSharingMember] = nil

	expectedDisqualifiedKeys := map[MemberID]*ephemeral.PrivateKey{
		disqualifiedSharingMember1: member.ephemeralKeyPairs[disqualifiedSharingMember1].PrivateKey,
		disqualifiedSharingMember2: member.ephemeralKeyPairs[disqualifiedSharingMember2].PrivateKey,
	}

	result, err := member.RevealDisqualifiedMembersKeys()
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := &DisqualifiedEphemeralKeysMessage{
		senderID:    member.ID,
		privateKeys: expectedDisqualifiedKeys,
	}

	if !reflect.DeepEqual(expectedResult, result) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedResult, result)
	}
}

func TestRecoverDisqualifiedShares(t *testing.T) {
	threshold := 2
	groupSize := 6

	members, err := initializeReconstructingMembersGroup(threshold, groupSize, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Recovering member:
	member1 := members[0]
	// Other members:
	member2 := members[1]
	member3 := members[2]
	member4 := members[3]
	otherMembers := []*ReconstructingMember{member2, member3, member4}
	// Disqualified members:
	member5 := members[4]
	member6 := members[5]
	disqualifiedMembers := []*ReconstructingMember{member5, member6}

	var disqualifiedEphemeralKeysMessages []*DisqualifiedEphemeralKeysMessage
	for _, otherMember := range otherMembers {
		revealedKeys := make(map[MemberID]*ephemeral.PrivateKey)
		for _, disqualifiedMember := range disqualifiedMembers {
			revealedKeys[disqualifiedMember.ID] = otherMember.ephemeralKeyPairs[disqualifiedMember.ID].PrivateKey
		}
		disqualifiedEphemeralKeysMessages = append(
			disqualifiedEphemeralKeysMessages,
			&DisqualifiedEphemeralKeysMessage{
				senderID:    otherMember.ID,
				privateKeys: revealedKeys,
			},
		)
	}

	disqualifiedMemberShares := make(map[MemberID]map[MemberID]*big.Int)
	evidenceLog := member1.protocolConfig.evidenceLog

	for _, disqualifiedMember := range disqualifiedMembers {
		disqualifiedMemberShares[disqualifiedMember.ID] = make(map[MemberID]*big.Int)
		// Simulate message broadcasted by disqualified member in Phase 3.
		peerSharesMessage := newPeerSharesMessage(disqualifiedMember.ID)

		for _, otherMember := range otherMembers {
			// Simulate shares evaluation from Phase 3.
			shareS := disqualifiedMember.evaluateMemberShare(otherMember.ID, disqualifiedMember.secretCoefficients)
			disqualifiedMemberShares[disqualifiedMember.ID][otherMember.ID] = shareS

			peerSharesMessage.addShares(
				otherMember.ID,
				shareS,
				big.NewInt(0), // share T is not needed
				disqualifiedMember.symmetricKeys[otherMember.ID],
			)
		}
		evidenceLog.PutPeerSharesMessage(peerSharesMessage)
	}

	recoveredDisqualifiedShares, err := member1.recoverDisqualifiedShares(disqualifiedEphemeralKeysMessages)
	if err != nil {
		t.Fatal(err)
	}

	if len(recoveredDisqualifiedShares) != len(disqualifiedMembers) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			len(recoveredDisqualifiedShares),
			len(disqualifiedMembers),
		)
	}

	for _, recoveredDisqualifiedShare := range recoveredDisqualifiedShares {
		for _, disqualifiedMember := range disqualifiedMembers {
			if recoveredDisqualifiedShare.disqualifiedMemberID == disqualifiedMember.ID {
				expectedRecoveredDisqualifiedShares := &DisqualifiedShares{
					disqualifiedMemberID: disqualifiedMember.ID,
					peerSharesS:          disqualifiedMemberShares[disqualifiedMember.ID],
				}

				if !reflect.DeepEqual(
					expectedRecoveredDisqualifiedShares,
					recoveredDisqualifiedShare,
				) {
					t.Fatalf("\nexpected: %v\nactual:   %v\n",
						expectedRecoveredDisqualifiedShares,
						recoveredDisqualifiedShare,
					)
				}
			}
		}
	}
}

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
	expectedIndividualPrivateKey1 := disqualifiedMember1.individualPrivateKey()
	expectedIndividualPrivateKey2 := disqualifiedMember2.individualPrivateKey()

	allDisqualifiedShares := disqualifyMembers(group, disqualifiedMembersIDs)

	for _, m := range group {
		if !contains(disqualifiedMembersIDs, m.ID) {
			m.reconstructIndividualPrivateKeys(allDisqualifiedShares)

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
	dkg := &DKG{P: big.NewInt(179), Q: big.NewInt(89), evidenceLog: newDkgEvidenceLog()}
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
		reconstructingMember.reconstructIndividualPublicKeys()

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
	dkg := &DKG{P: big.NewInt(1907), Q: big.NewInt(953), evidenceLog: newDkgEvidenceLog()}

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

func TestReconstructDisqualifiedIndividualKeys(t *testing.T) {
	threshold := 2
	groupSize := 6

	members, err := initializeReconstructingMembersGroup(threshold, groupSize, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Recovering member:
	member1 := members[0]
	// Other members:
	member2 := members[1]
	member3 := members[2]
	member4 := members[3]
	otherMembers := []*ReconstructingMember{member2, member3, member4}
	// Disqualified members:
	member5 := members[4]
	member6 := members[5]
	disqualifiedMembers := []*ReconstructingMember{member5, member6}

	var disqualifiedEphemeralKeysMessages []*DisqualifiedEphemeralKeysMessage
	for _, otherMember := range otherMembers {
		revealedKeys := make(map[MemberID]*ephemeral.PrivateKey)
		for _, disqualifiedMember := range disqualifiedMembers {
			revealedKeys[disqualifiedMember.ID] = otherMember.ephemeralKeyPairs[disqualifiedMember.ID].PrivateKey
		}
		disqualifiedEphemeralKeysMessages = append(
			disqualifiedEphemeralKeysMessages,
			&DisqualifiedEphemeralKeysMessage{
				senderID:    otherMember.ID,
				privateKeys: revealedKeys,
			},
		)
	}

	evidenceLog := member1.protocolConfig.evidenceLog
	for _, disqualifiedMember := range disqualifiedMembers {
		// Simulate message broadcasted by disqualified member in Phase 3.
		peerSharesMessage := newPeerSharesMessage(disqualifiedMember.ID)

		for _, otherMember := range otherMembers {
			// Evaluate shares which were calculated in Phase 3.
			shareS := disqualifiedMember.evaluateMemberShare(otherMember.ID, disqualifiedMember.secretCoefficients)

			peerSharesMessage.addShares(
				otherMember.ID,
				shareS,
				big.NewInt(0), // share T is not needed
				disqualifiedMember.symmetricKeys[otherMember.ID],
			)
		}
		evidenceLog.PutPeerSharesMessage(peerSharesMessage)
	}

	member1.ReconstructDisqualifiedIndividualKeys(disqualifiedEphemeralKeysMessages)

	for _, disqualifiedMember := range disqualifiedMembers {
		if disqualifiedMember.individualPrivateKey().
			Cmp(member1.reconstructedIndividualPrivateKeys[disqualifiedMember.ID]) != 0 {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				disqualifiedMember.individualPrivateKey(),
				member1.reconstructedIndividualPrivateKeys[disqualifiedMember.ID],
			)
		}

		if member1.reconstructedIndividualPublicKeys[disqualifiedMember.ID].
			Cmp(disqualifiedMember.individualPublicKey()) != 0 {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				disqualifiedMember.individualPrivateKey(),
				member1.reconstructedIndividualPublicKeys[disqualifiedMember.ID],
			)
		}
	}
}

func initializeRevealingMembersGroup(threshold, groupSize int, dkg *DKG) (
	[]*RevealingMember, error) {
	pointsJustifyingMembers, err := initializePointsJustifyingMemberGroup(threshold, groupSize, dkg)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var revealingMembers []*RevealingMember
	for _, pjm := range pointsJustifyingMembers {
		revealingMembers = append(revealingMembers, pjm.InitializeRevealing())
	}

	return revealingMembers, nil
}

func initializeReconstructingMembersGroup(threshold, groupSize int, dkg *DKG) (
	[]*ReconstructingMember, error) {
	revealingMembers, err := initializeRevealingMembersGroup(threshold, groupSize, dkg)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var reconstructingMembers []*ReconstructingMember
	for _, rm := range revealingMembers {
		reconstructingMembers = append(reconstructingMembers,
			rm.InitializeReconstruction())
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
