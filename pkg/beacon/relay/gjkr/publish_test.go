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
		Hash:        []byte(fmt.Sprintf("%v", result1)),
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
		Hash:        []byte(fmt.Sprintf("%v", result2)),
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
		Hash:        []byte(fmt.Sprintf("%v", result1)),
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
		Hash:        []byte(fmt.Sprintf("%v", result2)),
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
