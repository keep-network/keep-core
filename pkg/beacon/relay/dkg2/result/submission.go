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
        // Represents the member's position for submission.
	index gjkr.MemberID

	// Interface to access submission methods specific to the chain.
	chainHandle chain.Handle
	// Predefined step for each submitting window. The value is used to determine
	// the eligible submitting member.
	blockStep uint64
}

// SubmitDKGResult sends a result containing i.a. group public key and signatures
// supporting this result to the blockchain.
//
// It checks if the result has already been published to the blockchain with
// request ID specific for current DKG execution. If not, it determines if the
// current member is eligable to result submission. If allowed, it submits the
// results to the blockchain.
//
// User allowance to publish is determined based on the user's index and block
// step.
//
// When member is waiting for their round the function keeps tracking results being
// submitted to the blockchain. If any result is submitted for the current
// request ID, the current member finishes the phase immediately, without
// submitting its own result.
//
// It returns chain block height of the moment when the result was successfully
// submitted on chain by the member. In case of failure or result already
// submitted by another member it returns `0`.
//
// See Phase 14 of the protocol specification.
func (sm *SubmittingMember) SubmitDKGResult(
	requestID *big.Int,
	result *relayChain.DKGResult,
	signatures map[gjkr.MemberID]operator.Signature,
) (uint64, error) {
	onSubmittedResultChan := make(chan *event.DKGResultSubmission)

	chainRelay := sm.chainHandle.ThresholdRelay()
	subscription, err := chainRelay.OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			onSubmittedResultChan <- event
		},
	)
	if err != nil {
		close(onSubmittedResultChan)
		return 0, fmt.Errorf(
			"could not watch for DKG result publications [%v]",
			err,
		)
	}

	blockCounter, err := sm.chainHandle.BlockCounter()
	if err != nil {
		return 0, err
	}

	// Check if any result has already been published to the chain with current
	// request ID.
	alreadyPublished, err := chainRelay.IsDKGResultSubmitted(requestID)
	if err != nil {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
		return 0, fmt.Errorf(
			"could not check if the result is already published [%v]",
			err,
		)
	}

	// Someone who was ahead of us in the queue published the result. Giving up.
	if alreadyPublished {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
		// TODO: Should we return block height of the moment the result was submitted
		// or current block?
		currentBlock, err := blockCounter.CurrentBlock()
		if err != nil {
			return 0, err
		}

		return uint64(currentBlock), nil
	}

	// Waits until the current member is eligible to submit a result to the
	// blockchain.
	// TODO: Check if we need to use BlockHeighWaiter. To do that we would need
	// to pass block height when previous phase ended so we can synchronize.
	eligibleToSubmitWaiter, err := blockCounter.BlockWaiter(
		int((sm.index - 1)) * int(sm.blockStep),
	)
	if err != nil {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
		return 0, fmt.Errorf("block waiter failure [%v]", err)
	}

	for {
		select {
		case eligibleToSubmitBlock := <-eligibleToSubmitWaiter:
			// Member becomes eligible to submit the result.
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
				OnFailure(func(err error) {
					// Block height when member became eligible to submit.
					blockHeight <- uint64(eligibleToSubmitBlock)
					errorChannel <- err
				}).
				OnSuccess(func(
					dkgResultPublishedEvent *event.DKGResultSubmission,
				) {
					// Block height when result was successfully submitted.
					blockHeight <- dkgResultPublishedEvent.BlockNumber
					errorChannel <- nil
				})
			return <-blockHeight, <-errorChannel
		case publishedResultEvent := <-onSubmittedResultChan:
			// A result has been submitted by other member.
			if publishedResultEvent.RequestID.Cmp(requestID) == 0 {
				subscription.Unsubscribe()
				close(onSubmittedResultChan)
				return publishedResultEvent.BlockNumber, nil // leave without publishing the result
			}
		}
	}
}
