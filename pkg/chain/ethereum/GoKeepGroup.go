package ethereum

import (
	"bufio"
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
func NewKeepGroup(pv *ethereumChain) (rv *KeepGroup, err error) {

	// Proxy Address
	ContractAddressHex := pv.config.ContractAddresses["KeepGroup"]
	contractAddress := common.HexToAddress(ContractAddressHex)

	krbTransactor, err := gen.NewKeepGroupImplV1Transactor(contractAddress,
		pv.client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconTranactor contract: %s",
			err)
		return
	}

	file, err := os.Open(pv.config.Account.KeyFile)
	if err != nil {
		log.Printf("Failed to open keyfile: %v, %s", err, pv.config.Account.KeyFile)
		return
	}

	optsTransactor, err := bind.NewTransactor(bufio.NewReader(file),
		pv.config.Account.KeyFilePassword)
	if err != nil {
		log.Printf("Failed to read keyfile: %v, %s", err, pv.config.Account.KeyFile)
		return
	}

	krbCaller, err := gen.NewKeepGroupImplV1Caller(contractAddress, pv.client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconCaller contract: %s", err)
		return
	}

	optsCaller := &bind.CallOpts{
		Pending: false,
		From:    contractAddress,
		Context: nil,
	}

	krbContract, err := gen.NewKeepGroupImplV1(contractAddress, pv.client)
	if err != nil {
		log.Printf("Failed to instantiate contract object: %v at address: %s",
			err, ContractAddressHex)
		return
	}

	return &KeepGroup{
		name:            "KeepGroup", // "KeepGroupImplV1",
		provider:        pv,
		transactor:      krbTransactor,
		transactorOpts:  optsTransactor,
		caller:          krbCaller,
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
func (kg *KeepGroup) DissolveGroup(groupPubKey []byte) (*types.Transaction,
	error) {
	// function dissolveGroup(bytes32 _groupPubKey) public onlyOwner
	// returns(bool) {
	return kg.transactor.DissolveGroup(kg.transactorOpts, ToByte32(groupPubKey))
}

// CreateGroup starts a new group with the specified public key.
func (kg *KeepGroup) CreateGroup(groupPubKey []byte) (*types.Transaction,
	error) {
	// function createGroup(bytes32 _groupPubKey) public returns(bool) {
	return kg.transactor.CreateGroup(kg.transactorOpts, ToByte32(groupPubKey))
}

// NumberOfGroups returns the number of groups.
func (kg *KeepGroup) NumberOfGroups() (ng int, err error) {
	ngBig, err := kg.caller.NumberOfGroups(kg.callerOpts)
	// function numberOfGroups() public view returns(uint256) {
	if err == nil {
		ng = int(ngBig.Int64())
	}
	return
}

// GetGroupPubKey take the number of the group and returns that groups
// public key.
func (kg *KeepGroup) GetGroupPubKey(groupIndex int) (pub []byte,
	err error) {
	iBigGroupNumber := big.NewInt(int64(groupIndex))
	// function getGroupPubKey(uint256 _groupIndex) public view returns(bytes32) {
	tmp, err := kg.caller.GetGroupPubKey(kg.callerOpts, iBigGroupNumber)
	if err == nil {
		pub = tmp[:]
	}
	return
}

// GroupThreshold returns the group threshold.  This is the number
// of members that have to report a value to create a new signature.
func (kg *KeepGroup) GroupThreshold() (thr int, err error) {
	// function groupThreshold() public view returns(uint256) {
	thrBig, err := kg.caller.GroupThreshold(kg.callerOpts)
	if err == nil {
		thr = int(thrBig.Int64())
	}
	return
}

// GroupSize returns the number of members that are required
// to form a group.
func (kg *KeepGroup) GroupSize() (thr int, err error) {
	// function groupSize() public view returns(uint256) {
	thrBig, err := kg.caller.GroupSize(kg.callerOpts)
	if err == nil {
		thr = int(thrBig.Int64())
	}
	return
}

// GetGroupMemberPubKey returns the public key for group number i at location
// in group j
func (kg *KeepGroup) GetGroupMemberPubKey(i, j int) (pub []byte, err error) {
	iBig := big.NewInt(int64(i))
	jBig := big.NewInt(int64(j))
	// function getGroupMemberPubKey(uint256 _i, uint256 _j) public view
	// returns(bytes32) {
	tmp, err := kg.caller.GetGroupMemberPubKey(kg.callerOpts, iBig, jBig)
	if err == nil {
		pub = tmp[:]
	}
	return
}

// IsMember returns true if the member is a part of the specified group
func (kg *KeepGroup) IsMember(groupPubKey, memberPubKey []byte) (rv bool,
	err error) {
	// function isMember(bytes32 _groupPubKey, bytes32 _memberPubKey)
	// public view returns(bool) {
	rv, err = kg.caller.IsMember(kg.callerOpts, ToByte32(groupPubKey),
		ToByte32(memberPubKey))
	return
}

// FxGroupCompleteEvent defines the function that is called upon
// group completion
type FxGroupCompleteEvent func(GroupPubKey []byte)

// WatchGroupCompleteEvent create a watch for the group completion event
func (kg *KeepGroup) WatchGroupCompleteEvent(success FxGroupCompleteEvent,
	fail FxError) (err error) {
	name := "GroupCompleteEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupCompleteEvent, 10)
	event, err := kg.contract.WatchGroupCompleteEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
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
	return
}

// FxGroupErrorCode defines a function to watch for errors
type FxGroupErrorCode func(Code uint8)

// WatchGroupErrorCode creates a watch for the GroupErrorCode event
func (kg *KeepGroup) WatchGroupErrorCode(success FxGroupErrorCode,
	fail FxError) (err error) {
	name := "GroupErrorCode"
	sink := make(chan *gen.KeepGroupImplV1GroupErrorCode, 10)
	event, err := kg.contract.WatchGroupErrorCode(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
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
	return
}

// FxGroupExistsEvent defines the function that is called when creating
// a group.  Exists is true when the group already exists.
type FxGroupExistsEvent func(GroupPubKey []byte, Exists bool)

// WatchGroupExistsEvent watches for the GroupExists event.
func (kg *KeepGroup) WatchGroupExistsEvent(success FxGroupExistsEvent,
	fail FxError) (err error) {
	name := "GroupExistsEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupExistsEvent, 10)
	event, err := kg.contract.WatchGroupExistsEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
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
	return
}

// FxGroupStartedEvent defiens the function that is called when
// watching for started groups.
type FxGroupStartedEvent func(GroupPubKey []byte)

// WatchGroupStartedEvent watch for GroupStartedEvent
func (kg *KeepGroup) WatchGroupStartedEvent(success FxGroupStartedEvent,
	fail FxError) (err error) {
	name := "GroupStartedEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupStartedEvent, 10)
	event, err := kg.contract.WatchGroupStartedEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
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
	return
}
