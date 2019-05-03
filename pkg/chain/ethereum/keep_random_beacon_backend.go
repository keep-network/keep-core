package ethereum

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain/gen/abi"
	"github.com/keep-network/keep-core/pkg/subscription"
)

// keepRandomBeaconBackend connection information for interface to
// Keep Random Beacon backend contract.
type keepRandomBeaconBackend struct {
	caller          *abi.KeepRandomBeaconBackendCaller
	callerOpts      *bind.CallOpts
	transactor      *abi.KeepRandomBeaconBackendTransactor
	transactorOpts  *bind.TransactOpts
	contract        *abi.KeepRandomBeaconBackend
	contractAddress common.Address
}

// Important note on watching for Ethereum contract events in Go.
//
// In calls to an abigen generated watch, the first parameter is a filter.
// In example code, to see all events, one must pass in an empty filter.
// In geth, this doesn't work. An empty filter will result in an incorrect
// bloom filter to be selected for the Ethereum search code.
// Rather, to watch for events requires a 'nil' as the first parameter.
//
// For example:
//  	filter := nil
//  	eventSubscription, err := kg.contract.SomeContractSomeEvent(
//			filter,
//			eventChan,
//		)
//
// Will exhibit our desired behavior of selecting an empty filter.
//
// This is different from node.js/web3 code where a 'nil' is treated the same
// as an empty filter.

// NewKeepRandomBeaconBackend creates the necessary connections and configurations
// for accessing the Keep Random Beacon backend contract.
func newKeepRandomBeaconBackend(chainConfig *ethereumChain) (*keepRandomBeaconBackend, error) {
	contractAddressHex, exists := chainConfig.config.ContractAddresses["KeepRandomBeaconBackend"]
	if !exists {
		return nil, fmt.Errorf(
			"no address information for 'KeepRandomBeaconBackend' in configuration",
		)
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	keepRandomBeaconBackendTransactor, err := abi.NewKeepRandomBeaconBackendTransactor(
		contractAddress,
		chainConfig.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconTranactor contract: [%v]",
			err,
		)
	}

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

	optsTransactor := bind.NewKeyedTransactor(
		chainConfig.accountKey.PrivateKey,
	)

	keepRandomBeaconBackendCaller, err := abi.NewKeepRandomBeaconBackendCaller(contractAddress, chainConfig.client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconCaller contract: [%v]",
			err,
		)
	}

	optsCaller := &bind.CallOpts{
		From: contractAddress,
	}

	keepRandomBeaconBackendContract, err := abi.NewKeepRandomBeaconBackend(contractAddress, chainConfig.client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate contract object: %s at address: [%v]",
			err,
			contractAddressHex,
		)
	}

	return &keepRandomBeaconBackend{
		transactor:      keepRandomBeaconBackendTransactor,
		transactorOpts:  optsTransactor,
		caller:          keepRandomBeaconBackendCaller,
		callerOpts:      optsCaller,
		contract:        keepRandomBeaconBackendContract,
		contractAddress: contractAddress,
	}, nil
}

// Initialized calls the contract and returns true if the contract
// has had its Initialize method called.
func (kg *keepRandomBeaconBackend) Initialized() (bool, error) {
	return kg.caller.Initialized(kg.callerOpts)
}

// GroupThreshold returns the group threshold. This is the number
// of members that have to report a value to create a new signature.
func (kg *keepRandomBeaconBackend) GroupThreshold() (int, error) {
	requiredThresholdMembers, err := kg.caller.GroupThreshold(kg.callerOpts)
	if err != nil {
		return -1, err
	}
	return int(requiredThresholdMembers.Int64()), nil
}

// GroupSize returns the number of members that are required
// to form a group.
func (kg *keepRandomBeaconBackend) GroupSize() (int, error) {
	groupSize, err := kg.caller.GroupSize(kg.callerOpts)
	if err != nil {
		return -1, err
	}
	return int(groupSize.Int64()), nil
}

func (kg *keepRandomBeaconBackend) TicketInitialSubmissionTimeout() (*big.Int, error) {
	return kg.caller.TicketInitialSubmissionTimeout(kg.callerOpts)
}

func (kg *keepRandomBeaconBackend) TicketReactiveSubmissionTimeout() (*big.Int, error) {
	return kg.caller.TicketReactiveSubmissionTimeout(kg.callerOpts)
}

func (kg *keepRandomBeaconBackend) TicketChallengeTimeout() (*big.Int, error) {
	return kg.caller.TicketChallengeTimeout(kg.callerOpts)
}

func (kg *keepRandomBeaconBackend) ResultPublicationBlockStep() (*big.Int, error) {
	return kg.caller.ResultPublicationBlockStep(kg.callerOpts)
}

func (kg *keepRandomBeaconBackend) MinimumStake() (*big.Int, error) {
	return kg.caller.MinimumStake(kg.callerOpts)
}

func (kg *keepRandomBeaconBackend) TokenSupply() (*big.Int, error) {
	return kg.caller.TokenSupply(kg.callerOpts)
}

func (kg *keepRandomBeaconBackend) NaturalThreshold() (*big.Int, error) {
	return kg.caller.NaturalThreshold(kg.callerOpts)
}

// HasMinimumStake returns true if the specified address has sufficient
// state to participate.
func (kg *keepRandomBeaconBackend) HasMinimumStake(
	address common.Address,
) (bool, error) {
	return kg.caller.HasMinimumStake(kg.callerOpts, address)
}

func (kg *keepRandomBeaconBackend) SubmitTicket(
	ticket *relaychain.Ticket,
) (*types.Transaction, error) {
	return kg.transactor.SubmitTicket(
		kg.transactorOpts,
		ticket.Value,
		ticket.Proof.StakerValue,
		ticket.Proof.VirtualStakerIndex,
	)
}

func (kg *keepRandomBeaconBackend) SelectedParticipants() ([]common.Address, error) {
	return kg.caller.SelectedParticipants(kg.callerOpts)
}

func (kg *keepRandomBeaconBackend) IsDkgResultSubmitted(requestID *big.Int) (bool, error) {
	return kg.caller.IsDkgResultSubmitted(kg.callerOpts, requestID)
}

func (kg *keepRandomBeaconBackend) SubmitDKGResult(
	requestID *big.Int,
	submitterMemberIndex *big.Int,
	result *relaychain.DKGResult,
	signatures []byte,
	signingMembersIndexes []*big.Int,
) (*types.Transaction, error) {
	return kg.transactor.SubmitDkgResult(
		kg.transactorOpts,
		requestID,
		submitterMemberIndex,
		result.GroupPublicKey,
		result.Disqualified,
		result.Inactive,
		signatures,
		signingMembersIndexes,
	)
}

type dkgResultPublishedEventFunc func(
	requestID *big.Int,
	groupPublicKey []byte,
	blockNumber uint64,
)

func (kg *keepRandomBeaconBackend) WatchDKGResultPublishedEvent(
	success dkgResultPublishedEventFunc,
	fail errorCallback,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconBackendDkgResultPublishedEvent)
	eventSubscription, err := kg.contract.WatchDkgResultPublishedEvent(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return nil, fmt.Errorf(
			"could not create watch for DkgResultPublished event: [%v]",
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
					event.RequestId,
					event.GroupPubKey,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case err := <-eventSubscription.Err():
				fail(err)
				return
			}
		}
	}()

	return subscription.NewEventSubscription(func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}), nil
}

// SubmitRelayEntry submits a group signature for consideration.
func (kg *keepRandomBeaconBackend) SubmitRelayEntry(
	requestID *big.Int,
	groupPubKey []byte,
	previousEntry *big.Int,
	groupSignature *big.Int,
	seed *big.Int,
) (*types.Transaction, error) {
	return kg.transactor.RelayEntry(
		kg.transactorOpts,
		requestID,
		groupSignature,
		groupPubKey,
		previousEntry,
		seed,
	)
}

// relayEntryGeneratedFunc type of function called for
// RelayEntryGenerated event.
type relayEntryGeneratedFunc func(
	requestID *big.Int,
	requestResponse *big.Int,
	requestGroupPubKey []byte,
	previousEntry *big.Int,
	seed *big.Int,
	blockNumber uint64,
)


// WatchRelayEntryGenerated watches for event.
func (kg *keepRandomBeaconBackend) WatchRelayEntryGenerated(
	success relayEntryGeneratedFunc,
	fail errorCallback,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconBackendRelayEntryGenerated)
	eventSubscription, err := kg.contract.WatchRelayEntryGenerated(
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
					event.RequestGroupPubKey,
					event.PreviousEntry,
					event.Seed,
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
func (kg *keepRandomBeaconBackend) WatchRelayEntryRequested(
	success relayEntryRequestedFunc,
	fail errorCallback,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconBackendRelayEntryRequested)
	eventSubscription, err := kg.contract.WatchRelayEntryRequested(
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
