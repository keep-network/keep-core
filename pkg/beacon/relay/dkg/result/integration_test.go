// Package result_test contains integration tests for the full roundtrip of
// result-publication-specific parts of DKG.
package result_test

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg/result"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/dkgtest"
	"github.com/keep-network/keep-core/pkg/net"
)

func TestExecute_IA_members24_phase13(t *testing.T) {
	t.Parallel()

	groupSize := 5
	honestThreshold := 3
	seed := dkgtest.RandomSeed(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		hashSignatureMessage, ok := msg.(*result.DKGResultHashSignatureMessage)
		if ok && (hashSignatureMessage.SenderID() == group.MemberIndex(2) ||
			hashSignatureMessage.SenderID() == group.MemberIndex(4)) {
			return nil
		}

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
	dkgtest.AssertResultSupportingMembers(t, result, []group.MemberIndex{1, 3, 5}...)
}
