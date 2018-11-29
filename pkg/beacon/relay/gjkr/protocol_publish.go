package gjkr

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
)

// Result returns a result of distributed key generation. It takes generated
// group public key along with disqualified and inactive members and returns
// it in a Result struct.
//
// Additional validation to check if number of disqualified and inactive members
// is greater than half of the configured dishonest threshold. If so the group
// is to weak and the result is set to a failure.
func (pm *PublishingMember) Result() *result.Result {
	group := pm.group
	disqualifiedMembers := group.DisqualifiedMembers() // DQ
	inactiveMembers := group.InactiveMembers()         // IA

	// if nPlayers(IA + DQ) > T/2:
	if len(disqualifiedMembers)+len(inactiveMembers) > (group.dishonestThreshold / 2) {
		// Result.failure(disqualified = DQ)
		return &result.Result{
			Success:      false,
			Disqualified: disqualifiedMembers,
		}
	}

	// Result.success(pubkey = Y, inactive = IA, disqualified = DQ)
	return &result.Result{
		Success:        true,
		GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
		Disqualified:   disqualifiedMembers,
		Inactive:       inactiveMembers,
	}
}

// PublishResult sends a result containing i.a. group public key to the blockchain.
// It checks if the result has already been published to the blockchain. If not
// it determines if the current member is eligable to result submission. If allowed
// it submits the results to the blockchain. The function returns result published
// to the blockchain containing ID of the member who published it.
//
// See Phase 13 of the protocol specification.
func (pm *PublishingMember) PublishResult() (*event.PublishedResult, error) {
	chainRelay := pm.protocolConfig.ChainHandle().ThresholdRelay()

	onPublishedResultChan := make(chan *event.PublishedResult)
	chainRelay.OnResultPublished(func(publishedResult *event.PublishedResult) {
		onPublishedResultChan <- publishedResult
	})

	resultToPublish := pm.Result()

	blockCounter, err := pm.protocolConfig.ChainHandle().BlockCounter()
	if err != nil {
		return nil, fmt.Errorf("block counter failure [%v]", err)
	}

	// Waits until the current member is eligable to submit a result to the
	// blockchain.
	eligibleToSubmitWaiter, err := blockCounter.BlockWaiter(
		pm.PublishingIndex() * pm.protocolConfig.chain.blockStep)
	if err != nil {
		return nil, fmt.Errorf("block waiter failure [%v]", err)
	}

	// Check if the result is already published on the chain.
	if publishedResult := chainRelay.IsResultPublished(resultToPublish); publishedResult != nil {
		return publishedResult, nil
	}

	for {
		select {
		case <-eligibleToSubmitWaiter:
			publishedResultChan := make(chan *event.PublishedResult)
			errors := make(chan error)

			chainRelay.SubmitResult(pm.ID, resultToPublish).
				OnComplete(func(publishedResult *event.PublishedResult, err error) {
					publishedResultChan <- publishedResult
					errors <- err
				})
			return <-publishedResultChan, <-errors

		case newResult := <-onPublishedResultChan:
			// Check if published result matches a result the current member
			// wants to publish.
			if reflect.DeepEqual(resultToPublish, newResult.Result) {
				return newResult, nil
			}
		}
	}
}
