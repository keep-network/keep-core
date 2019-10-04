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
	signing chain.Signing,
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

	tickets := append(initialSubmissionTickets, reactiveSubmissionTickets...)

	submissionTimeout, err := blockCounter.BlockHeightWaiter(
		startBlockHeight + chainConfig.TicketReactiveSubmissionTimeout,
	)
	if err != nil {
		return err
	}

	var (
		errorChannel         = make(chan error, len(tickets))
		quitTicketSubmission = make(chan struct{}, 1)
	)

	go submitTickets(
		tickets,
		relayChain,
		quitTicketSubmission,
		errorChannel,
	)

	for {
		select {
		case err := <-errorChannel:
			logger.Errorf(
				"error during ticket submission: [%v]",
				err,
			)
		case submissionEndBlockHeight := <-submissionTimeout:
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
				GroupSelectionEndBlock: submissionEndBlockHeight,
			})

			return nil
		}
	}
}

// submitTickets submits tickets to the chain. It checks to see if the submission
// period is over in between ticket submits.
func submitTickets(
	tickets []*ticket,
	relayChain relaychain.GroupSelectionInterface,
	quit <-chan struct{},
	errCh chan<- error,
) {
	for _, ticket := range tickets {
		select {
		case <-quit:
			// Exit this loop when we get a signal from quit.
			return
		default:
			chainTicket, err := toChainTicket(ticket)
			if err != nil {
				errCh <- err
				continue
			}

			relayChain.SubmitTicket(chainTicket).OnFailure(
				func(err error) { errCh <- err },
			)
		}
	}
}

func toChainTicket(ticket *ticket) (*relaychain.Ticket, error) {
	return &relaychain.Ticket{
		Value: ticket.intValue(),
		Proof: &relaychain.TicketProof{
			StakerValue:        new(big.Int).SetBytes(ticket.proof.stakerValue),
			VirtualStakerIndex: ticket.proof.virtualStakerIndex,
		},
	}, nil
}
