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

// ManagedGrantABI is the input ABI used to generate the binding from.
const ManagedGrantABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenGrant\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_grantManager\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_grantId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_grantee\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"cancelledRequestedGrantee\",\"type\":\"address\"}],\"name\":\"GranteeReassignmentCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previouslyRequestedGrantee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRequestedGrantee\",\"type\":\"address\"}],\"name\":\"GranteeReassignmentChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldGrantee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newGrantee\",\"type\":\"address\"}],\"name\":\"GranteeReassignmentConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newGrantee\",\"type\":\"address\"}],\"name\":\"GranteeReassignmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokensWithdrawn\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[],\"name\":\"cancelReassignmentRequest\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"cancelStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newGrantee\",\"type\":\"address\"}],\"name\":\"changeReassignmentRequest\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newGrantee\",\"type\":\"address\"}],\"name\":\"confirmGranteeReassignment\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"grantId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"grantManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"grantee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"recoverStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newGrantee\",\"type\":\"address\"}],\"name\":\"requestGranteeReassignment\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"requestedNewGrantee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"stake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractERC20Burnable\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenGrant\",\"outputs\":[{\"internalType\":\"contractTokenGrant\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"undelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ManagedGrant is an auto generated Go binding around an Ethereum contract.
type ManagedGrant struct {
	ManagedGrantCaller     // Read-only binding to the contract
	ManagedGrantTransactor // Write-only binding to the contract
	ManagedGrantFilterer   // Log filterer for contract events
}

// ManagedGrantCaller is an auto generated read-only Go binding around an Ethereum contract.
type ManagedGrantCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagedGrantTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ManagedGrantTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagedGrantFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ManagedGrantFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagedGrantSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ManagedGrantSession struct {
	Contract     *ManagedGrant     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ManagedGrantCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ManagedGrantCallerSession struct {
	Contract *ManagedGrantCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ManagedGrantTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ManagedGrantTransactorSession struct {
	Contract     *ManagedGrantTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ManagedGrantRaw is an auto generated low-level Go binding around an Ethereum contract.
type ManagedGrantRaw struct {
	Contract *ManagedGrant // Generic contract binding to access the raw methods on
}

// ManagedGrantCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ManagedGrantCallerRaw struct {
	Contract *ManagedGrantCaller // Generic read-only contract binding to access the raw methods on
}

// ManagedGrantTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ManagedGrantTransactorRaw struct {
	Contract *ManagedGrantTransactor // Generic write-only contract binding to access the raw methods on
}

// NewManagedGrant creates a new instance of ManagedGrant, bound to a specific deployed contract.
func NewManagedGrant(address common.Address, backend bind.ContractBackend) (*ManagedGrant, error) {
	contract, err := bindManagedGrant(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ManagedGrant{ManagedGrantCaller: ManagedGrantCaller{contract: contract}, ManagedGrantTransactor: ManagedGrantTransactor{contract: contract}, ManagedGrantFilterer: ManagedGrantFilterer{contract: contract}}, nil
}

// NewManagedGrantCaller creates a new read-only instance of ManagedGrant, bound to a specific deployed contract.
func NewManagedGrantCaller(address common.Address, caller bind.ContractCaller) (*ManagedGrantCaller, error) {
	contract, err := bindManagedGrant(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ManagedGrantCaller{contract: contract}, nil
}

// NewManagedGrantTransactor creates a new write-only instance of ManagedGrant, bound to a specific deployed contract.
func NewManagedGrantTransactor(address common.Address, transactor bind.ContractTransactor) (*ManagedGrantTransactor, error) {
	contract, err := bindManagedGrant(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ManagedGrantTransactor{contract: contract}, nil
}

// NewManagedGrantFilterer creates a new log filterer instance of ManagedGrant, bound to a specific deployed contract.
func NewManagedGrantFilterer(address common.Address, filterer bind.ContractFilterer) (*ManagedGrantFilterer, error) {
	contract, err := bindManagedGrant(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ManagedGrantFilterer{contract: contract}, nil
}

// bindManagedGrant binds a generic wrapper to an already deployed contract.
func bindManagedGrant(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ManagedGrantABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagedGrant *ManagedGrantRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ManagedGrant.Contract.ManagedGrantCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagedGrant *ManagedGrantRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagedGrant.Contract.ManagedGrantTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagedGrant *ManagedGrantRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagedGrant.Contract.ManagedGrantTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagedGrant *ManagedGrantCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ManagedGrant.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagedGrant *ManagedGrantTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagedGrant.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagedGrant *ManagedGrantTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagedGrant.Contract.contract.Transact(opts, method, params...)
}

// GrantId is a free data retrieval call binding the contract method 0xf30c0e08.
//
// Solidity: function grantId() constant returns(uint256)
func (_ManagedGrant *ManagedGrantCaller) GrantId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ManagedGrant.contract.Call(opts, out, "grantId")
	return *ret0, err
}

// GrantId is a free data retrieval call binding the contract method 0xf30c0e08.
//
// Solidity: function grantId() constant returns(uint256)
func (_ManagedGrant *ManagedGrantSession) GrantId() (*big.Int, error) {
	return _ManagedGrant.Contract.GrantId(&_ManagedGrant.CallOpts)
}

// GrantId is a free data retrieval call binding the contract method 0xf30c0e08.
//
// Solidity: function grantId() constant returns(uint256)
func (_ManagedGrant *ManagedGrantCallerSession) GrantId() (*big.Int, error) {
	return _ManagedGrant.Contract.GrantId(&_ManagedGrant.CallOpts)
}

// GrantManager is a free data retrieval call binding the contract method 0x545c156b.
//
// Solidity: function grantManager() constant returns(address)
func (_ManagedGrant *ManagedGrantCaller) GrantManager(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ManagedGrant.contract.Call(opts, out, "grantManager")
	return *ret0, err
}

// GrantManager is a free data retrieval call binding the contract method 0x545c156b.
//
// Solidity: function grantManager() constant returns(address)
func (_ManagedGrant *ManagedGrantSession) GrantManager() (common.Address, error) {
	return _ManagedGrant.Contract.GrantManager(&_ManagedGrant.CallOpts)
}

// GrantManager is a free data retrieval call binding the contract method 0x545c156b.
//
// Solidity: function grantManager() constant returns(address)
func (_ManagedGrant *ManagedGrantCallerSession) GrantManager() (common.Address, error) {
	return _ManagedGrant.Contract.GrantManager(&_ManagedGrant.CallOpts)
}

// Grantee is a free data retrieval call binding the contract method 0xd5f52076.
//
// Solidity: function grantee() constant returns(address)
func (_ManagedGrant *ManagedGrantCaller) Grantee(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ManagedGrant.contract.Call(opts, out, "grantee")
	return *ret0, err
}

// Grantee is a free data retrieval call binding the contract method 0xd5f52076.
//
// Solidity: function grantee() constant returns(address)
func (_ManagedGrant *ManagedGrantSession) Grantee() (common.Address, error) {
	return _ManagedGrant.Contract.Grantee(&_ManagedGrant.CallOpts)
}

// Grantee is a free data retrieval call binding the contract method 0xd5f52076.
//
// Solidity: function grantee() constant returns(address)
func (_ManagedGrant *ManagedGrantCallerSession) Grantee() (common.Address, error) {
	return _ManagedGrant.Contract.Grantee(&_ManagedGrant.CallOpts)
}

// RequestedNewGrantee is a free data retrieval call binding the contract method 0x289f51b1.
//
// Solidity: function requestedNewGrantee() constant returns(address)
func (_ManagedGrant *ManagedGrantCaller) RequestedNewGrantee(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ManagedGrant.contract.Call(opts, out, "requestedNewGrantee")
	return *ret0, err
}

// RequestedNewGrantee is a free data retrieval call binding the contract method 0x289f51b1.
//
// Solidity: function requestedNewGrantee() constant returns(address)
func (_ManagedGrant *ManagedGrantSession) RequestedNewGrantee() (common.Address, error) {
	return _ManagedGrant.Contract.RequestedNewGrantee(&_ManagedGrant.CallOpts)
}

// RequestedNewGrantee is a free data retrieval call binding the contract method 0x289f51b1.
//
// Solidity: function requestedNewGrantee() constant returns(address)
func (_ManagedGrant *ManagedGrantCallerSession) RequestedNewGrantee() (common.Address, error) {
	return _ManagedGrant.Contract.RequestedNewGrantee(&_ManagedGrant.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_ManagedGrant *ManagedGrantCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ManagedGrant.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_ManagedGrant *ManagedGrantSession) Token() (common.Address, error) {
	return _ManagedGrant.Contract.Token(&_ManagedGrant.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_ManagedGrant *ManagedGrantCallerSession) Token() (common.Address, error) {
	return _ManagedGrant.Contract.Token(&_ManagedGrant.CallOpts)
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_ManagedGrant *ManagedGrantCaller) TokenGrant(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ManagedGrant.contract.Call(opts, out, "tokenGrant")
	return *ret0, err
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_ManagedGrant *ManagedGrantSession) TokenGrant() (common.Address, error) {
	return _ManagedGrant.Contract.TokenGrant(&_ManagedGrant.CallOpts)
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_ManagedGrant *ManagedGrantCallerSession) TokenGrant() (common.Address, error) {
	return _ManagedGrant.Contract.TokenGrant(&_ManagedGrant.CallOpts)
}

// CancelReassignmentRequest is a paid mutator transaction binding the contract method 0x0af8d820.
//
// Solidity: function cancelReassignmentRequest() returns()
func (_ManagedGrant *ManagedGrantTransactor) CancelReassignmentRequest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "cancelReassignmentRequest")
}

// CancelReassignmentRequest is a paid mutator transaction binding the contract method 0x0af8d820.
//
// Solidity: function cancelReassignmentRequest() returns()
func (_ManagedGrant *ManagedGrantSession) CancelReassignmentRequest() (*types.Transaction, error) {
	return _ManagedGrant.Contract.CancelReassignmentRequest(&_ManagedGrant.TransactOpts)
}

// CancelReassignmentRequest is a paid mutator transaction binding the contract method 0x0af8d820.
//
// Solidity: function cancelReassignmentRequest() returns()
func (_ManagedGrant *ManagedGrantTransactorSession) CancelReassignmentRequest() (*types.Transaction, error) {
	return _ManagedGrant.Contract.CancelReassignmentRequest(&_ManagedGrant.TransactOpts)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_ManagedGrant *ManagedGrantTransactor) CancelStake(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "cancelStake", _operator)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_ManagedGrant *ManagedGrantSession) CancelStake(_operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.CancelStake(&_ManagedGrant.TransactOpts, _operator)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_ManagedGrant *ManagedGrantTransactorSession) CancelStake(_operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.CancelStake(&_ManagedGrant.TransactOpts, _operator)
}

// ChangeReassignmentRequest is a paid mutator transaction binding the contract method 0x8669c083.
//
// Solidity: function changeReassignmentRequest(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantTransactor) ChangeReassignmentRequest(opts *bind.TransactOpts, _newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "changeReassignmentRequest", _newGrantee)
}

// ChangeReassignmentRequest is a paid mutator transaction binding the contract method 0x8669c083.
//
// Solidity: function changeReassignmentRequest(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantSession) ChangeReassignmentRequest(_newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.ChangeReassignmentRequest(&_ManagedGrant.TransactOpts, _newGrantee)
}

// ChangeReassignmentRequest is a paid mutator transaction binding the contract method 0x8669c083.
//
// Solidity: function changeReassignmentRequest(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantTransactorSession) ChangeReassignmentRequest(_newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.ChangeReassignmentRequest(&_ManagedGrant.TransactOpts, _newGrantee)
}

// ConfirmGranteeReassignment is a paid mutator transaction binding the contract method 0x6254d76c.
//
// Solidity: function confirmGranteeReassignment(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantTransactor) ConfirmGranteeReassignment(opts *bind.TransactOpts, _newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "confirmGranteeReassignment", _newGrantee)
}

// ConfirmGranteeReassignment is a paid mutator transaction binding the contract method 0x6254d76c.
//
// Solidity: function confirmGranteeReassignment(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantSession) ConfirmGranteeReassignment(_newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.ConfirmGranteeReassignment(&_ManagedGrant.TransactOpts, _newGrantee)
}

// ConfirmGranteeReassignment is a paid mutator transaction binding the contract method 0x6254d76c.
//
// Solidity: function confirmGranteeReassignment(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantTransactorSession) ConfirmGranteeReassignment(_newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.ConfirmGranteeReassignment(&_ManagedGrant.TransactOpts, _newGrantee)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_ManagedGrant *ManagedGrantTransactor) RecoverStake(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "recoverStake", _operator)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_ManagedGrant *ManagedGrantSession) RecoverStake(_operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.RecoverStake(&_ManagedGrant.TransactOpts, _operator)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_ManagedGrant *ManagedGrantTransactorSession) RecoverStake(_operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.RecoverStake(&_ManagedGrant.TransactOpts, _operator)
}

// RequestGranteeReassignment is a paid mutator transaction binding the contract method 0x5bd8c8e6.
//
// Solidity: function requestGranteeReassignment(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantTransactor) RequestGranteeReassignment(opts *bind.TransactOpts, _newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "requestGranteeReassignment", _newGrantee)
}

// RequestGranteeReassignment is a paid mutator transaction binding the contract method 0x5bd8c8e6.
//
// Solidity: function requestGranteeReassignment(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantSession) RequestGranteeReassignment(_newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.RequestGranteeReassignment(&_ManagedGrant.TransactOpts, _newGrantee)
}

// RequestGranteeReassignment is a paid mutator transaction binding the contract method 0x5bd8c8e6.
//
// Solidity: function requestGranteeReassignment(address _newGrantee) returns()
func (_ManagedGrant *ManagedGrantTransactorSession) RequestGranteeReassignment(_newGrantee common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.RequestGranteeReassignment(&_ManagedGrant.TransactOpts, _newGrantee)
}

// Stake is a paid mutator transaction binding the contract method 0x3e12170f.
//
// Solidity: function stake(address _stakingContract, uint256 _amount, bytes _extraData) returns()
func (_ManagedGrant *ManagedGrantTransactor) Stake(opts *bind.TransactOpts, _stakingContract common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "stake", _stakingContract, _amount, _extraData)
}

// Stake is a paid mutator transaction binding the contract method 0x3e12170f.
//
// Solidity: function stake(address _stakingContract, uint256 _amount, bytes _extraData) returns()
func (_ManagedGrant *ManagedGrantSession) Stake(_stakingContract common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ManagedGrant.Contract.Stake(&_ManagedGrant.TransactOpts, _stakingContract, _amount, _extraData)
}

// Stake is a paid mutator transaction binding the contract method 0x3e12170f.
//
// Solidity: function stake(address _stakingContract, uint256 _amount, bytes _extraData) returns()
func (_ManagedGrant *ManagedGrantTransactorSession) Stake(_stakingContract common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ManagedGrant.Contract.Stake(&_ManagedGrant.TransactOpts, _stakingContract, _amount, _extraData)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_ManagedGrant *ManagedGrantTransactor) Undelegate(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "undelegate", _operator)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_ManagedGrant *ManagedGrantSession) Undelegate(_operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.Undelegate(&_ManagedGrant.TransactOpts, _operator)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_ManagedGrant *ManagedGrantTransactorSession) Undelegate(_operator common.Address) (*types.Transaction, error) {
	return _ManagedGrant.Contract.Undelegate(&_ManagedGrant.TransactOpts, _operator)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ManagedGrant *ManagedGrantTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagedGrant.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ManagedGrant *ManagedGrantSession) Withdraw() (*types.Transaction, error) {
	return _ManagedGrant.Contract.Withdraw(&_ManagedGrant.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ManagedGrant *ManagedGrantTransactorSession) Withdraw() (*types.Transaction, error) {
	return _ManagedGrant.Contract.Withdraw(&_ManagedGrant.TransactOpts)
}

// ManagedGrantGranteeReassignmentCancelledIterator is returned from FilterGranteeReassignmentCancelled and is used to iterate over the raw logs and unpacked data for GranteeReassignmentCancelled events raised by the ManagedGrant contract.
type ManagedGrantGranteeReassignmentCancelledIterator struct {
	Event *ManagedGrantGranteeReassignmentCancelled // Event containing the contract specifics and raw log

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
func (it *ManagedGrantGranteeReassignmentCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagedGrantGranteeReassignmentCancelled)
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
		it.Event = new(ManagedGrantGranteeReassignmentCancelled)
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
func (it *ManagedGrantGranteeReassignmentCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagedGrantGranteeReassignmentCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagedGrantGranteeReassignmentCancelled represents a GranteeReassignmentCancelled event raised by the ManagedGrant contract.
type ManagedGrantGranteeReassignmentCancelled struct {
	CancelledRequestedGrantee common.Address
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterGranteeReassignmentCancelled is a free log retrieval operation binding the contract event 0x079bdabd0be4cb9ac4abae5f8236a5ffa1e3b2cb8292e13683586aa56e5bdcaa.
//
// Solidity: event GranteeReassignmentCancelled(address cancelledRequestedGrantee)
func (_ManagedGrant *ManagedGrantFilterer) FilterGranteeReassignmentCancelled(opts *bind.FilterOpts) (*ManagedGrantGranteeReassignmentCancelledIterator, error) {

	logs, sub, err := _ManagedGrant.contract.FilterLogs(opts, "GranteeReassignmentCancelled")
	if err != nil {
		return nil, err
	}
	return &ManagedGrantGranteeReassignmentCancelledIterator{contract: _ManagedGrant.contract, event: "GranteeReassignmentCancelled", logs: logs, sub: sub}, nil
}

// WatchGranteeReassignmentCancelled is a free log subscription operation binding the contract event 0x079bdabd0be4cb9ac4abae5f8236a5ffa1e3b2cb8292e13683586aa56e5bdcaa.
//
// Solidity: event GranteeReassignmentCancelled(address cancelledRequestedGrantee)
func (_ManagedGrant *ManagedGrantFilterer) WatchGranteeReassignmentCancelled(opts *bind.WatchOpts, sink chan<- *ManagedGrantGranteeReassignmentCancelled) (event.Subscription, error) {

	logs, sub, err := _ManagedGrant.contract.WatchLogs(opts, "GranteeReassignmentCancelled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagedGrantGranteeReassignmentCancelled)
				if err := _ManagedGrant.contract.UnpackLog(event, "GranteeReassignmentCancelled", log); err != nil {
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

// ParseGranteeReassignmentCancelled is a log parse operation binding the contract event 0x079bdabd0be4cb9ac4abae5f8236a5ffa1e3b2cb8292e13683586aa56e5bdcaa.
//
// Solidity: event GranteeReassignmentCancelled(address cancelledRequestedGrantee)
func (_ManagedGrant *ManagedGrantFilterer) ParseGranteeReassignmentCancelled(log types.Log) (*ManagedGrantGranteeReassignmentCancelled, error) {
	event := new(ManagedGrantGranteeReassignmentCancelled)
	if err := _ManagedGrant.contract.UnpackLog(event, "GranteeReassignmentCancelled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ManagedGrantGranteeReassignmentChangedIterator is returned from FilterGranteeReassignmentChanged and is used to iterate over the raw logs and unpacked data for GranteeReassignmentChanged events raised by the ManagedGrant contract.
type ManagedGrantGranteeReassignmentChangedIterator struct {
	Event *ManagedGrantGranteeReassignmentChanged // Event containing the contract specifics and raw log

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
func (it *ManagedGrantGranteeReassignmentChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagedGrantGranteeReassignmentChanged)
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
		it.Event = new(ManagedGrantGranteeReassignmentChanged)
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
func (it *ManagedGrantGranteeReassignmentChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagedGrantGranteeReassignmentChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagedGrantGranteeReassignmentChanged represents a GranteeReassignmentChanged event raised by the ManagedGrant contract.
type ManagedGrantGranteeReassignmentChanged struct {
	PreviouslyRequestedGrantee common.Address
	NewRequestedGrantee        common.Address
	Raw                        types.Log // Blockchain specific contextual infos
}

// FilterGranteeReassignmentChanged is a free log retrieval operation binding the contract event 0x609d93008f7d94dde9255b321cebd574e91bce227f60c61dc2162004af43892c.
//
// Solidity: event GranteeReassignmentChanged(address previouslyRequestedGrantee, address newRequestedGrantee)
func (_ManagedGrant *ManagedGrantFilterer) FilterGranteeReassignmentChanged(opts *bind.FilterOpts) (*ManagedGrantGranteeReassignmentChangedIterator, error) {

	logs, sub, err := _ManagedGrant.contract.FilterLogs(opts, "GranteeReassignmentChanged")
	if err != nil {
		return nil, err
	}
	return &ManagedGrantGranteeReassignmentChangedIterator{contract: _ManagedGrant.contract, event: "GranteeReassignmentChanged", logs: logs, sub: sub}, nil
}

// WatchGranteeReassignmentChanged is a free log subscription operation binding the contract event 0x609d93008f7d94dde9255b321cebd574e91bce227f60c61dc2162004af43892c.
//
// Solidity: event GranteeReassignmentChanged(address previouslyRequestedGrantee, address newRequestedGrantee)
func (_ManagedGrant *ManagedGrantFilterer) WatchGranteeReassignmentChanged(opts *bind.WatchOpts, sink chan<- *ManagedGrantGranteeReassignmentChanged) (event.Subscription, error) {

	logs, sub, err := _ManagedGrant.contract.WatchLogs(opts, "GranteeReassignmentChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagedGrantGranteeReassignmentChanged)
				if err := _ManagedGrant.contract.UnpackLog(event, "GranteeReassignmentChanged", log); err != nil {
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

// ParseGranteeReassignmentChanged is a log parse operation binding the contract event 0x609d93008f7d94dde9255b321cebd574e91bce227f60c61dc2162004af43892c.
//
// Solidity: event GranteeReassignmentChanged(address previouslyRequestedGrantee, address newRequestedGrantee)
func (_ManagedGrant *ManagedGrantFilterer) ParseGranteeReassignmentChanged(log types.Log) (*ManagedGrantGranteeReassignmentChanged, error) {
	event := new(ManagedGrantGranteeReassignmentChanged)
	if err := _ManagedGrant.contract.UnpackLog(event, "GranteeReassignmentChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ManagedGrantGranteeReassignmentConfirmedIterator is returned from FilterGranteeReassignmentConfirmed and is used to iterate over the raw logs and unpacked data for GranteeReassignmentConfirmed events raised by the ManagedGrant contract.
type ManagedGrantGranteeReassignmentConfirmedIterator struct {
	Event *ManagedGrantGranteeReassignmentConfirmed // Event containing the contract specifics and raw log

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
func (it *ManagedGrantGranteeReassignmentConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagedGrantGranteeReassignmentConfirmed)
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
		it.Event = new(ManagedGrantGranteeReassignmentConfirmed)
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
func (it *ManagedGrantGranteeReassignmentConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagedGrantGranteeReassignmentConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagedGrantGranteeReassignmentConfirmed represents a GranteeReassignmentConfirmed event raised by the ManagedGrant contract.
type ManagedGrantGranteeReassignmentConfirmed struct {
	OldGrantee common.Address
	NewGrantee common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterGranteeReassignmentConfirmed is a free log retrieval operation binding the contract event 0xc859af038814fe1f4f2cd57c83de6406ddb0cd2dd3c2294f5995244f3f012c59.
//
// Solidity: event GranteeReassignmentConfirmed(address oldGrantee, address newGrantee)
func (_ManagedGrant *ManagedGrantFilterer) FilterGranteeReassignmentConfirmed(opts *bind.FilterOpts) (*ManagedGrantGranteeReassignmentConfirmedIterator, error) {

	logs, sub, err := _ManagedGrant.contract.FilterLogs(opts, "GranteeReassignmentConfirmed")
	if err != nil {
		return nil, err
	}
	return &ManagedGrantGranteeReassignmentConfirmedIterator{contract: _ManagedGrant.contract, event: "GranteeReassignmentConfirmed", logs: logs, sub: sub}, nil
}

// WatchGranteeReassignmentConfirmed is a free log subscription operation binding the contract event 0xc859af038814fe1f4f2cd57c83de6406ddb0cd2dd3c2294f5995244f3f012c59.
//
// Solidity: event GranteeReassignmentConfirmed(address oldGrantee, address newGrantee)
func (_ManagedGrant *ManagedGrantFilterer) WatchGranteeReassignmentConfirmed(opts *bind.WatchOpts, sink chan<- *ManagedGrantGranteeReassignmentConfirmed) (event.Subscription, error) {

	logs, sub, err := _ManagedGrant.contract.WatchLogs(opts, "GranteeReassignmentConfirmed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagedGrantGranteeReassignmentConfirmed)
				if err := _ManagedGrant.contract.UnpackLog(event, "GranteeReassignmentConfirmed", log); err != nil {
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

// ParseGranteeReassignmentConfirmed is a log parse operation binding the contract event 0xc859af038814fe1f4f2cd57c83de6406ddb0cd2dd3c2294f5995244f3f012c59.
//
// Solidity: event GranteeReassignmentConfirmed(address oldGrantee, address newGrantee)
func (_ManagedGrant *ManagedGrantFilterer) ParseGranteeReassignmentConfirmed(log types.Log) (*ManagedGrantGranteeReassignmentConfirmed, error) {
	event := new(ManagedGrantGranteeReassignmentConfirmed)
	if err := _ManagedGrant.contract.UnpackLog(event, "GranteeReassignmentConfirmed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ManagedGrantGranteeReassignmentRequestedIterator is returned from FilterGranteeReassignmentRequested and is used to iterate over the raw logs and unpacked data for GranteeReassignmentRequested events raised by the ManagedGrant contract.
type ManagedGrantGranteeReassignmentRequestedIterator struct {
	Event *ManagedGrantGranteeReassignmentRequested // Event containing the contract specifics and raw log

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
func (it *ManagedGrantGranteeReassignmentRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagedGrantGranteeReassignmentRequested)
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
		it.Event = new(ManagedGrantGranteeReassignmentRequested)
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
func (it *ManagedGrantGranteeReassignmentRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagedGrantGranteeReassignmentRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagedGrantGranteeReassignmentRequested represents a GranteeReassignmentRequested event raised by the ManagedGrant contract.
type ManagedGrantGranteeReassignmentRequested struct {
	NewGrantee common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterGranteeReassignmentRequested is a free log retrieval operation binding the contract event 0xafd45caea5ec652705350c2ba80247df0fa5dfadf447277da0f21465a74d8fe0.
//
// Solidity: event GranteeReassignmentRequested(address newGrantee)
func (_ManagedGrant *ManagedGrantFilterer) FilterGranteeReassignmentRequested(opts *bind.FilterOpts) (*ManagedGrantGranteeReassignmentRequestedIterator, error) {

	logs, sub, err := _ManagedGrant.contract.FilterLogs(opts, "GranteeReassignmentRequested")
	if err != nil {
		return nil, err
	}
	return &ManagedGrantGranteeReassignmentRequestedIterator{contract: _ManagedGrant.contract, event: "GranteeReassignmentRequested", logs: logs, sub: sub}, nil
}

// WatchGranteeReassignmentRequested is a free log subscription operation binding the contract event 0xafd45caea5ec652705350c2ba80247df0fa5dfadf447277da0f21465a74d8fe0.
//
// Solidity: event GranteeReassignmentRequested(address newGrantee)
func (_ManagedGrant *ManagedGrantFilterer) WatchGranteeReassignmentRequested(opts *bind.WatchOpts, sink chan<- *ManagedGrantGranteeReassignmentRequested) (event.Subscription, error) {

	logs, sub, err := _ManagedGrant.contract.WatchLogs(opts, "GranteeReassignmentRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagedGrantGranteeReassignmentRequested)
				if err := _ManagedGrant.contract.UnpackLog(event, "GranteeReassignmentRequested", log); err != nil {
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

// ParseGranteeReassignmentRequested is a log parse operation binding the contract event 0xafd45caea5ec652705350c2ba80247df0fa5dfadf447277da0f21465a74d8fe0.
//
// Solidity: event GranteeReassignmentRequested(address newGrantee)
func (_ManagedGrant *ManagedGrantFilterer) ParseGranteeReassignmentRequested(log types.Log) (*ManagedGrantGranteeReassignmentRequested, error) {
	event := new(ManagedGrantGranteeReassignmentRequested)
	if err := _ManagedGrant.contract.UnpackLog(event, "GranteeReassignmentRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ManagedGrantTokensWithdrawnIterator is returned from FilterTokensWithdrawn and is used to iterate over the raw logs and unpacked data for TokensWithdrawn events raised by the ManagedGrant contract.
type ManagedGrantTokensWithdrawnIterator struct {
	Event *ManagedGrantTokensWithdrawn // Event containing the contract specifics and raw log

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
func (it *ManagedGrantTokensWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagedGrantTokensWithdrawn)
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
		it.Event = new(ManagedGrantTokensWithdrawn)
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
func (it *ManagedGrantTokensWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagedGrantTokensWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagedGrantTokensWithdrawn represents a TokensWithdrawn event raised by the ManagedGrant contract.
type ManagedGrantTokensWithdrawn struct {
	Destination common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokensWithdrawn is a free log retrieval operation binding the contract event 0x6352c5382c4a4578e712449ca65e83cdb392d045dfcf1cad9615189db2da244b.
//
// Solidity: event TokensWithdrawn(address destination, uint256 amount)
func (_ManagedGrant *ManagedGrantFilterer) FilterTokensWithdrawn(opts *bind.FilterOpts) (*ManagedGrantTokensWithdrawnIterator, error) {

	logs, sub, err := _ManagedGrant.contract.FilterLogs(opts, "TokensWithdrawn")
	if err != nil {
		return nil, err
	}
	return &ManagedGrantTokensWithdrawnIterator{contract: _ManagedGrant.contract, event: "TokensWithdrawn", logs: logs, sub: sub}, nil
}

// WatchTokensWithdrawn is a free log subscription operation binding the contract event 0x6352c5382c4a4578e712449ca65e83cdb392d045dfcf1cad9615189db2da244b.
//
// Solidity: event TokensWithdrawn(address destination, uint256 amount)
func (_ManagedGrant *ManagedGrantFilterer) WatchTokensWithdrawn(opts *bind.WatchOpts, sink chan<- *ManagedGrantTokensWithdrawn) (event.Subscription, error) {

	logs, sub, err := _ManagedGrant.contract.WatchLogs(opts, "TokensWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagedGrantTokensWithdrawn)
				if err := _ManagedGrant.contract.UnpackLog(event, "TokensWithdrawn", log); err != nil {
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
// Solidity: event TokensWithdrawn(address destination, uint256 amount)
func (_ManagedGrant *ManagedGrantFilterer) ParseTokensWithdrawn(log types.Log) (*ManagedGrantTokensWithdrawn, error) {
	event := new(ManagedGrantTokensWithdrawn)
	if err := _ManagedGrant.contract.UnpackLog(event, "TokensWithdrawn", log); err != nil {
		return nil, err
	}
	return event, nil
}
