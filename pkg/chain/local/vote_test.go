package local

import (
	//	"context"
	//	"math/big"
	"fmt"
	"math/big"
	"testing"
	//	"time"
	//	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

func TestVote(t *testing.T) {

	// func (c *localChain) Vote(requestID *big.Int, dkgResultHash []byte) {
	tests := map[string]struct {
		requestID *big.Int
		dkgResult *relaychain.DKGResult
	}{
		"test increment of votes when match occures": {
			requestID: big.NewInt(40000),
			dkgResult: &relaychain.DKGResult{
				Success:        true,
				GroupPublicKey: big.NewInt(1001),
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
		},
	}

	local := Connect(10, 4).ThresholdRelay()

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			fmt.Printf("%sTest %+v %s\n", MiscLib.ColorGreen, test, MiscLib.ColorReset)
			promise := local.SubmitDKGResult(test.requestID, test.dkgResult)
			_ = promise
			// requestID *big.Int, resultToPublish *relaychain.DKGResult,
			local.Vote(test.requestID, test.dkgResult.Hash())
			subs := local.GetDKGSubmissions(test.requestID)
			fmt.Printf("%s\tSubmissions: %+v%s\n", MiscLib.ColorCyan, subs, MiscLib.ColorReset)
			fmt.Printf("%s\tVotes: %+v%s\n", MiscLib.ColorCyan, godebug.SVarI(subs), MiscLib.ColorReset)
			if subs.Submissions[0].Votes != 2 {
				fmt.Printf("%s\tError: %+v%s\n", MiscLib.ColorRed, subs, MiscLib.ColorReset)
			}
		})
	}

}
