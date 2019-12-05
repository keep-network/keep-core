package entry

import (
	"fmt"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
)

type relayEntrySubmitter struct {
	chain        relayChain.Interface
	blockCounter chain.BlockCounter

	index group.MemberIndex
}

// submitRelayEntry submits the provided relay entry data to the chain.
// Group members tries to submit in the order specified by their indexes.
// Group member with index 1 tries to submit as the first one, group member 2
// tries to submit after a few blocks if member 1 did not submit and so on.
// Relay entry submit process starts at block height defined by startBlockheight
// parameter.
func (res *relayEntrySubmitter) submitRelayEntry(
	newEntry []byte,
	groupPublicKey []byte,
	startBlockHeight uint64,
) error {
	config, err := res.chain.GetConfig()
	if err != nil {
		return fmt.Errorf(
			"could not fetch chain's config: [%v]",
			err,
		)
	}

	onSubmittedResultChan := make(chan uint64)

	subscription, err := res.chain.OnRelayEntrySubmitted(
		func(event *event.EntrySubmitted) {
			onSubmittedResultChan <- event.BlockNumber
		},
	)
	if err != nil {
		close(onSubmittedResultChan)
		return fmt.Errorf("could not watch for relay entry submissions: [%v]", err)
	}

	returnWithError := func(err error) error {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
		return err
	}

	// TODO: we should eventually check if entry has been already submitted
	// but we may skip this check for V1.

	// Wait until the current member is eligible to submit the entry.
	eligibleToSubmitWaiter, err := res.waitForSubmissionEligibility(
		startBlockHeight,
		config.ResultPublicationBlockStep,
	)
	if err != nil {
		return returnWithError(
			fmt.Errorf("wait for eligibility failure: [%v]", err),
		)
	}

	for {
		select {
		case blockNumber := <-eligibleToSubmitWaiter:
			// Member becomes eligible to submit the result.
			errorChannel := make(chan error)
			defer close(errorChannel)

			subscription.Unsubscribe()
			close(onSubmittedResultChan)

			logger.Infof(
				"[member:%v] submitting relay entry [0x%x] on behalf of group "+
					"[0x%x] at block [%v]",
				res.index,
				newEntry,
				groupPublicKey,
				blockNumber,
			)

			res.chain.SubmitRelayEntry(newEntry).OnComplete(
				func(entry *event.EntrySubmitted, err error) {
					if err == nil {
						logger.Infof(
							"[member:%v] successfully submitted relay entry at block: [%v]",
							res.index,
							entry.BlockNumber,
						)
					}
					errorChannel <- err
				})
			return <-errorChannel
		case blockNumber := <-onSubmittedResultChan:
			logger.Infof(
				"[member:%v] leaving; relay entry submitted by other member at block [%v]",
				res.index,
				blockNumber,
			)
			return returnWithError(nil)
		}
	}
}

// waitForSubmissionEligibility waits until the current member is eligible to
// submit entry to the blockchain. First member is eligible to submit straight
// away, each following member is eligible after pre-defined block step.
func (res *relayEntrySubmitter) waitForSubmissionEligibility(
	startBlockHeight uint64,
	blockStep uint64,
) (<-chan uint64, error) {
	// (member_index - 1) * T_step
	blockWaitTime := (uint64(res.index) - 1) * blockStep

	eligibleBlockHeight := startBlockHeight + blockWaitTime
	logger.Infof(
		"[member:%v] waiting for block [%v] to submit",
		res.index,
		eligibleBlockHeight,
	)

	waiter, err := res.blockCounter.BlockHeightWaiter(eligibleBlockHeight)
	if err != nil {
		return nil, fmt.Errorf("block height waiter failure: [%v]", err)
	}

	return waiter, err
}
