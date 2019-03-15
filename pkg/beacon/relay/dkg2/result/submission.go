package result

import (
	"fmt"
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

// SubmittingMember represents a member submitting a DKG result to the
// blockchain along with signatures received from other group members supporting
// the result.
type SubmittingMember struct {
	index gjkr.MemberID

	chainHandle chain.Handle
	// Predefined step for each submitting window. The value is used to determine
	// eligible submitting member.
	blockStep uint32
}

// SubmitDKGResult is ... TODO: write documentation
func (sm *SubmittingMember) SubmitDKGResult(
	requestID *big.Int,
	result *relayChain.DKGResult,
	signatures map[gjkr.MemberID]operator.Signature,
) (int64, error) {
	chainRelay := sm.chainHandle.ThresholdRelay()
	blockCounter, err := sm.chainHandle.BlockCounter()
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
		int((sm.index - 1)) * int(sm.blockStep),
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
				uint32(sm.index),
				result,
				signatures,
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
