package ethereum

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gen "github.com/keep-network/keep-core/pkg/chain/gen"
)

// KeepGroup connection information for interface to KeepGroup contract
type KeepGroup struct {
	provider        *ethereumChain
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
func newKeepGroup(pv *ethereumChain) (*KeepGroup, error) {

	// Proxy Address
	ContractAddressHex, ok := pv.config.ContractAddresses["KeepGroupImplV1"]
	if !ok {
		return nil, fmt.Errorf("Error: no address information for 'KeepGroup' in configuration\n")
	}
	contractAddress := common.HexToAddress(ContractAddressHex)

	groupTransactor, err := gen.NewKeepGroupImplV1Transactor(contractAddress,
		pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to instantiate a KeepRelayBeaconTranactor contract: %s",
			err)
	}

	file, err := os.Open(pv.config.Account.KeyFile)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to open keyfile: %s, %s", err, pv.config.Account.KeyFile)
	}

	optsTransactor, err := bind.NewTransactor(
		bufio.NewReader(file),
		pv.config.Account.KeyFilePassword)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to read keyfile: %s, %s", err, pv.config.Account.KeyFile)
	}

	groupCaller, err := gen.NewKeepGroupImplV1Caller(contractAddress, pv.client)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to instantiate a KeepRelayBeaconCaller contract: %s", err)
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
			err, ContractAddressHex)
	}

	return &KeepGroup{
		name:            "KeepGroup", // "KeepGroupImplV1",
		provider:        pv,
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
func (kg *KeepGroup) Initialized() (bool, error) {
	return kg.caller.Initialized(kg.callerOpts)
}

// DissolveGroup breaks up the group that is associated with the public key.
func (kg *KeepGroup) DissolveGroup(
	groupPubKey []byte,
) (*types.Transaction, error) {
	return kg.transactor.DissolveGroup(kg.transactorOpts, ToByte32(groupPubKey))
}

// CreateGroup starts a new group with the specified public key.
func (kg *KeepGroup) CreateGroup(
	groupPubKey []byte,
) (*types.Transaction, error) {
	return kg.transactor.CreateGroup(kg.transactorOpts, ToByte32(groupPubKey))
}

// GetGroupPubKey take the number of the group and returns that groups
// public key.
func (kg *KeepGroup) GetGroupPubKey(
	groupIndex int,
) ([]byte, error) {
	var pub []byte
	iBigGroupNumber := big.NewInt(int64(groupIndex))
	tmp, err := kg.caller.GetGroupPubKey(kg.callerOpts, iBigGroupNumber)
	if err == nil {
		pub = tmp[:]
	}
	return pub, err
}

// GroupThreshold returns the group threshold.  This is the number
// of members that have to report a value to create a new signature.
func (kg *KeepGroup) GroupThreshold() (int, error) {
	thrBig, err := kg.caller.GroupThreshold(kg.callerOpts)
	if err != nil {
		return 0, err
	}
	return int(thrBig.Int64()), nil
}

// GroupSize returns the number of members that are required
// to form a group.
func (kg *KeepGroup) GroupSize() (int, error) {
	sizeBig, err := kg.caller.GroupSize(kg.callerOpts)
	if err != nil {
		return 0, err
	}
	return int(sizeBig.Int64()), nil
}

// GetGroupMemberPubKey returns the public key for group number i at location
// in group j
func (kg *KeepGroup) GetGroupMemberPubKey(i, j int) ([]byte, error) {
	var pub []byte
	iBig := big.NewInt(int64(i))
	jBig := big.NewInt(int64(j))
	tmp, err := kg.caller.GetGroupMemberPubKey(kg.callerOpts, iBig, jBig)
	if err == nil {
		pub = tmp[:]
	}
	return pub, err
}

// IsMember returns true if the member is a part of the specified group
func (kg *KeepGroup) IsMember(
	groupPubKey, memberPubKey []byte,
) (bool, error) {
	return kg.caller.IsMember(kg.callerOpts, ToByte32(groupPubKey),
		ToByte32(memberPubKey))
}

// groupCompleteEvent defines the function that is called upon
// group completion
type groupCompleteEvent func(GroupPubKey []byte)

// WatchGroupCompleteEvent create a watch for the group completion event
func (kg *KeepGroup) WatchGroupCompleteEvent(
	success groupCompleteEvent,
	fail errorCallback,
) error {
	name := "GroupCompleteEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupCompleteEvent, 10)
	event, err := kg.contract.WatchGroupCompleteEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return err
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
func (kg *KeepGroup) WatchGroupErrorCode(
	success groupErrorCode,
	fail errorCallback,
) error {
	name := "GroupErrorCode"
	sink := make(chan *gen.KeepGroupImplV1GroupErrorCode, 10)
	event, err := kg.contract.WatchGroupErrorCode(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return err
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
type groupExistsEvent func(GroupPubKey []byte, Exists bool)

// WatchGroupExistsEvent watches for the GroupExists event.
func (kg *KeepGroup) WatchGroupExistsEvent(
	success groupExistsEvent,
	fail errorCallback,
) error {
	name := "GroupExistsEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupExistsEvent, 10)
	event, err := kg.contract.WatchGroupExistsEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return err
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
type groupStartedEvent func(GroupPubKey []byte)

// WatchGroupStartedEvent watch for GroupStartedEvent
func (kg *KeepGroup) WatchGroupStartedEvent(
	success groupStartedEvent,
	fail errorCallback,
) error {
	name := "GroupStartedEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupStartedEvent, 10)
	event, err := kg.contract.WatchGroupStartedEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return err
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
