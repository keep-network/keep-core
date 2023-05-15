package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/contract"
	"github.com/keep-network/keep-core/pkg/maintainer"
)

// Definitions of contract names.
const (
	LightRelayContractName                = "LightRelay"
	LightRelayMaintainerProxyContractName = "LightRelayMaintainerProxy"
)

// BitcoinDifficultyChain represents a Bitcoin difficulty-specific chain handle.
type BitcoinDifficultyChain struct {
	*baseChain

	lightRelay                *contract.LightRelay
	lightRelayMaintainerProxy *contract.LightRelayMaintainerProxy
}

// NewBitcoinDifficultyChain construct a new instance of the Bitcoin difficulty
// - specific Ethereum chain handle.
func NewBitcoinDifficultyChain(
	ethereumConfig ethereum.Config,
	maintainerConfig maintainer.Config,
	baseChain *baseChain,
) (*BitcoinDifficultyChain, error) {
	lightRelayAddress, err := ethereumConfig.ContractAddress(
		LightRelayContractName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to LightRelay contract: [%v]",
			err,
		)
	}

	lightRelay, err :=
		contract.NewLightRelay(
			lightRelayAddress,
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
			"failed to attach to LightRelay contract: [%v]",
			err,
		)
	}

	// If the Bitcoin difficulty should be updated directly via LightRelay,
	// quit early without creating a handle to LightRelayMaintainerProxy.
	if maintainerConfig.DisableBitcoinDifficultyProxy {
		return &BitcoinDifficultyChain{
			baseChain:                 baseChain,
			lightRelay:                lightRelay,
			lightRelayMaintainerProxy: nil,
		}, nil
	}

	// The Bitcoin difficulty should be updated via LightRelayMaintainerProxy.
	lightRelayMaintainerProxyAddress, err := ethereumConfig.ContractAddress(
		LightRelayMaintainerProxyContractName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to LightRelayMaintainerProxy contract: [%v]",
			err,
		)
	}

	lightRelayMaintainerProxy, err :=
		contract.NewLightRelayMaintainerProxy(
			lightRelayMaintainerProxyAddress,
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
			"failed to attach to LightRelayMaintainerProxy contract: [%v]",
			err,
		)
	}

	retrievedLightRelayAddress, err := lightRelayMaintainerProxy.LightRelay()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to retrieve the relay address from LightRelayMaintainerProxy "+
				"contract: [%v]",
			err,
		)
	}

	// Verify the LightRelay set in LightRelayMaintainerProxy is the same
	// instance as the one set in the client via configuration.
	if lightRelayAddress != retrievedLightRelayAddress {
		return nil, fmt.Errorf("mismatch between LightRelay addresses")
	}

	return &BitcoinDifficultyChain{
		baseChain:                 baseChain,
		lightRelay:                lightRelay,
		lightRelayMaintainerProxy: lightRelayMaintainerProxy,
	}, nil
}

// Ready checks whether the relay is active (i.e. genesis has been performed).
// Note that if the relay is used by querying the current and previous epoch
// difficulty, at least one retarget needs to be provided after genesis;
// otherwise the prevEpochDifficulty will be uninitialised and zero.
func (bdc *BitcoinDifficultyChain) Ready() (bool, error) {
	return bdc.lightRelay.Ready()
}

// IsAuthorized checks whether the given address has been authorized to submit
// a retarget directly to LightRelay. This function should be used when
// retargetting via LightRelayMaintainerProxy is disabled.
func (bdc *BitcoinDifficultyChain) IsAuthorized(address chain.Address) (bool, error) {
	authorizationRequired, err := bdc.lightRelay.AuthorizationRequired()
	if err != nil {
		return false, fmt.Errorf(
			"cannot check whether authorization is required to submit "+
				"block headers: [%w]",
			err,
		)
	}

	if !authorizationRequired {
		return true, nil
	}

	return bdc.lightRelay.IsAuthorized(
		common.HexToAddress(address.String()),
	)
}

// IsAuthorizedForRefund checks whether the given address has been authorized to
// submit a retarget via LightRelayMaintainerProxy. This function should be used
// when retargetting via LightRelayMaintainerProxy is not disabled.
func (bdc *BitcoinDifficultyChain) IsAuthorizedForRefund(address chain.Address) (bool, error) {
	return bdc.lightRelayMaintainerProxy.IsAuthorized(
		common.HexToAddress(address.String()),
	)
}

// Retarget adds a new epoch to the relay by providing a proof of the difficulty
// before and after the retarget. The cost of calling this function is not
// refunded to the caller.
func (bdc *BitcoinDifficultyChain) Retarget(headers []*bitcoin.BlockHeader) error {
	var serializedHeaders []byte
	for _, header := range headers {
		serializedHeader := header.Serialize()
		serializedHeaders = append(serializedHeaders, serializedHeader[:]...)
	}

	// Update Bitcoin difficulty directly via LightRelay.
	_, err := bdc.lightRelay.Retarget(serializedHeaders)
	return err
}

// RetargetWithRefund adds a new epoch to the relay by providing a proof of the
// difficulty before and after the retarget. The cost of calling this function
// is refunded to the caller.
func (bdc *BitcoinDifficultyChain) RetargetWithRefund(headers []*bitcoin.BlockHeader) error {
	var serializedHeaders []byte
	for _, header := range headers {
		serializedHeader := header.Serialize()
		serializedHeaders = append(serializedHeaders, serializedHeader[:]...)
	}

	gasEstimate, err := bdc.lightRelayMaintainerProxy.RetargetGasEstimate(
		serializedHeaders,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to estimate gas for retarget with refund: [%w]",
			err,
		)
	}

	// Add 20% to the gas estimate as the transaction tends to fail with the
	// original gas estimate.
	gasEstimateWithMargin := float64(gasEstimate) * float64(1.2)

	// Update Bitcoin difficulty via LightRelayMaintainerProxy.
	_, err = bdc.lightRelayMaintainerProxy.Retarget(
		serializedHeaders,
		ethutil.TransactionOptions{
			GasLimit: uint64(gasEstimateWithMargin),
		},
	)
	return err
}

// CurrentEpoch returns the number of the latest difficulty epoch which is
// proven to the relay. If the genesis epoch's number is set correctly, and
// retargets along the way have been legitimate, this equals the height of
// the block starting the most recent epoch, divided by 2016.
func (bdc *BitcoinDifficultyChain) CurrentEpoch() (uint64, error) {
	return bdc.lightRelay.CurrentEpoch()
}

// ProofLength returns the number of blocks required for each side of a retarget
// proof.
func (bdc *BitcoinDifficultyChain) ProofLength() (uint64, error) {
	return bdc.lightRelay.ProofLength()
}
