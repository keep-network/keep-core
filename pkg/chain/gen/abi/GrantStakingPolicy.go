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

// GrantStakingPolicyABI is the input ABI used to generate the binding from.
const GrantStakingPolicyABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_now\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"grantedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cliff\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawn\",\"type\":\"uint256\"}],\"name\":\"getStakeableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// GrantStakingPolicy is an auto generated Go binding around an Ethereum contract.
type GrantStakingPolicy struct {
	GrantStakingPolicyCaller     // Read-only binding to the contract
	GrantStakingPolicyTransactor // Write-only binding to the contract
	GrantStakingPolicyFilterer   // Log filterer for contract events
}

// GrantStakingPolicyCaller is an auto generated read-only Go binding around an Ethereum contract.
type GrantStakingPolicyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GrantStakingPolicyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GrantStakingPolicyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GrantStakingPolicyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GrantStakingPolicyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GrantStakingPolicySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GrantStakingPolicySession struct {
	Contract     *GrantStakingPolicy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// GrantStakingPolicyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GrantStakingPolicyCallerSession struct {
	Contract *GrantStakingPolicyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// GrantStakingPolicyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GrantStakingPolicyTransactorSession struct {
	Contract     *GrantStakingPolicyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// GrantStakingPolicyRaw is an auto generated low-level Go binding around an Ethereum contract.
type GrantStakingPolicyRaw struct {
	Contract *GrantStakingPolicy // Generic contract binding to access the raw methods on
}

// GrantStakingPolicyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GrantStakingPolicyCallerRaw struct {
	Contract *GrantStakingPolicyCaller // Generic read-only contract binding to access the raw methods on
}

// GrantStakingPolicyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GrantStakingPolicyTransactorRaw struct {
	Contract *GrantStakingPolicyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGrantStakingPolicy creates a new instance of GrantStakingPolicy, bound to a specific deployed contract.
func NewGrantStakingPolicy(address common.Address, backend bind.ContractBackend) (*GrantStakingPolicy, error) {
	contract, err := bindGrantStakingPolicy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GrantStakingPolicy{GrantStakingPolicyCaller: GrantStakingPolicyCaller{contract: contract}, GrantStakingPolicyTransactor: GrantStakingPolicyTransactor{contract: contract}, GrantStakingPolicyFilterer: GrantStakingPolicyFilterer{contract: contract}}, nil
}

// NewGrantStakingPolicyCaller creates a new read-only instance of GrantStakingPolicy, bound to a specific deployed contract.
func NewGrantStakingPolicyCaller(address common.Address, caller bind.ContractCaller) (*GrantStakingPolicyCaller, error) {
	contract, err := bindGrantStakingPolicy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GrantStakingPolicyCaller{contract: contract}, nil
}

// NewGrantStakingPolicyTransactor creates a new write-only instance of GrantStakingPolicy, bound to a specific deployed contract.
func NewGrantStakingPolicyTransactor(address common.Address, transactor bind.ContractTransactor) (*GrantStakingPolicyTransactor, error) {
	contract, err := bindGrantStakingPolicy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GrantStakingPolicyTransactor{contract: contract}, nil
}

// NewGrantStakingPolicyFilterer creates a new log filterer instance of GrantStakingPolicy, bound to a specific deployed contract.
func NewGrantStakingPolicyFilterer(address common.Address, filterer bind.ContractFilterer) (*GrantStakingPolicyFilterer, error) {
	contract, err := bindGrantStakingPolicy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GrantStakingPolicyFilterer{contract: contract}, nil
}

// bindGrantStakingPolicy binds a generic wrapper to an already deployed contract.
func bindGrantStakingPolicy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GrantStakingPolicyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GrantStakingPolicy *GrantStakingPolicyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GrantStakingPolicy.Contract.GrantStakingPolicyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GrantStakingPolicy *GrantStakingPolicyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GrantStakingPolicy.Contract.GrantStakingPolicyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GrantStakingPolicy *GrantStakingPolicyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GrantStakingPolicy.Contract.GrantStakingPolicyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GrantStakingPolicy *GrantStakingPolicyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GrantStakingPolicy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GrantStakingPolicy *GrantStakingPolicyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GrantStakingPolicy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GrantStakingPolicy *GrantStakingPolicyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GrantStakingPolicy.Contract.contract.Transact(opts, method, params...)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_GrantStakingPolicy *GrantStakingPolicyCaller) GetStakeableAmount(opts *bind.CallOpts, _now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _GrantStakingPolicy.contract.Call(opts, out, "getStakeableAmount", _now, grantedAmount, duration, start, cliff, withdrawn)
	return *ret0, err
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_GrantStakingPolicy *GrantStakingPolicySession) GetStakeableAmount(_now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	return _GrantStakingPolicy.Contract.GetStakeableAmount(&_GrantStakingPolicy.CallOpts, _now, grantedAmount, duration, start, cliff, withdrawn)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_GrantStakingPolicy *GrantStakingPolicyCallerSession) GetStakeableAmount(_now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	return _GrantStakingPolicy.Contract.GetStakeableAmount(&_GrantStakingPolicy.CallOpts, _now, grantedAmount, duration, start, cliff, withdrawn)
}
