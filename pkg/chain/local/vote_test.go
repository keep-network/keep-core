package local

import (
	"fmt"
	"math/big"
	"testing"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/pschlump/MiscLib"
)

func TestVote(t *testing.T) {

	tests := map[string]struct {
		requestID *big.Int
		dkgResult *relaychain.DKGResult
		expected  int
	}{
		"test increment of votes when match occures": {
			requestID: big.NewInt(40000),
			dkgResult: &relaychain.DKGResult{
				Success:        true,
				GroupPublicKey: big.NewInt(1001),
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
			expected: 2,
		},
	}

	local := Connect(10, 4).ThresholdRelay()

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			promise := local.SubmitDKGResult(test.requestID, test.dkgResult)
			_ = promise
			local.Vote(test.requestID, test.dkgResult.Hash())
			subs := local.GetDKGSubmissions(test.requestID)
			actual := subs.Submissions[0].Votes
			if test.expected != actual {
				fmt.Printf("%s\tError: %+v%s\n", MiscLib.ColorRed, subs, MiscLib.ColorReset)
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					expected,
					actual,
				)
			}
		})
	}

}
