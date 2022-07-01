// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SortitionPoolMetaData contains all meta data concerning the SortitionPool contract.
var SortitionPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20WithPermit\",\"name\":\"_rewardToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_poolWeightDivisor\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"ids\",\"type\":\"uint32[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"until\",\"type\":\"uint256\"}],\"name\":\"IneligibleForRewards\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"}],\"name\":\"RewardEligibilityRestored\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"operator\",\"type\":\"uint32\"}],\"name\":\"canRestoreRewardEligibility\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getAvailableRewards\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"id\",\"type\":\"uint32\"}],\"name\":\"getIDOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"ids\",\"type\":\"uint32[]\"}],\"name\":\"getIDOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getOperatorID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getPoolWeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ineligibleEarnedRewards\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"authorizedStake\",\"type\":\"uint256\"}],\"name\":\"insertOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"operator\",\"type\":\"uint32\"}],\"name\":\"isEligibleForRewards\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isLocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isOperatorInPool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isOperatorRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"authorizedStake\",\"type\":\"uint256\"}],\"name\":\"isOperatorUpToDate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"operatorsInPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolWeightDivisor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"restoreRewardEligibility\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardToken\",\"outputs\":[{\"internalType\":\"contractIERC20WithPermit\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"operator\",\"type\":\"uint32\"}],\"name\":\"rewardsEligibilityRestorableAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"groupSize\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"seed\",\"type\":\"bytes32\"}],\"name\":\"selectGroup\",\"outputs\":[{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"operators\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"until\",\"type\":\"uint256\"}],\"name\":\"setRewardIneligibility\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalWeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"authorizedStake\",\"type\":\"uint256\"}],\"name\":\"updateOperatorStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawIneligible\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"}],\"name\":\"withdrawRewards\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SortitionPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use SortitionPoolMetaData.ABI instead.
var SortitionPoolABI = SortitionPoolMetaData.ABI

// SortitionPool is an auto generated Go binding around an Ethereum contract.
type SortitionPool struct {
	SortitionPoolCaller     // Read-only binding to the contract
	SortitionPoolTransactor // Write-only binding to the contract
	SortitionPoolFilterer   // Log filterer for contract events
}

// SortitionPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type SortitionPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SortitionPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SortitionPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SortitionPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SortitionPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SortitionPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SortitionPoolSession struct {
	Contract     *SortitionPool    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SortitionPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SortitionPoolCallerSession struct {
	Contract *SortitionPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SortitionPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SortitionPoolTransactorSession struct {
	Contract     *SortitionPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SortitionPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type SortitionPoolRaw struct {
	Contract *SortitionPool // Generic contract binding to access the raw methods on
}

// SortitionPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SortitionPoolCallerRaw struct {
	Contract *SortitionPoolCaller // Generic read-only contract binding to access the raw methods on
}

// SortitionPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SortitionPoolTransactorRaw struct {
	Contract *SortitionPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSortitionPool creates a new instance of SortitionPool, bound to a specific deployed contract.
func NewSortitionPool(address common.Address, backend bind.ContractBackend) (*SortitionPool, error) {
	contract, err := bindSortitionPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SortitionPool{SortitionPoolCaller: SortitionPoolCaller{contract: contract}, SortitionPoolTransactor: SortitionPoolTransactor{contract: contract}, SortitionPoolFilterer: SortitionPoolFilterer{contract: contract}}, nil
}

// NewSortitionPoolCaller creates a new read-only instance of SortitionPool, bound to a specific deployed contract.
func NewSortitionPoolCaller(address common.Address, caller bind.ContractCaller) (*SortitionPoolCaller, error) {
	contract, err := bindSortitionPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SortitionPoolCaller{contract: contract}, nil
}

// NewSortitionPoolTransactor creates a new write-only instance of SortitionPool, bound to a specific deployed contract.
func NewSortitionPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*SortitionPoolTransactor, error) {
	contract, err := bindSortitionPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SortitionPoolTransactor{contract: contract}, nil
}

// NewSortitionPoolFilterer creates a new log filterer instance of SortitionPool, bound to a specific deployed contract.
func NewSortitionPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*SortitionPoolFilterer, error) {
	contract, err := bindSortitionPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SortitionPoolFilterer{contract: contract}, nil
}

// bindSortitionPool binds a generic wrapper to an already deployed contract.
func bindSortitionPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SortitionPoolABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SortitionPool *SortitionPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SortitionPool.Contract.SortitionPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SortitionPool *SortitionPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SortitionPool.Contract.SortitionPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SortitionPool *SortitionPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SortitionPool.Contract.SortitionPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SortitionPool *SortitionPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SortitionPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SortitionPool *SortitionPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SortitionPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SortitionPool *SortitionPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SortitionPool.Contract.contract.Transact(opts, method, params...)
}

// CanRestoreRewardEligibility is a free data retrieval call binding the contract method 0x69ba3c33.
//
// Solidity: function canRestoreRewardEligibility(uint32 operator) view returns(bool)
func (_SortitionPool *SortitionPoolCaller) CanRestoreRewardEligibility(opts *bind.CallOpts, operator uint32) (bool, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "canRestoreRewardEligibility", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanRestoreRewardEligibility is a free data retrieval call binding the contract method 0x69ba3c33.
//
// Solidity: function canRestoreRewardEligibility(uint32 operator) view returns(bool)
func (_SortitionPool *SortitionPoolSession) CanRestoreRewardEligibility(operator uint32) (bool, error) {
	return _SortitionPool.Contract.CanRestoreRewardEligibility(&_SortitionPool.CallOpts, operator)
}

// CanRestoreRewardEligibility is a free data retrieval call binding the contract method 0x69ba3c33.
//
// Solidity: function canRestoreRewardEligibility(uint32 operator) view returns(bool)
func (_SortitionPool *SortitionPoolCallerSession) CanRestoreRewardEligibility(operator uint32) (bool, error) {
	return _SortitionPool.Contract.CanRestoreRewardEligibility(&_SortitionPool.CallOpts, operator)
}

// GetAvailableRewards is a free data retrieval call binding the contract method 0x873e31fa.
//
// Solidity: function getAvailableRewards(address operator) view returns(uint96)
func (_SortitionPool *SortitionPoolCaller) GetAvailableRewards(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "getAvailableRewards", operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAvailableRewards is a free data retrieval call binding the contract method 0x873e31fa.
//
// Solidity: function getAvailableRewards(address operator) view returns(uint96)
func (_SortitionPool *SortitionPoolSession) GetAvailableRewards(operator common.Address) (*big.Int, error) {
	return _SortitionPool.Contract.GetAvailableRewards(&_SortitionPool.CallOpts, operator)
}

// GetAvailableRewards is a free data retrieval call binding the contract method 0x873e31fa.
//
// Solidity: function getAvailableRewards(address operator) view returns(uint96)
func (_SortitionPool *SortitionPoolCallerSession) GetAvailableRewards(operator common.Address) (*big.Int, error) {
	return _SortitionPool.Contract.GetAvailableRewards(&_SortitionPool.CallOpts, operator)
}

// GetIDOperator is a free data retrieval call binding the contract method 0x8871ca5d.
//
// Solidity: function getIDOperator(uint32 id) view returns(address)
func (_SortitionPool *SortitionPoolCaller) GetIDOperator(opts *bind.CallOpts, id uint32) (common.Address, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "getIDOperator", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIDOperator is a free data retrieval call binding the contract method 0x8871ca5d.
//
// Solidity: function getIDOperator(uint32 id) view returns(address)
func (_SortitionPool *SortitionPoolSession) GetIDOperator(id uint32) (common.Address, error) {
	return _SortitionPool.Contract.GetIDOperator(&_SortitionPool.CallOpts, id)
}

// GetIDOperator is a free data retrieval call binding the contract method 0x8871ca5d.
//
// Solidity: function getIDOperator(uint32 id) view returns(address)
func (_SortitionPool *SortitionPoolCallerSession) GetIDOperator(id uint32) (common.Address, error) {
	return _SortitionPool.Contract.GetIDOperator(&_SortitionPool.CallOpts, id)
}

// GetIDOperators is a free data retrieval call binding the contract method 0xf7f9a8fa.
//
// Solidity: function getIDOperators(uint32[] ids) view returns(address[])
func (_SortitionPool *SortitionPoolCaller) GetIDOperators(opts *bind.CallOpts, ids []uint32) ([]common.Address, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "getIDOperators", ids)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetIDOperators is a free data retrieval call binding the contract method 0xf7f9a8fa.
//
// Solidity: function getIDOperators(uint32[] ids) view returns(address[])
func (_SortitionPool *SortitionPoolSession) GetIDOperators(ids []uint32) ([]common.Address, error) {
	return _SortitionPool.Contract.GetIDOperators(&_SortitionPool.CallOpts, ids)
}

// GetIDOperators is a free data retrieval call binding the contract method 0xf7f9a8fa.
//
// Solidity: function getIDOperators(uint32[] ids) view returns(address[])
func (_SortitionPool *SortitionPoolCallerSession) GetIDOperators(ids []uint32) ([]common.Address, error) {
	return _SortitionPool.Contract.GetIDOperators(&_SortitionPool.CallOpts, ids)
}

// GetOperatorID is a free data retrieval call binding the contract method 0x5a48b46b.
//
// Solidity: function getOperatorID(address operator) view returns(uint32)
func (_SortitionPool *SortitionPoolCaller) GetOperatorID(opts *bind.CallOpts, operator common.Address) (uint32, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "getOperatorID", operator)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetOperatorID is a free data retrieval call binding the contract method 0x5a48b46b.
//
// Solidity: function getOperatorID(address operator) view returns(uint32)
func (_SortitionPool *SortitionPoolSession) GetOperatorID(operator common.Address) (uint32, error) {
	return _SortitionPool.Contract.GetOperatorID(&_SortitionPool.CallOpts, operator)
}

// GetOperatorID is a free data retrieval call binding the contract method 0x5a48b46b.
//
// Solidity: function getOperatorID(address operator) view returns(uint32)
func (_SortitionPool *SortitionPoolCallerSession) GetOperatorID(operator common.Address) (uint32, error) {
	return _SortitionPool.Contract.GetOperatorID(&_SortitionPool.CallOpts, operator)
}

// GetPoolWeight is a free data retrieval call binding the contract method 0x5757ed5b.
//
// Solidity: function getPoolWeight(address operator) view returns(uint256)
func (_SortitionPool *SortitionPoolCaller) GetPoolWeight(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "getPoolWeight", operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPoolWeight is a free data retrieval call binding the contract method 0x5757ed5b.
//
// Solidity: function getPoolWeight(address operator) view returns(uint256)
func (_SortitionPool *SortitionPoolSession) GetPoolWeight(operator common.Address) (*big.Int, error) {
	return _SortitionPool.Contract.GetPoolWeight(&_SortitionPool.CallOpts, operator)
}

// GetPoolWeight is a free data retrieval call binding the contract method 0x5757ed5b.
//
// Solidity: function getPoolWeight(address operator) view returns(uint256)
func (_SortitionPool *SortitionPoolCallerSession) GetPoolWeight(operator common.Address) (*big.Int, error) {
	return _SortitionPool.Contract.GetPoolWeight(&_SortitionPool.CallOpts, operator)
}

// IneligibleEarnedRewards is a free data retrieval call binding the contract method 0xa7a7d391.
//
// Solidity: function ineligibleEarnedRewards() view returns(uint96)
func (_SortitionPool *SortitionPoolCaller) IneligibleEarnedRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "ineligibleEarnedRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IneligibleEarnedRewards is a free data retrieval call binding the contract method 0xa7a7d391.
//
// Solidity: function ineligibleEarnedRewards() view returns(uint96)
func (_SortitionPool *SortitionPoolSession) IneligibleEarnedRewards() (*big.Int, error) {
	return _SortitionPool.Contract.IneligibleEarnedRewards(&_SortitionPool.CallOpts)
}

// IneligibleEarnedRewards is a free data retrieval call binding the contract method 0xa7a7d391.
//
// Solidity: function ineligibleEarnedRewards() view returns(uint96)
func (_SortitionPool *SortitionPoolCallerSession) IneligibleEarnedRewards() (*big.Int, error) {
	return _SortitionPool.Contract.IneligibleEarnedRewards(&_SortitionPool.CallOpts)
}

// IsEligibleForRewards is a free data retrieval call binding the contract method 0x5222580a.
//
// Solidity: function isEligibleForRewards(uint32 operator) view returns(bool)
func (_SortitionPool *SortitionPoolCaller) IsEligibleForRewards(opts *bind.CallOpts, operator uint32) (bool, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "isEligibleForRewards", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEligibleForRewards is a free data retrieval call binding the contract method 0x5222580a.
//
// Solidity: function isEligibleForRewards(uint32 operator) view returns(bool)
func (_SortitionPool *SortitionPoolSession) IsEligibleForRewards(operator uint32) (bool, error) {
	return _SortitionPool.Contract.IsEligibleForRewards(&_SortitionPool.CallOpts, operator)
}

// IsEligibleForRewards is a free data retrieval call binding the contract method 0x5222580a.
//
// Solidity: function isEligibleForRewards(uint32 operator) view returns(bool)
func (_SortitionPool *SortitionPoolCallerSession) IsEligibleForRewards(operator uint32) (bool, error) {
	return _SortitionPool.Contract.IsEligibleForRewards(&_SortitionPool.CallOpts, operator)
}

// IsLocked is a free data retrieval call binding the contract method 0xa4e2d634.
//
// Solidity: function isLocked() view returns(bool)
func (_SortitionPool *SortitionPoolCaller) IsLocked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "isLocked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsLocked is a free data retrieval call binding the contract method 0xa4e2d634.
//
// Solidity: function isLocked() view returns(bool)
func (_SortitionPool *SortitionPoolSession) IsLocked() (bool, error) {
	return _SortitionPool.Contract.IsLocked(&_SortitionPool.CallOpts)
}

// IsLocked is a free data retrieval call binding the contract method 0xa4e2d634.
//
// Solidity: function isLocked() view returns(bool)
func (_SortitionPool *SortitionPoolCallerSession) IsLocked() (bool, error) {
	return _SortitionPool.Contract.IsLocked(&_SortitionPool.CallOpts)
}

// IsOperatorInPool is a free data retrieval call binding the contract method 0xf7186ce0.
//
// Solidity: function isOperatorInPool(address operator) view returns(bool)
func (_SortitionPool *SortitionPoolCaller) IsOperatorInPool(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "isOperatorInPool", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperatorInPool is a free data retrieval call binding the contract method 0xf7186ce0.
//
// Solidity: function isOperatorInPool(address operator) view returns(bool)
func (_SortitionPool *SortitionPoolSession) IsOperatorInPool(operator common.Address) (bool, error) {
	return _SortitionPool.Contract.IsOperatorInPool(&_SortitionPool.CallOpts, operator)
}

// IsOperatorInPool is a free data retrieval call binding the contract method 0xf7186ce0.
//
// Solidity: function isOperatorInPool(address operator) view returns(bool)
func (_SortitionPool *SortitionPoolCallerSession) IsOperatorInPool(operator common.Address) (bool, error) {
	return _SortitionPool.Contract.IsOperatorInPool(&_SortitionPool.CallOpts, operator)
}

// IsOperatorRegistered is a free data retrieval call binding the contract method 0x6b1906f8.
//
// Solidity: function isOperatorRegistered(address operator) view returns(bool)
func (_SortitionPool *SortitionPoolCaller) IsOperatorRegistered(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "isOperatorRegistered", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperatorRegistered is a free data retrieval call binding the contract method 0x6b1906f8.
//
// Solidity: function isOperatorRegistered(address operator) view returns(bool)
func (_SortitionPool *SortitionPoolSession) IsOperatorRegistered(operator common.Address) (bool, error) {
	return _SortitionPool.Contract.IsOperatorRegistered(&_SortitionPool.CallOpts, operator)
}

// IsOperatorRegistered is a free data retrieval call binding the contract method 0x6b1906f8.
//
// Solidity: function isOperatorRegistered(address operator) view returns(bool)
func (_SortitionPool *SortitionPoolCallerSession) IsOperatorRegistered(operator common.Address) (bool, error) {
	return _SortitionPool.Contract.IsOperatorRegistered(&_SortitionPool.CallOpts, operator)
}

// IsOperatorUpToDate is a free data retrieval call binding the contract method 0x4de824f0.
//
// Solidity: function isOperatorUpToDate(address operator, uint256 authorizedStake) view returns(bool)
func (_SortitionPool *SortitionPoolCaller) IsOperatorUpToDate(opts *bind.CallOpts, operator common.Address, authorizedStake *big.Int) (bool, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "isOperatorUpToDate", operator, authorizedStake)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperatorUpToDate is a free data retrieval call binding the contract method 0x4de824f0.
//
// Solidity: function isOperatorUpToDate(address operator, uint256 authorizedStake) view returns(bool)
func (_SortitionPool *SortitionPoolSession) IsOperatorUpToDate(operator common.Address, authorizedStake *big.Int) (bool, error) {
	return _SortitionPool.Contract.IsOperatorUpToDate(&_SortitionPool.CallOpts, operator, authorizedStake)
}

// IsOperatorUpToDate is a free data retrieval call binding the contract method 0x4de824f0.
//
// Solidity: function isOperatorUpToDate(address operator, uint256 authorizedStake) view returns(bool)
func (_SortitionPool *SortitionPoolCallerSession) IsOperatorUpToDate(operator common.Address, authorizedStake *big.Int) (bool, error) {
	return _SortitionPool.Contract.IsOperatorUpToDate(&_SortitionPool.CallOpts, operator, authorizedStake)
}

// OperatorsInPool is a free data retrieval call binding the contract method 0xe7bfd899.
//
// Solidity: function operatorsInPool() view returns(uint256)
func (_SortitionPool *SortitionPoolCaller) OperatorsInPool(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "operatorsInPool")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OperatorsInPool is a free data retrieval call binding the contract method 0xe7bfd899.
//
// Solidity: function operatorsInPool() view returns(uint256)
func (_SortitionPool *SortitionPoolSession) OperatorsInPool() (*big.Int, error) {
	return _SortitionPool.Contract.OperatorsInPool(&_SortitionPool.CallOpts)
}

// OperatorsInPool is a free data retrieval call binding the contract method 0xe7bfd899.
//
// Solidity: function operatorsInPool() view returns(uint256)
func (_SortitionPool *SortitionPoolCallerSession) OperatorsInPool() (*big.Int, error) {
	return _SortitionPool.Contract.OperatorsInPool(&_SortitionPool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SortitionPool *SortitionPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SortitionPool *SortitionPoolSession) Owner() (common.Address, error) {
	return _SortitionPool.Contract.Owner(&_SortitionPool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SortitionPool *SortitionPoolCallerSession) Owner() (common.Address, error) {
	return _SortitionPool.Contract.Owner(&_SortitionPool.CallOpts)
}

// PoolWeightDivisor is a free data retrieval call binding the contract method 0x43a3db30.
//
// Solidity: function poolWeightDivisor() view returns(uint256)
func (_SortitionPool *SortitionPoolCaller) PoolWeightDivisor(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "poolWeightDivisor")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolWeightDivisor is a free data retrieval call binding the contract method 0x43a3db30.
//
// Solidity: function poolWeightDivisor() view returns(uint256)
func (_SortitionPool *SortitionPoolSession) PoolWeightDivisor() (*big.Int, error) {
	return _SortitionPool.Contract.PoolWeightDivisor(&_SortitionPool.CallOpts)
}

// PoolWeightDivisor is a free data retrieval call binding the contract method 0x43a3db30.
//
// Solidity: function poolWeightDivisor() view returns(uint256)
func (_SortitionPool *SortitionPoolCallerSession) PoolWeightDivisor() (*big.Int, error) {
	return _SortitionPool.Contract.PoolWeightDivisor(&_SortitionPool.CallOpts)
}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_SortitionPool *SortitionPoolCaller) RewardToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "rewardToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_SortitionPool *SortitionPoolSession) RewardToken() (common.Address, error) {
	return _SortitionPool.Contract.RewardToken(&_SortitionPool.CallOpts)
}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_SortitionPool *SortitionPoolCallerSession) RewardToken() (common.Address, error) {
	return _SortitionPool.Contract.RewardToken(&_SortitionPool.CallOpts)
}

// RewardsEligibilityRestorableAt is a free data retrieval call binding the contract method 0xbe51fb4a.
//
// Solidity: function rewardsEligibilityRestorableAt(uint32 operator) view returns(uint256)
func (_SortitionPool *SortitionPoolCaller) RewardsEligibilityRestorableAt(opts *bind.CallOpts, operator uint32) (*big.Int, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "rewardsEligibilityRestorableAt", operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardsEligibilityRestorableAt is a free data retrieval call binding the contract method 0xbe51fb4a.
//
// Solidity: function rewardsEligibilityRestorableAt(uint32 operator) view returns(uint256)
func (_SortitionPool *SortitionPoolSession) RewardsEligibilityRestorableAt(operator uint32) (*big.Int, error) {
	return _SortitionPool.Contract.RewardsEligibilityRestorableAt(&_SortitionPool.CallOpts, operator)
}

// RewardsEligibilityRestorableAt is a free data retrieval call binding the contract method 0xbe51fb4a.
//
// Solidity: function rewardsEligibilityRestorableAt(uint32 operator) view returns(uint256)
func (_SortitionPool *SortitionPoolCallerSession) RewardsEligibilityRestorableAt(operator uint32) (*big.Int, error) {
	return _SortitionPool.Contract.RewardsEligibilityRestorableAt(&_SortitionPool.CallOpts, operator)
}

// SelectGroup is a free data retrieval call binding the contract method 0x6c2530b9.
//
// Solidity: function selectGroup(uint256 groupSize, bytes32 seed) view returns(uint32[])
func (_SortitionPool *SortitionPoolCaller) SelectGroup(opts *bind.CallOpts, groupSize *big.Int, seed [32]byte) ([]uint32, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "selectGroup", groupSize, seed)

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

// SelectGroup is a free data retrieval call binding the contract method 0x6c2530b9.
//
// Solidity: function selectGroup(uint256 groupSize, bytes32 seed) view returns(uint32[])
func (_SortitionPool *SortitionPoolSession) SelectGroup(groupSize *big.Int, seed [32]byte) ([]uint32, error) {
	return _SortitionPool.Contract.SelectGroup(&_SortitionPool.CallOpts, groupSize, seed)
}

// SelectGroup is a free data retrieval call binding the contract method 0x6c2530b9.
//
// Solidity: function selectGroup(uint256 groupSize, bytes32 seed) view returns(uint32[])
func (_SortitionPool *SortitionPoolCallerSession) SelectGroup(groupSize *big.Int, seed [32]byte) ([]uint32, error) {
	return _SortitionPool.Contract.SelectGroup(&_SortitionPool.CallOpts, groupSize, seed)
}

// TotalWeight is a free data retrieval call binding the contract method 0x96c82e57.
//
// Solidity: function totalWeight() view returns(uint256)
func (_SortitionPool *SortitionPoolCaller) TotalWeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SortitionPool.contract.Call(opts, &out, "totalWeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalWeight is a free data retrieval call binding the contract method 0x96c82e57.
//
// Solidity: function totalWeight() view returns(uint256)
func (_SortitionPool *SortitionPoolSession) TotalWeight() (*big.Int, error) {
	return _SortitionPool.Contract.TotalWeight(&_SortitionPool.CallOpts)
}

// TotalWeight is a free data retrieval call binding the contract method 0x96c82e57.
//
// Solidity: function totalWeight() view returns(uint256)
func (_SortitionPool *SortitionPoolCallerSession) TotalWeight() (*big.Int, error) {
	return _SortitionPool.Contract.TotalWeight(&_SortitionPool.CallOpts)
}

// InsertOperator is a paid mutator transaction binding the contract method 0x241a4188.
//
// Solidity: function insertOperator(address operator, uint256 authorizedStake) returns()
func (_SortitionPool *SortitionPoolTransactor) InsertOperator(opts *bind.TransactOpts, operator common.Address, authorizedStake *big.Int) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "insertOperator", operator, authorizedStake)
}

// InsertOperator is a paid mutator transaction binding the contract method 0x241a4188.
//
// Solidity: function insertOperator(address operator, uint256 authorizedStake) returns()
func (_SortitionPool *SortitionPoolSession) InsertOperator(operator common.Address, authorizedStake *big.Int) (*types.Transaction, error) {
	return _SortitionPool.Contract.InsertOperator(&_SortitionPool.TransactOpts, operator, authorizedStake)
}

// InsertOperator is a paid mutator transaction binding the contract method 0x241a4188.
//
// Solidity: function insertOperator(address operator, uint256 authorizedStake) returns()
func (_SortitionPool *SortitionPoolTransactorSession) InsertOperator(operator common.Address, authorizedStake *big.Int) (*types.Transaction, error) {
	return _SortitionPool.Contract.InsertOperator(&_SortitionPool.TransactOpts, operator, authorizedStake)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_SortitionPool *SortitionPoolTransactor) Lock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "lock")
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_SortitionPool *SortitionPoolSession) Lock() (*types.Transaction, error) {
	return _SortitionPool.Contract.Lock(&_SortitionPool.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_SortitionPool *SortitionPoolTransactorSession) Lock() (*types.Transaction, error) {
	return _SortitionPool.Contract.Lock(&_SortitionPool.TransactOpts)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address sender, uint256 amount, address token, bytes ) returns()
func (_SortitionPool *SortitionPoolTransactor) ReceiveApproval(opts *bind.TransactOpts, sender common.Address, amount *big.Int, token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "receiveApproval", sender, amount, token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address sender, uint256 amount, address token, bytes ) returns()
func (_SortitionPool *SortitionPoolSession) ReceiveApproval(sender common.Address, amount *big.Int, token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _SortitionPool.Contract.ReceiveApproval(&_SortitionPool.TransactOpts, sender, amount, token, arg3)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address sender, uint256 amount, address token, bytes ) returns()
func (_SortitionPool *SortitionPoolTransactorSession) ReceiveApproval(sender common.Address, amount *big.Int, token common.Address, arg3 []byte) (*types.Transaction, error) {
	return _SortitionPool.Contract.ReceiveApproval(&_SortitionPool.TransactOpts, sender, amount, token, arg3)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SortitionPool *SortitionPoolTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SortitionPool *SortitionPoolSession) RenounceOwnership() (*types.Transaction, error) {
	return _SortitionPool.Contract.RenounceOwnership(&_SortitionPool.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SortitionPool *SortitionPoolTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SortitionPool.Contract.RenounceOwnership(&_SortitionPool.TransactOpts)
}

// RestoreRewardEligibility is a paid mutator transaction binding the contract method 0xb2f3db4d.
//
// Solidity: function restoreRewardEligibility(address operator) returns()
func (_SortitionPool *SortitionPoolTransactor) RestoreRewardEligibility(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "restoreRewardEligibility", operator)
}

// RestoreRewardEligibility is a paid mutator transaction binding the contract method 0xb2f3db4d.
//
// Solidity: function restoreRewardEligibility(address operator) returns()
func (_SortitionPool *SortitionPoolSession) RestoreRewardEligibility(operator common.Address) (*types.Transaction, error) {
	return _SortitionPool.Contract.RestoreRewardEligibility(&_SortitionPool.TransactOpts, operator)
}

// RestoreRewardEligibility is a paid mutator transaction binding the contract method 0xb2f3db4d.
//
// Solidity: function restoreRewardEligibility(address operator) returns()
func (_SortitionPool *SortitionPoolTransactorSession) RestoreRewardEligibility(operator common.Address) (*types.Transaction, error) {
	return _SortitionPool.Contract.RestoreRewardEligibility(&_SortitionPool.TransactOpts, operator)
}

// SetRewardIneligibility is a paid mutator transaction binding the contract method 0x942f6892.
//
// Solidity: function setRewardIneligibility(uint32[] operators, uint256 until) returns()
func (_SortitionPool *SortitionPoolTransactor) SetRewardIneligibility(opts *bind.TransactOpts, operators []uint32, until *big.Int) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "setRewardIneligibility", operators, until)
}

// SetRewardIneligibility is a paid mutator transaction binding the contract method 0x942f6892.
//
// Solidity: function setRewardIneligibility(uint32[] operators, uint256 until) returns()
func (_SortitionPool *SortitionPoolSession) SetRewardIneligibility(operators []uint32, until *big.Int) (*types.Transaction, error) {
	return _SortitionPool.Contract.SetRewardIneligibility(&_SortitionPool.TransactOpts, operators, until)
}

// SetRewardIneligibility is a paid mutator transaction binding the contract method 0x942f6892.
//
// Solidity: function setRewardIneligibility(uint32[] operators, uint256 until) returns()
func (_SortitionPool *SortitionPoolTransactorSession) SetRewardIneligibility(operators []uint32, until *big.Int) (*types.Transaction, error) {
	return _SortitionPool.Contract.SetRewardIneligibility(&_SortitionPool.TransactOpts, operators, until)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SortitionPool *SortitionPoolTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SortitionPool *SortitionPoolSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SortitionPool.Contract.TransferOwnership(&_SortitionPool.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SortitionPool *SortitionPoolTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SortitionPool.Contract.TransferOwnership(&_SortitionPool.TransactOpts, newOwner)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_SortitionPool *SortitionPoolTransactor) Unlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "unlock")
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_SortitionPool *SortitionPoolSession) Unlock() (*types.Transaction, error) {
	return _SortitionPool.Contract.Unlock(&_SortitionPool.TransactOpts)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_SortitionPool *SortitionPoolTransactorSession) Unlock() (*types.Transaction, error) {
	return _SortitionPool.Contract.Unlock(&_SortitionPool.TransactOpts)
}

// UpdateOperatorStatus is a paid mutator transaction binding the contract method 0xdc7520c5.
//
// Solidity: function updateOperatorStatus(address operator, uint256 authorizedStake) returns()
func (_SortitionPool *SortitionPoolTransactor) UpdateOperatorStatus(opts *bind.TransactOpts, operator common.Address, authorizedStake *big.Int) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "updateOperatorStatus", operator, authorizedStake)
}

// UpdateOperatorStatus is a paid mutator transaction binding the contract method 0xdc7520c5.
//
// Solidity: function updateOperatorStatus(address operator, uint256 authorizedStake) returns()
func (_SortitionPool *SortitionPoolSession) UpdateOperatorStatus(operator common.Address, authorizedStake *big.Int) (*types.Transaction, error) {
	return _SortitionPool.Contract.UpdateOperatorStatus(&_SortitionPool.TransactOpts, operator, authorizedStake)
}

// UpdateOperatorStatus is a paid mutator transaction binding the contract method 0xdc7520c5.
//
// Solidity: function updateOperatorStatus(address operator, uint256 authorizedStake) returns()
func (_SortitionPool *SortitionPoolTransactorSession) UpdateOperatorStatus(operator common.Address, authorizedStake *big.Int) (*types.Transaction, error) {
	return _SortitionPool.Contract.UpdateOperatorStatus(&_SortitionPool.TransactOpts, operator, authorizedStake)
}

// WithdrawIneligible is a paid mutator transaction binding the contract method 0xa9649414.
//
// Solidity: function withdrawIneligible(address recipient) returns()
func (_SortitionPool *SortitionPoolTransactor) WithdrawIneligible(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "withdrawIneligible", recipient)
}

// WithdrawIneligible is a paid mutator transaction binding the contract method 0xa9649414.
//
// Solidity: function withdrawIneligible(address recipient) returns()
func (_SortitionPool *SortitionPoolSession) WithdrawIneligible(recipient common.Address) (*types.Transaction, error) {
	return _SortitionPool.Contract.WithdrawIneligible(&_SortitionPool.TransactOpts, recipient)
}

// WithdrawIneligible is a paid mutator transaction binding the contract method 0xa9649414.
//
// Solidity: function withdrawIneligible(address recipient) returns()
func (_SortitionPool *SortitionPoolTransactorSession) WithdrawIneligible(recipient common.Address) (*types.Transaction, error) {
	return _SortitionPool.Contract.WithdrawIneligible(&_SortitionPool.TransactOpts, recipient)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xe20981ca.
//
// Solidity: function withdrawRewards(address operator, address beneficiary) returns(uint96)
func (_SortitionPool *SortitionPoolTransactor) WithdrawRewards(opts *bind.TransactOpts, operator common.Address, beneficiary common.Address) (*types.Transaction, error) {
	return _SortitionPool.contract.Transact(opts, "withdrawRewards", operator, beneficiary)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xe20981ca.
//
// Solidity: function withdrawRewards(address operator, address beneficiary) returns(uint96)
func (_SortitionPool *SortitionPoolSession) WithdrawRewards(operator common.Address, beneficiary common.Address) (*types.Transaction, error) {
	return _SortitionPool.Contract.WithdrawRewards(&_SortitionPool.TransactOpts, operator, beneficiary)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xe20981ca.
//
// Solidity: function withdrawRewards(address operator, address beneficiary) returns(uint96)
func (_SortitionPool *SortitionPoolTransactorSession) WithdrawRewards(operator common.Address, beneficiary common.Address) (*types.Transaction, error) {
	return _SortitionPool.Contract.WithdrawRewards(&_SortitionPool.TransactOpts, operator, beneficiary)
}

// SortitionPoolIneligibleForRewardsIterator is returned from FilterIneligibleForRewards and is used to iterate over the raw logs and unpacked data for IneligibleForRewards events raised by the SortitionPool contract.
type SortitionPoolIneligibleForRewardsIterator struct {
	Event *SortitionPoolIneligibleForRewards // Event containing the contract specifics and raw log

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
func (it *SortitionPoolIneligibleForRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortitionPoolIneligibleForRewards)
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
		it.Event = new(SortitionPoolIneligibleForRewards)
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
func (it *SortitionPoolIneligibleForRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortitionPoolIneligibleForRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortitionPoolIneligibleForRewards represents a IneligibleForRewards event raised by the SortitionPool contract.
type SortitionPoolIneligibleForRewards struct {
	Ids   []uint32
	Until *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterIneligibleForRewards is a free log retrieval operation binding the contract event 0x01f5838e3dde8cf4817b958fe95be92bdfeccb34317e1d9f58d1cfe5230de231.
//
// Solidity: event IneligibleForRewards(uint32[] ids, uint256 until)
func (_SortitionPool *SortitionPoolFilterer) FilterIneligibleForRewards(opts *bind.FilterOpts) (*SortitionPoolIneligibleForRewardsIterator, error) {

	logs, sub, err := _SortitionPool.contract.FilterLogs(opts, "IneligibleForRewards")
	if err != nil {
		return nil, err
	}
	return &SortitionPoolIneligibleForRewardsIterator{contract: _SortitionPool.contract, event: "IneligibleForRewards", logs: logs, sub: sub}, nil
}

// WatchIneligibleForRewards is a free log subscription operation binding the contract event 0x01f5838e3dde8cf4817b958fe95be92bdfeccb34317e1d9f58d1cfe5230de231.
//
// Solidity: event IneligibleForRewards(uint32[] ids, uint256 until)
func (_SortitionPool *SortitionPoolFilterer) WatchIneligibleForRewards(opts *bind.WatchOpts, sink chan<- *SortitionPoolIneligibleForRewards) (event.Subscription, error) {

	logs, sub, err := _SortitionPool.contract.WatchLogs(opts, "IneligibleForRewards")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortitionPoolIneligibleForRewards)
				if err := _SortitionPool.contract.UnpackLog(event, "IneligibleForRewards", log); err != nil {
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

// ParseIneligibleForRewards is a log parse operation binding the contract event 0x01f5838e3dde8cf4817b958fe95be92bdfeccb34317e1d9f58d1cfe5230de231.
//
// Solidity: event IneligibleForRewards(uint32[] ids, uint256 until)
func (_SortitionPool *SortitionPoolFilterer) ParseIneligibleForRewards(log types.Log) (*SortitionPoolIneligibleForRewards, error) {
	event := new(SortitionPoolIneligibleForRewards)
	if err := _SortitionPool.contract.UnpackLog(event, "IneligibleForRewards", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortitionPoolOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SortitionPool contract.
type SortitionPoolOwnershipTransferredIterator struct {
	Event *SortitionPoolOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SortitionPoolOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortitionPoolOwnershipTransferred)
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
		it.Event = new(SortitionPoolOwnershipTransferred)
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
func (it *SortitionPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortitionPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortitionPoolOwnershipTransferred represents a OwnershipTransferred event raised by the SortitionPool contract.
type SortitionPoolOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SortitionPool *SortitionPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SortitionPoolOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SortitionPool.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SortitionPoolOwnershipTransferredIterator{contract: _SortitionPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SortitionPool *SortitionPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SortitionPoolOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SortitionPool.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortitionPoolOwnershipTransferred)
				if err := _SortitionPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SortitionPool *SortitionPoolFilterer) ParseOwnershipTransferred(log types.Log) (*SortitionPoolOwnershipTransferred, error) {
	event := new(SortitionPoolOwnershipTransferred)
	if err := _SortitionPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortitionPoolRewardEligibilityRestoredIterator is returned from FilterRewardEligibilityRestored and is used to iterate over the raw logs and unpacked data for RewardEligibilityRestored events raised by the SortitionPool contract.
type SortitionPoolRewardEligibilityRestoredIterator struct {
	Event *SortitionPoolRewardEligibilityRestored // Event containing the contract specifics and raw log

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
func (it *SortitionPoolRewardEligibilityRestoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortitionPoolRewardEligibilityRestored)
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
		it.Event = new(SortitionPoolRewardEligibilityRestored)
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
func (it *SortitionPoolRewardEligibilityRestoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortitionPoolRewardEligibilityRestoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortitionPoolRewardEligibilityRestored represents a RewardEligibilityRestored event raised by the SortitionPool contract.
type SortitionPoolRewardEligibilityRestored struct {
	Operator common.Address
	Id       uint32
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRewardEligibilityRestored is a free log retrieval operation binding the contract event 0xe61e9f0f049b3bfae1ae903a5e3018c02a008aa0d238ffddf23a4fb4c0278536.
//
// Solidity: event RewardEligibilityRestored(address indexed operator, uint32 indexed id)
func (_SortitionPool *SortitionPoolFilterer) FilterRewardEligibilityRestored(opts *bind.FilterOpts, operator []common.Address, id []uint32) (*SortitionPoolRewardEligibilityRestoredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _SortitionPool.contract.FilterLogs(opts, "RewardEligibilityRestored", operatorRule, idRule)
	if err != nil {
		return nil, err
	}
	return &SortitionPoolRewardEligibilityRestoredIterator{contract: _SortitionPool.contract, event: "RewardEligibilityRestored", logs: logs, sub: sub}, nil
}

// WatchRewardEligibilityRestored is a free log subscription operation binding the contract event 0xe61e9f0f049b3bfae1ae903a5e3018c02a008aa0d238ffddf23a4fb4c0278536.
//
// Solidity: event RewardEligibilityRestored(address indexed operator, uint32 indexed id)
func (_SortitionPool *SortitionPoolFilterer) WatchRewardEligibilityRestored(opts *bind.WatchOpts, sink chan<- *SortitionPoolRewardEligibilityRestored, operator []common.Address, id []uint32) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _SortitionPool.contract.WatchLogs(opts, "RewardEligibilityRestored", operatorRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortitionPoolRewardEligibilityRestored)
				if err := _SortitionPool.contract.UnpackLog(event, "RewardEligibilityRestored", log); err != nil {
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

// ParseRewardEligibilityRestored is a log parse operation binding the contract event 0xe61e9f0f049b3bfae1ae903a5e3018c02a008aa0d238ffddf23a4fb4c0278536.
//
// Solidity: event RewardEligibilityRestored(address indexed operator, uint32 indexed id)
func (_SortitionPool *SortitionPoolFilterer) ParseRewardEligibilityRestored(log types.Log) (*SortitionPoolRewardEligibilityRestored, error) {
	event := new(SortitionPoolRewardEligibilityRestored)
	if err := _SortitionPool.contract.UnpackLog(event, "RewardEligibilityRestored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
