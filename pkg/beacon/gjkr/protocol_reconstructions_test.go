package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/group"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
)

func TestRevealMisbehavedMembersKeys(t *testing.T) {
	dishonestThreshold := 3
	groupSize := 8

	members, err := initializeRevealingMembersGroup(dishonestThreshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}
	firstMember := members[0]

	disqualifiedSharingMember1 := group.MemberIndex(2)
	disqualifiedSharingMember2 := group.MemberIndex(3)
	disqualifiedNotSharingMember := group.MemberIndex(6)
	firstMember.group.MarkMemberAsDisqualified(disqualifiedSharingMember1)
	firstMember.group.MarkMemberAsDisqualified(disqualifiedSharingMember2)
	firstMember.group.MarkMemberAsDisqualified(disqualifiedNotSharingMember)

	// Simulate a case where member is disqualified in Phase 5.
	delete(firstMember.receivedQualifiedSharesS, disqualifiedNotSharingMember)

	// Simulate a case where members are disqualified in Phase 8.
	delete(firstMember.receivedValidPeerPublicKeySharePoints, disqualifiedSharingMember1)
	delete(firstMember.receivedValidPeerPublicKeySharePoints, disqualifiedSharingMember2)

	expectedDisqualifiedKeys := map[group.MemberIndex]*ephemeral.PrivateKey{
		disqualifiedSharingMember1: firstMember.ephemeralKeyPairs[disqualifiedSharingMember1].PrivateKey,
		disqualifiedSharingMember2: firstMember.ephemeralKeyPairs[disqualifiedSharingMember2].PrivateKey,
	}
	expectedResult := &MisbehavedEphemeralKeysMessage{
		senderID:    firstMember.ID,
		privateKeys: expectedDisqualifiedKeys,
	}

	result, err := firstMember.RevealMisbehavedMembersKeys()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedResult, result) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedResult, result)
	}
}

func TestRevealMisbehavedMembersShares(t *testing.T) {
	dishonestThreshold := 2
	groupSize := 6

	members, err := initializeReconstructingMembersGroup(dishonestThreshold, groupSize)
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

	misbehavedEphemeralKeysMessages, err := generateMisbehavedEphemeralKeysMessages(otherMembers, disqualifiedMembers)
	if err != nil {
		t.Fatal(err)
	}
	expectedDisqualifiedShares := generateDisqualifiedMemberShares(member1, otherMembers, disqualifiedMembers)

	// Simulate a case when `invalidRevealingMember` reveals invalid ephemeral
	// private key for `clearedMember`, so the `invalidRevealingMember` gets
	// disqualified and disqualified share for this pair of members is not recovered.
	invalidRevealingMember := member3
	clearedMember := member5
	for _, message := range misbehavedEphemeralKeysMessages {
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

	// Fill `expectedMembersForReconstruction` slice stored in `member1` state
	// with disqualified members ids. Without this, `member1` will not be able
	// to add their own shares received from disqualified members.
	for _, member := range disqualifiedMembers {
		member1.expectedMembersForReconstruction = append(
			member1.expectedMembersForReconstruction,
			member.ID,
		)
	}

	// TEST
	recoveredDisqualifiedShares, err := member1.revealMisbehavedMembersShares(misbehavedEphemeralKeysMessages)
	if err != nil {
		t.Fatal(err)
	}

	expectedDisqualifiedMemberIDs := make([]group.MemberIndex, 0)
	for _, disqualifiedMember := range disqualifiedMembers {
		expectedDisqualifiedMemberIDs = append(expectedDisqualifiedMemberIDs, disqualifiedMember.ID)
	}
	expectedDisqualifiedMemberIDs = append(expectedDisqualifiedMemberIDs, invalidRevealingMember.ID)
	if !reflect.DeepEqual(expectedDisqualifiedMemberIDs, member1.group.DisqualifiedMemberIDs()) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			expectedDisqualifiedMemberIDs,
			member1.group.DisqualifiedMemberIDs(),
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
			if recoveredDisqualifiedShare.misbehavedMemberID == disqualifiedMember.ID {
				expectedRecoveredDisqualifiedShares := &misbehavedShares{
					misbehavedMemberID: disqualifiedMember.ID,
					peerSharesS:        expectedDisqualifiedShares[disqualifiedMember.ID],
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

func generateMisbehavedEphemeralKeysMessages(
	otherMembers, disqualifiedMembers []*ReconstructingMember,
) ([]*MisbehavedEphemeralKeysMessage, error) {
	var misbehavedEphemeralKeysMessages []*MisbehavedEphemeralKeysMessage
	for _, otherMember := range otherMembers {
		for _, disqualifiedMember := range disqualifiedMembers {
			otherMember.group.MarkMemberAsDisqualified(disqualifiedMember.ID)
			delete(otherMember.receivedValidPeerPublicKeySharePoints, disqualifiedMember.ID)
		}
		misbehavedEphemeralKeysMessage, err := otherMember.RevealMisbehavedMembersKeys()
		if err != nil {
			return nil, err
		}
		misbehavedEphemeralKeysMessages = append(
			misbehavedEphemeralKeysMessages,
			misbehavedEphemeralKeysMessage,
		)
	}
	return misbehavedEphemeralKeysMessages, nil
}

func generateDisqualifiedMemberShares(
	currentMember *ReconstructingMember,
	otherMembers, disqualifiedMembers []*ReconstructingMember,
) map[group.MemberIndex]map[group.MemberIndex]*big.Int {
	disqualifiedMemberShares := make(map[group.MemberIndex]map[group.MemberIndex]*big.Int)

	for _, disqualifiedMember := range disqualifiedMembers {
		disqualifiedMemberShares[disqualifiedMember.ID] = make(map[group.MemberIndex]*big.Int)
		// Simulate message broadcasted by disqualified member in Phase 3.
		peerSharesMessage := newPeerSharesMessage(disqualifiedMember.ID)
		commitments := make([]*bn256.G1, 0)

		for i, otherMember := range otherMembers {
			// Simulate shares evaluation from Phase 3.
			shareS := disqualifiedMember.evaluateMemberShare(
				otherMember.ID,
				disqualifiedMember.secretCoefficients,
			)
			disqualifiedMemberShares[disqualifiedMember.ID][otherMember.ID] = shareS

			peerSharesMessage.addShares(
				otherMember.ID,
				shareS,
				shareS, // In the sake of simplicity shareT == shareS
				disqualifiedMember.symmetricKeys[otherMember.ID],
			)

			coefficient := disqualifiedMember.secretCoefficients[i]
			commitments = append(
				commitments,
				// Same coefficient is used as shareT == shareS
				currentMember.calculateCommitment(coefficient, coefficient),
			)
		}
		currentMember.evidenceLog.PutPeerSharesMessage(peerSharesMessage)

		// Prepare commitments valid with simulated PeerSharesMessage as this
		// condition is checked before a share is added to te reconstruction process.
		currentMember.receivedPeerCommitments[disqualifiedMember.ID] =
			commitments

		// Add current member own shareS received from disqualified member
		disqualifiedMemberShares[disqualifiedMember.ID][currentMember.ID] =
			currentMember.receivedQualifiedSharesS[disqualifiedMember.ID]
	}
	return disqualifiedMemberShares
}

func TestReconstructIndividualPrivateKeys(t *testing.T) {
	dishonestThreshold := 2
	groupSize := 5

	disqualifiedMembersIDs := []group.MemberIndex{3, 5}

	group, err := initializeReconstructingMembersGroup(dishonestThreshold, groupSize)
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

func contains(slice []group.MemberIndex, value group.MemberIndex) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}

func TestCalculateReconstructedIndividualPublicKeys(t *testing.T) {
	groupSize := 3
	dishonestThreshold := 1

	disqualifiedMembersIDs := []int{4, 5} // m

	reconstructedIndividualPrivateKeys := make( // z_m
		map[group.MemberIndex]*big.Int,
		len(disqualifiedMembersIDs),
	)
	reconstructedIndividualPrivateKeys[4] = big.NewInt(14) // z_4
	reconstructedIndividualPrivateKeys[5] = big.NewInt(15) // z_5

	expectedIndividualPublicKeys := make( // y_m = g^{z_m}
		map[group.MemberIndex]*bn256.G2,
		len(disqualifiedMembersIDs),
	)
	expectedIndividualPublicKeys[4] = new(bn256.G2).ScalarBaseMult(
		reconstructedIndividualPrivateKeys[4],
	)
	expectedIndividualPublicKeys[5] = new(bn256.G2).ScalarBaseMult(
		reconstructedIndividualPrivateKeys[5],
	)

	members, err := initializeReconstructingMembersGroup(dishonestThreshold, groupSize)
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

func TestReconstructMisbehavedIndividualKeys(t *testing.T) {
	dishonestThreshold := 2
	groupSize := 6

	members, err := initializeReconstructingMembersGroup(dishonestThreshold, groupSize)
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

	// Disqualified members must be also disqualified
	// from the recovering member's perspective
	member1.group.MarkMemberAsDisqualified(member5.ID)
	member1.group.MarkMemberAsDisqualified(member6.ID)

	var misbehavedEphemeralKeysMessages []*MisbehavedEphemeralKeysMessage
	for _, otherMember := range otherMembers {
		revealedKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
		for _, disqualifiedMember := range disqualifiedMembers {
			revealedKeys[disqualifiedMember.ID] = otherMember.ephemeralKeyPairs[disqualifiedMember.ID].PrivateKey
		}
		misbehavedEphemeralKeysMessages = append(
			misbehavedEphemeralKeysMessages,
			&MisbehavedEphemeralKeysMessage{
				senderID:    otherMember.ID,
				privateKeys: revealedKeys,
			},
		)
	}

	for _, disqualifiedMember := range disqualifiedMembers {
		// Simulate message broadcasted by disqualified member in Phase 3.
		peerSharesMessage := newPeerSharesMessage(disqualifiedMember.ID)
		commitments := make([]*bn256.G1, 0)

		for i, otherMember := range otherMembers {
			// Evaluate shares which were calculated in Phase 3.
			shareS := disqualifiedMember.evaluateMemberShare(
				otherMember.ID,
				disqualifiedMember.secretCoefficients,
			)

			peerSharesMessage.addShares(
				otherMember.ID,
				shareS,
				shareS, // In the sake of simplicity shareT == shareS
				disqualifiedMember.symmetricKeys[otherMember.ID],
			)

			coefficient := disqualifiedMember.secretCoefficients[i]
			commitments = append(
				commitments,
				// Same coefficient is used as shareT == shareS
				member1.calculateCommitment(coefficient, coefficient),
			)
		}
		member1.evidenceLog.PutPeerSharesMessage(peerSharesMessage)

		// Prepare commitments valid with simulated PeerSharesMessage as this
		// condition is checked before a share is added to te reconstruction process.
		member1.receivedPeerCommitments[disqualifiedMember.ID] = commitments
	}

	member1.ReconstructMisbehavedIndividualKeys(misbehavedEphemeralKeysMessages)

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
	dishonestThreshold, groupSize int,
) ([]*RevealingMember, error) {
	pointsJustifyingMembers, err := initializePointsJustifyingMemberGroup(
		dishonestThreshold,
		groupSize,
	)
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
	dishonestThreshold,
	groupSize int,
) ([]*ReconstructingMember, error) {
	revealingMembers, err := initializeRevealingMembersGroup(
		dishonestThreshold,
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
	disqualifiedMembersIDs []group.MemberIndex) []*misbehavedShares {
	allDisqualifiedShares := make([]*misbehavedShares, len(disqualifiedMembersIDs))
	for i, disqualifiedMemberID := range disqualifiedMembersIDs {
		sharesReceivedFromDisqualifiedMember := make(map[group.MemberIndex]*big.Int,
			len(members)-len(disqualifiedMembersIDs))
		// for each group member
		for _, m := range members {
			// if the member has not been disqualified
			if !contains(disqualifiedMembersIDs, m.ID) {
				// collect all shares which this member received from disqualified
				// member and store them in sharesReceivedFromDisqualifiedMember
				for peerID, receivedShare := range m.receivedQualifiedSharesS {
					if peerID == disqualifiedMemberID {
						sharesReceivedFromDisqualifiedMember[m.ID] = receivedShare
						break
					}
				}
			}
		}
		allDisqualifiedShares[i] = &misbehavedShares{
			misbehavedMemberID: disqualifiedMemberID,
			peerSharesS:        sharesReceivedFromDisqualifiedMember,
		}
	}

	return allDisqualifiedShares
}
