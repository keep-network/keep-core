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

// NewKeepGroup creates the necessary connections and configurations for accessing the KeepGroup contract.
func NewKeepGroup(pv *ethereumChain) (rv *KeepGroup, err error) {

	ContractAddressHex := pv.config.ContractAddresses["KeepGroup"] // Proxy Address
	contractAddress := common.HexToAddress(ContractAddressHex)

	krbTransactor, err := gen.NewKeepGroupImplV1Transactor(contractAddress, pv.client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconTranactor contract: %s", err)
		return
	}

	file, err := os.Open(pv.config.Account.KeyFile)
	if err != nil {
		log.Printf("Failed to open keyfile: %v, %s", err, pv.config.Account.KeyFile)
		return
	}

	optsTransactor, err := bind.NewTransactor(bufio.NewReader(file), pv.config.Account.KeyFilePassword)
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
		log.Printf("Failed to instantiate contract object: %v at address: %s", err, ContractAddressHex)
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

// Initialized calls the contract and returns true if the contract has had its Initialize method called.
func (kg *KeepGroup) Initialized() (bool, error) {
	return kg.caller.Initialized(kg.callerOpts)
}

// SetGroupThreshold sets the group threshold - this is the number of members reporting a generated
// value out of the group size.
func (kg *KeepGroup) SetGroupThreshold(groupThreshold int) (tx *types.Transaction, err error) {
	thr := big.NewInt(int64(groupThreshold))
	// function setGroupThreshold(uint256 _groupThreshold) public onlyOwner {
	tx, err = kg.transactor.SetGroupThreshold(kg.transactorOpts, thr)
	return
}

// GroupExists TODO
func (kg *KeepGroup) GroupExists(groupPubKey []byte) (*types.Transaction, error) {
	//    function groupExists(bytes32 _groupPubKey) public {
	return kg.transactor.GroupExists(kg.transactorOpts, ToByte32(groupPubKey))
}

// AddMemberToGroup adds a new member to at group.
func (kg *KeepGroup) AddMemberToGroup(groupPubKey, memberPubKey []byte) (*types.Transaction, error) {
	// function addMemberToGroup(bytes32 _groupPubKey, bytes32 _memberPubKey) public isStaked returns(bool) {
	return kg.transactor.AddMemberToGroup(kg.transactorOpts, ToByte32(groupPubKey), ToByte32(memberPubKey))
}

// DissolveGroup breaks up the group that is associated with the public key.
func (kg *KeepGroup) DissolveGroup(groupPubKey []byte) (*types.Transaction, error) {
	// function disolveGroup(bytes32 _groupPubKey) public onlyOwner returns(bool) {
	return kg.transactor.DisolveGroup(kg.transactorOpts, ToByte32(groupPubKey))
}

// CreateGroup starts a new group with the specified public key.
func (kg *KeepGroup) CreateGroup(groupPubKey []byte) (*types.Transaction, error) {
	// function createGroup(bytes32 _groupPubKey) public returns(bool) {
	return kg.transactor.CreateGroup(kg.transactorOpts, ToByte32(groupPubKey))
}

// GetNumberOfGroups returns the number of groups.
func (kg *KeepGroup) GetNumberOfGroups() (ng int, err error) {
	ngBig, err := kg.caller.GetNumberOfGroups(kg.callerOpts)
	// function getNumberOfGroups() public view returns(uint256) {
	if err == nil {
		ng = int(ngBig.Int64())
	}
	return
}

// GetGroupNMembers returns the Nth member of a group.  This is the number of the member.
func (kg *KeepGroup) GetGroupNMembers(groupNumber int) (nm int, err error) {
	iBigGroupNumber := big.NewInt(int64(groupNumber))
	// function getGroupNMembers(uint256 _i) public view returns(uint256) {
	ngBig, err := kg.caller.GetGroupNMembers(kg.callerOpts, iBigGroupNumber)
	if err == nil {
		nm = int(ngBig.Int64())
	}
	return
}

// GetGroupPubKey take the number of the group and returns that groups public key.
func (kg *KeepGroup) GetGroupPubKey(groupNumber int) (pub []byte, err error) {
	iBigGroupNumber := big.NewInt(int64(groupNumber))
	// function getGroupPubKey(uint256 _i) public view returns(bytes32) {
	tmp, err := kg.caller.GetGroupPubKey(kg.callerOpts, iBigGroupNumber)
	if err == nil {
		pub = tmp[:]
	}
	return
}

// SetGroupSize sets the number of members that is needed to form a group.
func (kg *KeepGroup) SetGroupSize(groupSize int) (tx *types.Transaction, err error) {
	thr := big.NewInt(int64(groupSize))
	// function setGroupSize(uint256 _groupSize) public onlyOwner {
	tx, err = kg.transactor.SetGroupSize(kg.transactorOpts, thr)
	return
}

// GetGroupThreshold returns the group threshold.  This is the number of members that
// have to report a value to create a new signature.
func (kg *KeepGroup) GetGroupThreshold() (thr int, err error) {
	// function getGroupThreshold() public view returns(uint256) {
	thrBig, err := kg.caller.GetGroupThreshold(kg.callerOpts)
	if err == nil {
		thr = int(thrBig.Int64())
	}
	return
}

// GetGroupSize returns the number of members that are required to form a group.
func (kg *KeepGroup) GetGroupSize() (thr int, err error) {
	// function getGroupSize() public view returns(uint256) {
	thrBig, err := kg.caller.GetGroupSize(kg.callerOpts)
	if err == nil {
		thr = int(thrBig.Int64())
	}
	return
}

// GetGroupNumber returns the number of a group given the public key of the group.
func (kg *KeepGroup) GetGroupNumber(groupPubKey []byte) (nm int, err error) {
	//    function getGroupNumber(bytes32 _groupPubKey) public view returns(uint256) {
	nmBig, err := kg.caller.GetGroupNumber(kg.callerOpts, ToByte32(groupPubKey))
	if err == nil {
		nm = int(nmBig.Int64())
	}
	return
}

// GetAllGroupMembers returns a set of public keys for each member of a group.
func (kg *KeepGroup) GetAllGroupMembers(groupPubKey []byte) (keySet [][]byte, err error) {
	groupNumber, err := kg.GetGroupNumber(groupPubKey)
	if err != nil {
		return
	}
	nm, err := kg.GetGroupNMembers(groupNumber)
	if err != nil {
		return
	}
	keySet = make([][]byte, 0, nm)
	for i := 0; i < nm; i++ {
		memberKey, err0 := kg.GetGroupPubKey(i)
		if err0 != nil {
			err = err0
			return
		}
		keySet = append(keySet, memberKey)
	}
	return
}

// GetGroupMemberPubKey returns the public key for group number i at location in group j
func (kg *KeepGroup) GetGroupMemberPubKey(i, j int) (pub []byte, err error) {
	iBig := big.NewInt(int64(i))
	jBig := big.NewInt(int64(j))
	// function getGroupMemberPubKey(uint256 _i, uint256 _j) public view returns(bytes32) {
	tmp, err := kg.caller.GetGroupMemberPubKey(kg.callerOpts, iBig, jBig)
	if err == nil {
		pub = tmp[:]
	}
	return
}

// GroupIsComplete return true if the group is complete
func (kg *KeepGroup) GroupIsComplete(groupPubKey []byte) (rv bool, err error) {
	// function groupIsComplete(bytes32 _groupPubKey) public view returns(bool) {
	rv, err = kg.caller.GroupIsComplete(kg.callerOpts, ToByte32(groupPubKey))
	return
}

// IsMember returns true if the member is a part of the specified group
func (kg *KeepGroup) IsMember(groupPubKey, memberPubKey []byte) (rv bool, err error) {
	// function isMember(bytes32 _groupPubKey, bytes32 _memberPubKey) public view returns(bool) {
	rv, err = kg.caller.IsMember(kg.callerOpts, ToByte32(groupPubKey), ToByte32(memberPubKey))
	return
}

// FxGroupCompleteEvent defines the function that is called upon group completion
type FxGroupCompleteEvent func(GroupPubKey []byte)

// WatchGroupCompleteEvent create a watch for the group completion event
func (kg *KeepGroup) WatchGroupCompleteEvent(success FxGroupCompleteEvent, fail FxError) (err error) {
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
func (kg *KeepGroup) WatchGroupErrorCode(success FxGroupErrorCode, fail FxError) (err error) {
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

// FxGroupExistsEvent defines the function that is called when creating a group.  Exists is
// true when the group already exists.
type FxGroupExistsEvent func(GroupPubKey []byte, Exists bool)

// WatchGroupExistsEvent watches for the GroupExists event.
func (kg *KeepGroup) WatchGroupExistsEvent(success FxGroupExistsEvent, fail FxError) (err error) {
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

// FxGroupStartedEvent defiens the function that is called when watching for started groups.
type FxGroupStartedEvent func(GroupPubKey []byte)

// WatchGroupStartedEvent watch for GroupStartedEvent
func (kg *KeepGroup) WatchGroupStartedEvent(success FxGroupStartedEvent, fail FxError) (err error) {
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
