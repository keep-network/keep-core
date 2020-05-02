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

// PermissiveStakingPolicyABI is the input ABI used to generate the binding from.
const PermissiveStakingPolicyABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_now\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"grantedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cliff\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawn\",\"type\":\"uint256\"}],\"name\":\"getStakeableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// PermissiveStakingPolicy is an auto generated Go binding around an Ethereum contract.
type PermissiveStakingPolicy struct {
	PermissiveStakingPolicyCaller     // Read-only binding to the contract
	PermissiveStakingPolicyTransactor // Write-only binding to the contract
	PermissiveStakingPolicyFilterer   // Log filterer for contract events
}

// PermissiveStakingPolicyCaller is an auto generated read-only Go binding around an Ethereum contract.
type PermissiveStakingPolicyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PermissiveStakingPolicyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PermissiveStakingPolicyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PermissiveStakingPolicyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PermissiveStakingPolicyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PermissiveStakingPolicySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PermissiveStakingPolicySession struct {
	Contract     *PermissiveStakingPolicy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// PermissiveStakingPolicyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PermissiveStakingPolicyCallerSession struct {
	Contract *PermissiveStakingPolicyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// PermissiveStakingPolicyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PermissiveStakingPolicyTransactorSession struct {
	Contract     *PermissiveStakingPolicyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// PermissiveStakingPolicyRaw is an auto generated low-level Go binding around an Ethereum contract.
type PermissiveStakingPolicyRaw struct {
	Contract *PermissiveStakingPolicy // Generic contract binding to access the raw methods on
}

// PermissiveStakingPolicyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PermissiveStakingPolicyCallerRaw struct {
	Contract *PermissiveStakingPolicyCaller // Generic read-only contract binding to access the raw methods on
}

// PermissiveStakingPolicyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PermissiveStakingPolicyTransactorRaw struct {
	Contract *PermissiveStakingPolicyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPermissiveStakingPolicy creates a new instance of PermissiveStakingPolicy, bound to a specific deployed contract.
func NewPermissiveStakingPolicy(address common.Address, backend bind.ContractBackend) (*PermissiveStakingPolicy, error) {
	contract, err := bindPermissiveStakingPolicy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PermissiveStakingPolicy{PermissiveStakingPolicyCaller: PermissiveStakingPolicyCaller{contract: contract}, PermissiveStakingPolicyTransactor: PermissiveStakingPolicyTransactor{contract: contract}, PermissiveStakingPolicyFilterer: PermissiveStakingPolicyFilterer{contract: contract}}, nil
}

// NewPermissiveStakingPolicyCaller creates a new read-only instance of PermissiveStakingPolicy, bound to a specific deployed contract.
func NewPermissiveStakingPolicyCaller(address common.Address, caller bind.ContractCaller) (*PermissiveStakingPolicyCaller, error) {
	contract, err := bindPermissiveStakingPolicy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PermissiveStakingPolicyCaller{contract: contract}, nil
}

// NewPermissiveStakingPolicyTransactor creates a new write-only instance of PermissiveStakingPolicy, bound to a specific deployed contract.
func NewPermissiveStakingPolicyTransactor(address common.Address, transactor bind.ContractTransactor) (*PermissiveStakingPolicyTransactor, error) {
	contract, err := bindPermissiveStakingPolicy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PermissiveStakingPolicyTransactor{contract: contract}, nil
}

// NewPermissiveStakingPolicyFilterer creates a new log filterer instance of PermissiveStakingPolicy, bound to a specific deployed contract.
func NewPermissiveStakingPolicyFilterer(address common.Address, filterer bind.ContractFilterer) (*PermissiveStakingPolicyFilterer, error) {
	contract, err := bindPermissiveStakingPolicy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PermissiveStakingPolicyFilterer{contract: contract}, nil
}

// bindPermissiveStakingPolicy binds a generic wrapper to an already deployed contract.
func bindPermissiveStakingPolicy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PermissiveStakingPolicyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PermissiveStakingPolicy *PermissiveStakingPolicyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PermissiveStakingPolicy.Contract.PermissiveStakingPolicyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PermissiveStakingPolicy *PermissiveStakingPolicyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissiveStakingPolicy.Contract.PermissiveStakingPolicyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PermissiveStakingPolicy *PermissiveStakingPolicyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PermissiveStakingPolicy.Contract.PermissiveStakingPolicyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PermissiveStakingPolicy *PermissiveStakingPolicyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PermissiveStakingPolicy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PermissiveStakingPolicy *PermissiveStakingPolicyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissiveStakingPolicy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PermissiveStakingPolicy *PermissiveStakingPolicyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PermissiveStakingPolicy.Contract.contract.Transact(opts, method, params...)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_PermissiveStakingPolicy *PermissiveStakingPolicyCaller) GetStakeableAmount(opts *bind.CallOpts, _now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PermissiveStakingPolicy.contract.Call(opts, out, "getStakeableAmount", _now, grantedAmount, duration, start, cliff, withdrawn)
	return *ret0, err
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_PermissiveStakingPolicy *PermissiveStakingPolicySession) GetStakeableAmount(_now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	return _PermissiveStakingPolicy.Contract.GetStakeableAmount(&_PermissiveStakingPolicy.CallOpts, _now, grantedAmount, duration, start, cliff, withdrawn)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_PermissiveStakingPolicy *PermissiveStakingPolicyCallerSession) GetStakeableAmount(_now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	return _PermissiveStakingPolicy.Contract.GetStakeableAmount(&_PermissiveStakingPolicy.CallOpts, _now, grantedAmount, duration, start, cliff, withdrawn)
}
