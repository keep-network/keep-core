package relay

import (
	"encoding/hex"
	"fmt"
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
)

// SubmitTicketsForGroupSelection takes the previous beacon value and attempts to
// generate the appropriate number of tickets for the staker. After ticket
// generation begins an interactive process, where the staker submits tickets
// that fall under the natural threshold, while challenging tickets on chain
// that fail verification. Submission ends at the end of the submission period.
//
// See the group selection protocol specification for more information.
func SubmitTicketsForGroupSelection(
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	signing chain.Signing,
	chainConfig *config.Chain,
	staker chain.Staker,
	newEntry *big.Int,
	startBlockHeight uint64,
	onGroupSelected func(*groupselection.Result),
) error {
	availableStake, err := staker.Stake()
	if err != nil {
		return err
	}
	tickets, err :=
		groupselection.GenerateTickets(
			newEntry.Bytes(),
			staker.ID(),
			availableStake,
			chainConfig.MinimumStake,
		)
	if err != nil {
		return err
	}

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

	// submit all tickets
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

			go onGroupSelected(&groupselection.Result{
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
	tickets []*groupselection.Ticket,
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

func toChainTicket(ticket *groupselection.Ticket) (*relaychain.Ticket, error) {
	return &relaychain.Ticket{
		Value: ticket.Value.Int(),
		Proof: &relaychain.TicketProof{
			StakerValue:        new(big.Int).SetBytes(ticket.Proof.StakerValue),
			VirtualStakerIndex: ticket.Proof.VirtualStakerIndex,
		},
	}, nil
}
