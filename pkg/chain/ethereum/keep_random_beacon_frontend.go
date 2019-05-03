package ethereum

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/keep-network/keep-core/pkg/chain/gen/abi"
	"github.com/keep-network/keep-core/pkg/subscription"
)

// KeepRandomBeaconFrontend connection information for interface to the contract.
type KeepRandomBeaconFrontend struct {
	caller            *abi.KeepRandomBeaconFrontendImplV1Caller
	callerOptions     *bind.CallOpts
	transactor        *abi.KeepRandomBeaconFrontendImplV1Transactor
	transactorOptions *bind.TransactOpts
	contract          *abi.KeepRandomBeaconFrontendImplV1
	contractAddress   common.Address
}

// NewKeepRandomBeaconFrontend creates the necessary connections and configurations for
// accessing the contract.
func newKeepRandomBeaconFrontend(chainConfig *ethereumChain) (*KeepRandomBeaconFrontend, error) {
	contractAddressHex, exists := chainConfig.config.ContractAddresses["KeepRandomBeaconFrontend"]
	if !exists {
		return nil, fmt.Errorf(
			"no address information for 'KeepRandomBeaconFrontend' in configuration",
		)
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	if chainConfig.accountKey == nil {
		key, err := DecryptKeyFile(
			chainConfig.config.Account.KeyFile,
			chainConfig.config.Account.KeyFilePassword,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to read KeyFile: %s: [%v]",
				chainConfig.config.Account.KeyFile,
				err,
			)
		}
		chainConfig.accountKey = key
	}

	beaconCaller, err := abi.NewKeepRandomBeaconFrontendImplV1Caller(
		contractAddress,
		chainConfig.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconCaller contract: [%v]",
			err,
		)
	}

	callerOptions := &bind.CallOpts{
		From: contractAddress,
	}

	beaconTransactor, err := abi.NewKeepRandomBeaconFrontendImplV1Transactor(
		contractAddress,
		chainConfig.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconTranactor contract: [%v]",
			err,
		)
	}

	transactorOptions := bind.NewKeyedTransactor(
		chainConfig.accountKey.PrivateKey,
	)

	randomBeaconFrontendContract, err := abi.NewKeepRandomBeaconFrontendImplV1(
		contractAddress,
		chainConfig.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate contract at address: %s [%v]",
			contractAddressHex,
			err,
		)
	}

	return &KeepRandomBeaconFrontend{
		caller:            beaconCaller,
		callerOptions:     callerOptions,
		transactor:        beaconTransactor,
		transactorOptions: transactorOptions,
		contract:          randomBeaconFrontendContract,
		contractAddress:   contractAddress,
	}, nil
}

// Initialized calls the contract and returns true if the contract has
// had its Initialize method called.
func (krb *KeepRandomBeaconFrontend) Initialized() (bool, error) {
	return krb.caller.Initialized(krb.callerOptions)
}

// RequestRelayEntry requests a new entry in the threshold relay.
func (krb *KeepRandomBeaconFrontend) RequestRelayEntry(
	rawseed []byte,
) (*types.Transaction, error) {
	seed := big.NewInt(0).SetBytes(rawseed)
	newTransactorOptions := *krb.transactorOptions
	newTransactorOptions.Value = big.NewInt(2)
	return krb.transactor.RequestRelayEntry(&newTransactorOptions, seed)
}

// relayEntryRequestedFunc type of function called for
// RelayEntryRequested event.
type relayEntryRequestedFunc func(
	requestID *big.Int,
	payment *big.Int,
	previousEntry *big.Int,
	seed *big.Int,
	groupPublicKey []byte,
	blockNumber uint64,
)

// WatchRelayEntryRequested watches for event RelayEntryRequested.
func (krb *KeepRandomBeaconFrontend) WatchRelayEntryRequested(
	success relayEntryRequestedFunc,
	fail errorCallback,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconFrontendImplV1RelayEntryRequested)
	eventSubscription, err := krb.contract.WatchRelayEntryRequested(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for RelayEntryRequested events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.RequestID,
					event.Payment,
					event.PreviousEntry,
					event.Seed,
					event.GroupPublicKey,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}
