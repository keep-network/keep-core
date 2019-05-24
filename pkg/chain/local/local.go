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
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/subscription"
)

var seedGroupPublicKey = []byte("seed to group public key")
var seedRelayEntry = big.NewInt(123456789)
var groupActiveTime = uint64(10)
var relayRequestTimeout = uint64(8)

// Chain is an extention of chain.Handle interface which exposes
// additional functions useful for testing.
type Chain interface {
	chain.Handle

	// GetDKGResult returns DKG result submitted to the chain for the given
	// request ID.
	GetDKGResult(requestID *big.Int) *relaychain.DKGResult
}

type localGroup struct {
	groupPublicKey          []byte
	registrationBlockHeight uint64
}

type localChain struct {
	relayConfig *relayconfig.Chain

	groupRelayEntriesMutex sync.Mutex
	groupRelayEntries      map[string]*big.Int

	submittedResultsMutex sync.Mutex
	// Map of submitted DKG Results. Key is a RequestID of the specific DKG
	// execution.
	submittedResults map[string]*relaychain.DKGResult

	groups []localGroup

	handlerMutex                  sync.Mutex
	relayEntryHandlers            map[int]func(entry *event.Entry)
	relayRequestHandlers          map[int]func(request *event.Request)
	groupSelectionStartedHandlers map[int]func(groupSelectionStart *event.GroupSelectionStart)
	groupRegisteredHandlers       map[int]func(groupRegistration *event.GroupRegistration)
	resultSubmissionHandlers      map[int]func(submission *event.DKGResultSubmission)

	requestID   int64
	latestValue *big.Int

	simulatedHeight uint64
	stakeMonitor    chain.StakeMonitor
	blockCounter    chain.BlockCounter

	tickets      []*relaychain.Ticket
	ticketsMutex sync.Mutex
}

func (c *localChain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}

func (c *localChain) StakeMonitor() (chain.StakeMonitor, error) {
	return c.stakeMonitor, nil
}

func (c *localChain) GetKeys() (*operator.PrivateKey, *operator.PublicKey) {
	privateKey, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		panic(err)
	}
	return privateKey, publicKey
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
		BlockNumber: c.simulatedHeight,
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

func (c *localChain) SubmitRelayEntry(entry *event.Entry) *async.RelayEntryPromise {
	c.ticketsMutex.Lock()
	c.tickets = make([]*relaychain.Ticket, 0)
	c.ticketsMutex.Unlock()

	relayEntryPromise := &async.RelayEntryPromise{}

	c.groupRelayEntriesMutex.Lock()
	defer c.groupRelayEntriesMutex.Unlock()

	existing, exists := c.groupRelayEntries[string(entry.GroupPubKey)+entry.RequestID.String()]
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
	c.groupRelayEntries[string(entry.GroupPubKey)+entry.RequestID.String()] = entry.Value

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

func (c *localChain) OnGroupSelectionStarted(
	handler func(entry *event.GroupSelectionStart),
) (subscription.EventSubscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()
	c.groupSelectionStartedHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.groupSelectionStartedHandlers, handlerID)
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
func Connect(groupSize int, threshold int, minimumStake *big.Int) Chain {
	bc, _ := blockCounter()

	tokenSupply, naturalThreshold := calculateGroupSelectionParameters(
		groupSize,
		minimumStake,
	)

	currentBlock, _ := bc.CurrentBlock()
	group := localGroup{
		groupPublicKey:          seedGroupPublicKey,
		registrationBlockHeight: currentBlock,
	}

	return &localChain{
		relayConfig: &relayconfig.Chain{
			GroupSize:                       groupSize,
			Threshold:                       threshold,
			TicketInitialSubmissionTimeout:  2,
			TicketReactiveSubmissionTimeout: 3,
			TicketChallengeTimeout:          4,
			ResultPublicationBlockStep:      3,
			MinimumStake:                    minimumStake,
			TokenSupply:                     tokenSupply,
			NaturalThreshold:                naturalThreshold,
		},
		groupRelayEntries:        make(map[string]*big.Int),
		submittedResults:         make(map[string]*relaychain.DKGResult),
		relayEntryHandlers:       make(map[int]func(request *event.Entry)),
		relayRequestHandlers:     make(map[int]func(request *event.Request)),
		groupRegisteredHandlers:  make(map[int]func(groupRegistration *event.GroupRegistration)),
		resultSubmissionHandlers: make(map[int]func(submission *event.DKGResultSubmission)),
		blockCounter:             bc,
		stakeMonitor:             NewStakeMonitor(minimumStake),
		tickets:                  make([]*relaychain.Ticket, 0),
		latestValue:              seedRelayEntry,
		groups:                   []localGroup{group},
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

func selectGroup(entry *big.Int, numberOfGroups int) int {
	if numberOfGroups == 0 {
		return 0
	}

	return int(new(big.Int).Mod(entry, big.NewInt(int64(numberOfGroups))).Int64())
}

// RequestRelayEntry simulates calling to start the random generation process.
func (c *localChain) RequestRelayEntry(seed *big.Int) *async.RelayRequestPromise {
	promise := &async.RelayRequestPromise{}

	selectedIdx := selectGroup(c.latestValue, len(c.groups))

	request := &event.Request{
		RequestID:      big.NewInt(c.requestID),
		Payment:        big.NewInt(1),
		PreviousEntry:  c.latestValue,
		GroupPublicKey: c.groups[selectedIdx].groupPublicKey,
		Seed:           seed,
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

func (c *localChain) IsStaleGroup(groupPublicKey []byte) (bool, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	bc, _ := blockCounter()
	bc.WaitForBlockHeight(c.simulatedHeight)
	currentBlock, err := bc.CurrentBlock()

	if err != nil {
		return false, fmt.Errorf("could not determine current block: [%v]", err)
	}

	for _, group := range c.groups {
		if bytes.Compare(group.groupPublicKey, groupPublicKey) == 0 {
			return group.registrationBlockHeight+groupActiveTime+relayRequestTimeout < currentBlock, nil
		}
	}

	return true, nil
}

// IsDKGResultPublished simulates check if the result was already submitted to a
// chain.
func (c *localChain) IsDKGResultSubmitted(requestID *big.Int) (bool, error) {
	c.submittedResultsMutex.Lock()
	defer c.submittedResultsMutex.Unlock()

	return c.submittedResults[requestID.String()] != nil, nil
}

// SubmitDKGResult submits the result to a chain.
func (c *localChain) SubmitDKGResult(
	requestID *big.Int,
	participantIndex group.MemberIndex,
	resultToPublish *relaychain.DKGResult,
	signatures map[group.MemberIndex]operator.Signature,
) *async.DKGResultSubmissionPromise {
	c.submittedResultsMutex.Lock()
	defer c.submittedResultsMutex.Unlock()

	dkgResultPublicationPromise := &async.DKGResultSubmissionPromise{}

	_, ok := c.submittedResults[requestID.String()]
	if ok {
		dkgResultPublicationPromise.Fail(fmt.Errorf(
			"result for request ID [%v] is already submitted",
			requestID,
		))
		return dkgResultPublicationPromise
	}

	c.submittedResults[requestID.String()] = resultToPublish

	currentBlock, err := c.blockCounter.CurrentBlock()
	if err != nil {
		dkgResultPublicationPromise.Fail(fmt.Errorf("cannot read current block"))
		return dkgResultPublicationPromise
	}

	dkgResultPublicationEvent := &event.DKGResultSubmission{
		RequestID:      requestID,
		MemberIndex:    uint32(participantIndex),
		GroupPublicKey: resultToPublish.GroupPublicKey[:],
		BlockNumber:    currentBlock,
	}

	myGroup := localGroup{
		groupPublicKey:          resultToPublish.GroupPublicKey,
		registrationBlockHeight: currentBlock,
	}
	c.groups = append(c.groups, myGroup)

	groupRegistrationEvent := &event.GroupRegistration{
		GroupPublicKey: resultToPublish.GroupPublicKey[:],
		RequestID:      requestID,
		BlockNumber:    currentBlock,
	}

	c.handlerMutex.Lock()
	for _, handler := range c.resultSubmissionHandlers {
		go func(handler func(*event.DKGResultSubmission), dkgResultPublication *event.DKGResultSubmission) {
			handler(dkgResultPublicationEvent)
		}(handler, dkgResultPublicationEvent)
	}

	for _, handler := range c.groupRegisteredHandlers {
		go func(handler func(*event.GroupRegistration), groupRegistration *event.GroupRegistration) {
			handler(groupRegistrationEvent)
		}(handler, groupRegistrationEvent)
	}
	c.handlerMutex.Unlock()

	err = dkgResultPublicationPromise.Fulfill(dkgResultPublicationEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "promise fulfill failed [%v].\n", err)
	}

	return dkgResultPublicationPromise
}

func (c *localChain) OnDKGResultSubmitted(
	handler func(dkgResultPublication *event.DKGResultSubmission),
) (subscription.EventSubscription, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := rand.Int()
	c.resultSubmissionHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.resultSubmissionHandlers, handlerID)
	}), nil
}

func (c *localChain) GetDKGResult(requestID *big.Int) *relaychain.DKGResult {
	return c.submittedResults[requestID.String()]
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
