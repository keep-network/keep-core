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

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/subscription"
)

type localChain struct {
	relayConfig *relayconfig.Chain

	groupRegistrationsMutex sync.Mutex
	groupRegistrations      map[string][]byte

	groupRelayEntriesMutex sync.Mutex
	groupRelayEntries      map[string]*big.Int

	submittedResultsMutex sync.Mutex
	// Map of submitted DKG Results. Key is a RequestID of the specific DKG
	// execution.
	submittedResults map[string][]*relaychain.DKGResult

	handlerMutex                 sync.Mutex
	relayEntryHandlers           map[int]func(entry *event.Entry)
	relayRequestHandlers         map[int]func(request *event.Request)
	groupRegisteredHandlers      map[int]func(groupRegistration *event.GroupRegistration)
	dkgResultPublicationHandlers map[int]func(dkgResultPublication *event.DKGResultPublication)

	requestID   int64
	latestValue *big.Int

	simulatedHeight int64
	stakeMonitor    chain.StakeMonitor
	blockCounter    chain.BlockCounter

	stakerList []string

	// Track the submitted votes.
	submissionsMutex sync.Mutex
	submissions      map[string]*relaychain.DKGSubmissions
	voteHandler      []func(dkgResultVote *event.DKGResultVote)

	tickets      []*relaychain.Ticket
	ticketsMutex sync.Mutex
}

// GetDKGSubmissions returns the current set of submissions for the requestID.
func (c *localChain) GetDKGSubmissions(requestID *big.Int) *relaychain.DKGSubmissions {
	c.submissionsMutex.Lock()
	defer c.submissionsMutex.Unlock()
	return c.submissions[requestID.String()]
}

// DKGResultVote places a vote for dkgResultHash and causes OnDKGResultVote event to occurs.
func (c *localChain) DKGResultVote(
	requestID *big.Int,
	dkgResultHash []byte,
) *async.DKGResultVotePromise {
	dkgResultVotePromise := &async.DKGResultVotePromise{}

	c.submissionsMutex.Lock()
	defer c.submissionsMutex.Unlock()
	submissions, ok := c.submissions[requestID.String()]
	if !ok {
		err := dkgResultVotePromise.Fail(
			fmt.Errorf("no submissions for given request id"),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "promise fail failed [%v]\n", err)
		}
		return dkgResultVotePromise
	}

	for _, submission := range submissions.DKGSubmissions {
		if bytes.Equal(submission.DKGResult.Hash(), dkgResultHash) {
			submission.Votes++
			dkgResultVote := &event.DKGResultVote{
				RequestID: requestID,
			}
			c.handlerMutex.Lock()
			for _, handler := range c.voteHandler {
				go func(handler func(*event.DKGResultVote), dkgResultVote *event.DKGResultVote) {
					handler(dkgResultVote)
				}(handler, dkgResultVote)
			}
			c.handlerMutex.Unlock()

			err := dkgResultVotePromise.Fulfill(dkgResultVote)
			if err != nil {
				fmt.Fprintf(os.Stderr, "promise fulfill failed [%v]\n", err)
			}

			break
		}
	}

	return dkgResultVotePromise
}

// OnDKGResultVote sets up to call the passed handler function when a vote occurs.
func (c *localChain) OnDKGResultVote(handler func(dkgResultVote *event.DKGResultVote)) {
	c.handlerMutex.Lock()
	c.voteHandler = append(c.voteHandler, handler)
	c.handlerMutex.Unlock()
}

func (c *localChain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}

func (c *localChain) StakeMonitor() (chain.StakeMonitor, error) {
	return c.stakeMonitor, nil
}

func (c *localChain) GetConfig() (*relayconfig.Chain, error) {
	return c.relayConfig, nil
}

func (c *localChain) SubmitTicket(ticket *relaychain.Ticket) *async.GroupTicketPromise {
	promise := &async.GroupTicketPromise{}

	c.ticketsMutex.Lock()
	defer c.ticketsMutex.Unlock()

	c.tickets = append(c.tickets, ticket)
	sort.SliceStable(c.tickets, func(i, j int) bool {
		return c.tickets[i].Value.Cmp(c.tickets[j].Value) == -1
	})

	promise.Fulfill(&event.GroupTicketSubmission{
		TicketValue: ticket.Value,
	})

	return promise
}

func (c *localChain) GetSelectedParticipants() ([]relaychain.StakerAddress, error) {
	c.ticketsMutex.Lock()
	defer c.ticketsMutex.Unlock()

	selectTickets := func() []*relaychain.Ticket {
		if len(c.tickets) <= c.relayConfig.GroupSize {
			return c.tickets
		}

		selectedTickets := make([]*relaychain.Ticket, c.relayConfig.GroupSize)
		copy(selectedTickets, c.tickets)
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
	groupID := requestID.String()

	groupRegistrationPromise := &async.GroupRegistrationPromise{}
	groupRegistration := &event.GroupRegistration{
		GroupPublicKey:        groupPublicKey,
		RequestID:             requestID,
		ActivationBlockHeight: big.NewInt(c.simulatedHeight),
	}

	c.groupRegistrationsMutex.Lock()
	defer c.groupRegistrationsMutex.Unlock()
	if existing, exists := c.groupRegistrations[groupID]; exists {
		if bytes.Compare(existing, groupPublicKey) != 0 {
			err := fmt.Errorf(
				"mismatched public key for [%s], submission failed; \n"+
					"[%v] vs [%v]",
				groupID,
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
	c.groupRegistrations[groupID] = groupPublicKey

	groupRegistrationPromise.Fulfill(groupRegistration)

	c.handlerMutex.Lock()
	for _, handler := range c.groupRegisteredHandlers {
		go func(handler func(groupRegistration *event.GroupRegistration), registration *event.GroupRegistration) {
			handler(registration)
		}(handler, groupRegistration)
	}
	c.handlerMutex.Unlock()

	atomic.AddInt64(&c.simulatedHeight, 1)

	return groupRegistrationPromise
}

func (c *localChain) SubmitRelayEntry(entry *event.Entry) *async.RelayEntryPromise {
	c.ticketsMutex.Lock()
	c.tickets = make([]*relaychain.Ticket, 0)
	c.ticketsMutex.Unlock()

	relayEntryPromise := &async.RelayEntryPromise{}

	c.groupRelayEntriesMutex.Lock()
	defer c.groupRelayEntriesMutex.Unlock()

	existing, exists := c.groupRelayEntries[entry.GroupID.String()+entry.RequestID.String()]
	if exists {
		if existing.Cmp(entry.Value) != 0 {
			err := fmt.Errorf(
				"mismatched signature for [%v], submission failed; \n"+
					"[%v] vs [%v]\n",
				entry.GroupID,
				existing,
				entry.Value,
			)

			relayEntryPromise.Fail(err)
		} else {
			relayEntryPromise.Fulfill(entry)
		}

		return relayEntryPromise
	}
	c.groupRelayEntries[entry.GroupID.String()+entry.RequestID.String()] = entry.Value

	c.handlerMutex.Lock()
	for _, handler := range c.relayEntryHandlers {
		go func(handler func(entry *event.Entry), entry *event.Entry) {
			handler(entry)
		}(handler, entry)
	}
	c.handlerMutex.Unlock()

	c.latestValue = entry.Value
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

func (c *localChain) OnGroupRegistered(
	handler func(groupRegistration *event.GroupRegistration),
) (subscription.EventSubscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()

	c.groupRegisteredHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.groupRegisteredHandlers, handlerID)
	}), nil
}

func (c *localChain) ThresholdRelay() relaychain.Interface {
	return relaychain.Interface(c)
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
		groupRegistrationsMutex: sync.Mutex{},
		groupRelayEntries:       make(map[string]*big.Int),

		groupRegistrations: make(map[string][]byte),

		submittedResults:             make(map[string][]*relaychain.DKGResult),
		dkgResultPublicationHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
		blockCounter:                 bc,
		submissions:                  make(map[string]*relaychain.DKGSubmissions),

		relayEntryHandlers:      make(map[int]func(request *event.Entry)),
		relayRequestHandlers:    make(map[int]func(request *event.Request)),
		groupRegisteredHandlers: make(map[int]func(groupRegistration *event.GroupRegistration)),
		stakeMonitor:            NewStakeMonitor(minimumStake),
		tickets:                 make([]*relaychain.Ticket, 0),
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

// RequestRelayEntry simulates calling to start the random generation process.
func (c *localChain) RequestRelayEntry(
	blockReward, seed *big.Int,
) *async.RelayRequestPromise {
	promise := &async.RelayRequestPromise{}

	request := &event.Request{
		PreviousValue: c.latestValue,
		RequestID:     big.NewInt(c.requestID),
		Payment:       big.NewInt(1),
		BlockReward:   blockReward,
		Seed:          seed,
	}
	atomic.AddInt64(&c.simulatedHeight, 1)
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

// IsDKGResultPublished simulates check if the result was already submitted to a
// chain.
func (c *localChain) IsDKGResultPublished(requestID *big.Int) (bool, error) {
	c.submittedResultsMutex.Lock()
	defer c.submittedResultsMutex.Unlock()

	return c.submittedResults[requestID.String()] != nil, nil
}

// SubmitDKGResult submits the result to a chain.
func (c *localChain) SubmitDKGResult(
	requestID *big.Int,
	resultToPublish *relaychain.DKGResult,
) *async.DKGResultPublicationPromise {
	c.submittedResultsMutex.Lock()
	defer c.submittedResultsMutex.Unlock()

	dkgResultPublicationPromise := &async.DKGResultPublicationPromise{}

	for publishedRequestID, publishedResults := range c.submittedResults {
		if publishedRequestID == requestID.String() {
			for _, publishedResult := range publishedResults {
				if publishedResult.Equals(resultToPublish) {
					dkgResultPublicationPromise.Fail(fmt.Errorf("result already submitted"))
					return dkgResultPublicationPromise
				}
			}
		}
	}

	c.submittedResults[requestID.String()] = append(c.submittedResults[requestID.String()], resultToPublish)

	c.submissionsMutex.Lock()
	if c.submissions == nil {
		c.submissions = make(map[string]*relaychain.DKGSubmissions)
	}
	if _, ok := c.submissions[requestID.String()]; !ok {
		c.submissions[requestID.String()] = &relaychain.DKGSubmissions{
			DKGSubmissions: []*relaychain.DKGSubmission{
				{
					DKGResult: resultToPublish,
					Votes:     1,
				},
			},
		}
	}
	c.submissionsMutex.Unlock()

	dkgResultPublicationEvent := &event.DKGResultPublication{
		RequestID:      requestID,
		GroupPublicKey: resultToPublish.GroupPublicKey[:],
	}

	c.handlerMutex.Lock()
	for _, handler := range c.dkgResultPublicationHandlers {
		go func(handler func(*event.DKGResultPublication), dkgResultPublication *event.DKGResultPublication) {
			handler(dkgResultPublicationEvent)
		}(handler, dkgResultPublicationEvent)
	}
	c.handlerMutex.Unlock()

	err := dkgResultPublicationPromise.Fulfill(dkgResultPublicationEvent)
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
	c.dkgResultPublicationHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.dkgResultPublicationHandlers, handlerID)
	}), nil
}
