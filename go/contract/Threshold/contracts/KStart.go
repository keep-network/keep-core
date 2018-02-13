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
const KStartABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"seed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"blockReward\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"payment\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"requestID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_RequestID\",\"type\":\"uint256\"},{\"name\":\"_Valid\",\"type\":\"bool\"},{\"name\":\"_theRandomNumber\",\"type\":\"uint256\"}],\"name\":\"randomNumberComplete\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"complete\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"number\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"isStaked\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRandomNumber\",\"outputs\":[{\"name\":\"status\",\"type\":\"uint256\"},{\"name\":\"theRandomNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_RequestID\",\"type\":\"uint256\"},{\"name\":\"_payment\",\"type\":\"uint256\"},{\"name\":\"_blockReward\",\"type\":\"uint256\"},{\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"requestRandomNumber\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"testName\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"RequestID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"payment\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"blockReward\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"seed\",\"type\":\"uint256\"}],\"name\":\"RelayRequestReady\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"RequestID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"Valid\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"RandomNumber\",\"type\":\"uint256\"}],\"name\":\"RandomNumberReady\",\"type\":\"event\"}]"

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

// BlockReward is a free data retrieval call binding the contract method 0x2fd6085e.
//
// Solidity: function blockReward( address) constant returns(uint256)
func (_KStart *KStartCaller) BlockReward(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "blockReward", arg0)
	return *ret0, err
}

// BlockReward is a free data retrieval call binding the contract method 0x2fd6085e.
//
// Solidity: function blockReward( address) constant returns(uint256)
func (_KStart *KStartSession) BlockReward(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.BlockReward(&_KStart.CallOpts, arg0)
}

// BlockReward is a free data retrieval call binding the contract method 0x2fd6085e.
//
// Solidity: function blockReward( address) constant returns(uint256)
func (_KStart *KStartCallerSession) BlockReward(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.BlockReward(&_KStart.CallOpts, arg0)
}

// Complete is a free data retrieval call binding the contract method 0x93af0292.
//
// Solidity: function complete( address) constant returns(uint256)
func (_KStart *KStartCaller) Complete(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "complete", arg0)
	return *ret0, err
}

// Complete is a free data retrieval call binding the contract method 0x93af0292.
//
// Solidity: function complete( address) constant returns(uint256)
func (_KStart *KStartSession) Complete(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.Complete(&_KStart.CallOpts, arg0)
}

// Complete is a free data retrieval call binding the contract method 0x93af0292.
//
// Solidity: function complete( address) constant returns(uint256)
func (_KStart *KStartCallerSession) Complete(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.Complete(&_KStart.CallOpts, arg0)
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() constant returns(status uint256, theRandomNumber uint256)
func (_KStart *KStartCaller) GetRandomNumber(opts *bind.CallOpts) (struct {
	Status          *big.Int
	TheRandomNumber *big.Int
}, error) {
	ret := new(struct {
		Status          *big.Int
		TheRandomNumber *big.Int
	})
	out := ret
	err := _KStart.contract.Call(opts, out, "getRandomNumber")
	return *ret, err
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() constant returns(status uint256, theRandomNumber uint256)
func (_KStart *KStartSession) GetRandomNumber() (struct {
	Status          *big.Int
	TheRandomNumber *big.Int
}, error) {
	return _KStart.Contract.GetRandomNumber(&_KStart.CallOpts)
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() constant returns(status uint256, theRandomNumber uint256)
func (_KStart *KStartCallerSession) GetRandomNumber() (struct {
	Status          *big.Int
	TheRandomNumber *big.Int
}, error) {
	return _KStart.Contract.GetRandomNumber(&_KStart.CallOpts)
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked( uint256) constant returns(bool)
func (_KStart *KStartCaller) IsStaked(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "isStaked", arg0)
	return *ret0, err
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked( uint256) constant returns(bool)
func (_KStart *KStartSession) IsStaked(arg0 *big.Int) (bool, error) {
	return _KStart.Contract.IsStaked(&_KStart.CallOpts, arg0)
}

// IsStaked is a free data retrieval call binding the contract method 0xbaa51f86.
//
// Solidity: function isStaked( uint256) constant returns(bool)
func (_KStart *KStartCallerSession) IsStaked(arg0 *big.Int) (bool, error) {
	return _KStart.Contract.IsStaked(&_KStart.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_KStart *KStartCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_KStart *KStartSession) Name() (string, error) {
	return _KStart.Contract.Name(&_KStart.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_KStart *KStartCallerSession) Name() (string, error) {
	return _KStart.Contract.Name(&_KStart.CallOpts)
}

// Number is a free data retrieval call binding the contract method 0xb154be85.
//
// Solidity: function number( address) constant returns(uint256)
func (_KStart *KStartCaller) Number(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "number", arg0)
	return *ret0, err
}

// Number is a free data retrieval call binding the contract method 0xb154be85.
//
// Solidity: function number( address) constant returns(uint256)
func (_KStart *KStartSession) Number(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.Number(&_KStart.CallOpts, arg0)
}

// Number is a free data retrieval call binding the contract method 0xb154be85.
//
// Solidity: function number( address) constant returns(uint256)
func (_KStart *KStartCallerSession) Number(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.Number(&_KStart.CallOpts, arg0)
}

// Payment is a free data retrieval call binding the contract method 0x3b92f3df.
//
// Solidity: function payment( address) constant returns(uint256)
func (_KStart *KStartCaller) Payment(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "payment", arg0)
	return *ret0, err
}

// Payment is a free data retrieval call binding the contract method 0x3b92f3df.
//
// Solidity: function payment( address) constant returns(uint256)
func (_KStart *KStartSession) Payment(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.Payment(&_KStart.CallOpts, arg0)
}

// Payment is a free data retrieval call binding the contract method 0x3b92f3df.
//
// Solidity: function payment( address) constant returns(uint256)
func (_KStart *KStartCallerSession) Payment(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.Payment(&_KStart.CallOpts, arg0)
}

// RequestID is a free data retrieval call binding the contract method 0x6ee1c916.
//
// Solidity: function requestID( address) constant returns(uint256)
func (_KStart *KStartCaller) RequestID(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "requestID", arg0)
	return *ret0, err
}

// RequestID is a free data retrieval call binding the contract method 0x6ee1c916.
//
// Solidity: function requestID( address) constant returns(uint256)
func (_KStart *KStartSession) RequestID(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.RequestID(&_KStart.CallOpts, arg0)
}

// RequestID is a free data retrieval call binding the contract method 0x6ee1c916.
//
// Solidity: function requestID( address) constant returns(uint256)
func (_KStart *KStartCallerSession) RequestID(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.RequestID(&_KStart.CallOpts, arg0)
}

// Seed is a free data retrieval call binding the contract method 0x063c67b7.
//
// Solidity: function seed( address) constant returns(uint256)
func (_KStart *KStartCaller) Seed(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KStart.contract.Call(opts, out, "seed", arg0)
	return *ret0, err
}

// Seed is a free data retrieval call binding the contract method 0x063c67b7.
//
// Solidity: function seed( address) constant returns(uint256)
func (_KStart *KStartSession) Seed(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.Seed(&_KStart.CallOpts, arg0)
}

// Seed is a free data retrieval call binding the contract method 0x063c67b7.
//
// Solidity: function seed( address) constant returns(uint256)
func (_KStart *KStartCallerSession) Seed(arg0 common.Address) (*big.Int, error) {
	return _KStart.Contract.Seed(&_KStart.CallOpts, arg0)
}

// RandomNumberComplete is a paid mutator transaction binding the contract method 0x774b25ad.
//
// Solidity: function randomNumberComplete(_RequestID uint256, _Valid bool, _theRandomNumber uint256) returns()
func (_KStart *KStartTransactor) RandomNumberComplete(opts *bind.TransactOpts, _RequestID *big.Int, _Valid bool, _theRandomNumber *big.Int) (*types.Transaction, error) {
	return _KStart.contract.Transact(opts, "randomNumberComplete", _RequestID, _Valid, _theRandomNumber)
}

// RandomNumberComplete is a paid mutator transaction binding the contract method 0x774b25ad.
//
// Solidity: function randomNumberComplete(_RequestID uint256, _Valid bool, _theRandomNumber uint256) returns()
func (_KStart *KStartSession) RandomNumberComplete(_RequestID *big.Int, _Valid bool, _theRandomNumber *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RandomNumberComplete(&_KStart.TransactOpts, _RequestID, _Valid, _theRandomNumber)
}

// RandomNumberComplete is a paid mutator transaction binding the contract method 0x774b25ad.
//
// Solidity: function randomNumberComplete(_RequestID uint256, _Valid bool, _theRandomNumber uint256) returns()
func (_KStart *KStartTransactorSession) RandomNumberComplete(_RequestID *big.Int, _Valid bool, _theRandomNumber *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RandomNumberComplete(&_KStart.TransactOpts, _RequestID, _Valid, _theRandomNumber)
}

// RequestRandomNumber is a paid mutator transaction binding the contract method 0xf74320e4.
//
// Solidity: function requestRandomNumber(_RequestID uint256, _payment uint256, _blockReward uint256, _seed uint256) returns()
func (_KStart *KStartTransactor) RequestRandomNumber(opts *bind.TransactOpts, _RequestID *big.Int, _payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KStart.contract.Transact(opts, "requestRandomNumber", _RequestID, _payment, _blockReward, _seed)
}

// RequestRandomNumber is a paid mutator transaction binding the contract method 0xf74320e4.
//
// Solidity: function requestRandomNumber(_RequestID uint256, _payment uint256, _blockReward uint256, _seed uint256) returns()
func (_KStart *KStartSession) RequestRandomNumber(_RequestID *big.Int, _payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RequestRandomNumber(&_KStart.TransactOpts, _RequestID, _payment, _blockReward, _seed)
}

// RequestRandomNumber is a paid mutator transaction binding the contract method 0xf74320e4.
//
// Solidity: function requestRandomNumber(_RequestID uint256, _payment uint256, _blockReward uint256, _seed uint256) returns()
func (_KStart *KStartTransactorSession) RequestRandomNumber(_RequestID *big.Int, _payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _KStart.Contract.RequestRandomNumber(&_KStart.TransactOpts, _RequestID, _payment, _blockReward, _seed)
}

// KStartRandomNumberReadyIterator is returned from FilterRandomNumberReady and is used to iterate over the raw logs and unpacked data for RandomNumberReady events raised by the KStart contract.
type KStartRandomNumberReadyIterator struct {
	Event *KStartRandomNumberReady // Event containing the contract specifics and raw log

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
func (it *KStartRandomNumberReadyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KStartRandomNumberReady)
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
		it.Event = new(KStartRandomNumberReady)
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
func (it *KStartRandomNumberReadyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KStartRandomNumberReadyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KStartRandomNumberReady represents a RandomNumberReady event raised by the KStart contract.
type KStartRandomNumberReady struct {
	RequestID    *big.Int
	Valid        bool
	RandomNumber *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRandomNumberReady is a free log retrieval operation binding the contract event 0x27cea675966d567c01a18e4d6e53abe51602de7cea1f854038ce6ef584338f5a.
//
// Solidity: event RandomNumberReady(RequestID uint256, Valid bool, RandomNumber uint256)
func (_KStart *KStartFilterer) FilterRandomNumberReady(opts *bind.FilterOpts) (*KStartRandomNumberReadyIterator, error) {

	logs, sub, err := _KStart.contract.FilterLogs(opts, "RandomNumberReady")
	if err != nil {
		return nil, err
	}
	return &KStartRandomNumberReadyIterator{contract: _KStart.contract, event: "RandomNumberReady", logs: logs, sub: sub}, nil
}

// WatchRandomNumberReady is a free log subscription operation binding the contract event 0x27cea675966d567c01a18e4d6e53abe51602de7cea1f854038ce6ef584338f5a.
//
// Solidity: event RandomNumberReady(RequestID uint256, Valid bool, RandomNumber uint256)
func (_KStart *KStartFilterer) WatchRandomNumberReady(opts *bind.WatchOpts, sink chan<- *KStartRandomNumberReady) (event.Subscription, error) {

	logs, sub, err := _KStart.contract.WatchLogs(opts, "RandomNumberReady")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KStartRandomNumberReady)
				if err := _KStart.contract.UnpackLog(event, "RandomNumberReady", log); err != nil {
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

// KStartRelayRequestReadyIterator is returned from FilterRelayRequestReady and is used to iterate over the raw logs and unpacked data for RelayRequestReady events raised by the KStart contract.
type KStartRelayRequestReadyIterator struct {
	Event *KStartRelayRequestReady // Event containing the contract specifics and raw log

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
func (it *KStartRelayRequestReadyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KStartRelayRequestReady)
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
		it.Event = new(KStartRelayRequestReady)
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
func (it *KStartRelayRequestReadyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KStartRelayRequestReadyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KStartRelayRequestReady represents a RelayRequestReady event raised by the KStart contract.
type KStartRelayRequestReady struct {
	RequestID   *big.Int
	Payment     *big.Int
	BlockReward *big.Int
	Seed        *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRelayRequestReady is a free log retrieval operation binding the contract event 0x74416580a5d574042861c6b2b7affa9581f905de2e358b973ef26eb550ab7a95.
//
// Solidity: event RelayRequestReady(RequestID uint256, payment uint256, blockReward uint256, seed uint256)
func (_KStart *KStartFilterer) FilterRelayRequestReady(opts *bind.FilterOpts) (*KStartRelayRequestReadyIterator, error) {

	logs, sub, err := _KStart.contract.FilterLogs(opts, "RelayRequestReady")
	if err != nil {
		return nil, err
	}
	return &KStartRelayRequestReadyIterator{contract: _KStart.contract, event: "RelayRequestReady", logs: logs, sub: sub}, nil
}

// WatchRelayRequestReady is a free log subscription operation binding the contract event 0x74416580a5d574042861c6b2b7affa9581f905de2e358b973ef26eb550ab7a95.
//
// Solidity: event RelayRequestReady(RequestID uint256, payment uint256, blockReward uint256, seed uint256)
func (_KStart *KStartFilterer) WatchRelayRequestReady(opts *bind.WatchOpts, sink chan<- *KStartRelayRequestReady) (event.Subscription, error) {

	logs, sub, err := _KStart.contract.WatchLogs(opts, "RelayRequestReady")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KStartRelayRequestReady)
				if err := _KStart.contract.UnpackLog(event, "RelayRequestReady", log); err != nil {
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
