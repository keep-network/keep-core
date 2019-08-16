/*
  Integration tests for the full DKG affecting GJKR-specific parts.
*/
package gjkr_test

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/dkgtest"
	"github.com/keep-network/keep-core/pkg/net"
)

func TestExecute_HappyPath(t *testing.T) {
	t.Parallel()

	groupSize := 5
	threshold := 3

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}

	result, err := dkgtest.RunTest(groupSize, threshold, interceptor)
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

func TestExecute_IA_member1_ephemeralKeyGenerationPhase1(t *testing.T) {
	t.Parallel()

	groupSize := 5
	threshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		publicKeyMessage, ok := msg.(*gjkr.EphemeralPublicKeyMessage)
		if ok && publicKeyMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, threshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
}

func TestExecute_IA_member1and2_commitmentPhase3(t *testing.T) {
	t.Parallel()

	groupSize := 7
	threshold := 4

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

	result, err := dkgtest.RunTest(groupSize, threshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1), group.MemberIndex(2))
	dkgtest.AssertValidGroupPublicKey(t, result)
}

func TestExecute_IA_member1_sharesAndCommitmentsVerificationPhase4(t *testing.T) {
	t.Parallel()

	groupSize := 3
	threshold := 2

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		accusationsMessage, ok := msg.(*gjkr.SecretSharesAccusationsMessage)
		if ok && accusationsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, threshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
}

func TestExecute_IA_member1_publicKeySharePointsCalculationPhase7(t *testing.T) {
	t.Parallel()

	groupSize := 5
	threshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		sharePointsMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && sharePointsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, threshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
}

func TestExecute_IA_member1_publicKeySharePointsVerificationPhase8(t *testing.T) {
	t.Parallel()

	groupSize := 5
	threshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		accusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		if ok && accusationsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, threshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-1)
	dkgtest.AssertMemberFailuresCount(t, result, 1)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(1))
	dkgtest.AssertValidGroupPublicKey(t, result)
}

func TestExecute_IA_member3and5_disqualifiedMembersKeysRevealingPhase10(t *testing.T) {
	t.Parallel()

	groupSize := 5
	threshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		disqualifiedKeysMessage, ok := msg.(*gjkr.DisqualifiedEphemeralKeysMessage)
		if ok && (disqualifiedKeysMessage.SenderID() == group.MemberIndex(3) ||
			disqualifiedKeysMessage.SenderID() == group.MemberIndex(5)) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, threshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize-2)
	dkgtest.AssertMemberFailuresCount(t, result, 2)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result, group.MemberIndex(3), group.MemberIndex(5))
	dkgtest.AssertValidGroupPublicKey(t, result)
}
