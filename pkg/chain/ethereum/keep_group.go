package ethereum

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	ethereumabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/chain/gen/abi"
	"github.com/keep-network/keep-core/pkg/subscription"
)

// keepGroup connection information for interface to KeepGroup contract.
type keepGroup struct {
	caller          *abi.KeepGroupImplV1Caller
	callerOpts      *bind.CallOpts
	errorResolver   *ethutil.ErrorResolver
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
		key, err := ethutil.DecryptKeyFile(
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

	contractAbi, err := ethereumabi.JSON(strings.NewReader(abi.KeepGroupImplV1ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &keepGroup{
		transactor:      groupTransactor,
		transactorOpts:  optsTransactor,
		caller:          groupCaller,
		callerOpts:      optsCaller,
		errorResolver:   ethutil.NewErrorResolver(chainConfig.client, &contractAbi, &contractAddress),
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

func (kg *keepGroup) TicketInitialSubmissionTimeout() (*big.Int, error) {
	return kg.caller.TicketInitialSubmissionTimeout(kg.callerOpts)
}

func (kg *keepGroup) TicketReactiveSubmissionTimeout() (*big.Int, error) {
	return kg.caller.TicketReactiveSubmissionTimeout(kg.callerOpts)
}

func (kg *keepGroup) TicketChallengeTimeout() (*big.Int, error) {
	return kg.caller.TicketChallengeTimeout(kg.callerOpts)
}

func (kg *keepGroup) ResultPublicationBlockStep() (*big.Int, error) {
	return kg.caller.ResultPublicationBlockStep(kg.callerOpts)
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
	transaction, err := kg.transactor.SubmitTicket(
		kg.transactorOpts,
		ticket.Value,
		ticket.Proof.StakerValue,
		ticket.Proof.VirtualStakerIndex,
	)

	if err != nil {
		return transaction, kg.errorResolver.ResolveError(
			err,
			nil,
			"submitTicket",
			ticket.Value,
			ticket.Proof.StakerValue,
			ticket.Proof.VirtualStakerIndex,
		)
	}

	return transaction, err
}

func (kg *keepGroup) SelectedParticipants() ([]common.Address, error) {
	return kg.caller.SelectedParticipants(kg.callerOpts)
}

func (kg *keepGroup) IsDkgResultSubmitted(requestID *big.Int) (bool, error) {
	return kg.caller.IsDkgResultSubmitted(kg.callerOpts, requestID)
}

// Checks if a group with the given public key is registered on-chain.
func (kg *keepGroup) IsGroupRegistered(groupPublicKey []byte) (bool, error) {
	return kg.caller.IsGroupRegistered(kg.callerOpts, groupPublicKey)
}

func (kg *keepGroup) SubmitDKGResult(
	requestID *big.Int,
	submitterMemberIndex *big.Int,
	result *relaychain.DKGResult,
	signatures []byte,
	signingMembersIndexes []*big.Int,
) (*types.Transaction, error) {
	transaction, err := kg.transactor.SubmitDkgResult(
		kg.transactorOpts,
		requestID,
		submitterMemberIndex,
		result.GroupPublicKey,
		result.Disqualified,
		result.Inactive,
		signatures,
		signingMembersIndexes,
	)

	if err != nil {
		return transaction, kg.errorResolver.ResolveError(
			err,
			nil,
			"submitDkgResult",
			requestID,
			submitterMemberIndex,
			result.GroupPublicKey,
			result.Disqualified,
			result.Inactive,
			signatures,
			signingMembersIndexes,
		)
	}

	return transaction, err
}

type dkgResultPublishedEventFunc func(
	requestID *big.Int,
	groupPublicKey []byte,
	blockNumber uint64,
)

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
