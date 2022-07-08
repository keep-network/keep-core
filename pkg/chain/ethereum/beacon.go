package ethereum

import (
	"bytes"
	"fmt"
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/subscription"
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
	sortitionPool *contract.SortitionPool

	chainConfig *beaconchain.Config
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
		contract.NewSortitionPool(
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
			"failed to attach to SortitionPool contract: [%v]",
			err,
		)
	}

	chainConfig, err := fetchChainConfig(randomBeacon)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch the chain config: [%v]",
			err,
		)
	}

	return &BeaconChain{
		Chain:         baseChain,
		randomBeacon:  randomBeacon,
		sortitionPool: sortitionPool,
		chainConfig:   chainConfig,
	}, nil
}

// GetConfig returns the expected configuration of the random beacon.
func (bc *BeaconChain) GetConfig() *beaconchain.Config {
	return bc.chainConfig
}

// OperatorToStakingProvider returns the staking provider address for the
// current operator. If the staking provider has not been registered for the
// operator, the returned address is empty and the boolean flag is set to false
// If the staking provider has been registered, the address is not empty and the
// boolean flag indicates true.
func (bc *BeaconChain) OperatorToStakingProvider() (chain.Address, bool, error) {
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

// JoinSortitionPool executes a transaction to have the current operator join
// the sortition pool.
func (bc *BeaconChain) JoinSortitionPool() error {
	_, err := bc.randomBeacon.JoinSortitionPool()
	return err
}

func (bc *BeaconChain) SelectGroup(seed *big.Int) ([]chain.Address, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) OnGroupRegistered(
	handler func(groupRegistration *event.GroupRegistration),
) subscription.EventSubscription {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) IsGroupRegistered(groupPublicKey []byte) (bool, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) IsStaleGroup(groupPublicKey []byte) (bool, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) GetGroupMembers(
	groupPublicKey []byte,
) ([]chain.Address, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) OnDKGStarted(
	handler func(event *event.DKGStarted),
) subscription.EventSubscription {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) SubmitDKGResult(
	participantIndex beaconchain.GroupMemberIndex,
	dkgResult *beaconchain.DKGResult,
	signatures map[beaconchain.GroupMemberIndex][]byte,
) *async.EventDKGResultSubmissionPromise {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) OnDKGResultSubmitted(
	handler func(event *event.DKGResultSubmission),
) subscription.EventSubscription {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) CalculateDKGResultHash(
	dkgResult *beaconchain.DKGResult,
) (beaconchain.DKGResultHash, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) SubmitRelayEntry(entry []byte) *async.EventEntrySubmittedPromise {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) OnRelayEntrySubmitted(
	handler func(entry *event.EntrySubmitted),
) subscription.EventSubscription {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) OnRelayEntryRequested(
	handler func(request *event.Request),
) subscription.EventSubscription {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) ReportRelayEntryTimeout() error {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) IsEntryInProgress() (bool, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) CurrentRequestStartBlock() (*big.Int, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) CurrentRequestPreviousEntry() ([]byte, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

func (bc *BeaconChain) CurrentRequestGroupPublicKey() ([]byte, error) {
	// TODO: Implementation.
	panic("not implemented yet")
}

// fetchChainConfig fetches the on-chain random beacon config.
// TODO: Adjust to the random beacon v2 requirements.
func fetchChainConfig(
	randomBeacon *contract.RandomBeacon,
) (*beaconchain.Config, error) {
	groupSize := 64
	honestThreshold := 33
	resultPublicationBlockStep := 6
	relayEntryTimeout := groupSize * resultPublicationBlockStep

	return &beaconchain.Config{
		GroupSize:                  groupSize,
		HonestThreshold:            honestThreshold,
		ResultPublicationBlockStep: uint64(resultPublicationBlockStep),
		RelayEntryTimeout:          uint64(relayEntryTimeout),
	}, nil
}
