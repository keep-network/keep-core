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
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

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

func (kg *KeepGroup) Initialized() (bool, error) {
	return kg.caller.Initialized(kg.callerOpts)
}

func (kg *KeepGroup) SetGroupThreshold(groupThreshold int) (tx *types.Transaction, err error) {
	thr := big.NewInt(int64(groupThreshold))
	// function setGroupThreshold(uint256 _groupThreshold) public onlyOwner {
	tx, err = kg.transactor.SetGroupThreshold(kg.transactorOpts, thr)
	return
}

func (kg *KeepGroup) GroupExists(groupPubKey []byte) (*types.Transaction, error) {
	//    function groupExists(bytes32 _groupPubKey) public {
	return kg.transactor.GroupExists(kg.transactorOpts, ToByte32(groupPubKey))
}

func (kg *KeepGroup) AddMemberToGroup(groupPubKey, memberPubKey []byte) (*types.Transaction, error) {
	// function addMemberToGroup(bytes32 _groupPubKey, bytes32 _memberPubKey) public isStaked returns(bool) {
	return kg.transactor.AddMemberToGroup(kg.transactorOpts, ToByte32(groupPubKey), ToByte32(memberPubKey))
}

func (kg *KeepGroup) DissolveGroup(groupPubKey []byte) (*types.Transaction, error) {
	// function disolveGroup(bytes32 _groupPubKey) public onlyOwner returns(bool) {
	return kg.transactor.DisolveGroup(kg.transactorOpts, ToByte32(groupPubKey))
}

func (kg *KeepGroup) CreateGroup(groupPubKey []byte) (*types.Transaction, error) {
	// function createGroup(bytes32 _groupPubKey) public returns(bool) {
	return kg.transactor.CreateGroup(kg.transactorOpts, ToByte32(groupPubKey))
}

func (kg *KeepGroup) GetNumberOfGroups() (ng int, err error) {
	ngBig, err := kg.caller.GetNumberOfGroups(kg.callerOpts)
	// function getNumberOfGroups() public view returns(uint256) {
	if err == nil {
		ng = int(ngBig.Int64())
	}
	return
}

func (kg *KeepGroup) GetGroupNMembers(groupNumber int) (nm int, err error) {
	iBigGroupNumber := big.NewInt(int64(groupNumber))
	// function getGroupNMembers(uint256 _i) public view returns(uint256) {
	ngBig, err := kg.caller.GetGroupNMembers(kg.callerOpts, iBigGroupNumber)
	if err == nil {
		nm = int(ngBig.Int64())
	}
	return
}

func (kg *KeepGroup) GetGroupPubKey(groupNumber int) (pub []byte, err error) {
	iBigGroupNumber := big.NewInt(int64(groupNumber))
	// function getGroupPubKey(uint256 _i) public view returns(bytes32) {
	tmp, err := kg.caller.GetGroupPubKey(kg.callerOpts, iBigGroupNumber)
	if err == nil {
		pub = tmp[:]
	}
	return
}

func (kg *KeepGroup) GetGroupNumber(groupPubKey []byte) (nm int, err error) {
	//    function getGroupNumber(bytes32 _groupPubKey) public view returns(uint256) {
	nmBig, err := kg.caller.GetGroupNumber(kg.callerOpts, ToByte32(groupPubKey))
	if err == nil {
		nm = int(nmBig.Int64())
	}
	return
}

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

func (kg *KeepGroup) GroupIsComplete(groupPubKey []byte) (rv bool, err error) {
	// function groupIsComplete(bytes32 _groupPubKey) public view returns(bool) {
	rv, err = kg.caller.GroupIsComplete(kg.callerOpts, ToByte32(groupPubKey))
	return
}

func (kg *KeepGroup) GroupExistsView(groupPubKey []byte) (rv bool, err error) {
	// function groupExistsView(bytes32 _groupPubKey) public view returns(bool) {
	rv, err = kg.caller.GroupExistsView(kg.callerOpts, ToByte32(groupPubKey))
	return
}

func (kg *KeepGroup) IsMember(groupPubKey, memberPubKey []byte) (rv bool, err error) {
	// function isMember(bytes32 _groupPubKey, bytes32 _memberPubKey) public view returns(bool) {
	rv, err = kg.caller.IsMember(kg.callerOpts, ToByte32(groupPubKey), ToByte32(memberPubKey))
	return
}

type FxGroupCompleteEvent func(GroupPubKey []byte)

func (kg *KeepGroup) WatchGroupCompleteEvent(success FxGroupCompleteEvent, fail FxError) (err error) {
	name := "GroupCompleteEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupCompleteEvent, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := kg.contract.WatchGroupCompleteEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				if db1 {
					fmt.Printf("%sGot a [%s] event! Yea!%s, %+v\n", MiscLib.ColorGreen, name, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))
				}
				success(rn.GroupPubKey[:])

			case ee := <-event.Err():
				if db1 {
					fmt.Printf("%sGot an error: %s%s\n", MiscLib.ColorYellow, ee, MiscLib.ColorReset)
				}
				fail(ee)
			}
		}
	}()
	return
}

type FxGroupErrorCode func(Code uint8)

func (kg *KeepGroup) WatchGroupErrorCode(success FxGroupErrorCode, fail FxError) (err error) {
	name := "GroupErrorCode"
	sink := make(chan *gen.KeepGroupImplV1GroupErrorCode, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := kg.contract.WatchGroupErrorCode(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				if db1 {
					fmt.Printf("%sGot a [%s] event! Yea!%s, %+v\n", MiscLib.ColorGreen, name, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))
				}
				success(rn.Code)

			case ee := <-event.Err():
				if db1 {
					fmt.Printf("%sGot an error: %s%s\n", MiscLib.ColorYellow, ee, MiscLib.ColorReset)
				}
				fail(ee)
			}
		}
	}()
	return
}

type FxGroupExistsEvent func(GroupPubKey []byte, Exists bool)

func (kg *KeepGroup) WatchGroupExistsEvent(success FxGroupExistsEvent, fail FxError) (err error) {
	name := "GroupExistsEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupExistsEvent, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := kg.contract.WatchGroupExistsEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				if db1 {
					fmt.Printf("%sGot a [%s] event! Yea!%s, %+v\n", MiscLib.ColorGreen, name, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))
				}
				success(rn.GroupPubKey[:], rn.Exists)

			case ee := <-event.Err():
				if db1 {
					fmt.Printf("%sGot an error: %s%s\n", MiscLib.ColorYellow, ee, MiscLib.ColorReset)
				}
				fail(ee)
			}
		}
	}()
	return
}

type FxGroupStartedEvent func(GroupPubKey []byte)

func (kg *KeepGroup) WatchGroupStartedEvent(success FxGroupStartedEvent, fail FxError) (err error) {
	name := "GroupStartedEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupStartedEvent, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := kg.contract.WatchGroupStartedEvent(nil, sink)
	if err != nil {
		log.Printf("Error creating watch for %s events: %s", name, err)
		return
	}
	go func() {
		for {
			select {
			case rn := <-sink:
				if db1 {
					fmt.Printf("%sGot a [%s] event! Yea!%s, %+v\n", MiscLib.ColorGreen, name, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))
				}
				success(rn.GroupPubKey[:])

			case ee := <-event.Err():
				if db1 {
					fmt.Printf("%sGot an error: %s%s\n", MiscLib.ColorYellow, ee, MiscLib.ColorReset)
				}
				fail(ee)
			}
		}
	}()
	return
}

/*
GoKeepGroup.go:20:6: exported type KeepGroup should have comment or be unexported
GoKeepGroup.go:33:1: exported function NewKeepGroup should have comment or be unexported
GoKeepGroup.go:94:1: exported method KeepGroup.Initialized should have comment or be unexported
GoKeepGroup.go:98:1: exported method KeepGroup.SetGroupThreshold should have comment or be unexported
GoKeepGroup.go:105:1: exported method KeepGroup.GroupExists should have comment or be unexported
GoKeepGroup.go:110:1: exported method KeepGroup.AddMemberToGroup should have comment or be unexported
GoKeepGroup.go:115:1: exported method KeepGroup.DissolveGroup should have comment or be unexported
GoKeepGroup.go:120:1: exported method KeepGroup.CreateGroup should have comment or be unexported
GoKeepGroup.go:125:1: exported method KeepGroup.GetNumberOfGroups should have comment or be unexported
GoKeepGroup.go:134:1: exported method KeepGroup.GetGroupNMembers should have comment or be unexported
GoKeepGroup.go:144:1: exported method KeepGroup.GetGroupPubKey should have comment or be unexported
GoKeepGroup.go:154:1: exported method KeepGroup.GetGroupNumber should have comment or be unexported
GoKeepGroup.go:163:1: exported method KeepGroup.GetAllGroupMembers should have comment or be unexported
GoKeepGroup.go:184:1: exported method KeepGroup.GetGroupMemberPubKey should have comment or be unexported
GoKeepGroup.go:195:1: exported method KeepGroup.GroupIsComplete should have comment or be unexported
GoKeepGroup.go:201:1: exported method KeepGroup.GroupExistsView should have comment or be unexported
GoKeepGroup.go:207:1: exported method KeepGroup.IsMember should have comment or be unexported
GoKeepGroup.go:213:6: exported type FxGroupCompleteEvent should have comment or be unexported
GoKeepGroup.go:215:1: exported method KeepGroup.WatchGroupCompleteEvent should have comment or be unexported
GoKeepGroup.go:247:6: exported type FxGroupErrorCode should have comment or be unexported
GoKeepGroup.go:249:1: exported method KeepGroup.WatchGroupErrorCode should have comment or be unexported
GoKeepGroup.go:281:6: exported type FxGroupExistsEvent should have comment or be unexported
GoKeepGroup.go:283:1: exported method KeepGroup.WatchGroupExistsEvent should have comment or be unexported
GoKeepGroup.go:315:6: exported type FxGroupStartedEvent should have comment or be unexported
GoKeepGroup.go:317:1: exported method KeepGroup.WatchGroupStartedEvent should have comment or be unexported
*/
