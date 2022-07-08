package ethereum

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/random-beacon/gen/contract"
)

// Definitions of contract names.
const (
	RandomBeaconContractName = "RandomBeacon"
)

// BeaconChain represents a beacon-specific chain handle.
type BeaconChain struct {
	*Chain

	randomBeacon  *contract.RandomBeacon
	sortitionPool *contract.BeaconSortitionPool
}

// newBeaconChain construct a new instance of the beacon-specific Ethereum
// chain handle.
func newBeaconChain(
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
			"failed to attach to RandomBeacon contract: [%v]",
			err,
		)
	}

	sortitionPoolAddress, err := randomBeacon.SortitionPool()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get sortition pool address: [%v]",
			err,
		)
	}

	sortitionPool, err :=
		contract.NewBeaconSortitionPool(
			sortitionPoolAddress,
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
			"failed to attach to BeaconSortitionPool contract: [%v]",
			err,
		)
	}

	return &BeaconChain{
		Chain:         baseChain,
		randomBeacon:  randomBeacon,
		sortitionPool: sortitionPool,
	}, nil
}

// CurrentOperatorToStakingProvider returns the staking provider address for the
// current operator. If the staking provider has not been registered for the
// operator, the returned address is empty and the boolean flag is set to false
// If the staking provider has been registered, the address is not empty and the
// boolean flag indicates true.
func (bc *BeaconChain) CurrentOperatorToStakingProvider() (chain.Address, bool, error) {
	stakingProvider, err := bc.randomBeacon.OperatorToStakingProvider(bc.key.Address)
	if err != nil {
		return "", false, fmt.Errorf(
			"failed to map operator %v to a staking provider: [%v]",
			bc.key.Address,
			err,
		)
	}

	if bytes.Equal(
		stakingProvider.Bytes(),
		bytes.Repeat([]byte{0}, common.AddressLength),
	) {
		return "", false, nil
	}

	return chain.Address(stakingProvider.Hex()), true, nil
}

func (bc *BeaconChain) OperatorToStakingProvider(
	operator chain.Address,
) (chain.Address, error) {
	stakingProvider, err := bc.randomBeacon.OperatorToStakingProvider(
		common.HexToAddress(operator.String()),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to map operator %v to a staking provider: [%v]",
			operator,
			err,
		)
	}
	return chain.Address(stakingProvider.Hex()), nil
}

// EligibleStake returns the current value of the staking provider's eligible
// stake. Eligible stake is defined as the currently authorized stake minus the
// pending authorization decrease. Eligible stake is what is used for operator's
// weight in the sortition pool. If the authorized stake minus the pending
// authorization decrease is below the minimum authorization, eligible stake
// is 0.
func (bc *BeaconChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	eligibleStake, err := bc.randomBeacon.EligibleStake(common.HexToAddress(stakingProvider.String()))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get eligible stake for staking provider %s: [%w]",
			stakingProvider,
			err,
		)
	}

	return eligibleStake, nil
}

// IsPoolLocked returns true if the sortition pool is locked and no state
// changes are allowed.
func (bc *BeaconChain) IsPoolLocked() (bool, error) {
	return bc.sortitionPool.IsLocked()
}

// IsOperatorInPool returns true if the current operator is registered in the
// sortition pool.
func (bc *BeaconChain) IsOperatorInPool() (bool, error) {
	return bc.randomBeacon.IsOperatorInPool(bc.key.Address)
}

// IsOperatorUpToDate checks if the operator's authorized stake is in sync
// with operator's weight in the sortition pool.
// If the operator's authorized stake is not in sync with sortition pool
// weight, function returns false.
// If the operator is not in the sortition pool and their authorized stake
// is non-zero, function returns false.
func (bc *BeaconChain) IsOperatorUpToDate() (bool, error) {
	return bc.randomBeacon.IsOperatorUpToDate(bc.key.Address)
}

// JoinSortitionPool executes a transaction to have the current operator join
// the sortition pool.
func (bc *BeaconChain) JoinSortitionPool() error {
	_, err := bc.randomBeacon.JoinSortitionPool()
	return err
}

// UpdateOperatorStatus executes a transaction to update the current
// operator's state in the sortition pool.
func (bc *BeaconChain) UpdateOperatorStatus() error {
	_, err := bc.randomBeacon.UpdateOperatorStatus(bc.key.Address)
	return err
}

func (bc *BeaconChain) StakeMonitor() (chain.StakeMonitor, error) {
	return newStakeMonitor(bc, bc.Chain), nil
}
