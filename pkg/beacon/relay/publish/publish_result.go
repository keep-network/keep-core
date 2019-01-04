package publish

import (
	"fmt"
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
)

// Publisher submits distributed key generation result to a blockchain.
type Publisher struct {
	// ID of distributed key generation execution.
	RequestID *big.Int
	// Handle to interact with a blockchain.
	chainHandle chain.Handle
	// Sequential number of the current member in the publishing group.
	// The value is used to determine eligible publishing member. Indexing starts
	// with 1. Relates to DKG Phase 13.
	publishingIndex int
	// Predefined step for each publishing window. The value is used to determine
	// eligible publishing member. Relates to DKG Phase 13.
	blockStep int
}

// PublishDKGResult sends a result containing i.a. group public key to the blockchain.
// It checks if the result has already been published to the blockchain with
// request ID specific for current DKG execution. If not, it determines if the
// current member is eligable to result submission. If allowed, it submits the
// results to the blockchain.
// When member is waiting for their round the function keeps tracking results being
// published to the blockchain. If any result is published for the current request
// ID the phase is finished.
//
// See Phase 13 of the protocol specification.
func (pm *Publisher) PublishDKGResult(resultToPublish *relayChain.DKGResult) error {
	chainRelay := pm.chainHandle.ThresholdRelay()

	onPublishedResultChan := make(chan *event.DKGResultPublication)
	subscription := chainRelay.OnDKGResultPublished(func(publishedResult *event.DKGResultPublication) {
		onPublishedResultChan <- publishedResult
	})
	defer subscription.Unsubscribe()

	// Check if any result has already been published to the chain with current
	// request ID.
	if chainRelay.IsDKGResultPublished(pm.RequestID) {
		return nil
	}

	blockCounter, err := pm.chainHandle.BlockCounter()
	if err != nil {
		return fmt.Errorf("could not initialize block counter [%v]", err)
	}

	// Waits until the current member is eligible to submit a result to the
	// blockchain.
	eligibleToSubmitWaiter, err := blockCounter.BlockWaiter(
		(pm.publishingIndex - 1) * pm.blockStep,
	)
	if err != nil {
		return fmt.Errorf("block waiter failure [%v]", err)
	}

	for {
		select {
		case <-eligibleToSubmitWaiter:
			errors := make(chan error)
			chainRelay.SubmitDKGResult(pm.RequestID, resultToPublish).
				OnComplete(func(resultPublicationEvent *event.DKGResultPublication, err error) {
					errors <- err
				})
			return <-errors
		case publishedResultEvent := <-onPublishedResultChan:
			if publishedResultEvent.RequestID.Cmp(pm.RequestID) == 0 {
				return nil
			}
		}
	}
}
