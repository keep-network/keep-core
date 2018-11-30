package publish

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestPublishResult(t *testing.T) {
	threshold := 2
	groupSize := 5
	blockStep := 2 // T_step

	chainHandle, initialBlock, err := initChainHandle(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	resultToPublish := &result.Result{
		Success:        false,
		GroupPublicKey: big.NewInt(12345),
		Disqualified:   []int{1, 2},
		Inactive:       []int{5},
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
				publishingIndex: test.publishingIndex,
				chainHandle:     chainHandle,
				blockStep:       blockStep,
			}

			expectedPublishedResult := &event.PublishedResult{
				PublisherID: publisher.ID,
				Result:      resultToPublish,
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

			if chainRelay.IsResultPublished(resultToPublish) != nil {
				t.Fatalf("result is already published on chain")
			}
			// TEST
			publishedResult, err := publisher.PublishResult(resultToPublish)
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
			if !reflect.DeepEqual(expectedPublishedResult, publishedResult) {
				t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult)
			}
			if chainRelay.IsResultPublished(resultToPublish) == nil {
				t.Fatalf("result is not published on chain")
			}
		})
	}
}

func TestPublishResult_AlreadyPublished(t *testing.T) {
	threshold := 2
	groupSize := 5
	blockStep := 2 // T_step

	chainHandle, _, err := initChainHandle(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	publisher1 := &Publisher{
		ID:              1,
		publishingIndex: 0,
		chainHandle:     chainHandle,
		blockStep:       blockStep,
	}
	publisher2 := &Publisher{
		ID:              2,
		publishingIndex: 1,
		chainHandle:     chainHandle,
		blockStep:       blockStep,
	}

	resultToPublish := &result.Result{
		GroupPublicKey: big.NewInt(12345),
	}
	expectedPublishedResult := &event.PublishedResult{
		PublisherID: publisher1.ID,
		Result:      resultToPublish,
	}

	chainRelay := chainHandle.ThresholdRelay()

	if chainRelay.IsResultPublished(resultToPublish) != nil {
		t.Fatalf("result is already published on chain")
	}

	// Case: Member 1 publishes a result.
	// Expected: A new result is published successfully by member 1.
	publishedResult1, err := publisher1.PublishResult(resultToPublish)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedPublishedResult, publishedResult1) {
		t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult1)
	}
	if chainRelay.IsResultPublished(resultToPublish) == nil {
		t.Fatalf("result is not published on chain")
	}

	// Case: Member 1 publishes the same result once again.
	// Expected: A new result is not published, function returns result published
	// already in previous step.
	publishedResult2, err := publisher1.PublishResult(resultToPublish)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedPublishedResult, publishedResult2) {
		t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult2)
	}
	if chainRelay.IsResultPublished(resultToPublish) == nil {
		t.Fatalf("result is not published on chain")
	}

	// Case: Member 2 publishes the same result as member 1 already did.
	// Expected: A new result is not published, function returns result published
	// already by member 1.
	var expectedError error
	expectedError = nil

	publishedResult3, err := publisher2.PublishResult(resultToPublish)
	if !reflect.DeepEqual(expectedPublishedResult, publishedResult3) {
		t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult3)
	}
	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", expectedError, err)
	}
	if chainRelay.IsResultPublished(resultToPublish) == nil {
		t.Fatalf("result is not published on chain")
	}
}

// This tests runs result publication concurrently by two members.
// Member with lower index gets to publish the result to chain. For the second
// member loop should be aborted and result published by the first member should
// be returned.
func TestPublishResult_ConcurrentExecution(t *testing.T) {
	threshold := 2
	groupSize := 5
	blockStep := 2 // T_step

	publisher1 := &Publisher{
		ID:              2,
		publishingIndex: 1,
		blockStep:       blockStep,
	}
	publisher2 := &Publisher{
		ID:              5,
		publishingIndex: 4,
		blockStep:       blockStep,
	}

	var tests = map[string]struct {
		resultToPublish1         *result.Result
		resultToPublish2         *result.Result
		expectedPublishedResult1 *event.PublishedResult
		expectedPublishedResult2 *event.PublishedResult
		expectedDuration         int
	}{
		"two members publish the same results": {
			resultToPublish1: &result.Result{
				GroupPublicKey: big.NewInt(101),
			},
			resultToPublish2: &result.Result{
				GroupPublicKey: big.NewInt(101),
			},
			expectedPublishedResult1: &event.PublishedResult{
				PublisherID: publisher1.ID,
				Result: &result.Result{
					GroupPublicKey: big.NewInt(101),
				},
			},
			expectedPublishedResult2: &event.PublishedResult{
				PublisherID: publisher1.ID,
				Result: &result.Result{
					GroupPublicKey: big.NewInt(101),
				},
			},
			expectedDuration: publisher1.publishingIndex * blockStep, // index * t_step
		},
		"two members publish different results": {
			resultToPublish1: &result.Result{
				GroupPublicKey: big.NewInt(201),
			},
			resultToPublish2: &result.Result{
				GroupPublicKey: big.NewInt(202),
			},
			expectedPublishedResult1: &event.PublishedResult{
				PublisherID: publisher1.ID,
				Result: &result.Result{
					GroupPublicKey: big.NewInt(201),
				},
			},
			expectedPublishedResult2: &event.PublishedResult{
				PublisherID: publisher2.ID,
				Result: &result.Result{
					GroupPublicKey: big.NewInt(202),
				},
			},
			expectedDuration: publisher2.publishingIndex * blockStep, // index * t_step
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			chainHandle, initialBlock, err := initChainHandle(threshold, groupSize)
			if err != nil {
				t.Fatal(err)
			}
			publisher1.chainHandle = chainHandle
			publisher2.chainHandle = chainHandle

			expectedBlockEnd := initialBlock + test.expectedDuration

			result1Chan := make(chan *event.PublishedResult)
			result2Chan := make(chan *event.PublishedResult)

			go func() {
				publishedResult1, err := publisher1.PublishResult(test.resultToPublish1)
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				result1Chan <- publishedResult1
			}()

			go func() {
				publishedResult2, err := publisher2.PublishResult(test.resultToPublish2)
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				result2Chan <- publishedResult2
			}()

			if result1 := <-result1Chan; !reflect.DeepEqual(result1, test.expectedPublishedResult1) {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedPublishedResult1, result1)
			}
			if result2 := <-result2Chan; !reflect.DeepEqual(result2, test.expectedPublishedResult2) {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedPublishedResult2, result2)
			}

			blockCounter, err := chainHandle.BlockCounter()
			if err != nil {
				t.Fatalf("unexpected error [%v]", err)
			}
			currentBlock, err := blockCounter.CurrentBlock()
			if err != nil {
				t.Fatalf("unexpected error [%v]", err)
			}

			if expectedBlockEnd != currentBlock {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedBlockEnd, currentBlock)
			}
		})
	}
}

func initChainHandle(threshold, groupSize int) (chainHandle chain.Handle, initialBlock int, err error) {
	chainHandle = local.Connect(groupSize, threshold)
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, -1, err
	}
	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, -1, err
	}

	initialBlock, err = blockCounter.CurrentBlock()
	if err != nil {
		return nil, -1, err
	}
	return
}
