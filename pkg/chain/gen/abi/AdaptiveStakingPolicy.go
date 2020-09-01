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

// AdaptiveStakingPolicyABI is the input ABI used to generate the binding from.
const AdaptiveStakingPolicyABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumMultiplier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_stakeaheadTime\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_useCliff\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_now\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"grantedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cliff\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawn\",\"type\":\"uint256\"}],\"name\":\"getStakeableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AdaptiveStakingPolicy is an auto generated Go binding around an Ethereum contract.
type AdaptiveStakingPolicy struct {
	AdaptiveStakingPolicyCaller     // Read-only binding to the contract
	AdaptiveStakingPolicyTransactor // Write-only binding to the contract
	AdaptiveStakingPolicyFilterer   // Log filterer for contract events
}

// AdaptiveStakingPolicyCaller is an auto generated read-only Go binding around an Ethereum contract.
type AdaptiveStakingPolicyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveStakingPolicyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AdaptiveStakingPolicyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveStakingPolicyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AdaptiveStakingPolicyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveStakingPolicySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AdaptiveStakingPolicySession struct {
	Contract     *AdaptiveStakingPolicy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AdaptiveStakingPolicyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AdaptiveStakingPolicyCallerSession struct {
	Contract *AdaptiveStakingPolicyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// AdaptiveStakingPolicyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AdaptiveStakingPolicyTransactorSession struct {
	Contract     *AdaptiveStakingPolicyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// AdaptiveStakingPolicyRaw is an auto generated low-level Go binding around an Ethereum contract.
type AdaptiveStakingPolicyRaw struct {
	Contract *AdaptiveStakingPolicy // Generic contract binding to access the raw methods on
}

// AdaptiveStakingPolicyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AdaptiveStakingPolicyCallerRaw struct {
	Contract *AdaptiveStakingPolicyCaller // Generic read-only contract binding to access the raw methods on
}

// AdaptiveStakingPolicyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AdaptiveStakingPolicyTransactorRaw struct {
	Contract *AdaptiveStakingPolicyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAdaptiveStakingPolicy creates a new instance of AdaptiveStakingPolicy, bound to a specific deployed contract.
func NewAdaptiveStakingPolicy(address common.Address, backend bind.ContractBackend) (*AdaptiveStakingPolicy, error) {
	contract, err := bindAdaptiveStakingPolicy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AdaptiveStakingPolicy{AdaptiveStakingPolicyCaller: AdaptiveStakingPolicyCaller{contract: contract}, AdaptiveStakingPolicyTransactor: AdaptiveStakingPolicyTransactor{contract: contract}, AdaptiveStakingPolicyFilterer: AdaptiveStakingPolicyFilterer{contract: contract}}, nil
}

// NewAdaptiveStakingPolicyCaller creates a new read-only instance of AdaptiveStakingPolicy, bound to a specific deployed contract.
func NewAdaptiveStakingPolicyCaller(address common.Address, caller bind.ContractCaller) (*AdaptiveStakingPolicyCaller, error) {
	contract, err := bindAdaptiveStakingPolicy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AdaptiveStakingPolicyCaller{contract: contract}, nil
}

// NewAdaptiveStakingPolicyTransactor creates a new write-only instance of AdaptiveStakingPolicy, bound to a specific deployed contract.
func NewAdaptiveStakingPolicyTransactor(address common.Address, transactor bind.ContractTransactor) (*AdaptiveStakingPolicyTransactor, error) {
	contract, err := bindAdaptiveStakingPolicy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AdaptiveStakingPolicyTransactor{contract: contract}, nil
}

// NewAdaptiveStakingPolicyFilterer creates a new log filterer instance of AdaptiveStakingPolicy, bound to a specific deployed contract.
func NewAdaptiveStakingPolicyFilterer(address common.Address, filterer bind.ContractFilterer) (*AdaptiveStakingPolicyFilterer, error) {
	contract, err := bindAdaptiveStakingPolicy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AdaptiveStakingPolicyFilterer{contract: contract}, nil
}

// bindAdaptiveStakingPolicy binds a generic wrapper to an already deployed contract.
func bindAdaptiveStakingPolicy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AdaptiveStakingPolicyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AdaptiveStakingPolicy.Contract.AdaptiveStakingPolicyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdaptiveStakingPolicy.Contract.AdaptiveStakingPolicyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AdaptiveStakingPolicy.Contract.AdaptiveStakingPolicyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AdaptiveStakingPolicy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdaptiveStakingPolicy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AdaptiveStakingPolicy.Contract.contract.Transact(opts, method, params...)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicyCaller) GetStakeableAmount(opts *bind.CallOpts, _now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AdaptiveStakingPolicy.contract.Call(opts, out, "getStakeableAmount", _now, grantedAmount, duration, start, cliff, withdrawn)
	return *ret0, err
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicySession) GetStakeableAmount(_now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	return _AdaptiveStakingPolicy.Contract.GetStakeableAmount(&_AdaptiveStakingPolicy.CallOpts, _now, grantedAmount, duration, start, cliff, withdrawn)
}

// GetStakeableAmount is a free data retrieval call binding the contract method 0xdab40935.
//
// Solidity: function getStakeableAmount(uint256 _now, uint256 grantedAmount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn) constant returns(uint256)
func (_AdaptiveStakingPolicy *AdaptiveStakingPolicyCallerSession) GetStakeableAmount(_now *big.Int, grantedAmount *big.Int, duration *big.Int, start *big.Int, cliff *big.Int, withdrawn *big.Int) (*big.Int, error) {
	return _AdaptiveStakingPolicy.Contract.GetStakeableAmount(&_AdaptiveStakingPolicy.CallOpts, _now, grantedAmount, duration, start, cliff, withdrawn)
}
