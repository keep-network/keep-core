package result

import (
	"fmt"
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

// SubmittingMember represents a member submitting a DKG result to the
// blockchain along with signatures received from other group members supporting
// the result.
type SubmittingMember struct {
	// Represents the member's position for submission.
	index member.Index

	// Predefined step for each submitting window. The value is used to determine
	// the eligible submitting member.
	blockStep uint64
}

// SubmitDKGResult sends a result, which contains the group public key and
// signatures, to the chain.
//
// It checks if the result has already been published to the blockchain with
// the request ID specific to the current DKG execution. If not, it determines if
// the current member is eligible to submit a result. If allowed, it submits
// the result to the chain.
//
// A user's turn to publish is determined based on the user's index and block
// step.
//
// If a result is submitted for the current request ID and it's accepted by the
// chain, the current member finishes the phase immediately, without submitting
// their own result.
//
// It returns the on-chain block height of the moment when the result was
// successfully submitted on chain by the member. In case of failure or result
// already submitted by another member it returns `0`.
//
// See Phase 14 of the protocol specification.
func (sm *SubmittingMember) SubmitDKGResult(
	requestID *big.Int,
	result *relayChain.DKGResult,
	signatures map[member.Index]operator.Signature,
	chainHandle chain.Handle,
) error {
	onSubmittedResultChan := make(chan *event.DKGResultSubmission)

	chainRelay := chainHandle.ThresholdRelay()
	subscription, err := chainRelay.OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			onSubmittedResultChan <- event
		},
	)
	if err != nil {
		close(onSubmittedResultChan)
		return fmt.Errorf(
			"could not watch for DKG result publications [%v]",
			err,
		)
	}

	returnWithError := func(err error) error {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
		return err
	}

	// Check if any result has already been published to the chain with current
	// request ID.
	alreadyPublished, err := chainRelay.IsDKGResultSubmitted(requestID)
	if err != nil {
		return returnWithError(
			fmt.Errorf(
				"could not check if the result is already published [%v]",
				err,
			),
		)
	}

	// Someone who was ahead of us in the queue published the result. Giving up.
	if alreadyPublished {
		return returnWithError(nil)
	}

	// Wait until the current member is eligible to submit the result.
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return returnWithError(err)
	}
	eligibleToSubmitWaiter, err := sm.waitForSubmissionEligibility(blockCounter)
	if err != nil {
		return returnWithError(
			fmt.Errorf("wait for eligibility failure [%v]", err),
		)
	}

	for {
		select {
		case <-eligibleToSubmitWaiter:
			// Member becomes eligible to submit the result.
			errorChannel := make(chan error)
			defer close(errorChannel)

			subscription.Unsubscribe()
			close(onSubmittedResultChan)

			chainRelay.SubmitDKGResult(
				requestID,
				sm.index,
				result,
				signatures,
			).
				OnComplete(func(
					dkgResultPublishedEvent *event.DKGResultSubmission,
					err error,
				) {
					errorChannel <- nil
				})
			return <-errorChannel
		case publishedResultEvent := <-onSubmittedResultChan:
			if publishedResultEvent.RequestID.Cmp(requestID) == 0 {
				// A result has been submitted by other member. Leave without
				// publishing the result.
				return returnWithError(nil)
			}
		}
	}
}

// waitForSubmissionEligibility waits until the current member is eligible to
// submit a result to the blockchain. First member is eligible to submit straight
// away, each following member is eligible after pre-defined block step.
func (sm *SubmittingMember) waitForSubmissionEligibility(blockCounter chain.BlockCounter) (<-chan int, error) {
	eligibleToSubmitWaiter, err := blockCounter.BlockWaiter(
		(int(sm.index) - 1) * int(sm.blockStep), // T_init + (member_index - 1) * T_step
	)
	if err != nil {
		return nil, fmt.Errorf("block waiter failure [%v]", err)
	}

	return eligibleToSubmitWaiter, err
}
