/*
  Integration tests for the full DKG affecting result publication parts.
*/
package result_test

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg/result"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/dkgtest"
	"github.com/keep-network/keep-core/pkg/net"
)

func TestExecute_IA_member2and4_DKGResultSigningPhase13(t *testing.T) {
	t.Parallel()

	groupSize := 5
	threshold := 3

	interceptorRules := func(msg net.TaggedMarshaler) net.TaggedMarshaler {

		hashSignatureMessage, ok := msg.(*result.DKGResultHashSignatureMessage)
		if ok && (hashSignatureMessage.SenderID() == group.MemberIndex(2) ||
			hashSignatureMessage.SenderID() == group.MemberIndex(4)) {
			return nil
		}

		return msg
	}

	result, err := dkgtest.RunTest(groupSize, threshold, interceptorRules)
	if err != nil {
		t.Fatal(err)
	}

	dkgtest.AssertDkgResultPublished(t, result)
	dkgtest.AssertSuccessfulSignersCount(t, result, groupSize)
	dkgtest.AssertMemberFailuresCount(t, result, 0)
	dkgtest.AssertSamePublicKey(t, result)
	dkgtest.AssertNoDisqualifiedMembers(t, result)
	dkgtest.AssertInactiveMembers(t, result)
	dkgtest.AssertValidGroupPublicKey(t, result)
}
