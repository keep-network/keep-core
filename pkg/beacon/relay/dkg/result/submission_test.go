package result

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/operator"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

func TestSubmitDKGResult(t *testing.T) {
	threshold := 2
	groupSize := 5

	chainHandle, initialBlock, err := initChainHandle(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	config, err := chainHandle.ThresholdRelay().GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	result := &relayChain.DKGResult{
		GroupPublicKey: []byte{123, 45},
	}
	signatures := map[group.MemberIndex]operator.Signature{
		1: operator.Signature{101},
		2: operator.Signature{102},
		3: operator.Signature{103},
	}

	tStep := config.ResultPublicationBlockStep

	var tests = map[string]struct {
		memberIndex     int
		expectedTimeEnd uint64
	}{
		"first member eligible to submit straight away": {
			memberIndex:     1,
			expectedTimeEnd: initialBlock, // T_now < T_init + T_step
		},
		"second member eligible to submit after T_step block passed": {
			memberIndex:     2,
			expectedTimeEnd: initialBlock + tStep, // T_now = T_init + T_step
		},
		"fourth member eligable to submit after T_dkg + 2*T_step passed": {
			memberIndex:     4,
			expectedTimeEnd: initialBlock + 3*tStep, // T_now = T_init + 3*T_step
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			requestID := big.NewInt(101)

			member := &SubmittingMember{
				index: group.MemberIndex(test.memberIndex),
			}

			// Reinitialize chain to reset block counter
			chainHandle, _, err := initChainHandle(threshold, groupSize)
			if err != nil {
				t.Fatalf("chain initialization failed [%v]", err)
			}

			relayChain := chainHandle.ThresholdRelay()

			isSubmitted, err := relayChain.IsDKGResultSubmitted(requestID)
			if err != nil {
				t.Fatal(err)
			}

			if isSubmitted {
				t.Fatalf("result is already submitted to the chain")
			}

			blockCounter, _ := chainHandle.BlockCounter()

			err = member.SubmitDKGResult(
				requestID,
				result,
				signatures,
				relayChain,
				blockCounter,
			)
			if err != nil {
				t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
			}

			currentBlock, _ := blockCounter.CurrentBlock()
			if currentBlock < test.expectedTimeEnd {
				t.Errorf(
					"invalid current block\nexpected: >= %v\nactual:      %v\n",
					test.expectedTimeEnd,
					currentBlock,
				)
			}
			isSubmitted, err = relayChain.IsDKGResultSubmitted(requestID)
			if err != nil {
				t.Fatal(err)
			}
			if !isSubmitted {
				t.Error("result is not submitted to the chain")
			}
		})
	}
}

// This tests runs result publication concurrently by two members.
// Member with lower index gets to publish the result to chain. For the second
// member loop should be aborted and result published by the first member should
// be returned.
func TestConcurrentPublishResult(t *testing.T) {
	threshold := 2
	groupSize := 5

	member1 := &SubmittingMember{
		index: group.MemberIndex(1), // P1
	}
	member2 := &SubmittingMember{
		index: group.MemberIndex(4), // P4
	}

	signatures := map[group.MemberIndex]operator.Signature{
		1: operator.Signature{101},
		2: operator.Signature{102},
		3: operator.Signature{103},
	}

	var tests = map[string]struct {
		resultToPublish1  *relayChain.DKGResult
		resultToPublish2  *relayChain.DKGResult
		requestID1        *big.Int
		requestID2        *big.Int
		expectedDuration1 func(tStep uint64) uint64 // index * t_step
		expectedDuration2 func(tStep uint64) uint64 // index * t_step
	}{
		"two members publish the same results": {
			resultToPublish1: &relayChain.DKGResult{
				GroupPublicKey: []byte{101},
			},
			resultToPublish2: &relayChain.DKGResult{
				GroupPublicKey: []byte{101},
			},
			requestID1:        big.NewInt(11),
			requestID2:        big.NewInt(11),
			expectedDuration1: func(tStep uint64) uint64 { return 0 }, // (P1-1) * t_step
			expectedDuration2: func(tStep uint64) uint64 { return 0 }, // result already published by member 1 -1
		},
		"two members publish different results": {
			resultToPublish1: &relayChain.DKGResult{
				GroupPublicKey: []byte{201},
			},
			resultToPublish2: &relayChain.DKGResult{
				GroupPublicKey: []byte{202},
			},
			requestID1:        big.NewInt(11),
			requestID2:        big.NewInt(11),
			expectedDuration1: func(tStep uint64) uint64 { return 0 }, // (P1-1) * t_step
			expectedDuration2: func(tStep uint64) uint64 { return 0 }, // result already published by member 1 -1
		},
		"two members publish the same results for different Request IDs": {
			resultToPublish1: &relayChain.DKGResult{
				GroupPublicKey: []byte{101},
			},
			resultToPublish2: &relayChain.DKGResult{
				GroupPublicKey: []byte{101},
			},
			requestID1:        big.NewInt(12),
			requestID2:        big.NewInt(13),
			expectedDuration1: func(tStep uint64) uint64 { return 0 }, // (P1-1) * t_step
			expectedDuration2: func(tStep uint64) uint64 {
				return (uint64(member2.index) - 1) * tStep // (P4-1) * t_step
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			chainHandle, initialBlock, err :=
				initChainHandle(threshold, groupSize)
			if err != nil {
				t.Fatal(err)
			}

			config, err := chainHandle.ThresholdRelay().GetConfig()
			if err != nil {
				t.Fatal(err)
			}

			tStep := config.ResultPublicationBlockStep

			expectedBlockEnd1 := initialBlock + test.expectedDuration1(tStep)
			expectedBlockEnd2 := initialBlock + test.expectedDuration2(tStep)

			result1Chan := make(chan uint64)
			defer close(result1Chan)
			result2Chan := make(chan uint64)
			defer close(result2Chan)

			go func() {
				blockCounter, _ := chainHandle.BlockCounter()

				err := member1.SubmitDKGResult(
					test.requestID1,
					test.resultToPublish1,
					signatures,
					chainHandle.ThresholdRelay(),
					blockCounter,
				)
				if err != nil {
					t.Fatal(err)
				}

				currentBlock, _ := blockCounter.CurrentBlock()
				result1Chan <- currentBlock
			}()

			go func() {
				blockCounter, _ := chainHandle.BlockCounter()

				err := member2.SubmitDKGResult(
					test.requestID2,
					test.resultToPublish2,
					signatures,
					chainHandle.ThresholdRelay(),
					blockCounter,
				)
				if err != nil {
					t.Fatal(err)
				}

				currentBlock, _ := blockCounter.CurrentBlock()
				result2Chan <- currentBlock
			}()

			if result1 := <-result1Chan; result1 != expectedBlockEnd1 {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedBlockEnd1, result1)
			}
			if result2 := <-result2Chan; result2 != expectedBlockEnd2 {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedBlockEnd2, result2)
			}
		})
	}
}

func initChainHandle(threshold, groupSize int) (chain.Handle, uint64, error) {
	chainHandle := local.Connect(groupSize, threshold, big.NewInt(200))

	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, 0, err
	}
	initialBlockChan, err := blockCounter.BlockHeightWaiter(1)
	if err != nil {
		return nil, 0, err
	}

	return chainHandle, <-initialBlockChan, nil
}
