package local

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"sync"

	"github.com/ipfs/go-log"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
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
	relayConfig *relaychain.Config

	groups []localGroup

	lastSubmittedDKGResult           *relaychain.DKGResult
	lastSubmittedDKGResultSignatures map[relaychain.GroupMemberIndex][]byte
	lastSubmittedRelayEntry          []byte

	handlerMutex             sync.Mutex
	relayEntryHandlers       map[int]func(entry *event.EntrySubmitted)
	relayRequestHandlers     map[int]func(request *event.Request)
	groupRegisteredHandlers  map[int]func(groupRegistration *event.GroupRegistration)
	dkgStartedHandlers       map[int]func(submission *event.DKGStarted)
	resultSubmissionHandlers map[int]func(submission *event.DKGResultSubmission)

	simulatedHeight uint64
	stakeMonitor    chain.StakeMonitor
	blockCounter    chain.BlockCounter

	relayEntryTimeoutReportsMutex sync.Mutex
	relayEntryTimeoutReports      []uint64

	operatorPrivateKey *operator.PrivateKey

	minimumStake *big.Int
}

func (c *localChain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}

func (c *localChain) StakeMonitor() (chain.StakeMonitor, error) {
	return c.stakeMonitor, nil
}

func (c *localChain) Signing() chain.Signing {
	return newSigner(c.operatorPrivateKey)
}

func (c *localChain) GetKeys() (*operator.PrivateKey, *operator.PublicKey, error) {
	return c.operatorPrivateKey, &c.operatorPrivateKey.PublicKey, nil
}

func (c *localChain) GetConfig() *relaychain.Config {
	return c.relayConfig
}

func (c *localChain) SubmitRelayEntry(newEntry []byte) *async.EventEntrySubmittedPromise {
	relayEntryPromise := &async.EventEntrySubmittedPromise{}

	currentBlock, err := c.blockCounter.CurrentBlock()
	if err != nil {
		failErr := relayEntryPromise.Fail(
			fmt.Errorf("cannot read current block: [%v]", err),
		)
		if failErr != nil {
			logger.Errorf("failed to fail promise: [%v]", failErr)
		}

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

	err = relayEntryPromise.Fulfill(entry)
	if err != nil {
		logger.Errorf("failed to fulfill promise: [%v]", err)
	}

	c.lastSubmittedRelayEntry = newEntry

	return relayEntryPromise
}

func (c *localChain) OnRelayEntrySubmitted(
	handler func(entry *event.EntrySubmitted),
) subscription.EventSubscription {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := generateHandlerID()
	c.relayEntryHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.relayEntryHandlers, handlerID)
	})
}

func (c *localChain) GetLastRelayEntry() []byte {
	return c.lastSubmittedRelayEntry
}

func (c *localChain) OnRelayEntryRequested(
	handler func(request *event.Request),
) subscription.EventSubscription {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := generateHandlerID()
	c.relayRequestHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.relayRequestHandlers, handlerID)
	})
}

func (c *localChain) SelectGroup(seed *big.Int) ([]relaychain.StakerAddress, error) {
	panic("not implemented")
}

func (c *localChain) OnGroupRegistered(
	handler func(groupRegistration *event.GroupRegistration),
) subscription.EventSubscription {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := generateHandlerID()

	c.groupRegisteredHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.groupRegisteredHandlers, handlerID)
	})
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
	operatorPrivateKey, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		panic(err)
	}

	return ConnectWithKey(groupSize, honestThreshold, minimumStake, operatorPrivateKey)
}

// ConnectWithKey initializes a local stub implementation of the chain
// interfaces for testing.
func ConnectWithKey(
	groupSize int,
	honestThreshold int,
	minimumStake *big.Int,
	operatorPrivateKey *operator.PrivateKey,
) Chain {
	bc, _ := BlockCounter()

	currentBlock, _ := bc.CurrentBlock()
	group := localGroup{
		groupPublicKey:          seedGroupPublicKey,
		registrationBlockHeight: currentBlock,
	}

	resultPublicationBlockStep := uint64(3)

	return &localChain{
		relayConfig: &relaychain.Config{
			GroupSize:                  groupSize,
			HonestThreshold:            honestThreshold,
			ResultPublicationBlockStep: resultPublicationBlockStep,
			RelayEntryTimeout:          resultPublicationBlockStep * uint64(groupSize),
		},
		relayEntryHandlers:       make(map[int]func(request *event.EntrySubmitted)),
		relayRequestHandlers:     make(map[int]func(request *event.Request)),
		groupRegisteredHandlers:  make(map[int]func(groupRegistration *event.GroupRegistration)),
		dkgStartedHandlers:       make(map[int]func(submission *event.DKGStarted)),
		resultSubmissionHandlers: make(map[int]func(submission *event.DKGResultSubmission)),
		blockCounter:             bc,
		stakeMonitor:             NewStakeMonitor(minimumStake),
		groups:                   []localGroup{group},
		operatorPrivateKey:       operatorPrivateKey,
		minimumStake:             minimumStake,
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

	err := bc.WaitForBlockHeight(c.simulatedHeight)
	if err != nil {
		logger.Errorf("could not wait for block height: [%v]", err)
	}

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
		err := dkgResultPublicationPromise.Fail(fmt.Errorf(
			"failed to submit result with [%v] signatures for honest threshold [%v]",
			len(signatures),
			c.relayConfig.HonestThreshold,
		))
		if err != nil {
			logger.Errorf("failed to fail promise: [%v]", err)
		}

		return dkgResultPublicationPromise
	}

	currentBlock, err := c.blockCounter.CurrentBlock()
	if err != nil {
		failErr := dkgResultPublicationPromise.Fail(
			fmt.Errorf("cannot read current block: [%v]", err),
		)
		if failErr != nil {
			logger.Errorf("failed to fail promise: [%v]", failErr)
		}

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
		logger.Errorf("failed to fulfill promise: [%v]", err)
	}

	return dkgResultPublicationPromise
}

func (c *localChain) OnDKGStarted(
	handler func(event *event.DKGStarted),
) subscription.EventSubscription {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := generateHandlerID()
	c.dkgStartedHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.dkgStartedHandlers, handlerID)
	})
}

func (c *localChain) OnDKGResultSubmitted(
	handler func(dkgResultPublication *event.DKGResultSubmission),
) subscription.EventSubscription {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	handlerID := generateHandlerID()
	c.resultSubmissionHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		c.handlerMutex.Lock()
		defer c.handlerMutex.Unlock()

		delete(c.resultSubmissionHandlers, handlerID)
	})
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

func (c *localChain) IsEntryInProgress() (bool, error) {
	panic("not implemented")
}

func (c *localChain) CurrentRequestStartBlock() (*big.Int, error) {
	panic("not implemented")
}

func (c *localChain) CurrentRequestPreviousEntry() ([]byte, error) {
	panic("not implemented")
}

func (c *localChain) CurrentRequestGroupPublicKey() ([]byte, error) {
	panic("not implemented")
}

func (c *localChain) GetRelayEntryTimeoutReports() []uint64 {
	return c.relayEntryTimeoutReports
}

func (c *localChain) MinimumStake() (*big.Int, error) {
	return c.minimumStake, nil
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

func generateHandlerID() int {
	// #nosec G404 (insecure random number source (rand))
	// Local chain implementation doesn't require secure randomness.
	return rand.Int()
}
