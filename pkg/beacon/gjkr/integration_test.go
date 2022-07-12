// Package gjkr_test contains integration tests for the full roundtrip
// of GJKR-specific parts of DKG.
package gjkr_test

import (
	"math/big"
	"sync"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/beacon/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/group"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/internal/dkgtest"
	"github.com/keep-network/keep-core/pkg/net"
)

func TestExecute_HappyPath(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize)
	dkgtest.AssertMemberFailuresCount(t, result, 0)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoMisbehavingMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
}

func TestExecute_IA_member1_phase1(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

func TestExecute_IA_members12_phase3(t *testing.T) {
	t.Parallel()

	groupSize := 7
	honestThreshold := 4
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
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

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5, 6, 7}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(1), group.MemberIndex(2))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5, 6, 7}...)
}

func TestExecute_IA_member1_phase4(t *testing.T) {
	t.Parallel()

	groupSize := 3
	honestThreshold := 2
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		if ok && accusationsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3}...)
}

func TestExecute_IA_member1_phase7(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		sharePointsMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && sharePointsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

func TestExecute_IA_member1_phase8(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		accusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		if ok && accusationsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

func TestExecute_IA_members35_phase10(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		misbehavedKeysMessage, ok := msg.(*gjkr.MisbehavedEphemeralKeysMessage)
		if ok && (misbehavedKeysMessage.SenderID() == group.MemberIndex(3) ||
			misbehavedKeysMessage.SenderID() == group.MemberIndex(5)) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 4, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(3), group.MemberIndex(5))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 4, 6}...)
}

// Phase 2 test case - a member sends an invalid ephemeral public key message.
// Message payload doesn't contain public keys for all other group members.
// Sender of the invalid message is disqualified by all of the receivers.
func TestExecute_DQ_member1_invalidMessage_phase2(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			publicKeyMessage.RemovePublicKey(group.MemberIndex(2))
			return publicKeyMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(1))
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
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		commitmentsMessage, ok := msg.(*gjkr.MemberCommitmentsMessage)
		if ok && commitmentsMessage.SenderID() == group.MemberIndex(5) {
			commitmentsMessage.RemoveCommitment(1)
			return commitmentsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 4}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(5))
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
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		if ok && sharesMessage.SenderID() == group.MemberIndex(4) {
			sharesMessage.RemoveShares(group.MemberIndex(1))
			return sharesMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(4))
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
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
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

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(3))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 4, 5}...)
}

// Phase 5 test case - a member misbehaved by sending shares which
// cannot be decrypted by the receiver. The receiver makes an accusation
// which is confirmed by others so the misbehaving member is marked
// as disqualified in phase 5.
func TestExecute_DQ_member2_cannotDecryptTheirShares_phase5(t *testing.T) {
	t.Parallel()

	// Do not change groupSize and honestThreshold, such values are chosen
	// intentionally to check the behavior on the minimum honest threshold.
	groupSize := 3
	honestThreshold := 2
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
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

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 3}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(2))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 3}...)
}

// Phase 5 test case - a member misbehaved by sending invalid commitment
// to another member. It becomes accused by the receiver of the
// invalid commitment. The accuser is right and the misbehaving member
// is marked as disqualified in phase 5.
func TestExecute_DQ_member5_inconsistentShares_phase5(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
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

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 4}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(5))
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
	seed := dkgtest.RandomSeed(t)

	manInTheMiddle, err := newManInTheMiddle(
		group.MemberIndex(4), // sender
		groupSize,
		honestThreshold,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		manInTheMiddle.interceptCommunication(msg)

		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		// Accuser performs false accusation against the accused member
		// using the ephemeral private key generated before. We replace
		// the whole accusedMemberKeys because the real member covered by
		// the MiM performs their own accusations.
		if ok && accusationsMessage.SenderID() == group.MemberIndex(4) {
			accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			accusedMembersKeys[group.MemberIndex(1)] =
				manInTheMiddle.ephemeralKeyPairs[group.MemberIndex(1)].PrivateKey
			accusationsMessage.SetAccusedMemberKeys(accusedMembersKeys)
			return accusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(4))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 3, 5}...)
}

// Phase 5 test case - a member misbehaved by performing an accusation against
// an inactive member. The accusation could not be resolved by other members so
// the accuser is disqualified in phase 5.
func TestExecute_DQ_member2_accusesInactiveMember_phase5(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	manInTheMiddle, err := newManInTheMiddle(
		group.MemberIndex(2), // sender
		groupSize,
		honestThreshold,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		manInTheMiddle.interceptCommunication(msg)

		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		// Drop message from accused member in order to simulate its inactivity.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		// Accuser performs accusation against the inactive accused member
		// using the ephemeral private key generated before.
		// Man-in-the-middle intercepts and replaces entire communication
		// with member 2.
		if ok && accusationsMessage.SenderID() == group.MemberIndex(2) {
			accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			accusedMembersKeys[group.MemberIndex(1)] =
				manInTheMiddle.ephemeralKeyPairs[group.MemberIndex(1)].PrivateKey
			accusationsMessage.SetAccusedMemberKeys(accusedMembersKeys)
			return accusationsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{1, 2}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5, 6}...)
}

// Phase 8 test case - a member sends an invalid member public key share points
// message. Message payload doesn't contain correct number of public key share
// points. Sender of the invalid message is disqualified by all of the receivers.
func TestExecute_DQ_member2_invalidMessage_phase8(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		sharePointsMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && sharePointsMessage.SenderID() == group.MemberIndex(2) {
			sharePointsMessage.RemovePublicKeyShare(0)
			return sharePointsMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(2))
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
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
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

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 3, 4, 6, 7}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{2, 5}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 3, 4, 6, 7}...)
}

// Phase 9 test case - some members misbehaved by sending in phase 7
// invalid public key shares to another members. They became accused in phase 8
// by the receivers of the invalid public key shares. The accusers are right
// and the misbehaving members are marked as disqualified in phase 9.
func TestExecute_DQ_members14_invalidPublicKeyShare_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
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

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 5, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{1, 4}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 5, 6}...)
}

// Phase 9 test case - a member misbehaved by performing a false accusation
// against another member. The accusation is checked by another members
// and because it is unfounded, the accuser is disqualified in phase 9.
func TestExecute_DQ_member4_falseAccusation_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	manInTheMiddle, err := newManInTheMiddle(
		group.MemberIndex(4), // sender
		groupSize,
		honestThreshold,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		pointsAccusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		// Modify default man-in-the-middle behavior: accuser performs false
		// accusation against the accused member using the ephemeral private key
		// generated before.
		if ok && pointsAccusationsMessage.SenderID() == group.MemberIndex(4) {
			accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			accusedMembersKeys[group.MemberIndex(1)] =
				manInTheMiddle.ephemeralKeyPairs[group.MemberIndex(1)].PrivateKey
			pointsAccusationsMessage.SetAccusedMemberKeys(accusedMembersKeys)
			return pointsAccusationsMessage
		}

		manInTheMiddle.interceptCommunication(msg)

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(4))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 3, 5}...)
}

// Phase 9 test case - a member misbehaved by performing an accusation against
// an inactive member. The accusation could not be resolved by other members so
// the accuser is disqualified in phase 9.
func TestExecute_DQ_member2_accusesInactiveMember_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	manInTheMiddle, err := newManInTheMiddle(
		group.MemberIndex(2), // sender
		groupSize,
		honestThreshold,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		pointsAccusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		// Modify default man-in-the-middle behavior: accuser performs
		// accusation against the inactive member
		// using the ephemeral private key generated before.
		if ok && pointsAccusationsMessage.SenderID() == group.MemberIndex(2) {
			accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			accusedMembersKeys[group.MemberIndex(1)] =
				manInTheMiddle.ephemeralKeyPairs[group.MemberIndex(1)].PrivateKey
			pointsAccusationsMessage.SetAccusedMemberKeys(accusedMembersKeys)
			return pointsAccusationsMessage
		}

		manInTheMiddle.interceptCommunication(msg)

		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		// Drop message from accused member in order to simulate its inactivity.
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{1, 2}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5, 6}...)
}

// Phase 9 test case - a member misbehaved by sending shares which
// cannot be decrypted by the receiver. The receiver breaks the protocol
// and doesn't complain about this situation in phase 4. Instead of,
// the receiver performs an accusation in phase 8. In result both are marked
// as disqualified in phase 9.
func TestExecute_DQ_members12_cannotDecryptTheirShares_phase9(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	manInTheMiddle, err := newManInTheMiddle(
		group.MemberIndex(1), // sender
		groupSize,
		honestThreshold,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		pointsAccusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		// Modify default man-in-the-middle behavior: accuser performs
		// accusation against the member only in phase 8 using the ephemeral
		// private key generated before.
		if ok && pointsAccusationsMessage.SenderID() == group.MemberIndex(1) {
			accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			accusedMembersKeys[group.MemberIndex(2)] =
				manInTheMiddle.ephemeralKeyPairs[group.MemberIndex(2)].PrivateKey
			pointsAccusationsMessage.SetAccusedMemberKeys(accusedMembersKeys)
			return pointsAccusationsMessage
		}

		manInTheMiddle.interceptCommunication(msg)

		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		// Accused member misbehaves by sending shares which cannot be
		// decrypted by the accuser.
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

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{1, 2}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5, 6}...)
}

// Phase 11 test case - a member misbehaved by not revealing its private key
// generated for the sake of communication with a disqualified QUAL member.
// After phase 9, all group members should have the same view on who is disqualified.
// Not revealing key of any disqualified member from QUAL is considered as misbehaviour
// and leads to disqualification of the member which was supposed to reveal the key.
func TestExecute_DQ_member2_notRevealsDisqualifiedQualMemberKey_phase11(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		// Member 1 misbehaves by sending an invalid public key share points message.
		// As result, they should be disqualified by other members.
		publicKeyShareMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && publicKeyShareMessage.SenderID() == group.MemberIndex(1) {
			publicKeyShareMessage.RemovePublicKeyShare(2)
			return publicKeyShareMessage
		}

		// Member 2 does not reveal private key generated for the sake
		// of communication with disqualified member 1 and member 1
		// is in QUAL set so it has been disqualified after phase 5.
		// For this reason, member 2 should be also disqualified.
		misbehavedKeysMessage, ok := msg.(*gjkr.MisbehavedEphemeralKeysMessage)
		if ok && misbehavedKeysMessage.SenderID() == group.MemberIndex(2) {
			misbehavedKeysMessage.RemovePrivateKey(group.MemberIndex(1))
			return misbehavedKeysMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{1, 2}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5, 6}...)
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
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		misbehavedKeysMessage, ok := msg.(*gjkr.MisbehavedEphemeralKeysMessage)
		if ok && misbehavedKeysMessage.SenderID() == group.MemberIndex(2) {
			randomKeyPair, _ := ephemeral.GenerateKeyPair()
			misbehavedKeysMessage.SetPrivateKey(
				group.MemberIndex(3),
				randomKeyPair.PrivateKey,
			)
			return misbehavedKeysMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(2))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 3, 4, 5}...)
}

// Phase 11 test case - a member reveal an ephemeral private key which doesn't
// correspond to the previously broadcast public key, generated for the sake of
// communication with the disqualified member. Due to such behaviour, the
// revealing member is marked as disqualified in phase 11.
func TestExecute_DQ_member5_revealsWrongPrivateKey_phase11(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		// Member 4 misbehaves by sending invalid public key shares.
		// As a result it become disqualified in phase 9.
		publicKeyShareMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && publicKeyShareMessage.SenderID() == group.MemberIndex(4) {
			publicKeyShareMessage.SetPublicKeyShare(
				1,
				new(bn256.G2).ScalarBaseMult(big.NewInt(5843)),
			)
			return publicKeyShareMessage
		}

		// Member 5 should reveal private key generated for the sake of
		// communication with disqualified member 4. Instead of revealing
		// private key matching previously announced public key, member 5
		// reveals some other key. As a result, member 5 should be disqualified.
		misbehavedKeysMessage, ok := msg.(*gjkr.MisbehavedEphemeralKeysMessage)
		if ok && misbehavedKeysMessage.SenderID() == group.MemberIndex(5) {
			randomKeyPair, _ := ephemeral.GenerateKeyPair()
			misbehavedKeysMessage.SetPrivateKey(
				group.MemberIndex(4),
				randomKeyPair.PrivateKey,
			)
			return misbehavedKeysMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 3, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{4, 5}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 3, 6}...)
}

// Phase 11 test case - a member misbehaved by revealing private key generated for
// the sake of communication with a member which became inactive before phase 5.
// Shares reconstruction could not be resolved by other members so the revealing
// member is disqualified in phase 11.
func TestExecute_DQ_member2_revealsInactiveNonQualMemberKey_phase11(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	manInTheMiddle, err := newManInTheMiddle(
		group.MemberIndex(2), // sender
		groupSize,
		honestThreshold,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		misbehavedKeysMessage, ok := msg.(*gjkr.MisbehavedEphemeralKeysMessage)
		// Modify default man-in-the-middle behavior: reveals private key
		// generated for the sake of communication with member marked as
		// inactive before phase 5
		if ok && misbehavedKeysMessage.SenderID() == group.MemberIndex(2) {
			privateKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			privateKeys[group.MemberIndex(1)] =
				manInTheMiddle.ephemeralKeyPairs[group.MemberIndex(1)].PrivateKey
			misbehavedKeysMessage.SetPrivateKeys(privateKeys)
			return misbehavedKeysMessage
		}

		manInTheMiddle.interceptCommunication(msg)

		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		// Drop message from revealed member in order to simulate its inactivity.
		if ok && sharesMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{3, 4, 5, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{1, 2}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{3, 4, 5, 6}...)
}

// Phase 11 test case - a member misbehaved by sending shares which
// cannot be decrypted by the receiver. For this reason, member is disqualified
// in phase 5.
// In phase 10, the receiver reveals private key generated for the sake of
// communication with that member.
// This is a protocol violation because only disqualified sharing members (QUAL)
// should be included when revealing keys generated for them.
// As a result, the revealing member is disqualified in phase 11.
func TestExecute_DQ_member3_revealsDisqualifiedNonQualMemberKey_phase11(t *testing.T) {
	t.Parallel()

	groupSize := 6
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	manInTheMiddle, err := newManInTheMiddle(
		group.MemberIndex(3), // sender
		groupSize,
		honestThreshold,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		// Modify default man-in-the-middle behavior: revealing member perform a
		// justified accusation against member 4 which sent invalid shares.
		if ok && accusationsMessage.SenderID() == group.MemberIndex(3) {
			accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			accusedMembersKeys[group.MemberIndex(4)] =
				manInTheMiddle.ephemeralKeyPairs[group.MemberIndex(4)].PrivateKey
			accusationsMessage.SetAccusedMemberKeys(accusedMembersKeys)
			return accusationsMessage
		}

		misbehavedKeysMessage, ok := msg.(*gjkr.MisbehavedEphemeralKeysMessage)
		// Modify default man-in-the-middle behavior: revealing member reveals
		// private key generated for the sake of communication with member 4.
		if ok && misbehavedKeysMessage.SenderID() == group.MemberIndex(3) {
			privateKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
			privateKeys[group.MemberIndex(4)] =
				manInTheMiddle.ephemeralKeyPairs[group.MemberIndex(4)].PrivateKey
			misbehavedKeysMessage.SetPrivateKeys(privateKeys)
			return misbehavedKeysMessage
		}

		manInTheMiddle.interceptCommunication(msg)

		sharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
		// Revealed member misbehaves by sending shares which cannot be
		// decrypted by the revealing member.
		if ok && sharesMessage.SenderID() == group.MemberIndex(4) {
			sharesMessage.SetShares(
				3,
				[]byte{0x00},
				[]byte{0x00},
			)
			return sharesMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{1, 2, 5, 6}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, []group.MemberIndex{3, 4}...)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 5, 6}...)
}

func TestExecute_InvalidMemberIndex(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			// Set an non-existing sender id.
			publicKeyMessage.SetSenderID(group.MemberIndex(groupSize + 1))
			return publicKeyMessage
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 4, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertMisbehavingMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

// manInTheMiddle is a helper tool allowing to easily intercept communication
// of a chosen member with the rest of the members for all phases of DKG.
// Man in the middle sets up symmetric keys, member shares, and commitments
// that can be later accessed in test. It also publishes matching public
// key share points. By default, man-in-the-middle drops all accusations from
// the original member. This behavior can be overwitten in test before
// man-in-the-middle intercepts.
//
// In a test not using manInTheMiddle we don't have an access to symmetric keys
// of a chosen member established with the rest of the members. Knowledge of
// these symmetric keys is required to test some scenarios.
type manInTheMiddle struct {
	senderIndex   group.MemberIndex
	receiverIndex group.MemberIndex

	seed *big.Int

	// phase 1
	ephemeralKeyPairs map[group.MemberIndex]*ephemeral.KeyPair

	// phase 2
	symmetricKeys      map[group.MemberIndex]*ephemeral.SymmetricEcdhKey
	symmetricKeysMutex sync.Mutex

	// phase 3
	sharesS     map[group.MemberIndex]*big.Int
	sharesT     map[group.MemberIndex]*big.Int
	commitments []*bn256.G1

	// phase 7
	publicKeySharePoints []*bn256.G2
}

// newManInTheMiddle creates a new instance of manInTheMiddle tool.
// It will intercept messages sent from the sender with the given index.
func newManInTheMiddle(
	senderIndex group.MemberIndex,
	groupSize, honestThreshold int,
	seed *big.Int,
) (*manInTheMiddle, error) {
	ephemeralKeyPairs := make(map[group.MemberIndex]*ephemeral.KeyPair, groupSize-1)
	sharesS := make(map[group.MemberIndex]*big.Int, groupSize-1)
	sharesT := make(map[group.MemberIndex]*big.Int, groupSize-1)

	dishonestThreshold := groupSize - honestThreshold
	coefficientsA, _ := gjkr.GeneratePolynomial(dishonestThreshold)
	coefficientsB, _ := gjkr.GeneratePolynomial(dishonestThreshold)

	for i := 1; i <= groupSize; i++ {
		receiverIndex := group.MemberIndex(i)
		if receiverIndex == senderIndex {
			continue
		}

		// ephemeral key pair - we'll create symmetric key between
		// sender and receiver using this key pair in phase 1 and 2
		// of the protocol
		keyPair, err := ephemeral.GenerateKeyPair()
		if err != nil {
			return nil, err
		}
		ephemeralKeyPairs[receiverIndex] = keyPair

		// shares used in the third phase of the protocol
		// those shares will be encrypted with the symmetric key and
		// used as shares of sender generated for the receiver
		sharesS[receiverIndex] = gjkr.EvaluateMemberShare(receiverIndex, coefficientsA)
		sharesT[receiverIndex] = gjkr.EvaluateMemberShare(receiverIndex, coefficientsB)
	}

	// commitments to the coefficients of the generated polynomial used
	// to evaluate shares for all of the receivers
	// those commitments will be used to alter the MemberCommitmentsMessage
	// sent by original sender
	commitments := make([]*bn256.G1, len(coefficientsA))
	H := altbn128.G1HashToPoint(seed.Bytes())
	for k := range commitments {
		// G * s + H * t
		commitments[k] = new(bn256.G1).Add(
			new(bn256.G1).ScalarBaseMult(coefficientsA[k]),
			new(bn256.G1).ScalarMult(H, coefficientsB[k]),
		)
	}

	// publicKeySharePoints calculated using the coefficients of the generated
	// polynomial used to evaluate shares for all of the receivers
	// those points will be used to alter the MemberPublicKeySharePointsMessage
	// sent by original sender
	publicKeySharePoints := make([]*bn256.G2, len(coefficientsA))
	for i, a := range coefficientsA {
		publicKeySharePoints[i] = new(bn256.G2).ScalarBaseMult(a)
	}

	return &manInTheMiddle{
		senderIndex: senderIndex,
		seed:        seed,

		ephemeralKeyPairs: ephemeralKeyPairs,

		symmetricKeys:      make(map[group.MemberIndex]*ephemeral.SymmetricEcdhKey),
		symmetricKeysMutex: sync.Mutex{},

		sharesS:     sharesS,
		sharesT:     sharesT,
		commitments: commitments,

		publicKeySharePoints: publicKeySharePoints,
	}, nil
}

// interceptCommunication intercepts communication for all phases of DKG protocol.
// Man in the middle sets up symmetric keys, member shares, and commitments.
// It also publishes matching public key share points. Man-in-the-middle drops
// all accusations from the original member.
// This behavior can be overwitten in test before man-in-the-middle intercepts.
func (mitm *manInTheMiddle) interceptCommunication(
	msg net.TaggedMarshaler,
) net.TaggedMarshaler {

	publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
	// Phase 1:
	// Original sender broadcasts EphemeralPublicKeyMessage.
	// We intercept that message and replace the public key generated for
	// each receiver with public key generated earlier by the man in the middle.
	if ok && publicKeyMessage.SenderID() == mitm.senderIndex {
		for receiverIndex, ephemeralKeyPair := range mitm.ephemeralKeyPairs {
			publicKeyMessage.SetPublicKey(
				receiverIndex,
				ephemeralKeyPair.PublicKey,
			)
		}
		return publicKeyMessage
	}
	// Phase 2:
	// The rest of the members generated and broadcast ephemeral public keys.
	// We follow the protocol and perform ECDH against given member's public
	// key generated for the sake of communication with our sender and private
	// ephemeral key generated earlier by the man in the middle.
	if ok && publicKeyMessage.SenderID() != mitm.senderIndex {
		keyPair := mitm.ephemeralKeyPairs[publicKeyMessage.SenderID()]
		symmetricKey := keyPair.PrivateKey.Ecdh(
			publicKeyMessage.GetPublicKey(mitm.senderIndex),
		)

		mitm.symmetricKeysMutex.Lock()
		mitm.symmetricKeys[publicKeyMessage.SenderID()] = symmetricKey
		mitm.symmetricKeysMutex.Unlock()

		return publicKeyMessage
	}

	// Phase 3:
	// Original sender broadcasts PeerSharesMessage and MemberCommitmentsMessage.
	// We intercept those messages and replace shares and commitments with the
	// ones generated by man-in-the-middle.
	// We do that because each receiver established symmetric key with the
	// man-in-the-middle and not with the original sender. Thus, we need to
	// encrypt shares using that symmetric key and regenerate commitments
	// to match the shares.
	peerSharesMessage, ok := msg.(*gjkr.PeerSharesMessage)
	if ok && peerSharesMessage.SenderID() == mitm.senderIndex {
		for receiverIndex, shareS := range mitm.sharesS {
			peerSharesMessage.AddShares(
				receiverIndex,
				shareS,
				mitm.sharesT[receiverIndex],
				mitm.symmetricKeys[receiverIndex],
			)
		}
		return peerSharesMessage
	}
	commitmentsMessage, ok := msg.(*gjkr.MemberCommitmentsMessage)
	if ok && commitmentsMessage.SenderID() == mitm.senderIndex {
		for i, commitment := range mitm.commitments {
			commitmentsMessage.SetCommitment(i, commitment)
		}
		return commitmentsMessage
	}

	// Phase 7:
	// Original sender broadcasts MemberPublicKeySharePointsMessage.
	// We intercept this message and replace public key share points with the
	// ones generated by man-in-the-middle.
	publicKeySharePointsMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
	if ok && publicKeySharePointsMessage.SenderID() == mitm.senderIndex {
		for i, publicKeyShare := range mitm.publicKeySharePoints {
			publicKeySharePointsMessage.SetPublicKeyShare(i, publicKeyShare)
		}
		return publicKeySharePointsMessage
	}

	// Phases 4, 8, 10:
	// By design, man-in-the-middle does not perform any accusations and do not
	// resolve them. All accusation are dropped by default. This behavior may be
	// overwritten by test.
	secretSharesAccusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
	if ok && secretSharesAccusationsMessage.SenderID() == mitm.senderIndex {
		accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
		secretSharesAccusationsMessage.SetAccusedMemberKeys(accusedMembersKeys)
		return secretSharesAccusationsMessage
	}
	pointsAccusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
	if ok && pointsAccusationsMessage.SenderID() == mitm.senderIndex {
		accusedMembersKeys := make(map[group.MemberIndex]*ephemeral.PrivateKey)
		pointsAccusationsMessage.SetAccusedMemberKeys(accusedMembersKeys)
		return pointsAccusationsMessage
	}
	misbehavedKeysMessage, ok := msg.(*gjkr.MisbehavedEphemeralKeysMessage)
	if ok && misbehavedKeysMessage.SenderID() == mitm.senderIndex {
		return nil
	}

	return msg
}
