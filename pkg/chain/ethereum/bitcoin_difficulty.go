package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/contract"
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
	LightRelayMaintainerProxy *contract.LightRelayMaintainerProxy
}

// NewBitcoinDifficultyChain construct a new instance of the Bitcoin difficulty
// - specific Ethereum chain handle.
func NewBitcoinDifficultyChain(
	config ethereum.Config,
	baseChain *baseChain,
) (*BitcoinDifficultyChain, error) {
	lightRelayAddress, err := config.ContractAddress(LightRelayContractName)
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

	lightRelayMaintainerProxyAddress, err := config.ContractAddress(
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
		LightRelayMaintainerProxy: lightRelayMaintainerProxy,
	}, nil
}

// Ready checks whether the relay is active (i.e. genesis has been performed).
// Note that if the relay is used by querying the current and previous epoch
// difficulty, at least one retarget needs to be provided after genesis;
// otherwise the prevEpochDifficulty will be uninitialised and zero.
func (bdc *BitcoinDifficultyChain) Ready() (bool, error) {
	return bdc.lightRelay.Ready()
}

// AuthorizationRequired checks whether the relay requires the address
// submitting a retarget to be authorised in advance by governance.
func (bdc *BitcoinDifficultyChain) AuthorizationRequired() (bool, error) {
	return bdc.lightRelay.AuthorizationRequired()
}

// IsAuthorized checks whether the given address has been authorised by
// governance to submit a retarget.
func (bdc *BitcoinDifficultyChain) IsAuthorized(address chain.Address) (bool, error) {
	return bdc.lightRelay.IsAuthorized(common.HexToAddress(address.String()))
}

// Retarget adds a new epoch to the relay by providing a proof of the difficulty
// before and after the retarget.
func (bdc *BitcoinDifficultyChain) Retarget(headers []*bitcoin.BlockHeader) error {
	var serializedHeaders []byte
	for _, header := range headers {
		serializedHeader := header.Serialize()
		serializedHeaders = append(serializedHeaders, serializedHeader[:]...)
	}
	_, err := bdc.LightRelayMaintainerProxy.Retarget(serializedHeaders)
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
