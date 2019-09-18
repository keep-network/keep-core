/*
  Integration tests for the full DKG affecting GJKR-specific parts.
*/
package gjkr_test

import (
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
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
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertNoInactiveMembers(t, result)
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
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
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
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1), group.MemberIndex(2))
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
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
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
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
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
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 4, 5}...)
}

func TestExecute_IA_members35_phase10(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		disqualifiedKeysMessage, ok := msg.(*gjkr.DisqualifiedEphemeralKeysMessage)
		if ok && (disqualifiedKeysMessage.SenderID() == group.MemberIndex(3) ||
			disqualifiedKeysMessage.SenderID() == group.MemberIndex(5)) {
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
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(3))
	dkgtest.AssertInactiveMembers(t, result)
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
	dkgtest.AssertDisqualifiedMembers(t, result, group.MemberIndex(4))
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 2, 3, 5}...)
}

// TODO Test case Phase 5: 'accuser accuse an inactive member ->
//  expected result: disqualify accuser'.
//  This case is difficult to implement for now because it needs
//  access to member internals. In order to make an accusation against inactive
//  member, there is a need to obtain ephemeral private key for the accused
//  member which is stored in accuser internal map called 'ephemeralKeyPairs'.

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
	dkgtest.AssertSuccessfulSigners(t, result, []group.MemberIndex{2, 3, 5}...)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertDisqualifiedMembers(t, result, []group.MemberIndex{1, 4}...)
	dkgtest.AssertNoInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{2, 3, 5}...)
}

// TODO Test case Phase 9: 'public key share valid ->
//  expected result: disqualify accuser'.
//  This case is difficult to implement for now because it needs
//  access to member internals. In order to make a false accusation
//  there is a need to obtain ephemeral private key for the accused member which
//  is stored in accuser internal map called 'ephemeralKeyPairs'.

// TODO Test case Phase 9: 'accuser accuse an inactive member ->
//  expected result: disqualify accuser'.
//  Public key share points broadcast in the previous phase are necessary to
//  resolve an accusation against the member. Member marked as inactive in any
//  previous phase should not be accused because the accusation can't be resolved.
//  This case is difficult to implement for now because it needs
//  access to member internals. In order to make an accusation against inactive
//  member, there is a need to obtain ephemeral private key for the accused
//  member which is stored in accuser internal map called 'ephemeralKeyPairs'.

// TODO Test case Phase 9: 'cannot decrypt shares ->
//  expected result: disqualify both'.
//  Only happens if the complainer failed to complain earlier
//  and thus both violated the protocol.
//  This case is difficult to implement for now because it needs
//  access to member internals. In order to screw up shares decryption
//  in this phase, there is a need to alter an already received message
//  which is stored in the evidence log.

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

	result, err := dkgtest.RunTest(groupSize, honestThreshold, seed, interceptor)
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

// manInTheMiddle is a helper tool allowing to easily intercept communication
// of a chosen member with the rest of the members, for the first three phases
// of DKG protocol and to set up symmetric keys, member shares, and commitments
// that can be later accessed in test.
//
// In a test not using manInTheMiddle we don't have an access to symmetric keys
// of a chosen member established with the rest of the members. Knowledge of
// these symmetric keys is required to test some scenarios.
type manInTheMiddle struct {
	senderIndex   group.MemberIndex
	receiverIndex group.MemberIndex

	seed *big.Int

	ephemeralKeyPairs map[group.MemberIndex]*ephemeral.KeyPair
	symmetricKeys     map[group.MemberIndex]*ephemeral.SymmetricEcdhKey

	sharesS     map[group.MemberIndex]*big.Int
	sharesT     map[group.MemberIndex]*big.Int
	commitments []*bn256.G1
}

// newManInTheMiddle creates a new instance of manInTheMiddle tool.
// It will intercept messages sent from the sender with the given index and
// modify the part of the message intended for the rest of the members.
// It will intercept the symmetric key handshake as well as generate
// peer shares and commitments based on the established symmetric key.
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

	return &manInTheMiddle{
		senderIndex: senderIndex,
		seed:        seed,

		ephemeralKeyPairs: ephemeralKeyPairs,
		symmetricKeys:     make(map[group.MemberIndex]*ephemeral.SymmetricEcdhKey),

		sharesS:     sharesS,
		sharesT:     sharesT,
		commitments: commitments,
	}, nil
}

// interceptCommunication intercepts the first three phases of DKG protocol
// to set up symmetric keys between a chosen sender and the rest of the members
// such that it can be later accessed and used in test.
// It also intercepts commitments and peer shares messages and modify the
// original values, replacing them with new ones generated based on the
// established (intercepted) symmetric keys.
func (mitm *manInTheMiddle) interceptCommunication(
	msg net.TaggedMarshaler,
) net.TaggedMarshaler {

	publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
	// Act 1:
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
	// Act 2:
	// The rest of the members generated and broadcast ephemeral public keys.
	// We follow the protocol and perform ECDH against given member's public
	// key generated for the sake of communication with our sender and private
	// ephemeral key generated earlier by the man in the middle.
	if ok && publicKeyMessage.SenderID() != mitm.senderIndex {
		mitm.symmetricKeys[publicKeyMessage.SenderID()] =
			mitm.ephemeralKeyPairs[publicKeyMessage.SenderID()].PrivateKey.Ecdh(
				publicKeyMessage.GetPublicKey(mitm.senderIndex),
			)
		return publicKeyMessage
	}

	// Act 3:
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

	return msg
}
