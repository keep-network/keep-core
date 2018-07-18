package local

import (
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
	groupRelayEntries      map[int64][32]byte

	handlerMutex               sync.Mutex
	relayEntryHandlers         []func(entry *event.Entry)
	relayRequestHandlers       []func(request *event.Request)
	groupRegisteredHandlers    []func(key *event.GroupRegistration)
	stakerRegistrationHandlers []func(staker *event.StakerRegistration)

	simulatedHeight int64
	blockCounter    chain.BlockCounter

	stakerList []string
}

func (c *localChain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}

func (c *localChain) GetConfig() (relayconfig.Chain, error) {
	return c.relayConfig, nil
}

func (c *localChain) SubmitGroupPublicKey(
	groupID string,
	key [96]byte,
) *async.GroupRegistrationPromise {
	groupRegistrationPromise := &async.GroupRegistrationPromise{}
	c.groupRegistrationsMutex.Lock()
	defer c.groupRegistrationsMutex.Unlock()
	if existing, exists := c.groupRegistrations[groupID]; exists && existing != key {
		fmt.Fprintf(
			os.Stderr,
			"mismatched public key for [%s], submission failed; \n"+
				"[%v] vs [%v]\n",
			groupID,
			existing,
			key,
		)
		return groupRegistrationPromise
	}
	c.groupRegistrations[groupID] = key

	registration := &event.GroupRegistration{
		GroupPublicKey:        []byte(groupID),
		RequestID:             big.NewInt(c.simulatedHeight),
		ActivationBlockHeight: big.NewInt(c.simulatedHeight),
	}
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

	existing, exists := c.groupRelayEntries[entry.GroupID.Int64()]
	if exists && existing != entry.Value {
		err := fmt.Errorf(
			"mismatched signature for [%v], submission failed; \n"+
				"[%v] vs [%v]\n",
			entry.GroupID,
			existing,
			entry.Value,
		)

		relayEntryPromise.Fail(err)

		return relayEntryPromise
	}
	c.groupRelayEntries[entry.GroupID.Int64()] = entry.Value

	c.handlerMutex.Lock()
	for _, handler := range c.relayEntryHandlers {
		go func(handler func(entry *event.Entry), entry *event.Entry) {
			handler(entry)
		}(handler, entry)
	}
	c.handlerMutex.Unlock()

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
		groupRelayEntries:       make(map[int64][32]byte),
		groupRegistrations:      make(map[string][96]byte),
		blockCounter:            bc,
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
		RequestID:   big.NewInt(c.simulatedHeight),
		Payment:     big.NewInt(1),
		BlockReward: blockReward,
		Seed:        seed,
	}
	atomic.AddInt64(&c.simulatedHeight, 1)

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
