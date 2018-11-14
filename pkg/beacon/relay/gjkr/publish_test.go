package gjkr

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestPublishResult(t *testing.T) {
	members, err := initializePublishingMembersGroup(1, 3)
	if err != nil {
		t.Fatalf("%s", err)
	}
	member := members[0]

	chainRelay := member.protocolConfig.chain.ThresholdRelay()

	result1 := &result.Result{GroupPublicKey: big.NewInt(1)}
	expectedEvent := &event.PublishedResult{
		PublisherID: member.ID,
		// Hash:        []byte(fmt.Sprintf("%v", result1)),
		Hash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148, 51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
	}
	if chainRelay.IsResultPublished(result1) {
		t.Fatalf("Result is already published on chain")
	}

	eventPublish, err := member.PublishResult(result1, 5)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedEvent, eventPublish) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedEvent, eventPublish)
	}

	if !chainRelay.IsResultPublished(result1) {
		t.Fatalf("Result should be published on chain")
	}

	result2 := &result.Result{GroupPublicKey: big.NewInt(2)}
	if chainRelay.IsResultPublished(result2) {
		t.Fatalf("Result is already published on chain")
	}
	eventPublish2, err := member.PublishResult(result2, 5)
	expectedEvent2 := &event.PublishedResult{
		PublisherID: member.ID,
		// Hash:        []byte(fmt.Sprintf("%v", result2)),
		Hash: []byte{44, 154, 182, 135, 19, 160, 13, 133, 181, 132, 146, 42, 153, 216, 10, 81, 67, 136, 165, 213, 19, 43, 219, 181, 4, 110, 168, 199, 224, 88, 149, 80},
	}
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedEvent2, eventPublish2) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedEvent2, eventPublish2)
	}

	if !chainRelay.IsResultPublished(result2) {
		t.Fatalf("Result should be published on chain")
	}
}

func TestPublishResult2(t *testing.T) {
	members, err := initializePublishingMembersGroup(1, 3)
	if err != nil {
		t.Fatalf("%s", err)
	}

	member := members[0]

	result1 := &result.Result{GroupPublicKey: big.NewInt(20001)}
	expectedEvent1 := &event.PublishedResult{
		PublisherID: members[0].ID,
		// Hash:        []byte(fmt.Sprintf("%v", result1)),
		Hash: []byte{17, 137, 135, 209, 119, 129, 62, 207, 107, 14, 232, 183, 212, 85, 145, 250, 177, 214, 29, 131, 210, 38, 166, 15, 30, 249, 96, 53, 131, 87, 139, 200},
	}
	eventPublish1, err := member.PublishResult(result1, 5)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedEvent1, eventPublish1) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedEvent1, eventPublish1)
	}

	member = members[1]
	eventPublish21, err := member.PublishResult(result1, 5)
	// expectedError := fmt.Errorf("sad")
	// if !reflect.DeepEqual(expectedError, err) {
	// 	t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	// }
	if eventPublish21 != nil {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", nil, eventPublish21)
	}

	result2 := &result.Result{GroupPublicKey: big.NewInt(20002)}
	expectedEvent2 := &event.PublishedResult{
		PublisherID: members[1].ID,
		// Hash:        []byte(fmt.Sprintf("%v", result2)),
		Hash: []byte{152, 74, 143, 75, 44, 214, 60, 242, 96, 1, 6, 138, 243, 191, 180, 171, 172, 83, 149, 128, 39, 129, 211, 169, 247, 41, 67, 149, 219, 31, 128, 105},
	}
	eventPublish2, err := member.PublishResult(result2, 5)
	if err != nil {
		t.Fatalf("\nexpected: %s\nactual:   %s\n", "", err)
	}
	if !reflect.DeepEqual(expectedEvent2, eventPublish2) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n", expectedEvent2, eventPublish2)
	}

}

func initializePublishingMembersGroup(threshold, groupSize int) ([]*PublishingMember, error) {
	chain := local.Connect(10, 4)
	blockCounter, err := chain.BlockCounter()
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

	dkg := &DKG{
		chain:              chain,
		expectedDuration:   4,
		blockStep:          1,
		initialBlockHeight: initialBlockHeight,
	}

	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	var members []*PublishingMember

	for i := 1; i <= groupSize; i++ {
		id := i
		members = append(members,
			&PublishingMember{
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
			})
		group.RegisterMemberID(id)
	}
	return members, nil
}
