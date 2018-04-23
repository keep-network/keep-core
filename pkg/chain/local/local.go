package local

import (
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain"
)

type handleStub int

func (handleStub) BlockCounter() chain.BlockCounter {
	return BlockCounter()
}

func (handleStub) GetConfig() beacon.Config {
	return beacon.Config{GroupSize: 10, Threshold: 4}
}

func (handleStub) RandomBeacon() beacon.ChainInterface {
	return handleStub(0)
}

// InitLocal initializes a local stub implementation of the chain interfaces for
// testing.
func InitLocal() chain.Handle {
	return handleStub(0)
}
