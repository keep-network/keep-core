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

// BeaconRewardsABI is the input ABI used to generate the binding from.
const BeaconRewardsABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakingContract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keep\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountTransferred\",\"type\":\"uint256\"}],\"name\":\"UpgradeFinalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRewardsContract\",\"type\":\"address\"}],\"name\":\"UpgradeInitiated\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"allocateRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dispensedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"eligibleButTerminated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"eligibleForReward\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"eligibleForReward\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"endOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeRewardsUpgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"firstIntervalStart\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"funded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"getAllocatedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getIntervalCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"getIntervalWeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newRewardsContract\",\"type\":\"address\"}],\"name\":\"initiateRewardsUpgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervalKeepsProcessed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"intervalOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervalWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"isAllocated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"isFinished\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"isTerminated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"keepsInInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"markAsFunded\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumKeepsPerInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newRewardsContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keepIdentifier\",\"type\":\"bytes32\"}],\"name\":\"receiveReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"receiveReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"keepIdentifiers\",\"type\":\"bytes32[]\"}],\"name\":\"receiveRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"groupIndices\",\"type\":\"uint256[]\"}],\"name\":\"receiveRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keepIdentifier\",\"type\":\"bytes32\"}],\"name\":\"reportTermination\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"reportTermination\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"groupIndices\",\"type\":\"uint256[]\"}],\"name\":\"reportTerminations\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"keepIdentifiers\",\"type\":\"bytes32[]\"}],\"name\":\"reportTerminations\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keep\",\"type\":\"bytes32\"}],\"name\":\"rewardClaimed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"startOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"termLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractKeepToken\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"unallocatedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradeFinalizedTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradeInitiatedTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// BeaconRewards is an auto generated Go binding around an Ethereum contract.
type BeaconRewards struct {
	BeaconRewardsCaller     // Read-only binding to the contract
	BeaconRewardsTransactor // Write-only binding to the contract
	BeaconRewardsFilterer   // Log filterer for contract events
}

// BeaconRewardsCaller is an auto generated read-only Go binding around an Ethereum contract.
type BeaconRewardsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconRewardsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BeaconRewardsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconRewardsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BeaconRewardsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconRewardsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BeaconRewardsSession struct {
	Contract     *BeaconRewards    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BeaconRewardsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BeaconRewardsCallerSession struct {
	Contract *BeaconRewardsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// BeaconRewardsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BeaconRewardsTransactorSession struct {
	Contract     *BeaconRewardsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BeaconRewardsRaw is an auto generated low-level Go binding around an Ethereum contract.
type BeaconRewardsRaw struct {
	Contract *BeaconRewards // Generic contract binding to access the raw methods on
}

// BeaconRewardsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BeaconRewardsCallerRaw struct {
	Contract *BeaconRewardsCaller // Generic read-only contract binding to access the raw methods on
}

// BeaconRewardsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BeaconRewardsTransactorRaw struct {
	Contract *BeaconRewardsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBeaconRewards creates a new instance of BeaconRewards, bound to a specific deployed contract.
func NewBeaconRewards(address common.Address, backend bind.ContractBackend) (*BeaconRewards, error) {
	contract, err := bindBeaconRewards(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconRewards{BeaconRewardsCaller: BeaconRewardsCaller{contract: contract}, BeaconRewardsTransactor: BeaconRewardsTransactor{contract: contract}, BeaconRewardsFilterer: BeaconRewardsFilterer{contract: contract}}, nil
}

// NewBeaconRewardsCaller creates a new read-only instance of BeaconRewards, bound to a specific deployed contract.
func NewBeaconRewardsCaller(address common.Address, caller bind.ContractCaller) (*BeaconRewardsCaller, error) {
	contract, err := bindBeaconRewards(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconRewardsCaller{contract: contract}, nil
}

// NewBeaconRewardsTransactor creates a new write-only instance of BeaconRewards, bound to a specific deployed contract.
func NewBeaconRewardsTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconRewardsTransactor, error) {
	contract, err := bindBeaconRewards(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconRewardsTransactor{contract: contract}, nil
}

// NewBeaconRewardsFilterer creates a new log filterer instance of BeaconRewards, bound to a specific deployed contract.
func NewBeaconRewardsFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconRewardsFilterer, error) {
	contract, err := bindBeaconRewards(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconRewardsFilterer{contract: contract}, nil
}

// bindBeaconRewards binds a generic wrapper to an already deployed contract.
func bindBeaconRewards(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BeaconRewardsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconRewards *BeaconRewardsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconRewards.Contract.BeaconRewardsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconRewards *BeaconRewardsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconRewards.Contract.BeaconRewardsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconRewards *BeaconRewardsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconRewards.Contract.BeaconRewardsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconRewards *BeaconRewardsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconRewards.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconRewards *BeaconRewardsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconRewards.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconRewards *BeaconRewardsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconRewards.Contract.contract.Transact(opts, method, params...)
}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) DispensedRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "dispensedRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) DispensedRewards() (*big.Int, error) {
	return _BeaconRewards.Contract.DispensedRewards(&_BeaconRewards.CallOpts)
}

// DispensedRewards is a free data retrieval call binding the contract method 0xe1fcd5e8.
//
// Solidity: function dispensedRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) DispensedRewards() (*big.Int, error) {
	return _BeaconRewards.Contract.DispensedRewards(&_BeaconRewards.CallOpts)
}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) EligibleButTerminated(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "eligibleButTerminated", _keep)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) EligibleButTerminated(_keep [32]byte) (bool, error) {
	return _BeaconRewards.Contract.EligibleButTerminated(&_BeaconRewards.CallOpts, _keep)
}

// EligibleButTerminated is a free data retrieval call binding the contract method 0x97a1e434.
//
// Solidity: function eligibleButTerminated(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) EligibleButTerminated(_keep [32]byte) (bool, error) {
	return _BeaconRewards.Contract.EligibleButTerminated(&_BeaconRewards.CallOpts, _keep)
}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) EligibleForReward(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "eligibleForReward", _keep)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) EligibleForReward(_keep [32]byte) (bool, error) {
	return _BeaconRewards.Contract.EligibleForReward(&_BeaconRewards.CallOpts, _keep)
}

// EligibleForReward is a free data retrieval call binding the contract method 0x12f57063.
//
// Solidity: function eligibleForReward(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) EligibleForReward(_keep [32]byte) (bool, error) {
	return _BeaconRewards.Contract.EligibleForReward(&_BeaconRewards.CallOpts, _keep)
}

// EligibleForReward0 is a free data retrieval call binding the contract method 0x4611ae04.
//
// Solidity: function eligibleForReward(uint256 groupIndex) view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) EligibleForReward0(opts *bind.CallOpts, groupIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "eligibleForReward0", groupIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EligibleForReward0 is a free data retrieval call binding the contract method 0x4611ae04.
//
// Solidity: function eligibleForReward(uint256 groupIndex) view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) EligibleForReward0(groupIndex *big.Int) (bool, error) {
	return _BeaconRewards.Contract.EligibleForReward0(&_BeaconRewards.CallOpts, groupIndex)
}

// EligibleForReward0 is a free data retrieval call binding the contract method 0x4611ae04.
//
// Solidity: function eligibleForReward(uint256 groupIndex) view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) EligibleForReward0(groupIndex *big.Int) (bool, error) {
	return _BeaconRewards.Contract.EligibleForReward0(&_BeaconRewards.CallOpts, groupIndex)
}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) EndOf(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "endOf", interval)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) EndOf(interval *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.EndOf(&_BeaconRewards.CallOpts, interval)
}

// EndOf is a free data retrieval call binding the contract method 0xb965933e.
//
// Solidity: function endOf(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) EndOf(interval *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.EndOf(&_BeaconRewards.CallOpts, interval)
}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) FirstIntervalStart(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "firstIntervalStart")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) FirstIntervalStart() (*big.Int, error) {
	return _BeaconRewards.Contract.FirstIntervalStart(&_BeaconRewards.CallOpts)
}

// FirstIntervalStart is a free data retrieval call binding the contract method 0x95aeca48.
//
// Solidity: function firstIntervalStart() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) FirstIntervalStart() (*big.Int, error) {
	return _BeaconRewards.Contract.FirstIntervalStart(&_BeaconRewards.CallOpts)
}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) Funded(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "funded")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) Funded() (bool, error) {
	return _BeaconRewards.Contract.Funded(&_BeaconRewards.CallOpts)
}

// Funded is a free data retrieval call binding the contract method 0xf3a504f2.
//
// Solidity: function funded() view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) Funded() (bool, error) {
	return _BeaconRewards.Contract.Funded(&_BeaconRewards.CallOpts)
}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) GetAllocatedRewards(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "getAllocatedRewards", interval)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) GetAllocatedRewards(interval *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.GetAllocatedRewards(&_BeaconRewards.CallOpts, interval)
}

// GetAllocatedRewards is a free data retrieval call binding the contract method 0x35c9ef71.
//
// Solidity: function getAllocatedRewards(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) GetAllocatedRewards(interval *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.GetAllocatedRewards(&_BeaconRewards.CallOpts, interval)
}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) GetIntervalCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "getIntervalCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) GetIntervalCount() (*big.Int, error) {
	return _BeaconRewards.Contract.GetIntervalCount(&_BeaconRewards.CallOpts)
}

// GetIntervalCount is a free data retrieval call binding the contract method 0x183affa1.
//
// Solidity: function getIntervalCount() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) GetIntervalCount() (*big.Int, error) {
	return _BeaconRewards.Contract.GetIntervalCount(&_BeaconRewards.CallOpts)
}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) GetIntervalWeight(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "getIntervalWeight", interval)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) GetIntervalWeight(interval *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.GetIntervalWeight(&_BeaconRewards.CallOpts, interval)
}

// GetIntervalWeight is a free data retrieval call binding the contract method 0x36544951.
//
// Solidity: function getIntervalWeight(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) GetIntervalWeight(interval *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.GetIntervalWeight(&_BeaconRewards.CallOpts, interval)
}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) IntervalKeepsProcessed(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "intervalKeepsProcessed", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) IntervalKeepsProcessed(arg0 *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.IntervalKeepsProcessed(&_BeaconRewards.CallOpts, arg0)
}

// IntervalKeepsProcessed is a free data retrieval call binding the contract method 0x6a823878.
//
// Solidity: function intervalKeepsProcessed(uint256 ) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) IntervalKeepsProcessed(arg0 *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.IntervalKeepsProcessed(&_BeaconRewards.CallOpts, arg0)
}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) IntervalOf(opts *bind.CallOpts, timestamp *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "intervalOf", timestamp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) IntervalOf(timestamp *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.IntervalOf(&_BeaconRewards.CallOpts, timestamp)
}

// IntervalOf is a free data retrieval call binding the contract method 0xb369af71.
//
// Solidity: function intervalOf(uint256 timestamp) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) IntervalOf(timestamp *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.IntervalOf(&_BeaconRewards.CallOpts, timestamp)
}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) IntervalWeights(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "intervalWeights", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) IntervalWeights(arg0 *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.IntervalWeights(&_BeaconRewards.CallOpts, arg0)
}

// IntervalWeights is a free data retrieval call binding the contract method 0x92672b35.
//
// Solidity: function intervalWeights(uint256 ) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) IntervalWeights(arg0 *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.IntervalWeights(&_BeaconRewards.CallOpts, arg0)
}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) IsAllocated(opts *bind.CallOpts, interval *big.Int) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "isAllocated", interval)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) IsAllocated(interval *big.Int) (bool, error) {
	return _BeaconRewards.Contract.IsAllocated(&_BeaconRewards.CallOpts, interval)
}

// IsAllocated is a free data retrieval call binding the contract method 0x0a1bb1b1.
//
// Solidity: function isAllocated(uint256 interval) view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) IsAllocated(interval *big.Int) (bool, error) {
	return _BeaconRewards.Contract.IsAllocated(&_BeaconRewards.CallOpts, interval)
}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) IsFinished(opts *bind.CallOpts, interval *big.Int) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "isFinished", interval)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) IsFinished(interval *big.Int) (bool, error) {
	return _BeaconRewards.Contract.IsFinished(&_BeaconRewards.CallOpts, interval)
}

// IsFinished is a free data retrieval call binding the contract method 0xe8ba6509.
//
// Solidity: function isFinished(uint256 interval) view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) IsFinished(interval *big.Int) (bool, error) {
	return _BeaconRewards.Contract.IsFinished(&_BeaconRewards.CallOpts, interval)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) IsOwner() (bool, error) {
	return _BeaconRewards.Contract.IsOwner(&_BeaconRewards.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) IsOwner() (bool, error) {
	return _BeaconRewards.Contract.IsOwner(&_BeaconRewards.CallOpts)
}

// IsTerminated is a free data retrieval call binding the contract method 0xad36d6cc.
//
// Solidity: function isTerminated(uint256 groupIndex) view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) IsTerminated(opts *bind.CallOpts, groupIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "isTerminated", groupIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTerminated is a free data retrieval call binding the contract method 0xad36d6cc.
//
// Solidity: function isTerminated(uint256 groupIndex) view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) IsTerminated(groupIndex *big.Int) (bool, error) {
	return _BeaconRewards.Contract.IsTerminated(&_BeaconRewards.CallOpts, groupIndex)
}

// IsTerminated is a free data retrieval call binding the contract method 0xad36d6cc.
//
// Solidity: function isTerminated(uint256 groupIndex) view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) IsTerminated(groupIndex *big.Int) (bool, error) {
	return _BeaconRewards.Contract.IsTerminated(&_BeaconRewards.CallOpts, groupIndex)
}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) MinimumKeepsPerInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "minimumKeepsPerInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) MinimumKeepsPerInterval() (*big.Int, error) {
	return _BeaconRewards.Contract.MinimumKeepsPerInterval(&_BeaconRewards.CallOpts)
}

// MinimumKeepsPerInterval is a free data retrieval call binding the contract method 0xc70cee6a.
//
// Solidity: function minimumKeepsPerInterval() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) MinimumKeepsPerInterval() (*big.Int, error) {
	return _BeaconRewards.Contract.MinimumKeepsPerInterval(&_BeaconRewards.CallOpts)
}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() view returns(address)
func (_BeaconRewards *BeaconRewardsCaller) NewRewardsContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "newRewardsContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() view returns(address)
func (_BeaconRewards *BeaconRewardsSession) NewRewardsContract() (common.Address, error) {
	return _BeaconRewards.Contract.NewRewardsContract(&_BeaconRewards.CallOpts)
}

// NewRewardsContract is a free data retrieval call binding the contract method 0x81b83c33.
//
// Solidity: function newRewardsContract() view returns(address)
func (_BeaconRewards *BeaconRewardsCallerSession) NewRewardsContract() (common.Address, error) {
	return _BeaconRewards.Contract.NewRewardsContract(&_BeaconRewards.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BeaconRewards *BeaconRewardsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BeaconRewards *BeaconRewardsSession) Owner() (common.Address, error) {
	return _BeaconRewards.Contract.Owner(&_BeaconRewards.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BeaconRewards *BeaconRewardsCallerSession) Owner() (common.Address, error) {
	return _BeaconRewards.Contract.Owner(&_BeaconRewards.CallOpts)
}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsCaller) RewardClaimed(opts *bind.CallOpts, _keep [32]byte) (bool, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "rewardClaimed", _keep)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsSession) RewardClaimed(_keep [32]byte) (bool, error) {
	return _BeaconRewards.Contract.RewardClaimed(&_BeaconRewards.CallOpts, _keep)
}

// RewardClaimed is a free data retrieval call binding the contract method 0xce2830ad.
//
// Solidity: function rewardClaimed(bytes32 _keep) view returns(bool)
func (_BeaconRewards *BeaconRewardsCallerSession) RewardClaimed(_keep [32]byte) (bool, error) {
	return _BeaconRewards.Contract.RewardClaimed(&_BeaconRewards.CallOpts, _keep)
}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) StartOf(opts *bind.CallOpts, interval *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "startOf", interval)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) StartOf(interval *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.StartOf(&_BeaconRewards.CallOpts, interval)
}

// StartOf is a free data retrieval call binding the contract method 0xb41e7f47.
//
// Solidity: function startOf(uint256 interval) view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) StartOf(interval *big.Int) (*big.Int, error) {
	return _BeaconRewards.Contract.StartOf(&_BeaconRewards.CallOpts, interval)
}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) TermLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "termLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) TermLength() (*big.Int, error) {
	return _BeaconRewards.Contract.TermLength(&_BeaconRewards.CallOpts)
}

// TermLength is a free data retrieval call binding the contract method 0x79b18cd1.
//
// Solidity: function termLength() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) TermLength() (*big.Int, error) {
	return _BeaconRewards.Contract.TermLength(&_BeaconRewards.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_BeaconRewards *BeaconRewardsCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_BeaconRewards *BeaconRewardsSession) Token() (common.Address, error) {
	return _BeaconRewards.Contract.Token(&_BeaconRewards.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_BeaconRewards *BeaconRewardsCallerSession) Token() (common.Address, error) {
	return _BeaconRewards.Contract.Token(&_BeaconRewards.CallOpts)
}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) TotalRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "totalRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) TotalRewards() (*big.Int, error) {
	return _BeaconRewards.Contract.TotalRewards(&_BeaconRewards.CallOpts)
}

// TotalRewards is a free data retrieval call binding the contract method 0x0e15561a.
//
// Solidity: function totalRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) TotalRewards() (*big.Int, error) {
	return _BeaconRewards.Contract.TotalRewards(&_BeaconRewards.CallOpts)
}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) UnallocatedRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "unallocatedRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) UnallocatedRewards() (*big.Int, error) {
	return _BeaconRewards.Contract.UnallocatedRewards(&_BeaconRewards.CallOpts)
}

// UnallocatedRewards is a free data retrieval call binding the contract method 0x2b6fdc7e.
//
// Solidity: function unallocatedRewards() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) UnallocatedRewards() (*big.Int, error) {
	return _BeaconRewards.Contract.UnallocatedRewards(&_BeaconRewards.CallOpts)
}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) UpgradeFinalizedTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "upgradeFinalizedTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) UpgradeFinalizedTimestamp() (*big.Int, error) {
	return _BeaconRewards.Contract.UpgradeFinalizedTimestamp(&_BeaconRewards.CallOpts)
}

// UpgradeFinalizedTimestamp is a free data retrieval call binding the contract method 0xce2b3e3c.
//
// Solidity: function upgradeFinalizedTimestamp() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) UpgradeFinalizedTimestamp() (*big.Int, error) {
	return _BeaconRewards.Contract.UpgradeFinalizedTimestamp(&_BeaconRewards.CallOpts)
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCaller) UpgradeInitiatedTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconRewards.contract.Call(opts, &out, "upgradeInitiatedTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() view returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) UpgradeInitiatedTimestamp() (*big.Int, error) {
	return _BeaconRewards.Contract.UpgradeInitiatedTimestamp(&_BeaconRewards.CallOpts)
}

// UpgradeInitiatedTimestamp is a free data retrieval call binding the contract method 0x95131526.
//
// Solidity: function upgradeInitiatedTimestamp() view returns(uint256)
func (_BeaconRewards *BeaconRewardsCallerSession) UpgradeInitiatedTimestamp() (*big.Int, error) {
	return _BeaconRewards.Contract.UpgradeInitiatedTimestamp(&_BeaconRewards.CallOpts)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_BeaconRewards *BeaconRewardsTransactor) AllocateRewards(opts *bind.TransactOpts, interval *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "allocateRewards", interval)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_BeaconRewards *BeaconRewardsSession) AllocateRewards(interval *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.AllocateRewards(&_BeaconRewards.TransactOpts, interval)
}

// AllocateRewards is a paid mutator transaction binding the contract method 0x28fc33c7.
//
// Solidity: function allocateRewards(uint256 interval) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) AllocateRewards(interval *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.AllocateRewards(&_BeaconRewards.TransactOpts, interval)
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_BeaconRewards *BeaconRewardsTransactor) FinalizeRewardsUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "finalizeRewardsUpgrade")
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_BeaconRewards *BeaconRewardsSession) FinalizeRewardsUpgrade() (*types.Transaction, error) {
	return _BeaconRewards.Contract.FinalizeRewardsUpgrade(&_BeaconRewards.TransactOpts)
}

// FinalizeRewardsUpgrade is a paid mutator transaction binding the contract method 0xe27000e7.
//
// Solidity: function finalizeRewardsUpgrade() returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) FinalizeRewardsUpgrade() (*types.Transaction, error) {
	return _BeaconRewards.Contract.FinalizeRewardsUpgrade(&_BeaconRewards.TransactOpts)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_BeaconRewards *BeaconRewardsTransactor) InitiateRewardsUpgrade(opts *bind.TransactOpts, _newRewardsContract common.Address) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "initiateRewardsUpgrade", _newRewardsContract)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_BeaconRewards *BeaconRewardsSession) InitiateRewardsUpgrade(_newRewardsContract common.Address) (*types.Transaction, error) {
	return _BeaconRewards.Contract.InitiateRewardsUpgrade(&_BeaconRewards.TransactOpts, _newRewardsContract)
}

// InitiateRewardsUpgrade is a paid mutator transaction binding the contract method 0xce21d8cf.
//
// Solidity: function initiateRewardsUpgrade(address _newRewardsContract) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) InitiateRewardsUpgrade(_newRewardsContract common.Address) (*types.Transaction, error) {
	return _BeaconRewards.Contract.InitiateRewardsUpgrade(&_BeaconRewards.TransactOpts, _newRewardsContract)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_BeaconRewards *BeaconRewardsTransactor) KeepsInInterval(opts *bind.TransactOpts, interval *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "keepsInInterval", interval)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_BeaconRewards *BeaconRewardsSession) KeepsInInterval(interval *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.KeepsInInterval(&_BeaconRewards.TransactOpts, interval)
}

// KeepsInInterval is a paid mutator transaction binding the contract method 0xaeae364d.
//
// Solidity: function keepsInInterval(uint256 interval) returns(uint256)
func (_BeaconRewards *BeaconRewardsTransactorSession) KeepsInInterval(interval *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.KeepsInInterval(&_BeaconRewards.TransactOpts, interval)
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_BeaconRewards *BeaconRewardsTransactor) MarkAsFunded(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "markAsFunded")
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_BeaconRewards *BeaconRewardsSession) MarkAsFunded() (*types.Transaction, error) {
	return _BeaconRewards.Contract.MarkAsFunded(&_BeaconRewards.TransactOpts)
}

// MarkAsFunded is a paid mutator transaction binding the contract method 0x0da8882f.
//
// Solidity: function markAsFunded() returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) MarkAsFunded() (*types.Transaction, error) {
	return _BeaconRewards.Contract.MarkAsFunded(&_BeaconRewards.TransactOpts)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReceiveApproval(opts *bind.TransactOpts, _from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "receiveApproval", _from, _value, _token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_BeaconRewards *BeaconRewardsSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveApproval(&_BeaconRewards.TransactOpts, _from, _value, _token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes ) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveApproval(&_BeaconRewards.TransactOpts, _from, _value, _token, arg3)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReceiveReward(opts *bind.TransactOpts, keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "receiveReward", keepIdentifier)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_BeaconRewards *BeaconRewardsSession) ReceiveReward(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveReward(&_BeaconRewards.TransactOpts, keepIdentifier)
}

// ReceiveReward is a paid mutator transaction binding the contract method 0x68ec9c42.
//
// Solidity: function receiveReward(bytes32 keepIdentifier) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReceiveReward(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveReward(&_BeaconRewards.TransactOpts, keepIdentifier)
}

// ReceiveReward0 is a paid mutator transaction binding the contract method 0xae10d5ea.
//
// Solidity: function receiveReward(uint256 groupIndex) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReceiveReward0(opts *bind.TransactOpts, groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "receiveReward0", groupIndex)
}

// ReceiveReward0 is a paid mutator transaction binding the contract method 0xae10d5ea.
//
// Solidity: function receiveReward(uint256 groupIndex) returns()
func (_BeaconRewards *BeaconRewardsSession) ReceiveReward0(groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveReward0(&_BeaconRewards.TransactOpts, groupIndex)
}

// ReceiveReward0 is a paid mutator transaction binding the contract method 0xae10d5ea.
//
// Solidity: function receiveReward(uint256 groupIndex) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReceiveReward0(groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveReward0(&_BeaconRewards.TransactOpts, groupIndex)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReceiveRewards(opts *bind.TransactOpts, keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "receiveRewards", keepIdentifiers)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_BeaconRewards *BeaconRewardsSession) ReceiveRewards(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveRewards(&_BeaconRewards.TransactOpts, keepIdentifiers)
}

// ReceiveRewards is a paid mutator transaction binding the contract method 0x45016312.
//
// Solidity: function receiveRewards(bytes32[] keepIdentifiers) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReceiveRewards(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveRewards(&_BeaconRewards.TransactOpts, keepIdentifiers)
}

// ReceiveRewards0 is a paid mutator transaction binding the contract method 0x73c6d862.
//
// Solidity: function receiveRewards(uint256[] groupIndices) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReceiveRewards0(opts *bind.TransactOpts, groupIndices []*big.Int) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "receiveRewards0", groupIndices)
}

// ReceiveRewards0 is a paid mutator transaction binding the contract method 0x73c6d862.
//
// Solidity: function receiveRewards(uint256[] groupIndices) returns()
func (_BeaconRewards *BeaconRewardsSession) ReceiveRewards0(groupIndices []*big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveRewards0(&_BeaconRewards.TransactOpts, groupIndices)
}

// ReceiveRewards0 is a paid mutator transaction binding the contract method 0x73c6d862.
//
// Solidity: function receiveRewards(uint256[] groupIndices) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReceiveRewards0(groupIndices []*big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReceiveRewards0(&_BeaconRewards.TransactOpts, groupIndices)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BeaconRewards *BeaconRewardsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BeaconRewards *BeaconRewardsSession) RenounceOwnership() (*types.Transaction, error) {
	return _BeaconRewards.Contract.RenounceOwnership(&_BeaconRewards.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BeaconRewards.Contract.RenounceOwnership(&_BeaconRewards.TransactOpts)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReportTermination(opts *bind.TransactOpts, keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "reportTermination", keepIdentifier)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_BeaconRewards *BeaconRewardsSession) ReportTermination(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReportTermination(&_BeaconRewards.TransactOpts, keepIdentifier)
}

// ReportTermination is a paid mutator transaction binding the contract method 0x86d954c0.
//
// Solidity: function reportTermination(bytes32 keepIdentifier) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReportTermination(keepIdentifier [32]byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReportTermination(&_BeaconRewards.TransactOpts, keepIdentifier)
}

// ReportTermination0 is a paid mutator transaction binding the contract method 0xaa5c63a9.
//
// Solidity: function reportTermination(uint256 groupIndex) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReportTermination0(opts *bind.TransactOpts, groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "reportTermination0", groupIndex)
}

// ReportTermination0 is a paid mutator transaction binding the contract method 0xaa5c63a9.
//
// Solidity: function reportTermination(uint256 groupIndex) returns()
func (_BeaconRewards *BeaconRewardsSession) ReportTermination0(groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReportTermination0(&_BeaconRewards.TransactOpts, groupIndex)
}

// ReportTermination0 is a paid mutator transaction binding the contract method 0xaa5c63a9.
//
// Solidity: function reportTermination(uint256 groupIndex) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReportTermination0(groupIndex *big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReportTermination0(&_BeaconRewards.TransactOpts, groupIndex)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0x6c21aba7.
//
// Solidity: function reportTerminations(uint256[] groupIndices) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReportTerminations(opts *bind.TransactOpts, groupIndices []*big.Int) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "reportTerminations", groupIndices)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0x6c21aba7.
//
// Solidity: function reportTerminations(uint256[] groupIndices) returns()
func (_BeaconRewards *BeaconRewardsSession) ReportTerminations(groupIndices []*big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReportTerminations(&_BeaconRewards.TransactOpts, groupIndices)
}

// ReportTerminations is a paid mutator transaction binding the contract method 0x6c21aba7.
//
// Solidity: function reportTerminations(uint256[] groupIndices) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReportTerminations(groupIndices []*big.Int) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReportTerminations(&_BeaconRewards.TransactOpts, groupIndices)
}

// ReportTerminations0 is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_BeaconRewards *BeaconRewardsTransactor) ReportTerminations0(opts *bind.TransactOpts, keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "reportTerminations0", keepIdentifiers)
}

// ReportTerminations0 is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_BeaconRewards *BeaconRewardsSession) ReportTerminations0(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReportTerminations0(&_BeaconRewards.TransactOpts, keepIdentifiers)
}

// ReportTerminations0 is a paid mutator transaction binding the contract method 0xab7239c0.
//
// Solidity: function reportTerminations(bytes32[] keepIdentifiers) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) ReportTerminations0(keepIdentifiers [][32]byte) (*types.Transaction, error) {
	return _BeaconRewards.Contract.ReportTerminations0(&_BeaconRewards.TransactOpts, keepIdentifiers)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BeaconRewards *BeaconRewardsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BeaconRewards.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BeaconRewards *BeaconRewardsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BeaconRewards.Contract.TransferOwnership(&_BeaconRewards.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BeaconRewards *BeaconRewardsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BeaconRewards.Contract.TransferOwnership(&_BeaconRewards.TransactOpts, newOwner)
}

// BeaconRewardsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BeaconRewards contract.
type BeaconRewardsOwnershipTransferredIterator struct {
	Event *BeaconRewardsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BeaconRewardsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconRewardsOwnershipTransferred)
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
		it.Event = new(BeaconRewardsOwnershipTransferred)
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
func (it *BeaconRewardsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconRewardsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconRewardsOwnershipTransferred represents a OwnershipTransferred event raised by the BeaconRewards contract.
type BeaconRewardsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BeaconRewards *BeaconRewardsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BeaconRewardsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BeaconRewards.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BeaconRewardsOwnershipTransferredIterator{contract: _BeaconRewards.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BeaconRewards *BeaconRewardsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BeaconRewardsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BeaconRewards.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconRewardsOwnershipTransferred)
				if err := _BeaconRewards.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BeaconRewards *BeaconRewardsFilterer) ParseOwnershipTransferred(log types.Log) (*BeaconRewardsOwnershipTransferred, error) {
	event := new(BeaconRewardsOwnershipTransferred)
	if err := _BeaconRewards.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconRewardsRewardReceivedIterator is returned from FilterRewardReceived and is used to iterate over the raw logs and unpacked data for RewardReceived events raised by the BeaconRewards contract.
type BeaconRewardsRewardReceivedIterator struct {
	Event *BeaconRewardsRewardReceived // Event containing the contract specifics and raw log

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
func (it *BeaconRewardsRewardReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconRewardsRewardReceived)
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
		it.Event = new(BeaconRewardsRewardReceived)
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
func (it *BeaconRewardsRewardReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconRewardsRewardReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconRewardsRewardReceived represents a RewardReceived event raised by the BeaconRewards contract.
type BeaconRewardsRewardReceived struct {
	Keep   [32]byte
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardReceived is a free log retrieval operation binding the contract event 0xf8fcbb083cc485f0dc726d4235dbec7c9c5c03af58254f353627503577216170.
//
// Solidity: event RewardReceived(bytes32 keep, uint256 amount)
func (_BeaconRewards *BeaconRewardsFilterer) FilterRewardReceived(opts *bind.FilterOpts) (*BeaconRewardsRewardReceivedIterator, error) {

	logs, sub, err := _BeaconRewards.contract.FilterLogs(opts, "RewardReceived")
	if err != nil {
		return nil, err
	}
	return &BeaconRewardsRewardReceivedIterator{contract: _BeaconRewards.contract, event: "RewardReceived", logs: logs, sub: sub}, nil
}

// WatchRewardReceived is a free log subscription operation binding the contract event 0xf8fcbb083cc485f0dc726d4235dbec7c9c5c03af58254f353627503577216170.
//
// Solidity: event RewardReceived(bytes32 keep, uint256 amount)
func (_BeaconRewards *BeaconRewardsFilterer) WatchRewardReceived(opts *bind.WatchOpts, sink chan<- *BeaconRewardsRewardReceived) (event.Subscription, error) {

	logs, sub, err := _BeaconRewards.contract.WatchLogs(opts, "RewardReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconRewardsRewardReceived)
				if err := _BeaconRewards.contract.UnpackLog(event, "RewardReceived", log); err != nil {
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
func (_BeaconRewards *BeaconRewardsFilterer) ParseRewardReceived(log types.Log) (*BeaconRewardsRewardReceived, error) {
	event := new(BeaconRewardsRewardReceived)
	if err := _BeaconRewards.contract.UnpackLog(event, "RewardReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconRewardsUpgradeFinalizedIterator is returned from FilterUpgradeFinalized and is used to iterate over the raw logs and unpacked data for UpgradeFinalized events raised by the BeaconRewards contract.
type BeaconRewardsUpgradeFinalizedIterator struct {
	Event *BeaconRewardsUpgradeFinalized // Event containing the contract specifics and raw log

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
func (it *BeaconRewardsUpgradeFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconRewardsUpgradeFinalized)
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
		it.Event = new(BeaconRewardsUpgradeFinalized)
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
func (it *BeaconRewardsUpgradeFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconRewardsUpgradeFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconRewardsUpgradeFinalized represents a UpgradeFinalized event raised by the BeaconRewards contract.
type BeaconRewardsUpgradeFinalized struct {
	AmountTransferred *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgradeFinalized is a free log retrieval operation binding the contract event 0x89f5b857dbf5a32a4fa55162349dc1146850c9b76ff05bedf3a47b3f14095c0a.
//
// Solidity: event UpgradeFinalized(uint256 amountTransferred)
func (_BeaconRewards *BeaconRewardsFilterer) FilterUpgradeFinalized(opts *bind.FilterOpts) (*BeaconRewardsUpgradeFinalizedIterator, error) {

	logs, sub, err := _BeaconRewards.contract.FilterLogs(opts, "UpgradeFinalized")
	if err != nil {
		return nil, err
	}
	return &BeaconRewardsUpgradeFinalizedIterator{contract: _BeaconRewards.contract, event: "UpgradeFinalized", logs: logs, sub: sub}, nil
}

// WatchUpgradeFinalized is a free log subscription operation binding the contract event 0x89f5b857dbf5a32a4fa55162349dc1146850c9b76ff05bedf3a47b3f14095c0a.
//
// Solidity: event UpgradeFinalized(uint256 amountTransferred)
func (_BeaconRewards *BeaconRewardsFilterer) WatchUpgradeFinalized(opts *bind.WatchOpts, sink chan<- *BeaconRewardsUpgradeFinalized) (event.Subscription, error) {

	logs, sub, err := _BeaconRewards.contract.WatchLogs(opts, "UpgradeFinalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconRewardsUpgradeFinalized)
				if err := _BeaconRewards.contract.UnpackLog(event, "UpgradeFinalized", log); err != nil {
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
func (_BeaconRewards *BeaconRewardsFilterer) ParseUpgradeFinalized(log types.Log) (*BeaconRewardsUpgradeFinalized, error) {
	event := new(BeaconRewardsUpgradeFinalized)
	if err := _BeaconRewards.contract.UnpackLog(event, "UpgradeFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconRewardsUpgradeInitiatedIterator is returned from FilterUpgradeInitiated and is used to iterate over the raw logs and unpacked data for UpgradeInitiated events raised by the BeaconRewards contract.
type BeaconRewardsUpgradeInitiatedIterator struct {
	Event *BeaconRewardsUpgradeInitiated // Event containing the contract specifics and raw log

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
func (it *BeaconRewardsUpgradeInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconRewardsUpgradeInitiated)
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
		it.Event = new(BeaconRewardsUpgradeInitiated)
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
func (it *BeaconRewardsUpgradeInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconRewardsUpgradeInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconRewardsUpgradeInitiated represents a UpgradeInitiated event raised by the BeaconRewards contract.
type BeaconRewardsUpgradeInitiated struct {
	NewRewardsContract common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpgradeInitiated is a free log retrieval operation binding the contract event 0x4508cb0974b8dda93dec7c130aa186db2705fcee58bc00d1334726761835aca1.
//
// Solidity: event UpgradeInitiated(address newRewardsContract)
func (_BeaconRewards *BeaconRewardsFilterer) FilterUpgradeInitiated(opts *bind.FilterOpts) (*BeaconRewardsUpgradeInitiatedIterator, error) {

	logs, sub, err := _BeaconRewards.contract.FilterLogs(opts, "UpgradeInitiated")
	if err != nil {
		return nil, err
	}
	return &BeaconRewardsUpgradeInitiatedIterator{contract: _BeaconRewards.contract, event: "UpgradeInitiated", logs: logs, sub: sub}, nil
}

// WatchUpgradeInitiated is a free log subscription operation binding the contract event 0x4508cb0974b8dda93dec7c130aa186db2705fcee58bc00d1334726761835aca1.
//
// Solidity: event UpgradeInitiated(address newRewardsContract)
func (_BeaconRewards *BeaconRewardsFilterer) WatchUpgradeInitiated(opts *bind.WatchOpts, sink chan<- *BeaconRewardsUpgradeInitiated) (event.Subscription, error) {

	logs, sub, err := _BeaconRewards.contract.WatchLogs(opts, "UpgradeInitiated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconRewardsUpgradeInitiated)
				if err := _BeaconRewards.contract.UnpackLog(event, "UpgradeInitiated", log); err != nil {
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
func (_BeaconRewards *BeaconRewardsFilterer) ParseUpgradeInitiated(log types.Log) (*BeaconRewardsUpgradeInitiated, error) {
	event := new(BeaconRewardsUpgradeInitiated)
	if err := _BeaconRewards.contract.UnpackLog(event, "UpgradeInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
