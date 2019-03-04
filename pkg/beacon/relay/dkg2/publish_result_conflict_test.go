package dkg2

import (
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// This tests executes result conflict resolution by a single member and validates
// simple phase flow.
func TestResultConflictResolution_SingleMemberTests(t *testing.T) {
	groupSize := 10         // N
	dishonestThreshold := 4 // M

	publishingIndex := 1
	blockStep := uint64(2)
	conflictDuration := uint64(10)

	dkgResult1 := &relayChain.DKGResult{
		GroupPublicKey: []byte{100},
	}

	dkgResult2 := &relayChain.DKGResult{
		GroupPublicKey: []byte{200},
	}

	var chain relayChain.Interface

	var tests = map[string]struct {
		correctResult *relayChain.DKGResult
		prerequisites func(requestID *big.Int)

		expectedBlockNumber   int64
		expectedDuration      uint64
		expectedError         error
		submissionsValidation func(submission dkgResultsVotes)
	}{
		"run result conflict resolution when there are no any previous submissions": {
			correctResult: dkgResult1,
			prerequisites: func(requestID *big.Int) {
			},
			expectedBlockNumber: -1,
			expectedDuration:    0,
			expectedError:       fmt.Errorf("nothing submitted"),
		},
		"run result conflict resolution when there are previous submissions but the correct result is not submitted yet": {
			correctResult: dkgResult1,
			prerequisites: func(requestID *big.Int) {
				chain.SubmitDKGResult(requestID, dkgResult2).OnFailure(func(err error) {
					t.Fatal(err)
				}) // Result 1: Votes = 0 | Result 2: Votes = 1
			},
			expectedBlockNumber: int64(blockStep + 1),
			expectedDuration:    conflictDuration,
			expectedError:       nil,
			submissionsValidation: func(submissions dkgResultsVotes) {
				dkgResult1Hash, _ := chain.CalculateDKGResultHash(dkgResult1)
				if !submissions.contains(dkgResult1Hash) {
					t.Error("submissions should contain the result")
				}
			},
		},
		"run result conflict resolution when the leading result has majority of votes": {
			correctResult: dkgResult1,
			prerequisites: func(requestID *big.Int) {
				dkgResult2Hash, _ := chain.CalculateDKGResultHash(dkgResult2)

				chain.SubmitDKGResult(requestID, dkgResult1).
					OnFailure(func(err error) {
						t.Fatal(err)

					}) // Result 1: Votes = 1 | Result 2: Votes = 0
				chain.SubmitDKGResult(requestID, dkgResult2).
					OnFailure(func(err error) {
						t.Fatal(err)
					}) // Result 1: Votes = 1 | Result 2: Votes = 1
				chain.VoteOnDKGResult(requestID, 91, dkgResult2Hash).
					OnFailure(func(err error) {
						t.Fatal(err)
					}) // Result 1: Votes = 1 | Result 2: Votes = 2
				chain.VoteOnDKGResult(requestID, 92, dkgResult2Hash).
					OnFailure(func(err error) {
						t.Fatal(err)
					}) // Result 1: Votes = 1 | Result 2: Votes = 3
				chain.VoteOnDKGResult(requestID, 93, dkgResult2Hash).
					OnFailure(func(err error) {
						t.Fatal(err)
					}) // Result 1: Votes = 1 | Result 2: Votes = 4
				chain.VoteOnDKGResult(requestID, 94, dkgResult2Hash).
					OnFailure(func(err error) {
						t.Fatal(err)
					}) // Result 1: Votes = 1 | Result 2: Votes = 5
			},
			expectedBlockNumber: -1,
			expectedDuration:    0,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			requestID := big.NewInt(1)

			chainHandle, blockCounter, initialBlock, err := initChainHandle(dishonestThreshold, groupSize)
			if err != nil {
				t.Fatal(err)
			}
			chain = chainHandle.ThresholdRelay()

			publisher := &Publisher{
				RequestID:          requestID,
				publishingIndex:    publishingIndex,
				dishonestThreshold: dishonestThreshold,
				blockStep:          blockStep,
				blockCounter:       blockCounter,
				conflictDuration:   conflictDuration,
			}

			if test.prerequisites != nil {
				test.prerequisites(requestID)
			}

			// Wait for `blockStep` to simulate Phase 13.
			blockStartWaiter, err := blockCounter.BlockHeightWaiter(int(uint64(initialBlock) + blockStep))
			if err != nil {
				t.Fatal(err)
			}

			// TEST
			blockStart := uint64(<-blockStartWaiter)

			blockNumber, actualError := publisher.resultConflictResolution(
				dkgResult1,
				chain,
				blockStart,
			)

			blockEnd, err := blockCounter.CurrentBlock()
			if err != nil {
				t.Fatal(err)
			}
			duration := uint64(blockEnd) - blockStart

			// VALIDATIONS
			if !reflect.DeepEqual(actualError, test.expectedError) {
				t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedError, actualError)
			}

			if duration != test.expectedDuration {
				t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedDuration, duration)
			}

			if blockNumber != int64(test.expectedBlockNumber) {
				t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedBlockNumber, blockNumber)
			}

			if test.submissionsValidation != nil {
				dkgResultsVotes := dkgResultsVotes(chain.GetDKGResultsVotes(requestID))
				test.submissionsValidation(dkgResultsVotes)
			}
		})
	}
}

// This test simulates concurrent execution of the result conflict resolution by
// multiple members. Each member supports a specific result and starts the phase
// at different time defined by block number. It validates following scenario:
// Block|  Votes   | Comment
//      | R1 R2 R3 | R = Result, P = Player
//      | 1  0  0  | Initial state, R1 already submitted
//    0 | 2  0  0  | P1 joins and votes for R1
//    1 | 2  1  0  | P2 joins and P2 submitts R2
//    3 | 2  1  0  | P3 joins, result R1 is already leading, so no action required
//    4 | 3  2  0  | P4 joins and votes for R2 [2 2 0] -> P3 votes for R1 [3 2 0]
//    6 | 3  2  0  | P5 joins and votes for R3
func TestResultConflictResolution_Concurrent(t *testing.T) {
	type TestStats struct {
		blockNumber int64
		duration    uint64
	}

	groupSize := 10         // N
	dishonestThreshold := 6 // M
	blockStep := uint64(2)
	conflictDuration := uint64(10)

	requestID := big.NewInt(1)
	dkgResult1 := &relayChain.DKGResult{GroupPublicKey: []byte{10}}
	dkgResult2 := &relayChain.DKGResult{GroupPublicKey: []byte{20}}
	dkgResult3 := &relayChain.DKGResult{GroupPublicKey: []byte{30}}

	// Test steps configured in a map, where key is a publisher's index.
	var testSteps = map[int]struct {
		startOnBlock  uint64 // Simulates different times when members joins the execution.
		correctResult *relayChain.DKGResult

		expectedDuration uint64
	}{
		1: {
			startOnBlock:     0,
			correctResult:    dkgResult1,
			expectedDuration: conflictDuration,
		},
		2: {
			startOnBlock:     1,
			correctResult:    dkgResult2,
			expectedDuration: conflictDuration - 1,
		},
		3: {
			startOnBlock:     3,
			correctResult:    dkgResult1,
			expectedDuration: conflictDuration - 3,
		},
		4: {
			startOnBlock:     4,
			correctResult:    dkgResult2,
			expectedDuration: conflictDuration - 4,
		},
		5: {
			startOnBlock:     6,
			correctResult:    dkgResult3,
			expectedDuration: conflictDuration - 6,
		},
	}
	expectedBlockNumber := int64(blockStep + 7)

	chainHandle, blockCounter, _, err := initChainHandle(dishonestThreshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}
	chain := chainHandle.ThresholdRelay()

	dkgResult1Hash, _ := chain.CalculateDKGResultHash(dkgResult1)
	dkgResult2Hash, _ := chain.CalculateDKGResultHash(dkgResult2)
	dkgResult3Hash, _ := chain.CalculateDKGResultHash(dkgResult3)

	expectedDKGResultsVotes := relayChain.DKGResultsVotes{
		dkgResult1Hash: 3,
		dkgResult2Hash: 2,
		dkgResult3Hash: 1,
	}

	// waitGroup is used to find when tests executed in goroutines completes.
	var waitGroup sync.WaitGroup
	var initialBlock uint64

	conflictResolutionTest := func(
		publishingIndex int,
		startOnBlock uint64,
		correctDKGResult *relayChain.DKGResult,
		testStats *sync.Map,
	) {
		defer waitGroup.Done()

		publisher := &Publisher{
			RequestID:          requestID,
			publishingIndex:    publishingIndex,
			dishonestThreshold: dishonestThreshold,
			blockStep:          blockStep,
			blockCounter:       blockCounter,
			conflictDuration:   conflictDuration,
		}

		// Simulate publishers to join at different times.
		startBlock := initialBlock + startOnBlock
		blockCounter.WaitForBlockHeight(int(startBlock))
		if err != nil {
			t.Error(err)
		}

		blockNumber, err := publisher.resultConflictResolution(
			correctDKGResult,
			chain,
			initialBlock,
		)
		if err != nil {
			t.Error(err)
		}
		endBlock, err := blockCounter.CurrentBlock()
		if err != nil {
			t.Fatal(err)
		}

		testStats.Store(publishingIndex, &TestStats{blockNumber, uint64(endBlock) - startBlock})
	}

	// Simulate Phase 13 - submit initial DKG result and wait some blocks.
	// TODO: When SubmitDKGResult supports member index set it to `0`
	chain.SubmitDKGResult(requestID, dkgResult1).OnFailure(func(err error) {
		t.Fatal(err)
	})
	startBlockWaiter, err := blockCounter.BlockWaiter(int(blockStep))
	initialBlock = uint64(<-startBlockWaiter)

	// Tests
	waitGroup.Add(len(testSteps))
	actualTestStats := &sync.Map{} //stores each test's execution results
	for publishingIndex, test := range testSteps {
		go conflictResolutionTest(
			publishingIndex,
			test.startOnBlock,
			test.correctResult,
			actualTestStats,
		)
	}
	waitGroup.Wait() // waits for tests execution completion

	// Validations
	actualTestStats.Range(func(publishingIndex, actualStats interface{}) bool {
		if !reflect.DeepEqual(
			actualStats.(*TestStats).blockNumber,
			expectedBlockNumber,
		) {
			t.Errorf(
				"invalid block number for publisher %+v\nexpected: %+v\nactual:   %v\n",
				publishingIndex,
				expectedBlockNumber,
				actualStats.(*TestStats).blockNumber,
			)
		}
		if !reflect.DeepEqual(
			actualStats.(*TestStats).duration,
			testSteps[publishingIndex.(int)].expectedDuration,
		) {
			t.Errorf(
				"invalid duration for publisher %+v\nexpected: %+v\nactual:   %v\n",
				publishingIndex,
				testSteps[publishingIndex.(int)].expectedDuration,
				actualStats.(*TestStats).duration,
			)
		}

		return true
	})

	dkgResultsVotes := dkgResultsVotes(chain.GetDKGResultsVotes(requestID))
	for dkgResultHash, votes := range dkgResultsVotes {
		if dkgResultsVotes[dkgResultHash] != votes {
			t.Errorf("invalid votes for hash [%x]\nexpected: %v\nactual:   %v\n", dkgResultHash, expectedDKGResultsVotes, dkgResultsVotes)
		}
	}
}
