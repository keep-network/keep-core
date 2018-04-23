package ethereumblock

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain"
)

type handleStub int

func (handleStub) BlockCounter() chain.BlockCounter {
	client, err := rpc.Dial("ws://192.168.0.158:8546")
	if err != nil {
		fmt.Println("Error Connecting to Server", err)

	}
	return BlockCounter(client)
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
