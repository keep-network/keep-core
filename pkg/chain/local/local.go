package local

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/subscription"
)

type localChain struct {
	relayConfig *relayconfig.Chain

	relayEntriesMutex sync.Mutex
	relayEntries      map[string]*big.Int

	groupTicketsMutex       sync.Mutex
	groupTickets            []*relaychain.Ticket
	groupRegistrationsMutex sync.Mutex
	groupRegistrations      map[string][]byte

	dkgSubmittedResultsMutex        sync.Mutex
	dkgSubmittedResults             map[string][]*relaychain.DKGResult
	dkgResultsVotesMutex            sync.Mutex
	dkgResultsVotes                 map[string]relaychain.DKGResultsVotes
	dkgAlreadySubmittedOrVotedMutex sync.Mutex
	dkgAlreadySubmittedOrVoted      map[string][]int

	handlerMutex                sync.Mutex
	relayRequestHandlers        map[int]func(request *event.Request)
	relayEntryHandlers          map[int]func(entry *event.Entry)
	groupRegistrationHandlers   map[int]func(groupRegistration *event.GroupRegistration)
	dkgResultSubmissionHandlers map[int]func(dkgResultPublication *event.DKGResultPublication)
	dkgResultVoteHandler        map[int]func(dkgResultVote *event.DKGResultVote)

	requestID     int64
	previousValue *big.Int

	simulatedHeight uint64

	stakeMonitor chain.StakeMonitor
	blockCounter chain.BlockCounter
}

func (c *localChain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}

func (c *localChain) StakeMonitor() (chain.StakeMonitor, error) {
	return c.stakeMonitor, nil
}

func (c *localChain) ThresholdRelay() relaychain.Interface {
	return relaychain.Interface(c)
}

// RequestRelayEntry simulates calling to start the random generation process.
func (c *localChain) RequestRelayEntry(
	blockReward, seed *big.Int,
) *async.RelayRequestPromise {
	promise := &async.RelayRequestPromise{}

	request := &event.Request{
		PreviousValue: c.previousValue,
		RequestID:     big.NewInt(c.requestID),
		Payment:       big.NewInt(1),
		BlockReward:   blockReward,
		Seed:          seed,
	}
	atomic.AddUint64(&c.simulatedHeight, 1)
	atomic.AddInt64(&c.requestID, 1)

	c.handlerMutex.Lock()
	for _, handler := range c.relayRequestHandlers {
		go func(handler func(*event.Request), request *event.Request) {
			handler(request)
		}(handler, request)
	}
	c.handlerMutex.Unlock()

	promise.Fulfill(request)

	return promise
}

func (c *localChain) SubmitRelayEntry(entry *event.Entry) *async.RelayEntryPromise {
	// Cleans up submitted tickets after group selection execution.
	c.groupTicketsMutex.Lock()
	c.groupTickets = make([]*relaychain.Ticket, 0)
	c.groupTicketsMutex.Unlock()

	relayEntryPromise := &async.RelayEntryPromise{}

	c.relayEntriesMutex.Lock()
	defer c.relayEntriesMutex.Unlock()

	existing, exists := c.relayEntries[string(entry.GroupPubKey)+entry.RequestID.String()]
	if exists {
		if existing.Cmp(entry.Value) != 0 {
			err := fmt.Errorf(
				"mismatched signature for [%v], submission failed; \n"+
					"[%v] vs [%v]\n",
				entry.GroupPubKey,
				existing,
				entry.Value,
			)

			relayEntryPromise.Fail(err)
		} else {
			relayEntryPromise.Fulfill(entry)
		}

		return relayEntryPromise
	}
	c.relayEntries[string(entry.GroupPubKey)+entry.RequestID.String()] = entry.Value

	c.handlerMutex.Lock()
	for _, handler := range c.relayEntryHandlers {
		go func(handler func(entry *event.Entry), entry *event.Entry) {
			handler(entry)
		}(handler, entry)
	}
	c.handlerMutex.Unlock()

	c.previousValue = entry.Value
	relayEntryPromise.Fulfill(entry)

	return relayEntryPromise
}

func (c *localChain) OnRelayEntryGenerated(
	handler func(entry *event.Entry),
) (subscription.EventSubscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()
	c.relayEntryHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.relayEntryHandlers, handlerID)
	}), nil
}

func (c *localChain) OnRelayEntryRequested(
	handler func(request *event.Request),
) (subscription.EventSubscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()
	c.relayRequestHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.relayRequestHandlers, handlerID)
	}), nil
}

func (c *localChain) SubmitTicket(ticket *relaychain.Ticket) *async.GroupTicketPromise {
	promise := &async.GroupTicketPromise{}

	c.groupTicketsMutex.Lock()
	defer c.groupTicketsMutex.Unlock()

	c.groupTickets = append(c.groupTickets, ticket)
	sort.SliceStable(c.groupTickets, func(i, j int) bool {
		return c.groupTickets[i].Value.Cmp(c.groupTickets[j].Value) == -1
	})

	promise.Fulfill(&event.GroupTicketSubmission{
		TicketValue: ticket.Value,
		BlockNumber: c.simulatedHeight,
	})

	return promise
}

func (c *localChain) GetSelectedParticipants() ([]relaychain.StakerAddress, error) {
	c.groupTicketsMutex.Lock()
	defer c.groupTicketsMutex.Unlock()

	selectTickets := func() []*relaychain.Ticket {
		if len(c.groupTickets) <= c.relayConfig.GroupSize {
			return c.groupTickets
		}

		selectedTickets := make([]*relaychain.Ticket, c.relayConfig.GroupSize)
		copy(selectedTickets, c.groupTickets)
		return selectedTickets
	}

	selectedTickets := selectTickets()

	selectedParticipants := make([]relaychain.StakerAddress, len(selectedTickets))
	for i, ticket := range selectedTickets {
		selectedParticipants[i] = ticket.Proof.StakerValue.Bytes()
	}

	return selectedParticipants, nil
}

func (c *localChain) SubmitGroupPublicKey(
	requestID *big.Int,
	groupPublicKey []byte,
) *async.GroupRegistrationPromise {
	groupPubKey := requestID.String()

	groupRegistrationPromise := &async.GroupRegistrationPromise{}
	groupRegistration := &event.GroupRegistration{
		GroupPublicKey: groupPublicKey,
		RequestID:      requestID,
		BlockNumber:    c.simulatedHeight,
	}

	c.groupRegistrationsMutex.Lock()
	defer c.groupRegistrationsMutex.Unlock()
	if existing, exists := c.groupRegistrations[groupPubKey]; exists {
		if bytes.Compare(existing, groupPublicKey) != 0 {
			err := fmt.Errorf(
				"mismatched public key for [%s], submission failed; \n"+
					"[%v] vs [%v]",
				groupPubKey,
				existing,
				groupPublicKey,
			)
			fmt.Fprintf(os.Stderr, err.Error())

			groupRegistrationPromise.Fail(err)
		} else {
			groupRegistrationPromise.Fulfill(groupRegistration)
		}

		return groupRegistrationPromise
	}
	c.groupRegistrations[groupPubKey] = groupPublicKey

	groupRegistrationPromise.Fulfill(groupRegistration)

	c.handlerMutex.Lock()
	for _, handler := range c.groupRegistrationHandlers {
		go func(handler func(groupRegistration *event.GroupRegistration), registration *event.GroupRegistration) {
			handler(registration)
		}(handler, groupRegistration)
	}
	c.handlerMutex.Unlock()

	atomic.AddUint64(&c.simulatedHeight, 1)

	return groupRegistrationPromise
}

func (c *localChain) OnGroupRegistered(
	handler func(groupRegistration *event.GroupRegistration),
) (subscription.EventSubscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()

	c.groupRegistrationHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.groupRegistrationHandlers, handlerID)
	}), nil
}

// SubmitDKGResult submits the result to a chain.
func (c *localChain) SubmitDKGResult(
	requestID *big.Int,
	resultToPublish *relaychain.DKGResult,
) *async.DKGResultPublicationPromise {
	c.dkgSubmittedResultsMutex.Lock()
	defer c.dkgSubmittedResultsMutex.Unlock()

	c.dkgResultsVotesMutex.Lock()
	defer c.dkgResultsVotesMutex.Unlock()

	dkgResultPublicationPromise := &async.DKGResultPublicationPromise{}

	// Submit DKG result.
	dkgResults, ok := c.dkgSubmittedResults[requestID.String()]
	// Initialize map entry if it is a first submission for given request ID.
	if !ok {
		dkgResults = []*relaychain.DKGResult{}
	}

	// If the result is already submitted for the given request ID.
	for _, dkgResult := range dkgResults {
		if dkgResult.Equals(resultToPublish) {
			dkgResultPublicationPromise.Fail(fmt.Errorf("result already submitted"))
			return dkgResultPublicationPromise
		}
	}

	c.dkgSubmittedResults[requestID.String()] = append(dkgResults, resultToPublish)

	// Vote on DKG result.
	dkgResultsVotes, ok := c.dkgResultsVotes[requestID.String()]
	// Initialize map entry if it is a first submission for given request ID.
	if !ok {
		c.dkgResultsVotes[requestID.String()] = relaychain.DKGResultsVotes{}
	}

	resultToPublishHash, err := c.CalculateDKGResultHash(resultToPublish)
	if err != nil {
		dkgResultPublicationPromise.Fail(fmt.Errorf("hash calculation failed"))
		return dkgResultPublicationPromise
	}

	// If the result is already registered in votes map for the given request ID,
	// it shouldn't be possible to register it again. One should use OnDKGResultVote
	// function to submit a vote.
	if _, ok := dkgResultsVotes[resultToPublishHash]; ok {
		err := dkgResultPublicationPromise.Fail(
			fmt.Errorf("result already registered in votes map, use vote function to submit a vote"),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "promise fail failed [%v]\n", err)
		}
		return dkgResultPublicationPromise
	}

	// Register initial vote for newly submitted result.
	c.dkgResultsVotes[requestID.String()][resultToPublishHash] = 1

	// Emit an event on DKG result submission.
	currentBlock, err := c.blockCounter.CurrentBlock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get current block [%v]\n", err)
	}

	dkgResultSubmissionEvent := &event.DKGResultPublication{
		RequestID:      requestID,
		GroupPublicKey: resultToPublish.GroupPublicKey,
		BlockNumber:    uint64(currentBlock),
	}

	c.handlerMutex.Lock()
	for _, handler := range c.dkgResultSubmissionHandlers {
		go handler(dkgResultSubmissionEvent)
	}
	c.handlerMutex.Unlock()

	// Fulfill the promise.
	err = dkgResultPublicationPromise.Fulfill(dkgResultSubmissionEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "promise fulfill failed [%v].\n", err)
	}

	return dkgResultPublicationPromise
}

func (c *localChain) OnDKGResultPublished(
	handler func(dkgResultPublication *event.DKGResultPublication),
) (subscription.EventSubscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()
	c.dkgResultSubmissionHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.dkgResultSubmissionHandlers, handlerID)
	}), nil
}

// IsDKGResultPublished simulates check if the result was already submitted to a
// chain.
func (c *localChain) IsDKGResultPublished(requestID *big.Int) (bool, error) {
	c.dkgSubmittedResultsMutex.Lock()
	defer c.dkgSubmittedResultsMutex.Unlock()

	if submissions, ok := c.dkgSubmittedResults[requestID.String()]; ok {
		return len(submissions) > 0, nil
	}

	return false, nil
}

// CalculateDKGResultHash calculates a 256-bit hash of the DKG result.
func (c *localChain) CalculateDKGResultHash(
	dkgResult *relaychain.DKGResult,
) (relaychain.DKGResultHash, error) {
	encodedDKGResult := fmt.Sprint(dkgResult)
	dkgResultHash := relaychain.DKGResultHash(
		sha3.Sum256([]byte(encodedDKGResult)),
	)

	return dkgResultHash, nil
}

// GetDKGResultsVotes returns a map containing number of votes for each DKG
// result hash registered under specific request ID.
func (c *localChain) GetDKGResultsVotes(
	requestID *big.Int,
) relaychain.DKGResultsVotes {
	c.dkgResultsVotesMutex.Lock()
	defer c.dkgResultsVotesMutex.Unlock()

	if submissions, ok := c.dkgResultsVotes[requestID.String()]; ok {
		return submissions
	}

	return relaychain.DKGResultsVotes{}
}

// VoteOnDKGResult registers a vote for the DKG result hash.
func (c *localChain) VoteOnDKGResult(
	requestID *big.Int,
	memberIndex int,
	dkgResultHash relaychain.DKGResultHash,
) *async.DKGResultVotePromise {
	c.dkgResultsVotesMutex.Lock()
	defer c.dkgResultsVotesMutex.Unlock()

	c.dkgAlreadySubmittedOrVotedMutex.Lock()
	defer c.dkgAlreadySubmittedOrVotedMutex.Unlock()

	dkgResultVotePromise := &async.DKGResultVotePromise{}

	// Member cannot vote if already submitted DKG result or voted before.
	alreadySubmittedOrVoted, ok := c.dkgAlreadySubmittedOrVoted[requestID.String()]
	if !ok {
		err := dkgResultVotePromise.Fail(
			fmt.Errorf("no registered submissions or votes for given request id"),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "promise fail failed [%v]\n", err)
		}
		return dkgResultVotePromise
	}
	for _, index := range alreadySubmittedOrVoted {
		if memberIndex == index {
			err := dkgResultVotePromise.Fail(
				fmt.Errorf("this member already submitted or voted on DKG result for given request id"),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "promise fail failed [%v]\n", err)
			}
			return dkgResultVotePromise
		}
	}

	dkgResultsVotes, ok := c.dkgResultsVotes[requestID.String()]
	// If there are no result submissions registered in results votes map.
	if !ok || len(dkgResultsVotes) == 0 {
		err := dkgResultVotePromise.Fail(
			fmt.Errorf("no registered dkg results votes for given request id"),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "promise fail failed [%v]\n", err)
		}
		return dkgResultVotePromise
	}

	// It's not possible to vote on a dkg result which submission has not been
	// registered yet.
	if _, ok = dkgResultsVotes[dkgResultHash]; !ok {
		err := dkgResultVotePromise.Fail(
			fmt.Errorf("result hash is not registered in dkg results votes map"),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "promise fail failed [%v]\n", err)
		}
		return dkgResultVotePromise
	}

	// Register vote for the DKG result hash.
	dkgResultsVotes[dkgResultHash]++

	// Register that member already voted.
	c.dkgAlreadySubmittedOrVoted[requestID.String()] = append(
		alreadySubmittedOrVoted,
		memberIndex,
	)

	// Emit an event on DKG result vote.
	currentBlock, err := c.blockCounter.CurrentBlock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get current block [%v]\n", err)
	}

	dkgResultVote := &event.DKGResultVote{
		RequestID:     requestID,
		MemberIndex:   memberIndex,
		DKGResultHash: dkgResultHash,
		BlockNumber:   uint64(currentBlock),
	}

	c.handlerMutex.Lock()
	for _, handler := range c.dkgResultVoteHandler {
		go handler(dkgResultVote)
	}
	c.handlerMutex.Unlock()

	// Fulfill the promise.
	err = dkgResultVotePromise.Fulfill(dkgResultVote)
	if err != nil {
		fmt.Fprintf(os.Stderr, "promise fulfill failed [%v]\n", err)
	}

	return dkgResultVotePromise
}

// OnDKGResultVote registers a callback that is invoked when an on-chain
// notification of a new, valid vote is seen.
func (c *localChain) OnDKGResultVote(
	handler func(dkgResultVote *event.DKGResultVote),
) (subscription.EventSubscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()
	c.dkgResultVoteHandler[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.dkgResultVoteHandler, handlerID)
	}), nil
}

func (c *localChain) GetConfig() (*relayconfig.Chain, error) {
	return c.relayConfig, nil
}

// Connect initializes a local stub implementation of the chain interfaces
// for testing.
func Connect(groupSize int, threshold int, minimumStake *big.Int) chain.Handle {
	bc, _ := blockCounter()

	tokenSupply, naturalThreshold := calculateGroupSelectionParameters(
		groupSize,
		minimumStake,
	)

	return &localChain{
		relayConfig: &relayconfig.Chain{
			GroupSize:                       groupSize,
			Threshold:                       threshold,
			TicketInitialSubmissionTimeout:  2,
			TicketReactiveSubmissionTimeout: 3,
			TicketChallengeTimeout:          4,
			MinimumStake:                    minimumStake,
			TokenSupply:                     tokenSupply,
			NaturalThreshold:                naturalThreshold,
		},

		relayEntries: make(map[string]*big.Int),

		groupTickets:       make([]*relaychain.Ticket, 0),
		groupRegistrations: make(map[string][]byte),

		dkgSubmittedResults:        make(map[string][]*relaychain.DKGResult),
		dkgResultsVotes:            make(map[string]relaychain.DKGResultsVotes),
		dkgAlreadySubmittedOrVoted: make(map[string][]int),

		relayRequestHandlers:        make(map[int]func(request *event.Request)),
		relayEntryHandlers:          make(map[int]func(request *event.Entry)),
		groupRegistrationHandlers:   make(map[int]func(groupRegistration *event.GroupRegistration)),
		dkgResultSubmissionHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
		dkgResultVoteHandler:        make(map[int]func(dkgResultVote *event.DKGResultVote)),

		stakeMonitor: NewStakeMonitor(minimumStake),
		blockCounter: bc,
	}
}

func calculateGroupSelectionParameters(groupSize int, minimumStake *big.Int) (
	tokenSupply *big.Int,
	naturalThreshold *big.Int,
) {
	// (2^256)-1
	ticketsSpace := new(big.Int).Sub(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
		big.NewInt(1),
	)

	// 10^9
	tokenSupply = new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)

	// groupSize * ( ticketsSpace / (tokenSupply / minimumStake) )
	naturalThreshold = new(big.Int).Mul(
		big.NewInt(int64(groupSize)),
		new(big.Int).Div(
			ticketsSpace,
			new(big.Int).Div(tokenSupply, minimumStake),
		),
	)

	return tokenSupply, naturalThreshold
}
