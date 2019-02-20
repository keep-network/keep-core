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

// keepGroup connection information for interface to KeepGroup contract.
type keepGroup struct {
	caller          *abi.KeepGroupImplV1Caller
	callerOpts      *bind.CallOpts
	transactor      *abi.KeepGroupImplV1Transactor
	transactorOpts  *bind.TransactOpts
	contract        *abi.KeepGroupImplV1
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

// NewKeepGroup creates the necessary connections and configurations
// for accessing the KeepGroup contract.
func newKeepGroup(chainConfig *ethereumChain) (*keepGroup, error) {
	contractAddressHex, exists := chainConfig.config.ContractAddresses["KeepGroup"]
	if !exists {
		return nil, fmt.Errorf(
			"no address information for 'KeepGroup' in configuration",
		)
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	groupTransactor, err := abi.NewKeepGroupImplV1Transactor(
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

	groupCaller, err := abi.NewKeepGroupImplV1Caller(contractAddress, chainConfig.client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconCaller contract: [%v]",
			err,
		)
	}

	optsCaller := &bind.CallOpts{
		From: contractAddress,
	}

	groupContract, err := abi.NewKeepGroupImplV1(contractAddress, chainConfig.client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate contract object: %s at address: [%v]",
			err,
			contractAddressHex,
		)
	}

	return &keepGroup{
		transactor:      groupTransactor,
		transactorOpts:  optsTransactor,
		caller:          groupCaller,
		callerOpts:      optsCaller,
		contract:        groupContract,
		contractAddress: contractAddress,
	}, nil
}

// Initialized calls the contract and returns true if the contract
// has had its Initialize method called.
func (kg *keepGroup) Initialized() (bool, error) {
	return kg.caller.Initialized(kg.callerOpts)
}

// GroupThreshold returns the group threshold.  This is the number
// of members that have to report a value to create a new signature.
func (kg *keepGroup) GroupThreshold() (int, error) {
	requiredThresholdMembers, err := kg.caller.GroupThreshold(kg.callerOpts)
	if err != nil {
		return -1, err
	}
	return int(requiredThresholdMembers.Int64()), nil
}

// GroupSize returns the number of members that are required
// to form a group.
func (kg *keepGroup) GroupSize() (int, error) {
	groupSize, err := kg.caller.GroupSize(kg.callerOpts)
	if err != nil {
		return -1, err
	}
	return int(groupSize.Int64()), nil
}

func (kg *keepGroup) TicketInitialSubmissionTimeout() (int, error) {
	ticketInitialSubmissionTimeout, err :=
		kg.caller.TicketInitialSubmissionTimeout(kg.callerOpts)
	if err != nil {
		return -1, err
	}
	return int(ticketInitialSubmissionTimeout.Int64()), nil
}

func (kg *keepGroup) TicketReactiveSubmissionTimeout() (int, error) {
	ticketReactiveSubmissionTimeout, err :=
		kg.caller.TicketReactiveSubmissionTimeout(kg.callerOpts)
	if err != nil {
		return -1, err
	}
	return int(ticketReactiveSubmissionTimeout.Int64()), nil
}

func (kg *keepGroup) TicketChallengeTimeout() (int, error) {
	ticketChallengeTimeout, err := kg.caller.TicketChallengeTimeout(kg.callerOpts)
	if err != nil {
		return -1, err
	}
	return int(ticketChallengeTimeout.Int64()), nil
}

func (kg *keepGroup) MinimumStake() (*big.Int, error) {
	return kg.caller.MinimumStake(kg.callerOpts)
}

func (kg *keepGroup) TokenSupply() (*big.Int, error) {
	return kg.caller.TokenSupply(kg.callerOpts)
}

func (kg *keepGroup) NaturalThreshold() (*big.Int, error) {
	return kg.caller.NaturalThreshold(kg.callerOpts)
}

// HasMinimumStake returns true if the specified address has sufficient
// state to participate.
func (kg *keepGroup) HasMinimumStake(
	address common.Address,
) (bool, error) {
	return kg.caller.HasMinimumStake(kg.callerOpts, address)
}

func (kg *keepGroup) SubmitTicket(
	ticket *relaychain.Ticket,
) (*types.Transaction, error) {
	return kg.transactor.SubmitTicket(
		kg.transactorOpts,
		ticket.Value,
		ticket.Proof.StakerValue,
		ticket.Proof.VirtualStakerIndex,
	)
}

func (kg *keepGroup) SelectedParticipants() ([]common.Address, error) {
	return kg.caller.SelectedParticipants(kg.callerOpts)
}

func (kg *keepGroup) IsDkgResultSubmitted(requestID *big.Int) (bool, error) {
	return kg.caller.IsDkgResultSubmitted(kg.callerOpts, requestID)
}

func (kg *keepGroup) SubmitDKGResult(
	requestID *big.Int,
	result *relaychain.DKGResult,
) (*types.Transaction, error) {
	return kg.transactor.SubmitDkgResult(
		kg.transactorOpts,
		requestID,
		result.Success,
		result.GroupPublicKey,
		result.Disqualified,
		result.Inactive,
	)
}

type dkgResultPublishedEventFunc func(requestID *big.Int, groupPublicKey []byte)

func (kg *keepGroup) WatchDKGResultPublishedEvent(
	success dkgResultPublishedEventFunc,
	fail errorCallback,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepGroupImplV1DkgResultPublishedEvent)
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
				success(event.RequestId, event.GroupPubKey)
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

// SubmitGroupPublicKey upon completion of a signature make the contract
// call to put it on chain.
func (kg *keepGroup) SubmitGroupPublicKey(
	groupPublicKey []byte,
	requestID *big.Int,
) (*types.Transaction, error) {
	return kg.transactor.SubmitGroupPublicKey(kg.transactorOpts, groupPublicKey, requestID)
}

// submitGroupPublicKeyEventFunc type of function called for
// SubmitGroupPublicKeyEvent event.
type submitGroupPublicKeyEventFunc func(
	groupPublicKey []byte,
	requestID *big.Int,
	activationBlockHeight *big.Int,
)

// WatchSubmitGroupPublicKeyEvent watches for event SubmitGroupPublicKeyEvent.
func (kg *keepGroup) WatchSubmitGroupPublicKeyEvent(
	success submitGroupPublicKeyEventFunc,
	fail errorCallback,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepGroupImplV1SubmitGroupPublicKeyEvent)
	eventSubscription, err := kg.contract.WatchSubmitGroupPublicKeyEvent(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return nil, fmt.Errorf(
			"could not create watch for SubmitGroupPublicKeyEvent event: [%v]",
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
					event.GroupPublicKey,
					event.RequestID,
					event.ActivationBlockHeight,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
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
