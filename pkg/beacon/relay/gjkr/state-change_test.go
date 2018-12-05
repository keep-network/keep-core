package gjkr

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
)

func TestChallengeStateChagne(t *testing.T) {
	fmt.Printf("Test of func (vc *ChallengeState) ChallengeStateChange( - direct\n")

	type ops struct {
		EventName     string
		resultHash    []byte
		correctResult int
		currentResult result.Result
		finalState    int
	}

	var tests = map[string]struct {
		runTest            bool  //
		groupSize          int   //
		thresholdSize      int   //
		dishonestThreshold int   // T_Max
		votingDuration     int   // T_Conflict
		Ops                []ops //
	}{
		"001 Correct Result is Lead Result, Test Case 00003.": {
			runTest:            true,
			thresholdSize:      8,
			groupSize:          15,
			dishonestThreshold: 8,
			votingDuration:     4,
			Ops: []ops{
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						// Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
			},
			// finalState: 1,
		},
		"002 Correct Result is Lead Result, Test Case 00003 - Verify that successisive tests do not impact each other.": {
			runTest:            true,
			thresholdSize:      8,
			groupSize:          15,
			dishonestThreshold: 8,
			votingDuration:     4,
			Ops: []ops{
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 3,
				},
			},
		},
		"003 Correct Result is Lead Result, First Plus Conflict 00001.": {
			runTest:            true,
			thresholdSize:      8,
			groupSize:          15,
			dishonestThreshold: 2,
			votingDuration:     0,
			Ops: []ops{
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName: "sleep",
				},
				{
					EventName:  "check-state",
					finalState: 3,
				},
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
			},
			// finalState: 3,
		},
		"004 Correct Result is Lead Result, Longer sleep over many blocks.": {
			runTest:            true,
			thresholdSize:      8,
			groupSize:          15,
			dishonestThreshold: 8,
			votingDuration:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 3,
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 1,
				},
			},
		},
		"005 Correct Result is Lead Result, Test voting for stuff.": {
			// } else if m1.AllVotes[leadResult].Votes > vc.dishonestThreshold {
			runTest:            true,
			thresholdSize:      8,
			groupSize:          15,
			dishonestThreshold: 8,
			votingDuration:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 3,
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{EventName: "sleep"},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{EventName: "dump-vc"},
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{EventName: "dump-vc"},
				{
					EventName:  "check-state",
					finalState: 2,
				},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 1,
				},
			},
		},
		"006 Correct Result is Lead Result, submitted a value.": {
			// } else if inResult(correctResult, m1.AllVotes, m1.AllResults) {
			runTest:            true,
			thresholdSize:      8,
			groupSize:          15,
			dishonestThreshold: 8,
			votingDuration:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 3,
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Challenge",
					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{EventName: "dump-vc"},
				{
					EventName:  "check-state",
					finalState: 7,
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						//	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 1,
				},
			},
		},
		"007 Correct Result is Lead Result, Submitted result multiple tims to result in an 'already submitted' state.": {
			// } else if inResult(correctResult, m1.AllVotes, m1.AllResults) {
			// then...
			// check that results that already submitted are handled correctly.
			runTest:            true,
			thresholdSize:      8,
			groupSize:          15,
			dishonestThreshold: 8,
			votingDuration:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						// Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 3,
				},
				{
					EventName: "Challenge",
					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						// Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				// {EventName: "dump-vc"},
				{
					EventName:  "check-state",
					finalState: 7,
				},
				{
					EventName: "sleep",
				},
				{ // test the ability to detec already submitted.
					EventName: "Challenge",
					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						// Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{EventName: "dump-vc"},
				{
					EventName:  "check-state",
					finalState: 5,
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						// Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 1,
				},
			},
		},
		"008 Correct Result is Lead Result, Test the else case.": {
			runTest:            true,
			thresholdSize:      8,
			groupSize:          15,
			dishonestThreshold: 8,
			votingDuration:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						// 	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 3,
				},
				{
					EventName: "Challenge",
					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						// 	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 7,
				},
				{EventName: "sleep"},
				{ // test the ability to detec already submitted.
					EventName: "Challenge",
					resultHash: []byte{139, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						// 	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{139, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{EventName: "dump-vc"},
				{
					EventName:  "check-state",
					finalState: 5,
				},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{ // test the ability to detec already submitted.
					EventName: "Challenge",
					resultHash: []byte{139, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						// 	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{139, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{EventName: "dump-vc"},
				{
					EventName:  "check-state",
					finalState: 1,
				},
				{EventName: "sleep"},
				{
					EventName: "Challenge",
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						// 	Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					finalState: 1,
				},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName[3:], func(t *testing.T) {
			if !test.runTest {
				return
			}
			threshold := 1                /* xyzzy */
			groupSize := 3                /* xyzzy */
			expectedProtocolDuration := 3 // T_dkg	/* xyzzy */
			blockStep := 2                // T_step	/* xyzzy */
			chain, err := initChain(threshold, groupSize, expectedProtocolDuration, blockStep)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}

			fmt.Printf("Running Test %s\n", testName)
			dkg := &DKG{
				chain: chain,
				P:     big.NewInt(100), //
				Q:     big.NewInt(9),   //
			}
			vc, err := dkg.NewChallengeState()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			vc.dishonestThreshold = test.dishonestThreshold // M_Max
			vc.votingDuration = test.votingDuration         // T_conflict
			vc.TFirst = 1                                   // T_First -

			for pc, anOp := range test.Ops {
				switch anOp.EventName {
				case "Challenge":
					vc.ChallengeStateChange(&anOp.currentResult, anOp.resultHash, anOp.correctResult)
				case "Vote":
					vc.VoteForHash(anOp.resultHash)
				case "Block":
					// vc.VoteForhash(anOp.resultHash)
				case "sleep":
					fmt.Printf(";")
					time.Sleep(1 * time.Second)
					fmt.Printf(":")
				case "check-state":
					cs := stateReachedValue
					if anOp.finalState != cs {
						t.Errorf("Did not reach the correct state in pc=%d with test=[%s], expected %d found %d\n", pc, testName, anOp.finalState, cs)
					}
				case "dump-vc":
					fmt.Printf("Test [%s] at end: dkg: %s\n", testName, ConvertToJSON(vc))
				}
			}

			fmt.Printf("Test [%s] at end: dkg: %s\n", testName, ConvertToJSON(vc))
		})
	}

}
