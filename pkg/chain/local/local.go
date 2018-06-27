package local

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

type localChain struct {
	relayConfig relayconfig.Chain

	groupPublicKeysMutex sync.Mutex
	groupPublicKeys      map[string][96]byte

	groupRelayEntriesMutex sync.Mutex
	groupRelayEntries      map[*big.Int][32]byte

	handlerMutex                     sync.Mutex
	groupPublicKeyFailureHandlers    []func(string, string)
	groupPublicKeySubmissionHandlers []func(string, *big.Int)

	blockCounter chain.BlockCounter
}

func (c *localChain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}

func (c *localChain) GetConfig() (relayconfig.Chain, error) {
	return c.relayConfig, nil
}

func (c *localChain) SubmitGroupPublicKey(groupID string, key [96]byte) error {
	c.groupPublicKeysMutex.Lock()
	defer c.groupPublicKeysMutex.Unlock()
	if existing, exists := c.groupPublicKeys[groupID]; exists && existing != key {
		errorMsg := fmt.Sprintf(
			"mismatched public key for [%s], submission failed; \n"+
				"[%v] vs [%v]\n",
			groupID,
			existing,
			key,
		)

		c.handlerMutex.Lock()
		for _, handler := range c.groupPublicKeyFailureHandlers {
			handler(groupID, errorMsg)
		}
		c.handlerMutex.Unlock()

		return nil
	}
	c.groupPublicKeys[groupID] = key

	c.handlerMutex.Lock()
	for _, handler := range c.groupPublicKeySubmissionHandlers {
		handler(groupID, &big.Int{})
	}
	c.handlerMutex.Unlock()

	return nil
}

func (c *localChain) OnGroupPublicKeySubmissionFailed(
	handler func(string, string),
) error {
	c.handlerMutex.Lock()
	c.groupPublicKeyFailureHandlers = append(c.groupPublicKeyFailureHandlers, handler)
	c.handlerMutex.Unlock()

	return nil
}

func (c *localChain) OnGroupPublicKeySubmitted(
	handler func(groupID string, activationBlock *big.Int),
) error {
	c.handlerMutex.Lock()
	c.groupPublicKeySubmissionHandlers = append(
		c.groupPublicKeySubmissionHandlers,
		handler,
	)
	c.handlerMutex.Unlock()

	return nil
}

func (c *localChain) SubmitRelayEntry(entry *relay.Entry) *async.RelayEntryPromise {
	relayEntryPromise := &async.RelayEntryPromise{}

	c.groupRelayEntriesMutex.Lock()
	defer c.groupRelayEntriesMutex.Unlock()

	if existing, exists := c.groupRelayEntries[entry.GroupID]; exists && existing != entry.Value {
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
	c.groupRelayEntries[entry.GroupID] = entry.Value

	relayEntryPromise.Fulfill(&relay.Entry{
		RequestID:     big.NewInt(int64(0)),
		Value:         entry.Value,
		GroupID:       entry.GroupID,
		PreviousEntry: big.NewInt(int64(0)),
		Timestamp:     time.Now().UTC(),
	})

	return relayEntryPromise
}

func (c *localChain) ThresholdRelay() relaychain.Interface {
	return relaychain.Interface(c)
}

// Connect initializes a local stub implementation of the chain interfaces for
// testing.
func Connect(groupSize int, threshold int) chain.Handle {
	bc, _ := blockCounter()

	return &localChain{
		relayConfig: relayconfig.Chain{
			GroupSize: groupSize,
			Threshold: threshold,
		},
		groupPublicKeysMutex: sync.Mutex{},
		groupRelayEntries:    make(map[*big.Int][32]byte),
		groupPublicKeys:      make(map[string][96]byte),
		blockCounter:         bc,
	}
}
