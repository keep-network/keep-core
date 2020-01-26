// Package groupselection implements the random beacon group selection protocol
// - an interactive, ticket-based method of selecting a candidate group from
// the set of all stakers given a pseudorandom seed value.
package groupselection

import (
	"fmt"
	"math/big"

	"github.com/ipfs/go-log"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/chain"
)

var logger = log.Logger("keep-groupselection")

// CandidateToNewGroup attempts to generate and submit tickets for the staker to
// join a new group.
//
// There are two phases of ticket submission:
// - initial ticket submission,
// - reactive ticket submission.
//
// During the initial ticket submission, only tickets with a value below the
// natural threshold are submitted to the chain. Those tickets have the highest
// chance of being selected to the group and this way we minimize staker's
// gas expenditure.
//
// During the reactive ticket submission, all other staker's tickets are
// submitted. Reactive ticket submission is skipped if during the initial
// ticket submission there was enough tickets submitted to a chain to form
// a group. Those tickets could be submitted by any stakers participating in
// a new group selection.
//
// The function never submits more tickets than the group size.
func CandidateToNewGroup(
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	chainConfig *config.Chain,
	staker chain.Staker,
	newEntry *big.Int,
	startBlockHeight uint64,
	onGroupSelected func(*Result),
) error {
	availableStake, err := staker.Stake()
	if err != nil {
		return err
	}
	initialSubmissionTickets, reactiveSubmissionTickets, err :=
		generateTickets(
			newEntry.Bytes(),
			staker.ID(),
			availableStake,
			chainConfig.MinimumStake,
			naturalThreshold(chainConfig),
		)
	if err != nil {
		return err
	}

	logger.Infof(
		"generated [%v] tickets for initial submission phase and [%v] "+
			"tickets for reactive submission phase",
		len(initialSubmissionTickets),
		len(reactiveSubmissionTickets),
	)

	return startTicketSubmission(
		initialSubmissionTickets,
		reactiveSubmissionTickets,
		relayChain,
		blockCounter,
		chainConfig,
		startBlockHeight,
		onGroupSelected,
	)
}

func startTicketSubmission(
	initialSubmissionTickets []*ticket,
	reactiveSubmissionTickets []*ticket,
	relayChain relaychain.GroupSelectionInterface,
	blockCounter chain.BlockCounter,
	chainConfig *config.Chain,
	startBlockHeight uint64,
	onGroupSelected func(*Result),
) error {
	initialSubmissionTimeout, err := blockCounter.BlockHeightWaiter(
		startBlockHeight + chainConfig.TicketSubmissionTimeout/2,
	)
	if err != nil {
		return err
	}

	reactiveSubmissionTimeout, err := blockCounter.BlockHeightWaiter(
		startBlockHeight + chainConfig.TicketSubmissionTimeout,
	)
	if err != nil {
		return err
	}

	// Buffer quit signals - we never know if the goroutine finished
	// before we try to cancel it. The initial ticket submission may be
	// cancelled right after the initial submission timeout and after the
	// reactive submission timeout and it is possible it already completed.
	// Hence, we buffer two quit signals. The reactive ticket submission
	// is cancelled right after the reactive submission timeout. Here as well,
	// we do not know if the goroutine already completed, so we need to buffer
	// one quit signal.
	quitInitialTicketSubmission := make(chan struct{}, 2)
	quitReactiveTicketSubmission := make(chan struct{}, 1)

	// Check how many tickets with values below the natural threshold has been
	// generated and compare this number with the group size. Decide how many
	// tickets should be submitted. It does not make sense to submit more
	// tickets than the group size.
	var numberOfTicketsToSubmit int
	if len(initialSubmissionTickets) > chainConfig.GroupSize {
		numberOfTicketsToSubmit = chainConfig.GroupSize
	} else {
		numberOfTicketsToSubmit = len(initialSubmissionTickets)
	}

	logger.Infof(
		"entering initial ticket submission phase with [%v] tickets",
		numberOfTicketsToSubmit,
	)

	// Submit tickets with values below the natural threshold.
	// Do not submit more tickets than the group size.
	go submitTickets(
		initialSubmissionTickets[:numberOfTicketsToSubmit],
		relayChain,
		quitInitialTicketSubmission,
	)

	for {
		select {
		case initialSubmissionEndBlockHeight := <-initialSubmissionTimeout:
			// Initial ticket submission phase has ended. We need to determine
			// the total number of tickets submitted by all stakers who
			// candidate to a new group and decide whether to stop or to
			// enter reactive ticket submission.

			logger.Infof(
				"initial ticket submission ended at block [%v]",
				initialSubmissionEndBlockHeight,
			)

			ticketsCount, err := relayChain.GetSubmittedTicketsCount()
			if err != nil {
				return fmt.Errorf(
					"could not get submitted tickets count: [%v]",
					err,
				)
			}

			groupSize := big.NewInt(int64(chainConfig.GroupSize))
			if ticketsCount.Cmp(groupSize) >= 0 {
				// If there has been enough tickets submitted to form a new
				// group we stop ticket submission skipping the reactive ticket
				// submission phase.
				logger.Infof(
					"[%v] tickets submitted by group member candidates; "+
						"skipping reactive submission",
					ticketsCount,
				)

				quitInitialTicketSubmission <- struct{}{}
			} else {
				// If there has been not enough tickets submitted to form a new
				// group, we enter reactive ticket submission where we'll submit
				// remaining tickets. Note we are not stopping the goroutine
				// potentially still submitting tickets with values below the
				// initial threshold.
				// The number of remaining tickets is never larger than the
				// group size, including tickets with values below the natural
				// threshold.

				// Check how many tickets have been generated and compare this
				// value with the group size. Decide how many tickets should be
				// submitted. It does not make sense to submit more tickets
				// than the group size.
				if len(initialSubmissionTickets)+
					len(reactiveSubmissionTickets) > chainConfig.GroupSize {
					numberOfTicketsToSubmit = chainConfig.GroupSize -
						len(initialSubmissionTickets)
				} else {
					numberOfTicketsToSubmit = len(reactiveSubmissionTickets)
				}

				logger.Infof(
					"[%v] tickets submitted by group member candidates; "+
						"entering reactive submission phase with [%v] "+
						"additional tickets",
					ticketsCount,
					numberOfTicketsToSubmit,
				)

				// Submit tickets with values above the natural threshold.
				// Do not submit more tickets than the group size including
				// tickets with values below the natural threshold.
				go submitTickets(
					reactiveSubmissionTickets[:numberOfTicketsToSubmit],
					relayChain,
					quitReactiveTicketSubmission,
				)
			}

		case reactiveSubmissionEndBlockHeight := <-reactiveSubmissionTimeout:
			// Reactive ticket submission phase has ended. We need to quit two
			// potentially still running ticket submission goroutines, figure
			// out which stakers have been selected to the group and trigger
			// appropriate callback.

			logger.Infof(
				"reactive ticket submission ended at block [%v]",
				reactiveSubmissionEndBlockHeight,
			)

			quitInitialTicketSubmission <- struct{}{}
			quitReactiveTicketSubmission <- struct{}{}

			selectedStakers, err := relayChain.GetSelectedParticipants()
			if err != nil {
				return fmt.Errorf(
					"could not fetch selected participants after submission timeout [%v]",
					err,
				)
			}

			go onGroupSelected(&Result{
				SelectedStakers:        selectedStakers,
				GroupSelectionEndBlock: reactiveSubmissionEndBlockHeight,
			})

			return nil
		}
	}
}

// naturalThreshold is the value for group size of N under which N virtual
// stakers tickets would be expected to fall below if the tokens were optimally
// staked, and the tickets values were evenly distributed in the domain of the
// pseudorandom function.
//
// natural threshold =
// (group size * number of all possible ticket values) /
// (token supply / min stake)
func naturalThreshold(chainConfig *config.Chain) *big.Int {
	// (2^256)-1
	ticketsSpace := new(big.Int).Sub(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
		big.NewInt(1),
	)

	// 10^27
	tokenSupply := new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil)

	// groupSize * ( ticketsSpace / (tokenSupply / minimumStake) )
	return new(big.Int).Mul(
		big.NewInt(int64(chainConfig.GroupSize)),
		new(big.Int).Div(
			ticketsSpace,
			new(big.Int).Div(tokenSupply, chainConfig.MinimumStake),
		),
	)
}
