package local

import (
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
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

type localChain struct {
	relayConfig *relayconfig.Chain

	groupRegistrationsMutex sync.Mutex
	groupRegistrations      map[string][96]byte

	groupRelayEntriesMutex sync.Mutex
	groupRelayEntries      map[string]*big.Int

	submittedResultsMutex sync.Mutex
	// Map of submitted DKG Results. Key is a RequestID of the specific DKG
	// execution.
	submittedResults map[*big.Int][]*relaychain.DKGResult

	handlerMutex                 sync.Mutex
	relayEntryHandlers           []func(entry *event.Entry)
	relayRequestHandlers         []func(request *event.Request)
	groupRegisteredHandlers      []func(key *event.GroupRegistration)
	dkgResultPublicationHandlers map[int]func(dkgResultPublication *event.DKGResultPublication)

	requestID   int64
	latestValue *big.Int

	simulatedHeight int64
	stakeMonitor    chain.StakeMonitor
	blockCounter    chain.BlockCounter

	tickets      []*groupselection.Ticket
	ticketsMutex sync.Mutex
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

func (c *localChain) SubmitTicket(ticket *groupselection.Ticket) *async.GroupTicketPromise {
	promise := &async.GroupTicketPromise{}

	c.ticketsMutex.Lock()
	defer c.ticketsMutex.Unlock()

	c.tickets = append(c.tickets, ticket)
	sort.SliceStable(c.tickets, func(i, j int) bool {
		return c.tickets[i].Value.Int().Cmp(c.tickets[j].Value.Int()) == -1
	})

	promise.Fulfill(ticket)

	return promise
}

func (c *localChain) SubmitChallenge(
	ticket *groupselection.TicketChallenge,
) *async.GroupTicketChallengePromise {
	promise := &async.GroupTicketChallengePromise{}
	promise.Fail(fmt.Errorf("function not implemented"))
	return promise
}

func (c *localChain) GetOrderedTickets() ([]*groupselection.Ticket, error) {
	c.ticketsMutex.Lock()
	defer c.ticketsMutex.Unlock()

	return c.tickets, nil
}

func (c *localChain) SubmitGroupPublicKey(
	requestID *big.Int,
	key [96]byte,
) *async.GroupRegistrationPromise {
	groupID := requestID.String()

	groupRegistrationPromise := &async.GroupRegistrationPromise{}
	registration := &event.GroupRegistration{
		GroupPublicKey:        key[:],
		RequestID:             requestID,
		ActivationBlockHeight: big.NewInt(c.simulatedHeight),
	}

	c.groupRegistrationsMutex.Lock()
	defer c.groupRegistrationsMutex.Unlock()
	if existing, exists := c.groupRegistrations[groupID]; exists {
		if existing != key {
			err := fmt.Errorf(
				"mismatched public key for [%s], submission failed; \n"+
					"[%v] vs [%v]",
				groupID,
				existing,
				key,
			)
			fmt.Fprintf(os.Stderr, err.Error())

			groupRegistrationPromise.Fail(err)
		} else {
			groupRegistrationPromise.Fulfill(registration)
		}

		return groupRegistrationPromise
	}
	c.groupRegistrations[groupID] = key

	groupRegistrationPromise.Fulfill(registration)

	c.handlerMutex.Lock()
	for _, handler := range c.groupRegisteredHandlers {
		go func(handler func(registration *event.GroupRegistration), registration *event.GroupRegistration) {
			handler(registration)
		}(handler, registration)
	}
	c.handlerMutex.Unlock()

	atomic.AddInt64(&c.simulatedHeight, 1)

	return groupRegistrationPromise
}

func (c *localChain) SubmitRelayEntry(entry *event.Entry) *async.RelayEntryPromise {
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

func (c *localChain) OnRelayEntryGenerated(handler func(entry *event.Entry)) {
	c.handlerMutex.Lock()
	c.relayEntryHandlers = append(
		c.relayEntryHandlers,
		handler,
	)
	c.handlerMutex.Unlock()
}

func (c *localChain) OnRelayEntryRequested(handler func(request *event.Request)) {
	c.handlerMutex.Lock()
	c.relayRequestHandlers = append(
		c.relayRequestHandlers,
		handler,
	)
	c.handlerMutex.Unlock()
}

func (c *localChain) OnGroupRegistered(handler func(key *event.GroupRegistration)) {
	c.handlerMutex.Lock()
	c.groupRegisteredHandlers = append(
		c.groupRegisteredHandlers,
		handler,
	)
	c.handlerMutex.Unlock()
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
		groupRegistrationsMutex:      sync.Mutex{},
		groupRelayEntries:            make(map[string]*big.Int),
		groupRegistrations:           make(map[string][96]byte),
		submittedResults:             make(map[*big.Int][]*relaychain.DKGResult),
		dkgResultPublicationHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
		blockCounter:                 bc,
		stakeMonitor:                 NewStakeMonitor(minimumStake),
		tickets:                      make([]*groupselection.Ticket, 0),
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

	return c.submittedResults[requestID] != nil, nil
}

// SubmitDKGResult submits the result to a chain.
func (c *localChain) SubmitDKGResult(
	requestID *big.Int, resultToPublish *relaychain.DKGResult,
) *async.DKGResultPublicationPromise {
	c.submittedResultsMutex.Lock()
	defer c.submittedResultsMutex.Unlock()

	dkgResultPublicationPromise := &async.DKGResultPublicationPromise{}

	for publishedRequestID, publishedResults := range c.submittedResults {
		if publishedRequestID.Cmp(requestID) == 0 {
			for _, publishedResult := range publishedResults {
				if publishedResult.Equals(resultToPublish) {
					dkgResultPublicationPromise.Fail(fmt.Errorf("result already submitted"))
					return dkgResultPublicationPromise
				}
			}
		}
	}

	c.submittedResults[requestID] = append(c.submittedResults[requestID], resultToPublish)

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
) (event.Subscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()
	c.dkgResultPublicationHandlers[handlerID] = handler

	return event.NewSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.dkgResultPublicationHandlers, handlerID)
	}), nil
}
