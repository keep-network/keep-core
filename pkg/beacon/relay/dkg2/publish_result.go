package dkg2

import (
	"fmt"
	"math/big"
	"sync"

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

	_, err = publisher.resultConflictResolution(result, chainRelay, uint64(blockHeight))
	if err != nil {
		return fmt.Errorf("result conflict resolution failed [%v]", err)
	}

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

// resultConflictResolution executes conflict resolution if the member considers
// other than currently leading on-chain DKG result to be a corrent one.
//
// If the correct DKG result is not yet submitted to the chain the member will
// submit it. Otherwise member will vote for a hash of the currently submitted
// result.
//
// Each member is allowed to exactly one submission or vote. If the member
// submitted a DKG result in the previous phase, their won't be able to submit
// or vote again in this phase.
//
// It requires starting block height to be provided as a reference when the
// phase begins. The value is block height of the previous phase end.
//
// It returns chain block height of the moment when the last DKG result was
// submitted or voted. In case the result has the majority of votes it returns
// current block height. In case of failure it returns `-1`.
//
// See Phase 14 of the protocol specification.
func (pm *Publisher) resultConflictResolution(
	correctResult *relayChain.DKGResult,
	chainRelay relayChain.Interface,
	startingBlockHeight uint64,
) (int64, error) {
	onVoteChan := make(chan *event.DKGResultVote)
	defer close(onVoteChan)
	onVoteSubscription, err := chainRelay.OnDKGResultVote(
		func(vote *event.DKGResultVote) {
			onVoteChan <- vote
		},
	)
	if err != nil {
		return -1, fmt.Errorf("could not watch for DKG result vote [%v]", err)
	}
	defer onVoteSubscription.Unsubscribe()

	onSubmissionChan := make(chan *event.DKGResultPublication)
	defer close(onSubmissionChan)
	onSubmissionSubscription, err := chainRelay.OnDKGResultPublished(
		func(result *event.DKGResultPublication) {
			onSubmissionChan <- result
		},
	)
	if err != nil {
		return -1, fmt.Errorf("could not watch for DKG result vote [%v]", err)
	}
	defer onSubmissionSubscription.Unsubscribe()

	errorChannel := make(chan error)
	defer close(errorChannel)

	resultsVotes := dkgResultsVotes(chainRelay.GetDKGResultsVotes(pm.RequestID))
	if resultsVotes == nil || len(resultsVotes) == 0 {
		return -1, fmt.Errorf("nothing submitted")
	}

	if resultsVotes.leadHasEnoughVotes(pm.dishonestThreshold) {
		currentBlock, err := pm.blockCounter.CurrentBlock()
		if err != nil {
			return -1, err
		}

		fmt.Printf(
			"[publisher: %v] Lead has enough votes.\n",
			pm.publishingIndex,
		)
		return int64(currentBlock), nil
	}

	blockNumber := int64(-1)

	correctResultHash, err := chainRelay.CalculateDKGResultHash(correctResult)
	if err != nil {
		return -1, fmt.Errorf("could not calculate dkg result hash [%v]", err)
	}

	if !resultsVotes.contains(correctResultHash) && !pm.alreadySubmitted {
		blockNumberChan := make(chan uint64)
		fmt.Printf(
			"[publisher: %v] Initial check: Result not submitted yet.\n",
			pm.publishingIndex,
		)
		onSubmissionSubscription.Unsubscribe()

		chainRelay.SubmitDKGResult(pm.RequestID, correctResult).
			OnComplete(
				func(dkgResultPublishedEvent *event.DKGResultPublication, err error) {
					if dkgResultPublishedEvent != nil {
						blockNumberChan <- dkgResultPublishedEvent.BlockNumber
					}
					errorChannel <- nil
				},
			)
		blockNumber = int64(<-blockNumberChan)
		pm.alreadySubmitted = true

		err := <-errorChannel
		if err != nil {
			return -1, err
		}
	}

	// TODO: Timeout should be extended after receiving last minute votes.
	phaseTimeout := int(startingBlockHeight + pm.conflictDuration)
	phaseDurationWaiter, err := pm.blockCounter.BlockHeightWaiter(phaseTimeout)
	if err != nil {
		return -1, fmt.Errorf("block waiter failure [%v]", err)
	}

	// Returns already submitted
	votesAndSubmissions := func(
		blockNumber uint64,
		chainRelay relayChain.Interface,
	) (bool, int64, error) {
		blockNumberChan := make(chan uint64)

		submissions := dkgResultsVotes(chainRelay.GetDKGResultsVotes(pm.RequestID))

		if submissions.leadHasEnoughVotes(pm.dishonestThreshold) {
			fmt.Printf("[publisher: %v] Lead has enough votes.\n", pm.publishingIndex)
			return false, int64(blockNumber), nil
		}

		if submissions.isStrictlyLeading(correctResultHash) {
			fmt.Printf(
				"[publisher: %v] Result is the only lead.\n",
				pm.publishingIndex,
			)
			return false, int64(blockNumber), nil
		}

		if submissions.contains(correctResultHash) {
			fmt.Printf("[publisher: %v] Vote for the result.\n", pm.publishingIndex)
			chainRelay.VoteOnDKGResult(
				pm.RequestID, pm.publishingIndex,
				correctResultHash,
			).OnComplete(func(dkgResultVote *event.DKGResultVote, err error) {
				if dkgResultVote != nil {
					blockNumberChan <- dkgResultVote.BlockNumber
				}
				errorChannel <- nil
			})
			return true, int64(<-blockNumberChan), <-errorChannel
		}

		fmt.Printf("[publisher: %v] Submit the result.\n", pm.publishingIndex)
		chainRelay.SubmitDKGResult(pm.RequestID, correctResult).
			OnComplete(
				func(dkgResultPublishedEvent *event.DKGResultPublication, err error) {
					if dkgResultPublishedEvent != nil {
						blockNumberChan <- dkgResultPublishedEvent.BlockNumber
					}
					errorChannel <- nil
				},
			)
		return true, int64(<-blockNumberChan), <-errorChannel
	}

	votesAndSubmissionsMutex := &sync.Mutex{}

	for {
		select {
		case <-phaseDurationWaiter:
			fmt.Printf("[publisher: %v] Result conflict resolution timeout.\n", pm.publishingIndex)
			return blockNumber, nil
		case vote := <-onVoteChan:
			votesAndSubmissionsMutex.Lock()

			if vote.RequestID.Cmp(pm.RequestID) == 0 {
				fmt.Printf("[publisher: %v] Vote event received.\n", pm.publishingIndex)
				blockNumber = int64(vote.BlockNumber)

				if !pm.alreadySubmitted {
					pm.alreadySubmitted, blockNumber, err = votesAndSubmissions(vote.BlockNumber, chainRelay)
					if err != nil {
						votesAndSubmissionsMutex.Unlock()
						return blockNumber, err
					}
				}
			}
			votesAndSubmissionsMutex.Unlock()
		case submission := <-onSubmissionChan:
			votesAndSubmissionsMutex.Lock()

			if submission.RequestID.Cmp(pm.RequestID) == 0 {
				fmt.Printf("[publisher: %v] Submission event received.\n", pm.publishingIndex)
				blockNumber = int64(submission.BlockNumber)

				if !pm.alreadySubmitted {
					pm.alreadySubmitted, blockNumber, err = votesAndSubmissions(submission.BlockNumber, chainRelay)
					if err != nil {
						votesAndSubmissionsMutex.Unlock()
						return blockNumber, err
					}
				}
			}
			votesAndSubmissionsMutex.Unlock()
		}
	}
}
