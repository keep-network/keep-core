package local

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"sync"
	"sync/atomic"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

type localChain struct {
	relayConfig relayconfig.Chain

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
	stakerRegistrationHandlers   []func(staker *event.StakerRegistration)
	dkgResultPublicationHandlers []func(dkgResultPublication *event.DKGResultPublication)

	requestID   int64
	latestValue *big.Int

	simulatedHeight int64
	stakeMonitor    chain.StakeMonitor
	blockCounter    chain.BlockCounter

	stakerList []string

	// Track the submitted votes -
	submissionsMutex sync.Mutex
	submissions      map[*big.Int]relaychain.Submissions
	// vote handler
	voteHandler []func(dkgResultVote *event.DKGResultVote)
}

/*
From "submissions.go"
====================================================================
// Submissions - PHASE 14
type Submissions struct {
	requestID   *big.Int
	Submissions []*Submission
}

// Submission - PHASE 14
type Submission struct {
	DKGResult *DKGResult
	Votes     int
}

// Lead returns a submission with the highest number of votes.
func (s *Submissions) Lead() *Submission {
	top := -1
	topPos := 0
	for pos, aSubmission := range s.submissions {
		if top < aSubmission.Votes {
			topPos = pos
			top = aSubmission.Votes
		}
	}
	return s.submissions[topPos]
}

func (s *Submissions) Contains(result *DKGResult) bool {
	// TODO Implement
	return false
}
====================================================================
*/

// PHASE 14
// from relaychain:
//	// SubmitDKGResult sends DKG result to a chain.
//	SubmitDKGResult(requestID *big.Int, dkgResult *DKGResult) *async.DKGResultPublicationPromise
//
// Sets:	submittedResults map[*big.Int][]*relaychain.DKGResult
func (c *localChain) GetDKGSubmissions(requestID *big.Int) *relaychain.Submissions {
	c.submissionsMutex.Lock()
	defer c.submissionsMutex.Unlock()
	x := c.submissions[requestID]
	return &x
}

// Vote places a vote for dkgResultHash and causes OnDKGResultVote event to occurs.
// PHASE 14
// func (c *localChain) SubmitDKGResult( -- submits a result - iof we vote for a result
// that is non-existent - do we just ignore - or error - or create?
func (c *localChain) Vote(requestID *big.Int, dkgResultHash []byte) {
	c.submissionsMutex.Lock()
	defer c.submissionsMutex.Unlock()
	x, ok := c.submissions[requestID]
	if !ok {
		fmt.Printf("Missing requestID in c.submissions - early return\n")
		return
	}
	for pos, sub := range x.Submissions {
		fmt.Printf("At pos=%d\n", pos)
		if bytes.Equal(sub.DKGResult.Hash(), dkgResultHash) {
			fmt.Printf("At pos=%d, match for DKGResult\n", pos)
			sub.Votes++
			x.Submissions[pos] = sub
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
			return
		}
	}
}

// PHASE 14
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

func (c *localChain) GetConfig() (relayconfig.Chain, error) {
	return c.relayConfig, nil
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

	// Like This... xyzzy100
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
		if existing != entry.Value {
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

func (c *localChain) OnStakerAdded(handler func(staker *event.StakerRegistration)) {
	c.handlerMutex.Lock()
	c.stakerRegistrationHandlers = append(
		c.stakerRegistrationHandlers,
		handler,
	)
	c.handlerMutex.Unlock()
}

func (c *localChain) ThresholdRelay() relaychain.Interface {
	return relaychain.Interface(c)
}

// Connect initializes a local stub implementation of the chain interfaces
// for testing.
func Connect(groupSize int, threshold int) chain.Handle {
	bc, _ := blockCounter()

	return &localChain{
		relayConfig: relayconfig.Chain{
			GroupSize: groupSize,
			Threshold: threshold,
		},
		groupRegistrationsMutex: sync.Mutex{},
		groupRelayEntries:       make(map[string]*big.Int),
		groupRegistrations:      make(map[string][96]byte),
		submittedResults:        make(map[*big.Int][]*relaychain.DKGResult),
		blockCounter:            bc,
		stakeMonitor:            NewStakeMonitor(),
		submissions:             make(map[*big.Int]relaychain.Submissions),
	}
}

// AddStaker is a temporary function for Milestone 1 that
// adds a staker to the group contract.
func (c *localChain) AddStaker(
	groupMemberID string,
) *async.StakerRegistrationPromise {
	onStakerAddedPromise := &async.StakerRegistrationPromise{}
	index := len(c.stakerList)
	c.stakerList = append(c.stakerList, groupMemberID)

	c.handlerMutex.Lock()
	for _, handler := range c.stakerRegistrationHandlers {
		go func(handler func(staker *event.StakerRegistration), groupMemberID string, index int) {
			handler(&event.StakerRegistration{
				GroupMemberID: groupMemberID,
				Index:         index,
			})
		}(handler, groupMemberID, index)
	}
	c.handlerMutex.Unlock()

	err := onStakerAddedPromise.Fulfill(&event.StakerRegistration{
		Index:         index,
		GroupMemberID: string(groupMemberID),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Promise Fulfill failed [%v].\n", err)
	}

	return onStakerAddedPromise
}

// GetStakerList is a temporary function for Milestone 1 that
// gets back the list of stakers.
func (c *localChain) GetStakerList() ([]string, error) {
	return c.stakerList, nil
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
func (c *localChain) IsDKGResultPublished(
	requestID *big.Int, result *relaychain.DKGResult,
) bool {
	for publishedRequestID, publishedResults := range c.submittedResults {
		if publishedRequestID.Cmp(requestID) == 0 {
			for _, publishedResult := range publishedResults {
				if publishedResult.Equals(result) {
					return true
				}
			}
			return false
		}
	}
	return false
}

/*
From "submissions.go"
====================================================================
// Submissions - PHASE 14
type Submissions struct {
	requestID   *big.Int
	Submissions []*Submission
}

// Submission - PHASE 14
type Submission struct {
	DKGResult *DKGResult
	Votes     int
}
====================================================================
*/
// SubmitDKGResult submits the result to a chain.
func (c *localChain) SubmitDKGResult(
	requestID *big.Int, resultToPublish *relaychain.DKGResult,
) *async.DKGResultPublicationPromise {
	c.submittedResultsMutex.Lock()
	defer c.submittedResultsMutex.Unlock()

	dkgResultPublicationPromise := &async.DKGResultPublicationPromise{}

	if c.IsDKGResultPublished(requestID, resultToPublish) {
		dkgResultPublicationPromise.Fail(fmt.Errorf("result already submitted"))
		return dkgResultPublicationPromise
	}

	c.submittedResults[requestID] = append(c.submittedResults[requestID], resultToPublish)

	c.submissionsMutex.Lock()
	if c.submissions == nil {
		c.submissions = make(map[*big.Int]relaychain.Submissions)
	}
	if _, ok := c.submissions[requestID]; !ok {
		fmt.Printf("Addin in a 1st vote\n")
		// submissions      map[*big.Int]relaychain.Submissions
		ss := relaychain.Submissions{
			Submissions: []*relaychain.Submission{
				{
					DKGResult: &relaychain.DKGResult{
						Success:        true,
						GroupPublicKey: big.NewInt(1001), // where is this from TODO - FIXME
						Disqualified:   []bool{},
						Inactive:       []bool{},
					},
					Votes: 1,
				},
			},
		}
		c.submissions[requestID] = ss
	}
	c.submissionsMutex.Unlock()

	dkgResultPublicationEvent := &event.DKGResultPublication{RequestID: requestID}

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
) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()

	c.dkgResultPublicationHandlers = append(
		c.dkgResultPublicationHandlers,
		handler,
	)
}
