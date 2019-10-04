// Package groupselection implements the random beacon group selection protocol
// - an interactive, ticket-based method of selecting a candidate group from
// the set of all stakers given a pseudorandom seed value.
package groupselection

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ipfs/go-log"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/chain"
)

var logger = log.Logger("keep-groupselection")

// Result represents the result of group selection protocol. It contains the
// list of all stakers selected to the candidate group as well as the number of
// block at which the group selection protocol completed.
type Result struct {
	SelectedStakers        [][]byte
	GroupSelectionEndBlock uint64
}

// SubmitTickets attempts to generate and submit tickets for the staker to join
// a new candidate group.
func SubmitTickets(
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
			chainConfig.NaturalThreshold,
		)
	if err != nil {
		return err
	}

	logger.Info(
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
		startBlockHeight + chainConfig.TicketInitialSubmissionTimeout,
	)
	if err != nil {
		return err
	}

	reactiveSubmissionTimeout, err := blockCounter.BlockHeightWaiter(
		startBlockHeight + chainConfig.TicketReactiveSubmissionTimeout,
	)
	if err != nil {
		return err
	}

	quitTicketSubmission := make(chan struct{}, 1)

	go submitTickets(
		initialSubmissionTickets,
		relayChain,
		quitTicketSubmission,
	)

	for {
		select {
		case initialSubmissionEndBlockHeight := <-initialSubmissionTimeout:
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
				logger.Infof(
					"[%v] tickets submitted by candidates; "+
						"skipping reactive submission",
					ticketsCount,
				)

				quitTicketSubmission <- struct{}{}
				return nil
			}

			logger.Infof(
				"[%v] tickets submitted by candidates; "+
					"entering reactive submission",
				ticketsCount,
			)

			go submitTickets(
				reactiveSubmissionTickets,
				relayChain,
				quitTicketSubmission,
			)

		case reactiveSubmissionEndBlockHeight := <-reactiveSubmissionTimeout:
			logger.Infof(
				"reactive ticket submission ended at block [%v]",
				reactiveSubmissionEndBlockHeight,
			)

			quitTicketSubmission <- struct{}{}

			selectedParticipants, err := relayChain.GetSelectedParticipants()
			if err != nil {
				return fmt.Errorf(
					"could not fetch selected participants after submission timeout [%v]",
					err,
				)
			}

			selectedStakers := make([][]byte, len(selectedParticipants))
			for i, participant := range selectedParticipants {
				selectedStakers[i] = participant
				logger.Infof("new group member: [0x%v]", hex.EncodeToString(participant))
			}

			go onGroupSelected(&Result{
				SelectedStakers:        selectedStakers,
				GroupSelectionEndBlock: reactiveSubmissionEndBlockHeight,
			})

			return nil
		}
	}
}
