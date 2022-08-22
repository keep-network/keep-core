package ethereum

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/contract"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

// Definitions of contract names.
const (
	WalletRegistryContractName = "WalletRegistry"
)

// TbtcChain represents a TBTC-specific chain handle.
type TbtcChain struct {
	*baseChain

	walletRegistry *contract.WalletRegistry

	mockWalletRegistry *mockWalletRegistry
	sortitionPool      *contract.EcdsaSortitionPool
}

// NewTbtcChain construct a new instance of the TBTC-specific Ethereum
// chain handle.
func newTbtcChain(
	config ethereum.Config,
	baseChain *baseChain,
) (*TbtcChain, error) {
	walletRegistryAddress, err := config.ContractAddress(WalletRegistryContractName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve %s contract address: [%v]",
			WalletRegistryContractName,
			err,
		)
	}

	walletRegistry, err :=
		contract.NewWalletRegistry(
			walletRegistryAddress,
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
			"failed to attach to WalletRegistry contract: [%v]",
			err,
		)
	}

	sortitionPoolAddress, err := walletRegistry.SortitionPool()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get sortition pool address: [%v]",
			err,
		)
	}

	sortitionPool, err :=
		contract.NewEcdsaSortitionPool(
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
			"failed to attach to EcdsaSortitionPool contract: [%v]",
			err,
		)
	}

	return &TbtcChain{
		baseChain:          baseChain,
		walletRegistry:     walletRegistry,
		mockWalletRegistry: newMockWalletRegistry(baseChain.blockCounter),
		sortitionPool:      sortitionPool,
	}, nil
}

// GetConfig returns the expected configuration of the TBTC module.
func (tc *TbtcChain) GetConfig() *tbtc.ChainConfig {
	groupSize := 100
	honestThreshold := 51
	resultPublicationBlockStep := 1

	return &tbtc.ChainConfig{
		GroupSize:                  groupSize,
		HonestThreshold:            honestThreshold,
		ResultPublicationBlockStep: uint64(resultPublicationBlockStep),
	}
}

// IsRecognized checks whether the given operator is recognized by the TbtcChain
// as eligible to join the network. If the operator has a stake delegation or
// had a stake delegation in the past, it will be recognized.
func (tc *TbtcChain) IsRecognized(operatorPublicKey *operator.PublicKey) (bool, error) {
	operatorAddress, err := operatorPublicKeyToChainAddress(operatorPublicKey)
	if err != nil {
		return false, fmt.Errorf(
			"cannot convert from operator key to chain address: [%v]",
			err,
		)
	}

	stakingProvider, err := tc.walletRegistry.OperatorToStakingProvider(
		operatorAddress,
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to map operator [%v] to a staking provider: [%v]",
			operatorAddress,
			err,
		)
	}

	if (stakingProvider == common.Address{}) {
		return false, nil
	}

	// Check if the staking provider has an owner. This check ensures that there
	// is/was a stake delegation for the given staking provider.
	_, _, _, hasStakeDelegation, err := tc.baseChain.RolesOf(
		chain.Address(stakingProvider.Hex()),
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check stake delegation for staking provider [%v]: [%v]",
			stakingProvider,
			err,
		)
	}

	if !hasStakeDelegation {
		return false, nil
	}

	return true, nil
}

// OperatorToStakingProvider returns the staking provider address for the
// operator. If the staking provider has not been registered for the
// operator, the returned address is empty and the boolean flag is set to
// false. If the staking provider has been registered, the address is not
// empty and the boolean flag indicates true.
func (tc *TbtcChain) OperatorToStakingProvider() (chain.Address, bool, error) {
	stakingProvider, err := tc.walletRegistry.OperatorToStakingProvider(tc.key.Address)
	if err != nil {
		return "", false, fmt.Errorf(
			"failed to map operator [%v] to a staking provider: [%v]",
			tc.key.Address,
			err,
		)
	}

	if (stakingProvider == common.Address{}) {
		return "", false, nil
	}

	return chain.Address(stakingProvider.Hex()), true, nil
}

// EligibleStake returns the current value of the staking provider's
// eligible stake. Eligible stake is defined as the currently authorized
// stake minus the pending authorization decrease. Eligible stake
// is what is used for operator's weight in the sortition pool.
// If the authorized stake minus the pending authorization decrease
// is below the minimum authorization, eligible stake is 0.
func (tc *TbtcChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	eligibleStake, err := tc.walletRegistry.EligibleStake(
		common.HexToAddress(stakingProvider.String()),
	)
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
func (tc *TbtcChain) IsPoolLocked() (bool, error) {
	return tc.sortitionPool.IsLocked()
}

// IsOperatorInPool returns true if the operator is registered in
// the sortition pool.
func (tc *TbtcChain) IsOperatorInPool() (bool, error) {
	return tc.walletRegistry.IsOperatorInPool(tc.key.Address)
}

// IsOperatorUpToDate checks if the operator's authorized stake is in sync
// with operator's weight in the sortition pool.
// If the operator's authorized stake is not in sync with sortition pool
// weight, function returns false.
// If the operator is not in the sortition pool and their authorized stake
// is non-zero, function returns false.
func (tc *TbtcChain) IsOperatorUpToDate() (bool, error) {
	return tc.walletRegistry.IsOperatorUpToDate(tc.key.Address)
}

// JoinSortitionPool executes a transaction to have the operator join the
// sortition pool.
func (tc *TbtcChain) JoinSortitionPool() error {
	_, err := tc.walletRegistry.JoinSortitionPool()
	return err
}

// UpdateOperatorStatus executes a transaction to update the operator's
// state in the sortition pool.
func (tc *TbtcChain) UpdateOperatorStatus() error {
	_, err := tc.walletRegistry.UpdateOperatorStatus(tc.key.Address)
	return err
}

// IsEligibleForRewards checks whether the operator is eligible for rewards
// or not.
func (tc *TbtcChain) IsEligibleForRewards() (bool, error) {
	return tc.sortitionPool.IsEligibleForRewards(tc.key.Address)
}

// Checks whether the operator is able to restore their eligibility for
// rewards right away.
func (tc *TbtcChain) CanRestoreRewardEligibility() (bool, error) {
	return tc.sortitionPool.CanRestoreRewardEligibility(tc.key.Address)
}

// Restores reward eligibility for the operator.
func (tc *TbtcChain) RestoreRewardEligibility() error {
	_, err := tc.sortitionPool.RestoreRewardEligibility(tc.key.Address)
	return err
}

// SelectGroup returns the group members for the group generated by
// the given seed. This function can return an error if the beacon chain's
// state does not allow for group selection at the moment.
func (tc *TbtcChain) SelectGroup(seed *big.Int) ([]chain.Address, error) {
	groupSize := big.NewInt(int64(tc.GetConfig().GroupSize))
	seedBytes := [32]byte{}
	seed.FillBytes(seedBytes[:])

	// TODO: Replace with a call to the WalletRegistry.selectGroup function.
	operatorsIDs, err := tc.sortitionPool.SelectGroup(groupSize, seedBytes)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot select group in the sortition pool: [%v]",
			err,
		)
	}

	operatorsAddresses, err := tc.sortitionPool.GetIDOperators(operatorsIDs)
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

// TODO: Implement a real OnDKGStarted function.
func (tc *TbtcChain) OnDKGStarted(
	handler func(event *tbtc.DKGStartedEvent),
) subscription.EventSubscription {
	return tc.mockWalletRegistry.OnDKGStarted(handler)
}

// TODO: Implement a real OnDKGResultSubmitted event subscription. The current
//       implementation just pipes the DKG submission event generated within
//       SubmitDKGResult to the handlers registered in the
//       dkgResultSubmissionHandlers map.
func (tc *TbtcChain) OnDKGResultSubmitted(
	handler func(event *tbtc.DKGResultSubmittedEvent),
) subscription.EventSubscription {
	return tc.mockWalletRegistry.OnDKGResultSubmitted(handler)
}

// SignResult signs the provided DKG result. It returns the information
// pertaining to the signing process: public key, signature, result hash.
func (tc *TbtcChain) SignResult(result *dkg.Result) (*dkg.SignedResult, error) {
	resultHash, err := tc.calculateDKGResultHash(result)
	if err != nil {
		return nil, fmt.Errorf(
			"dkg result hash calculation failed [%v]",
			err,
		)
	}

	signature, err := tc.Signing().Sign(resultHash[:])
	if err != nil {
		return nil, fmt.Errorf(
			"dkg result hash signing failed [%v]",
			err,
		)
	}

	return &dkg.SignedResult{
		PublicKey:  tc.Signing().PublicKey(),
		Signature:  signature,
		ResultHash: resultHash,
	}, nil
}

// VerifySignature verifies if the signature was generated from the provided
// DKG result has using the provided public key.
func (tc *TbtcChain) VerifySignature(signedResult *dkg.SignedResult) (bool, error) {
	return tc.Signing().VerifyWithPublicKey(
		signedResult.ResultHash[:],
		signedResult.Signature,
		signedResult.PublicKey,
	)
}

// SubmitResult submits the DKG result to the chain, along with signatures
// over result hash from group participants supporting the result.
func (tc *TbtcChain) SubmitResult(
	result *dkg.Result,
	signatures map[group.MemberIndex][]byte,
	startBlockNumber uint64,
	memberIndex group.MemberIndex,
) error {
	config := tc.GetConfig()

	// Chain rejects the result if it has less than 25% safety margin.
	// If there are not enough signatures to preserve the margin, it does not
	// make sense to submit the result.
	signatureThreshold := config.HonestThreshold +
		(config.GroupSize-config.HonestThreshold)/2
	if len(signatures) < signatureThreshold {
		return fmt.Errorf(
			"could not submit result with [%v] signatures for signature threshold [%v]",
			len(signatures),
			signatureThreshold,
		)
	}

	onSubmittedResultChan := make(chan uint64)

	subscription := tc.OnDKGResultSubmitted(
		func(event *tbtc.DKGResultSubmittedEvent) {
			onSubmittedResultChan <- event.BlockNumber
		},
	)

	returnWithError := func(err error) error {
		subscription.Unsubscribe()
		close(onSubmittedResultChan)
		return err
	}

	groupPublicKeyBytes, err := result.GroupPublicKeyBytes()
	if err != nil {
		return returnWithError(
			fmt.Errorf(
				"could not extract public key bytes from the result: [%v]",
				err,
			),
		)
	}

	alreadySubmitted, err := tc.isDKGResultSubmitted(groupPublicKeyBytes)
	if err != nil {
		return returnWithError(
			fmt.Errorf(
				"could not check if the result is already submitted: [%v]",
				err,
			),
		)
	}

	// Someone who was ahead of us in the queue submitted the result. Giving up.
	if alreadySubmitted {
		return returnWithError(nil)
	}

	// Wait until the current member is eligible to submit the result.
	eligibleToSubmitWaiter, err := tc.waitForSubmissionEligibility(
		startBlockNumber,
		memberIndex,
	)
	if err != nil {
		return returnWithError(
			fmt.Errorf("wait for eligibility failure: [%v]", err),
		)
	}

	for {
		select {
		case /*blockNumber :=*/ <-eligibleToSubmitWaiter:
			// Member becomes eligible to submit the result.
			subscription.Unsubscribe()
			close(onSubmittedResultChan)

			// TODO: What to do with logging?
			// sm.logger.Infof(
			// 	"[member:%v] submitting DKG result with public key [0x%x] and "+
			// 		"[%v] supporting member signatures at block [%v]",
			// 	sm.memberIndex,
			// 	result.GroupPublicKeyBytes,
			// 	len(signatures),
			// 	blockNumber,
			// )

			return tc.submitDKGResult(
				memberIndex,
				result,
				signatures,
			)
		case /*blockNumber :=*/ <-onSubmittedResultChan:
			// TODO: What to do with logging
			// sm.logger.Infof(
			// 	"[member:%v] leaving; DKG result submitted by other member at block [%v]",
			// 	sm.memberIndex,
			// 	blockNumber,
			// )
			// A result has been submitted by other member. Leave without
			// publishing the result.
			return returnWithError(nil)
		}
	}
}

// TODO: Implement a real isDKGResultSubmitted function.
// TODO: Check what the argument should be.
func (tc *TbtcChain) isDKGResultSubmitted(groupPublicKey []byte) (bool, error) {
	return false, nil
}

// calculateDKGResultHash calculates Keccak-256 hash of the DKG result. Operation
// is performed off-chain.
//
// It first encodes the result using solidity ABI and then calculates Keccak-256
// hash over it. This corresponds to the DKG result hash calculation on-chain.
// Hashes calculated off-chain and on-chain must always match.
func (tc *TbtcChain) calculateDKGResultHash(
	result *dkg.Result,
) (dkg.ResultHash, error) {
	groupPublicKeyBytes, err := result.GroupPublicKeyBytes()
	if err != nil {
		return dkg.ResultHash{}, err
	}

	// Encode DKG result to the format matched with Solidity keccak256(abi.encodePacked(...))
	hash := crypto.Keccak256(groupPublicKeyBytes, result.MisbehavedMembersBytes())
	return dkg.ResultHashFromBytes(hash)
}

// TODO: Implement a real submitDKGResult action. The current implementation
//       just creates and pipes the DKG submission event to the handlers
//       registered in the dkgResultSubmissionHandlers map.
func (tc *TbtcChain) submitDKGResult(
	memberIndex group.MemberIndex,
	result *dkg.Result,
	signatures map[group.MemberIndex][]byte,
) error {
	return tc.mockWalletRegistry.SubmitDKGResult(
		memberIndex,
		result,
		signatures,
	)
}

// TODO: Temporary mock that simulates the behavior of the WalletRegistry
//       contract. Should be removed eventually.
type mockWalletRegistry struct {
	blockCounter chain.BlockCounter

	dkgResultSubmissionHandlersMutex sync.Mutex
	dkgResultSubmissionHandlers      map[int]func(submission *tbtc.DKGResultSubmittedEvent)

	currentDkgMutex      sync.RWMutex
	currentDkgStartBlock *big.Int

	activeGroupMutex         sync.RWMutex
	activeGroup              []byte
	activeGroupOperableBlock *big.Int
}

func newMockWalletRegistry(blockCounter chain.BlockCounter) *mockWalletRegistry {
	return &mockWalletRegistry{blockCounter: blockCounter}
}

func (mwr *mockWalletRegistry) OnDKGStarted(
	handler func(event *tbtc.DKGStartedEvent),
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	blocksChan := mwr.blockCounter.WatchBlocks(ctx)

	go func() {
		for {
			select {
			case block := <-blocksChan:
				// Generate an event every 500th block starting from block 250.
				// The shift is done in order to avoid overlapping with beacon
				// DKG test loop.
				shift := uint64(250)
				if block >= shift && (block-shift)%500 == 0 {
					// The seed is keccak256(block).
					blockBytes := make([]byte, 8)
					binary.BigEndian.PutUint64(blockBytes, block)
					seedBytes := crypto.Keccak256(blockBytes)
					seed := new(big.Int).SetBytes(seedBytes)

					go handler(&tbtc.DKGStartedEvent{
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

func (mwr *mockWalletRegistry) OnDKGResultSubmitted(
	handler func(event *tbtc.DKGResultSubmittedEvent),
) subscription.EventSubscription {
	mwr.dkgResultSubmissionHandlersMutex.Lock()
	defer mwr.dkgResultSubmissionHandlersMutex.Unlock()

	// #nosec G404 (insecure random number source (rand))
	// Temporary test implementation doesn't require secure randomness.
	handlerID := rand.Int()

	mwr.dkgResultSubmissionHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		mwr.dkgResultSubmissionHandlersMutex.Lock()
		defer mwr.dkgResultSubmissionHandlersMutex.Unlock()

		delete(mwr.dkgResultSubmissionHandlers, handlerID)
	})
}

func (mwr *mockWalletRegistry) SubmitDKGResult(
	memberIndex group.MemberIndex,
	result *dkg.Result,
	signatures map[group.MemberIndex][]byte,
) error {
	mwr.dkgResultSubmissionHandlersMutex.Lock()
	defer mwr.dkgResultSubmissionHandlersMutex.Unlock()

	mwr.currentDkgMutex.Lock()
	defer mwr.currentDkgMutex.Unlock()

	mwr.activeGroupMutex.Lock()
	defer mwr.activeGroupMutex.Unlock()

	// Abort if there is no DKG in progress. This check is needed to handle a
	// situation in which two operators of the same client attempt to submit
	// the DKG result.
	if mwr.currentDkgStartBlock == nil {
		return nil
	}

	blockNumber, err := mwr.blockCounter.CurrentBlock()
	if err != nil {
		return fmt.Errorf("failed to get the current block")
	}

	groupPublicKeyBytes, err := result.GroupPublicKeyBytes()
	if err != nil {
		return fmt.Errorf(
			"failed to extract group public key bytes from the result [%v]",
			err,
		)
	}

	for _, handler := range mwr.dkgResultSubmissionHandlers {
		go func(handler func(*tbtc.DKGResultSubmittedEvent)) {
			handler(&tbtc.DKGResultSubmittedEvent{
				MemberIndex:         uint32(memberIndex),
				GroupPublicKeyBytes: groupPublicKeyBytes,
				Misbehaved:          result.MisbehavedMembersBytes(),
				BlockNumber:         blockNumber,
			})
		}(handler)
	}

	mwr.activeGroup = groupPublicKeyBytes
	mwr.activeGroupOperableBlock = new(big.Int).Add(
		mwr.currentDkgStartBlock,
		big.NewInt(150),
	)
	mwr.currentDkgStartBlock = nil

	return nil
}

// waitForSubmissionEligibility waits until the current member is eligible to
// submit a result to the blockchain. First member is eligible to submit straight
// away, each following member is eligible after pre-defined block step.
func (tc *TbtcChain) waitForSubmissionEligibility(
	startBlockNumber uint64,
	memberIndex group.MemberIndex,
) (<-chan uint64, error) {
	// T_init + (member_index - 1) * T_step
	blockWaitTime := (uint64(memberIndex) - 1) *
		tc.GetConfig().ResultPublicationBlockStep

	eligibleBlockHeight := startBlockNumber + blockWaitTime

	// TODO: What should we do about the logging? Do not log?
	//       Pass the logger? Return the logging messages?

	// sm.logger.Infof(
	// 	"[member:%v] waiting for block [%v] to submit",
	// 	participantIndex,
	// 	eligibleBlockHeight,
	// )

	blockCounter, err := tc.BlockCounter()
	if err != nil {
		return nil, fmt.Errorf("could not get block counter [%v]", err)
	}

	waiter, err := blockCounter.BlockHeightWaiter(eligibleBlockHeight)
	if err != nil {
		return nil, fmt.Errorf("block height waiter failure [%v]", err)
	}

	return waiter, err
}
