package dkg2

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
	// Initialized BlockCounter which allows for the reading, counting, and
	// waiting of blocks for the purposes of synchronization.
	blockCounter chain.BlockCounter
	// Sequential number of the current member in the publishing group.
	// The value is used to determine eligible publishing member. Indexing starts
	// with `1`. Relates to DKG Phase 13.
	publishingIndex int
	// Predefined step for each publishing window. The value is used to determine
	// eligible publishing member. Relates to DKG Phase 13.
	blockStep uint64
	// Tracks if member still has a right to submit DKG result or vote to the
	// chain. Member is eligible to only one submission.
	alreadySubmitted bool
	// Result conflict resolution duration.
	conflictDuration uint64
	// Maximum number of malicious members.
	dishonestThreshold int
}

// executePublishing runs Distributed Key Generation result publication and voting,
// given unique identifier of DKG execution, a player index in the group, handler
// to interact with a chain and the Distributed Key Generation result in a format
// accepted by the chain.
func executePublishing(
	requestID *big.Int,
	publishingIndex int,
	dishonestThreshold int,
	chainRelay relayChain.Interface,
	blockCounter chain.BlockCounter,
	result *relayChain.DKGResult,
) error {
	if publishingIndex < 1 {
		return fmt.Errorf("publishing index must be >= 1")
	}

	publisher := &Publisher{
		RequestID:          requestID,
		blockCounter:       blockCounter,
		publishingIndex:    publishingIndex,
		blockStep:          1,
		dishonestThreshold: dishonestThreshold,
	}

	blockHeight, err := publisher.publishResult(result, chainRelay)
	if err != nil {
		return fmt.Errorf("result publication failed [%v]", err)
	}
	if blockHeight < 0 {
		return fmt.Errorf("block height is less than zero [%v]", blockHeight)
	}

	// TODO Execute Phase 14 here

	return nil
}

// publishResult sends a result containing i.a. group public key to the blockchain.
// It checks if the result has already been published to the blockchain with
// request ID specific for current DKG execution. If not, it determines if the
// current member is eligable to result submission. If allowed, it submits the
// results to the blockchain.
//
// User allowance to publish is determined based on the user's publishing index
// and publishing block step.
//
// When member is waiting for their round the function keeps tracking results being
// published to the blockchain. If any result is published for the current
// request ID, the current member finishes the phase immediately, without
// publishing its own result.
//
// It returns chain block height of the moment when the result was published on
// chain in case the result has been already published by another publisher it
// returns current block height. In case of failure it returns `-1`.
//
// See Phase 13 of the protocol specification.
func (pm *Publisher) publishResult(
	result *relayChain.DKGResult,
	chainRelay relayChain.Interface,
) (int64, error) {
	onPublishedResultChan := make(chan *event.DKGResultPublication)

	subscription, err := chainRelay.OnDKGResultPublished(
		func(publishedResult *event.DKGResultPublication) {
			onPublishedResultChan <- publishedResult
		},
	)
	if err != nil {
		close(onPublishedResultChan)
		return -1, fmt.Errorf(
			"could not watch for DKG result publications [%v]",
			err,
		)
	}

	// Check if any result has already been published to the chain with current
	// request ID.
	alreadyPublished, err := chainRelay.IsDKGResultPublished(pm.RequestID)
	if err != nil {
		subscription.Unsubscribe()
		close(onPublishedResultChan)
		return -1, fmt.Errorf(
			"could not check if the result is already published [%v]",
			err,
		)
	}

	// Someone who was ahead of us in the queue published the result. Giving up.
	if alreadyPublished {
		// TODO: Should `IsDKGResultPublished` return block height when the result was published?
		// We wouldn't have to return currentBlock then.
		currentBlock, err := pm.blockCounter.CurrentBlock()

		subscription.Unsubscribe()
		close(onPublishedResultChan)

		return int64(currentBlock), err
	}

	// Waits until the current member is eligible to submit a result to the
	// blockchain.
	eligibleToSubmitWaiter, err := pm.blockCounter.BlockWaiter(
		(pm.publishingIndex - 1) * int(pm.blockStep),
	)
	if err != nil {
		subscription.Unsubscribe()
		close(onPublishedResultChan)
		return -1, fmt.Errorf("block waiter failure [%v]", err)
	}

	for {
		select {
		case <-eligibleToSubmitWaiter:
			blockHeight := make(chan uint64)
			defer close(blockHeight)
			errorChannel := make(chan error)
			defer close(errorChannel)

			subscription.Unsubscribe()
			close(onPublishedResultChan)

			chainRelay.SubmitDKGResult(pm.RequestID /*, pm.publishingIndex,*/, result).
				OnSuccess(func(dkgResultPublishedEvent *event.DKGResultPublication) {
					blockHeight <- dkgResultPublishedEvent.BlockNumber
					// TODO: This is a temporary solution until DKG Phase 14 is
					// ready. We assume that only one DKG result is published in
					// DKG Phase 13 and submit it as a final group public key.

					chainRelay.SubmitGroupPublicKey(
						pm.RequestID,
						dkgResultPublishedEvent.GroupPublicKey,
					).OnSuccess(func(groupRegisteredEvent *event.GroupRegistration) {
						fmt.Printf(
							"Group public key submitted for requestID=[%v]\n",
							pm.RequestID,
						)
						errorChannel <- nil
					}).OnFailure(func(err error) {
						errorChannel <- err
					})
				}).
				OnFailure(func(err error) {
					errorChannel <- err
				})

			pm.alreadySubmitted = true
			return int64(<-blockHeight), <-errorChannel
		case publishedResultEvent := <-onPublishedResultChan:
			if publishedResultEvent.RequestID.Cmp(pm.RequestID) == 0 {
				subscription.Unsubscribe()
				close(onPublishedResultChan)
				return int64(publishedResultEvent.BlockNumber), nil // leave without publishing the result
			}
		}
	}
}
