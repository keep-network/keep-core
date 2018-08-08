package ethereum

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/keep-network/keep-core/pkg/chain/gen"
)

// keepGroup connection information for interface to KeepGroup contract.
type keepGroup struct {
	caller          *gen.KeepGroupImplV1Caller
	callerOpts      *bind.CallOpts
	transactor      *gen.KeepGroupImplV1Transactor
	transactorOpts  *bind.TransactOpts
	contract        *gen.KeepGroupImplV1
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
func newKeepGroup(pv *ethereumChain) (*keepGroup, error) {
	contractAddressHex, exists := pv.config.ContractAddresses["KeepGroupImplV1"]
	if !exists {
		return nil, fmt.Errorf(
			"no address information for 'KeepGroup' in configuration",
		)
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	groupTransactor, err := gen.NewKeepGroupImplV1Transactor(
		contractAddress,
		pv.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconTranactor contract: [%v]",
			err,
		)
	}

	file, err := os.Open(pv.config.Account.KeyFile)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open keyfile %s: [%v]",
			pv.config.Account.KeyFile,
			err,
		)
	}

	optsTransactor, err := bind.NewTransactor(
		bufio.NewReader(file),
		pv.config.Account.KeyFilePassword,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read keyfile: %s: [%v]",
			pv.config.Account.KeyFile,
			err,
		)
	}

	groupCaller, err := gen.NewKeepGroupImplV1Caller(contractAddress, pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconCaller contract: [%v]",
			err,
		)
	}

	optsCaller := &bind.CallOpts{
		From: contractAddress,
	}

	groupContract, err := gen.NewKeepGroupImplV1(contractAddress, pv.client)
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

// DissolveGroup breaks up the group that is associated with the public key.
func (kg *keepGroup) DissolveGroup(
	groupPubKey []byte,
) (*types.Transaction, error) {
	groupPubKeyArray, err := toByte32(groupPubKey)
	if err != nil {
		return nil, err
	}
	return kg.transactor.DissolveGroup(kg.transactorOpts, groupPubKeyArray)
}

// CreateGroup starts a new group with the specified public key.
func (kg *keepGroup) CreateGroup(
	groupPubKey []byte,
) (*types.Transaction, error) {
	groupPubKeyArray, err := toByte32(groupPubKey)
	if err != nil {
		return nil, err
	}
	return kg.transactor.CreateGroup(kg.transactorOpts, groupPubKeyArray)
}

// GroupThreshold returns the group threshold.  This is the number
// of members that have to report a value to create a new signature.
func (kg *keepGroup) GroupThreshold() (int, error) {
	requiredThresholdMembers, err := kg.caller.GroupThreshold(kg.callerOpts)
	if err != nil {
		return 0, err
	}
	return int(requiredThresholdMembers.Int64()), nil
}

// GroupSize returns the number of members that are required
// to form a group.
func (kg *keepGroup) GroupSize() (int, error) {
	groupSize, err := kg.caller.GroupSize(kg.callerOpts)
	if err != nil {
		return 0, err
	}
	return int(groupSize.Int64()), nil
}

// GetNStaker - temporary code for Milestone 1 - will return
// the number of stackers that have been added to the contract.
func (kg *keepGroup) GetNStaker() (int, error) {
	stakerCount, err := kg.caller.GetNStaker(kg.callerOpts)
	if err != nil {
		return 0, err
	}
	return int(stakerCount.Int64()), nil
}

// AddStaker - temporary code for Milestone 1 - will add a
// staker to the contract.
func (kg *keepGroup) AddStaker(
	groupMemberID string,
) (*types.Transaction, error) {
	groupMemberIDByte32, err := toByte32([]byte(groupMemberID))
	if err != nil {
		return nil, err
	}
	return kg.transactor.AddStaker(
		kg.transactorOpts,
		groupMemberIDByte32,
	)
}

// Just Do It
func (kg *keepGroup) ResetStaker() (*types.Transaction, error) {
	return kg.transactor.ResetStaker(kg.transactorOpts)
}

// function getStaker(uint32 _index) public view returns ( bytes32 ) {
func (kg *keepGroup) GetStaker(index int) ([]byte, error) {
	staker, err := kg.caller.GetStaker(kg.callerOpts, uint32(index))
	if err != nil {
		return []byte{}, err
	}
	return staker[:], nil
}

// GetGroupMemberPubKey returns the public key for group number i at location
// in group j.
func (kg *keepGroup) GetGroupMemberPubKey(
	groupIndex,
	memberIndex int,
) ([]byte, error) {
	groupIndexBig := big.NewInt(int64(groupIndex))
	memberIndexBig := big.NewInt(int64(memberIndex))
	groupMemberKey, err := kg.caller.GetGroupMemberPubKey(
		kg.callerOpts,
		groupIndexBig,
		memberIndexBig,
	)
	if err != nil {
		return nil, err
	}
	return groupMemberKey[:], nil
}

// IsMember returns true if the member is a part of the specified group.
func (kg *keepGroup) IsMember(
	groupPubKey, memberPubKey []byte,
) (bool, error) {
	groupPubKeyArray, err := toByte32(groupPubKey)
	if err != nil {
		return false, err
	}
	memberPubKeyArray, err := toByte32(memberPubKey)
	if err != nil {
		return false, err
	}
	return kg.caller.IsMember(
		kg.callerOpts,
		groupPubKeyArray,
		memberPubKeyArray,
	)
}

// groupCompleteEventFunc defines the function that is called upon
// group completion.
type groupCompleteEventFunc func(groupPubKey []byte)

// WatchGroupCompleteEvent create a watch for the group completion event.
func (kg *keepGroup) WatchGroupCompleteEvent(
	success groupCompleteEventFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepGroupImplV1GroupCompleteEvent)
	eventSubscription, err := kg.contract.WatchGroupCompleteEvent(nil, eventChan)
	if err != nil {
		close(eventChan)
		return fmt.Errorf(
			"error creating watch for GroupCompleteEvent events [%v]",
			err,
		)
	}
	go func() {
		defer close(eventChan)
		defer eventSubscription.Unsubscribe()
		for {
			select {
			case event := <-eventChan:
				success(event.GroupPubKey[:])
				return

			case err := <-eventSubscription.Err():
				fail(err)
				return
			}
		}
	}()
	return nil
}

// groupErrorCodeFunc defines a function to watch for errors.
type groupErrorCodeFunc func(Code uint8)

// WatchGroupErrorCode creates a watch for the GroupErrorCode event.
func (kg *keepGroup) WatchGroupErrorCode(
	success groupErrorCodeFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepGroupImplV1GroupErrorCode)
	eventSubscription, err := kg.contract.WatchGroupErrorCode(nil, eventChan)
	if err != nil {
		close(eventChan)
		return fmt.Errorf(
			"failed go create watch for GroupErrorCode events: [%v]",
			err,
		)
	}
	go func() {
		defer close(eventChan)
		defer eventSubscription.Unsubscribe()
		for {
			select {
			case event := <-eventChan:
				success(event.Code)
				return

			case err := <-eventSubscription.Err():
				fail(err)
				return
			}
		}
	}()
	return nil
}

// groupExistsEventFunc defines the function that is called when creating
// a group.  Exists is true when the group already exists.
type groupExistsEventFunc func(groupPubKey []byte, Exists bool)

// WatchGroupExistsEvent watches for the GroupExists event.
func (kg *keepGroup) WatchGroupExistsEvent(
	success groupExistsEventFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepGroupImplV1GroupExistsEvent)
	eventSubscription, err := kg.contract.WatchGroupExistsEvent(nil, eventChan)
	if err != nil {
		close(eventChan)
		return fmt.Errorf(
			"error creating watch for GropExistsEvent events [%v]",
			err,
		)
	}
	go func() {
		defer close(eventChan)
		defer eventSubscription.Unsubscribe()
		for {
			select {
			case event := <-eventChan:
				success(event.GroupPubKey[:], event.Exists)
				return

			case err := <-eventSubscription.Err():
				fail(err)
				return
			}
		}
	}()
	return nil
}

// groupStartedEventFunc defiens the function that is called when
// watching for started groups.
type groupStartedEventFunc func(groupPubKey []byte)

// WatchGroupStartedEvent watch for GroupStartedEvent
func (kg *keepGroup) WatchGroupStartedEvent(
	success groupStartedEventFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepGroupImplV1GroupStartedEvent)
	eventSubscription, err := kg.contract.WatchGroupStartedEvent(nil, eventChan)
	if err != nil {
		close(eventChan)
		return fmt.Errorf(
			"error creating watch for GorupStartedEvent events [%v]",
			err,
		)
	}
	go func() {
		defer close(eventChan)
		defer eventSubscription.Unsubscribe()
		for {
			select {
			case event := <-eventChan:
				success(event.GroupPubKey[:])
				return

			case err := <-eventSubscription.Err():
				fail(err)
				return
			}
		}
	}()
	return nil
}

// onStakerAddedFunc is the type of function called when an OnStakerAdded event
// is observed on-chain and reported to a watching handler.
type onStakerAddedFunc func(index int, groupMemberID []byte)

// WatchGroupStartedEvent watch for GroupStartedEvent
func (kg *keepGroup) WatchOnStakerAdded(
	success onStakerAddedFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepGroupImplV1OnStakerAdded)
	eventSubscription, err := kg.contract.WatchOnStakerAdded(nil, eventChan)
	if err != nil {
		close(eventChan)
		return fmt.Errorf("error creating watch for OnStakerAdded events [%v]", err)
	}
	go func() {
		defer close(eventChan)
		defer eventSubscription.Unsubscribe()
		for {
			select {
			case event := <-eventChan:
				success(int(event.Index), event.GroupMemberID[:])
				return

			case err := <-eventSubscription.Err():
				fail(err)
				return
			}
		}
	}()
	return nil
}
