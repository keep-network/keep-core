package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestRevealDisqualifiedMembersKeys(t *testing.T) {
	threshold := 2
	groupSize := 8

	members, err := initializeRevealingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}
	member := members[0]

	disqualifiedSharingMember1 := MemberID(2)
	disqualifiedSharingMember2 := MemberID(3)
	disqualifiedNotSharingMember := MemberID(6)
	member.group.MarkMemberAsDisqualified(disqualifiedSharingMember1)
	member.group.MarkMemberAsDisqualified(disqualifiedSharingMember2)
	member.group.MarkMemberAsDisqualified(disqualifiedNotSharingMember)

	// Simulate a case where member is disqualified in Phase 5.
	delete(member.receivedValidSharesS, disqualifiedNotSharingMember)

	expectedDisqualifiedKeys := map[MemberID]*ephemeral.PrivateKey{
		disqualifiedSharingMember1: member.ephemeralKeyPairs[disqualifiedSharingMember1].PrivateKey,
		disqualifiedSharingMember2: member.ephemeralKeyPairs[disqualifiedSharingMember2].PrivateKey,
	}
	expectedResult := &DisqualifiedEphemeralKeysMessage{
		senderID:    member.ID,
		privateKeys: expectedDisqualifiedKeys,
	}

	result, err := member.RevealDisqualifiedMembersKeys()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedResult, result) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedResult, result)
	}
}

func TestRecoverDisqualifiedShares(t *testing.T) {
	threshold := 2
	groupSize := 6

	members, err := initializeReconstructingMembersGroup(threshold, groupSize)
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

	disqualifiedEphemeralKeysMessages, err := generateDisqualifiedEphemeralKeysMessages(otherMembers, disqualifiedMembers)
	if err != nil {
		t.Fatal(err)
	}
	expectedDisqualifiedShares := generateDisqualifiedMemberShares(member1, otherMembers, disqualifiedMembers)

	// Simulate a case when `invalidRevealingMember` reveals invalid ephemeral
	// private key for `clearedMember`, so the `invalidRevealingMember` gets
	// disqualified and disqualified share for this pair of members is not recovered.
	invalidRevealingMember := member3
	clearedMember := member5
	for _, message := range disqualifiedEphemeralKeysMessages {
		if message.senderID == invalidRevealingMember.ID {
			newKeyPair, err := ephemeral.GenerateKeyPair()
			if err != nil {
				t.Fatal(err)
			}
			message.privateKeys[clearedMember.ID] = newKeyPair.PrivateKey
			break
		}
	}
	if _, ok := expectedDisqualifiedShares[clearedMember.ID][invalidRevealingMember.ID]; ok {
		delete(expectedDisqualifiedShares[clearedMember.ID], invalidRevealingMember.ID)
	}

	// TEST
	recoveredDisqualifiedShares, err := member1.recoverDisqualifiedShares(disqualifiedEphemeralKeysMessages)
	if err != nil {
		t.Fatal(err)
	}

	expectedDisqualifiedMemberIDs := make([]MemberID, 0)
	for _, disqualifiedMember := range disqualifiedMembers {
		expectedDisqualifiedMemberIDs = append(expectedDisqualifiedMemberIDs, disqualifiedMember.ID)
	}
	expectedDisqualifiedMemberIDs = append(expectedDisqualifiedMemberIDs, invalidRevealingMember.ID)
	if !reflect.DeepEqual(expectedDisqualifiedMemberIDs, member1.group.disqualifiedMemberIDs) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			expectedDisqualifiedMemberIDs,
			member1.group.disqualifiedMemberIDs,
		)
	}

	if len(recoveredDisqualifiedShares) != len(disqualifiedMembers) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			len(disqualifiedMembers),
			len(recoveredDisqualifiedShares),
		)
	}

	for _, recoveredDisqualifiedShare := range recoveredDisqualifiedShares {
		for _, disqualifiedMember := range disqualifiedMembers {
			if recoveredDisqualifiedShare.disqualifiedMemberID == disqualifiedMember.ID {
				expectedRecoveredDisqualifiedShares := &disqualifiedShares{
					disqualifiedMemberID: disqualifiedMember.ID,
					peerSharesS:          expectedDisqualifiedShares[disqualifiedMember.ID],
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

func generateDisqualifiedEphemeralKeysMessages(
	otherMembers, disqualifiedMembers []*ReconstructingMember,
) ([]*DisqualifiedEphemeralKeysMessage, error) {
	var disqualifiedEphemeralKeysMessages []*DisqualifiedEphemeralKeysMessage
	for _, otherMember := range otherMembers {
		for _, disqualifiedMember := range disqualifiedMembers {
			otherMember.group.MarkMemberAsDisqualified(disqualifiedMember.ID)
		}
		disqualifiedEphemeralKeysMessage, err := otherMember.RevealDisqualifiedMembersKeys()
		if err != nil {
			return nil, err
		}
		disqualifiedEphemeralKeysMessages = append(
			disqualifiedEphemeralKeysMessages,
			disqualifiedEphemeralKeysMessage,
		)
	}
	return disqualifiedEphemeralKeysMessages, nil
}

func generateDisqualifiedMemberShares(
	currentMember *ReconstructingMember,
	otherMembers, disqualifiedMembers []*ReconstructingMember,
) map[MemberID]map[MemberID]*big.Int {
	disqualifiedMemberShares := make(map[MemberID]map[MemberID]*big.Int)

	for _, disqualifiedMember := range disqualifiedMembers {
		disqualifiedMemberShares[disqualifiedMember.ID] = make(map[MemberID]*big.Int)
		// Simulate message broadcasted by disqualified member in Phase 3.
		peerSharesMessage := newPeerSharesMessage(disqualifiedMember.ID)

		for _, otherMember := range otherMembers {
			// Simulate shares evaluation from Phase 3.
			shareS := disqualifiedMember.evaluateMemberShare(
				otherMember.ID,
				disqualifiedMember.secretCoefficients,
			)
			disqualifiedMemberShares[disqualifiedMember.ID][otherMember.ID] = shareS

			peerSharesMessage.addShares(
				otherMember.ID,
				shareS,
				big.NewInt(0), // share T is not needed
				disqualifiedMember.symmetricKeys[otherMember.ID],
			)
		}
		currentMember.evidenceLog.PutPeerSharesMessage(peerSharesMessage)
	}
	return disqualifiedMemberShares
}

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
		reconstructingMember.reconstructIndividualPublicKeys()

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

func TestReconstructDisqualifiedIndividualKeys(t *testing.T) {
	threshold := 2
	groupSize := 6

	members, err := initializeReconstructingMembersGroup(threshold, groupSize)
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
		member1.evidenceLog.PutPeerSharesMessage(peerSharesMessage)
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

		if disqualifiedMember.individualPublicKey().String() !=
			member1.reconstructedIndividualPublicKeys[disqualifiedMember.ID].String() {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				disqualifiedMember.individualPrivateKey(),
				member1.reconstructedIndividualPublicKeys[disqualifiedMember.ID],
			)
		}
	}
}

func initializeRevealingMembersGroup(
	threshold, groupSize int,
) ([]*RevealingMember, error) {
	pointsJustifyingMembers, err := initializePointsJustifyingMemberGroup(threshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var revealingMembers []*RevealingMember
	for _, pjm := range pointsJustifyingMembers {
		revealingMembers = append(revealingMembers, pjm.InitializeRevealing())
	}

	return revealingMembers, nil
}

func initializeReconstructingMembersGroup(
	threshold,
	groupSize int,
) ([]*ReconstructingMember, error) {
	revealingMembers, err := initializeRevealingMembersGroup(
		threshold,
		groupSize,
	)
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
) []*disqualifiedShares {
	allDisqualifiedShares := make([]*disqualifiedShares, len(disqualifiedMembersIDs))
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
		allDisqualifiedShares[i] = &disqualifiedShares{
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
