package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/contract"
)

// Definitions of contract names.
const (
	WalletRegistryContractName = "WalletRegistry"
)

// TbtcChain represents a TBTC-specific chain handle.
type TbtcChain struct {
	*Chain

	walletRegistry *contract.WalletRegistry
}

// NewTbtcChain construct a new instance of the TBTC-specific Ethereum
// chain handle.
func newTbtcChain(
	config ethereum.Config,
	baseChain *Chain,
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
		Chain:          baseChain,
		walletRegistry: walletRegistry,
	}, nil
}

func (tc *TbtcChain) OperatorToStakingProvider(
	operator chain.Address,
) (chain.Address, error) {
	stakingProvider, err := tc.walletRegistry.OperatorToStakingProvider(
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

func (tc *TbtcChain) StakeMonitor() (chain.StakeMonitor, error) {
	return newStakeMonitor(tc, tc.Chain), nil
}
