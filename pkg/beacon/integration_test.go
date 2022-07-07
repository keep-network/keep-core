// Package relay_test contains integration tests for the whole random beacon
// roundtrip including DKG and threshold signing.
package beacon_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/entry"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/group"
	"github.com/keep-network/keep-core/pkg/bls"

	"github.com/keep-network/keep-core/pkg/internal/dkgtest"
	"github.com/keep-network/keep-core/pkg/internal/entrytest"
	"github.com/keep-network/keep-core/pkg/net"
)

const groupSize = 10
const honestThreshold = 6

func previousEntry() []byte {
	return previousEntryG1().Marshal()
}

func previousEntryG1() *bn256.G1 {
	return new(bn256.G1).ScalarBaseMult(big.NewInt(1328472189))
}

// Success: all members of the signing group participate in signing.
func TestAllMembersSigning(t *testing.T) {
	t.Parallel()

	signingMembersCount := groupSize
	dkgResult, signingResult := runTest(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryPublished(t, signingResult)
	entrytest.AssertNoSignerFailures(t, signingResult)

	groupPublicKey, err := getFirstGroupPublicKey(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	newEntry, err := signingResult.EntryValue()
	if err != nil {
		t.Fatal(err)
	}

	if !bls.VerifyG1(groupPublicKey, previousEntryG1(), newEntry) {
		t.Errorf("threshold signature failed BLS verification")
	}
}

// Success: honest threshold of the signing group members participate in
// signing.
func TestHonestThresholdMembersSigning(t *testing.T) {
	t.Parallel()

	signingMembersCount := honestThreshold
	dkgResult, signingResult := runTest(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryPublished(t, signingResult)
	entrytest.AssertNoSignerFailures(t, signingResult)

	groupPublicKey, err := getFirstGroupPublicKey(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	newEntry, err := signingResult.EntryValue()
	if err != nil {
		t.Fatal(err)
	}

	if !bls.VerifyG1(groupPublicKey, previousEntryG1(), newEntry) {
		t.Errorf("threshold signature failed BLS verification")
	}
}

// Failure: Less than honest threshold signing group members participate in
// signing.
func TestLessThanHonestThresholdMembersSigning(t *testing.T) {
	t.Parallel()

	signingMembersCount := honestThreshold - 1
	dkgResult, signingResult := runTest(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryNotPublished(t, signingResult)
	entrytest.AssertSignerFailuresCount(t, signingResult, signingMembersCount)
}

// Success: honest threshold of the signing group members participate in
// signing.
//
// In this scenario, one of the members doesn't send `MemberPublicKeySharePointsMessage`
// thus they become inactive at the beginning of phase 8 during DKG.
// This is problematic because that member provided valid shares in phase 3
// and all group members include that shares in their private key shares.
// Since that member did not provide public key share points, shares from that
// member are not included in the information we use to calculate the public key
// of the group. If we do not reconstruct and include shares of that member,
// we may end up with a situation when a signature does not match the
// group public key.
func TestInactiveMemberPublicKeySharesReconstructionAndSigning(t *testing.T) {
	t.Parallel()

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		sharePointsMessage, ok := msg.(*gjkr.MemberPublicKeySharePointsMessage)
		if ok && sharePointsMessage.SenderID() == group.MemberIndex(1) {
			return nil
		}

		return msg
	}

	signingMembersCount := honestThreshold
	dkgResult, signingResult := runTestWithInterceptor(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
		interceptor,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryPublished(t, signingResult)
	entrytest.AssertNoSignerFailures(t, signingResult)

	groupPublicKey, err := getFirstGroupPublicKey(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	newEntry, err := signingResult.EntryValue()
	if err != nil {
		t.Fatal(err)
	}

	if !bls.VerifyG1(groupPublicKey, previousEntryG1(), newEntry) {
		t.Errorf("threshold signature failed BLS verification")
	}
}

// Success: honest threshold of the signing group members participate in
// signing.
//
// In this scenario, one of the members doesn't send `PointsAccusationsMessage`
// thus they become inactive at the beginning of phase 9 during DKG.
// This is problematic because that member provided valid shares in phase 3
// and all group members include that shares in their private key shares.
// What is even more important, that member also provided valid public key share
// points during the phase 8 of DKG so other members have to decide whether they
// use those received public key share points or they qualify the inactive member
// to the reconstruction phase and reconstruct them on their own.
// Otherwise, we may end up with a situation when a signature does not match the
// group public key.
func TestInactivePointsAccusationsReconstructionAndSigning(t *testing.T) {
	t.Parallel()

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		pointsAccusationsMessage, ok := msg.(*gjkr.PointsAccusationsMessage)
		if ok && pointsAccusationsMessage.SenderID() == group.MemberIndex(3) {
			return nil
		}

		return msg
	}

	signingMembersCount := honestThreshold
	dkgResult, signingResult := runTestWithInterceptor(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
		interceptor,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryPublished(t, signingResult)
	entrytest.AssertNoSignerFailures(t, signingResult)

	groupPublicKey, err := getFirstGroupPublicKey(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	newEntry, err := signingResult.EntryValue()
	if err != nil {
		t.Fatal(err)
	}

	if !bls.VerifyG1(groupPublicKey, previousEntryG1(), newEntry) {
		t.Errorf("threshold signature failed BLS verification")
	}
}

func TestSigningWithInvalidSignatureShare(t *testing.T) {
	t.Parallel()

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		signatureShareMessage, ok := msg.(*entry.SignatureShareMessage)

		// Member 1 sends shares which could not be unmarshalled as a G1 point.
		if ok && signatureShareMessage.SenderID() == group.MemberIndex(1) {
			return entry.NewSignatureShareMessage(
				signatureShareMessage.SenderID(),
				[]byte{0, 1},
			)
		}

		// Member 2 sends a proper G1 point which is invalid in terms of
		// the current relay entry request.
		if ok && signatureShareMessage.SenderID() == group.MemberIndex(2) {
			_, randomG1, err := bn256.RandomG1(rand.Reader)
			if err != nil {
				t.Fatal(err)
			}

			return entry.NewSignatureShareMessage(
				signatureShareMessage.SenderID(),
				randomG1.Marshal(),
			)
		}

		return msg
	}

	signingMembersCount := groupSize
	dkgResult, signingResult := runTestWithInterceptor(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
		interceptor,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryPublished(t, signingResult)
	entrytest.AssertNoSignerFailures(t, signingResult)

	groupPublicKey, err := getFirstGroupPublicKey(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	newEntry, err := signingResult.EntryValue()
	if err != nil {
		t.Fatal(err)
	}

	if !bls.VerifyG1(groupPublicKey, previousEntryG1(), newEntry) {
		t.Errorf("threshold signature failed BLS verification")
	}
}

func runTest(t *testing.T, groupSize, honestThreshold, honestSignersCount int) (
	*dkgtest.Result,
	*entrytest.Result,
) {
	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}

	return runTestWithInterceptor(
		t,
		groupSize,
		honestThreshold,
		honestSignersCount,
		interceptor,
	)
}

func runTestWithInterceptor(
	t *testing.T,
	groupSize, honestThreshold, honestSignersCount int,
	interceptor func(msg net.TaggedMarshaler) net.TaggedMarshaler,
) (
	*dkgtest.Result,
	*entrytest.Result,
) {
	dkgSeed := dkgtest.RandomSeed(t)
	dkgResult, err := dkgtest.RunTest(groupSize, honestThreshold, dkgSeed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	signers := dkgResult.GetSigners()[0:honestSignersCount]

	signingResult, err := entrytest.RunTest(
		signers,
		honestThreshold,
		interceptor,
		previousEntry(),
	)
	if err != nil {
		t.Fatal(err)
	}

	return dkgResult, signingResult
}

func getFirstGroupPublicKey(result *dkgtest.Result) (*bn256.G2, error) {
	signers := result.GetSigners()
	if len(signers) == 0 {
		return nil, fmt.Errorf("no signers in result")
	}

	publicKey := new(bn256.G2)
	_, err := publicKey.Unmarshal(signers[0].GroupPublicKeyBytes())
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
