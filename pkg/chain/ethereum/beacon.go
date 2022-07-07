package ethereum

import (
	"context"
	"fmt"

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

// newBeaconChain construct a new instance of the beacon-specific Ethereum
// chain handle.
func newBeaconChain(
	ctx context.Context,
	config ethereum.Config,
	baseChain *Chain,
) (*BeaconChain, error) {
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
			baseChain.chainID,
			baseChain.key,
			baseChain.client,
			baseChain.nonceManager,
			baseChain.miningWaiter,
			baseChain.blockCounter,
			baseChain.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to %s contract: [%v]",
			RandomBeaconContractName,
			err,
		)
	}

	return &BeaconChain{
		Chain:        baseChain,
		randomBeacon: randomBeacon,
	}, nil
}
