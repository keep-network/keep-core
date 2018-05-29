package ethereum

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gen "github.com/keep-network/keep-core/pkg/chain/gen"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

type KeepGroup struct {
	Provider        *provider
	Caller          *gen.KeepGroupImplV1Caller
	CallerOpts      *bind.CallOpts
	Transactor      *gen.KeepGroupImplV1Transactor
	TransactorOpts  *bind.TransactOpts
	Contract        *gen.KeepGroupImplV1
	ABI             string
	ABIparsed       abi.ABI
	ContractAddress common.Address
	Name            string
}

func NewKeepGroup(pv *provider) (rv *KeepGroup, err error) {

	ContractAddressHex := pv.Config.ContractAddresses["KeepGroup"] // Proxy Address
	ContractAddress := common.HexToAddress(ContractAddressHex)

	krbTransactor, err := gen.NewKeepGroupImplV1Transactor(ContractAddress, pv.Client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconTranactor contract: %s", err)
		return
	}

	file, err := os.Open(pv.Config.Account.KeyFile)
	if err != nil {
		log.Printf("Failed to open keyfile: %v, %s", err, pv.Config.Account.KeyFile)
		return
	}

	optsTransactor, err := bind.NewTransactor(bufio.NewReader(file), pv.Config.Account.KeyFilePassword)
	if err != nil {
		log.Printf("Failed to read keyfile: %v, %s", err, pv.Config.Account.KeyFile)
		return
	}

	krbCaller, err := gen.NewKeepGroupImplV1Caller(ContractAddress, pv.Client)
	if err != nil {
		log.Printf("Failed to instantiate a KeepRelayBeaconCaller contract: %s", err)
		return
	}

	optsCaller := &bind.CallOpts{
		Pending: false,
		From:    ContractAddress,
		Context: nil,
	}

	parsed, err := abi.JSON(strings.NewReader(gen.KeepGroupImplV1ABI))
	if err != nil {
		log.Printf("Failed to parse ABI, error:%s", err)
		return
	}

	krbContract, err := gen.NewKeepGroupImplV1(ContractAddress, pv.Client)
	if err != nil {
		log.Printf("Failed to instantiate contract object: %v at address: %s", err, ContractAddressHex)
		return
	}

	return &KeepGroup{
		Name:            "KeepGroup", // "KeepGroupImplV1",
		Provider:        pv,
		Transactor:      krbTransactor,
		TransactorOpts:  optsTransactor,
		Caller:          krbCaller,
		CallerOpts:      optsCaller,
		Contract:        krbContract,
		ABI:             gen.KeepGroupImplV1ABI,
		ABIparsed:       parsed,
		ContractAddress: ContractAddress,
	}, nil
}

func (kg *KeepGroup) Initialized() (bool, error) {
	return kg.Caller.Initialized(kg.CallerOpts)
}

func (kg *KeepGroup) SetGroupThreshold(groupThreshold int) (tx *types.Transaction, err error) {
	thr := big.NewInt(int64(groupThreshold))
	// function setGroupThreshold(uint256 _groupThreshold) public onlyOwner {
	tx, err = kg.Transactor.SetGroupThreshold(kg.TransactorOpts, thr)
	return
}

func (kg *KeepGroup) GroupExists(groupPubKey []byte) (*types.Transaction, error) {
	//    function groupExists(bytes32 _groupPubKey) public {
	return kg.Transactor.GroupExists(kg.TransactorOpts, ToByte32(groupPubKey))
}

func (kg *KeepGroup) AddMemberToGroup(groupPubKey, memberPubKey []byte) (*types.Transaction, error) {
	// function addMemberToGroup(bytes32 _groupPubKey, bytes32 _memberPubKey) public isStaked returns(bool) {
	return kg.Transactor.AddMemberToGroup(kg.TransactorOpts, ToByte32(groupPubKey), ToByte32(memberPubKey))
}

func (kg *KeepGroup) DissolveGroup(groupPubKey []byte) (*types.Transaction, error) {
	// function disolveGroup(bytes32 _groupPubKey) public onlyOwner returns(bool) {
	return kg.Transactor.DisolveGroup(kg.TransactorOpts, ToByte32(groupPubKey))
}

func (kg *KeepGroup) CreateGroup(groupPubKey []byte) (*types.Transaction, error) {
	// function createGroup(bytes32 _groupPubKey) public returns(bool) {
	return kg.Transactor.CreateGroup(kg.TransactorOpts, ToByte32(groupPubKey))
}

func (kg *KeepGroup) GetNumberOfGroups() (ng int, err error) {
	ngBig, err := kg.Caller.GetNumberOfGroups(kg.CallerOpts)
	// function getNumberOfGroups() public view returns(uint256) {
	if err == nil {
		ng = int(ngBig.Int64())
	}
	return
}

func (kg *KeepGroup) GetGroupNMembers(groupNumber int) (nm int, err error) {
	iBigGroupNumber := big.NewInt(int64(groupNumber))
	// function getGroupNMembers(uint256 _i) public view returns(uint256) {
	ngBig, err := kg.Caller.GetGroupNMembers(kg.CallerOpts, iBigGroupNumber)
	if err == nil {
		nm = int(ngBig.Int64())
	}
	return
}

func (kg *KeepGroup) GetGroupPubKey(groupNumber int) (pub []byte, err error) {
	iBigGroupNumber := big.NewInt(int64(groupNumber))
	// function getGroupPubKey(uint256 _i) public view returns(bytes32) {
	tmp, err := kg.Caller.GetGroupPubKey(kg.CallerOpts, iBigGroupNumber)
	if err == nil {
		pub = tmp[:]
	}
	return
}

func (kg *KeepGroup) GetGroupNumber(groupPubKey []byte) (nm int, err error) {
	//    function getGroupNumber(bytes32 _groupPubKey) public view returns(uint256) {
	nmBig, err := kg.Caller.GetGroupNumber(kg.CallerOpts, ToByte32(groupPubKey))
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
	tmp, err := kg.Caller.GetGroupMemberPubKey(kg.CallerOpts, iBig, jBig)
	if err == nil {
		pub = tmp[:]
	}
	return
}

func (kg *KeepGroup) GroupIsComplete(groupPubKey []byte) (rv bool, err error) {
	// function groupIsComplete(bytes32 _groupPubKey) public view returns(bool) {
	rv, err = kg.Caller.GroupIsComplete(kg.CallerOpts, ToByte32(groupPubKey))
	return
}

func (kg *KeepGroup) GroupExistsView(groupPubKey []byte) (rv bool, err error) {
	// function groupExistsView(bytes32 _groupPubKey) public view returns(bool) {
	rv, err = kg.Caller.GroupExistsView(kg.CallerOpts, ToByte32(groupPubKey))
	return
}

func (kg *KeepGroup) IsMember(groupPubKey, memberPubKey []byte) (rv bool, err error) {
	// function isMember(bytes32 _groupPubKey, bytes32 _memberPubKey) public view returns(bool) {
	rv, err = kg.Caller.IsMember(kg.CallerOpts, ToByte32(groupPubKey), ToByte32(memberPubKey))
	return
}

type FxGroupCompleteEvent func(GroupPubKey []byte)

func (kg *KeepGroup) WatchGroupCompleteEvent(success FxGroupCompleteEvent, fail FxError) (err error) {
	name := "GroupCompleteEvent"
	sink := make(chan *gen.KeepGroupImplV1GroupCompleteEvent, 10)
	if db1 {
		fmt.Printf("Calling Watch for %s, %s\n", name, godebug.LF())
	}
	event, err := kg.Contract.WatchGroupCompleteEvent(nil, sink)
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
	event, err := kg.Contract.WatchGroupErrorCode(nil, sink)
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
	event, err := kg.Contract.WatchGroupExistsEvent(nil, sink)
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
	event, err := kg.Contract.WatchGroupStartedEvent(nil, sink)
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
