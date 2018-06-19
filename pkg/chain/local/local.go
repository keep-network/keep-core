package local

import (
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/callback"
	"github.com/keep-network/keep-core/pkg/chain"
)

type localChain struct {
	relayConfig                      relay.Config
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

func (c *localChain) GetConfig() (relay.Config, error) {
	return c.relayConfig, nil
}

func (ec *localChain) SubmitGroupPublicKey(
	groupID string,
	key [96]byte,
) *callback.Promise {
	return &callback.Promise{}
}

func (c *localChain) ThresholdRelay() relay.ChainInterface {
	return relay.ChainInterface(c)
}

// Connect initializes a local stub implementation of the chain interfaces for
// testing.
func Connect(groupSize int, threshold int) chain.Handle {
	bc, _ := blockCounter()

	return &localChain{
		relayConfig: relay.Config{
			GroupSize: groupSize,
			Threshold: threshold,
		},
		groupPublicKeysMutex: sync.Mutex{},
		groupPublicKeys:      make(map[string][96]byte),
		blockCounter:         bc,
	}
}
