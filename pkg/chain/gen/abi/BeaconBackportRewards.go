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

// BeaconBackportRewardsABI is the input ABI used to generate the binding from.
const BeaconBackportRewardsABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakingContract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keep\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountTransferred\",\"type\":\"uint256\"}],\"name\":\"UpgradeFinalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRewardsContract\",\"type\":\"address\"}],\"name\":\"UpgradeInitiated\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"allocateRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dispensedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"eligibleButTerminated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"eligibleForReward\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"endOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeRewardsUpgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"firstIntervalStart\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"funded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"getAllocatedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getIntervalCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"getIntervalWeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newRewardsContract\",\"type\":\"address\"}],\"name\":\"initiateRewardsUpgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervalKeepsProcessed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"intervalOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervalWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"isAllocated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"isFinished\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"keepsInInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"markAsFunded\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumKeepsPerInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newRewardsContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keepIdentifier\",\"type\":\"bytes32\"}],\"name\":\"receiveReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"receiveReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"keepIdentifiers\",\"type\":\"bytes32[]\"}],\"name\":\"receiveRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keepIdentifier\",\"type\":\"bytes32\"}],\"name\":\"reportTermination\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"keepIdentifiers\",\"type\":\"bytes32[]\"}],\"name\":\"reportTerminations\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"rewardClaimed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"startOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"termLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractKeepToken\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"unallocatedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradeFinalizedTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradeInitiatedTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// BeaconBackportRewards is an auto generated Go binding around an Ethereum contract.
type BeaconBackportRewards struct {
	BeaconBackportRewardsCaller     // Read-only binding to the contract
	BeaconBackportRewardsTransactor // Write-only binding to the contract
	BeaconBackportRewardsFilterer   // Log filterer for contract events
}

// BeaconBackportRewardsCaller is an auto generated read-only Go binding around an Ethereum contract.
type BeaconBackportRewardsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconBackportRewardsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BeaconBackportRewardsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconBackportRewardsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BeaconBackportRewardsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconBackportRewardsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BeaconBackportRewardsSession struct {
	Contract     *BeaconBackportRewards // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BeaconBackportRewardsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BeaconBackportRewardsCallerSession struct {
	Contract *BeaconBackportRewardsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// BeaconBackportRewardsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BeaconBackportRewardsTransactorSession struct {
	Contract     *BeaconBackportRewardsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// BeaconBackportRewardsRaw is an auto generated low-level Go binding around an Ethereum contract.
type BeaconBackportRewardsRaw struct {
	Contract *BeaconBackportRewards // Generic contract binding to access the raw methods on
}

// BeaconBackportRewardsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BeaconBackportRewardsCallerRaw struct {
	Contract *BeaconBackportRewardsCaller // Generic read-only contract binding to access the raw methods on
}

// BeaconBackportRewardsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BeaconBackportRewardsTransactorRaw struct {
	Contract *BeaconBackportRewardsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBeaconBackportRewards creates a new instance of BeaconBackportRewards, bound to a specific deployed contract.
func NewBeaconBackportRewards(address common.Address, backend bind.ContractBackend) (*BeaconBackportRewards, error) {
	contract, err := bindBeaconBackportRewards(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconBackportRewards{BeaconBackportRewardsCaller: BeaconBackportRewardsCaller{contract: contract}, BeaconBackportRewardsTransactor: BeaconBackportRewardsTransactor{contract: contract}, BeaconBackportRewardsFilterer: BeaconBackportRewardsFilterer{contract: contract}}, nil
}

// NewBeaconBackportRewardsCaller creates a new read-only instance of BeaconBackportRewards, bound to a specific deployed contract.
func NewBeaconBackportRewardsCaller(address common.Address, caller bind.ContractCaller) (*BeaconBackportRewardsCaller, error) {
	contract, err := bindBeaconBackportRewards(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconBackportRewardsCaller{contract: contract}, nil
}

// NewBeaconBackportRewardsTransactor creates a new write-only instance of BeaconBackportRewards, bound to a specific deployed contract.
func NewBeaconBackportRewardsTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconBackportRewardsTransactor, error) {
	contract, err := bindBeaconBackportRewards(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconBackportRewardsTransactor{contract: contract}, nil
}

// NewBeaconBackportRewardsFilterer creates a new log filterer instance of BeaconBackportRewards, bound to a specific deployed contract.
func NewBeaconBackportRewardsFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconBackportRewardsFilterer, error) {
	contract, err := bindBeaconBackportRewards(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconBackportRewardsFilterer{contract: contract}, nil
}

// bindBeaconBackportRewards binds a generic wrapper to an already deployed contract.
func bindBeaconBackportRewards(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BeaconBackportRewardsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconBackportRewards *BeaconBackportRewardsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BeaconBackportRewards.Contract.BeaconBackportRewardsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconBackportRewards *BeaconBackportRewardsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.BeaconBackportRewardsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconBackportRewards *BeaconBackportRewardsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.BeaconBackportRewardsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconBackportRewards *BeaconBackportRewardsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BeaconBackportRewards.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.contract.Transact(opts, method, params...)
}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) DispensedRewards(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "dispensedRewards")
	return *ret0, err
}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) DispensedRewards() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.DispensedRewards(&_BeaconBackportRewards.CallOpts)
}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) DispensedRewards() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.DispensedRewards(&_BeaconBackportRewards.CallOpts)
}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) EligibleButTerminated(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "eligibleButTerminated", _keep)
	return *ret0, err
}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) EligibleButTerminated(_keep [32]byte) (bool, error) {
	return _BeaconBackportRewards.Contract.EligibleButTerminated(&_BeaconBackportRewards.CallOpts, _keep)
}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) EligibleButTerminated(_keep [32]byte) (bool, error) {
	return _BeaconBackportRewards.Contract.EligibleButTerminated(&_BeaconBackportRewards.CallOpts, _keep)
}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) EligibleForReward(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "eligibleForReward", _keep)
	return *ret0, err
}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) EligibleForReward(_keep [32]byte) (bool, error) {
	return _BeaconBackportRewards.Contract.EligibleForReward(&_BeaconBackportRewards.CallOpts, _keep)
}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) EligibleForReward(_keep [32]byte) (bool, error) {
	return _BeaconBackportRewards.Contract.EligibleForReward(&_BeaconBackportRewards.CallOpts, _keep)
}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) EndOf(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "endOf", interval)
	return *ret0, err
}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) EndOf(interval *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.EndOf(&_BeaconBackportRewards.CallOpts, interval)
}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) EndOf(interval *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.EndOf(&_BeaconBackportRewards.CallOpts, interval)
}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) FirstIntervalStart(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "firstIntervalStart")
	return *ret0, err
}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) FirstIntervalStart() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.FirstIntervalStart(&_BeaconBackportRewards.CallOpts)
}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) FirstIntervalStart() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.FirstIntervalStart(&_BeaconBackportRewards.CallOpts)
}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) Funded(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "funded")
	return *ret0, err
}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) Funded() (bool, error) {
	return _BeaconBackportRewards.Contract.Funded(&_BeaconBackportRewards.CallOpts)
}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) Funded() (bool, error) {
	return _BeaconBackportRewards.Contract.Funded(&_BeaconBackportRewards.CallOpts)
}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) GetAllocatedRewards(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "getAllocatedRewards", interval)
	return *ret0, err
}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) GetAllocatedRewards(interval *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.GetAllocatedRewards(&_BeaconBackportRewards.CallOpts, interval)
}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) GetAllocatedRewards(interval *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.GetAllocatedRewards(&_BeaconBackportRewards.CallOpts, interval)
}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) GetIntervalCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "getIntervalCount")
	return *ret0, err
}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) GetIntervalCount() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.GetIntervalCount(&_BeaconBackportRewards.CallOpts)
}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) GetIntervalCount() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.GetIntervalCount(&_BeaconBackportRewards.CallOpts)
}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) GetIntervalWeight(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "getIntervalWeight", interval)
	return *ret0, err
}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) GetIntervalWeight(interval *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.GetIntervalWeight(&_BeaconBackportRewards.CallOpts, interval)
}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) GetIntervalWeight(interval *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.GetIntervalWeight(&_BeaconBackportRewards.CallOpts, interval)
}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) IntervalKeepsProcessed(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "intervalKeepsProcessed", arg0)
	return *ret0, err
}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) IntervalKeepsProcessed(arg0 *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.IntervalKeepsProcessed(&_BeaconBackportRewards.CallOpts, arg0)
}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) IntervalKeepsProcessed(arg0 *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.IntervalKeepsProcessed(&_BeaconBackportRewards.CallOpts, arg0)
}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) IntervalOf(opts *bind.CallOpts, timestamp *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "intervalOf", timestamp)
	return *ret0, err
}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) IntervalOf(timestamp *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.IntervalOf(&_BeaconBackportRewards.CallOpts, timestamp)
}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) IntervalOf(timestamp *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.IntervalOf(&_BeaconBackportRewards.CallOpts, timestamp)
}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) IntervalWeights(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "intervalWeights", arg0)
	return *ret0, err
}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) IntervalWeights(arg0 *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.IntervalWeights(&_BeaconBackportRewards.CallOpts, arg0)
}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) IntervalWeights(arg0 *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.IntervalWeights(&_BeaconBackportRewards.CallOpts, arg0)
}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) IsAllocated(opts *bind.CallOpts, interval *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "isAllocated", interval)
	return *ret0, err
}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) IsAllocated(interval *big.Int) (bool, error) {
	return _BeaconBackportRewards.Contract.IsAllocated(&_BeaconBackportRewards.CallOpts, interval)
}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) IsAllocated(interval *big.Int) (bool, error) {
	return _BeaconBackportRewards.Contract.IsAllocated(&_BeaconBackportRewards.CallOpts, interval)
}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) IsFinished(opts *bind.CallOpts, interval *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "isFinished", interval)
	return *ret0, err
}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) IsFinished(interval *big.Int) (bool, error) {
	return _BeaconBackportRewards.Contract.IsFinished(&_BeaconBackportRewards.CallOpts, interval)
}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) IsFinished(interval *big.Int) (bool, error) {
	return _BeaconBackportRewards.Contract.IsFinished(&_BeaconBackportRewards.CallOpts, interval)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) IsOwner() (bool, error) {
	return _BeaconBackportRewards.Contract.IsOwner(&_BeaconBackportRewards.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) IsOwner() (bool, error) {
	return _BeaconBackportRewards.Contract.IsOwner(&_BeaconBackportRewards.CallOpts)
}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) MinimumKeepsPerInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "minimumKeepsPerInterval")
	return *ret0, err
}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) MinimumKeepsPerInterval() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.MinimumKeepsPerInterval(&_BeaconBackportRewards.CallOpts)
}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) MinimumKeepsPerInterval() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.MinimumKeepsPerInterval(&_BeaconBackportRewards.CallOpts)
}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) NewRewardsContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "newRewardsContract")
	return *ret0, err
}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) NewRewardsContract() (common.Address, error) {
	return _BeaconBackportRewards.Contract.NewRewardsContract(&_BeaconBackportRewards.CallOpts)
}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) NewRewardsContract() (common.Address, error) {
	return _BeaconBackportRewards.Contract.NewRewardsContract(&_BeaconBackportRewards.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) Owner() (common.Address, error) {
	return _BeaconBackportRewards.Contract.Owner(&_BeaconBackportRewards.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) Owner() (common.Address, error) {
	return _BeaconBackportRewards.Contract.Owner(&_BeaconBackportRewards.CallOpts)
}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) RewardClaimed(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "rewardClaimed", _keep)
	return *ret0, err
}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) RewardClaimed(_keep [32]byte) (bool, error) {
	return _BeaconBackportRewards.Contract.RewardClaimed(&_BeaconBackportRewards.CallOpts, _keep)
}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) constant returns(bool)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) RewardClaimed(_keep [32]byte) (bool, error) {
	return _BeaconBackportRewards.Contract.RewardClaimed(&_BeaconBackportRewards.CallOpts, _keep)
}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) StartOf(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "startOf", interval)
	return *ret0, err
}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) StartOf(interval *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.StartOf(&_BeaconBackportRewards.CallOpts, interval)
}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) StartOf(interval *big.Int) (*big.Int, error) {
	return _BeaconBackportRewards.Contract.StartOf(&_BeaconBackportRewards.CallOpts, interval)
}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) TermLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "termLength")
	return *ret0, err
}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) TermLength() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.TermLength(&_BeaconBackportRewards.CallOpts)
}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) TermLength() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.TermLength(&_BeaconBackportRewards.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) Token() (common.Address, error) {
	return _BeaconBackportRewards.Contract.Token(&_BeaconBackportRewards.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) Token() (common.Address, error) {
	return _BeaconBackportRewards.Contract.Token(&_BeaconBackportRewards.CallOpts)
}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) TotalRewards(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "totalRewards")
	return *ret0, err
}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) TotalRewards() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.TotalRewards(&_BeaconBackportRewards.CallOpts)
}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) TotalRewards() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.TotalRewards(&_BeaconBackportRewards.CallOpts)
}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) UnallocatedRewards(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "unallocatedRewards")
	return *ret0, err
}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) UnallocatedRewards() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.UnallocatedRewards(&_BeaconBackportRewards.CallOpts)
}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) UnallocatedRewards() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.UnallocatedRewards(&_BeaconBackportRewards.CallOpts)
}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) UpgradeFinalizedTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "upgradeFinalizedTimestamp")
	return *ret0, err
}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) UpgradeFinalizedTimestamp() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.UpgradeFinalizedTimestamp(&_BeaconBackportRewards.CallOpts)
}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) UpgradeFinalizedTimestamp() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.UpgradeFinalizedTimestamp(&_BeaconBackportRewards.CallOpts)
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCaller) UpgradeInitiatedTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BeaconBackportRewards.contract.Call(opts, out, "upgradeInitiatedTimestamp")
	return *ret0, err
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) UpgradeInitiatedTimestamp() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.UpgradeInitiatedTimestamp(&_BeaconBackportRewards.CallOpts)
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsCallerSession) UpgradeInitiatedTimestamp() (*big.Int, error) {
	return _BeaconBackportRewards.Contract.UpgradeInitiatedTimestamp(&_BeaconBackportRewards.CallOpts)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) AllocateRewards(opts *bind.TransactOpts, interval *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "allocateRewards", interval)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) AllocateRewards(interval *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.AllocateRewards(&_BeaconBackportRewards.TransactOpts, interval)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) AllocateRewards(interval *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.AllocateRewards(&_BeaconBackportRewards.TransactOpts, interval)
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) FinalizeRewardsUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "finalizeRewardsUpgrade")
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) FinalizeRewardsUpgrade() (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.FinalizeRewardsUpgrade(&_BeaconBackportRewards.TransactOpts)
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) FinalizeRewardsUpgrade() (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.FinalizeRewardsUpgrade(&_BeaconBackportRewards.TransactOpts)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) InitiateRewardsUpgrade(opts *bind.TransactOpts, _newRewardsContract common.Address) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "initiateRewardsUpgrade", _newRewardsContract)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) InitiateRewardsUpgrade(_newRewardsContract common.Address) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.InitiateRewardsUpgrade(&_BeaconBackportRewards.TransactOpts, _newRewardsContract)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) InitiateRewardsUpgrade(_newRewardsContract common.Address) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.InitiateRewardsUpgrade(&_BeaconBackportRewards.TransactOpts, _newRewardsContract)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) KeepsInInterval(opts *bind.TransactOpts, interval *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "keepsInInterval", interval)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsSession) KeepsInInterval(interval *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.KeepsInInterval(&_BeaconBackportRewards.TransactOpts, interval)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) KeepsInInterval(interval *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.KeepsInInterval(&_BeaconBackportRewards.TransactOpts, interval)
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) MarkAsFunded(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "markAsFunded")
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) MarkAsFunded() (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.MarkAsFunded(&_BeaconBackportRewards.TransactOpts)
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) MarkAsFunded() (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.MarkAsFunded(&_BeaconBackportRewards.TransactOpts)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) ReceiveApproval(opts *bind.TransactOpts, _from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "receiveApproval", _from, _value, _token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReceiveApproval(&_BeaconBackportRewards.TransactOpts, _from, _value, _token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReceiveApproval(&_BeaconBackportRewards.TransactOpts, _from, _value, _token, arg3)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) ReceiveReward(opts *bind.TransactOpts, keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "receiveReward", keepIdentifier)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) ReceiveReward(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReceiveReward(&_BeaconBackportRewards.TransactOpts, keepIdentifier)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) ReceiveReward(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReceiveReward(&_BeaconBackportRewards.TransactOpts, keepIdentifier)
}

// ReceiveReward0 is a paid mutator transaction binding the contract method 0xae10d5ea.
//
// Solidity: function receiveReward(uint256 groupIndex) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) ReceiveReward0(opts *bind.TransactOpts, groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "receiveReward0", groupIndex)
}

// ReceiveReward0 is a paid mutator transaction binding the contract method 0xae10d5ea.
//
// Solidity: function receiveReward(uint256 groupIndex) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) ReceiveReward0(groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReceiveReward0(&_BeaconBackportRewards.TransactOpts, groupIndex)
}

// ReceiveReward0 is a paid mutator transaction binding the contract method 0xae10d5ea.
//
// Solidity: function receiveReward(uint256 groupIndex) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) ReceiveReward0(groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReceiveReward0(&_BeaconBackportRewards.TransactOpts, groupIndex)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) ReceiveRewards(opts *bind.TransactOpts, keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "receiveRewards", keepIdentifiers)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) ReceiveRewards(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReceiveRewards(&_BeaconBackportRewards.TransactOpts, keepIdentifiers)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) ReceiveRewards(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReceiveRewards(&_BeaconBackportRewards.TransactOpts, keepIdentifiers)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) RenounceOwnership() (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.RenounceOwnership(&_BeaconBackportRewards.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.RenounceOwnership(&_BeaconBackportRewards.TransactOpts)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) ReportTermination(opts *bind.TransactOpts, keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "reportTermination", keepIdentifier)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) ReportTermination(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReportTermination(&_BeaconBackportRewards.TransactOpts, keepIdentifier)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) ReportTermination(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReportTermination(&_BeaconBackportRewards.TransactOpts, keepIdentifier)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) ReportTerminations(opts *bind.TransactOpts, keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "reportTerminations", keepIdentifiers)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) ReportTerminations(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReportTerminations(&_BeaconBackportRewards.TransactOpts, keepIdentifiers)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) ReportTerminations(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.ReportTerminations(&_BeaconBackportRewards.TransactOpts, keepIdentifiers)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BeaconBackportRewards.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.TransferOwnership(&_BeaconBackportRewards.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BeaconBackportRewards *BeaconBackportRewardsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BeaconBackportRewards.Contract.TransferOwnership(&_BeaconBackportRewards.TransactOpts, newOwner)
}

// BeaconBackportRewardsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BeaconBackportRewards contract.
type BeaconBackportRewardsOwnershipTransferredIterator struct {
	Event *BeaconBackportRewardsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BeaconBackportRewardsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconBackportRewardsOwnershipTransferred)
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
		it.Event = new(BeaconBackportRewardsOwnershipTransferred)
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
func (it *BeaconBackportRewardsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconBackportRewardsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconBackportRewardsOwnershipTransferred represents a OwnershipTransferred event raised by the BeaconBackportRewards contract.
type BeaconBackportRewardsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BeaconBackportRewardsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BeaconBackportRewards.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BeaconBackportRewardsOwnershipTransferredIterator{contract: _BeaconBackportRewards.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BeaconBackportRewardsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BeaconBackportRewards.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconBackportRewardsOwnershipTransferred)
				if err := _BeaconBackportRewards.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) ParseOwnershipTransferred(log types.Log) (*BeaconBackportRewardsOwnershipTransferred, error) {
	event := new(BeaconBackportRewardsOwnershipTransferred)
	if err := _BeaconBackportRewards.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BeaconBackportRewardsRewardReceivedIterator is returned from FilterRewardReceived and is used to iterate over the raw logs and unpacked data for RewardReceived events raised by the BeaconBackportRewards contract.
type BeaconBackportRewardsRewardReceivedIterator struct {
	Event *BeaconBackportRewardsRewardReceived // Event containing the contract specifics and raw log

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
func (it *BeaconBackportRewardsRewardReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconBackportRewardsRewardReceived)
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
		it.Event = new(BeaconBackportRewardsRewardReceived)
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
func (it *BeaconBackportRewardsRewardReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconBackportRewardsRewardReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconBackportRewardsRewardReceived represents a RewardReceived event raised by the BeaconBackportRewards contract.
type BeaconBackportRewardsRewardReceived struct {
	Keep   [32]byte
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardReceived is a free log retrieval operation binding the contract event 0xf8fcbb083cc485f0dc726d4235dbec7c9c5c03af58254f353627503577216170.
//
// Solidity: event RewardReceived(bytes32 keep, uint256 amount)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) FilterRewardReceived(opts *bind.FilterOpts) (*BeaconBackportRewardsRewardReceivedIterator, error) {

	logs, sub, err := _BeaconBackportRewards.contract.FilterLogs(opts, "RewardReceived")
	if err != nil {
		return nil, err
	}
	return &BeaconBackportRewardsRewardReceivedIterator{contract: _BeaconBackportRewards.contract, event: "RewardReceived", logs: logs, sub: sub}, nil
}

// WatchRewardReceived is a free log subscription operation binding the contract event 0xf8fcbb083cc485f0dc726d4235dbec7c9c5c03af58254f353627503577216170.
//
// Solidity: event RewardReceived(bytes32 keep, uint256 amount)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) WatchRewardReceived(opts *bind.WatchOpts, sink chan<- *BeaconBackportRewardsRewardReceived) (event.Subscription, error) {

	logs, sub, err := _BeaconBackportRewards.contract.WatchLogs(opts, "RewardReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconBackportRewardsRewardReceived)
				if err := _BeaconBackportRewards.contract.UnpackLog(event, "RewardReceived", log); err != nil {
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

// ParseRewardReceived is a log parse operation binding the contract event 0xf8fcbb083cc485f0dc726d4235dbec7c9c5c03af58254f353627503577216170.
//
// Solidity: event RewardReceived(bytes32 keep, uint256 amount)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) ParseRewardReceived(log types.Log) (*BeaconBackportRewardsRewardReceived, error) {
	event := new(BeaconBackportRewardsRewardReceived)
	if err := _BeaconBackportRewards.contract.UnpackLog(event, "RewardReceived", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BeaconBackportRewardsUpgradeFinalizedIterator is returned from FilterUpgradeFinalized and is used to iterate over the raw logs and unpacked data for UpgradeFinalized events raised by the BeaconBackportRewards contract.
type BeaconBackportRewardsUpgradeFinalizedIterator struct {
	Event *BeaconBackportRewardsUpgradeFinalized // Event containing the contract specifics and raw log

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
func (it *BeaconBackportRewardsUpgradeFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconBackportRewardsUpgradeFinalized)
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
		it.Event = new(BeaconBackportRewardsUpgradeFinalized)
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
func (it *BeaconBackportRewardsUpgradeFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconBackportRewardsUpgradeFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconBackportRewardsUpgradeFinalized represents a UpgradeFinalized event raised by the BeaconBackportRewards contract.
type BeaconBackportRewardsUpgradeFinalized struct {
	AmountTransferred *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgradeFinalized is a free log retrieval operation binding the contract event 0x89f5b857dbf5a32a4fa55162349dc1146850c9b76ff05bedf3a47b3f14095c0a.
//
// Solidity: event UpgradeFinalized(uint256 amountTransferred)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) FilterUpgradeFinalized(opts *bind.FilterOpts) (*BeaconBackportRewardsUpgradeFinalizedIterator, error) {

	logs, sub, err := _BeaconBackportRewards.contract.FilterLogs(opts, "UpgradeFinalized")
	if err != nil {
		return nil, err
	}
	return &BeaconBackportRewardsUpgradeFinalizedIterator{contract: _BeaconBackportRewards.contract, event: "UpgradeFinalized", logs: logs, sub: sub}, nil
}

// WatchUpgradeFinalized is a free log subscription operation binding the contract event 0x89f5b857dbf5a32a4fa55162349dc1146850c9b76ff05bedf3a47b3f14095c0a.
//
// Solidity: event UpgradeFinalized(uint256 amountTransferred)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) WatchUpgradeFinalized(opts *bind.WatchOpts, sink chan<- *BeaconBackportRewardsUpgradeFinalized) (event.Subscription, error) {

	logs, sub, err := _BeaconBackportRewards.contract.WatchLogs(opts, "UpgradeFinalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconBackportRewardsUpgradeFinalized)
				if err := _BeaconBackportRewards.contract.UnpackLog(event, "UpgradeFinalized", log); err != nil {
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

// ParseUpgradeFinalized is a log parse operation binding the contract event 0x89f5b857dbf5a32a4fa55162349dc1146850c9b76ff05bedf3a47b3f14095c0a.
//
// Solidity: event UpgradeFinalized(uint256 amountTransferred)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) ParseUpgradeFinalized(log types.Log) (*BeaconBackportRewardsUpgradeFinalized, error) {
	event := new(BeaconBackportRewardsUpgradeFinalized)
	if err := _BeaconBackportRewards.contract.UnpackLog(event, "UpgradeFinalized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// BeaconBackportRewardsUpgradeInitiatedIterator is returned from FilterUpgradeInitiated and is used to iterate over the raw logs and unpacked data for UpgradeInitiated events raised by the BeaconBackportRewards contract.
type BeaconBackportRewardsUpgradeInitiatedIterator struct {
	Event *BeaconBackportRewardsUpgradeInitiated // Event containing the contract specifics and raw log

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
func (it *BeaconBackportRewardsUpgradeInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconBackportRewardsUpgradeInitiated)
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
		it.Event = new(BeaconBackportRewardsUpgradeInitiated)
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
func (it *BeaconBackportRewardsUpgradeInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconBackportRewardsUpgradeInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconBackportRewardsUpgradeInitiated represents a UpgradeInitiated event raised by the BeaconBackportRewards contract.
type BeaconBackportRewardsUpgradeInitiated struct {
	NewRewardsContract common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpgradeInitiated is a free log retrieval operation binding the contract event 0x4508cb0974b8dda93dec7c130aa186db2705fcee58bc00d1334726761835aca1.
//
// Solidity: event UpgradeInitiated(address newRewardsContract)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) FilterUpgradeInitiated(opts *bind.FilterOpts) (*BeaconBackportRewardsUpgradeInitiatedIterator, error) {

	logs, sub, err := _BeaconBackportRewards.contract.FilterLogs(opts, "UpgradeInitiated")
	if err != nil {
		return nil, err
	}
	return &BeaconBackportRewardsUpgradeInitiatedIterator{contract: _BeaconBackportRewards.contract, event: "UpgradeInitiated", logs: logs, sub: sub}, nil
}

// WatchUpgradeInitiated is a free log subscription operation binding the contract event 0x4508cb0974b8dda93dec7c130aa186db2705fcee58bc00d1334726761835aca1.
//
// Solidity: event UpgradeInitiated(address newRewardsContract)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) WatchUpgradeInitiated(opts *bind.WatchOpts, sink chan<- *BeaconBackportRewardsUpgradeInitiated) (event.Subscription, error) {

	logs, sub, err := _BeaconBackportRewards.contract.WatchLogs(opts, "UpgradeInitiated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconBackportRewardsUpgradeInitiated)
				if err := _BeaconBackportRewards.contract.UnpackLog(event, "UpgradeInitiated", log); err != nil {
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

// ParseUpgradeInitiated is a log parse operation binding the contract event 0x4508cb0974b8dda93dec7c130aa186db2705fcee58bc00d1334726761835aca1.
//
// Solidity: event UpgradeInitiated(address newRewardsContract)
func (_BeaconBackportRewards *BeaconBackportRewardsFilterer) ParseUpgradeInitiated(log types.Log) (*BeaconBackportRewardsUpgradeInitiated, error) {
	event := new(BeaconBackportRewardsUpgradeInitiated)
	if err := _BeaconBackportRewards.contract.UnpackLog(event, "UpgradeInitiated", log); err != nil {
		return nil, err
	}
	return event, nil
}
