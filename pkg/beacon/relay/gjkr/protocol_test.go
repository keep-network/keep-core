package gjkr

import (
	"fmt"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/bls"
)

func TestRoundTrip(t *testing.T) {
	threshold := 3
	groupSize := 5

	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	var sharesMessages []*PeerSharesMessage
	var commitmentsMessages []*MemberCommitmentsMessage
	for _, member := range committingMembers {
		sharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
		if err != nil {
			t.Fatalf("shares and commitments calculation failed [%s]", err)
		}
		sharesMessages = append(sharesMessages, sharesMessage)
		commitmentsMessages = append(commitmentsMessages, commitmentsMessage)

	}

	var commitmentVerifyingMembers []*CommitmentsVerifyingMember
	for _, cm := range committingMembers {
		commitmentVerifyingMembers = append(commitmentVerifyingMembers,
			cm.InitializeCommitmentsVerification())
	}

	for _, member := range commitmentVerifyingMembers {
		accusedSecretSharesMessage, err := member.VerifyReceivedSharesAndCommitmentsMessages(
			filterPeerSharesMessage(sharesMessages, member.ID),
			filterMemberCommitmentsMessages(commitmentsMessages, member.ID),
		)
		if err != nil {
			t.Fatalf("shares and commitments verification failed [%s]", err)
		}

		if len(accusedSecretSharesMessage.accusedMembersKeys) > 0 {
			t.Fatalf("\nexpected: 0 accusations\nactual:   %d\n",
				len(accusedSecretSharesMessage.accusedMembersKeys),
			)
		}
	}

	var qualifiedMembers []*QualifiedMember
	for _, cvm := range commitmentVerifyingMembers {
		qualifiedMembers = append(qualifiedMembers,
			cvm.InitializeSharesJustification().InitializeQualified())
	}

	for _, member := range qualifiedMembers {
		member.CombineMemberShares()
	}

	var sharingMembers []*SharingMember
	for _, qm := range qualifiedMembers {
		sharingMembers = append(sharingMembers, qm.InitializeSharing())
	}

	for _, member := range sharingMembers {
		if len(member.receivedValidSharesS) != groupSize-1 {
			t.Fatalf("\nexpected: %d received shares S\nactual:   %d\n",
				groupSize-1,
				len(member.receivedValidSharesS),
			)
		}
		if len(member.receivedValidSharesT) != groupSize-1 {
			t.Fatalf("\nexpected: %d received shares T\nactual:   %d\n",
				groupSize-1,
				len(member.receivedValidSharesT),
			)
		}
		member.CombineMemberShares()
	}

	publicKeySharePointsMessages := make([]*MemberPublicKeySharePointsMessage, groupSize)
	for i, member := range sharingMembers {
		publicKeySharePointsMessages[i] = member.CalculatePublicKeySharePoints()
	}

	for _, member := range sharingMembers {
		accusedPointsMessage, err := member.VerifyPublicKeySharePoints(
			filterMemberPublicKeySharePointsMessages(publicKeySharePointsMessages, member.ID),
		)
		if err != nil {
			t.Fatalf("public coefficients verification failed [%s]", err)
		}
		if len(accusedPointsMessage.accusedMembersKeys) > 0 {
			t.Fatalf("\nexpected: 0 accusations\nactual:   %d\n",
				len(accusedPointsMessage.accusedMembersKeys),
			)
		}
	}

	var combiningMembers []*CombiningMember
	for _, sm := range sharingMembers {
		combiningMembers = append(combiningMembers,
			sm.InitializePointsJustification().InitializeRevealing().
				InitializeReconstruction().InitializeCombining())
	}

	for _, member := range combiningMembers {
		member.CombineGroupPublicKey()
	}

	fmt.Printf("Combined group public key: %x\n", combiningMembers[0].groupPublicKey.Marshal())

	publicKeyShares := make([]*bls.PublicKeyShare, 0)

	for _, member := range combiningMembers {
		groupPublicKeyShare := new(bn256.G2).ScalarBaseMult(member.groupPrivateKeyShare)

		publicKeyShares = append(publicKeyShares, &bls.PublicKeyShare{
			I: int(member.ID),
			V: groupPublicKeyShare,
		})
	}

	recoveredPublicKey, err := bls.RecoverPublicKey(publicKeyShares, threshold)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Recovered group public key: %x\n", recoveredPublicKey.Marshal())

	if recoveredPublicKey.String() != combiningMembers[0].groupPublicKey.String() {
		t.Fatalf(
			"\nexpected: %v\nactual:   %x\n",
			combiningMembers[0].groupPublicKey.Marshal(),
			recoveredPublicKey.Marshal(),
		)
	}

}
