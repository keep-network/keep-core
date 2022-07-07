package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain/random-beacon/gen/contract"
)

// Definitions of contract names.
const (
	RandomBeaconContractName = "RandomBeacon"
)

// BeaconChain represents a beacon-specific chain handle.
type BeaconChain struct {
	*Chain

	randomBeacon *contract.RandomBeacon
}

// NewBeaconChain construct a new instance of the beacon-specific Ethereum
// chain handle.
func NewBeaconChain(
	ctx context.Context,
	config ethereum.Config,
	client *ethclient.Client,
) (*BeaconChain, error) {
	chain, err := NewChain(ctx, config, client)
	if err != nil {
		return nil, fmt.Errorf("cannot create base chain handle: [%v]", err)
	}

	randomBeaconAddress, err := config.ContractAddress(RandomBeaconContractName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve %s contract address: [%v]",
			RandomBeaconContractName,
			err,
		)
	}

	randomBeacon, err :=
		contract.NewRandomBeacon(
			randomBeaconAddress,
			chain.chainID,
			chain.key,
			chain.client,
			chain.nonceManager,
			chain.miningWaiter,
			chain.blockCounter,
			chain.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to %s contract: [%v]",
			RandomBeaconContractName,
			err,
		)
	}

	return &BeaconChain{
		Chain:        chain,
		randomBeacon: randomBeacon,
	}, nil
}
