// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package KStart

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

// KStartABI is the input ABI used to generate the binding from.
const KStartABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"groupID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_payment\",\"type\":\"uint256\"},{\"name\":\"_blockReward\",\"type\":\"uint256\"},{\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"requestRelay\",\"outputs\":[{\"name\":\"RequestID\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"blockReward\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_RequestID\",\"type\":\"uint256\"},{\"name\":\"_groupSignature\",\"type\":\"uint256\"},{\"name\":\"_groupID\",\"type\":\"uint256\"},{\"name\":\"_previousEntry\",\"type\":\"uint256\"}],\"name\":\"relayEntry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_LastValidRelayTxHash\",\"type\":\"uint256\"},{\"name\":\"_LastValidRelayBlock\",\"type\":\"uint256\"}],\"name\":\"relayEntryAcustation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_PK_G_i\",\"type\":\"uint256\"},{\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"submitGroupPublicKey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"signature\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"payment\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"seed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_RequestID\",\"type\":\"uint256\"}],\"name\":\"getRandomNumber\",\"outputs\":[{\"name\":\"theRandomNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_UserPublicKey\",\"type\":\"uint256\"}],\"name\":\"isStaked\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"GenRequestIDSequence\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"RequestID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Payment\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"BlockReward\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Seed\",\"type\":\"uint256\"}],\"name\":\"RequestRelayEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"RequestID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Signature\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"GroupID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"PreviousEntry\",\"type\":\"uint256\"}],\"name\":\"RelayEntryEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"LastValidRelayEntry\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"LastValidRelayTxHash\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"LastValidRelayBlock\",\"type\":\"uint256\"}],\"name\":\"RelayResetEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_PK_G_i\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_id\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_activationBlockHeight\",\"type\":\"uint256\"}],\"name\":\"SubmitGroupPublicKeyEvent\",\"type\":\"event\"}]"

// KStart is an auto generated Go binding around an Ethereum contract.
type KStart struct {
	KStartCaller     // Read-only binding to the contract
	KStartTransactor // Write-only binding to the contract
	KStartFilterer   // Log filterer for contract events
}

// KStartCaller is an auto generated read-only Go binding around an Ethereum contract.
type KStartCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KStartTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KStartTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KStartFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KStartFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KStartSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KStartSession struct {
	Contract     *KStart           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KStartCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KStartCallerSession struct {
	Contract *KStartCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// KStartTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KStartTransactorSession struct {
	Contract     *KStartTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KStartRaw is an auto generated low-level Go binding around an Ethereum contract.
type KStartRaw struct {
	Contract *KStart // Generic contract binding to access the raw methods on
}

// KStartCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KStartCallerRaw struct {
	Contract *KStartCaller // Generic read-only contract binding to access the raw methods on
}

// KStartTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KStartTransactorRaw struct {
	Contract *KStartTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKStart creates a new instance of KStart, bound to a specific deployed contract.
func NewKStart(address common.Address, backend bind.ContractBackend) (*KStart, error) {
	contract, err := bindKStart(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KStart{KStartCaller: KStartCaller{contract: contract}, KStartTransactor: KStartTransactor{contract: contract}, KStartFilterer: KStartFilterer{contract: contract}}, nil
}

// NewKStartCaller creates a new read-only instance of KStart, bound to a specific deployed contract.
func NewKStartCaller(address common.Address, caller bind.ContractCaller) (*KStartCaller, error) {
	contract, err := bindKStart(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KStartCaller{contract: contract}, nil
}

// NewKStartTransactor creates a new write-only instance of KStart, bound to a specific deployed contract.
func NewKStartTransactor(address common.Address, transactor bind.ContractTransactor) (*KStartTransactor, error) {
	contract, err := bindKStart(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KStartTransactor{contract: contract}, nil
}

// NewKStartFilterer creates a new log filterer instance of KStart, bound to a specific deployed contract.
func NewKStartFilterer(address common.Address, filterer bind.ContractFilterer) (*KStartFilterer, error) {
	contract, err := bindKStart(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KStartFilterer{contract: contract}, nil
}

// bindKStart binds a generic wrapper to an already deployed contract.
func bindKStart(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KStartABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KStart *KStartRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KStart.Contract.KStartCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KStart *KStartRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KStart.Contract.KStartTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KStart *KStartRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KStart.Contract.KStartTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KStart *KStartCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KStart.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KStart *KStartTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KStart.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KStart *KStartTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KStart.Contract.contract.Transact(opts, method, params...)
}

// GenRequestIDSequence is a free data retrieval call binding the contract method 0xe205f3d7.
//
// Solidity: function GenRequestIDSequence() constant returns(address)
func (_KStart *KStartCaller) GenRequestIDSequence(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "GenRequestIDSequence")
	return *ret0, err
}

// GenRequestIDSequence is a free data retrieval call binding the contract method 0xe205f3d7.
//
// Solidity: function GenRequestIDSequence() constant returns(address)
func (_KStart *KStartSession) GenRequestIDSequence() (common.Address, error) {
	return _KStart.Contract.GenRequestIDSequence(&_KStart.CallOpts)
}

// GenRequestIDSequence is a free data retrieval call binding the contract method 0xe205f3d7.
//
// Solidity: function GenRequestIDSequence() constant returns(address)
func (_KStart *KStartCallerSession) GenRequestIDSequence() (common.Address, error) {
	return _KStart.Contract.GenRequestIDSequence(&_KStart.CallOpts)
}

// BlockReward is a free data retrieval call binding the contract method 0x25d1818b.
//
// Solidity: function blockReward( uint256) constant returns(uint256)
func (_KStart *KStartCaller) BlockReward(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "blockReward", arg0)
	return *ret0, err
}

// BlockReward is a free data retrieval call binding the contract method 0x25d1818b.
//
// Solidity: function blockReward( uint256) constant returns(uint256)
func (_KStart *KStartSession) BlockReward(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.BlockReward(&_KStart.CallOpts, arg0)
}

// BlockReward is a free data retrieval call binding the contract method 0x25d1818b.
//
// Solidity: function blockReward( uint256) constant returns(uint256)
func (_KStart *KStartCallerSession) BlockReward(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.BlockReward(&_KStart.CallOpts, arg0)
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xb37217a4.
//
// Solidity: function getRandomNumber(_RequestID uint256) constant returns(theRandomNumber uint256)
func (_KStart *KStartCaller) GetRandomNumber(opts *bind.CallOpts, _RequestID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "getRandomNumber", _RequestID)
	return *ret0, err
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xb37217a4.
//
// Solidity: function getRandomNumber(_RequestID uint256) constant returns(theRandomNumber uint256)
func (_KStart *KStartSession) GetRandomNumber(_RequestID *big.Int) (*big.Int, error) {
	return _KStart.Contract.GetRandomNumber(&_KStart.CallOpts, _RequestID)
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xb37217a4.
//
// Solidity: function getRandomNumber(_RequestID uint256) constant returns(theRandomNumber uint256)
func (_KStart *KStartCallerSession) GetRandomNumber(_RequestID *big.Int) (*big.Int, error) {
	return _KStart.Contract.GetRandomNumber(&_KStart.CallOpts, _RequestID)
}

// GroupID is a free data retrieval call binding the contract method 0x104e2471.
//
// Solidity: function groupID( uint256) constant returns(uint256)
func (_KStart *KStartCaller) GroupID(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "groupID", arg0)
	return *ret0, err
}

// GroupID is a free data retrieval call binding the contract method 0x104e2471.
//
// Solidity: function groupID( uint256) constant returns(uint256)
func (_KStart *KStartSession) GroupID(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.GroupID(&_KStart.CallOpts, arg0)
}

// GroupID is a free data retrieval call binding the contract method 0x104e2471.
//
// Solidity: function groupID( uint256) constant returns(uint256)
func (_KStart *KStartCallerSession) GroupID(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.GroupID(&_KStart.CallOpts, arg0)
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked(_UserPublicKey uint256) constant returns(bool)
func (_KStart *KStartCaller) IsStaked(opts *bind.CallOpts, _UserPublicKey *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "isStaked", _UserPublicKey)
	return *ret0, err
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked(_UserPublicKey uint256) constant returns(bool)
func (_KStart *KStartSession) IsStaked(_UserPublicKey *big.Int) (bool, error) {
	return _KStart.Contract.IsStaked(&_KStart.CallOpts, _UserPublicKey)
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked(_UserPublicKey uint256) constant returns(bool)
func (_KStart *KStartCallerSession) IsStaked(_UserPublicKey *big.Int) (bool, error) {
	return _KStart.Contract.IsStaked(&_KStart.CallOpts, _UserPublicKey)
}

// Payment is a free data retrieval call binding the contract method 0x8b3c99e3.
//
// Solidity: function payment( uint256) constant returns(uint256)
func (_KStart *KStartCaller) Payment(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "payment", arg0)
	return *ret0, err
}

// Payment is a free data retrieval call binding the contract method 0x8b3c99e3.
//
// Solidity: function payment( uint256) constant returns(uint256)
func (_KStart *KStartSession) Payment(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.Payment(&_KStart.CallOpts, arg0)
}

// Payment is a free data retrieval call binding the contract method 0x8b3c99e3.
//
// Solidity: function payment( uint256) constant returns(uint256)
func (_KStart *KStartCallerSession) Payment(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.Payment(&_KStart.CallOpts, arg0)
}

// Seed is a free data retrieval call binding the contract method 0x95564837.
//
// Solidity: function seed( uint256) constant returns(uint256)
func (_KStart *KStartCaller) Seed(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "seed", arg0)
	return *ret0, err
}

// Seed is a free data retrieval call binding the contract method 0x95564837.
//
// Solidity: function seed( uint256) constant returns(uint256)
func (_KStart *KStartSession) Seed(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.Seed(&_KStart.CallOpts, arg0)
}

// Seed is a free data retrieval call binding the contract method 0x95564837.
//
// Solidity: function seed( uint256) constant returns(uint256)
func (_KStart *KStartCallerSession) Seed(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.Seed(&_KStart.CallOpts, arg0)
}

// Signature is a free data retrieval call binding the contract method 0x70629548.
//
// Solidity: function signature( uint256) constant returns(uint256)
func (_KStart *KStartCaller) Signature(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "signature", arg0)
	return *ret0, err
}

// Signature is a free data retrieval call binding the contract method 0x70629548.
//
// Solidity: function signature( uint256) constant returns(uint256)
func (_KStart *KStartSession) Signature(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.Signature(&_KStart.CallOpts, arg0)
}

// Signature is a free data retrieval call binding the contract method 0x70629548.
//
// Solidity: function signature( uint256) constant returns(uint256)
func (_KStart *KStartCallerSession) Signature(arg0 *big.Int) (*big.Int, error) {
	return _KStart.Contract.Signature(&_KStart.CallOpts, arg0)
}

// RelayEntry is a paid mutator transaction binding the contract method 0x262e691d.
//
// Solidity: function relayEntry(_RequestID uint256, _groupSignature uint256, _groupID uint256, _previousEntry uint256) returns()
func (_KStart *KStartTransactor) RelayEntry(opts *bind.TransactOpts, _RequestID *big.Int, _groupSignature *big.Int, _groupID *big.Int, _previousEntry *big.Int) (*types.Transaction, error) {
	return _KStart.contract.Transact(opts, "relayEntry", _RequestID, _groupSignature, _groupID, _previousEntry)
}

// RelayEntry is a paid mutator transaction binding the contract method 0x262e691d.
//
// Solidity: function relayEntry(_RequestID uint256, _groupSignature uint256, _groupID uint256, _previousEntry uint256) returns()
func (_KStart *KStartSession) RelayEntry(_RequestID *big.Int, _groupSignature *big.Int, _groupID *big.Int, _previousEntry *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RelayEntry(&_KStart.TransactOpts, _RequestID, _groupSignature, _groupID, _previousEntry)
}

// RelayEntry is a paid mutator transaction binding the contract method 0x262e691d.
//
// Solidity: function relayEntry(_RequestID uint256, _groupSignature uint256, _groupID uint256, _previousEntry uint256) returns()
func (_KStart *KStartTransactorSession) RelayEntry(_RequestID *big.Int, _groupSignature *big.Int, _groupID *big.Int, _previousEntry *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RelayEntry(&_KStart.TransactOpts, _RequestID, _groupSignature, _groupID, _previousEntry)
}

// RelayEntryAcustation is a paid mutator transaction binding the contract method 0x4f3f9ad1.
//
// Solidity: function relayEntryAcustation(_LastValidRelayTxHash uint256, _LastValidRelayBlock uint256) returns()
func (_KStart *KStartTransactor) RelayEntryAcustation(opts *bind.TransactOpts, _LastValidRelayTxHash *big.Int, _LastValidRelayBlock *big.Int) (*types.Transaction, error) {
	return _KStart.contract.Transact(opts, "relayEntryAcustation", _LastValidRelayTxHash, _LastValidRelayBlock)
}

// RelayEntryAcustation is a paid mutator transaction binding the contract method 0x4f3f9ad1.
//
// Solidity: function relayEntryAcustation(_LastValidRelayTxHash uint256, _LastValidRelayBlock uint256) returns()
func (_KStart *KStartSession) RelayEntryAcustation(_LastValidRelayTxHash *big.Int, _LastValidRelayBlock *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RelayEntryAcustation(&_KStart.TransactOpts, _LastValidRelayTxHash, _LastValidRelayBlock)
}

// RelayEntryAcustation is a paid mutator transaction binding the contract method 0x4f3f9ad1.
//
// Solidity: function relayEntryAcustation(_LastValidRelayTxHash uint256, _LastValidRelayBlock uint256) returns()
func (_KStart *KStartTransactorSession) RelayEntryAcustation(_LastValidRelayTxHash *big.Int, _LastValidRelayBlock *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RelayEntryAcustation(&_KStart.TransactOpts, _LastValidRelayTxHash, _LastValidRelayBlock)
}

// RequestRelay is a paid mutator transaction binding the contract method 0x13df0019.
//
// Solidity: function requestRelay(_payment uint256, _blockReward uint256, _seed uint256) returns(RequestID uint256)
func (_KStart *KStartTransactor) RequestRelay(opts *bind.TransactOpts, _payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KStart.contract.Transact(opts, "requestRelay", _payment, _blockReward, _seed)
}

// RequestRelay is a paid mutator transaction binding the contract method 0x13df0019.
//
// Solidity: function requestRelay(_payment uint256, _blockReward uint256, _seed uint256) returns(RequestID uint256)
func (_KStart *KStartSession) RequestRelay(_payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RequestRelay(&_KStart.TransactOpts, _payment, _blockReward, _seed)
}

// RequestRelay is a paid mutator transaction binding the contract method 0x13df0019.
//
// Solidity: function requestRelay(_payment uint256, _blockReward uint256, _seed uint256) returns(RequestID uint256)
func (_KStart *KStartTransactorSession) RequestRelay(_payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RequestRelay(&_KStart.TransactOpts, _payment, _blockReward, _seed)
}

// SubmitGroupPublicKey is a paid mutator transaction binding the contract method 0x598b70ac.
//
// Solidity: function submitGroupPublicKey(_PK_G_i uint256, _id uint256) returns()
func (_KStart *KStartTransactor) SubmitGroupPublicKey(opts *bind.TransactOpts, _PK_G_i *big.Int, _id *big.Int) (*types.Transaction, error) {
	return _KStart.contract.Transact(opts, "submitGroupPublicKey", _PK_G_i, _id)
}

// SubmitGroupPublicKey is a paid mutator transaction binding the contract method 0x598b70ac.
//
// Solidity: function submitGroupPublicKey(_PK_G_i uint256, _id uint256) returns()
func (_KStart *KStartSession) SubmitGroupPublicKey(_PK_G_i *big.Int, _id *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.SubmitGroupPublicKey(&_KStart.TransactOpts, _PK_G_i, _id)
}

// SubmitGroupPublicKey is a paid mutator transaction binding the contract method 0x598b70ac.
//
// Solidity: function submitGroupPublicKey(_PK_G_i uint256, _id uint256) returns()
func (_KStart *KStartTransactorSession) SubmitGroupPublicKey(_PK_G_i *big.Int, _id *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.SubmitGroupPublicKey(&_KStart.TransactOpts, _PK_G_i, _id)
}

// KStartRelayEntryEventIterator is returned from FilterRelayEntryEvent and is used to iterate over the raw logs and unpacked data for RelayEntryEvent events raised by the KStart contract.
type KStartRelayEntryEventIterator struct {
	Event *KStartRelayEntryEvent // Event containing the contract specifics and raw log

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
func (it *KStartRelayEntryEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KStartRelayEntryEvent)
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
		it.Event = new(KStartRelayEntryEvent)
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
func (it *KStartRelayEntryEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KStartRelayEntryEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KStartRelayEntryEvent represents a RelayEntryEvent event raised by the KStart contract.
type KStartRelayEntryEvent struct {
	RequestID     *big.Int
	Signature     *big.Int
	GroupID       *big.Int
	PreviousEntry *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRelayEntryEvent is a free log retrieval operation binding the contract event 0x7b5431a25ea2792ba39a15c6b7ee8d6ef7aaa9bcc3d920ec3d997408a24b983b.
//
// Solidity: event RelayEntryEvent(RequestID uint256, Signature uint256, GroupID uint256, PreviousEntry uint256)
func (_KStart *KStartFilterer) FilterRelayEntryEvent(opts *bind.FilterOpts) (*KStartRelayEntryEventIterator, error) {

	logs, sub, err := _KStart.contract.FilterLogs(opts, "RelayEntryEvent")
	if err != nil {
		return nil, err
	}
	return &KStartRelayEntryEventIterator{contract: _KStart.contract, event: "RelayEntryEvent", logs: logs, sub: sub}, nil
}

// WatchRelayEntryEvent is a free log subscription operation binding the contract event 0x7b5431a25ea2792ba39a15c6b7ee8d6ef7aaa9bcc3d920ec3d997408a24b983b.
//
// Solidity: event RelayEntryEvent(RequestID uint256, Signature uint256, GroupID uint256, PreviousEntry uint256)
func (_KStart *KStartFilterer) WatchRelayEntryEvent(opts *bind.WatchOpts, sink chan<- *KStartRelayEntryEvent) (event.Subscription, error) {

	logs, sub, err := _KStart.contract.WatchLogs(opts, "RelayEntryEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KStartRelayEntryEvent)
				if err := _KStart.contract.UnpackLog(event, "RelayEntryEvent", log); err != nil {
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

// KStartRelayResetEventIterator is returned from FilterRelayResetEvent and is used to iterate over the raw logs and unpacked data for RelayResetEvent events raised by the KStart contract.
type KStartRelayResetEventIterator struct {
	Event *KStartRelayResetEvent // Event containing the contract specifics and raw log

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
func (it *KStartRelayResetEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KStartRelayResetEvent)
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
		it.Event = new(KStartRelayResetEvent)
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
func (it *KStartRelayResetEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KStartRelayResetEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KStartRelayResetEvent represents a RelayResetEvent event raised by the KStart contract.
type KStartRelayResetEvent struct {
	LastValidRelayEntry  *big.Int
	LastValidRelayTxHash *big.Int
	LastValidRelayBlock  *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterRelayResetEvent is a free log retrieval operation binding the contract event 0x242ff963387d59926325b194e1d237810c5c08836df9074589e6208087fa49d3.
//
// Solidity: event RelayResetEvent(LastValidRelayEntry uint256, LastValidRelayTxHash uint256, LastValidRelayBlock uint256)
func (_KStart *KStartFilterer) FilterRelayResetEvent(opts *bind.FilterOpts) (*KStartRelayResetEventIterator, error) {

	logs, sub, err := _KStart.contract.FilterLogs(opts, "RelayResetEvent")
	if err != nil {
		return nil, err
	}
	return &KStartRelayResetEventIterator{contract: _KStart.contract, event: "RelayResetEvent", logs: logs, sub: sub}, nil
}

// WatchRelayResetEvent is a free log subscription operation binding the contract event 0x242ff963387d59926325b194e1d237810c5c08836df9074589e6208087fa49d3.
//
// Solidity: event RelayResetEvent(LastValidRelayEntry uint256, LastValidRelayTxHash uint256, LastValidRelayBlock uint256)
func (_KStart *KStartFilterer) WatchRelayResetEvent(opts *bind.WatchOpts, sink chan<- *KStartRelayResetEvent) (event.Subscription, error) {

	logs, sub, err := _KStart.contract.WatchLogs(opts, "RelayResetEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KStartRelayResetEvent)
				if err := _KStart.contract.UnpackLog(event, "RelayResetEvent", log); err != nil {
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

// KStartRequestRelayEventIterator is returned from FilterRequestRelayEvent and is used to iterate over the raw logs and unpacked data for RequestRelayEvent events raised by the KStart contract.
type KStartRequestRelayEventIterator struct {
	Event *KStartRequestRelayEvent // Event containing the contract specifics and raw log

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
func (it *KStartRequestRelayEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KStartRequestRelayEvent)
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
		it.Event = new(KStartRequestRelayEvent)
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
func (it *KStartRequestRelayEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KStartRequestRelayEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KStartRequestRelayEvent represents a RequestRelayEvent event raised by the KStart contract.
type KStartRequestRelayEvent struct {
	RequestID   *big.Int
	Payment     *big.Int
	BlockReward *big.Int
	Seed        *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRequestRelayEvent is a free log retrieval operation binding the contract event 0xd0dcee08a4b291c512dfded29428b657b0a9db9625866be2b3c6cfbaaf77561a.
//
// Solidity: event RequestRelayEvent(RequestID uint256, Payment uint256, BlockReward uint256, Seed uint256)
func (_KStart *KStartFilterer) FilterRequestRelayEvent(opts *bind.FilterOpts) (*KStartRequestRelayEventIterator, error) {

	logs, sub, err := _KStart.contract.FilterLogs(opts, "RequestRelayEvent")
	if err != nil {
		return nil, err
	}
	return &KStartRequestRelayEventIterator{contract: _KStart.contract, event: "RequestRelayEvent", logs: logs, sub: sub}, nil
}

// WatchRequestRelayEvent is a free log subscription operation binding the contract event 0xd0dcee08a4b291c512dfded29428b657b0a9db9625866be2b3c6cfbaaf77561a.
//
// Solidity: event RequestRelayEvent(RequestID uint256, Payment uint256, BlockReward uint256, Seed uint256)
func (_KStart *KStartFilterer) WatchRequestRelayEvent(opts *bind.WatchOpts, sink chan<- *KStartRequestRelayEvent) (event.Subscription, error) {

	logs, sub, err := _KStart.contract.WatchLogs(opts, "RequestRelayEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KStartRequestRelayEvent)
				if err := _KStart.contract.UnpackLog(event, "RequestRelayEvent", log); err != nil {
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

// KStartSubmitGroupPublicKeyEventIterator is returned from FilterSubmitGroupPublicKeyEvent and is used to iterate over the raw logs and unpacked data for SubmitGroupPublicKeyEvent events raised by the KStart contract.
type KStartSubmitGroupPublicKeyEventIterator struct {
	Event *KStartSubmitGroupPublicKeyEvent // Event containing the contract specifics and raw log

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
func (it *KStartSubmitGroupPublicKeyEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KStartSubmitGroupPublicKeyEvent)
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
		it.Event = new(KStartSubmitGroupPublicKeyEvent)
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
func (it *KStartSubmitGroupPublicKeyEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KStartSubmitGroupPublicKeyEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KStartSubmitGroupPublicKeyEvent represents a SubmitGroupPublicKeyEvent event raised by the KStart contract.
type KStartSubmitGroupPublicKeyEvent struct {
	PK_G_i                *big.Int
	Id                    *big.Int
	ActivationBlockHeight *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterSubmitGroupPublicKeyEvent is a free log retrieval operation binding the contract event 0x519ddcf35a3c14355496a7de6c423a2534e4e8653a4488178d86db2b7ccf01e2.
//
// Solidity: event SubmitGroupPublicKeyEvent(_PK_G_i uint256, _id uint256, _activationBlockHeight uint256)
func (_KStart *KStartFilterer) FilterSubmitGroupPublicKeyEvent(opts *bind.FilterOpts) (*KStartSubmitGroupPublicKeyEventIterator, error) {

	logs, sub, err := _KStart.contract.FilterLogs(opts, "SubmitGroupPublicKeyEvent")
	if err != nil {
		return nil, err
	}
	return &KStartSubmitGroupPublicKeyEventIterator{contract: _KStart.contract, event: "SubmitGroupPublicKeyEvent", logs: logs, sub: sub}, nil
}

// WatchSubmitGroupPublicKeyEvent is a free log subscription operation binding the contract event 0x519ddcf35a3c14355496a7de6c423a2534e4e8653a4488178d86db2b7ccf01e2.
//
// Solidity: event SubmitGroupPublicKeyEvent(_PK_G_i uint256, _id uint256, _activationBlockHeight uint256)
func (_KStart *KStartFilterer) WatchSubmitGroupPublicKeyEvent(opts *bind.WatchOpts, sink chan<- *KStartSubmitGroupPublicKeyEvent) (event.Subscription, error) {

	logs, sub, err := _KStart.contract.WatchLogs(opts, "SubmitGroupPublicKeyEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KStartSubmitGroupPublicKeyEvent)
				if err := _KStart.contract.UnpackLog(event, "SubmitGroupPublicKeyEvent", log); err != nil {
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
