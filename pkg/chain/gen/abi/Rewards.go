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

// RewardsABI is the input ABI used to generate the binding from.
const RewardsABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_firstIntervalStart\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_intervalWeights\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_termLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minimumKeepsPerInterval\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keep\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountTransferred\",\"type\":\"uint256\"}],\"name\":\"UpgradeFinalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRewardsContract\",\"type\":\"address\"}],\"name\":\"UpgradeInitiated\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"allocateRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dispensedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"eligibleButTerminated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"eligibleForReward\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"endOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeRewardsUpgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"firstIntervalStart\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"funded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"getAllocatedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getIntervalCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"getIntervalWeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newRewardsContract\",\"type\":\"address\"}],\"name\":\"initiateRewardsUpgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervalKeepsProcessed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"intervalOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervalWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"isAllocated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"isFinished\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"keepsInInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"markAsFunded\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumKeepsPerInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newRewardsContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keepIdentifier\",\"type\":\"bytes32\"}],\"name\":\"receiveReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"keepIdentifiers\",\"type\":\"bytes32[]\"}],\"name\":\"receiveRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keepIdentifier\",\"type\":\"bytes32\"}],\"name\":\"reportTermination\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"keepIdentifiers\",\"type\":\"bytes32[]\"}],\"name\":\"reportTerminations\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"rewardClaimed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"startOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"termLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractKeepToken\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"unallocatedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradeFinalizedTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradeInitiatedTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Rewards is an auto generated Go binding around an Ethereum contract.
type Rewards struct {
	RewardsCaller     // Read-only binding to the contract
	RewardsTransactor // Write-only binding to the contract
	RewardsFilterer   // Log filterer for contract events
}

// RewardsCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewardsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewardsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewardsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewardsSession struct {
	Contract     *Rewards          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewardsCallerSession struct {
	Contract *RewardsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// RewardsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewardsTransactorSession struct {
	Contract     *RewardsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// RewardsRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewardsRaw struct {
	Contract *Rewards // Generic contract binding to access the raw methods on
}

// RewardsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewardsCallerRaw struct {
	Contract *RewardsCaller // Generic read-only contract binding to access the raw methods on
}

// RewardsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewardsTransactorRaw struct {
	Contract *RewardsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRewards creates a new instance of Rewards, bound to a specific deployed contract.
func NewRewards(address common.Address, backend bind.ContractBackend) (*Rewards, error) {
	contract, err := bindRewards(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rewards{RewardsCaller: RewardsCaller{contract: contract}, RewardsTransactor: RewardsTransactor{contract: contract}, RewardsFilterer: RewardsFilterer{contract: contract}}, nil
}

// NewRewardsCaller creates a new read-only instance of Rewards, bound to a specific deployed contract.
func NewRewardsCaller(address common.Address, caller bind.ContractCaller) (*RewardsCaller, error) {
	contract, err := bindRewards(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardsCaller{contract: contract}, nil
}

// NewRewardsTransactor creates a new write-only instance of Rewards, bound to a specific deployed contract.
func NewRewardsTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardsTransactor, error) {
	contract, err := bindRewards(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardsTransactor{contract: contract}, nil
}

// NewRewardsFilterer creates a new log filterer instance of Rewards, bound to a specific deployed contract.
func NewRewardsFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardsFilterer, error) {
	contract, err := bindRewards(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardsFilterer{contract: contract}, nil
}

// bindRewards binds a generic wrapper to an already deployed contract.
func bindRewards(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RewardsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rewards *RewardsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Rewards.Contract.RewardsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rewards *RewardsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewards.Contract.RewardsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rewards *RewardsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rewards.Contract.RewardsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rewards *RewardsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Rewards.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rewards *RewardsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewards.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rewards *RewardsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rewards.Contract.contract.Transact(opts, method, params...)
}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() constant returns(uint256)
func (_Rewards *RewardsCaller) DispensedRewards(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "dispensedRewards")
	return *ret0, err
}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() constant returns(uint256)
func (_Rewards *RewardsSession) DispensedRewards() (*big.Int, error) {
	return _Rewards.Contract.DispensedRewards(&_Rewards.CallOpts)
}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() constant returns(uint256)
func (_Rewards *RewardsCallerSession) DispensedRewards() (*big.Int, error) {
	return _Rewards.Contract.DispensedRewards(&_Rewards.CallOpts)
}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsCaller) EligibleButTerminated(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "eligibleButTerminated", _keep)
	return *ret0, err
}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsSession) EligibleButTerminated(_keep [32]byte) (bool, error) {
	return _Rewards.Contract.EligibleButTerminated(&_Rewards.CallOpts, _keep)
}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsCallerSession) EligibleButTerminated(_keep [32]byte) (bool, error) {
	return _Rewards.Contract.EligibleButTerminated(&_Rewards.CallOpts, _keep)
}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsCaller) EligibleForReward(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "eligibleForReward", _keep)
	return *ret0, err
}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsSession) EligibleForReward(_keep [32]byte) (bool, error) {
	return _Rewards.Contract.EligibleForReward(&_Rewards.CallOpts, _keep)
}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsCallerSession) EligibleForReward(_keep [32]byte) (bool, error) {
	return _Rewards.Contract.EligibleForReward(&_Rewards.CallOpts, _keep)
}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsCaller) EndOf(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "endOf", interval)
	return *ret0, err
}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsSession) EndOf(interval *big.Int) (*big.Int, error) {
	return _Rewards.Contract.EndOf(&_Rewards.CallOpts, interval)
}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsCallerSession) EndOf(interval *big.Int) (*big.Int, error) {
	return _Rewards.Contract.EndOf(&_Rewards.CallOpts, interval)
}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() constant returns(uint256)
func (_Rewards *RewardsCaller) FirstIntervalStart(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "firstIntervalStart")
	return *ret0, err
}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() constant returns(uint256)
func (_Rewards *RewardsSession) FirstIntervalStart() (*big.Int, error) {
	return _Rewards.Contract.FirstIntervalStart(&_Rewards.CallOpts)
}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() constant returns(uint256)
func (_Rewards *RewardsCallerSession) FirstIntervalStart() (*big.Int, error) {
	return _Rewards.Contract.FirstIntervalStart(&_Rewards.CallOpts)
}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() constant returns(bool)
func (_Rewards *RewardsCaller) Funded(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "funded")
	return *ret0, err
}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() constant returns(bool)
func (_Rewards *RewardsSession) Funded() (bool, error) {
	return _Rewards.Contract.Funded(&_Rewards.CallOpts)
}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() constant returns(bool)
func (_Rewards *RewardsCallerSession) Funded() (bool, error) {
	return _Rewards.Contract.Funded(&_Rewards.CallOpts)
}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsCaller) GetAllocatedRewards(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "getAllocatedRewards", interval)
	return *ret0, err
}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsSession) GetAllocatedRewards(interval *big.Int) (*big.Int, error) {
	return _Rewards.Contract.GetAllocatedRewards(&_Rewards.CallOpts, interval)
}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsCallerSession) GetAllocatedRewards(interval *big.Int) (*big.Int, error) {
	return _Rewards.Contract.GetAllocatedRewards(&_Rewards.CallOpts, interval)
}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() constant returns(uint256)
func (_Rewards *RewardsCaller) GetIntervalCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "getIntervalCount")
	return *ret0, err
}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() constant returns(uint256)
func (_Rewards *RewardsSession) GetIntervalCount() (*big.Int, error) {
	return _Rewards.Contract.GetIntervalCount(&_Rewards.CallOpts)
}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() constant returns(uint256)
func (_Rewards *RewardsCallerSession) GetIntervalCount() (*big.Int, error) {
	return _Rewards.Contract.GetIntervalCount(&_Rewards.CallOpts)
}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsCaller) GetIntervalWeight(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "getIntervalWeight", interval)
	return *ret0, err
}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsSession) GetIntervalWeight(interval *big.Int) (*big.Int, error) {
	return _Rewards.Contract.GetIntervalWeight(&_Rewards.CallOpts, interval)
}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsCallerSession) GetIntervalWeight(interval *big.Int) (*big.Int, error) {
	return _Rewards.Contract.GetIntervalWeight(&_Rewards.CallOpts, interval)
}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) constant returns(uint256)
func (_Rewards *RewardsCaller) IntervalKeepsProcessed(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "intervalKeepsProcessed", arg0)
	return *ret0, err
}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) constant returns(uint256)
func (_Rewards *RewardsSession) IntervalKeepsProcessed(arg0 *big.Int) (*big.Int, error) {
	return _Rewards.Contract.IntervalKeepsProcessed(&_Rewards.CallOpts, arg0)
}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) constant returns(uint256)
func (_Rewards *RewardsCallerSession) IntervalKeepsProcessed(arg0 *big.Int) (*big.Int, error) {
	return _Rewards.Contract.IntervalKeepsProcessed(&_Rewards.CallOpts, arg0)
}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) constant returns(uint256)
func (_Rewards *RewardsCaller) IntervalOf(opts *bind.CallOpts, timestamp *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "intervalOf", timestamp)
	return *ret0, err
}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) constant returns(uint256)
func (_Rewards *RewardsSession) IntervalOf(timestamp *big.Int) (*big.Int, error) {
	return _Rewards.Contract.IntervalOf(&_Rewards.CallOpts, timestamp)
}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) constant returns(uint256)
func (_Rewards *RewardsCallerSession) IntervalOf(timestamp *big.Int) (*big.Int, error) {
	return _Rewards.Contract.IntervalOf(&_Rewards.CallOpts, timestamp)
}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) constant returns(uint256)
func (_Rewards *RewardsCaller) IntervalWeights(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "intervalWeights", arg0)
	return *ret0, err
}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) constant returns(uint256)
func (_Rewards *RewardsSession) IntervalWeights(arg0 *big.Int) (*big.Int, error) {
	return _Rewards.Contract.IntervalWeights(&_Rewards.CallOpts, arg0)
}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) constant returns(uint256)
func (_Rewards *RewardsCallerSession) IntervalWeights(arg0 *big.Int) (*big.Int, error) {
	return _Rewards.Contract.IntervalWeights(&_Rewards.CallOpts, arg0)
}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) constant returns(bool)
func (_Rewards *RewardsCaller) IsAllocated(opts *bind.CallOpts, interval *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "isAllocated", interval)
	return *ret0, err
}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) constant returns(bool)
func (_Rewards *RewardsSession) IsAllocated(interval *big.Int) (bool, error) {
	return _Rewards.Contract.IsAllocated(&_Rewards.CallOpts, interval)
}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) constant returns(bool)
func (_Rewards *RewardsCallerSession) IsAllocated(interval *big.Int) (bool, error) {
	return _Rewards.Contract.IsAllocated(&_Rewards.CallOpts, interval)
}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) constant returns(bool)
func (_Rewards *RewardsCaller) IsFinished(opts *bind.CallOpts, interval *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "isFinished", interval)
	return *ret0, err
}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) constant returns(bool)
func (_Rewards *RewardsSession) IsFinished(interval *big.Int) (bool, error) {
	return _Rewards.Contract.IsFinished(&_Rewards.CallOpts, interval)
}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) constant returns(bool)
func (_Rewards *RewardsCallerSession) IsFinished(interval *big.Int) (bool, error) {
	return _Rewards.Contract.IsFinished(&_Rewards.CallOpts, interval)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Rewards *RewardsCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Rewards *RewardsSession) IsOwner() (bool, error) {
	return _Rewards.Contract.IsOwner(&_Rewards.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Rewards *RewardsCallerSession) IsOwner() (bool, error) {
	return _Rewards.Contract.IsOwner(&_Rewards.CallOpts)
}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() constant returns(uint256)
func (_Rewards *RewardsCaller) MinimumKeepsPerInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "minimumKeepsPerInterval")
	return *ret0, err
}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() constant returns(uint256)
func (_Rewards *RewardsSession) MinimumKeepsPerInterval() (*big.Int, error) {
	return _Rewards.Contract.MinimumKeepsPerInterval(&_Rewards.CallOpts)
}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() constant returns(uint256)
func (_Rewards *RewardsCallerSession) MinimumKeepsPerInterval() (*big.Int, error) {
	return _Rewards.Contract.MinimumKeepsPerInterval(&_Rewards.CallOpts)
}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() constant returns(address)
func (_Rewards *RewardsCaller) NewRewardsContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "newRewardsContract")
	return *ret0, err
}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() constant returns(address)
func (_Rewards *RewardsSession) NewRewardsContract() (common.Address, error) {
	return _Rewards.Contract.NewRewardsContract(&_Rewards.CallOpts)
}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() constant returns(address)
func (_Rewards *RewardsCallerSession) NewRewardsContract() (common.Address, error) {
	return _Rewards.Contract.NewRewardsContract(&_Rewards.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Rewards *RewardsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Rewards *RewardsSession) Owner() (common.Address, error) {
	return _Rewards.Contract.Owner(&_Rewards.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Rewards *RewardsCallerSession) Owner() (common.Address, error) {
	return _Rewards.Contract.Owner(&_Rewards.CallOpts)
}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsCaller) RewardClaimed(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "rewardClaimed", _keep)
	return *ret0, err
}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsSession) RewardClaimed(_keep [32]byte) (bool, error) {
	return _Rewards.Contract.RewardClaimed(&_Rewards.CallOpts, _keep)
}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) constant returns(bool)
func (_Rewards *RewardsCallerSession) RewardClaimed(_keep [32]byte) (bool, error) {
	return _Rewards.Contract.RewardClaimed(&_Rewards.CallOpts, _keep)
}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsCaller) StartOf(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "startOf", interval)
	return *ret0, err
}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsSession) StartOf(interval *big.Int) (*big.Int, error) {
	return _Rewards.Contract.StartOf(&_Rewards.CallOpts, interval)
}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) constant returns(uint256)
func (_Rewards *RewardsCallerSession) StartOf(interval *big.Int) (*big.Int, error) {
	return _Rewards.Contract.StartOf(&_Rewards.CallOpts, interval)
}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() constant returns(uint256)
func (_Rewards *RewardsCaller) TermLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "termLength")
	return *ret0, err
}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() constant returns(uint256)
func (_Rewards *RewardsSession) TermLength() (*big.Int, error) {
	return _Rewards.Contract.TermLength(&_Rewards.CallOpts)
}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() constant returns(uint256)
func (_Rewards *RewardsCallerSession) TermLength() (*big.Int, error) {
	return _Rewards.Contract.TermLength(&_Rewards.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Rewards *RewardsCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Rewards *RewardsSession) Token() (common.Address, error) {
	return _Rewards.Contract.Token(&_Rewards.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Rewards *RewardsCallerSession) Token() (common.Address, error) {
	return _Rewards.Contract.Token(&_Rewards.CallOpts)
}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() constant returns(uint256)
func (_Rewards *RewardsCaller) TotalRewards(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "totalRewards")
	return *ret0, err
}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() constant returns(uint256)
func (_Rewards *RewardsSession) TotalRewards() (*big.Int, error) {
	return _Rewards.Contract.TotalRewards(&_Rewards.CallOpts)
}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() constant returns(uint256)
func (_Rewards *RewardsCallerSession) TotalRewards() (*big.Int, error) {
	return _Rewards.Contract.TotalRewards(&_Rewards.CallOpts)
}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() constant returns(uint256)
func (_Rewards *RewardsCaller) UnallocatedRewards(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "unallocatedRewards")
	return *ret0, err
}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() constant returns(uint256)
func (_Rewards *RewardsSession) UnallocatedRewards() (*big.Int, error) {
	return _Rewards.Contract.UnallocatedRewards(&_Rewards.CallOpts)
}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() constant returns(uint256)
func (_Rewards *RewardsCallerSession) UnallocatedRewards() (*big.Int, error) {
	return _Rewards.Contract.UnallocatedRewards(&_Rewards.CallOpts)
}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() constant returns(uint256)
func (_Rewards *RewardsCaller) UpgradeFinalizedTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "upgradeFinalizedTimestamp")
	return *ret0, err
}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() constant returns(uint256)
func (_Rewards *RewardsSession) UpgradeFinalizedTimestamp() (*big.Int, error) {
	return _Rewards.Contract.UpgradeFinalizedTimestamp(&_Rewards.CallOpts)
}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() constant returns(uint256)
func (_Rewards *RewardsCallerSession) UpgradeFinalizedTimestamp() (*big.Int, error) {
	return _Rewards.Contract.UpgradeFinalizedTimestamp(&_Rewards.CallOpts)
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256)
func (_Rewards *RewardsCaller) UpgradeInitiatedTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rewards.contract.Call(opts, out, "upgradeInitiatedTimestamp")
	return *ret0, err
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256)
func (_Rewards *RewardsSession) UpgradeInitiatedTimestamp() (*big.Int, error) {
	return _Rewards.Contract.UpgradeInitiatedTimestamp(&_Rewards.CallOpts)
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() constant returns(uint256)
func (_Rewards *RewardsCallerSession) UpgradeInitiatedTimestamp() (*big.Int, error) {
	return _Rewards.Contract.UpgradeInitiatedTimestamp(&_Rewards.CallOpts)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_Rewards *RewardsTransactor) AllocateRewards(opts *bind.TransactOpts, interval *big.Int) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "allocateRewards", interval)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_Rewards *RewardsSession) AllocateRewards(interval *big.Int) (*types.Transaction, error) {
	return _Rewards.Contract.AllocateRewards(&_Rewards.TransactOpts, interval)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_Rewards *RewardsTransactorSession) AllocateRewards(interval *big.Int) (*types.Transaction, error) {
	return _Rewards.Contract.AllocateRewards(&_Rewards.TransactOpts, interval)
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_Rewards *RewardsTransactor) FinalizeRewardsUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "finalizeRewardsUpgrade")
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_Rewards *RewardsSession) FinalizeRewardsUpgrade() (*types.Transaction, error) {
	return _Rewards.Contract.FinalizeRewardsUpgrade(&_Rewards.TransactOpts)
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_Rewards *RewardsTransactorSession) FinalizeRewardsUpgrade() (*types.Transaction, error) {
	return _Rewards.Contract.FinalizeRewardsUpgrade(&_Rewards.TransactOpts)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_Rewards *RewardsTransactor) InitiateRewardsUpgrade(opts *bind.TransactOpts, _newRewardsContract common.Address) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "initiateRewardsUpgrade", _newRewardsContract)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_Rewards *RewardsSession) InitiateRewardsUpgrade(_newRewardsContract common.Address) (*types.Transaction, error) {
	return _Rewards.Contract.InitiateRewardsUpgrade(&_Rewards.TransactOpts, _newRewardsContract)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_Rewards *RewardsTransactorSession) InitiateRewardsUpgrade(_newRewardsContract common.Address) (*types.Transaction, error) {
	return _Rewards.Contract.InitiateRewardsUpgrade(&_Rewards.TransactOpts, _newRewardsContract)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_Rewards *RewardsTransactor) KeepsInInterval(opts *bind.TransactOpts, interval *big.Int) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "keepsInInterval", interval)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_Rewards *RewardsSession) KeepsInInterval(interval *big.Int) (*types.Transaction, error) {
	return _Rewards.Contract.KeepsInInterval(&_Rewards.TransactOpts, interval)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_Rewards *RewardsTransactorSession) KeepsInInterval(interval *big.Int) (*types.Transaction, error) {
	return _Rewards.Contract.KeepsInInterval(&_Rewards.TransactOpts, interval)
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_Rewards *RewardsTransactor) MarkAsFunded(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "markAsFunded")
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_Rewards *RewardsSession) MarkAsFunded() (*types.Transaction, error) {
	return _Rewards.Contract.MarkAsFunded(&_Rewards.TransactOpts)
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_Rewards *RewardsTransactorSession) MarkAsFunded() (*types.Transaction, error) {
	return _Rewards.Contract.MarkAsFunded(&_Rewards.TransactOpts)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_Rewards *RewardsTransactor) ReceiveApproval(opts *bind.TransactOpts, _from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "receiveApproval", _from, _value, _token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_Rewards *RewardsSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReceiveApproval(&_Rewards.TransactOpts, _from, _value, _token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_Rewards *RewardsTransactorSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReceiveApproval(&_Rewards.TransactOpts, _from, _value, _token, arg3)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_Rewards *RewardsTransactor) ReceiveReward(opts *bind.TransactOpts, keepIdentifier [32]byte) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "receiveReward", keepIdentifier)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_Rewards *RewardsSession) ReceiveReward(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReceiveReward(&_Rewards.TransactOpts, keepIdentifier)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_Rewards *RewardsTransactorSession) ReceiveReward(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReceiveReward(&_Rewards.TransactOpts, keepIdentifier)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_Rewards *RewardsTransactor) ReceiveRewards(opts *bind.TransactOpts, keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "receiveRewards", keepIdentifiers)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_Rewards *RewardsSession) ReceiveRewards(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReceiveRewards(&_Rewards.TransactOpts, keepIdentifiers)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_Rewards *RewardsTransactorSession) ReceiveRewards(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReceiveRewards(&_Rewards.TransactOpts, keepIdentifiers)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewards *RewardsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewards *RewardsSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rewards.Contract.RenounceOwnership(&_Rewards.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewards *RewardsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rewards.Contract.RenounceOwnership(&_Rewards.TransactOpts)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_Rewards *RewardsTransactor) ReportTermination(opts *bind.TransactOpts, keepIdentifier [32]byte) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "reportTermination", keepIdentifier)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_Rewards *RewardsSession) ReportTermination(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReportTermination(&_Rewards.TransactOpts, keepIdentifier)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_Rewards *RewardsTransactorSession) ReportTermination(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReportTermination(&_Rewards.TransactOpts, keepIdentifier)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_Rewards *RewardsTransactor) ReportTerminations(opts *bind.TransactOpts, keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "reportTerminations", keepIdentifiers)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_Rewards *RewardsSession) ReportTerminations(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReportTerminations(&_Rewards.TransactOpts, keepIdentifiers)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_Rewards *RewardsTransactorSession) ReportTerminations(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _Rewards.Contract.ReportTerminations(&_Rewards.TransactOpts, keepIdentifiers)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewards *RewardsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Rewards.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewards *RewardsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rewards.Contract.TransferOwnership(&_Rewards.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewards *RewardsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rewards.Contract.TransferOwnership(&_Rewards.TransactOpts, newOwner)
}

// RewardsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Rewards contract.
type RewardsOwnershipTransferredIterator struct {
	Event *RewardsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RewardsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardsOwnershipTransferred)
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
		it.Event = new(RewardsOwnershipTransferred)
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
func (it *RewardsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardsOwnershipTransferred represents a OwnershipTransferred event raised by the Rewards contract.
type RewardsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rewards *RewardsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RewardsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewards.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RewardsOwnershipTransferredIterator{contract: _Rewards.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rewards *RewardsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RewardsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewards.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardsOwnershipTransferred)
				if err := _Rewards.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Rewards *RewardsFilterer) ParseOwnershipTransferred(log types.Log) (*RewardsOwnershipTransferred, error) {
	event := new(RewardsOwnershipTransferred)
	if err := _Rewards.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RewardsRewardReceivedIterator is returned from FilterRewardReceived and is used to iterate over the raw logs and unpacked data for RewardReceived events raised by the Rewards contract.
type RewardsRewardReceivedIterator struct {
	Event *RewardsRewardReceived // Event containing the contract specifics and raw log

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
func (it *RewardsRewardReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardsRewardReceived)
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
		it.Event = new(RewardsRewardReceived)
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
func (it *RewardsRewardReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardsRewardReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardsRewardReceived represents a RewardReceived event raised by the Rewards contract.
type RewardsRewardReceived struct {
	Keep   [32]byte
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardReceived is a free log retrieval operation binding the contract event 0xf8fcbb083cc485f0dc726d4235dbec7c9c5c03af58254f353627503577216170.
//
// Solidity: event RewardReceived(bytes32 keep, uint256 amount)
func (_Rewards *RewardsFilterer) FilterRewardReceived(opts *bind.FilterOpts) (*RewardsRewardReceivedIterator, error) {

	logs, sub, err := _Rewards.contract.FilterLogs(opts, "RewardReceived")
	if err != nil {
		return nil, err
	}
	return &RewardsRewardReceivedIterator{contract: _Rewards.contract, event: "RewardReceived", logs: logs, sub: sub}, nil
}

// WatchRewardReceived is a free log subscription operation binding the contract event 0xf8fcbb083cc485f0dc726d4235dbec7c9c5c03af58254f353627503577216170.
//
// Solidity: event RewardReceived(bytes32 keep, uint256 amount)
func (_Rewards *RewardsFilterer) WatchRewardReceived(opts *bind.WatchOpts, sink chan<- *RewardsRewardReceived) (event.Subscription, error) {

	logs, sub, err := _Rewards.contract.WatchLogs(opts, "RewardReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardsRewardReceived)
				if err := _Rewards.contract.UnpackLog(event, "RewardReceived", log); err != nil {
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
func (_Rewards *RewardsFilterer) ParseRewardReceived(log types.Log) (*RewardsRewardReceived, error) {
	event := new(RewardsRewardReceived)
	if err := _Rewards.contract.UnpackLog(event, "RewardReceived", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RewardsUpgradeFinalizedIterator is returned from FilterUpgradeFinalized and is used to iterate over the raw logs and unpacked data for UpgradeFinalized events raised by the Rewards contract.
type RewardsUpgradeFinalizedIterator struct {
	Event *RewardsUpgradeFinalized // Event containing the contract specifics and raw log

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
func (it *RewardsUpgradeFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardsUpgradeFinalized)
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
		it.Event = new(RewardsUpgradeFinalized)
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
func (it *RewardsUpgradeFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardsUpgradeFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardsUpgradeFinalized represents a UpgradeFinalized event raised by the Rewards contract.
type RewardsUpgradeFinalized struct {
	AmountTransferred *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgradeFinalized is a free log retrieval operation binding the contract event 0x89f5b857dbf5a32a4fa55162349dc1146850c9b76ff05bedf3a47b3f14095c0a.
//
// Solidity: event UpgradeFinalized(uint256 amountTransferred)
func (_Rewards *RewardsFilterer) FilterUpgradeFinalized(opts *bind.FilterOpts) (*RewardsUpgradeFinalizedIterator, error) {

	logs, sub, err := _Rewards.contract.FilterLogs(opts, "UpgradeFinalized")
	if err != nil {
		return nil, err
	}
	return &RewardsUpgradeFinalizedIterator{contract: _Rewards.contract, event: "UpgradeFinalized", logs: logs, sub: sub}, nil
}

// WatchUpgradeFinalized is a free log subscription operation binding the contract event 0x89f5b857dbf5a32a4fa55162349dc1146850c9b76ff05bedf3a47b3f14095c0a.
//
// Solidity: event UpgradeFinalized(uint256 amountTransferred)
func (_Rewards *RewardsFilterer) WatchUpgradeFinalized(opts *bind.WatchOpts, sink chan<- *RewardsUpgradeFinalized) (event.Subscription, error) {

	logs, sub, err := _Rewards.contract.WatchLogs(opts, "UpgradeFinalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardsUpgradeFinalized)
				if err := _Rewards.contract.UnpackLog(event, "UpgradeFinalized", log); err != nil {
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
func (_Rewards *RewardsFilterer) ParseUpgradeFinalized(log types.Log) (*RewardsUpgradeFinalized, error) {
	event := new(RewardsUpgradeFinalized)
	if err := _Rewards.contract.UnpackLog(event, "UpgradeFinalized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RewardsUpgradeInitiatedIterator is returned from FilterUpgradeInitiated and is used to iterate over the raw logs and unpacked data for UpgradeInitiated events raised by the Rewards contract.
type RewardsUpgradeInitiatedIterator struct {
	Event *RewardsUpgradeInitiated // Event containing the contract specifics and raw log

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
func (it *RewardsUpgradeInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardsUpgradeInitiated)
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
		it.Event = new(RewardsUpgradeInitiated)
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
func (it *RewardsUpgradeInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardsUpgradeInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardsUpgradeInitiated represents a UpgradeInitiated event raised by the Rewards contract.
type RewardsUpgradeInitiated struct {
	NewRewardsContract common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpgradeInitiated is a free log retrieval operation binding the contract event 0x4508cb0974b8dda93dec7c130aa186db2705fcee58bc00d1334726761835aca1.
//
// Solidity: event UpgradeInitiated(address newRewardsContract)
func (_Rewards *RewardsFilterer) FilterUpgradeInitiated(opts *bind.FilterOpts) (*RewardsUpgradeInitiatedIterator, error) {

	logs, sub, err := _Rewards.contract.FilterLogs(opts, "UpgradeInitiated")
	if err != nil {
		return nil, err
	}
	return &RewardsUpgradeInitiatedIterator{contract: _Rewards.contract, event: "UpgradeInitiated", logs: logs, sub: sub}, nil
}

// WatchUpgradeInitiated is a free log subscription operation binding the contract event 0x4508cb0974b8dda93dec7c130aa186db2705fcee58bc00d1334726761835aca1.
//
// Solidity: event UpgradeInitiated(address newRewardsContract)
func (_Rewards *RewardsFilterer) WatchUpgradeInitiated(opts *bind.WatchOpts, sink chan<- *RewardsUpgradeInitiated) (event.Subscription, error) {

	logs, sub, err := _Rewards.contract.WatchLogs(opts, "UpgradeInitiated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardsUpgradeInitiated)
				if err := _Rewards.contract.UnpackLog(event, "UpgradeInitiated", log); err != nil {
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
func (_Rewards *RewardsFilterer) ParseUpgradeInitiated(log types.Log) (*RewardsUpgradeInitiated, error) {
	event := new(RewardsUpgradeInitiated)
	if err := _Rewards.contract.UnpackLog(event, "UpgradeInitiated", log); err != nil {
		return nil, err
	}
	return event, nil
}
