package gjkr

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

func TestChallengeStateChagne(t *testing.T) {
	fmt.Printf("Test of func (vc *DKG) ChallengeStateChange(EventName string, group *big.Int, resultHash []byte, correctResult int) {\n")
	// func (vc *DKG) ChallengeStateChange(EventName string, group *big.Int, resultHash []byte, correctResult int) {

	type ops struct {
		// vc.ChallengeStateChange(test.EventName, test.Group, &test.currentResult, test.resultHash, test.correctResult)
		EventName     string
		Group         *big.Int
		resultHash    []byte
		correctResult int
		currentResult result.Result
		finalState    int
	}

	var tests = map[string]struct {
		runTest       bool
		groupSize     int
		thresholdSize int
		TMax          int
		TConflict     int
		Ops           []ops
	}{
		"001 Correct Result is Lead Result, Test Case 00003.": {
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			TMax:          8,
			TConflict:     4,
			Ops: []ops{
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			TMax:          8,
			TConflict:     4,
			Ops: []ops{
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 3,
				},
			},
		},
		"003 Correct Result is Lead Result, First Plus Conflict 00001.": {
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			TMax:          2,
			TConflict:     0,
			Ops: []ops{
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
					Group:      big.NewInt(10000),
					finalState: 3,
				},
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			TMax:          8,
			TConflict:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
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
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 1,
				},
			},
		},
		"005 Correct Result is Lead Result, Test voting for stuff.": {
			// } else if m1.AllVotes[leadResult].Votes > vc.TMax {
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			TMax:          8,
			TConflict:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 3,
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{EventName: "sleep"},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{EventName: "dump-vc"},
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
					Group:      big.NewInt(10000),
					finalState: 2,
				},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{EventName: "sleep"},
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 1,
				},
			},
		},
		"006 Correct Result is Lead Result, submitted a value.": {
			// } else if inResult(correctResult, m1.AllVotes, m1.AllResults) {
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			TMax:          8,
			TConflict:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 3,
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Vote",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
					Group:      big.NewInt(10000),
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
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 1,
				},
			},
		},
		"007 Correct Result is Lead Result, Submitted result multiple tims to result in an 'already submitted' state.": {
			// } else if inResult(correctResult, m1.AllVotes, m1.AllResults) {
			// then...
			// check that results that already submitted are handled correctly.
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			TMax:          8,
			TConflict:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 3,
				},
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
					Group:      big.NewInt(10000),
					finalState: 7,
				},
				{
					EventName: "sleep",
				},
				{ // test the ability to detec already submitted.
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
					Group:      big.NewInt(10000),
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
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 1,
				},
			},
		},
		"008 Correct Result is Lead Result, Test the else case.": {
			// Not working yet  - xyzzy - TODO
			runTest:       true,
			thresholdSize: 8,
			groupSize:     15,
			TMax:          8,
			TConflict:     6,
			Ops: []ops{
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 3,
				},
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
					finalState: 7,
				},
				{EventName: "sleep"},
				{ // test the ability to detec already submitted.
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{139, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
					Group:      big.NewInt(10000),
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
				{ // test the ability to detec already submitted.
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{139, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 1,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
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
					Group:      big.NewInt(10000),
					finalState: 1,
				},
				{
					EventName: "sleep",
				},
				{
					EventName: "Challenge",
					Group:     big.NewInt(10000),
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					correctResult: 0,
					currentResult: result.Result{
						Type:           result.PerfectSuccess,
						GroupPublicKey: big.NewInt(10000),
						Disqualified:   []int{},
						Inactive:       []int{},
						HashValue: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
							51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
					},
				},
				{
					EventName:  "check-state",
					Group:      big.NewInt(10000),
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
			fmt.Printf("%sRunning Test %s%s\n", MiscLib.ColorCyan, testName, MiscLib.ColorReset)
			var chainRelay chain.Interface
			var err error
			var vc *DKG
			var members []*PublishingMember
			var member *PublishingMember
			members, err = initializePublishingMembersGroup2(test.thresholdSize, test.groupSize)
			if err != nil {
				t.Errorf("%s", err)
			}
			member = members[0] // WHY? Fix This FIXME xyzzy
			chainRelay = member.protocolConfig.chain.ThresholdRelay()
			_ = chainRelay // xyzzy FIXME
			err = setup2CurrentBlockHeight(member.protocolConfig.chain)
			if err != nil {
				t.Fatalf("Failed to setup chain: %v", err)
			}
			vc = &DKG{
				P:         big.NewInt(100), //
				Q:         big.NewInt(9),   //
				TMax:      test.TMax,       // M_Max
				TConflict: test.TConflict,  // T_conflict
				TNow:      1,               // T_Now - Current block height at the current time
				TFirst:    1,               // T_First -
			}
			vc.PerGroup = make(map[string]ValidationGroupState)

			for pc, anOp := range test.Ops {
				switch anOp.EventName {
				case "Challenge", "Vote":
					vc.ChallengeStateChange(anOp.EventName, anOp.Group, &anOp.currentResult, anOp.resultHash, anOp.correctResult)
				case "Block":
					vc.ChallengeStateChange(anOp.EventName, nil, &anOp.currentResult, anOp.resultHash, anOp.correctResult)
				case "sleep":
					fmt.Printf(";")
					time.Sleep(1 * time.Second)
					fmt.Printf(":")
				case "check-state":
					groupAsKey := fmt.Sprintf("%s", anOp.Group)
					cs := vc.PerGroup[groupAsKey].CurrentState
					if anOp.finalState != cs {
						t.Errorf("Did not reach the correct state in pc=%d with test=[%s], expected %d found %d\n", pc, testName, anOp.finalState, cs)
					}
				case "dump-vc":
					fmt.Printf("Test [%s] at end: dkg: %s\n", testName, godebug.SVarI(vc))
				}
			}

			fmt.Printf("Test [%s] at end: dkg: %s\n", testName, godebug.SVarI(vc))
		})
	}

}
