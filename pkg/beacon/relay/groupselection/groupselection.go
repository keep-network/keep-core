package groupselection

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"

	"github.com/ipfs/go-log"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/chain"
)

var logger = log.Logger("keep-groupselection")

// SubmitTickets takes the previous beacon value and attempts to
// generate the appropriate number of tickets for the staker. After ticket
// generation begins an interactive process, where the staker submits tickets
// that fall under the natural threshold, while challenging tickets on chain
// that fail verification. Submission ends at the end of the submission period.
//
// See the group selection protocol specification for more information.
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
	tickets, err :=
		generateTickets(
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
		Value: ticket.Value.int(),
		Proof: &relaychain.TicketProof{
			StakerValue:        new(big.Int).SetBytes(ticket.Proof.StakerValue),
			VirtualStakerIndex: ticket.Proof.VirtualStakerIndex,
		},
	}, nil
}

// generateTickets generates a set of tickets for the given staker and relay
// entry value given the specified minimum stake. Returns the resulting
// tickets in sorted order, or an error if there were issues computing the
// tickets.
func generateTickets(
	beaconValue []byte, // V_i
	stakerValue []byte, // Q_j
	availableStake *big.Int, // S_j
	minimumStake *big.Int,
) ([]*ticket, error) {
	stakingWeight := (&big.Int{}).Quo(availableStake, minimumStake) // W_j

	tickets := make(tickets, 0)
	for virtualStaker := int64(1); virtualStaker <= stakingWeight.Int64(); virtualStaker++ {
		ticket, err := newTicket(beaconValue, stakerValue, big.NewInt(virtualStaker)) // prf
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	sort.Stable(tickets)

	return tickets, nil
}

// Result represents ordered, selected tickets from those submitted to the chain.
type Result struct {
	SelectedStakers        [][]byte
	GroupSelectionEndBlock uint64
}
