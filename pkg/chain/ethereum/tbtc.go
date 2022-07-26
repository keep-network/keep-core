package ethereum

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/contract"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"math/big"
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

	return &TbtcChain{
		baseChain:          baseChain,
		walletRegistry:     walletRegistry,
		mockWalletRegistry: newMockWalletRegistry(baseChain.blockCounter),
	}, nil
}

// GetConfig returns the expected configuration of the TBTC module.
func (tc *TbtcChain) GetConfig() *tbtc.ChainConfig {
	groupSize := 100
	honestThreshold := 51

	return &tbtc.ChainConfig{
		GroupSize:       groupSize,
		HonestThreshold: honestThreshold,
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

func (tc *TbtcChain) OperatorToStakingProvider() (chain.Address, bool, error) {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) IsPoolLocked() (bool, error) {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) IsOperatorInPool() (bool, error) {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) IsOperatorUpToDate() (bool, error) {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) JoinSortitionPool() error {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) UpdateOperatorStatus() error {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) IsEligibleForRewards() (bool, error) {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) CanRestoreRewardEligibility() (bool, error) {
	//TODO: Implementation.
	panic("not implemented yet")
}

func (tc *TbtcChain) RestoreRewardEligibility() error {
	//TODO: Implementation.
	panic("not implemented yet")
}

// TODO: Implement a real SelectGroup function.
func (tc *TbtcChain) SelectGroup(seed *big.Int) ([]chain.Address, error) {
	_, operatorPublicKey, err := tc.OperatorKeyPair()
	if err != nil {
		return nil, err
	}

	operatorAddress, err := tc.Signing().PublicKeyToAddress(operatorPublicKey)
	if err != nil {
		return nil, err
	}

	groupOperators := make([]chain.Address, tc.GetConfig().GroupSize)
	for i := range groupOperators {
		groupOperators[i] = operatorAddress
	}

	return groupOperators, nil
}

// TODO: Implement a real OnDKGStarted function.
func (tc *TbtcChain) OnDKGStarted(
	handler func(event *tbtc.DKGStartedEvent),
) subscription.EventSubscription {
	return tc.mockWalletRegistry.OnDKGStarted(handler)
}

// TODO: Temporary mock that simulates the behavior of the WalletRegistry
//       contract. Should be removed eventually.
type mockWalletRegistry struct {
	blockCounter chain.BlockCounter
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
				// Generate an event every 500th block.
				if block%500 == 0 {
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
