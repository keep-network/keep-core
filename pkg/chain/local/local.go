package local

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"sync"

	"github.com/ipfs/go-log/v2"

	crand "crypto/rand"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/subscription"
	"golang.org/x/crypto/sha3"
)

var logger = log.Logger("keep-chain-local")

var seedGroupPublicKey = []byte("seed to group public key")
var seedRelayEntry = big.NewInt(123456789)
var groupActiveTime = uint64(10)
var relayRequestTimeout = uint64(8)

// Chain is an extention of chain.Handle interface which exposes
// additional functions useful for testing.
type Chain interface {
	chain.Handle

	// GetLastDKGResult returns the last DKG result submitted to the chain
	// as well as all the signatures that supported that result.
	GetLastDKGResult() (*relaychain.DKGResult, map[relaychain.GroupMemberIndex][]byte)

	// GetLastRelayEntry returns the last relay entry submitted to the chain.
	GetLastRelayEntry() []byte

	// GetRelayEntryTimeoutReports returns an array of blocks which denote at what
	// block a relay entry timeout occured.
	GetRelayEntryTimeoutReports() []uint64
}

type localGroup struct {
	groupPublicKey          []byte
	registrationBlockHeight uint64
}

type localChain struct {
	relayConfig *relayconfig.Chain

	groups []localGroup

	lastSubmittedDKGResult           *relaychain.DKGResult
	lastSubmittedDKGResultSignatures map[relaychain.GroupMemberIndex][]byte
	lastSubmittedRelayEntry          []byte

	handlerMutex                  sync.Mutex
	relayEntryHandlers            map[int]func(entry *event.EntrySubmitted)
	relayRequestHandlers          map[int]func(request *event.Request)
	groupSelectionStartedHandlers map[int]func(groupSelectionStart *event.GroupSelectionStart)
	groupRegisteredHandlers       map[int]func(groupRegistration *event.GroupRegistration)
	resultSubmissionHandlers      map[int]func(submission *event.DKGResultSubmission)

	simulatedHeight uint64
	stakeMonitor    chain.StakeMonitor
	blockCounter    chain.BlockCounter

	tickets      []*relaychain.Ticket
	ticketsMutex sync.Mutex

	relayEntryTimeoutReportsMutex sync.Mutex
	relayEntryTimeoutReports      []uint64

	operatorKey *ecdsa.PrivateKey
}

func (c *localChain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}

func (c *localChain) StakeMonitor() (chain.StakeMonitor, error) {
	return c.stakeMonitor, nil
}

func (c *localChain) Signing() chain.Signing {
	return &localSigning{c.operatorKey}
}

func (c *localChain) GetKeys() (*operator.PrivateKey, *operator.PublicKey) {
	return c.operatorKey, &c.operatorKey.PublicKey
}

func (c *localChain) GetConfig() (*relayconfig.Chain, error) {
	return c.relayConfig, nil
}

func (c *localChain) SubmitTicket(ticket *relaychain.Ticket) *async.EventGroupTicketSubmissionPromise {
	promise := &async.EventGroupTicketSubmissionPromise{}

	c.ticketsMutex.Lock()
	defer c.ticketsMutex.Unlock()

	c.tickets = append(c.tickets, ticket)
	sort.SliceStable(c.tickets, func(i, j int) bool {
		// Ticket value bytes are interpreted as a big-endian unsigned integers.
		iValue := new(big.Int).SetBytes(c.tickets[i].Value[:])
		jValue := new(big.Int).SetBytes(c.tickets[j].Value[:])

		return iValue.Cmp(jValue) == -1
	})

	_ = promise.Fulfill(&event.GroupTicketSubmission{
		TicketValue: new(big.Int).SetBytes(ticket.Value[:]),
		BlockNumber: c.simulatedHeight,
	})

	return promise
}

func (c *localChain) GetSubmittedTickets() ([]uint64, error) {
	tickets := make([]uint64, len(c.tickets))

	for i := range tickets {
		tickets[i] = binary.BigEndian.Uint64(c.tickets[i].Value[:])
	}

	return tickets, nil
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

func (c *localChain) SubmitRelayEntry(newEntry []byte) *async.EventEntrySubmittedPromise {
	c.ticketsMutex.Lock()
	c.tickets = make([]*relaychain.Ticket, 0)
	c.ticketsMutex.Unlock()

	relayEntryPromise := &async.EventEntrySubmittedPromise{}

	currentBlock, err := c.blockCounter.CurrentBlock()
	if err != nil {
		relayEntryPromise.Fail(fmt.Errorf("cannot read current block"))
		return relayEntryPromise
	}

	entry := &event.EntrySubmitted{
		BlockNumber: currentBlock,
	}

	c.handlerMutex.Lock()
	for _, handler := range c.relayEntryHandlers {
		go func(handler func(entry *event.EntrySubmitted), entry *event.EntrySubmitted) {
			handler(entry)
		}(handler, entry)
	}
	c.handlerMutex.Unlock()

	relayEntryPromise.Fulfill(entry)

	c.lastSubmittedRelayEntry = newEntry

	return relayEntryPromise
}

func (c *localChain) OnRelayEntrySubmitted(
	handler func(entry *event.EntrySubmitted),
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

func (c *localChain) GetLastRelayEntry() []byte {
	return c.lastSubmittedRelayEntry
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

// Connect initializes a local stub implementation of the chain
// interfaces for testing. It uses auto-generated operator key.
func Connect(
	groupSize int,
	honestThreshold int,
	minimumStake *big.Int,
) Chain {
	operatorKey, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		panic(err)
	}

	return ConnectWithKey(groupSize, honestThreshold, minimumStake, operatorKey)
}

// ConnectWithKey initializes a local stub implementation of the chain
// interfaces for testing.
func ConnectWithKey(
	groupSize int,
	honestThreshold int,
	minimumStake *big.Int,
	operatorKey *ecdsa.PrivateKey,
) Chain {
	bc, _ := BlockCounter()

	currentBlock, _ := bc.CurrentBlock()
	group := localGroup{
		groupPublicKey:          seedGroupPublicKey,
		registrationBlockHeight: currentBlock,
	}

	resultPublicationBlockStep := uint64(3)

	return &localChain{
		relayConfig: &relayconfig.Chain{
			GroupSize:                  groupSize,
			HonestThreshold:            honestThreshold,
			TicketSubmissionTimeout:    6,
			ResultPublicationBlockStep: resultPublicationBlockStep,
			MinimumStake:               minimumStake,
			RelayEntryTimeout:          resultPublicationBlockStep * uint64(groupSize),
		},
		relayEntryHandlers:       make(map[int]func(request *event.EntrySubmitted)),
		relayRequestHandlers:     make(map[int]func(request *event.Request)),
		groupRegisteredHandlers:  make(map[int]func(groupRegistration *event.GroupRegistration)),
		resultSubmissionHandlers: make(map[int]func(submission *event.DKGResultSubmission)),
		blockCounter:             bc,
		stakeMonitor:             NewStakeMonitor(minimumStake),
		tickets:                  make([]*relaychain.Ticket, 0),
		groups:                   []localGroup{group},
		operatorKey:              operatorKey,
	}
}

func selectGroup(entry *big.Int, numberOfGroups int) int {
	if numberOfGroups == 0 {
		return 0
	}

	return int(new(big.Int).Mod(entry, big.NewInt(int64(numberOfGroups))).Int64())
}

func (c *localChain) IsStaleGroup(groupPublicKey []byte) (bool, error) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	bc, _ := BlockCounter()
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

func (c *localChain) GetGroupMembers(groupPublicKey []byte) (
	[]relaychain.StakerAddress,
	error,
) {
	return nil, nil // no-op
}

func (c *localChain) IsGroupRegistered(groupPublicKey []byte) (bool, error) {
	for _, group := range c.groups {
		if bytes.Compare(group.groupPublicKey, groupPublicKey) == 0 {
			return true, nil
		}
	}
	return false, nil
}

// SubmitDKGResult submits the result to a chain.
func (c *localChain) SubmitDKGResult(
	participantIndex relaychain.GroupMemberIndex,
	resultToPublish *relaychain.DKGResult,
	signatures map[relaychain.GroupMemberIndex][]byte,
) *async.EventDKGResultSubmissionPromise {
	dkgResultPublicationPromise := &async.EventDKGResultSubmissionPromise{}

	if len(signatures) < c.relayConfig.HonestThreshold {
		dkgResultPublicationPromise.Fail(fmt.Errorf(
			"failed to submit result with [%v] signatures for honest threshold [%v]",
			len(signatures),
			c.relayConfig.HonestThreshold,
		))
		return dkgResultPublicationPromise
	}

	currentBlock, err := c.blockCounter.CurrentBlock()
	if err != nil {
		dkgResultPublicationPromise.Fail(fmt.Errorf("cannot read current block"))
		return dkgResultPublicationPromise
	}

	dkgResultPublicationEvent := &event.DKGResultSubmission{
		MemberIndex:    uint32(participantIndex),
		GroupPublicKey: resultToPublish.GroupPublicKey[:],
		Misbehaved:     resultToPublish.Misbehaved,
		BlockNumber:    currentBlock,
	}

	myGroup := localGroup{
		groupPublicKey:          resultToPublish.GroupPublicKey,
		registrationBlockHeight: currentBlock,
	}
	c.groups = append(c.groups, myGroup)
	c.lastSubmittedDKGResult = resultToPublish
	c.lastSubmittedDKGResultSignatures = signatures

	groupRegistrationEvent := &event.GroupRegistration{
		GroupPublicKey: resultToPublish.GroupPublicKey[:],
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
		logger.Errorf("failed to fulfill promise: [%v].", err)
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

func (c *localChain) GetLastDKGResult() (
	*relaychain.DKGResult,
	map[relaychain.GroupMemberIndex][]byte,
) {
	return c.lastSubmittedDKGResult, c.lastSubmittedDKGResultSignatures
}

func (c *localChain) ReportRelayEntryTimeout() error {
	c.relayEntryTimeoutReportsMutex.Lock()
	defer c.relayEntryTimeoutReportsMutex.Unlock()

	currentBlock, err := c.blockCounter.CurrentBlock()
	if err != nil {
		return err
	}

	c.relayEntryTimeoutReports = append(c.relayEntryTimeoutReports, currentBlock)
	return nil
}

func (c *localChain) GetRelayEntryTimeoutReports() []uint64 {
	return c.relayEntryTimeoutReports
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
