package gjkr

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestResult(t *testing.T) {
	threshold := 4
	groupSize := 8
	blockStep := 2 // T_step

	members, err := initializePublishingMembersGroup(threshold, groupSize, blockStep)
	if err != nil {
		t.Fatalf("%s", err)
	}

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
			for _, member := range members {
				member.group.disqualifiedMemberIDs = test.disqualifiedMemberIDs
				member.group.inactiveMemberIDs = test.inactiveMemberIDs

				resultToPublish := member.Result()

				if !reflect.DeepEqual(test.expectedResult, resultToPublish) {
					t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, resultToPublish)
				}
			}
		})
	}
}

func TestPublishResult(t *testing.T) {
	threshold := 2
	groupSize := 5
	blockStep := 2 // T_step

	members, err := initializePublishingMembersGroup(threshold, groupSize, blockStep)
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
			test.publisher.protocolConfig.chain, err = initChain(threshold, groupSize, blockStep)
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
	blockStep := 2 // T_step

	members, err := initializePublishingMembersGroup(threshold, groupSize, blockStep)
	if err != nil {
		t.Fatalf("%s", err)
	}

	publisher1 := members[0]
	publisher2 := members[1]

	expectedResult := publisher1.Result()

	chainRelay := publisher1.protocolConfig.ChainHandle().ThresholdRelay()

	if chainRelay.IsResultPublished(expectedResult) != nil {
		t.Fatalf("result is already published on chain")
	}

	expectedPublishedResult := &event.PublishedResult{
		PublisherID: publisher1.ID,
		Result:      expectedResult,
	}

	// Case: Member 1 publishes a result.
	// Expected: A new result is published successfully by member 1.
	publishedResult1, err := publisher1.PublishResult()
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedPublishedResult, publishedResult1) {
		t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult1)
	}
	if chainRelay.IsResultPublished(expectedResult) == nil {
		t.Fatalf("result is not published on chain")
	}

	// Case: Member 1 publishes the same result once again.
	// Expected: A new result is not published, function returns result published
	// already in previous step.
	publishedResult2, err := publisher1.PublishResult()
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedPublishedResult, publishedResult2) {
		t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult2)
	}
	if chainRelay.IsResultPublished(expectedResult) == nil {
		t.Fatalf("result is not published on chain")
	}

	// Case: Member 2 publishes the same result as member 1 already did.
	// Expected: A new result is not published, function returns result published
	// already by member 1.
	var expectedError error
	expectedError = nil

	publishedResult3, err := publisher2.PublishResult()
	if !reflect.DeepEqual(expectedPublishedResult, publishedResult3) {
		t.Fatalf("invalid published result\nexpected: %v\nactual:   %v\n", expectedPublishedResult, publishedResult3)
	}
	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", expectedError, err)
	}
	if chainRelay.IsResultPublished(expectedResult) == nil {
		t.Fatalf("result is not published on chain")
	}
}

func initializePublishingMembersGroup(threshold, groupSize, blockStep int,
) ([]*PublishingMember, error) {
	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	chain, err := initChain(threshold, groupSize, blockStep)
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

func initChain(threshold, groupSize, blockStep int) (*Chain, error) {
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
		handle:             chainHandle,
		blockStep:          blockStep,          // T_step
		initialBlockHeight: initialBlockHeight, // T_init
	}, nil
}
