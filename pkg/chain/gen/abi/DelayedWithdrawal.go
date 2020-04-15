// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// DelayedWithdrawalABI is the input ABI used to generate the binding from.
const DelayedWithdrawalABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"finishWithdrawal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"initiateWithdrawal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// DelayedWithdrawal is an auto generated Go binding around an Ethereum contract.
type DelayedWithdrawal struct {
	DelayedWithdrawalCaller     // Read-only binding to the contract
	DelayedWithdrawalTransactor // Write-only binding to the contract
	DelayedWithdrawalFilterer   // Log filterer for contract events
}

// DelayedWithdrawalCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelayedWithdrawalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelayedWithdrawalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelayedWithdrawalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelayedWithdrawalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelayedWithdrawalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelayedWithdrawalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelayedWithdrawalSession struct {
	Contract     *DelayedWithdrawal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// DelayedWithdrawalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelayedWithdrawalCallerSession struct {
	Contract *DelayedWithdrawalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// DelayedWithdrawalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelayedWithdrawalTransactorSession struct {
	Contract     *DelayedWithdrawalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// DelayedWithdrawalRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelayedWithdrawalRaw struct {
	Contract *DelayedWithdrawal // Generic contract binding to access the raw methods on
}

// DelayedWithdrawalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelayedWithdrawalCallerRaw struct {
	Contract *DelayedWithdrawalCaller // Generic read-only contract binding to access the raw methods on
}

// DelayedWithdrawalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelayedWithdrawalTransactorRaw struct {
	Contract *DelayedWithdrawalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelayedWithdrawal creates a new instance of DelayedWithdrawal, bound to a specific deployed contract.
func NewDelayedWithdrawal(address common.Address, backend bind.ContractBackend) (*DelayedWithdrawal, error) {
	contract, err := bindDelayedWithdrawal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DelayedWithdrawal{DelayedWithdrawalCaller: DelayedWithdrawalCaller{contract: contract}, DelayedWithdrawalTransactor: DelayedWithdrawalTransactor{contract: contract}, DelayedWithdrawalFilterer: DelayedWithdrawalFilterer{contract: contract}}, nil
}

// NewDelayedWithdrawalCaller creates a new read-only instance of DelayedWithdrawal, bound to a specific deployed contract.
func NewDelayedWithdrawalCaller(address common.Address, caller bind.ContractCaller) (*DelayedWithdrawalCaller, error) {
	contract, err := bindDelayedWithdrawal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelayedWithdrawalCaller{contract: contract}, nil
}

// NewDelayedWithdrawalTransactor creates a new write-only instance of DelayedWithdrawal, bound to a specific deployed contract.
func NewDelayedWithdrawalTransactor(address common.Address, transactor bind.ContractTransactor) (*DelayedWithdrawalTransactor, error) {
	contract, err := bindDelayedWithdrawal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelayedWithdrawalTransactor{contract: contract}, nil
}

// NewDelayedWithdrawalFilterer creates a new log filterer instance of DelayedWithdrawal, bound to a specific deployed contract.
func NewDelayedWithdrawalFilterer(address common.Address, filterer bind.ContractFilterer) (*DelayedWithdrawalFilterer, error) {
	contract, err := bindDelayedWithdrawal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelayedWithdrawalFilterer{contract: contract}, nil
}

// bindDelayedWithdrawal binds a generic wrapper to an already deployed contract.
func bindDelayedWithdrawal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DelayedWithdrawalABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelayedWithdrawal *DelayedWithdrawalRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DelayedWithdrawal.Contract.DelayedWithdrawalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelayedWithdrawal *DelayedWithdrawalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.DelayedWithdrawalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelayedWithdrawal *DelayedWithdrawalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.DelayedWithdrawalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelayedWithdrawal *DelayedWithdrawalCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DelayedWithdrawal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelayedWithdrawal *DelayedWithdrawalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelayedWithdrawal *DelayedWithdrawalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DelayedWithdrawal *DelayedWithdrawalCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DelayedWithdrawal.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DelayedWithdrawal *DelayedWithdrawalSession) IsOwner() (bool, error) {
	return _DelayedWithdrawal.Contract.IsOwner(&_DelayedWithdrawal.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DelayedWithdrawal *DelayedWithdrawalCallerSession) IsOwner() (bool, error) {
	return _DelayedWithdrawal.Contract.IsOwner(&_DelayedWithdrawal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DelayedWithdrawal *DelayedWithdrawalCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DelayedWithdrawal.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DelayedWithdrawal *DelayedWithdrawalSession) Owner() (common.Address, error) {
	return _DelayedWithdrawal.Contract.Owner(&_DelayedWithdrawal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DelayedWithdrawal *DelayedWithdrawalCallerSession) Owner() (common.Address, error) {
	return _DelayedWithdrawal.Contract.Owner(&_DelayedWithdrawal.CallOpts)
}

// FinishWithdrawal is a paid mutator transaction binding the contract method 0xa5150bbc.
//
// Solidity: function finishWithdrawal(address payee) returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactor) FinishWithdrawal(opts *bind.TransactOpts, payee common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.contract.Transact(opts, "finishWithdrawal", payee)
}

// FinishWithdrawal is a paid mutator transaction binding the contract method 0xa5150bbc.
//
// Solidity: function finishWithdrawal(address payee) returns()
func (_DelayedWithdrawal *DelayedWithdrawalSession) FinishWithdrawal(payee common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.FinishWithdrawal(&_DelayedWithdrawal.TransactOpts, payee)
}

// FinishWithdrawal is a paid mutator transaction binding the contract method 0xa5150bbc.
//
// Solidity: function finishWithdrawal(address payee) returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactorSession) FinishWithdrawal(payee common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.FinishWithdrawal(&_DelayedWithdrawal.TransactOpts, payee)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address sender) returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactor) Initialize(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.contract.Transact(opts, "initialize", sender)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address sender) returns()
func (_DelayedWithdrawal *DelayedWithdrawalSession) Initialize(sender common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.Initialize(&_DelayedWithdrawal.TransactOpts, sender)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address sender) returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactorSession) Initialize(sender common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.Initialize(&_DelayedWithdrawal.TransactOpts, sender)
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactor) InitiateWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelayedWithdrawal.contract.Transact(opts, "initiateWithdrawal")
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns()
func (_DelayedWithdrawal *DelayedWithdrawalSession) InitiateWithdrawal() (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.InitiateWithdrawal(&_DelayedWithdrawal.TransactOpts)
}

// InitiateWithdrawal is a paid mutator transaction binding the contract method 0xb51d1d4f.
//
// Solidity: function initiateWithdrawal() returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactorSession) InitiateWithdrawal() (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.InitiateWithdrawal(&_DelayedWithdrawal.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelayedWithdrawal.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DelayedWithdrawal *DelayedWithdrawalSession) RenounceOwnership() (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.RenounceOwnership(&_DelayedWithdrawal.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.RenounceOwnership(&_DelayedWithdrawal.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DelayedWithdrawal *DelayedWithdrawalSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.TransferOwnership(&_DelayedWithdrawal.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DelayedWithdrawal *DelayedWithdrawalTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DelayedWithdrawal.Contract.TransferOwnership(&_DelayedWithdrawal.TransactOpts, newOwner)
}

// DelayedWithdrawalOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DelayedWithdrawal contract.
type DelayedWithdrawalOwnershipTransferredIterator struct {
	Event *DelayedWithdrawalOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DelayedWithdrawalOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelayedWithdrawalOwnershipTransferred)
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
		it.Event = new(DelayedWithdrawalOwnershipTransferred)
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
func (it *DelayedWithdrawalOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelayedWithdrawalOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelayedWithdrawalOwnershipTransferred represents a OwnershipTransferred event raised by the DelayedWithdrawal contract.
type DelayedWithdrawalOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DelayedWithdrawal *DelayedWithdrawalFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DelayedWithdrawalOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DelayedWithdrawal.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DelayedWithdrawalOwnershipTransferredIterator{contract: _DelayedWithdrawal.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DelayedWithdrawal *DelayedWithdrawalFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DelayedWithdrawalOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DelayedWithdrawal.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelayedWithdrawalOwnershipTransferred)
				if err := _DelayedWithdrawal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DelayedWithdrawal *DelayedWithdrawalFilterer) ParseOwnershipTransferred(log types.Log) (*DelayedWithdrawalOwnershipTransferred, error) {
	event := new(DelayedWithdrawalOwnershipTransferred)
	if err := _DelayedWithdrawal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
