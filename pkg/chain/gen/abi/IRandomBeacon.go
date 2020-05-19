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

// IRandomBeaconABI is the input ABI used to generate the binding from.
const IRandomBeaconABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"entry\",\"type\":\"uint256\"}],\"name\":\"RelayEntryGenerated\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackGas\",\"type\":\"uint256\"}],\"name\":\"entryFeeEstimate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"requestRelayEntry\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"callbackGas\",\"type\":\"uint256\"}],\"name\":\"requestRelayEntry\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// IRandomBeacon is an auto generated Go binding around an Ethereum contract.
type IRandomBeacon struct {
	IRandomBeaconCaller     // Read-only binding to the contract
	IRandomBeaconTransactor // Write-only binding to the contract
	IRandomBeaconFilterer   // Log filterer for contract events
}

// IRandomBeaconCaller is an auto generated read-only Go binding around an Ethereum contract.
type IRandomBeaconCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRandomBeaconTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IRandomBeaconTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRandomBeaconFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IRandomBeaconFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRandomBeaconSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IRandomBeaconSession struct {
	Contract     *IRandomBeacon    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IRandomBeaconCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IRandomBeaconCallerSession struct {
	Contract *IRandomBeaconCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// IRandomBeaconTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IRandomBeaconTransactorSession struct {
	Contract     *IRandomBeaconTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IRandomBeaconRaw is an auto generated low-level Go binding around an Ethereum contract.
type IRandomBeaconRaw struct {
	Contract *IRandomBeacon // Generic contract binding to access the raw methods on
}

// IRandomBeaconCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IRandomBeaconCallerRaw struct {
	Contract *IRandomBeaconCaller // Generic read-only contract binding to access the raw methods on
}

// IRandomBeaconTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IRandomBeaconTransactorRaw struct {
	Contract *IRandomBeaconTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIRandomBeacon creates a new instance of IRandomBeacon, bound to a specific deployed contract.
func NewIRandomBeacon(address common.Address, backend bind.ContractBackend) (*IRandomBeacon, error) {
	contract, err := bindIRandomBeacon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IRandomBeacon{IRandomBeaconCaller: IRandomBeaconCaller{contract: contract}, IRandomBeaconTransactor: IRandomBeaconTransactor{contract: contract}, IRandomBeaconFilterer: IRandomBeaconFilterer{contract: contract}}, nil
}

// NewIRandomBeaconCaller creates a new read-only instance of IRandomBeacon, bound to a specific deployed contract.
func NewIRandomBeaconCaller(address common.Address, caller bind.ContractCaller) (*IRandomBeaconCaller, error) {
	contract, err := bindIRandomBeacon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IRandomBeaconCaller{contract: contract}, nil
}

// NewIRandomBeaconTransactor creates a new write-only instance of IRandomBeacon, bound to a specific deployed contract.
func NewIRandomBeaconTransactor(address common.Address, transactor bind.ContractTransactor) (*IRandomBeaconTransactor, error) {
	contract, err := bindIRandomBeacon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IRandomBeaconTransactor{contract: contract}, nil
}

// NewIRandomBeaconFilterer creates a new log filterer instance of IRandomBeacon, bound to a specific deployed contract.
func NewIRandomBeaconFilterer(address common.Address, filterer bind.ContractFilterer) (*IRandomBeaconFilterer, error) {
	contract, err := bindIRandomBeacon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IRandomBeaconFilterer{contract: contract}, nil
}

// bindIRandomBeacon binds a generic wrapper to an already deployed contract.
func bindIRandomBeacon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IRandomBeaconABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRandomBeacon *IRandomBeaconRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IRandomBeacon.Contract.IRandomBeaconCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRandomBeacon *IRandomBeaconRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRandomBeacon.Contract.IRandomBeaconTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRandomBeacon *IRandomBeaconRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRandomBeacon.Contract.IRandomBeaconTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRandomBeacon *IRandomBeaconCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IRandomBeacon.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRandomBeacon *IRandomBeaconTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRandomBeacon.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRandomBeacon *IRandomBeaconTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRandomBeacon.Contract.contract.Transact(opts, method, params...)
}

// EntryFeeEstimate is a free data retrieval call binding the contract method 0xd13f1391.
//
// Solidity: function entryFeeEstimate(uint256 callbackGas) constant returns(uint256)
func (_IRandomBeacon *IRandomBeaconCaller) EntryFeeEstimate(opts *bind.CallOpts, callbackGas *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IRandomBeacon.contract.Call(opts, out, "entryFeeEstimate", callbackGas)
	return *ret0, err
}

// EntryFeeEstimate is a free data retrieval call binding the contract method 0xd13f1391.
//
// Solidity: function entryFeeEstimate(uint256 callbackGas) constant returns(uint256)
func (_IRandomBeacon *IRandomBeaconSession) EntryFeeEstimate(callbackGas *big.Int) (*big.Int, error) {
	return _IRandomBeacon.Contract.EntryFeeEstimate(&_IRandomBeacon.CallOpts, callbackGas)
}

// EntryFeeEstimate is a free data retrieval call binding the contract method 0xd13f1391.
//
// Solidity: function entryFeeEstimate(uint256 callbackGas) constant returns(uint256)
func (_IRandomBeacon *IRandomBeaconCallerSession) EntryFeeEstimate(callbackGas *big.Int) (*big.Int, error) {
	return _IRandomBeacon.Contract.EntryFeeEstimate(&_IRandomBeacon.CallOpts, callbackGas)
}

// RequestRelayEntry is a paid mutator transaction binding the contract method 0x2d53fe8b.
//
// Solidity: function requestRelayEntry() returns(uint256)
func (_IRandomBeacon *IRandomBeaconTransactor) RequestRelayEntry(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRandomBeacon.contract.Transact(opts, "requestRelayEntry")
}

// RequestRelayEntry is a paid mutator transaction binding the contract method 0x2d53fe8b.
//
// Solidity: function requestRelayEntry() returns(uint256)
func (_IRandomBeacon *IRandomBeaconSession) RequestRelayEntry() (*types.Transaction, error) {
	return _IRandomBeacon.Contract.RequestRelayEntry(&_IRandomBeacon.TransactOpts)
}

// RequestRelayEntry is a paid mutator transaction binding the contract method 0x2d53fe8b.
//
// Solidity: function requestRelayEntry() returns(uint256)
func (_IRandomBeacon *IRandomBeaconTransactorSession) RequestRelayEntry() (*types.Transaction, error) {
	return _IRandomBeacon.Contract.RequestRelayEntry(&_IRandomBeacon.TransactOpts)
}

// RequestRelayEntry0 is a paid mutator transaction binding the contract method 0xe1f95890.
//
// Solidity: function requestRelayEntry(address callbackContract, uint256 callbackGas) returns(uint256)
func (_IRandomBeacon *IRandomBeaconTransactor) RequestRelayEntry0(opts *bind.TransactOpts, callbackContract common.Address, callbackGas *big.Int) (*types.Transaction, error) {
	return _IRandomBeacon.contract.Transact(opts, "requestRelayEntry0", callbackContract, callbackGas)
}

// RequestRelayEntry0 is a paid mutator transaction binding the contract method 0xe1f95890.
//
// Solidity: function requestRelayEntry(address callbackContract, uint256 callbackGas) returns(uint256)
func (_IRandomBeacon *IRandomBeaconSession) RequestRelayEntry0(callbackContract common.Address, callbackGas *big.Int) (*types.Transaction, error) {
	return _IRandomBeacon.Contract.RequestRelayEntry0(&_IRandomBeacon.TransactOpts, callbackContract, callbackGas)
}

// RequestRelayEntry0 is a paid mutator transaction binding the contract method 0xe1f95890.
//
// Solidity: function requestRelayEntry(address callbackContract, uint256 callbackGas) returns(uint256)
func (_IRandomBeacon *IRandomBeaconTransactorSession) RequestRelayEntry0(callbackContract common.Address, callbackGas *big.Int) (*types.Transaction, error) {
	return _IRandomBeacon.Contract.RequestRelayEntry0(&_IRandomBeacon.TransactOpts, callbackContract, callbackGas)
}

// IRandomBeaconRelayEntryGeneratedIterator is returned from FilterRelayEntryGenerated and is used to iterate over the raw logs and unpacked data for RelayEntryGenerated events raised by the IRandomBeacon contract.
type IRandomBeaconRelayEntryGeneratedIterator struct {
	Event *IRandomBeaconRelayEntryGenerated // Event containing the contract specifics and raw log

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
func (it *IRandomBeaconRelayEntryGeneratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IRandomBeaconRelayEntryGenerated)
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
		it.Event = new(IRandomBeaconRelayEntryGenerated)
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
func (it *IRandomBeaconRelayEntryGeneratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IRandomBeaconRelayEntryGeneratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IRandomBeaconRelayEntryGenerated represents a RelayEntryGenerated event raised by the IRandomBeacon contract.
type IRandomBeaconRelayEntryGenerated struct {
	RequestId *big.Int
	Entry     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRelayEntryGenerated is a free log retrieval operation binding the contract event 0xddb7473cb5eaf3aaa2cf09ea9d573d4876aac471614f99773f32a4b5994ea7e1.
//
// Solidity: event RelayEntryGenerated(uint256 requestId, uint256 entry)
func (_IRandomBeacon *IRandomBeaconFilterer) FilterRelayEntryGenerated(opts *bind.FilterOpts) (*IRandomBeaconRelayEntryGeneratedIterator, error) {

	logs, sub, err := _IRandomBeacon.contract.FilterLogs(opts, "RelayEntryGenerated")
	if err != nil {
		return nil, err
	}
	return &IRandomBeaconRelayEntryGeneratedIterator{contract: _IRandomBeacon.contract, event: "RelayEntryGenerated", logs: logs, sub: sub}, nil
}

// WatchRelayEntryGenerated is a free log subscription operation binding the contract event 0xddb7473cb5eaf3aaa2cf09ea9d573d4876aac471614f99773f32a4b5994ea7e1.
//
// Solidity: event RelayEntryGenerated(uint256 requestId, uint256 entry)
func (_IRandomBeacon *IRandomBeaconFilterer) WatchRelayEntryGenerated(opts *bind.WatchOpts, sink chan<- *IRandomBeaconRelayEntryGenerated) (event.Subscription, error) {

	logs, sub, err := _IRandomBeacon.contract.WatchLogs(opts, "RelayEntryGenerated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IRandomBeaconRelayEntryGenerated)
				if err := _IRandomBeacon.contract.UnpackLog(event, "RelayEntryGenerated", log); err != nil {
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

// ParseRelayEntryGenerated is a log parse operation binding the contract event 0xddb7473cb5eaf3aaa2cf09ea9d573d4876aac471614f99773f32a4b5994ea7e1.
//
// Solidity: event RelayEntryGenerated(uint256 requestId, uint256 entry)
func (_IRandomBeacon *IRandomBeaconFilterer) ParseRelayEntryGenerated(log types.Log) (*IRandomBeaconRelayEntryGenerated, error) {
	event := new(IRandomBeaconRelayEntryGenerated)
	if err := _IRandomBeacon.contract.UnpackLog(event, "RelayEntryGenerated", log); err != nil {
		return nil, err
	}
	return event, nil
}
