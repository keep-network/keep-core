package local

import (
	"math/big"
	"testing"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// TestVote checks that Vote() works - but also uses
// SubmitDKGResult() and GetDKGSubmission().
func TestVote(t *testing.T) {

	group := big.NewInt(40000)
	tests := map[string]struct {
		requestID      *big.Int
		dkgResult      *relaychain.DKGResult
		requestIDset   *big.Int
		groupPublicKey []byte
		expected       int
	}{
		"test increment of votes when match occures": {
			requestID: group,
			dkgResult: &relaychain.DKGResult{
				Success:        true,
				GroupPublicKey: []byte{10, 1},
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
			requestIDset:   group,
			groupPublicKey: []byte{10, 1},
			expected:       2,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			local := Connect(10, 4).ThresholdRelay()
			promise := local.SubmitDKGResult(test.requestID, test.dkgResult)
			_ = promise // in this package promice is fulfilled immediatly - so can ignore it.
			local.DKGResultVote(test.requestID, test.dkgResult.Hash())
			subs := local.GetDKGSubmissions(test.requestID)
			actual := subs.DKGSubmissions[0].Votes
			if test.expected != actual {
				t.Errorf(
					"\nTest: %s\nexpected: [%v]\nactual:   [%v]",
					testName,
					test.expected,
					actual,
				)
			}
		})
	}

}
