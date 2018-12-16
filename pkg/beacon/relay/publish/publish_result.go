package publish

import (
	"fmt"
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// Publisher is a member submitting distributed key generation result to a
// blockchain.
type Publisher struct {
	ID gjkr.MemberID
	// ID of distributed key generation execution.
	RequestID *big.Int
	// Handle to interact with a blockchain.
	chainHandle chain.Handle
	// Sequential number of the current member in the publishing group.
	// The value is used to determine eligible publishing member. Relates to DKG
	// Phase 13.
	publishingIndex int
	// Predefined step for each publishing window. The value is used to determine
	// eligible publishing member. Relates to DKG Phase 13.
	blockStep int

	conflictDuration int // T_conflict
	votingThreshold  int // T_max
}

// PublishDKGResult sends a result containing i.a. group public key to the blockchain.
// It checks if the result has already been published to the blockchain with
// request ID specific for current DKG execution. If not it determines if the
// current member is eligable to result submission. If allowed it submits the
// results to the blockchain. The function returns result published
// to the blockchain containing ID of the member who published it.
//
// See Phase 13 of the protocol specification.
func (pm *Publisher) PublishDKGResult(resultToPublish *relayChain.DKGResult) error {
	chainRelay := pm.chainHandle.ThresholdRelay()

	onPublishedResultChan := make(chan *event.DKGResultPublication)
	chainRelay.OnDKGResultPublished(func(publishedResult *event.DKGResultPublication) {
		onPublishedResultChan <- publishedResult
	})

	// Check if the result has already been published to the chain.
	if chainRelay.IsDKGResultPublished(pm.RequestID, resultToPublish) {
		return nil // TODO What should we return here? Should it be an error?
	}

	blockCounter, err := pm.chainHandle.BlockCounter()
	if err != nil {
		return fmt.Errorf("block counter failure [%v]", err)
	}

	// Waits until the current member is eligible to submit a result to the
	// blockchain.
	eligibleToSubmitWaiter, err := blockCounter.BlockWaiter(
		pm.publishingIndex * pm.blockStep,
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
				if chainRelay.IsDKGResultPublished(pm.RequestID, resultToPublish) {
					return nil
				}
			}
		}
	}
}

/*
// Phase14 - PHASE 14
func (pm *Publisher) BlockCounter() (int, error) {
	blockCounter, err := pm.chainHandle.BlockCounter()
	if err != nil {
		return 0, fmt.Errorf("block counter failure [%v]", err)
	}
	return blockCounter, nil
}
*/

// Phase14 - PHASE 14
func (pm *Publisher) Phase14(correctResult *relayChain.DKGResult) error {
	chainRelay := pm.chainHandle.ThresholdRelay()

	fmt.Printf("%sAt: %s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
	onVoteChan := make(chan *event.DKGResultVote)
	chainRelay.OnDKGResultVote(func(vote *event.DKGResultVote) {
		fmt.Printf("%sAt: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
		onVoteChan <- vote
	})
	onSubmissionChan := make(chan *event.DKGResultPublication)
	chainRelay.OnDKGResultPublished(func(result *event.DKGResultPublication) {
		onSubmissionChan <- result
	})

	fmt.Printf("At: %s\n", godebug.LF())
	if pm.RequestID == nil {
		fmt.Printf("At: %s\n", godebug.LF())
	}
	submissions := chainRelay.GetDKGSubmissions(pm.RequestID)
	if submissions == nil {
		fmt.Printf("At: %s\n", godebug.LF())
		return fmt.Errorf("nothing submitted")
	}
	if !nOfVotesBelowThreshold(submissions, pm.votingThreshold) {
		fmt.Printf("At: %s\n", godebug.LF()) // <<<<<<<<<<<<<<<<
		return fmt.Errorf("voting threshold exceeded")
	}
	fmt.Printf("At: %s\n", godebug.LF())
	if !submissions.Contains(correctResult) {
		fmt.Printf("At: %s --- Submissions did not contain the 'correctResult' passed, %s\n", godebug.LF(), godebug.SVarI(correctResult))
		chainRelay.SubmitDKGResult(pm.RequestID, correctResult)
		return nil
	}

	fmt.Printf("At: %s\n", godebug.LF())
	blockCounter, err := pm.chainHandle.BlockCounter()
	if err != nil {
		fmt.Printf("At: %s\n", godebug.LF())
		return fmt.Errorf("block counter failure [%v]", err)
	}

	fmt.Printf("At: %s\n", godebug.LF())
	// firstBlock := 0 // T_First
	firstBlock, err := blockCounter.CurrentBlock() // T_First
	if err != nil {
		fmt.Printf("At: %s\n", godebug.LF())
		return fmt.Errorf("current block failure [%v]", err)
	}

	fmt.Printf("At: %s\n", godebug.LF())
	// NOTE: We wait for T_conflict blocks but the protocol specification states
	// that we should wait for block `T_first + T_conflict`. Need clarification.
	phaseDurationWaiter, err := blockCounter.BlockWaiter(firstBlock + pm.conflictDuration)
	if err != nil {
		fmt.Printf("At: %s\n", godebug.LF())
		return fmt.Errorf("block waiter failure [%v]", err)
	}

	fmt.Printf("At: %s\n", godebug.LF())
	votesAndSubmissions := func(chainRelay relayChain.Interface) (bool, error) {
		fmt.Printf("At: %s\n", godebug.LF())
		submissions := chainRelay.GetDKGSubmissions(pm.RequestID)
		if !nOfVotesBelowThreshold(submissions, pm.votingThreshold) {
			fmt.Printf("At: %s\n", godebug.LF())
			return true, fmt.Errorf("voting threshold exceeded")
		}

		fmt.Printf("At: %s\n", godebug.LF())
		if !submissions.Lead().DKGResult.Equals(correctResult) {
			fmt.Printf("At: %s\n", godebug.LF())
			chainRelay.Vote(pm.RequestID, correctResult.Hash())
			return true, nil
		} else if !submissions.Contains(correctResult) {
			fmt.Printf("At: %s\n", godebug.LF())
			chainRelay.SubmitDKGResult(pm.RequestID, correctResult)
			return true, nil
		}
		fmt.Printf("At: %s\n", godebug.LF())
		return false, nil
	}

	for {
		fmt.Printf("At: %s\n", godebug.LF())
		select {
		case <-phaseDurationWaiter:
			fmt.Printf("At: %s\n", godebug.LF())
			return nil
		case vote := <-onVoteChan:
			fmt.Printf("%sAt: %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
			if vote.RequestID.Cmp(pm.RequestID) == 0 {
				fmt.Printf("At: %s\n", godebug.LF())
				if result, err := votesAndSubmissions(chainRelay); result {
					fmt.Printf("At: %s\n", godebug.LF())
					return err
				}
			}
		case submission := <-onSubmissionChan:
			fmt.Printf("%sAt: %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
			if submission.RequestID.Cmp(pm.RequestID) == 0 {
				fmt.Printf("At: %s\n", godebug.LF())
				if result, err := votesAndSubmissions(chainRelay); result {
					fmt.Printf("%sAt: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
					return err
				}
			}
		}
	}

}

func nOfVotesBelowThreshold(submissions *relayChain.Submissions, votingThreshold int) bool {
	fmt.Printf("At: %s\n\tCalled From %s\n\tsubmissions=%s\n", godebug.LF(), godebug.LF(2), godebug.SVarI(submissions))
	if submissions.Submissions == nil {
		fmt.Printf("At: %s -- submissions= NIL - no votes yet\n", godebug.LF())
		return true
	}
	return submissions.Lead().Votes <= votingThreshold // leadResult.votes > M_max
}
