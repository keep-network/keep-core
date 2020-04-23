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

// ManagedGrantFactoryABI is the input ABI used to generate the binding from.
const ManagedGrantFactoryABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenGrant\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"grantAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"grantee\",\"type\":\"address\"}],\"name\":\"ManagedGrantCreated\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"grantee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cliff\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"policy\",\"type\":\"address\"}],\"name\":\"createManagedGrant\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_managedGrant\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractKeepToken\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenGrant\",\"outputs\":[{\"internalType\":\"contractTokenGrant\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ManagedGrantFactory is an auto generated Go binding around an Ethereum contract.
type ManagedGrantFactory struct {
	ManagedGrantFactoryCaller     // Read-only binding to the contract
	ManagedGrantFactoryTransactor // Write-only binding to the contract
	ManagedGrantFactoryFilterer   // Log filterer for contract events
}

// ManagedGrantFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ManagedGrantFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagedGrantFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ManagedGrantFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagedGrantFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ManagedGrantFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagedGrantFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ManagedGrantFactorySession struct {
	Contract     *ManagedGrantFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ManagedGrantFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ManagedGrantFactoryCallerSession struct {
	Contract *ManagedGrantFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ManagedGrantFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ManagedGrantFactoryTransactorSession struct {
	Contract     *ManagedGrantFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ManagedGrantFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ManagedGrantFactoryRaw struct {
	Contract *ManagedGrantFactory // Generic contract binding to access the raw methods on
}

// ManagedGrantFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ManagedGrantFactoryCallerRaw struct {
	Contract *ManagedGrantFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// ManagedGrantFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ManagedGrantFactoryTransactorRaw struct {
	Contract *ManagedGrantFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewManagedGrantFactory creates a new instance of ManagedGrantFactory, bound to a specific deployed contract.
func NewManagedGrantFactory(address common.Address, backend bind.ContractBackend) (*ManagedGrantFactory, error) {
	contract, err := bindManagedGrantFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ManagedGrantFactory{ManagedGrantFactoryCaller: ManagedGrantFactoryCaller{contract: contract}, ManagedGrantFactoryTransactor: ManagedGrantFactoryTransactor{contract: contract}, ManagedGrantFactoryFilterer: ManagedGrantFactoryFilterer{contract: contract}}, nil
}

// NewManagedGrantFactoryCaller creates a new read-only instance of ManagedGrantFactory, bound to a specific deployed contract.
func NewManagedGrantFactoryCaller(address common.Address, caller bind.ContractCaller) (*ManagedGrantFactoryCaller, error) {
	contract, err := bindManagedGrantFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ManagedGrantFactoryCaller{contract: contract}, nil
}

// NewManagedGrantFactoryTransactor creates a new write-only instance of ManagedGrantFactory, bound to a specific deployed contract.
func NewManagedGrantFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*ManagedGrantFactoryTransactor, error) {
	contract, err := bindManagedGrantFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ManagedGrantFactoryTransactor{contract: contract}, nil
}

// NewManagedGrantFactoryFilterer creates a new log filterer instance of ManagedGrantFactory, bound to a specific deployed contract.
func NewManagedGrantFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*ManagedGrantFactoryFilterer, error) {
	contract, err := bindManagedGrantFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ManagedGrantFactoryFilterer{contract: contract}, nil
}

// bindManagedGrantFactory binds a generic wrapper to an already deployed contract.
func bindManagedGrantFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ManagedGrantFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagedGrantFactory *ManagedGrantFactoryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ManagedGrantFactory.Contract.ManagedGrantFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagedGrantFactory *ManagedGrantFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagedGrantFactory.Contract.ManagedGrantFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagedGrantFactory *ManagedGrantFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagedGrantFactory.Contract.ManagedGrantFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagedGrantFactory *ManagedGrantFactoryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ManagedGrantFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagedGrantFactory *ManagedGrantFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagedGrantFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagedGrantFactory *ManagedGrantFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagedGrantFactory.Contract.contract.Transact(opts, method, params...)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_ManagedGrantFactory *ManagedGrantFactoryCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ManagedGrantFactory.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_ManagedGrantFactory *ManagedGrantFactorySession) Token() (common.Address, error) {
	return _ManagedGrantFactory.Contract.Token(&_ManagedGrantFactory.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_ManagedGrantFactory *ManagedGrantFactoryCallerSession) Token() (common.Address, error) {
	return _ManagedGrantFactory.Contract.Token(&_ManagedGrantFactory.CallOpts)
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_ManagedGrantFactory *ManagedGrantFactoryCaller) TokenGrant(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ManagedGrantFactory.contract.Call(opts, out, "tokenGrant")
	return *ret0, err
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_ManagedGrantFactory *ManagedGrantFactorySession) TokenGrant() (common.Address, error) {
	return _ManagedGrantFactory.Contract.TokenGrant(&_ManagedGrantFactory.CallOpts)
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_ManagedGrantFactory *ManagedGrantFactoryCallerSession) TokenGrant() (common.Address, error) {
	return _ManagedGrantFactory.Contract.TokenGrant(&_ManagedGrantFactory.CallOpts)
}

// CreateManagedGrant is a paid mutator transaction binding the contract method 0x319b84b7.
//
// Solidity: function createManagedGrant(address grantee, uint256 amount, uint256 duration, uint256 start, uint256 cliff, bool revocable, address policy) returns(address _managedGrant)
func (_ManagedGrantFactory *ManagedGrantFactoryTransactor) CreateManagedGrant(opts *bind.TransactOpts, grantee common.Address, amount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, revocable bool, policy common.Address) (*types.Transaction, error) {
	return _ManagedGrantFactory.contract.Transact(opts, "createManagedGrant", grantee, amount, duration, start, cliff, revocable, policy)
}

// CreateManagedGrant is a paid mutator transaction binding the contract method 0x319b84b7.
//
// Solidity: function createManagedGrant(address grantee, uint256 amount, uint256 duration, uint256 start, uint256 cliff, bool revocable, address policy) returns(address _managedGrant)
func (_ManagedGrantFactory *ManagedGrantFactorySession) CreateManagedGrant(grantee common.Address, amount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, revocable bool, policy common.Address) (*types.Transaction, error) {
	return _ManagedGrantFactory.Contract.CreateManagedGrant(&_ManagedGrantFactory.TransactOpts, grantee, amount, duration, start, cliff, revocable, policy)
}

// CreateManagedGrant is a paid mutator transaction binding the contract method 0x319b84b7.
//
// Solidity: function createManagedGrant(address grantee, uint256 amount, uint256 duration, uint256 start, uint256 cliff, bool revocable, address policy) returns(address _managedGrant)
func (_ManagedGrantFactory *ManagedGrantFactoryTransactorSession) CreateManagedGrant(grantee common.Address, amount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, revocable bool, policy common.Address) (*types.Transaction, error) {
	return _ManagedGrantFactory.Contract.CreateManagedGrant(&_ManagedGrantFactory.TransactOpts, grantee, amount, duration, start, cliff, revocable, policy)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _amount, address _token, bytes _extraData) returns()
func (_ManagedGrantFactory *ManagedGrantFactoryTransactor) ReceiveApproval(opts *bind.TransactOpts, _from common.Address, _amount *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _ManagedGrantFactory.contract.Transact(opts, "receiveApproval", _from, _amount, _token, _extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _amount, address _token, bytes _extraData) returns()
func (_ManagedGrantFactory *ManagedGrantFactorySession) ReceiveApproval(_from common.Address, _amount *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _ManagedGrantFactory.Contract.ReceiveApproval(&_ManagedGrantFactory.TransactOpts, _from, _amount, _token, _extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _amount, address _token, bytes _extraData) returns()
func (_ManagedGrantFactory *ManagedGrantFactoryTransactorSession) ReceiveApproval(_from common.Address, _amount *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _ManagedGrantFactory.Contract.ReceiveApproval(&_ManagedGrantFactory.TransactOpts, _from, _amount, _token, _extraData)
}

// ManagedGrantFactoryManagedGrantCreatedIterator is returned from FilterManagedGrantCreated and is used to iterate over the raw logs and unpacked data for ManagedGrantCreated events raised by the ManagedGrantFactory contract.
type ManagedGrantFactoryManagedGrantCreatedIterator struct {
	Event *ManagedGrantFactoryManagedGrantCreated // Event containing the contract specifics and raw log

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
func (it *ManagedGrantFactoryManagedGrantCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagedGrantFactoryManagedGrantCreated)
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
		it.Event = new(ManagedGrantFactoryManagedGrantCreated)
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
func (it *ManagedGrantFactoryManagedGrantCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagedGrantFactoryManagedGrantCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagedGrantFactoryManagedGrantCreated represents a ManagedGrantCreated event raised by the ManagedGrantFactory contract.
type ManagedGrantFactoryManagedGrantCreated struct {
	GrantAddress common.Address
	Grantee      common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterManagedGrantCreated is a free log retrieval operation binding the contract event 0xd670b8855e91e4ed32e465e934368080baed18d1aba14c28902bc9b934f3368a.
//
// Solidity: event ManagedGrantCreated(address grantAddress, address grantee)
func (_ManagedGrantFactory *ManagedGrantFactoryFilterer) FilterManagedGrantCreated(opts *bind.FilterOpts) (*ManagedGrantFactoryManagedGrantCreatedIterator, error) {

	logs, sub, err := _ManagedGrantFactory.contract.FilterLogs(opts, "ManagedGrantCreated")
	if err != nil {
		return nil, err
	}
	return &ManagedGrantFactoryManagedGrantCreatedIterator{contract: _ManagedGrantFactory.contract, event: "ManagedGrantCreated", logs: logs, sub: sub}, nil
}

// WatchManagedGrantCreated is a free log subscription operation binding the contract event 0xd670b8855e91e4ed32e465e934368080baed18d1aba14c28902bc9b934f3368a.
//
// Solidity: event ManagedGrantCreated(address grantAddress, address grantee)
func (_ManagedGrantFactory *ManagedGrantFactoryFilterer) WatchManagedGrantCreated(opts *bind.WatchOpts, sink chan<- *ManagedGrantFactoryManagedGrantCreated) (event.Subscription, error) {

	logs, sub, err := _ManagedGrantFactory.contract.WatchLogs(opts, "ManagedGrantCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagedGrantFactoryManagedGrantCreated)
				if err := _ManagedGrantFactory.contract.UnpackLog(event, "ManagedGrantCreated", log); err != nil {
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

// ParseManagedGrantCreated is a log parse operation binding the contract event 0xd670b8855e91e4ed32e465e934368080baed18d1aba14c28902bc9b934f3368a.
//
// Solidity: event ManagedGrantCreated(address grantAddress, address grantee)
func (_ManagedGrantFactory *ManagedGrantFactoryFilterer) ParseManagedGrantCreated(log types.Log) (*ManagedGrantFactoryManagedGrantCreated, error) {
	event := new(ManagedGrantFactoryManagedGrantCreated)
	if err := _ManagedGrantFactory.contract.UnpackLog(event, "ManagedGrantCreated", log); err != nil {
		return nil, err
	}
	return event, nil
}
