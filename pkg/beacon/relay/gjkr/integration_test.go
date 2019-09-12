/*
  Integration tests for the full DKG affecting GJKR-specific parts.
*/
package gjkr_test

import (
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/dkgtest"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestExecute_HappyPath(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize)
	dkgtest.AssertMemberFailuresCount(t, result, 0)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
}

func TestExecute_IA_member1_phase1(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

func TestExecute_IA_members12_phase3(t *testing.T) {
	t.Parallel()

	groupSize := 7
	honestThreshold := 4

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		// drop commitment message from member 1
		commitmentMessage, ok := msg.(*gjkr.MemberCommitmentsMessage)
		if ok && commitmentMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		// drop shares message from member 2
		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		if ok && sharesMessage.SenderID() == group.MemberIndex(2) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5, 6, 7}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1), group.MemberIndex(2))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5, 6, 7}...)
}

func TestExecute_IA_member1_phase4(t *testing.T) {
	t.Parallel()

	groupSize := 3
	honestThreshold := 2

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		if ok && accusationsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3}...)
}

func TestExecute_IA_member1_phase7(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		sharePointsMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && sharePointsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

func TestExecute_IA_member1_phase8(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		accusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		if ok && accusationsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

func TestExecute_IA_members35_phase10(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		disqualifiedKeysMessage, ok := msg.(*gjkr.DisqualifiedEphemeralKeysMessage)
		if ok && (disqualifiedKeysMessage.SenderID() == group.MemberIndex(3) ||
			disqualifiedKeysMessage.SenderID() == group.MemberIndex(5)) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 4}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(3), group.MemberIndex(5))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 4}...)
}

// Phase 2 test case - a member sends an invalid ephemeral public key message.
// Message payload doesn't contain public keys for all other group members.
// Sender of the invalid message is disqualified by all of the receivers.
func TestExecute_DQ_member1_invalidMessage_phase2(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			publicKeyMessage.RemovePublicKey(group.MemberIndex(2))
			return publicKeyMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

// Phase 4 test case - a member sends an invalid member commitments message.
// Message payload doesn't contain a correct number of commitments.
// Sender of the invalid message is disqualified by all of the receivers.
func TestExecute_DQ_member5_invalidCommitmentsMessage_phase4(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		commitmentsMessage, ok := msg.(*gjkr.MemberCommitmentsMessage)
		if ok && commitmentsMessage.SenderID() == group.MemberIndex(5) {
			commitmentsMessage.RemoveCommitment(1)
			return commitmentsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 4}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(5))
	dkgtest.AssertInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 3, 4}...)
}

// Phase 4 test case - a member sends an invalid peer shares message.
// Message payload doesn't contain shares for all other group members.
// Sender of the invalid message is disqualified by all of the receivers.
func TestExecute_DQ_member4_invalidSharesMessage_phase4(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		if ok && sharesMessage.SenderID() == group.MemberIndex(4) {
			sharesMessage.RemoveShares(group.MemberIndex(1))
			return sharesMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(4))
	dkgtest.AssertInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 3, 5}...)
}

// Phase 5 test case - a member performs an accusation but reveals an
// ephemeral private key which doesn't correspond to the previously broadcast
// public key, generated for the sake of communication with the accused member.
// Due to such behaviour, the accuser is marked as disqualified in phase 5.
func TestExecute_DQ_member3_revealsWrongPrivateKey_phase5(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		if ok && accusationsMessage.SenderID() == group.MemberIndex(3) {
			// accuser (member 3) reveals a random private key which doesn't
			// correspond to the previously broadcast public key
			// generated for the sake of communication with the member 1
			randomKeyPair, _ := ephemeral.GenerateKeyPair()
			accusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(1),
				randomKeyPair.PrivateKey,
			)
			return accusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(3))
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 4, 5}...)
}

// Phase 5 test case - a member misbehaved by sending shares which
// cannot be decrypted by the receiver. The receiver makes an accusation
// which is confirmed by others so the misbehaving member is marked
// as disqualified in phase 5.
func TestExecute_DQ_member2_cannotDecryptTheirShares_phase5(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		if ok && sharesMessage.SenderID() == group.MemberIndex(2) {
			sharesMessage.SetShares(
				1,
				[]byte{0x00},
				[]byte{0x00},
			)
			return sharesMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(2))
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 3, 4, 5}...)
}

// Phase 5 test case - a member misbehaved by sending invalid commitment
// to another member. It becomes accused by the receiver of the
// invalid commitment. The accuser is right and the misbehaving member
// is marked as disqualified in phase 5.
func TestExecute_DQ_member5_inconsistentShares_phase5(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		commitmentsMessage, ok := msg.(*gjkr.MemberCommitmentsMessage)
		if ok && commitmentsMessage.SenderID() == group.MemberIndex(5) {
			commitmentsMessage.SetCommitment(
				2,
				new(bn256.G1).ScalarBaseMult(big.NewInt(1337)),
			)
			return commitmentsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 4}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(5))
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 3, 4}...)
}

// Phase 5 test case - a member misbehaved by performing a false accusation
// against another member. The accusation is checked by another members
// and because it is unfounded, the accuser is disqualified in phase 5.
func TestExecute_DQ_member4_falseAccusation_phase5(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	// Prepare a key pair which will be used to replace original ephemeral keys,
	// generated by the accuser for the sake of communication with accused member.
	// This way one can make an accusation and reveal an ephemeral private
	// key which will be valid from the perspective of members resolving
	// the accusation.
	keyPair, _ := ephemeral.GenerateKeyPair()

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		// Replace the original ephemeral public key published by the accuser
		// with the ephemeral public key generated before.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(4) {
			publicKeyMessage.SetPublicKey(
				group.MemberIndex(1),
				keyPair.PublicKey,
			)
			return publicKeyMessage
		}

		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		// Because of key replacement performed earlier, the accuser and
		// the accused member don't use the same symmetric key. In result,
		// the accused member could not decrypt shares received from the accuser
		// and performs an accusation against the accuser. This accusation will
		// cause disqualification of the accuser so it should be dropped as this
		// scenario aims to test disqualification due to false accusation only.
		if ok && accusationsMessage.SenderID() == group.MemberIndex(1) {
			accusationsMessage.RemoveAccusedMemberKey(group.MemberIndex(4))
			return accusationsMessage
		}
		// Accuser performs false accusation against the accused member
		// using the ephemeral private key generated before.
		if ok && accusationsMessage.SenderID() == group.MemberIndex(4) {
			accusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(1),
				keyPair.PrivateKey,
			)
			return accusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(4))
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 3, 5}...)
}

// Phase 5 test case - a member misbehaved by performing an accusation against
// an inactive member. The accusation could not be resolved by other members so
// the accuser is disqualified in phase 5.
func TestExecute_DQ_member2_accusesInactiveMember_phase5(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	// Prepare a key pair which will be used to replace original ephemeral keys,
	// generated by the accuser for the sake of communication with accused member.
	// This way one can make an accusation and reveal an ephemeral private
	// key which will be valid from the perspective of members resolving
	// the accusation.
	keyPair, _ := ephemeral.GenerateKeyPair()

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		// Replace the original ephemeral public key published by the accuser
		// with the ephemeral public key generated before.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(2) {
			publicKeyMessage.SetPublicKey(
				group.MemberIndex(1),
				keyPair.PublicKey,
			)
			return publicKeyMessage
		}
		// Drop message from accused member in order to simulate its inactivity.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		// Accuser performs accusation against the inactive accused member
		// using the ephemeral private key generated before.
		if ok && accusationsMessage.SenderID() == group.MemberIndex(2) {
			accusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(1),
				keyPair.PrivateKey,
			)
			return accusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(2))
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5}...)
}

// Phase 8 test case - a member sends an invalid member public key share points
// message. Message payload doesn't contain correct number of public key share
// points. Sender of the invalid message is disqualified by all of the receivers.
func TestExecute_DQ_member2_invalidMessage_phase8(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		sharePointsMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && sharePointsMessage.SenderID() == group.MemberIndex(2) {
			sharePointsMessage.RemovePublicKeyShare(0)
			return sharePointsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(2))
	dkgtest.AssertInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 3, 4, 5}...)
}

// Phase 9 test case - some members perform an accusation but reveal
// ephemeral private keys which don't correspond to the previously broadcast
// public keys, generated for the sake of communication with the accused members.
// Due to such behaviour, the accusers are marked as disqualified in phase 9.
func TestExecute_DQ_members25_revealWrongPrivateKey_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 7
	honestThreshold := 4

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		accusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		if ok && accusationsMessage.SenderID() == group.MemberIndex(2) {
			randomKeyPair, _ := ephemeral.GenerateKeyPair()
			accusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(1),
				randomKeyPair.PrivateKey,
			)
			accusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(3),
				randomKeyPair.PrivateKey,
			)
			return accusationsMessage
		}

		if ok && accusationsMessage.SenderID() == group.MemberIndex(5) {
			randomKeyPair, _ := ephemeral.GenerateKeyPair()
			accusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(4),
				randomKeyPair.PrivateKey,
			)
			return accusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 3, 4, 6, 7}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, []group.MemberIndex{2, 5}...)
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 3, 4, 6, 7}...)
}

// Phase 9 test case - some members misbehaved by sending in phase 7
// invalid public key shares to another members. They became accused in phase 8
// by the receivers of the invalid public key shares. The accusers are right
// and the misbehaving members are marked as disqualified in phase 9.
func TestExecute_DQ_members14_invalidPublicKeyShare_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		publicKeyShareMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && publicKeyShareMessage.SenderID() == group.MemberIndex(1) {
			publicKeyShareMessage.SetPublicKeyShare(
				1,
				new(bn256.G2).ScalarBaseMult(big.NewInt(5843)),
			)
			return publicKeyShareMessage
		}

		if ok && publicKeyShareMessage.SenderID() == group.MemberIndex(4) {
			publicKeyShareMessage.SetPublicKeyShare(
				2,
				new(bn256.G2).ScalarBaseMult(big.NewInt(7456)),
			)
			return publicKeyShareMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, []group.MemberIndex{1, 4}...)
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 5}...)
}

// Phase 9 test case - a member misbehaved by performing a false accusation
// against another member. The accusation is checked by another members
// and because it is unfounded, the accuser is disqualified in phase 9.
func TestExecute_DQ_member4_falseAccusation_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	// Prepare a key pair which will be used to replace original ephemeral keys,
	// generated by the accuser for the sake of communication with accused member.
	// This way one can make an accusation and reveal an ephemeral private
	// key which will be valid from the perspective of members resolving
	// the accusation.
	keyPair, _ := ephemeral.GenerateKeyPair()

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		// Replace the original ephemeral public key published by the accuser
		// with the ephemeral public key generated before.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(4) {
			publicKeyMessage.SetPublicKey(
				group.MemberIndex(1),
				keyPair.PublicKey,
			)
			return publicKeyMessage
		}

		sharesAccusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		// Because of key replacement performed earlier, the accuser and
		// the accused member don't use the same symmetric key. In result,
		// they could not decrypt shares received from each other which
		// causes accusations in phase 4. As we aim to test disqualification
		// in phase 9, we must drop these accusations.
		if ok && sharesAccusationsMessage.SenderID() == group.MemberIndex(1) {
			sharesAccusationsMessage.RemoveAccusedMemberKey(group.MemberIndex(4))
			return sharesAccusationsMessage
		}
		if ok && sharesAccusationsMessage.SenderID() == group.MemberIndex(4) {
			sharesAccusationsMessage.RemoveAccusedMemberKey(group.MemberIndex(1))
			return sharesAccusationsMessage
		}

		pointsAccusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		// Accuser performs false accusation against the accused member
		// using the ephemeral private key generated before.
		if ok && pointsAccusationsMessage.SenderID() == group.MemberIndex(4) {
			pointsAccusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(1),
				keyPair.PrivateKey,
			)
			return pointsAccusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	// Member 4 should fail and be disqualified because of false accusation.
	// Member 1 should also fail because it is not able to recover share
	// from member 4 in phase 11 and its signature diverge from the others.
	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(4))
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 5}...)
}

// Phase 9 test case - a member misbehaved by performing an accusation against
// an inactive member. The accusation could not be resolved by other members so
// the accuser is disqualified in phase 9.
func TestExecute_DQ_member2_accusesInactiveMember_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	// Prepare a key pair which will be used to replace original ephemeral keys,
	// generated by the accuser for the sake of communication with accused member.
	// This way one can make an accusation and reveal an ephemeral private
	// key which will be valid from the perspective of members resolving
	// the accusation.
	keyPair, _ := ephemeral.GenerateKeyPair()

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		// Replace the original ephemeral public key published by the accuser
		// with the ephemeral public key generated before.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(2) {
			publicKeyMessage.SetPublicKey(
				group.MemberIndex(1),
				keyPair.PublicKey,
			)
			return publicKeyMessage
		}
		// Drop message from accused member in order to simulate its inactivity.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		pointsAccusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		// Accuser performs accusation against the inactive accused member
		// using the ephemeral private key generated before.
		if ok && pointsAccusationsMessage.SenderID() == group.MemberIndex(2) {
			pointsAccusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(1),
				keyPair.PrivateKey,
			)
			return pointsAccusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(2))
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5}...)
}

// Phase 9 test case - a member misbehaved by sending shares which
// cannot be decrypted by the receiver. The receiver breaks the protocol
// and doesn't complain about this situation in phase 4. Instead of,
// the receiver performs an accusation in phase 8. In result both are marked
// as disqualified in phase 9.
func TestExecute_DQ_members34_cannotDecryptTheirShares_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	// Prepare a key pair which will be used to replace original ephemeral keys,
	// generated by the accuser for the sake of communication with accused member.
	// This way one can make an accusation and reveal an ephemeral private
	// key which will be valid from the perspective of members resolving
	// the accusation.
	keyPair, _ := ephemeral.GenerateKeyPair()

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		// Replace the original ephemeral public key published by the accuser
		// with the ephemeral public key generated before.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(3) {
			publicKeyMessage.SetPublicKey(
				group.MemberIndex(4),
				keyPair.PublicKey,
			)
			return publicKeyMessage
		}

		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		// Accused member actually misbehaves by sending shares which
		// cannot be decrypted by the accuser.
		if ok && sharesMessage.SenderID() == group.MemberIndex(4) {
			sharesMessage.SetShares(
				3,
				[]byte{0x00},
				[]byte{0x00},
			)
			return sharesMessage
		}

		sharesAccusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		// Because of key replacement performed earlier, the accuser and
		// the accused member don't use the same symmetric key. In result,
		// they could not decrypt shares received from each other which
		// causes accusations in phase 4. As we aim to test disqualification
		// in phase 9, we must drop these accusations.
		if ok && sharesAccusationsMessage.SenderID() == group.MemberIndex(3) {
			sharesAccusationsMessage.RemoveAccusedMemberKey(group.MemberIndex(4))
			return sharesAccusationsMessage
		}
		if ok && sharesAccusationsMessage.SenderID() == group.MemberIndex(4) {
			sharesAccusationsMessage.RemoveAccusedMemberKey(group.MemberIndex(3))
			return sharesAccusationsMessage
		}

		pointsAccusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		// Accuser performs accusation against the accused member only in
		// phase 8 using the ephemeral private key generated before.
		if ok && pointsAccusationsMessage.SenderID() == group.MemberIndex(3) {
			pointsAccusationsMessage.SetAccusedMemberKey(
				group.MemberIndex(4),
				keyPair.PrivateKey,
			)
			return pointsAccusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, []group.MemberIndex{3, 4}...)
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 5}...)
}

// Phase 11 test case - a member misbehaved by revealing key of an operating
// member. The revealing member becomes disqualified by all other members which
// consider the member for which the key has been revealed as normally operating.
// After phase 9, all group members should have the same view on who
// is disqualified. Revealing key of non-disqualified members is forbidden and
// leads to disqualifying the revealing member.
func TestExecute_DQ_member2_revealedKeyOfOperatingMember_phase11(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		disqualifiedKeysMessage, ok := msg.(*gjkr.DisqualifiedEphemeralKeysMessage)
		if ok && disqualifiedKeysMessage.SenderID() == group.MemberIndex(2) {
			randomKeyPair, _ := ephemeral.GenerateKeyPair()
			disqualifiedKeysMessage.SetPrivateKey(
				group.MemberIndex(3),
				randomKeyPair.PrivateKey,
			)
			return disqualifiedKeysMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(2))
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 3, 4, 5}...)
}
