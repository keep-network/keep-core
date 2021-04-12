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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PhasedEscrowABI is the input ABI used to generate the binding from.
const PhasedEscrowABI = "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_token\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"}],\"name\":\"BeneficiaryUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokensWithdrawn\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"beneficiary\",\"outputs\":[{\"internalType\":\"contractIBeneficiaryContract\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIBeneficiaryContract\",\"name\":\"_beneficiary\",\"type\":\"address\"}],\"name\":\"setBeneficiary\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractEscrow\",\"name\":\"_escrow\",\"type\":\"address\"}],\"name\":\"withdrawFromEscrow\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// PhasedEscrow is an auto generated Go binding around an Ethereum contract.
type PhasedEscrow struct {
	PhasedEscrowCaller     // Read-only binding to the contract
	PhasedEscrowTransactor // Write-only binding to the contract
	PhasedEscrowFilterer   // Log filterer for contract events
}

// PhasedEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type PhasedEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PhasedEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PhasedEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PhasedEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PhasedEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PhasedEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PhasedEscrowSession struct {
	Contract     *PhasedEscrow     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PhasedEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PhasedEscrowCallerSession struct {
	Contract *PhasedEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PhasedEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PhasedEscrowTransactorSession struct {
	Contract     *PhasedEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PhasedEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type PhasedEscrowRaw struct {
	Contract *PhasedEscrow // Generic contract binding to access the raw methods on
}

// PhasedEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PhasedEscrowCallerRaw struct {
	Contract *PhasedEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// PhasedEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PhasedEscrowTransactorRaw struct {
	Contract *PhasedEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPhasedEscrow creates a new instance of PhasedEscrow, bound to a specific deployed contract.
func NewPhasedEscrow(address common.Address, backend bind.ContractBackend) (*PhasedEscrow, error) {
	contract, err := bindPhasedEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PhasedEscrow{PhasedEscrowCaller: PhasedEscrowCaller{contract: contract}, PhasedEscrowTransactor: PhasedEscrowTransactor{contract: contract}, PhasedEscrowFilterer: PhasedEscrowFilterer{contract: contract}}, nil
}

// NewPhasedEscrowCaller creates a new read-only instance of PhasedEscrow, bound to a specific deployed contract.
func NewPhasedEscrowCaller(address common.Address, caller bind.ContractCaller) (*PhasedEscrowCaller, error) {
	contract, err := bindPhasedEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PhasedEscrowCaller{contract: contract}, nil
}

// NewPhasedEscrowTransactor creates a new write-only instance of PhasedEscrow, bound to a specific deployed contract.
func NewPhasedEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*PhasedEscrowTransactor, error) {
	contract, err := bindPhasedEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PhasedEscrowTransactor{contract: contract}, nil
}

// NewPhasedEscrowFilterer creates a new log filterer instance of PhasedEscrow, bound to a specific deployed contract.
func NewPhasedEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*PhasedEscrowFilterer, error) {
	contract, err := bindPhasedEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PhasedEscrowFilterer{contract: contract}, nil
}

// bindPhasedEscrow binds a generic wrapper to an already deployed contract.
func bindPhasedEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PhasedEscrowABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PhasedEscrow *PhasedEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PhasedEscrow.Contract.PhasedEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PhasedEscrow *PhasedEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.PhasedEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PhasedEscrow *PhasedEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.PhasedEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PhasedEscrow *PhasedEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PhasedEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PhasedEscrow *PhasedEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PhasedEscrow *PhasedEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.contract.Transact(opts, method, params...)
}

// Beneficiary is a free data retrieval call binding the contract method 0x38af3eed.
//
// Solidity: function beneficiary() view returns(address)
func (_PhasedEscrow *PhasedEscrowCaller) Beneficiary(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PhasedEscrow.contract.Call(opts, &out, "beneficiary")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Beneficiary is a free data retrieval call binding the contract method 0x38af3eed.
//
// Solidity: function beneficiary() view returns(address)
func (_PhasedEscrow *PhasedEscrowSession) Beneficiary() (common.Address, error) {
	return _PhasedEscrow.Contract.Beneficiary(&_PhasedEscrow.CallOpts)
}

// Beneficiary is a free data retrieval call binding the contract method 0x38af3eed.
//
// Solidity: function beneficiary() view returns(address)
func (_PhasedEscrow *PhasedEscrowCallerSession) Beneficiary() (common.Address, error) {
	return _PhasedEscrow.Contract.Beneficiary(&_PhasedEscrow.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PhasedEscrow *PhasedEscrowCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PhasedEscrow.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PhasedEscrow *PhasedEscrowSession) IsOwner() (bool, error) {
	return _PhasedEscrow.Contract.IsOwner(&_PhasedEscrow.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PhasedEscrow *PhasedEscrowCallerSession) IsOwner() (bool, error) {
	return _PhasedEscrow.Contract.IsOwner(&_PhasedEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PhasedEscrow *PhasedEscrowCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PhasedEscrow.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PhasedEscrow *PhasedEscrowSession) Owner() (common.Address, error) {
	return _PhasedEscrow.Contract.Owner(&_PhasedEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PhasedEscrow *PhasedEscrowCallerSession) Owner() (common.Address, error) {
	return _PhasedEscrow.Contract.Owner(&_PhasedEscrow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PhasedEscrow *PhasedEscrowCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PhasedEscrow.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PhasedEscrow *PhasedEscrowSession) Token() (common.Address, error) {
	return _PhasedEscrow.Contract.Token(&_PhasedEscrow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PhasedEscrow *PhasedEscrowCallerSession) Token() (common.Address, error) {
	return _PhasedEscrow.Contract.Token(&_PhasedEscrow.CallOpts)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_PhasedEscrow *PhasedEscrowTransactor) ReceiveApproval(opts *bind.TransactOpts, _from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _PhasedEscrow.contract.Transact(opts, "receiveApproval", _from, _value, _token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_PhasedEscrow *PhasedEscrowSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.ReceiveApproval(&_PhasedEscrow.TransactOpts, _from, _value, _token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_PhasedEscrow *PhasedEscrowTransactorSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.ReceiveApproval(&_PhasedEscrow.TransactOpts, _from, _value, _token, arg3)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PhasedEscrow *PhasedEscrowTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PhasedEscrow.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PhasedEscrow *PhasedEscrowSession) RenounceOwnership() (*types.Transaction, error) {
	return _PhasedEscrow.Contract.RenounceOwnership(&_PhasedEscrow.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PhasedEscrow *PhasedEscrowTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PhasedEscrow.Contract.RenounceOwnership(&_PhasedEscrow.TransactOpts)
}

// SetBeneficiary is a paid mutator transaction binding the contract method 0x1c31f710.
//
// Solidity: function setBeneficiary(address _beneficiary) returns()
func (_PhasedEscrow *PhasedEscrowTransactor) SetBeneficiary(opts *bind.TransactOpts, _beneficiary common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.contract.Transact(opts, "setBeneficiary", _beneficiary)
}

// SetBeneficiary is a paid mutator transaction binding the contract method 0x1c31f710.
//
// Solidity: function setBeneficiary(address _beneficiary) returns()
func (_PhasedEscrow *PhasedEscrowSession) SetBeneficiary(_beneficiary common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.SetBeneficiary(&_PhasedEscrow.TransactOpts, _beneficiary)
}

// SetBeneficiary is a paid mutator transaction binding the contract method 0x1c31f710.
//
// Solidity: function setBeneficiary(address _beneficiary) returns()
func (_PhasedEscrow *PhasedEscrowTransactorSession) SetBeneficiary(_beneficiary common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.SetBeneficiary(&_PhasedEscrow.TransactOpts, _beneficiary)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PhasedEscrow *PhasedEscrowTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PhasedEscrow *PhasedEscrowSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.TransferOwnership(&_PhasedEscrow.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PhasedEscrow *PhasedEscrowTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.TransferOwnership(&_PhasedEscrow.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_PhasedEscrow *PhasedEscrowTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _PhasedEscrow.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_PhasedEscrow *PhasedEscrowSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.Withdraw(&_PhasedEscrow.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_PhasedEscrow *PhasedEscrowTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.Withdraw(&_PhasedEscrow.TransactOpts, amount)
}

// WithdrawFromEscrow is a paid mutator transaction binding the contract method 0x924ab158.
//
// Solidity: function withdrawFromEscrow(address _escrow) returns()
func (_PhasedEscrow *PhasedEscrowTransactor) WithdrawFromEscrow(opts *bind.TransactOpts, _escrow common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.contract.Transact(opts, "withdrawFromEscrow", _escrow)
}

// WithdrawFromEscrow is a paid mutator transaction binding the contract method 0x924ab158.
//
// Solidity: function withdrawFromEscrow(address _escrow) returns()
func (_PhasedEscrow *PhasedEscrowSession) WithdrawFromEscrow(_escrow common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.WithdrawFromEscrow(&_PhasedEscrow.TransactOpts, _escrow)
}

// WithdrawFromEscrow is a paid mutator transaction binding the contract method 0x924ab158.
//
// Solidity: function withdrawFromEscrow(address _escrow) returns()
func (_PhasedEscrow *PhasedEscrowTransactorSession) WithdrawFromEscrow(_escrow common.Address) (*types.Transaction, error) {
	return _PhasedEscrow.Contract.WithdrawFromEscrow(&_PhasedEscrow.TransactOpts, _escrow)
}

// PhasedEscrowBeneficiaryUpdatedIterator is returned from FilterBeneficiaryUpdated and is used to iterate over the raw logs and unpacked data for BeneficiaryUpdated events raised by the PhasedEscrow contract.
type PhasedEscrowBeneficiaryUpdatedIterator struct {
	Event *PhasedEscrowBeneficiaryUpdated // Event containing the contract specifics and raw log

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
func (it *PhasedEscrowBeneficiaryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PhasedEscrowBeneficiaryUpdated)
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
		it.Event = new(PhasedEscrowBeneficiaryUpdated)
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
func (it *PhasedEscrowBeneficiaryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PhasedEscrowBeneficiaryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PhasedEscrowBeneficiaryUpdated represents a BeneficiaryUpdated event raised by the PhasedEscrow contract.
type PhasedEscrowBeneficiaryUpdated struct {
	Beneficiary common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBeneficiaryUpdated is a free log retrieval operation binding the contract event 0xeee59a71c694e68368a1cb0d135c448051bbfb12289e6c2223b0ceb100c2321d.
//
// Solidity: event BeneficiaryUpdated(address beneficiary)
func (_PhasedEscrow *PhasedEscrowFilterer) FilterBeneficiaryUpdated(opts *bind.FilterOpts) (*PhasedEscrowBeneficiaryUpdatedIterator, error) {

	logs, sub, err := _PhasedEscrow.contract.FilterLogs(opts, "BeneficiaryUpdated")
	if err != nil {
		return nil, err
	}
	return &PhasedEscrowBeneficiaryUpdatedIterator{contract: _PhasedEscrow.contract, event: "BeneficiaryUpdated", logs: logs, sub: sub}, nil
}

// WatchBeneficiaryUpdated is a free log subscription operation binding the contract event 0xeee59a71c694e68368a1cb0d135c448051bbfb12289e6c2223b0ceb100c2321d.
//
// Solidity: event BeneficiaryUpdated(address beneficiary)
func (_PhasedEscrow *PhasedEscrowFilterer) WatchBeneficiaryUpdated(opts *bind.WatchOpts, sink chan<- *PhasedEscrowBeneficiaryUpdated) (event.Subscription, error) {

	logs, sub, err := _PhasedEscrow.contract.WatchLogs(opts, "BeneficiaryUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PhasedEscrowBeneficiaryUpdated)
				if err := _PhasedEscrow.contract.UnpackLog(event, "BeneficiaryUpdated", log); err != nil {
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

// ParseBeneficiaryUpdated is a log parse operation binding the contract event 0xeee59a71c694e68368a1cb0d135c448051bbfb12289e6c2223b0ceb100c2321d.
//
// Solidity: event BeneficiaryUpdated(address beneficiary)
func (_PhasedEscrow *PhasedEscrowFilterer) ParseBeneficiaryUpdated(log types.Log) (*PhasedEscrowBeneficiaryUpdated, error) {
	event := new(PhasedEscrowBeneficiaryUpdated)
	if err := _PhasedEscrow.contract.UnpackLog(event, "BeneficiaryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PhasedEscrowOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PhasedEscrow contract.
type PhasedEscrowOwnershipTransferredIterator struct {
	Event *PhasedEscrowOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PhasedEscrowOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PhasedEscrowOwnershipTransferred)
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
		it.Event = new(PhasedEscrowOwnershipTransferred)
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
func (it *PhasedEscrowOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PhasedEscrowOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PhasedEscrowOwnershipTransferred represents a OwnershipTransferred event raised by the PhasedEscrow contract.
type PhasedEscrowOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PhasedEscrow *PhasedEscrowFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PhasedEscrowOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PhasedEscrow.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PhasedEscrowOwnershipTransferredIterator{contract: _PhasedEscrow.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PhasedEscrow *PhasedEscrowFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PhasedEscrowOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PhasedEscrow.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PhasedEscrowOwnershipTransferred)
				if err := _PhasedEscrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_PhasedEscrow *PhasedEscrowFilterer) ParseOwnershipTransferred(log types.Log) (*PhasedEscrowOwnershipTransferred, error) {
	event := new(PhasedEscrowOwnershipTransferred)
	if err := _PhasedEscrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PhasedEscrowTokensWithdrawnIterator is returned from FilterTokensWithdrawn and is used to iterate over the raw logs and unpacked data for TokensWithdrawn events raised by the PhasedEscrow contract.
type PhasedEscrowTokensWithdrawnIterator struct {
	Event *PhasedEscrowTokensWithdrawn // Event containing the contract specifics and raw log

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
func (it *PhasedEscrowTokensWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PhasedEscrowTokensWithdrawn)
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
		it.Event = new(PhasedEscrowTokensWithdrawn)
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
func (it *PhasedEscrowTokensWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PhasedEscrowTokensWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PhasedEscrowTokensWithdrawn represents a TokensWithdrawn event raised by the PhasedEscrow contract.
type PhasedEscrowTokensWithdrawn struct {
	Beneficiary common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokensWithdrawn is a free log retrieval operation binding the contract event 0x6352c5382c4a4578e712449ca65e83cdb392d045dfcf1cad9615189db2da244b.
//
// Solidity: event TokensWithdrawn(address beneficiary, uint256 amount)
func (_PhasedEscrow *PhasedEscrowFilterer) FilterTokensWithdrawn(opts *bind.FilterOpts) (*PhasedEscrowTokensWithdrawnIterator, error) {

	logs, sub, err := _PhasedEscrow.contract.FilterLogs(opts, "TokensWithdrawn")
	if err != nil {
		return nil, err
	}
	return &PhasedEscrowTokensWithdrawnIterator{contract: _PhasedEscrow.contract, event: "TokensWithdrawn", logs: logs, sub: sub}, nil
}

// WatchTokensWithdrawn is a free log subscription operation binding the contract event 0x6352c5382c4a4578e712449ca65e83cdb392d045dfcf1cad9615189db2da244b.
//
// Solidity: event TokensWithdrawn(address beneficiary, uint256 amount)
func (_PhasedEscrow *PhasedEscrowFilterer) WatchTokensWithdrawn(opts *bind.WatchOpts, sink chan<- *PhasedEscrowTokensWithdrawn) (event.Subscription, error) {

	logs, sub, err := _PhasedEscrow.contract.WatchLogs(opts, "TokensWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PhasedEscrowTokensWithdrawn)
				if err := _PhasedEscrow.contract.UnpackLog(event, "TokensWithdrawn", log); err != nil {
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

// ParseTokensWithdrawn is a log parse operation binding the contract event 0x6352c5382c4a4578e712449ca65e83cdb392d045dfcf1cad9615189db2da244b.
//
// Solidity: event TokensWithdrawn(address beneficiary, uint256 amount)
func (_PhasedEscrow *PhasedEscrowFilterer) ParseTokensWithdrawn(log types.Log) (*PhasedEscrowTokensWithdrawn, error) {
	event := new(PhasedEscrowTokensWithdrawn)
	if err := _PhasedEscrow.contract.UnpackLog(event, "TokensWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
