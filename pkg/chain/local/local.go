package local

import (
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/callback"
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

func (c *localChain) GetConfig() (beacon.Config, error) {
	return c.beaconConfig, nil
}

func (ec *localChain) SubmitGroupPublicKey(
	groupID string,
	key [96]byte,
) *callback.Promise {
	return &callback.Promise{}
}

func (c *localChain) RandomBeacon() beacon.ChainInterface {
	return beacon.ChainInterface(c)
}

func (c *localChain) ThresholdRelay() relay.ChainInterface {
	return relay.ChainInterface(c)
}

// Connect initializes a local stub implementation of the chain interfaces for
// testing.
func Connect(groupSize int, threshold int) chain.Handle {
	bc, _ := blockCounter()

	return &localChain{
		beaconConfig: beacon.Config{
			GroupSize: groupSize,
			Threshold: threshold,
		},
		groupPublicKeysMutex: sync.Mutex{},
		groupPublicKeys:      make(map[string][96]byte),
		blockCounter:         bc,
	}
}
