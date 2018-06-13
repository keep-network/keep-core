package ethereum

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gen "github.com/keep-network/keep-core/pkg/chain/gen"
)

// keepGroup connection information for interface to KeepGroup contract
type keepGroup struct {
	caller          *gen.KeepGroupImplV1Caller
	callerOpts      *bind.CallOpts
	transactor      *gen.KeepGroupImplV1Transactor
	transactorOpts  *bind.TransactOpts
	contract        *gen.KeepGroupImplV1
	contractAddress common.Address
	name            string
}

// NewKeepGroup creates the necessary connections and configurations
// for accessing the KeepGroup contract.
func newKeepGroup(pv *ethereumChain) (*keepGroup, error) {
	// Proxy Address
	contractAddressHex, ok := pv.config.ContractAddresses["KeepGroupImplV1"]
	if !ok {
		return nil, fmt.Errorf("Error: no address information for 'KeepGroup' in configuration\n")
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	groupTransactor, err := gen.NewKeepGroupImplV1Transactor(contractAddress,
		pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to instantiate a KeepRelayBeaconTranactor contract: %s",
			err,
		)
	}

	file, err := os.Open(pv.config.Account.KeyFile)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to open keyfile: %s, %s",
			err,
			pv.config.Account.KeyFile,
		)
	}

	optsTransactor, err := bind.NewTransactor(
		bufio.NewReader(file),
		pv.config.Account.KeyFilePassword,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to read keyfile: %s, %s",
			err,
			pv.config.Account.KeyFile,
		)
	}

	groupCaller, err := gen.NewKeepGroupImplV1Caller(contractAddress, pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to instantiate a KeepRelayBeaconCaller contract: %s",
			err,
		)
	}

	optsCaller := &bind.CallOpts{
		Pending: false,
		From:    contractAddress,
		Context: nil,
	}

	krbContract, err := gen.NewKeepGroupImplV1(contractAddress, pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to instantiate contract object: %s at address: %s",
			err,
			contractAddressHex,
		)
	}

	return &keepGroup{
		name:            "KeepGroup", // "KeepGroupImplV1",
		transactor:      groupTransactor,
		transactorOpts:  optsTransactor,
		caller:          groupCaller,
		callerOpts:      optsCaller,
		contract:        krbContract,
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
	groupPubKeyArray, err := ToByte32(groupPubKey)
	if err != nil {
		return nil, err
	}
	return kg.transactor.DissolveGroup(kg.transactorOpts, groupPubKeyArray)
}

// CreateGroup starts a new group with the specified public key.
func (kg *keepGroup) CreateGroup(
	groupPubKey []byte,
) (*types.Transaction, error) {
	groupPubKeyArray, err := ToByte32(groupPubKey)
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

// GetGroupMemberPubKey returns the public key for group number i at location
// in group j
func (kg *keepGroup) GetGroupMemberPubKey(i, j int) ([]byte, error) {
	iBig := big.NewInt(int64(i))
	jBig := big.NewInt(int64(j))
	groupMemberKey, err := kg.caller.GetGroupMemberPubKey(kg.callerOpts, iBig, jBig)
	if err != nil {
		return nil, err
	}
	return groupMemberKey[:], nil
}

// IsMember returns true if the member is a part of the specified group
func (kg *keepGroup) IsMember(
	groupPubKey, memberPubKey []byte,
) (bool, error) {
	groupPubKeyArray, err := ToByte32(groupPubKey)
	if err != nil {
		return false, err
	}
	memberPubKeyArray, err := ToByte32(memberPubKey)
	if err != nil {
		return false, err
	}
	return kg.caller.IsMember(
		kg.callerOpts,
		groupPubKeyArray,
		memberPubKeyArray,
	)
}

// groupCompleteEvent defines the function that is called upon
// group completion
type groupCompleteEvent func(groupPubKey []byte)

// WatchGroupCompleteEvent create a watch for the group completion event
func (kg *keepGroup) WatchGroupCompleteEvent(
	success groupCompleteEvent,
	fail errorCallback,
) error {
	name := "GroupCompleteEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupCompleteEvent, 10)
	event, err := kg.contract.WatchGroupCompleteEvent(nil, sink)
	if err != nil {
		return fmt.Errorf("error creating watch for [%s] events [%s]", name, err)
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				success(rn.GroupPubKey[:])

			case ee := <-event.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// groupErrorCode defines a function to watch for errors
type groupErrorCode func(Code uint8)

// WatchGroupErrorCode creates a watch for the GroupErrorCode event
func (kg *keepGroup) WatchGroupErrorCode(
	success groupErrorCode,
	fail errorCallback,
) error {
	name := "GroupErrorCode"
	sink := make(chan *gen.KeepGroupImplV1GroupErrorCode, 10)
	event, err := kg.contract.WatchGroupErrorCode(nil, sink)
	if err != nil {
		return fmt.Errorf("Error creating watch for %s events: %s", name, err)
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				success(rn.Code)

			case ee := <-event.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// groupExistsEvent defines the function that is called when creating
// a group.  Exists is true when the group already exists.
type groupExistsEvent func(groupPubKey []byte, Exists bool)

// WatchGroupExistsEvent watches for the GroupExists event.
func (kg *keepGroup) WatchGroupExistsEvent(
	success groupExistsEvent,
	fail errorCallback,
) error {
	name := "GroupExistsEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupExistsEvent, 10)
	event, err := kg.contract.WatchGroupExistsEvent(nil, sink)
	if err != nil {
		return fmt.Errorf("error creating watch for [%s] events [%s]", name, err)
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				success(rn.GroupPubKey[:], rn.Exists)

			case ee := <-event.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// groupStartedEvent defiens the function that is called when
// watching for started groups.
type groupStartedEvent func(groupPubKey []byte)

// WatchGroupStartedEvent watch for GroupStartedEvent
func (kg *keepGroup) WatchGroupStartedEvent(
	success groupStartedEvent,
	fail errorCallback,
) error {
	name := "GroupStartedEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupStartedEvent, 10)
	event, err := kg.contract.WatchGroupStartedEvent(nil, sink)
	if err != nil {
		return fmt.Errorf("error creating watch for [%s] events [%s]", name, err)
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				success(rn.GroupPubKey[:])

			case ee := <-event.Err():
				fail(ee)
			}
		}
	}()
	return nil
}
