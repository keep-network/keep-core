package local

import (
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/chain"
)

type localChain struct {
	beaconConfig         beacon.Config
	groupPublicKeysMutex sync.Mutex
	groupPublicKeys      map[string][96]byte
	blockCounter         chain.BlockCounter
}

func (c *localChain) BlockCounter() chain.BlockCounter {
	return c.blockCounter
}

func (c *localChain) GetConfig() beacon.Config {
	return c.beaconConfig
}

func (c *localChain) SubmitGroupPublicKey(groupID string, key [96]byte) error {
	c.groupPublicKeysMutex.Lock()
	defer c.groupPublicKeysMutex.Unlock()
	if existing, exists := c.groupPublicKeys[groupID]; exists && existing != key {
		return fmt.Errorf(
			"mismatched public key for [%s], submission failed; \n"+
				"[%v] vs [%v]",
			groupID,
			existing,
			key)
	}
	c.groupPublicKeys[groupID] = key

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
	return &localChain{
		beaconConfig:         beacon.Config{GroupSize: 10, Threshold: 4},
		groupPublicKeysMutex: sync.Mutex{},
		groupPublicKeys:      make(map[string][96]byte),
		blockCounter:         blockCounter()}
}
