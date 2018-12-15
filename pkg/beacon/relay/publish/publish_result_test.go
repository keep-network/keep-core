package publish

import (
	"math/big"
	"reflect"
	"testing"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestPublishDKGResult(t *testing.T) {
	return // xyzzy
	threshold := 2
	groupSize := 5
	blockStep := 2 // T_step

	chainHandle, initialBlock, err := initChainHandle(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	resultToPublish := &relayChain.DKGResult{
		GroupPublicKey: big.NewInt(12345),
	}

	var tests = map[string]struct {
		publishingIndex int
		expectedTimeEnd int
	}{
		"first member eligible to publish straight away": {
			publishingIndex: 0,
			expectedTimeEnd: initialBlock, // T_now < T_init + T_step
		},
		"second member eligible to publish after T_step block passed": {
			publishingIndex: 1,
			expectedTimeEnd: initialBlock + blockStep, // T_now = T_init + T_step
		},
		"fourth member eligable to publish after T_dkg + 2*T_step passed": {
			publishingIndex: 3,
			expectedTimeEnd: initialBlock + 3*blockStep, // T_now = T_init + 3*T_step
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			publisher := &Publisher{
				ID:              gjkr.MemberID(test.publishingIndex + 1),
				RequestID:       big.NewInt(101),
				publishingIndex: test.publishingIndex,
				chainHandle:     chainHandle,
				blockStep:       blockStep,
			}

			// Reinitialize chain to reset block counter
			publisher.chainHandle, initialBlock, err = initChainHandle(threshold, groupSize)
			if err != nil {
				t.Fatalf("chain initialization failed [%v]", err)
			}

			chainRelay := publisher.chainHandle.ThresholdRelay()
			blockCounter, err := publisher.chainHandle.BlockCounter()
			if err != nil {
				t.Fatalf("unexpected error [%v]", err)
			}

			if chainRelay.IsDKGResultPublished(
				publisher.RequestID,
				resultToPublish,
			) {
				t.Fatalf("result is already published on chain")
			}
			// TEST
			err = publisher.PublishDKGResult(resultToPublish)
			if err != nil {
				t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
			}
			currentBlock, err := blockCounter.CurrentBlock()
			if err != nil {
				t.Fatalf("unexpected error [%v]", err)
			}
			if test.expectedTimeEnd != currentBlock {
				t.Fatalf("invalid current block\nexpected: %v\nactual:   %v\n", test.expectedTimeEnd, currentBlock)
			}
			if !chainRelay.IsDKGResultPublished(
				publisher.RequestID,
				resultToPublish,
			) {
				t.Fatalf("result is not published on chain")
			}
		})
	}
}

func TestPublishDKGResult_AlreadyPublished(t *testing.T) {
	return // xyzzy
	threshold := 2
	groupSize := 5
	blockStep := 2 // T_step

	chainHandle, _, err := initChainHandle(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	publisher1 := &Publisher{
		ID:              1,
		RequestID:       big.NewInt(101),
		publishingIndex: 0,
		chainHandle:     chainHandle,
		blockStep:       blockStep,
	}
	publisher2 := &Publisher{
		ID:              2,
		RequestID:       big.NewInt(101),
		publishingIndex: 1,
		chainHandle:     chainHandle,
		blockStep:       blockStep,
	}

	resultToPublish := &relayChain.DKGResult{
		GroupPublicKey: big.NewInt(12345),
	}

	chainRelay := chainHandle.ThresholdRelay()

	if chainRelay.IsDKGResultPublished(publisher1.RequestID, resultToPublish) {
		t.Fatalf("result is already published on chain")
	}

	// Case: Member 1 publishes a result.
	// Expected: A new result is published successfully by member 1.
	err = publisher1.PublishDKGResult(resultToPublish)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !chainRelay.IsDKGResultPublished(publisher1.RequestID, resultToPublish) {
		t.Fatalf("result is already published on chain")
	}

	// Case: Member 1 publishes the same result once again.
	// Expected: A new result is not published, function returns result published
	// already in previous step.
	err = publisher1.PublishDKGResult(resultToPublish)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !chainRelay.IsDKGResultPublished(publisher1.RequestID, resultToPublish) {
		t.Fatalf("result is not published on chain")
	}

	// Case: Member 2 publishes the same result as member 1 already did.
	// Expected: A new result is not published, function returns result published
	// already by member 1.
	var expectedError error
	expectedError = nil

	if !chainRelay.IsDKGResultPublished(publisher2.RequestID, resultToPublish) {
		t.Fatalf("result is not published on chain")
	}

	err = publisher2.PublishDKGResult(resultToPublish)
	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}

	if !chainRelay.IsDKGResultPublished(publisher2.RequestID, resultToPublish) {
		t.Fatalf("result is not published on chain")
	}
}

// This tests runs result publication concurrently by two members.
// Member with lower index gets to publish the result to chain. For the second
// member loop should be aborted and result published by the first member should
// be returned.
func TestPublishDKGResult_ConcurrentExecution(t *testing.T) {
	return // xyzzy
	threshold := 2
	groupSize := 5
	blockStep := 2 // t_step

	publisher1 := &Publisher{
		ID:              2,
		publishingIndex: 1, // P1
		blockStep:       blockStep,
	}
	publisher2 := &Publisher{
		ID:              5,
		publishingIndex: 4, // P2
		blockStep:       blockStep,
	}

	var tests = map[string]struct {
		resultToPublish1  *relayChain.DKGResult
		resultToPublish2  *relayChain.DKGResult
		requestID1        *big.Int
		requestID2        *big.Int
		expectedDuration1 int // index * t_step
		expectedDuration2 int // index * t_step
	}{
		"two members publish the same results": {
			resultToPublish1: &relayChain.DKGResult{
				GroupPublicKey: big.NewInt(101),
			},
			resultToPublish2: &relayChain.DKGResult{
				GroupPublicKey: big.NewInt(101),
			},
			requestID1:        big.NewInt(11),
			requestID2:        big.NewInt(11),
			expectedDuration1: publisher1.publishingIndex * blockStep, // P1 * t_step
			expectedDuration2: publisher1.publishingIndex * blockStep, // P1 * t_step
		},
		"two members publish different results": {
			resultToPublish1: &relayChain.DKGResult{
				GroupPublicKey: big.NewInt(201),
			},
			resultToPublish2: &relayChain.DKGResult{
				GroupPublicKey: big.NewInt(202),
			},
			requestID1:        big.NewInt(11),
			requestID2:        big.NewInt(11),
			expectedDuration1: publisher1.publishingIndex * blockStep, // P1 * t_step
			expectedDuration2: publisher2.publishingIndex * blockStep, // P2 * t_step
		},
		"two members publish the same results for different Request IDs": {
			resultToPublish1: &relayChain.DKGResult{
				GroupPublicKey: big.NewInt(101),
			},
			resultToPublish2: &relayChain.DKGResult{
				GroupPublicKey: big.NewInt(101),
			},
			requestID1:        big.NewInt(12),
			requestID2:        big.NewInt(13),
			expectedDuration1: publisher1.publishingIndex * blockStep, // P1 * t_step
			expectedDuration2: publisher2.publishingIndex * blockStep, // P1 * t_step
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			publisher1.RequestID = test.requestID1
			publisher2.RequestID = test.requestID2

			chainHandle, initialBlock, err := initChainHandle(threshold, groupSize)
			if err != nil {
				t.Fatal(err)
			}
			publisher1.chainHandle = chainHandle
			publisher2.chainHandle = chainHandle

			expectedBlockEnd1 := initialBlock + test.expectedDuration1
			expectedBlockEnd2 := initialBlock + test.expectedDuration2

			result1Chan := make(chan int)
			result2Chan := make(chan int)

			blockCounter, err := chainHandle.BlockCounter()
			if err != nil {
				t.Fatalf("unexpected error [%v]", err)
			}

			go func() {
				err := publisher1.PublishDKGResult(test.resultToPublish1)
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				currentBlock, err := blockCounter.CurrentBlock()
				if err != nil {
					t.Fatalf("unexpected error [%v]", err)
				}

				result1Chan <- currentBlock
			}()

			go func() {
				err := publisher2.PublishDKGResult(test.resultToPublish2)
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				currentBlock, err := blockCounter.CurrentBlock()
				if err != nil {
					t.Fatalf("unexpected error [%v]", err)
				}

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

func initChainHandle(threshold, groupSize int) (chainHandle chain.Handle, initialBlock int, err error) {
	chainHandle = local.Connect(groupSize, threshold)
	blockCounter, err := chainHandle.BlockCounter() // PJS - save blockCounter?
	if err != nil {
		return nil, -1, err
	}
	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, -1, err
	}

	initialBlock, err = blockCounter.CurrentBlock() // PJS - need CurrentBlock to make this work
	if err != nil {
		return nil, -1, err
	}
	return
}
