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

// KeepRandomBeaconServiceABI is the input ABI used to generate the binding from.
const KeepRandomBeaconServiceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"UpgradeCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UpgradeStarted\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"adm\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"completeUpgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"initializationData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_newImplementation\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newAdmin\",\"type\":\"address\"}],\"name\":\"updateAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradeInitiatedTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_upgradeInitiatedTimestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradeTimeDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_upgradeTimeDelay\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KeepRandomBeaconService is an auto generated Go binding around an Ethereum contract.
type KeepRandomBeaconService struct {
	KeepRandomBeaconServiceCaller     // Read-only binding to the contract
	KeepRandomBeaconServiceTransactor // Write-only binding to the contract
	KeepRandomBeaconServiceFilterer   // Log filterer for contract events
}

// KeepRandomBeaconServiceCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeepRandomBeaconServiceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconServiceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeepRandomBeaconServiceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconServiceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeepRandomBeaconServiceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconServiceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeepRandomBeaconServiceSession struct {
	Contract     *KeepRandomBeaconService // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// KeepRandomBeaconServiceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeepRandomBeaconServiceCallerSession struct {
	Contract *KeepRandomBeaconServiceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// KeepRandomBeaconServiceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeepRandomBeaconServiceTransactorSession struct {
	Contract     *KeepRandomBeaconServiceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// KeepRandomBeaconServiceRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeepRandomBeaconServiceRaw struct {
	Contract *KeepRandomBeaconService // Generic contract binding to access the raw methods on
}

// KeepRandomBeaconServiceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeepRandomBeaconServiceCallerRaw struct {
	Contract *KeepRandomBeaconServiceCaller // Generic read-only contract binding to access the raw methods on
}

// KeepRandomBeaconServiceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeepRandomBeaconServiceTransactorRaw struct {
	Contract *KeepRandomBeaconServiceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeepRandomBeaconService creates a new instance of KeepRandomBeaconService, bound to a specific deployed contract.
func NewKeepRandomBeaconService(address common.Address, backend bind.ContractBackend) (*KeepRandomBeaconService, error) {
	contract, err := bindKeepRandomBeaconService(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconService{KeepRandomBeaconServiceCaller: KeepRandomBeaconServiceCaller{contract: contract}, KeepRandomBeaconServiceTransactor: KeepRandomBeaconServiceTransactor{contract: contract}, KeepRandomBeaconServiceFilterer: KeepRandomBeaconServiceFilterer{contract: contract}}, nil
}

// NewKeepRandomBeaconServiceCaller creates a new read-only instance of KeepRandomBeaconService, bound to a specific deployed contract.
func NewKeepRandomBeaconServiceCaller(address common.Address, caller bind.ContractCaller) (*KeepRandomBeaconServiceCaller, error) {
	contract, err := bindKeepRandomBeaconService(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceCaller{contract: contract}, nil
}

// NewKeepRandomBeaconServiceTransactor creates a new write-only instance of KeepRandomBeaconService, bound to a specific deployed contract.
func NewKeepRandomBeaconServiceTransactor(address common.Address, transactor bind.ContractTransactor) (*KeepRandomBeaconServiceTransactor, error) {
	contract, err := bindKeepRandomBeaconService(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceTransactor{contract: contract}, nil
}

// NewKeepRandomBeaconServiceFilterer creates a new log filterer instance of KeepRandomBeaconService, bound to a specific deployed contract.
func NewKeepRandomBeaconServiceFilterer(address common.Address, filterer bind.ContractFilterer) (*KeepRandomBeaconServiceFilterer, error) {
	contract, err := bindKeepRandomBeaconService(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceFilterer{contract: contract}, nil
}

// bindKeepRandomBeaconService binds a generic wrapper to an already deployed contract.
func bindKeepRandomBeaconService(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeepRandomBeaconServiceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRandomBeaconService *KeepRandomBeaconServiceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRandomBeaconService.Contract.KeepRandomBeaconServiceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRandomBeaconService *KeepRandomBeaconServiceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.KeepRandomBeaconServiceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRandomBeaconService *KeepRandomBeaconServiceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.KeepRandomBeaconServiceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRandomBeaconService.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRandomBeaconService *KeepRandomBeaconServiceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRandomBeaconService *KeepRandomBeaconServiceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address adm)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRandomBeaconService.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address adm)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) Admin() (common.Address, error) {
	return _KeepRandomBeaconService.Contract.Admin(&_KeepRandomBeaconService.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address adm)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCallerSession) Admin() (common.Address, error) {
	return _KeepRandomBeaconService.Contract.Admin(&_KeepRandomBeaconService.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRandomBeaconService.contract.Call(opts, out, "implementation")
	return *ret0, err
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) Implementation() (common.Address, error) {
	return _KeepRandomBeaconService.Contract.Implementation(&_KeepRandomBeaconService.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCallerSession) Implementation() (common.Address, error) {
	return _KeepRandomBeaconService.Contract.Implementation(&_KeepRandomBeaconService.CallOpts)
}

// InitializationData is a free data retrieval call binding the contract method 0x0f51d40d.
//
// Solidity: function initializationData(address ) constant returns(bytes)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCaller) InitializationData(opts *bind.CallOpts, arg0 common.Address) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _KeepRandomBeaconService.contract.Call(opts, out, "initializationData", arg0)
	return *ret0, err
}

// InitializationData is a free data retrieval call binding the contract method 0x0f51d40d.
//
// Solidity: function initializationData(address ) constant returns(bytes)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) InitializationData(arg0 common.Address) ([]byte, error) {
	return _KeepRandomBeaconService.Contract.InitializationData(&_KeepRandomBeaconService.CallOpts, arg0)
}

// InitializationData is a free data retrieval call binding the contract method 0x0f51d40d.
//
// Solidity: function initializationData(address ) constant returns(bytes)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCallerSession) InitializationData(arg0 common.Address) ([]byte, error) {
	return _KeepRandomBeaconService.Contract.InitializationData(&_KeepRandomBeaconService.CallOpts, arg0)
}

// NewImplementation is a free data retrieval call binding the contract method 0x8b677b03.
//
// Solidity: function newImplementation() constant returns(address _newImplementation)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCaller) NewImplementation(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRandomBeaconService.contract.Call(opts, out, "newImplementation")
	return *ret0, err
}

// NewImplementation is a free data retrieval call binding the contract method 0x8b677b03.
//
// Solidity: function newImplementation() constant returns(address _newImplementation)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) NewImplementation() (common.Address, error) {
	return _KeepRandomBeaconService.Contract.NewImplementation(&_KeepRandomBeaconService.CallOpts)
}

// NewImplementation is a free data retrieval call binding the contract method 0x8b677b03.
//
// Solidity: function newImplementation() constant returns(address _newImplementation)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCallerSession) NewImplementation() (common.Address, error) {
	return _KeepRandomBeaconService.Contract.NewImplementation(&_KeepRandomBeaconService.CallOpts)
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256 _upgradeInitiatedTimestamp)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCaller) UpgradeInitiatedTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconService.contract.Call(opts, out, "upgradeInitiatedTimestamp")
	return *ret0, err
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256 _upgradeInitiatedTimestamp)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) UpgradeInitiatedTimestamp() (*big.Int, error) {
	return _KeepRandomBeaconService.Contract.UpgradeInitiatedTimestamp(&_KeepRandomBeaconService.CallOpts)
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256 _upgradeInitiatedTimestamp)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCallerSession) UpgradeInitiatedTimestamp() (*big.Int, error) {
	return _KeepRandomBeaconService.Contract.UpgradeInitiatedTimestamp(&_KeepRandomBeaconService.CallOpts)
}

// UpgradeTimeDelay is a free data retrieval call binding the contract method 0x9260315b.
//
// Solidity: function upgradeTimeDelay() constant returns(uint256 _upgradeTimeDelay)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCaller) UpgradeTimeDelay(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconService.contract.Call(opts, out, "upgradeTimeDelay")
	return *ret0, err
}

// UpgradeTimeDelay is a free data retrieval call binding the contract method 0x9260315b.
//
// Solidity: function upgradeTimeDelay() constant returns(uint256 _upgradeTimeDelay)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) UpgradeTimeDelay() (*big.Int, error) {
	return _KeepRandomBeaconService.Contract.UpgradeTimeDelay(&_KeepRandomBeaconService.CallOpts)
}

// UpgradeTimeDelay is a free data retrieval call binding the contract method 0x9260315b.
//
// Solidity: function upgradeTimeDelay() constant returns(uint256 _upgradeTimeDelay)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceCallerSession) UpgradeTimeDelay() (*big.Int, error) {
	return _KeepRandomBeaconService.Contract.UpgradeTimeDelay(&_KeepRandomBeaconService.CallOpts)
}

// CompleteUpgrade is a paid mutator transaction binding the contract method 0x4d055667.
//
// Solidity: function completeUpgrade() returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceTransactor) CompleteUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconService.contract.Transact(opts, "completeUpgrade")
}

// CompleteUpgrade is a paid mutator transaction binding the contract method 0x4d055667.
//
// Solidity: function completeUpgrade() returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) CompleteUpgrade() (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.CompleteUpgrade(&_KeepRandomBeaconService.TransactOpts)
}

// CompleteUpgrade is a paid mutator transaction binding the contract method 0x4d055667.
//
// Solidity: function completeUpgrade() returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceTransactorSession) CompleteUpgrade() (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.CompleteUpgrade(&_KeepRandomBeaconService.TransactOpts)
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xe2f273bd.
//
// Solidity: function updateAdmin(address _newAdmin) returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceTransactor) UpdateAdmin(opts *bind.TransactOpts, _newAdmin common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconService.contract.Transact(opts, "updateAdmin", _newAdmin)
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xe2f273bd.
//
// Solidity: function updateAdmin(address _newAdmin) returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) UpdateAdmin(_newAdmin common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.UpdateAdmin(&_KeepRandomBeaconService.TransactOpts, _newAdmin)
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xe2f273bd.
//
// Solidity: function updateAdmin(address _newAdmin) returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceTransactorSession) UpdateAdmin(_newAdmin common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.UpdateAdmin(&_KeepRandomBeaconService.TransactOpts, _newAdmin)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x6fbc15e9.
//
// Solidity: function upgradeTo(address _newImplementation, bytes _data) returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceTransactor) UpgradeTo(opts *bind.TransactOpts, _newImplementation common.Address, _data []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconService.contract.Transact(opts, "upgradeTo", _newImplementation, _data)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x6fbc15e9.
//
// Solidity: function upgradeTo(address _newImplementation, bytes _data) returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceSession) UpgradeTo(_newImplementation common.Address, _data []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.UpgradeTo(&_KeepRandomBeaconService.TransactOpts, _newImplementation, _data)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x6fbc15e9.
//
// Solidity: function upgradeTo(address _newImplementation, bytes _data) returns()
func (_KeepRandomBeaconService *KeepRandomBeaconServiceTransactorSession) UpgradeTo(_newImplementation common.Address, _data []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconService.Contract.UpgradeTo(&_KeepRandomBeaconService.TransactOpts, _newImplementation, _data)
}

// KeepRandomBeaconServiceUpgradeCompletedIterator is returned from FilterUpgradeCompleted and is used to iterate over the raw logs and unpacked data for UpgradeCompleted events raised by the KeepRandomBeaconService contract.
type KeepRandomBeaconServiceUpgradeCompletedIterator struct {
	Event *KeepRandomBeaconServiceUpgradeCompleted // Event containing the contract specifics and raw log

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
func (it *KeepRandomBeaconServiceUpgradeCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconServiceUpgradeCompleted)
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
		it.Event = new(KeepRandomBeaconServiceUpgradeCompleted)
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
func (it *KeepRandomBeaconServiceUpgradeCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconServiceUpgradeCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconServiceUpgradeCompleted represents a UpgradeCompleted event raised by the KeepRandomBeaconService contract.
type KeepRandomBeaconServiceUpgradeCompleted struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgradeCompleted is a free log retrieval operation binding the contract event 0x2711ef8a6ba82dd467437054855e6edfc6424d0fb0127992b419c20f2d2e114e.
//
// Solidity: event UpgradeCompleted(address implementation)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceFilterer) FilterUpgradeCompleted(opts *bind.FilterOpts) (*KeepRandomBeaconServiceUpgradeCompletedIterator, error) {

	logs, sub, err := _KeepRandomBeaconService.contract.FilterLogs(opts, "UpgradeCompleted")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceUpgradeCompletedIterator{contract: _KeepRandomBeaconService.contract, event: "UpgradeCompleted", logs: logs, sub: sub}, nil
}

// WatchUpgradeCompleted is a free log subscription operation binding the contract event 0x2711ef8a6ba82dd467437054855e6edfc6424d0fb0127992b419c20f2d2e114e.
//
// Solidity: event UpgradeCompleted(address implementation)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceFilterer) WatchUpgradeCompleted(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconServiceUpgradeCompleted) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconService.contract.WatchLogs(opts, "UpgradeCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconServiceUpgradeCompleted)
				if err := _KeepRandomBeaconService.contract.UnpackLog(event, "UpgradeCompleted", log); err != nil {
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

// ParseUpgradeCompleted is a log parse operation binding the contract event 0x2711ef8a6ba82dd467437054855e6edfc6424d0fb0127992b419c20f2d2e114e.
//
// Solidity: event UpgradeCompleted(address implementation)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceFilterer) ParseUpgradeCompleted(log types.Log) (*KeepRandomBeaconServiceUpgradeCompleted, error) {
	event := new(KeepRandomBeaconServiceUpgradeCompleted)
	if err := _KeepRandomBeaconService.contract.UnpackLog(event, "UpgradeCompleted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconServiceUpgradeStartedIterator is returned from FilterUpgradeStarted and is used to iterate over the raw logs and unpacked data for UpgradeStarted events raised by the KeepRandomBeaconService contract.
type KeepRandomBeaconServiceUpgradeStartedIterator struct {
	Event *KeepRandomBeaconServiceUpgradeStarted // Event containing the contract specifics and raw log

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
func (it *KeepRandomBeaconServiceUpgradeStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconServiceUpgradeStarted)
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
		it.Event = new(KeepRandomBeaconServiceUpgradeStarted)
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
func (it *KeepRandomBeaconServiceUpgradeStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconServiceUpgradeStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconServiceUpgradeStarted represents a UpgradeStarted event raised by the KeepRandomBeaconService contract.
type KeepRandomBeaconServiceUpgradeStarted struct {
	Implementation common.Address
	Timestamp      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgradeStarted is a free log retrieval operation binding the contract event 0xe3ab72cd1d043aea9afd3c48aae4ee8f59a80381fec6c69ca502e2fb4b4f5c5a.
//
// Solidity: event UpgradeStarted(address implementation, uint256 timestamp)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceFilterer) FilterUpgradeStarted(opts *bind.FilterOpts) (*KeepRandomBeaconServiceUpgradeStartedIterator, error) {

	logs, sub, err := _KeepRandomBeaconService.contract.FilterLogs(opts, "UpgradeStarted")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceUpgradeStartedIterator{contract: _KeepRandomBeaconService.contract, event: "UpgradeStarted", logs: logs, sub: sub}, nil
}

// WatchUpgradeStarted is a free log subscription operation binding the contract event 0xe3ab72cd1d043aea9afd3c48aae4ee8f59a80381fec6c69ca502e2fb4b4f5c5a.
//
// Solidity: event UpgradeStarted(address implementation, uint256 timestamp)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceFilterer) WatchUpgradeStarted(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconServiceUpgradeStarted) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconService.contract.WatchLogs(opts, "UpgradeStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconServiceUpgradeStarted)
				if err := _KeepRandomBeaconService.contract.UnpackLog(event, "UpgradeStarted", log); err != nil {
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

// ParseUpgradeStarted is a log parse operation binding the contract event 0xe3ab72cd1d043aea9afd3c48aae4ee8f59a80381fec6c69ca502e2fb4b4f5c5a.
//
// Solidity: event UpgradeStarted(address implementation, uint256 timestamp)
func (_KeepRandomBeaconService *KeepRandomBeaconServiceFilterer) ParseUpgradeStarted(log types.Log) (*KeepRandomBeaconServiceUpgradeStarted, error) {
	event := new(KeepRandomBeaconServiceUpgradeStarted)
	if err := _KeepRandomBeaconService.contract.UnpackLog(event, "UpgradeStarted", log); err != nil {
		return nil, err
	}
	return event, nil
}
