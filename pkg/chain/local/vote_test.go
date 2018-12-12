package local

import (
	"math/big"
	"sync"
	"testing"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
)

// TestVote checks that Vote() works - but also uses
// SubmitDKGResult() and GetDKGSubmission().
func TestVote(t *testing.T) {

	group := big.NewInt(40000)
	tests := map[string]struct {
		requestID      *big.Int
		dkgResult      *relaychain.DKGResult
		requestIDset   *big.Int
		groupPublicKey *big.Int
		expected       int
	}{
		"test increment of votes when match occures": {
			requestID: group,
			dkgResult: &relaychain.DKGResult{
				Success:        true,
				GroupPublicKey: big.NewInt(1001),
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
			requestIDset:   group,
			groupPublicKey: big.NewInt(1001),
			expected:       2,
		},
		"test increment of votes when no match occures": {
			requestID: group,
			dkgResult: &relaychain.DKGResult{
				Success:        true,
				GroupPublicKey: big.NewInt(1002),
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
			requestIDset:   group,
			groupPublicKey: big.NewInt(1001),
			expected:       1,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			local := ConnectForTest(10, 4)
			local.groupPublicKeyMap[bigIntToHex(test.requestIDset)] = test.groupPublicKey
			promise := local.SubmitDKGResult(test.requestID, test.dkgResult)
			_ = promise // in this package promice is fulfilled immediatly - so can ignore it.
			local.Vote(test.requestID, test.dkgResult.Hash())
			subs := local.GetDKGSubmissions(test.requestID)
			actual := subs.Submissions[0].Votes
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

// Connect initializes a local stub implementation of the chain interfaces
// for testing.
func ConnectForTest(groupSize int, threshold int) *localChain {
	bc, _ := blockCounter()

	return &localChain{
		relayConfig: relayconfig.Chain{
			GroupSize: groupSize,
			Threshold: threshold,
		},
		groupRegistrationsMutex: sync.Mutex{},
		groupRelayEntries:       make(map[string]*big.Int),
		groupRegistrations:      make(map[string][96]byte),
		submittedResults:        make(map[string][]*relaychain.DKGResult),
		blockCounter:            bc,
		stakeMonitor:            NewStakeMonitor(),
		submissions:             make(map[string]relaychain.Submissions),
		groupPublicKeyMap:       make(map[string]*big.Int),
	}
}
