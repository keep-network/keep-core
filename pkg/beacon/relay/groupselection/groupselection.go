// Package groupselection implements the random beacon group selection protocol
// - an interactive, ticket-based method of selecting a candidate group from
// the set of all stakers given a pseudorandom seed value.
package groupselection

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ipfs/go-log"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/chain"
)

var logger = log.Logger("keep-groupselection")

// Recommended parameters all clients should use to minimize their expenses.
// It is not a must to obey but it is nice and polite. And being nice to others
// helps in reducing own costs because other clients should respect the same
// protocol.
const (
	// The number of blocks one round takes.
	roundDuration = uint64(6)

	// The delay in blocks after all rounds complete to ensure all transactions
	// are mined
	miningLag = uint64(12)
)

// Result represents the result of group selection protocol. It contains the
// list of all stakers selected to the candidate group as well as the number of
// block at which the group selection protocol completed.
type Result struct {
	SelectedStakers        []relaychain.StakerAddress
	GroupSelectionEndBlock uint64
}

// CandidateToNewGroup attempts to generate and submit tickets for the
// staker to join a new group.
//
// To minimize the submitter's cost by minimizing the number of redundant
// tickets that are not selected into the group, tickets are submitted in
// N rounds, each round taking 6 blocks.
// As the basic principle, the number of leading zeros in the ticket
// value is subtracted from the number of rounds to determine the round
// the ticket should be submitted in:
// - in round 0, tickets with N or more leading zeros are submitted
// - in round 1, tickets with N-1 or more leading zeros are submitted
// (...)
// - in round N, tickets with no leading zeros are submitted.
//
// In each round, group member candidate needs to monitor tickets
// submitted by other candidates and compare them against tickets of
// the candidate not yet submitted to determine if continuing with
// ticket submission still makes sense.
//
// After the last round, there is a 12 blocks mining lag allowing all
// outstanding ticket submissions to have a higher chance of being
// mined before the deadline.
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

	minimumStake, err := relayChain.MinimumStake()
	if err != nil {
		return err
	}

	tickets, err := generateTickets(
		newEntry.Bytes(),
		staker.Address(),
		availableStake,
		minimumStake,
	)
	if err != nil {
		return err
	}

	logger.Infof("starting ticket submission with [%v] tickets", len(tickets))

	err = submitTickets(
		tickets,
		relayChain,
		blockCounter,
		chainConfig,
		startBlockHeight,
	)
	if err != nil {
		logger.Errorf("ticket submission terminated with error: [%v]", err)
	}

	// Wait till the end of the ticket submission in case submitTickets failed
	// in the middle and there is still a chance we qualified to a group.
	ticketSubmissionTimeoutChannel, err := blockCounter.BlockHeightWaiter(
		startBlockHeight + chainConfig.TicketSubmissionTimeout,
	)
	if err != nil {
		return err
	}

	ticketSubmissionEndBlockHeight := <-ticketSubmissionTimeoutChannel

	logger.Infof(
		"ticket submission ended at block [%v]",
		ticketSubmissionEndBlockHeight,
	)

	selectedStakers, err := relayChain.GetSelectedParticipants()
	if err != nil {
		return fmt.Errorf(
			"could not fetch selected participants "+
				"after submission timeout [%v]",
			err,
		)
	}

	go onGroupSelected(&Result{
		SelectedStakers:        selectedStakers,
		GroupSelectionEndBlock: ticketSubmissionEndBlockHeight,
	})

	return nil
}

func submitTickets(
	tickets []*ticket,
	relayChain relaychain.GroupSelectionInterface,
	blockCounter chain.BlockCounter,
	chainConfig *config.Chain,
	startBlockHeight uint64,
) error {
	rounds, err := calculateRoundsCount(chainConfig.TicketSubmissionTimeout)
	if err != nil {
		return err
	}

	for roundIndex := uint64(0); roundIndex <= rounds; roundIndex++ {
		roundStartDelay := roundIndex * roundDuration
		roundStartBlock := startBlockHeight + roundStartDelay
		roundLeadingZeros := rounds - roundIndex

		logger.Infof(
			"ticket submission round [%v] will start at "+
				"block [%v] and cover tickets with [%v] leading zeros",
			roundIndex,
			roundStartBlock,
			roundLeadingZeros,
		)

		err := blockCounter.WaitForBlockHeight(roundStartBlock)
		if err != nil {
			return err
		}

		candidateTickets, err := roundCandidateTickets(
			relayChain,
			tickets,
			roundIndex,
			roundLeadingZeros,
			chainConfig.GroupSize,
		)
		if err != nil {
			return err
		}

		logger.Infof(
			"ticket submission round [%v] submitting "+
				"[%v] tickets",
			roundIndex,
			len(candidateTickets),
		)

		submitTicketsOnChain(candidateTickets, relayChain)
	}

	return nil
}

// calculateRoundsCount takes the on-chain ticket submission timeout
// and calculates the number of rounds for ticket submission. If it is not
// possible to use the recommended round duration and mining lag because the
// supplied timeout is too short, function returns an error.
func calculateRoundsCount(submissionTimeout uint64) (uint64, error) {
	if submissionTimeout-miningLag <= roundDuration {
		return 0, fmt.Errorf("submission timeout is too short")
	}

	return (submissionTimeout - miningLag) / roundDuration, nil
}

// roundCandidateTickets returns tickets which should be submitted in
// the given ticket submission round.
//
// Bear in mind that member tickets slice should be sorted in ascending
// order by their value.
func roundCandidateTickets(
	relayChain relaychain.GroupSelectionInterface,
	memberTickets []*ticket,
	roundIndex uint64,
	roundLeadingZeros uint64,
	groupSize int,
) ([]*ticket, error) {
	// Get unsorted submitted tickets from the chain.
	// This slice will be also filled by candidate tickets values in order to
	// compare subsequent member ticket values against all submitted tickets
	// so far and determine an optimal number of candidate tickets.
	submittedTickets, err := relayChain.GetSubmittedTickets()
	if err != nil {
		return nil, fmt.Errorf(
			"could not get submitted tickets: [%v]",
			err,
		)
	}

	candidateTickets := make([]*ticket, 0)

	for _, candidateTicket := range memberTickets {
		candidateTicketLeadingZeros := uint64(
			candidateTicket.leadingZeros(),
		)

		// Check if the given candidate ticket should be proceeded in
		// the current round.
		if roundIndex == 0 {
			if candidateTicketLeadingZeros < roundLeadingZeros {
				continue
			}
		} else {
			if candidateTicketLeadingZeros != roundLeadingZeros {
				continue
			}
		}

		// Sort submitted tickets slice in ascending order.
		sort.SliceStable(
			submittedTickets,
			func(i, j int) bool {
				return submittedTickets[i] < submittedTickets[j]
			},
		)

		shouldBeSubmitted := false
		candidateTicketValue := candidateTicket.intValue().Uint64()

		if len(submittedTickets) < groupSize {
			// If the submitted tickets count is less than the group
			// size the candidate ticket can be added unconditionally.
			submittedTickets = append(
				submittedTickets,
				candidateTicketValue,
			)
			shouldBeSubmitted = true
		} else {
			// If the submitted tickets count is equal to the group
			// size the candidate ticket can be added only if
			// it is smaller than the highest submitted ticket.
			// Note that, maximum length of submitted tickets slice
			// will be exceeded and will be trimmed in next
			// iteration.
			highestSubmittedTicket := submittedTickets[len(submittedTickets)-1]
			if candidateTicketValue < highestSubmittedTicket {
				submittedTickets[len(submittedTickets)-1] = candidateTicketValue
				shouldBeSubmitted = true
			}
		}

		// If current candidate ticket should not be submitted,
		// there is no sense to continue with next candidate tickets
		// because they will have higher value than the current one.
		if !shouldBeSubmitted {
			break
		}

		candidateTickets = append(
			candidateTickets,
			candidateTicket,
		)
	}

	return candidateTickets, nil
}
