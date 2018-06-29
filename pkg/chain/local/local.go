package local

import (
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon/chaintype"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	relayconfig "github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

type localChain struct {
	relayConfig          relayconfig.Chain
	groupPublicKeysMutex sync.Mutex
	groupPublicKeys      map[string][96]byte
	blockCounter         chain.BlockCounter
	simulatedHeight      int64
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
) *async.GroupPublicKeyPromise {
	groupKeyPromise := &async.GroupPublicKeyPromise{}
	c.groupPublicKeysMutex.Lock()
	defer c.groupPublicKeysMutex.Unlock()
	if existing, exists := c.groupPublicKeys[groupID]; exists && existing != key {
		fmt.Fprintf(
			os.Stderr,
			"mismatched public key for [%s], submission failed; \n"+
				"[%v] vs [%v]\n",
			groupID,
			existing,
			key,
		)
		return groupKeyPromise
	}
	c.groupPublicKeys[groupID] = key
	c.simulatedHeight++

	groupKeyPromise.Fulfill(&chaintype.GroupPublicKey{
		GroupPublicKey:        []byte(groupID),
		RequestID:             big.NewInt(c.simulatedHeight),
		ActivationBlockHeight: big.NewInt(c.simulatedHeight),
	})

	return groupKeyPromise
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
		groupPublicKeysMutex: sync.Mutex{},
		groupPublicKeys:      make(map[string][96]byte),
		blockCounter:         bc,
	}
}
