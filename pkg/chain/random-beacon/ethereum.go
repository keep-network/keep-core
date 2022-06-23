package beacon

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-common/pkg/chain/ethlike"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain/random-beacon/gen/contract"
	"github.com/keep-network/keep-core/pkg/chain/sortition"
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
func Connect(ec chain.Handle) (*EthereumHandle, error) {
	eh := &EthereumHandle{}
	eh.Chain = ec.(*ethereum.Chain)

	randomBeaconAddress, err := eh.Chain.Config().ContractAddress(RandomBeaconContractName)
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

// OperatorToStakingProvider ...
func (eh *EthereumHandle) OperatorToStakingProvider() (string, error) {
	stakingProvider, err := eh.randomBeacon.OperatorToStakingProvider(eh.Address())
	if err != nil {
		return "", fmt.Errorf(
			"failed to map operator %v to a staking provider: [%w]",
			eh.Address(),
			err,
		)
	}
	if bytes.Equal(
		stakingProvider.Bytes(),
		bytes.Repeat([]byte{0}, common.AddressLength),
	) {
		return "", sortition.ErrOperatorNotRegistered
	}

	return stakingProvider.Hex(), nil
}

// EligibleStake ...
func (eh *EthereumHandle) EligibleStake(stakingProvider string) (*big.Int, error) {
	eligibleStake, err := eh.randomBeacon.EligibleStake(common.HexToAddress(stakingProvider))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get eligible stake for staking provider %s: [%w]",
			stakingProvider,
			err,
		)
	}

	return eligibleStake, nil
}

// IsOperatorInPool checks if the operator is in the sortition pool.
func (eh *EthereumHandle) IsOperatorInPool() (bool, error) {
	return eh.randomBeacon.IsOperatorInPool(eh.Address())
}

// JoinSortitionPool joins the sortition pool.
func (eh *EthereumHandle) JoinSortitionPool() error {
	_, err := eh.randomBeacon.JoinSortitionPool()

	return err
}
