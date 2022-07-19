package gjkr

import (
	"fmt"
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/group"
)

func TestCombineGroupPublicKey(t *testing.T) {
	dishonestThreshold := 1
	groupSize := 3

	expectedGroupPublicKey := new(bn256.G2).ScalarBaseMult(
		big.NewInt(243), // 10 + 20 + 30 + 91 + 92
	)
	members, err := initializeCombiningMembersGroup(dishonestThreshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}
	member := members[0]

	// Member's public coefficients. Zeroth coefficient is member's individual
	// public key.
	member.publicKeySharePoints = []*bn256.G2{
		new(bn256.G2).ScalarBaseMult(big.NewInt(10)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(11)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(12)),
	}

	// Public coefficients received from peer members. Each peer member's zeroth
	// coefficient is their individual public key.
	member.receivedValidPeerPublicKeySharePoints[2] = []*bn256.G2{
		new(bn256.G2).ScalarBaseMult(big.NewInt(20)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(21)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(22)),
	}
	member.receivedValidPeerPublicKeySharePoints[3] = []*bn256.G2{
		new(bn256.G2).ScalarBaseMult(big.NewInt(30)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(31)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(32)),
	}

	// Reconstructed individual public keys for disqualified members.
	member.reconstructedIndividualPublicKeys[4] = new(bn256.G2).ScalarBaseMult(
		big.NewInt(91),
	)
	member.reconstructedIndividualPublicKeys[5] = new(bn256.G2).ScalarBaseMult(
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

func TestCombineGroupPublicKeyShares(t *testing.T) {
	dishonestThreshold := 1
	groupSize := 3

	members, err := initializeCombiningMembersGroup(dishonestThreshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	member := members[0]

	// Member's public coefficients.
	member.publicKeySharePoints = []*bn256.G2{
		new(bn256.G2).ScalarBaseMult(big.NewInt(10)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(11)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(12)),
	}

	// Public coefficients received from peer members.
	member.receivedValidPeerPublicKeySharePoints[2] = []*bn256.G2{
		new(bn256.G2).ScalarBaseMult(big.NewInt(20)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(21)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(22)),
	}
	member.receivedValidPeerPublicKeySharePoints[3] = []*bn256.G2{
		new(bn256.G2).ScalarBaseMult(big.NewInt(30)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(31)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(32)),
	}

	member.ComputeGroupPublicKeyShares()
	groupPublicKeyShares := <-member.groupPublicKeySharesChannel

	expectedGroupPublicKeySharesLength := 2 // groupSize - 1 (combining member)
	if len(groupPublicKeyShares) != expectedGroupPublicKeySharesLength {
		t.Fatalf(
			"incorrect group public key shares count"+
				"\nexpected: %v\nactual:   %v\n",
			expectedGroupPublicKeySharesLength,
			len(groupPublicKeyShares),
		)
	}

	// Member 2
	// part1 = 10 * 2^0 + 11 * 2^1 + 12 * 2^2 = 80
	// part2 = 20 * 2^0 + 21 * 2^1 + 22 * 2^2 = 150
	// part3 = 30 * 2^0 + 31 * 2^1 + 32 * 2^2 = 220
	// groupPublicKeyShare = part1 + part2 + part3 = 450
	expectedGroupPublicKeyShareMember2 := new(bn256.G2).ScalarBaseMult(big.NewInt(450))
	if groupPublicKeyShares[2].String() != expectedGroupPublicKeyShareMember2.String() {
		t.Fatalf(
			"incorrect group public share for "+
				"member 2\nexpected: %v\nactual:   %v\n",
			expectedGroupPublicKeyShareMember2,
			groupPublicKeyShares[2],
		)
	}

	// Member 3
	// part1 = 10 * 3^0 + 11 * 3^1 + 12 * 3^2 = 151
	// part2 = 20 * 3^0 + 21 * 3^1 + 22 * 3^2 = 281
	// part3 = 30 * 3^0 + 31 * 3^1 + 32 * 3^2 = 411
	// groupPublicKeyShare = part1 + part2 + part3 = 843
	expectedGroupPublicKeyShareMember3 := new(bn256.G2).ScalarBaseMult(big.NewInt(843))
	if groupPublicKeyShares[3].String() != expectedGroupPublicKeyShareMember3.String() {
		t.Fatalf(
			"incorrect group public share for "+
				"member 3\nexpected: %v\nactual:   %v\n",
			expectedGroupPublicKeyShareMember3,
			groupPublicKeyShares[3],
		)
	}
}

func TestCombineGroupPublicKeyShares_WithReconstruction(t *testing.T) {
	dishonestThreshold := 1
	groupSize := 3

	members, err := initializeCombiningMembersGroup(dishonestThreshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	member := members[0]

	// Member's public coefficients.
	member.publicKeySharePoints = []*bn256.G2{
		new(bn256.G2).ScalarBaseMult(big.NewInt(10)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(11)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(12)),
	}

	// Public coefficients received from peer members.
	member.receivedValidPeerPublicKeySharePoints[2] = []*bn256.G2{
		new(bn256.G2).ScalarBaseMult(big.NewInt(20)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(21)),
		new(bn256.G2).ScalarBaseMult(big.NewInt(22)),
	}

	// Simulate that member 3 didn't send their public key share points,
	// became inactive at the beginning of phase 8 and their shares have
	// been revealed in phase 11.
	member.group.MarkMemberAsInactive(3)
	delete(member.receivedValidPeerPublicKeySharePoints, 3)
	member.revealedMisbehavedMembersShares = []*misbehavedShares{{
		misbehavedMemberID: 3,
		peerSharesS: map[group.MemberIndex]*big.Int{
			2: big.NewInt(50),
		},
	}}

	member.ComputeGroupPublicKeyShares()
	groupPublicKeyShares := <-member.groupPublicKeySharesChannel

	expectedGroupPublicKeySharesLength := 1 // groupSize - 1 (combining member) - 1 (inactive member)
	if len(groupPublicKeyShares) != expectedGroupPublicKeySharesLength {
		t.Fatalf(
			"incorrect group public key shares count"+
				"\nexpected: %v\nactual:   %v\n",
			expectedGroupPublicKeySharesLength,
			len(groupPublicKeyShares),
		)
	}

	// Member 2
	// part1 = 10 * 2^0 + 11 * 2^1 + 12 * 2^2 = 80
	// part2 = 20 * 2^0 + 21 * 2^1 + 22 * 2^2 = 150
	// part3 = 50 (revealed share)
	// groupPublicKeyShare = part1 + part2 + part3 = 280
	expectedGroupPublicKeyShareMember2 := new(bn256.G2).ScalarBaseMult(big.NewInt(280))
	if groupPublicKeyShares[2].String() != expectedGroupPublicKeyShareMember2.String() {
		t.Fatalf(
			"incorrect group public share for "+
				"member 2\nexpected: %v\nactual:   %v\n",
			expectedGroupPublicKeyShareMember2,
			groupPublicKeyShares[2],
		)
	}
}

func initializeCombiningMembersGroup(
	dishonestThreshold,
	groupSize int,
) ([]*CombiningMember, error) {
	reconstructingMembers, err := initializeReconstructingMembersGroup(
		dishonestThreshold,
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
