package entry

import (
	"fmt"
	"math/big"
	"time"

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
	signingID *big.Int,
	newEntry *big.Int,
	previousEntry *big.Int,
	seed *big.Int,
	groupPublicKey []byte,
	startBlockHeight uint64,
) error {
	config, err := res.chain.GetConfig()
	if err != nil {
		return fmt.Errorf(
			"could not fetch chain's config [%v]",
			err,
		)
	}

	onSubmittedResultChan := make(chan *event.Entry)

	subscription, err := res.chain.OnSignatureSubmitted(
		func(event *event.Entry) {
			onSubmittedResultChan <- event
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

			fmt.Printf(
				"[member:%v] Submitting relay entry on behalf of the group [%v]...\n",
				res.index,
				groupPublicKey,
			)
			entry := &event.Entry{
				SigningId:     signingID,
				Value:         newEntry,
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
				GroupPubKey:   groupPublicKey,
				Seed:          seed,
			}

			res.chain.SubmitRelayEntry(entry).OnComplete(
				func(entry *event.Entry, err error) {
					if err == nil {
						fmt.Printf(
							"[member:%v] Relay entry for request [%v] successfully submitted at block [%v]\n",
							res.index,
							signingID,
							entry.BlockNumber,
						)
					}
					errorChannel <- err
				})
			return <-errorChannel
		case submittedEntryEvent := <-onSubmittedResultChan:
			if submittedEntryEvent.SigningId.Cmp(signingID) == 0 {
				fmt.Printf(
					"[member:%v] Relay entry submitted by other member, leaving.\n",
					res.index,
				)
				return returnWithError(nil)
			}
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
	fmt.Printf(
		"[member:%v] Waiting for block [%v] to submit...\n",
		res.index,
		eligibleBlockHeight,
	)

	waiter, err := res.blockCounter.BlockHeightWaiter(eligibleBlockHeight)
	if err != nil {
		return nil, fmt.Errorf("block height waiter failure [%v]", err)
	}

	return waiter, err
}
