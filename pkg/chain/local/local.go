package local

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/chain"
)

type localChain struct {
	beaconConfig                     beacon.Config
	groupPublicKeysMutex             sync.Mutex
	groupPublicKeys                  map[string][96]byte
	handlerMutex                     sync.Mutex
	groupPublicKeyFailureHandlers    []func(string, string)
	groupPublicKeySubmissionHandlers []func(string, *big.Int)
	blockCounter                     chain.BlockCounter
}

func (c *localChain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}

func (c *localChain) GetConfig() beacon.Config {
	return c.beaconConfig
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

func (c *localChain) OnGroupPublicKeySubmissionFailed(handler func(string, string)) error {
	c.handlerMutex.Lock()
	c.groupPublicKeyFailureHandlers = append(c.groupPublicKeyFailureHandlers, handler)
	c.handlerMutex.Unlock()

	return nil
}

func (c *localChain) OnGroupPublicKeySubmitted(handler func(groupID string, activationBlock *big.Int)) error {
	c.handlerMutex.Lock()
	c.groupPublicKeySubmissionHandlers = append(c.groupPublicKeySubmissionHandlers, handler)
	c.handlerMutex.Unlock()

	return nil
}

func (c *localChain) RandomBeacon() beacon.ChainInterface {
	return beacon.ChainInterface(c)
}

func (c *localChain) ThresholdRelay() relay.ChainInterface {
	return relay.ChainInterface(c)
}

// Connect initializes a local stub implementation of the chain interfaces for
// testing.
func Connect() chain.Handle {
	bc, _ := blockCounter()
	return &localChain{
		beaconConfig:         beacon.Config{GroupSize: 10, Threshold: 4},
		groupPublicKeysMutex: sync.Mutex{},
		groupPublicKeys:      make(map[string][96]byte),
		blockCounter:         bc,
	}
}
