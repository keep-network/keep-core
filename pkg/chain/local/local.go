package local

import (
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/entry"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

type localChain struct {
	relayConfig relayconfig.Chain

	groupRegistrationsMutex sync.Mutex
	groupRegistrations      map[string][96]byte

	groupRelayEntriesMutex sync.Mutex
	groupRelayEntries      map[int64][32]byte

	handlerMutex            sync.Mutex
	relayEntryHandlers      []func(entry *relay.Entry)
	relayRequestHandlers    []func(request *entry.Request)
	groupRegisteredHandlers []func(key *relay.GroupRegistration)

	simulatedHeight int64
	blockCounter    chain.BlockCounter
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
	c.simulatedHeight++

	groupRegistrationPromise.Fulfill(&relay.GroupRegistration{
		GroupPublicKey:        []byte(groupID),
		RequestID:             big.NewInt(c.simulatedHeight),
		ActivationBlockHeight: big.NewInt(c.simulatedHeight),
	})

	return groupRegistrationPromise
}

func (c *localChain) SubmitRelayEntry(entry *relay.Entry) *async.RelayEntryPromise {
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

	relayEntryPromise.Fulfill(&relay.Entry{
		RequestID:     entry.RequestID,
		Value:         entry.Value,
		GroupID:       entry.GroupID,
		PreviousEntry: entry.PreviousEntry,
		Timestamp:     time.Now().UTC(),
	})

	return relayEntryPromise
}

func (c *localChain) OnRelayEntryGenerated(handler func(entry *relay.Entry)) {
	c.handlerMutex.Lock()
	c.relayEntryHandlers = append(
		c.relayEntryHandlers,
		handler,
	)
	c.handlerMutex.Unlock()
}

func (c *localChain) OnRelayEntryRequested(handler func(request *entry.Request)) {
	c.handlerMutex.Lock()
	c.relayRequestHandlers = append(
		c.relayRequestHandlers,
		handler,
	)
	c.handlerMutex.Unlock()
}

func (c *localChain) OnGroupRegistered(handler func(key *relay.GroupRegistration)) {
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
