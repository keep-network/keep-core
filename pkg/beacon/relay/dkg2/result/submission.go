package result

import (
	"fmt"
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

// SubmitDKGResult is ... TODO: write documentation
func (rsm *ResultSubmittingMember) SubmitDKGResult(
	requestID *big.Int,
	result *relayChain.DKGResult,
) (int64, error) {
	chainRelay := rsm.chainHandle.ThresholdRelay()
	blockCounter, err := rsm.chainHandle.BlockCounter()
	if err != nil {
		return -1, err
	}

	onSubmittedResultChan := make(chan *event.DKGResultSubmission)

	subscription, err := chainRelay.OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			onSubmittedResultChan <- event
		},
	)
	if err != nil {
		close(onSubmittedResultChan)
		return -1, fmt.Errorf(
			"could not watch for DKG result publications [%v]",
			err,
		)
	}

	// Check if any result has already been published to the chain with current
	// request ID.
	alreadyPublished, err := chainRelay.IsDKGResultSubmitted(requestID)
	if err != nil {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
		return -1, fmt.Errorf(
			"could not check if the result is already published [%v]",
			err,
		)
	}

	// Someone who was ahead of us in the queue published the result. Giving up.
	if alreadyPublished {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
		return -1, nil
	}

	// Waits until the current member is eligible to submit a result to the
	// blockchain.
	eligibleToSubmitWaiter, err := blockCounter.BlockWaiter(
		int((rsm.index - 1)) * int(rsm.blockStep),
	)
	if err != nil {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
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
			close(onSubmittedResultChan)

			chainRelay.SubmitDKGResult(
				requestID,
				uint32(rsm.index),
				result,
				rsm.receivedValidResultSignatures,
			).
				OnComplete(func(
					dkgResultPublishedEvent *event.DKGResultSubmission,
					err error,
				) {
					if err != nil {
						errorChannel <- err
					}
					blockHeight <- dkgResultPublishedEvent.BlockNumber
					errorChannel <- nil
				})
			return int64(<-blockHeight), <-errorChannel
		case publishedResultEvent := <-onSubmittedResultChan:
			if publishedResultEvent.RequestID.Cmp(requestID) == 0 {
				subscription.Unsubscribe()
				close(onSubmittedResultChan)
				return int64(publishedResultEvent.BlockNumber), nil // leave without publishing the result
			}
		}
	}
}
