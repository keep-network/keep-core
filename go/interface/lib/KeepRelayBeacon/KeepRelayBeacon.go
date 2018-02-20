// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package KeepRelayBeacon

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// KeepRelayBeaconABI is the input ABI used to generate the binding from.
const KeepRelayBeaconABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"groupID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_payment\",\"type\":\"uint256\"},{\"name\":\"_blockReward\",\"type\":\"uint256\"},{\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"requestRelay\",\"outputs\":[{\"name\":\"RequestID\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"KStart\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"blockReward\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_RequestID\",\"type\":\"uint256\"},{\"name\":\"_groupSignature\",\"type\":\"uint256\"},{\"name\":\"_groupID\",\"type\":\"uint256\"},{\"name\":\"_previousEntry\",\"type\":\"uint256\"}],\"name\":\"relayEntry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_LastValidRelayTxHash\",\"type\":\"uint256\"},{\"name\":\"_LastValidRelayBlock\",\"type\":\"uint256\"}],\"name\":\"relayEntryAcustation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_PK_G_i\",\"type\":\"uint256\"},{\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"submitGroupPublicKey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"signature\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"payment\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"seed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_RequestID\",\"type\":\"uint256\"}],\"name\":\"getRandomNumber\",\"outputs\":[{\"name\":\"theRandomNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_UserPublicKey\",\"type\":\"uint256\"}],\"name\":\"isStaked\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"GenRequestIDSequence\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"RequestID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Payment\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"BlockReward\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Seed\",\"type\":\"uint256\"}],\"name\":\"RelayEntryRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"RequestID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Signature\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"GroupID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"PreviousEntry\",\"type\":\"uint256\"}],\"name\":\"RelayEntryGenerated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"LastValidRelayEntry\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"LastValidRelayTxHash\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"LastValidRelayBlock\",\"type\":\"uint256\"}],\"name\":\"RelayResetEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_PK_G_i\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_id\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_activationBlockHeight\",\"type\":\"uint256\"}],\"name\":\"SubmitGroupPublicKeyEvent\",\"type\":\"event\"}]"

// KeepRelayBeacon is an auto generated Go binding around an Ethereum contract.
type KeepRelayBeacon struct {
	KeepRelayBeaconCaller     // Read-only binding to the contract
	KeepRelayBeaconTransactor // Write-only binding to the contract
	KeepRelayBeaconFilterer   // Log filterer for contract events
}

// KeepRelayBeaconCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeepRelayBeaconCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRelayBeaconTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeepRelayBeaconTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRelayBeaconFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeepRelayBeaconFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRelayBeaconSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeepRelayBeaconSession struct {
	Contract     *KeepRelayBeacon  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeepRelayBeaconCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeepRelayBeaconCallerSession struct {
	Contract *KeepRelayBeaconCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// KeepRelayBeaconTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeepRelayBeaconTransactorSession struct {
	Contract     *KeepRelayBeaconTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// KeepRelayBeaconRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeepRelayBeaconRaw struct {
	Contract *KeepRelayBeacon // Generic contract binding to access the raw methods on
}

// KeepRelayBeaconCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeepRelayBeaconCallerRaw struct {
	Contract *KeepRelayBeaconCaller // Generic read-only contract binding to access the raw methods on
}

// KeepRelayBeaconTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeepRelayBeaconTransactorRaw struct {
	Contract *KeepRelayBeaconTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeepRelayBeacon creates a new instance of KeepRelayBeacon, bound to a specific deployed contract.
func NewKeepRelayBeacon(address common.Address, backend bind.ContractBackend) (*KeepRelayBeacon, error) {
	contract, err := bindKeepRelayBeacon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeepRelayBeacon{KeepRelayBeaconCaller: KeepRelayBeaconCaller{contract: contract}, KeepRelayBeaconTransactor: KeepRelayBeaconTransactor{contract: contract}, KeepRelayBeaconFilterer: KeepRelayBeaconFilterer{contract: contract}}, nil
}

// NewKeepRelayBeaconCaller creates a new read-only instance of KeepRelayBeacon, bound to a specific deployed contract.
func NewKeepRelayBeaconCaller(address common.Address, caller bind.ContractCaller) (*KeepRelayBeaconCaller, error) {
	contract, err := bindKeepRelayBeacon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRelayBeaconCaller{contract: contract}, nil
}

// NewKeepRelayBeaconTransactor creates a new write-only instance of KeepRelayBeacon, bound to a specific deployed contract.
func NewKeepRelayBeaconTransactor(address common.Address, transactor bind.ContractTransactor) (*KeepRelayBeaconTransactor, error) {
	contract, err := bindKeepRelayBeacon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRelayBeaconTransactor{contract: contract}, nil
}

// NewKeepRelayBeaconFilterer creates a new log filterer instance of KeepRelayBeacon, bound to a specific deployed contract.
func NewKeepRelayBeaconFilterer(address common.Address, filterer bind.ContractFilterer) (*KeepRelayBeaconFilterer, error) {
	contract, err := bindKeepRelayBeacon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeepRelayBeaconFilterer{contract: contract}, nil
}

// bindKeepRelayBeacon binds a generic wrapper to an already deployed contract.
func bindKeepRelayBeacon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeepRelayBeaconABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRelayBeacon *KeepRelayBeaconRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRelayBeacon.Contract.KeepRelayBeaconCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRelayBeacon *KeepRelayBeaconRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.KeepRelayBeaconTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRelayBeacon *KeepRelayBeaconRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.KeepRelayBeaconTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRelayBeacon *KeepRelayBeaconCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRelayBeacon.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRelayBeacon *KeepRelayBeaconTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRelayBeacon *KeepRelayBeaconTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.contract.Transact(opts, method, params...)
}

// GenRequestIDSequence is a free data retrieval call binding the contract method 0xe205f3d7.
//
// Solidity: function GenRequestIDSequence() constant returns(address)
func (_KeepRelayBeacon *KeepRelayBeaconCaller) GenRequestIDSequence(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRelayBeacon.contract.Call(opts, out, "GenRequestIDSequence")
	return *ret0, err
}

// GenRequestIDSequence is a free data retrieval call binding the contract method 0xe205f3d7.
//
// Solidity: function GenRequestIDSequence() constant returns(address)
func (_KeepRelayBeacon *KeepRelayBeaconSession) GenRequestIDSequence() (common.Address, error) {
	return _KeepRelayBeacon.Contract.GenRequestIDSequence(&_KeepRelayBeacon.CallOpts)
}

// GenRequestIDSequence is a free data retrieval call binding the contract method 0xe205f3d7.
//
// Solidity: function GenRequestIDSequence() constant returns(address)
func (_KeepRelayBeacon *KeepRelayBeaconCallerSession) GenRequestIDSequence() (common.Address, error) {
	return _KeepRelayBeacon.Contract.GenRequestIDSequence(&_KeepRelayBeacon.CallOpts)
}

// BlockReward is a free data retrieval call binding the contract method 0x25d1818b.
//
// Solidity: function blockReward( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCaller) BlockReward(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRelayBeacon.contract.Call(opts, out, "blockReward", arg0)
	return *ret0, err
}

// BlockReward is a free data retrieval call binding the contract method 0x25d1818b.
//
// Solidity: function blockReward( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconSession) BlockReward(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.BlockReward(&_KeepRelayBeacon.CallOpts, arg0)
}

// BlockReward is a free data retrieval call binding the contract method 0x25d1818b.
//
// Solidity: function blockReward( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCallerSession) BlockReward(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.BlockReward(&_KeepRelayBeacon.CallOpts, arg0)
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xb37217a4.
//
// Solidity: function getRandomNumber(_RequestID uint256) constant returns(theRandomNumber uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCaller) GetRandomNumber(opts *bind.CallOpts, _RequestID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRelayBeacon.contract.Call(opts, out, "getRandomNumber", _RequestID)
	return *ret0, err
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xb37217a4.
//
// Solidity: function getRandomNumber(_RequestID uint256) constant returns(theRandomNumber uint256)
func (_KeepRelayBeacon *KeepRelayBeaconSession) GetRandomNumber(_RequestID *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.GetRandomNumber(&_KeepRelayBeacon.CallOpts, _RequestID)
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xb37217a4.
//
// Solidity: function getRandomNumber(_RequestID uint256) constant returns(theRandomNumber uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCallerSession) GetRandomNumber(_RequestID *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.GetRandomNumber(&_KeepRelayBeacon.CallOpts, _RequestID)
}

// GroupID is a free data retrieval call binding the contract method 0x104e2471.
//
// Solidity: function groupID( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCaller) GroupID(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRelayBeacon.contract.Call(opts, out, "groupID", arg0)
	return *ret0, err
}

// GroupID is a free data retrieval call binding the contract method 0x104e2471.
//
// Solidity: function groupID( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconSession) GroupID(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.GroupID(&_KeepRelayBeacon.CallOpts, arg0)
}

// GroupID is a free data retrieval call binding the contract method 0x104e2471.
//
// Solidity: function groupID( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCallerSession) GroupID(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.GroupID(&_KeepRelayBeacon.CallOpts, arg0)
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked(_UserPublicKey uint256) constant returns(bool)
func (_KeepRelayBeacon *KeepRelayBeaconCaller) IsStaked(opts *bind.CallOpts, _UserPublicKey *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRelayBeacon.contract.Call(opts, out, "isStaked", _UserPublicKey)
	return *ret0, err
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked(_UserPublicKey uint256) constant returns(bool)
func (_KeepRelayBeacon *KeepRelayBeaconSession) IsStaked(_UserPublicKey *big.Int) (bool, error) {
	return _KeepRelayBeacon.Contract.IsStaked(&_KeepRelayBeacon.CallOpts, _UserPublicKey)
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked(_UserPublicKey uint256) constant returns(bool)
func (_KeepRelayBeacon *KeepRelayBeaconCallerSession) IsStaked(_UserPublicKey *big.Int) (bool, error) {
	return _KeepRelayBeacon.Contract.IsStaked(&_KeepRelayBeacon.CallOpts, _UserPublicKey)
}

// Payment is a free data retrieval call binding the contract method 0x8b3c99e3.
//
// Solidity: function payment( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCaller) Payment(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRelayBeacon.contract.Call(opts, out, "payment", arg0)
	return *ret0, err
}

// Payment is a free data retrieval call binding the contract method 0x8b3c99e3.
//
// Solidity: function payment( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconSession) Payment(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.Payment(&_KeepRelayBeacon.CallOpts, arg0)
}

// Payment is a free data retrieval call binding the contract method 0x8b3c99e3.
//
// Solidity: function payment( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCallerSession) Payment(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.Payment(&_KeepRelayBeacon.CallOpts, arg0)
}

// Seed is a free data retrieval call binding the contract method 0x95564837.
//
// Solidity: function seed( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCaller) Seed(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRelayBeacon.contract.Call(opts, out, "seed", arg0)
	return *ret0, err
}

// Seed is a free data retrieval call binding the contract method 0x95564837.
//
// Solidity: function seed( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconSession) Seed(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.Seed(&_KeepRelayBeacon.CallOpts, arg0)
}

// Seed is a free data retrieval call binding the contract method 0x95564837.
//
// Solidity: function seed( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCallerSession) Seed(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.Seed(&_KeepRelayBeacon.CallOpts, arg0)
}

// Signature is a free data retrieval call binding the contract method 0x70629548.
//
// Solidity: function signature( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCaller) Signature(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRelayBeacon.contract.Call(opts, out, "signature", arg0)
	return *ret0, err
}

// Signature is a free data retrieval call binding the contract method 0x70629548.
//
// Solidity: function signature( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconSession) Signature(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.Signature(&_KeepRelayBeacon.CallOpts, arg0)
}

// Signature is a free data retrieval call binding the contract method 0x70629548.
//
// Solidity: function signature( uint256) constant returns(uint256)
func (_KeepRelayBeacon *KeepRelayBeaconCallerSession) Signature(arg0 *big.Int) (*big.Int, error) {
	return _KeepRelayBeacon.Contract.Signature(&_KeepRelayBeacon.CallOpts, arg0)
}

// KStart is a paid mutator transaction binding the contract method 0x17dff6dc.
//
// Solidity: function KStart() returns()
func (_KeepRelayBeacon *KeepRelayBeaconTransactor) KStart(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRelayBeacon.contract.Transact(opts, "KStart")
}

// KStart is a paid mutator transaction binding the contract method 0x17dff6dc.
//
// Solidity: function KStart() returns()
func (_KeepRelayBeacon *KeepRelayBeaconSession) KStart() (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.KStart(&_KeepRelayBeacon.TransactOpts)
}

// KStart is a paid mutator transaction binding the contract method 0x17dff6dc.
//
// Solidity: function KStart() returns()
func (_KeepRelayBeacon *KeepRelayBeaconTransactorSession) KStart() (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.KStart(&_KeepRelayBeacon.TransactOpts)
}

// RelayEntry is a paid mutator transaction binding the contract method 0x262e691d.
//
// Solidity: function relayEntry(_RequestID uint256, _groupSignature uint256, _groupID uint256, _previousEntry uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconTransactor) RelayEntry(opts *bind.TransactOpts, _RequestID *big.Int, _groupSignature *big.Int, _groupID *big.Int, _previousEntry *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.contract.Transact(opts, "relayEntry", _RequestID, _groupSignature, _groupID, _previousEntry)
}

// RelayEntry is a paid mutator transaction binding the contract method 0x262e691d.
//
// Solidity: function relayEntry(_RequestID uint256, _groupSignature uint256, _groupID uint256, _previousEntry uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconSession) RelayEntry(_RequestID *big.Int, _groupSignature *big.Int, _groupID *big.Int, _previousEntry *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.RelayEntry(&_KeepRelayBeacon.TransactOpts, _RequestID, _groupSignature, _groupID, _previousEntry)
}

// RelayEntry is a paid mutator transaction binding the contract method 0x262e691d.
//
// Solidity: function relayEntry(_RequestID uint256, _groupSignature uint256, _groupID uint256, _previousEntry uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconTransactorSession) RelayEntry(_RequestID *big.Int, _groupSignature *big.Int, _groupID *big.Int, _previousEntry *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.RelayEntry(&_KeepRelayBeacon.TransactOpts, _RequestID, _groupSignature, _groupID, _previousEntry)
}

// RelayEntryAcustation is a paid mutator transaction binding the contract method 0x4f3f9ad1.
//
// Solidity: function relayEntryAcustation(_LastValidRelayTxHash uint256, _LastValidRelayBlock uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconTransactor) RelayEntryAcustation(opts *bind.TransactOpts, _LastValidRelayTxHash *big.Int, _LastValidRelayBlock *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.contract.Transact(opts, "relayEntryAcustation", _LastValidRelayTxHash, _LastValidRelayBlock)
}

// RelayEntryAcustation is a paid mutator transaction binding the contract method 0x4f3f9ad1.
//
// Solidity: function relayEntryAcustation(_LastValidRelayTxHash uint256, _LastValidRelayBlock uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconSession) RelayEntryAcustation(_LastValidRelayTxHash *big.Int, _LastValidRelayBlock *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.RelayEntryAcustation(&_KeepRelayBeacon.TransactOpts, _LastValidRelayTxHash, _LastValidRelayBlock)
}

// RelayEntryAcustation is a paid mutator transaction binding the contract method 0x4f3f9ad1.
//
// Solidity: function relayEntryAcustation(_LastValidRelayTxHash uint256, _LastValidRelayBlock uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconTransactorSession) RelayEntryAcustation(_LastValidRelayTxHash *big.Int, _LastValidRelayBlock *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.RelayEntryAcustation(&_KeepRelayBeacon.TransactOpts, _LastValidRelayTxHash, _LastValidRelayBlock)
}

// RequestRelay is a paid mutator transaction binding the contract method 0x13df0019.
//
// Solidity: function requestRelay(_payment uint256, _blockReward uint256, _seed uint256) returns(RequestID uint256)
func (_KeepRelayBeacon *KeepRelayBeaconTransactor) RequestRelay(opts *bind.TransactOpts, _payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.contract.Transact(opts, "requestRelay", _payment, _blockReward, _seed)
}

// RequestRelay is a paid mutator transaction binding the contract method 0x13df0019.
//
// Solidity: function requestRelay(_payment uint256, _blockReward uint256, _seed uint256) returns(RequestID uint256)
func (_KeepRelayBeacon *KeepRelayBeaconSession) RequestRelay(_payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.RequestRelay(&_KeepRelayBeacon.TransactOpts, _payment, _blockReward, _seed)
}

// RequestRelay is a paid mutator transaction binding the contract method 0x13df0019.
//
// Solidity: function requestRelay(_payment uint256, _blockReward uint256, _seed uint256) returns(RequestID uint256)
func (_KeepRelayBeacon *KeepRelayBeaconTransactorSession) RequestRelay(_payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.RequestRelay(&_KeepRelayBeacon.TransactOpts, _payment, _blockReward, _seed)
}

// SubmitGroupPublicKey is a paid mutator transaction binding the contract method 0x598b70ac.
//
// Solidity: function submitGroupPublicKey(_PK_G_i uint256, _id uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconTransactor) SubmitGroupPublicKey(opts *bind.TransactOpts, _PK_G_i *big.Int, _id *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.contract.Transact(opts, "submitGroupPublicKey", _PK_G_i, _id)
}

// SubmitGroupPublicKey is a paid mutator transaction binding the contract method 0x598b70ac.
//
// Solidity: function submitGroupPublicKey(_PK_G_i uint256, _id uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconSession) SubmitGroupPublicKey(_PK_G_i *big.Int, _id *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.SubmitGroupPublicKey(&_KeepRelayBeacon.TransactOpts, _PK_G_i, _id)
}

// SubmitGroupPublicKey is a paid mutator transaction binding the contract method 0x598b70ac.
//
// Solidity: function submitGroupPublicKey(_PK_G_i uint256, _id uint256) returns()
func (_KeepRelayBeacon *KeepRelayBeaconTransactorSession) SubmitGroupPublicKey(_PK_G_i *big.Int, _id *big.Int) (*types.Transaction, error) {
	return _KeepRelayBeacon.Contract.SubmitGroupPublicKey(&_KeepRelayBeacon.TransactOpts, _PK_G_i, _id)
}

// KeepRelayBeaconRelayEntryGeneratedIterator is returned from FilterRelayEntryGenerated and is used to iterate over the raw logs and unpacked data for RelayEntryGenerated events raised by the KeepRelayBeacon contract.
type KeepRelayBeaconRelayEntryGeneratedIterator struct {
	Event *KeepRelayBeaconRelayEntryGenerated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRelayBeaconRelayEntryGeneratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRelayBeaconRelayEntryGenerated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRelayBeaconRelayEntryGenerated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRelayBeaconRelayEntryGeneratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRelayBeaconRelayEntryGeneratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRelayBeaconRelayEntryGenerated represents a RelayEntryGenerated event raised by the KeepRelayBeacon contract.
type KeepRelayBeaconRelayEntryGenerated struct {
	RequestID     *big.Int
	Signature     *big.Int
	GroupID       *big.Int
	PreviousEntry *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRelayEntryGenerated is a free log retrieval operation binding the contract event 0x05f168e6dc72933f6c45c446d8f56a7aebb6e7fd04536360daff679e8a74279a.
//
// Solidity: event RelayEntryGenerated(RequestID uint256, Signature uint256, GroupID uint256, PreviousEntry uint256)
func (_KeepRelayBeacon *KeepRelayBeaconFilterer) FilterRelayEntryGenerated(opts *bind.FilterOpts) (*KeepRelayBeaconRelayEntryGeneratedIterator, error) {

	logs, sub, err := _KeepRelayBeacon.contract.FilterLogs(opts, "RelayEntryGenerated")
	if err != nil {
		return nil, err
	}
	return &KeepRelayBeaconRelayEntryGeneratedIterator{contract: _KeepRelayBeacon.contract, event: "RelayEntryGenerated", logs: logs, sub: sub}, nil
}

// WatchRelayEntryGenerated is a free log subscription operation binding the contract event 0x05f168e6dc72933f6c45c446d8f56a7aebb6e7fd04536360daff679e8a74279a.
//
// Solidity: event RelayEntryGenerated(RequestID uint256, Signature uint256, GroupID uint256, PreviousEntry uint256)
func (_KeepRelayBeacon *KeepRelayBeaconFilterer) WatchRelayEntryGenerated(opts *bind.WatchOpts, sink chan<- *KeepRelayBeaconRelayEntryGenerated) (event.Subscription, error) {

	logs, sub, err := _KeepRelayBeacon.contract.WatchLogs(opts, "RelayEntryGenerated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRelayBeaconRelayEntryGenerated)
				if err := _KeepRelayBeacon.contract.UnpackLog(event, "RelayEntryGenerated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// KeepRelayBeaconRelayEntryRequestedIterator is returned from FilterRelayEntryRequested and is used to iterate over the raw logs and unpacked data for RelayEntryRequested events raised by the KeepRelayBeacon contract.
type KeepRelayBeaconRelayEntryRequestedIterator struct {
	Event *KeepRelayBeaconRelayEntryRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRelayBeaconRelayEntryRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRelayBeaconRelayEntryRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRelayBeaconRelayEntryRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRelayBeaconRelayEntryRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRelayBeaconRelayEntryRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRelayBeaconRelayEntryRequested represents a RelayEntryRequested event raised by the KeepRelayBeacon contract.
type KeepRelayBeaconRelayEntryRequested struct {
	RequestID   *big.Int
	Payment     *big.Int
	BlockReward *big.Int
	Seed        *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRelayEntryRequested is a free log retrieval operation binding the contract event 0x934cbfb2fe04bb7630b83ee67109a80d5cd1a7244b202262c62531083d62315a.
//
// Solidity: event RelayEntryRequested(RequestID uint256, Payment uint256, BlockReward uint256, Seed uint256)
func (_KeepRelayBeacon *KeepRelayBeaconFilterer) FilterRelayEntryRequested(opts *bind.FilterOpts) (*KeepRelayBeaconRelayEntryRequestedIterator, error) {

	logs, sub, err := _KeepRelayBeacon.contract.FilterLogs(opts, "RelayEntryRequested")
	if err != nil {
		return nil, err
	}
	return &KeepRelayBeaconRelayEntryRequestedIterator{contract: _KeepRelayBeacon.contract, event: "RelayEntryRequested", logs: logs, sub: sub}, nil
}

// WatchRelayEntryRequested is a free log subscription operation binding the contract event 0x934cbfb2fe04bb7630b83ee67109a80d5cd1a7244b202262c62531083d62315a.
//
// Solidity: event RelayEntryRequested(RequestID uint256, Payment uint256, BlockReward uint256, Seed uint256)
func (_KeepRelayBeacon *KeepRelayBeaconFilterer) WatchRelayEntryRequested(opts *bind.WatchOpts, sink chan<- *KeepRelayBeaconRelayEntryRequested) (event.Subscription, error) {

	logs, sub, err := _KeepRelayBeacon.contract.WatchLogs(opts, "RelayEntryRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRelayBeaconRelayEntryRequested)
				if err := _KeepRelayBeacon.contract.UnpackLog(event, "RelayEntryRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// KeepRelayBeaconRelayResetEventIterator is returned from FilterRelayResetEvent and is used to iterate over the raw logs and unpacked data for RelayResetEvent events raised by the KeepRelayBeacon contract.
type KeepRelayBeaconRelayResetEventIterator struct {
	Event *KeepRelayBeaconRelayResetEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRelayBeaconRelayResetEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRelayBeaconRelayResetEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRelayBeaconRelayResetEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRelayBeaconRelayResetEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRelayBeaconRelayResetEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRelayBeaconRelayResetEvent represents a RelayResetEvent event raised by the KeepRelayBeacon contract.
type KeepRelayBeaconRelayResetEvent struct {
	LastValidRelayEntry  *big.Int
	LastValidRelayTxHash *big.Int
	LastValidRelayBlock  *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterRelayResetEvent is a free log retrieval operation binding the contract event 0x242ff963387d59926325b194e1d237810c5c08836df9074589e6208087fa49d3.
//
// Solidity: event RelayResetEvent(LastValidRelayEntry uint256, LastValidRelayTxHash uint256, LastValidRelayBlock uint256)
func (_KeepRelayBeacon *KeepRelayBeaconFilterer) FilterRelayResetEvent(opts *bind.FilterOpts) (*KeepRelayBeaconRelayResetEventIterator, error) {

	logs, sub, err := _KeepRelayBeacon.contract.FilterLogs(opts, "RelayResetEvent")
	if err != nil {
		return nil, err
	}
	return &KeepRelayBeaconRelayResetEventIterator{contract: _KeepRelayBeacon.contract, event: "RelayResetEvent", logs: logs, sub: sub}, nil
}

// WatchRelayResetEvent is a free log subscription operation binding the contract event 0x242ff963387d59926325b194e1d237810c5c08836df9074589e6208087fa49d3.
//
// Solidity: event RelayResetEvent(LastValidRelayEntry uint256, LastValidRelayTxHash uint256, LastValidRelayBlock uint256)
func (_KeepRelayBeacon *KeepRelayBeaconFilterer) WatchRelayResetEvent(opts *bind.WatchOpts, sink chan<- *KeepRelayBeaconRelayResetEvent) (event.Subscription, error) {

	logs, sub, err := _KeepRelayBeacon.contract.WatchLogs(opts, "RelayResetEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRelayBeaconRelayResetEvent)
				if err := _KeepRelayBeacon.contract.UnpackLog(event, "RelayResetEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// KeepRelayBeaconSubmitGroupPublicKeyEventIterator is returned from FilterSubmitGroupPublicKeyEvent and is used to iterate over the raw logs and unpacked data for SubmitGroupPublicKeyEvent events raised by the KeepRelayBeacon contract.
type KeepRelayBeaconSubmitGroupPublicKeyEventIterator struct {
	Event *KeepRelayBeaconSubmitGroupPublicKeyEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRelayBeaconSubmitGroupPublicKeyEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRelayBeaconSubmitGroupPublicKeyEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRelayBeaconSubmitGroupPublicKeyEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRelayBeaconSubmitGroupPublicKeyEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRelayBeaconSubmitGroupPublicKeyEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRelayBeaconSubmitGroupPublicKeyEvent represents a SubmitGroupPublicKeyEvent event raised by the KeepRelayBeacon contract.
type KeepRelayBeaconSubmitGroupPublicKeyEvent struct {
	PK_G_i                *big.Int
	Id                    *big.Int
	ActivationBlockHeight *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterSubmitGroupPublicKeyEvent is a free log retrieval operation binding the contract event 0x519ddcf35a3c14355496a7de6c423a2534e4e8653a4488178d86db2b7ccf01e2.
//
// Solidity: event SubmitGroupPublicKeyEvent(_PK_G_i uint256, _id uint256, _activationBlockHeight uint256)
func (_KeepRelayBeacon *KeepRelayBeaconFilterer) FilterSubmitGroupPublicKeyEvent(opts *bind.FilterOpts) (*KeepRelayBeaconSubmitGroupPublicKeyEventIterator, error) {

	logs, sub, err := _KeepRelayBeacon.contract.FilterLogs(opts, "SubmitGroupPublicKeyEvent")
	if err != nil {
		return nil, err
	}
	return &KeepRelayBeaconSubmitGroupPublicKeyEventIterator{contract: _KeepRelayBeacon.contract, event: "SubmitGroupPublicKeyEvent", logs: logs, sub: sub}, nil
}

// WatchSubmitGroupPublicKeyEvent is a free log subscription operation binding the contract event 0x519ddcf35a3c14355496a7de6c423a2534e4e8653a4488178d86db2b7ccf01e2.
//
// Solidity: event SubmitGroupPublicKeyEvent(_PK_G_i uint256, _id uint256, _activationBlockHeight uint256)
func (_KeepRelayBeacon *KeepRelayBeaconFilterer) WatchSubmitGroupPublicKeyEvent(opts *bind.WatchOpts, sink chan<- *KeepRelayBeaconSubmitGroupPublicKeyEvent) (event.Subscription, error) {

	logs, sub, err := _KeepRelayBeacon.contract.WatchLogs(opts, "SubmitGroupPublicKeyEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRelayBeaconSubmitGroupPublicKeyEvent)
				if err := _KeepRelayBeacon.contract.UnpackLog(event, "SubmitGroupPublicKeyEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
