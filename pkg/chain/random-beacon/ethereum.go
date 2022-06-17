package beacon

import (
	"fmt"

	"github.com/keep-network/keep-common/pkg/chain/ethlike"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain/random-beacon/gen/contract"
)

// Definitions of contract names.
const (
	RandomBeaconContractName = "RandomBeacon"
)

// EthereumHandle is a handle for interaction with the ethereum contracts of the
// Random Beacon module.
type EthereumHandle struct {
	*ethereum.Chain

	randomBeacon *contract.RandomBeacon
}

// Connect connects to chain.
func Connect(ec *ethereum.Chain) (*EthereumHandle, error) {
	eh := &EthereumHandle{}
	eh.Chain = ec

	randomBeaconAddress, err := ec.Config().ContractAddress(RandomBeaconContractName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve %s contract address: [%v]",
			RandomBeaconContractName,
			err,
		)
	}

	blockCounter, err := eh.BlockCounter()
	if err != nil {
		return nil, fmt.Errorf("failed to get block counter [%v]", err)
	}

	eh.randomBeacon, err =
		contract.NewRandomBeacon(
			randomBeaconAddress,
			eh.ChainID(),
			eh.AccountKey(),
			eh.Client(),
			eh.NonceManager(),
			eh.MiningWaiter(),
			blockCounter.(*ethlike.BlockCounter),
			eh.TransactionMutex(),
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to %s contract: [%v]",
			RandomBeaconContractName,
			err,
		)
	}

	return eh, nil
}

// JoinSortitionPool joins the sortition pool.
func (rb *EthereumHandle) JoinSortitionPool() error {
	_, err := rb.randomBeacon.JoinSortitionPool()

	return err
}
