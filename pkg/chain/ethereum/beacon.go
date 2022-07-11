package ethereum

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/subscription"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen/contract"
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

	chainConfig *beaconchain.Config

	// TODO: Temporary helper map. Should be removed once real-world DKG
	//       result submission is implemented.
	dkgResultSubmissionHandlersMutex sync.Mutex
	dkgResultSubmissionHandlers      map[int]func(submission *event.DKGResultSubmission)
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

	chainConfig, err := fetchChainConfig(randomBeacon)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch the chain config: [%v]",
			err,
		)
	}

	return &BeaconChain{
		Chain:                       baseChain,
		randomBeacon:                randomBeacon,
		sortitionPool:               sortitionPool,
		chainConfig:                 chainConfig,
		dkgResultSubmissionHandlers: make(map[int]func(submission *event.DKGResultSubmission)),
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

// SelectGroup returns the group members for the group generated by
// the given seed. This function can return an error if the beacon chain's
// state does not allow for group selection at the moment.
func (bc *BeaconChain) SelectGroup(seed *big.Int) ([]chain.Address, error) {
	groupSize := big.NewInt(int64(bc.chainConfig.GroupSize))
	seedBytes := [32]byte{}
	seed.FillBytes(seedBytes[:])

	// TODO: Replace with a call to the RandomBeacon.selectGroup function.
	operatorsIDs, err := bc.sortitionPool.SelectGroup(groupSize, seedBytes)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot select group in the sortition pool: [%v]",
			err,
		)
	}

	operatorsAddresses, err := bc.sortitionPool.GetIDOperators(operatorsIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert operators' IDs to addresses: [%v]",
			err,
		)
	}

	result := make([]chain.Address, len(operatorsAddresses))
	for i := range result {
		result[i] = chain.Address(operatorsAddresses[i].String())
	}

	return result, nil
}

func (bc *BeaconChain) OnGroupRegistered(
	handler func(groupRegistration *event.GroupRegistration),
) subscription.EventSubscription {
	// TODO: Implementation.
	panic("not implemented yet")
}

// TODO: Implement a real IsGroupRegistered function.
func (bc *BeaconChain) IsGroupRegistered(groupPublicKey []byte) (bool, error) {
	return false, nil
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

// TODO: Implement a real OnDKGStarted event subscription. The current
//       implementation generates a fake event every 500th block where the
//       seed is the keccak256 of the block number.
func (bc *BeaconChain) OnDKGStarted(
	handler func(event *event.DKGStarted),
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	blocksChan := bc.blockCounter.WatchBlocks(ctx)

	go func() {
		for {
			select {
			case block := <-blocksChan:
				// Generate an event every 500th block.
				if block%500 == 0 {
					// The seed is keccak256(block).
					blockBytes := make([]byte, 8)
					binary.BigEndian.PutUint64(blockBytes, block)
					seedBytes := crypto.Keccak256(blockBytes)
					seed := new(big.Int).SetBytes(seedBytes)

					handler(&event.DKGStarted{
						Seed:        seed,
						BlockNumber: block,
					})
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return subscription.NewEventSubscription(func() {
		cancelCtx()
	})
}

// TODO: Implement a real SubmitDKGResult action. The current implementation
//       just creates and pipes the DKG submission event to the handlers
//       registered in the dkgResultSubmissionHandlers map. Consider getting
//       rid of the result promise in favor of the fire-and-forget style.
func (bc *BeaconChain) SubmitDKGResult(
	participantIndex beaconchain.GroupMemberIndex,
	dkgResult *beaconchain.DKGResult,
	signatures map[beaconchain.GroupMemberIndex][]byte,
) *async.EventDKGResultSubmissionPromise {
	resultPublicationPromise := &async.EventDKGResultSubmissionPromise{}

	failPromise := func(err error) {
		failErr := resultPublicationPromise.Fail(err)
		if failErr != nil {
			logger.Errorf(
				"failed to fail promise for [%v]: [%v]",
				err,
				failErr,
			)
		}
	}

	publishedResult := make(chan *event.DKGResultSubmission)

	subscription := bc.OnDKGResultSubmitted(
		func(onChainEvent *event.DKGResultSubmission) {
			publishedResult <- onChainEvent
		},
	)

	go func() {
		for {
			select {
			case event, success := <-publishedResult:
				// Channel is closed when SubmitDKGResult failed.
				// When this happens, event is nil.
				if !success {
					return
				}

				subscription.Unsubscribe()
				close(publishedResult)

				err := resultPublicationPromise.Fulfill(event)
				if err != nil {
					logger.Errorf(
						"failed to fulfill promise: [%v]",
						err,
					)
				}

				return
			}
		}
	}()

	bc.dkgResultSubmissionHandlersMutex.Lock()
	defer bc.dkgResultSubmissionHandlersMutex.Unlock()

	blockNumber, err := bc.blockCounter.CurrentBlock()
	if err != nil {
		close(publishedResult)
		subscription.Unsubscribe()
		failPromise(err)
		return resultPublicationPromise
	}

	for _, handler := range bc.dkgResultSubmissionHandlers {
		go func(handler func(*event.DKGResultSubmission)) {
			handler(&event.DKGResultSubmission{
				MemberIndex:    uint32(participantIndex),
				GroupPublicKey: dkgResult.GroupPublicKey,
				Misbehaved:     dkgResult.Misbehaved,
				BlockNumber:    blockNumber,
			})
		}(handler)
	}

	return resultPublicationPromise
}

// TODO: Implement a real OnDKGResultSubmitted event subscription. The current
//       implementation just pipes the DKG submission event generated within
//       SubmitDKGResult to the handlers registered in the
//       dkgResultSubmissionHandlers map.
func (bc *BeaconChain) OnDKGResultSubmitted(
	handler func(event *event.DKGResultSubmission),
) subscription.EventSubscription {
	bc.dkgResultSubmissionHandlersMutex.Lock()
	defer bc.dkgResultSubmissionHandlersMutex.Unlock()

	// #nosec G404 (insecure random number source (rand))
	// Temporary test implementation doesn't require secure randomness.
	handlerID := rand.Int()

	bc.dkgResultSubmissionHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		bc.dkgResultSubmissionHandlersMutex.Lock()
		defer bc.dkgResultSubmissionHandlersMutex.Unlock()

		delete(bc.dkgResultSubmissionHandlers, handlerID)
	})
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