package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestPrepareResult(t *testing.T) {
	threshold := 4
	groupSize := 8
	expectedProtocolDuration := 3 // T_dkg
	blockStep := 2                // T_step

	members, err := initializePublishingMembersGroup(threshold, groupSize, expectedProtocolDuration, blockStep)
	if err != nil {
		t.Fatalf("%s", err)
	}

	publishingMember := members[0]

	var tests = map[string]struct {
		disqualifiedMemberIDs []int
		inactiveMemberIDs     []int
		expectedResult        *result.Result
	}{
		"no disqualified or inactive members - success": {
			expectedResult: &result.Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   nil,
				Inactive:       nil,
			},
		},
		"one disqualified member - success": {
			disqualifiedMemberIDs: []int{2},
			expectedResult: &result.Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   []int{2},
				Inactive:       nil,
			},
		},
		"two inactive members - success": {
			inactiveMemberIDs: []int{3, 7},
			expectedResult: &result.Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   nil,
				Inactive:       []int{3, 7},
			},
		},
		"more than half of threshold disqualified and inactive members - failure": {
			disqualifiedMemberIDs: []int{2},
			inactiveMemberIDs:     []int{3, 7},
			expectedResult: &result.Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   []int{2},
				Inactive:       nil, // in case of failure only disqualified members are slashed
			},
		},
		"more than half of threshold inactive members - failure": {
			inactiveMemberIDs: []int{3, 5, 7},
			expectedResult: &result.Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   nil,
				Inactive:       nil, // in case of failure only disqualified members are slashed
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			publishingMember.group.disqualifiedMemberIDs = test.disqualifiedMemberIDs
			publishingMember.group.inactiveMemberIDs = test.inactiveMemberIDs

			publishingMember.PrepareResult()

			if !reflect.DeepEqual(test.expectedResult, publishingMember.result) {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, publishingMember.result)
			}
		})
	}
}

func TestPublishResult(t *testing.T) {
	threshold := 2
	groupSize := 5
	expectedProtocolDuration := 3 // T_dkg
	blockStep := 2                // T_step

	members, err := initializePublishingMembersGroup(threshold, groupSize, expectedProtocolDuration, blockStep)
	if err != nil {
		t.Fatalf("%s", err)
	}
	initialBlock := members[0].protocolConfig.chain.initialBlockHeight // T_init

	var tests = map[string]struct {
		publisher       *PublishingMember
		expectedTimeEnd int
	}{
		"first member eligible to publish straight away": {
			publisher:       members[0],
			expectedTimeEnd: initialBlock, // T_now < T_init + T_step
		},
		"second member eligible to publish after T_step block passed": {
			publisher:       members[1],
			expectedTimeEnd: initialBlock + blockStep, // T_now = T_init + T_step
		},
		"fourth member eligable to publish after T_dkg + 2*T_step passed": {
			publisher:       members[3],
			expectedTimeEnd: initialBlock + 3*blockStep, // T_now = T_init + 3*T_step
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := test.publisher.Result()

			expectedPublishedResult := &event.PublishedResult{
				PublisherID: test.publisher.ID,
				Result:      result,
			}

			// Reinitialize chain to reset block counter
			test.publisher.protocolConfig.chain, err = initChain(threshold, groupSize, expectedProtocolDuration, blockStep)
			if err != nil {
				t.Fatalf("chain initialization failed [%v]", err)
			}

			chainRelay := test.publisher.protocolConfig.ChainHandle().ThresholdRelay()

			if chainRelay.IsResultPublished(result) != nil {
				t.Fatalf("result is already published on chain")
			}
			// TEST
			publishedResult, err := test.publisher.PublishResult()
			if err != nil {
				t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
			}
			currentBlock, err := test.publisher.protocolConfig.chain.CurrentBlock()
			if err != nil {
				t.Fatalf("unexpected error [%v]", err)
			}
			if test.expectedTimeEnd != currentBlock {
				t.Fatalf("invalid current block\nexpected: %v\nactual:   %v\n", test.expectedTimeEnd, currentBlock)
			}
			if !reflect.DeepEqual(expectedPublishedResult, publishedResult) {
				t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult)
			}
			if chainRelay.IsResultPublished(result) == nil {
				t.Fatalf("result is not published on chain")
			}
		})
	}
}

func TestPublishResult_AlreadyPublished(t *testing.T) {
	threshold := 2
	groupSize := 5
	expectedProtocolDuration := 3 // T_dkg
	blockStep := 2                // T_step

	members, err := initializePublishingMembersGroup(threshold, groupSize, expectedProtocolDuration, blockStep)
	if err != nil {
		t.Fatalf("%s", err)
	}
	publisher := members[0]

	result := &result.Result{GroupPublicKey: big.NewInt(13)}

	chainRelay := publisher.protocolConfig.ChainHandle().ThresholdRelay()

	if chainRelay.IsResultPublished(result) {
		t.Fatalf("result is already published on chain")
	}

	// Publish a result
	expectedPublishedResult := &event.PublishedResult{
		PublisherID: publisher.ID,
		Result:      []byte(fmt.Sprintf("%v", result)),
	}

	publishedResult, err := publisher.PublishResult(result)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedPublishedResult, publishedResult) {
		t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult)
	}
	if !chainRelay.IsResultPublished(result) {
		t.Fatalf("result is not published on chain")
	}

	// Publish the same result for the second time
	expectedPublishedResult = nil

	publishedResult, err = publisher.PublishResult(result)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedPublishedResult, publishedResult) {
		t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult)
	}
	if !chainRelay.IsResultPublished(result) {
		t.Fatalf("result is not published on chain")
	}
}

func TestDeterminePublishersIDs(t *testing.T) {
	threshold := 1
	groupSize := 3
	expectedProtocolDuration := 3 // T_dkg
	blockStep := 2                // T_step

	members, err := initializePublishingMembersGroup(threshold, groupSize, expectedProtocolDuration, blockStep)
	if err != nil {
		t.Fatalf("%s", err)
	}
	member := members[0]
	initialBlock := member.protocolConfig.chain.initialBlockHeight // T_init

	blockCounter, err := member.protocolConfig.ChainHandle().BlockCounter()
	if err != nil {
		t.Fatalf("getting block counter failed [%v]", err)
	}

	// Tests steps are strictly corelated so we define them in a slice to execute
	// them in exactly the same order as they are defined.
	var testSteps = []struct {
		name                  string
		waitForBlocks         int   // wait for number of blocks to pass before test step execution
		expectedPublishersIDs []int // expected evaluated publishers IDs
		expectedCurrentBlock  int   // expected current block number
	}{
		{
			name:                  "T_elapsed < T_dkg",
			waitForBlocks:         0,
			expectedPublishersIDs: []int{1},
			expectedCurrentBlock:  initialBlock, // T_init = 1
		},
		{
			name:                  "T_elapsed = T_dkg",
			waitForBlocks:         expectedProtocolDuration,
			expectedPublishersIDs: []int{1},
			expectedCurrentBlock:  initialBlock + expectedProtocolDuration, // T_init + T_dkg = 4
		},
		{
			name:                  "T_elapsed > T_dkg && T_over < T_step",
			waitForBlocks:         1,
			expectedPublishersIDs: []int{1, 2},
			expectedCurrentBlock:  initialBlock + expectedProtocolDuration + 1, // T_init + T_dkg + 1 = 5
		},
		{
			name:                  "T_elapsed > T_dkg && T_over > T_step",
			waitForBlocks:         blockStep,
			expectedPublishersIDs: []int{1, 2, 3},
			expectedCurrentBlock:  initialBlock + expectedProtocolDuration + 1 + blockStep, // T_init + T_dkg + 1 + T_step = 7
		},
		{
			name:                  "T_elapsed > T_dkg && T_over > 2*T_step",
			waitForBlocks:         blockStep,
			expectedPublishersIDs: []int{1, 2, 3},
			expectedCurrentBlock:  initialBlock + expectedProtocolDuration + 1 + 2*blockStep, // T_init + T_dkg + 1 + 2*T_step = 9
		},
	}

	for _, testStep := range testSteps {
		blockCounter.WaitForBlocks(testStep.waitForBlocks)

		// Execute function under test
		result, err := member.determinePublishersIDs()
		if err != nil {
			t.Fatalf("%s", err)
		}

		currentBlock, err := blockCounter.CurrentBlock()
		if err != nil {
			t.Fatalf("getting current block failed [%v]", err)
		}
		if currentBlock != testStep.expectedCurrentBlock {
			t.Fatalf("invalid current block for step: %s\nexpected: %v\nactual:   %v\n",
				testStep.name,
				testStep.expectedCurrentBlock,
				currentBlock,
			)

		}

		if !reflect.DeepEqual(testStep.expectedPublishersIDs, result) {
			t.Fatalf("invalid publishers IDs for step: %s\nexpected: %v\nactual:   %v\n",
				testStep.name,
				testStep.expectedPublishersIDs,
				result,
			)
		}
	}
}

func initializePublishingMembersGroup(
	threshold, groupSize, expectedProtocolDuration, blockStep int,
) ([]*PublishingMember, error) {
	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	chain, err := initChain(threshold, groupSize, expectedProtocolDuration, blockStep)
	if err != nil {
		return nil, err
	}

	dkg := &DKG{chain: chain}

	var members []*PublishingMember

	for i := 1; i <= groupSize; i++ {
		id := i
		members = append(members,
			&PublishingMember{
				PointsJustifyingMember: &PointsJustifyingMember{
					SharingMember: &SharingMember{
						QualifiedMember: &QualifiedMember{
							SharesJustifyingMember: &SharesJustifyingMember{
								CommittingMember: &CommittingMember{
									memberCore: &memberCore{
										ID:             id,
										group:          group,
										protocolConfig: dkg,
									},
								},
							},
						},
					},
				},
			})
		group.RegisterMemberID(id)
	}
	return members, nil
}

func initChain(
	threshold, groupSize, expectedProtocolDuration, blockStep int,
) (*Chain, error) {
	chainHandle := local.Connect(groupSize, threshold)
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, err
	}
	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, err
	}

	initialBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, err
	}

	return &Chain{
		handle:                   chainHandle,
		expectedProtocolDuration: expectedProtocolDuration, // T_dkg
		blockStep:                blockStep,                // T_step
		initialBlockHeight:       initialBlockHeight,       // T_init
	}, nil
}
