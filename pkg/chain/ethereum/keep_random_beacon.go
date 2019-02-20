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

// KeepRandomBeacon connection information for interface to the contract.
type KeepRandomBeacon struct {
	caller            *abi.KeepRandomBeaconImplV1Caller
	callerOptions     *bind.CallOpts
	transactor        *abi.KeepRandomBeaconImplV1Transactor
	transactorOptions *bind.TransactOpts
	contract          *abi.KeepRandomBeaconImplV1
	contractAddress   common.Address
}

// NewKeepRandomBeacon creates the necessary connections and configurations for
// accessing the contract.
func newKeepRandomBeacon(chainConfig *ethereumChain) (*KeepRandomBeacon, error) {
	contractAddressHex, exists := chainConfig.config.ContractAddresses["KeepRandomBeacon"]
	if !exists {
		return nil, fmt.Errorf(
			"no address information for 'KeepRandomBeacon' in configuration",
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

	beaconCaller, err := abi.NewKeepRandomBeaconImplV1Caller(
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

	beaconTransactor, err := abi.NewKeepRandomBeaconImplV1Transactor(
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

	randomBeaconContract, err := abi.NewKeepRandomBeaconImplV1(
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

	return &KeepRandomBeacon{
		caller:            beaconCaller,
		callerOptions:     callerOptions,
		transactor:        beaconTransactor,
		transactorOptions: transactorOptions,
		contract:          randomBeaconContract,
		contractAddress:   contractAddress,
	}, nil
}

// Initialized calls the contract and returns true if the contract has
// had its Initialize method called.
func (krb *KeepRandomBeacon) Initialized() (bool, error) {
	return krb.caller.Initialized(krb.callerOptions)
}

// RequestRelayEntry requests a new entry in the threshold relay.
func (krb *KeepRandomBeacon) RequestRelayEntry(
	blockReward *big.Int,
	rawseed []byte,
) (*types.Transaction, error) {
	seed := big.NewInt(0).SetBytes(rawseed)
	newTransactorOptions := *krb.transactorOptions
	newTransactorOptions.Value = big.NewInt(2)
	return krb.transactor.RequestRelayEntry(&newTransactorOptions, blockReward, seed)
}

// SubmitRelayEntry submits a group signature for consideration.
func (krb *KeepRandomBeacon) SubmitRelayEntry(
	requestID *big.Int,
	groupID *big.Int,
	previousEntry *big.Int,
	groupSignature *big.Int,
	seed *big.Int,
) (*types.Transaction, error) {
	return krb.transactor.RelayEntry(
		krb.transactorOptions,
		requestID,
		groupSignature,
		groupID,
		previousEntry,
		seed,
	)
}

// relayEntryRequestedFunc type of function called for
// RelayEntryRequested event.
type relayEntryRequestedFunc func(
	requestID *big.Int,
	payment *big.Int,
	blockReward *big.Int,
	seed *big.Int,
	blockNumber *big.Int,
)

// WatchRelayEntryRequested watches for event RelayEntryRequested.
func (krb *KeepRandomBeacon) WatchRelayEntryRequested(
	success relayEntryRequestedFunc,
	fail errorCallback,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconImplV1RelayEntryRequested)
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
					event.BlockReward,
					event.Seed,
					event.BlockNumber,
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

// relayEntryGeneratedFunc type of function called for
// RelayEntryGenerated event.
type relayEntryGeneratedFunc func(
	requestID *big.Int,
	requestResponse *big.Int,
	requestGroupID *big.Int,
	previousEntry *big.Int,
	blockNumber *big.Int,
	seed *big.Int,
)

// WatchRelayEntryGenerated watches for event.
func (krb *KeepRandomBeacon) WatchRelayEntryGenerated(
	success relayEntryGeneratedFunc,
	fail errorCallback,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconImplV1RelayEntryGenerated)
	eventSubscription, err := krb.contract.WatchRelayEntryGenerated(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for RelayEntryGenerated event: [%v]",
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
					event.RequestResponse,
					event.RequestGroupID,
					event.PreviousEntry,
					event.BlockNumber,
					event.Seed,
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
