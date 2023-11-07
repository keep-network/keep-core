package result

import (
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"

	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

func TestSubmitDKGResult(t *testing.T) {
	honestThreshold := 3
	groupSize := 5

	beaconChain, _, initialBlock, err := initChainHandle(honestThreshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	config := beaconChain.GetConfig()

	result := &beaconchain.DKGResult{
		GroupPublicKey: []byte{123, 45},
	}
	signatures := map[group.MemberIndex][]byte{
		1: []byte{101},
		2: []byte{102},
		3: []byte{103},
		4: []byte{104},
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
			member := &SubmittingMember{
				logger: &testutils.MockLogger{},
				index:  group.MemberIndex(test.memberIndex),
			}

			// Reinitialize chain to reset block counter
			beaconChain, blockCounter, initialBlockHeight, err := initChainHandle(
				honestThreshold,
				groupSize,
			)
			if err != nil {
				t.Fatalf("chain initialization failed [%v]", err)
			}

			isSubmitted, err := beaconChain.IsGroupRegistered(result.GroupPublicKey)
			if err != nil {
				t.Fatal(err)
			}

			if isSubmitted {
				t.Fatalf("result is already submitted to the chain")
			}

			err = member.SubmitDKGResult(
				result,
				signatures,
				beaconChain,
				blockCounter,
				initialBlockHeight,
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
			isSubmitted, err = beaconChain.IsGroupRegistered(result.GroupPublicKey)
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
	honestThreshold := 3
	groupSize := 5

	member1 := &SubmittingMember{
		logger: &testutils.MockLogger{},
		index:  group.MemberIndex(1), // P1
	}
	member2 := &SubmittingMember{
		logger: &testutils.MockLogger{},
		index:  group.MemberIndex(4), // P4
	}

	signatures := map[group.MemberIndex][]byte{
		1: []byte{101},
		2: []byte{102},
		3: []byte{103},
		4: []byte{104},
	}

	var tests = map[string]struct {
		resultToPublish1  *beaconchain.DKGResult
		resultToPublish2  *beaconchain.DKGResult
		expectedDuration1 func(tStep uint64) uint64 // index * t_step
		expectedDuration2 func(tStep uint64) uint64 // index * t_step
	}{
		"two members publish the same results": {
			resultToPublish1: &beaconchain.DKGResult{
				GroupPublicKey: []byte{101},
			},
			resultToPublish2: &beaconchain.DKGResult{
				GroupPublicKey: []byte{101},
			},
			expectedDuration1: func(tStep uint64) uint64 { return 0 }, // (P1-1) * t_step
			expectedDuration2: func(tStep uint64) uint64 { return 0 }, // result already published by member 1 -1
		},
		"two members publish different results": {
			resultToPublish1: &beaconchain.DKGResult{
				GroupPublicKey: []byte{201},
			},
			resultToPublish2: &beaconchain.DKGResult{
				GroupPublicKey: []byte{202},
			},
			expectedDuration1: func(tStep uint64) uint64 { return 0 }, // (P1-1) * t_step
			expectedDuration2: func(tStep uint64) uint64 { return 0 }, // result already published by member 1 -1
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			beaconChain, blockCounter, initialBlock, err :=
				initChainHandle(honestThreshold, groupSize)
			if err != nil {
				t.Fatal(err)
			}

			config := beaconChain.GetConfig()

			tStep := config.ResultPublicationBlockStep

			expectedBlockEnd1 := initialBlock + test.expectedDuration1(tStep)
			expectedBlockEnd2 := initialBlock + test.expectedDuration2(tStep)

			result1Chan := make(chan uint64)
			defer close(result1Chan)
			result2Chan := make(chan uint64)
			defer close(result2Chan)

			go func() {
				err := member1.SubmitDKGResult(
					test.resultToPublish1,
					signatures,
					beaconChain,
					blockCounter,
					initialBlock,
				)
				if err != nil {
					t.Error(err)
				}

				currentBlock, _ := blockCounter.CurrentBlock()
				result1Chan <- currentBlock
			}()

			go func() {
				err := member2.SubmitDKGResult(
					test.resultToPublish2,
					signatures,
					beaconChain,
					blockCounter,
					initialBlock,
				)
				if err != nil {
					t.Error(err)
				}

				currentBlock, _ := blockCounter.CurrentBlock()
				result2Chan <- currentBlock
			}()

			if result1 := <-result1Chan; result1 != expectedBlockEnd1 {
				t.Errorf("\nexpected: %v\nactual:   %v\n", expectedBlockEnd1, result1)
			}
			if result2 := <-result2Chan; result2 != expectedBlockEnd2 {
				t.Errorf("\nexpected: %v\nactual:   %v\n", expectedBlockEnd2, result2)
			}
		})
	}
}

func initChainHandle(honestThreshold int, groupSize int) (
	beaconchain.Interface,
	chain.BlockCounter,
	uint64,
	error,
) {
	chainHandle := local_v1.Connect(groupSize, honestThreshold)

	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, nil, 0, err
	}
	initialBlockChan, err := blockCounter.BlockHeightWaiter(1)
	if err != nil {
		return nil, nil, 0, err
	}

	return chainHandle, blockCounter, <-initialBlockChan, nil
}
