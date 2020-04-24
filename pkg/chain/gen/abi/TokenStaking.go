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

// TokenStakingABI is the input ABI used to generate the binding from.
const TokenStakingABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_registry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_initializationPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_undelegationPeriod\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"lockCreator\",\"type\":\"address\"}],\"name\":\"ExpiredLockReleased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"lockCreator\",\"type\":\"address\"}],\"name\":\"LockReleased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"recoveredAt\",\"type\":\"uint256\"}],\"name\":\"RecoveredStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"lockCreator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"until\",\"type\":\"uint256\"}],\"name\":\"StakeLocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokensSeized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokensSlashed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"undelegatedAt\",\"type\":\"uint256\"}],\"name\":\"Undelegated\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"activeStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"authorizeOperatorContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"authorizerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"cancelStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatedAuthoritySource\",\"type\":\"address\"}],\"name\":\"claimDelegatedAuthority\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"eligibleStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"getAuthoritySource\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"getDelegationInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"undelegatedAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getLocks\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"creators\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"expirations\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"hasMinimumStake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"initializationPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"isAuthorizedForOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isStakeLocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"lockStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"beneficiaryOf\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maximumLockDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumStakeBase\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumStakeSchedule\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumStakeScheduleStart\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumStakeSteps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"packedParams\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"authorizer\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"operatorsOf\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ownerOperators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"recoverStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractRegistry\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"releaseExpiredLock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToSeize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewardMultiplier\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tattletale\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"misbehavedOperators\",\"type\":\"address[]\"}],\"name\":\"seize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToSlash\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"misbehavedOperators\",\"type\":\"address[]\"}],\"name\":\"slash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractERC20Burnable\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"undelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_undelegationTimestamp\",\"type\":\"uint256\"}],\"name\":\"undelegateAt\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"undelegationPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"unlockStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenStaking is an auto generated Go binding around an Ethereum contract.
type TokenStaking struct {
	TokenStakingCaller     // Read-only binding to the contract
	TokenStakingTransactor // Write-only binding to the contract
	TokenStakingFilterer   // Log filterer for contract events
}

// TokenStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenStakingSession struct {
	Contract     *TokenStaking     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenStakingCallerSession struct {
	Contract *TokenStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TokenStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenStakingTransactorSession struct {
	Contract     *TokenStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TokenStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenStakingRaw struct {
	Contract *TokenStaking // Generic contract binding to access the raw methods on
}

// TokenStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenStakingCallerRaw struct {
	Contract *TokenStakingCaller // Generic read-only contract binding to access the raw methods on
}

// TokenStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenStakingTransactorRaw struct {
	Contract *TokenStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenStaking creates a new instance of TokenStaking, bound to a specific deployed contract.
func NewTokenStaking(address common.Address, backend bind.ContractBackend) (*TokenStaking, error) {
	contract, err := bindTokenStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenStaking{TokenStakingCaller: TokenStakingCaller{contract: contract}, TokenStakingTransactor: TokenStakingTransactor{contract: contract}, TokenStakingFilterer: TokenStakingFilterer{contract: contract}}, nil
}

// NewTokenStakingCaller creates a new read-only instance of TokenStaking, bound to a specific deployed contract.
func NewTokenStakingCaller(address common.Address, caller bind.ContractCaller) (*TokenStakingCaller, error) {
	contract, err := bindTokenStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenStakingCaller{contract: contract}, nil
}

// NewTokenStakingTransactor creates a new write-only instance of TokenStaking, bound to a specific deployed contract.
func NewTokenStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenStakingTransactor, error) {
	contract, err := bindTokenStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenStakingTransactor{contract: contract}, nil
}

// NewTokenStakingFilterer creates a new log filterer instance of TokenStaking, bound to a specific deployed contract.
func NewTokenStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenStakingFilterer, error) {
	contract, err := bindTokenStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenStakingFilterer{contract: contract}, nil
}

// bindTokenStaking binds a generic wrapper to an already deployed contract.
func bindTokenStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenStakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenStaking *TokenStakingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenStaking.Contract.TokenStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenStaking *TokenStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenStaking.Contract.TokenStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenStaking *TokenStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenStaking.Contract.TokenStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenStaking *TokenStakingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenStaking *TokenStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenStaking *TokenStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenStaking.Contract.contract.Transact(opts, method, params...)
}

// ActiveStake is a free data retrieval call binding the contract method 0x9557e0bb.
//
// Solidity: function activeStake(address _operator, address _operatorContract) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingCaller) ActiveStake(opts *bind.CallOpts, _operator common.Address, _operatorContract common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "activeStake", _operator, _operatorContract)
	return *ret0, err
}

// ActiveStake is a free data retrieval call binding the contract method 0x9557e0bb.
//
// Solidity: function activeStake(address _operator, address _operatorContract) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingSession) ActiveStake(_operator common.Address, _operatorContract common.Address) (*big.Int, error) {
	return _TokenStaking.Contract.ActiveStake(&_TokenStaking.CallOpts, _operator, _operatorContract)
}

// ActiveStake is a free data retrieval call binding the contract method 0x9557e0bb.
//
// Solidity: function activeStake(address _operator, address _operatorContract) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingCallerSession) ActiveStake(_operator common.Address, _operatorContract common.Address) (*big.Int, error) {
	return _TokenStaking.Contract.ActiveStake(&_TokenStaking.CallOpts, _operator, _operatorContract)
}

// AuthorizerOf is a free data retrieval call binding the contract method 0xfb1677b1.
//
// Solidity: function authorizerOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingCaller) AuthorizerOf(opts *bind.CallOpts, _operator common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "authorizerOf", _operator)
	return *ret0, err
}

// AuthorizerOf is a free data retrieval call binding the contract method 0xfb1677b1.
//
// Solidity: function authorizerOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingSession) AuthorizerOf(_operator common.Address) (common.Address, error) {
	return _TokenStaking.Contract.AuthorizerOf(&_TokenStaking.CallOpts, _operator)
}

// AuthorizerOf is a free data retrieval call binding the contract method 0xfb1677b1.
//
// Solidity: function authorizerOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingCallerSession) AuthorizerOf(_operator common.Address) (common.Address, error) {
	return _TokenStaking.Contract.AuthorizerOf(&_TokenStaking.CallOpts, _operator)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _address) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingCaller) BalanceOf(opts *bind.CallOpts, _address common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "balanceOf", _address)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _address) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingSession) BalanceOf(_address common.Address) (*big.Int, error) {
	return _TokenStaking.Contract.BalanceOf(&_TokenStaking.CallOpts, _address)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _address) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingCallerSession) BalanceOf(_address common.Address) (*big.Int, error) {
	return _TokenStaking.Contract.BalanceOf(&_TokenStaking.CallOpts, _address)
}

// EligibleStake is a free data retrieval call binding the contract method 0xafff33ef.
//
// Solidity: function eligibleStake(address _operator, address _operatorContract) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingCaller) EligibleStake(opts *bind.CallOpts, _operator common.Address, _operatorContract common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "eligibleStake", _operator, _operatorContract)
	return *ret0, err
}

// EligibleStake is a free data retrieval call binding the contract method 0xafff33ef.
//
// Solidity: function eligibleStake(address _operator, address _operatorContract) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingSession) EligibleStake(_operator common.Address, _operatorContract common.Address) (*big.Int, error) {
	return _TokenStaking.Contract.EligibleStake(&_TokenStaking.CallOpts, _operator, _operatorContract)
}

// EligibleStake is a free data retrieval call binding the contract method 0xafff33ef.
//
// Solidity: function eligibleStake(address _operator, address _operatorContract) constant returns(uint256 balance)
func (_TokenStaking *TokenStakingCallerSession) EligibleStake(_operator common.Address, _operatorContract common.Address) (*big.Int, error) {
	return _TokenStaking.Contract.EligibleStake(&_TokenStaking.CallOpts, _operator, _operatorContract)
}

// GetAuthoritySource is a free data retrieval call binding the contract method 0xcbe945dc.
//
// Solidity: function getAuthoritySource(address operatorContract) constant returns(address)
func (_TokenStaking *TokenStakingCaller) GetAuthoritySource(opts *bind.CallOpts, operatorContract common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "getAuthoritySource", operatorContract)
	return *ret0, err
}

// GetAuthoritySource is a free data retrieval call binding the contract method 0xcbe945dc.
//
// Solidity: function getAuthoritySource(address operatorContract) constant returns(address)
func (_TokenStaking *TokenStakingSession) GetAuthoritySource(operatorContract common.Address) (common.Address, error) {
	return _TokenStaking.Contract.GetAuthoritySource(&_TokenStaking.CallOpts, operatorContract)
}

// GetAuthoritySource is a free data retrieval call binding the contract method 0xcbe945dc.
//
// Solidity: function getAuthoritySource(address operatorContract) constant returns(address)
func (_TokenStaking *TokenStakingCallerSession) GetAuthoritySource(operatorContract common.Address) (common.Address, error) {
	return _TokenStaking.Contract.GetAuthoritySource(&_TokenStaking.CallOpts, operatorContract)
}

// GetDelegationInfo is a free data retrieval call binding the contract method 0xfab46d66.
//
// Solidity: function getDelegationInfo(address _operator) constant returns(uint256 amount, uint256 createdAt, uint256 undelegatedAt)
func (_TokenStaking *TokenStakingCaller) GetDelegationInfo(opts *bind.CallOpts, _operator common.Address) (struct {
	Amount        *big.Int
	CreatedAt     *big.Int
	UndelegatedAt *big.Int
}, error) {
	ret := new(struct {
		Amount        *big.Int
		CreatedAt     *big.Int
		UndelegatedAt *big.Int
	})
	out := ret
	err := _TokenStaking.contract.Call(opts, out, "getDelegationInfo", _operator)
	return *ret, err
}

// GetDelegationInfo is a free data retrieval call binding the contract method 0xfab46d66.
//
// Solidity: function getDelegationInfo(address _operator) constant returns(uint256 amount, uint256 createdAt, uint256 undelegatedAt)
func (_TokenStaking *TokenStakingSession) GetDelegationInfo(_operator common.Address) (struct {
	Amount        *big.Int
	CreatedAt     *big.Int
	UndelegatedAt *big.Int
}, error) {
	return _TokenStaking.Contract.GetDelegationInfo(&_TokenStaking.CallOpts, _operator)
}

// GetDelegationInfo is a free data retrieval call binding the contract method 0xfab46d66.
//
// Solidity: function getDelegationInfo(address _operator) constant returns(uint256 amount, uint256 createdAt, uint256 undelegatedAt)
func (_TokenStaking *TokenStakingCallerSession) GetDelegationInfo(_operator common.Address) (struct {
	Amount        *big.Int
	CreatedAt     *big.Int
	UndelegatedAt *big.Int
}, error) {
	return _TokenStaking.Contract.GetDelegationInfo(&_TokenStaking.CallOpts, _operator)
}

// GetLocks is a free data retrieval call binding the contract method 0x719f3089.
//
// Solidity: function getLocks(address operator) constant returns(address[] creators, uint256[] expirations)
func (_TokenStaking *TokenStakingCaller) GetLocks(opts *bind.CallOpts, operator common.Address) (struct {
	Creators    []common.Address
	Expirations []*big.Int
}, error) {
	ret := new(struct {
		Creators    []common.Address
		Expirations []*big.Int
	})
	out := ret
	err := _TokenStaking.contract.Call(opts, out, "getLocks", operator)
	return *ret, err
}

// GetLocks is a free data retrieval call binding the contract method 0x719f3089.
//
// Solidity: function getLocks(address operator) constant returns(address[] creators, uint256[] expirations)
func (_TokenStaking *TokenStakingSession) GetLocks(operator common.Address) (struct {
	Creators    []common.Address
	Expirations []*big.Int
}, error) {
	return _TokenStaking.Contract.GetLocks(&_TokenStaking.CallOpts, operator)
}

// GetLocks is a free data retrieval call binding the contract method 0x719f3089.
//
// Solidity: function getLocks(address operator) constant returns(address[] creators, uint256[] expirations)
func (_TokenStaking *TokenStakingCallerSession) GetLocks(operator common.Address) (struct {
	Creators    []common.Address
	Expirations []*big.Int
}, error) {
	return _TokenStaking.Contract.GetLocks(&_TokenStaking.CallOpts, operator)
}

// HasMinimumStake is a free data retrieval call binding the contract method 0x10a63ec0.
//
// Solidity: function hasMinimumStake(address staker, address operatorContract) constant returns(bool)
func (_TokenStaking *TokenStakingCaller) HasMinimumStake(opts *bind.CallOpts, staker common.Address, operatorContract common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "hasMinimumStake", staker, operatorContract)
	return *ret0, err
}

// HasMinimumStake is a free data retrieval call binding the contract method 0x10a63ec0.
//
// Solidity: function hasMinimumStake(address staker, address operatorContract) constant returns(bool)
func (_TokenStaking *TokenStakingSession) HasMinimumStake(staker common.Address, operatorContract common.Address) (bool, error) {
	return _TokenStaking.Contract.HasMinimumStake(&_TokenStaking.CallOpts, staker, operatorContract)
}

// HasMinimumStake is a free data retrieval call binding the contract method 0x10a63ec0.
//
// Solidity: function hasMinimumStake(address staker, address operatorContract) constant returns(bool)
func (_TokenStaking *TokenStakingCallerSession) HasMinimumStake(staker common.Address, operatorContract common.Address) (bool, error) {
	return _TokenStaking.Contract.HasMinimumStake(&_TokenStaking.CallOpts, staker, operatorContract)
}

// InitializationPeriod is a free data retrieval call binding the contract method 0xaed1ec72.
//
// Solidity: function initializationPeriod() constant returns(uint256)
func (_TokenStaking *TokenStakingCaller) InitializationPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "initializationPeriod")
	return *ret0, err
}

// InitializationPeriod is a free data retrieval call binding the contract method 0xaed1ec72.
//
// Solidity: function initializationPeriod() constant returns(uint256)
func (_TokenStaking *TokenStakingSession) InitializationPeriod() (*big.Int, error) {
	return _TokenStaking.Contract.InitializationPeriod(&_TokenStaking.CallOpts)
}

// InitializationPeriod is a free data retrieval call binding the contract method 0xaed1ec72.
//
// Solidity: function initializationPeriod() constant returns(uint256)
func (_TokenStaking *TokenStakingCallerSession) InitializationPeriod() (*big.Int, error) {
	return _TokenStaking.Contract.InitializationPeriod(&_TokenStaking.CallOpts)
}

// IsAuthorizedForOperator is a free data retrieval call binding the contract method 0xef1f9661.
//
// Solidity: function isAuthorizedForOperator(address _operator, address _operatorContract) constant returns(bool)
func (_TokenStaking *TokenStakingCaller) IsAuthorizedForOperator(opts *bind.CallOpts, _operator common.Address, _operatorContract common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "isAuthorizedForOperator", _operator, _operatorContract)
	return *ret0, err
}

// IsAuthorizedForOperator is a free data retrieval call binding the contract method 0xef1f9661.
//
// Solidity: function isAuthorizedForOperator(address _operator, address _operatorContract) constant returns(bool)
func (_TokenStaking *TokenStakingSession) IsAuthorizedForOperator(_operator common.Address, _operatorContract common.Address) (bool, error) {
	return _TokenStaking.Contract.IsAuthorizedForOperator(&_TokenStaking.CallOpts, _operator, _operatorContract)
}

// IsAuthorizedForOperator is a free data retrieval call binding the contract method 0xef1f9661.
//
// Solidity: function isAuthorizedForOperator(address _operator, address _operatorContract) constant returns(bool)
func (_TokenStaking *TokenStakingCallerSession) IsAuthorizedForOperator(_operator common.Address, _operatorContract common.Address) (bool, error) {
	return _TokenStaking.Contract.IsAuthorizedForOperator(&_TokenStaking.CallOpts, _operator, _operatorContract)
}

// IsStakeLocked is a free data retrieval call binding the contract method 0x335e91a1.
//
// Solidity: function isStakeLocked(address operator) constant returns(bool)
func (_TokenStaking *TokenStakingCaller) IsStakeLocked(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "isStakeLocked", operator)
	return *ret0, err
}

// IsStakeLocked is a free data retrieval call binding the contract method 0x335e91a1.
//
// Solidity: function isStakeLocked(address operator) constant returns(bool)
func (_TokenStaking *TokenStakingSession) IsStakeLocked(operator common.Address) (bool, error) {
	return _TokenStaking.Contract.IsStakeLocked(&_TokenStaking.CallOpts, operator)
}

// IsStakeLocked is a free data retrieval call binding the contract method 0x335e91a1.
//
// Solidity: function isStakeLocked(address operator) constant returns(bool)
func (_TokenStaking *TokenStakingCallerSession) IsStakeLocked(operator common.Address) (bool, error) {
	return _TokenStaking.Contract.IsStakeLocked(&_TokenStaking.CallOpts, operator)
}

// BeneficiaryOf is a free data retrieval call binding the contract method 0x1cdac873.
//
// Solidity: function beneficiaryOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingCaller) BeneficiaryOf(opts *bind.CallOpts, _operator common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "beneficiaryOf", _operator)
	return *ret0, err
}

// BeneficiaryOf is a free data retrieval call binding the contract method 0x1cdac873.
//
// Solidity: function beneficiaryOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingSession) BeneficiaryOf(_operator common.Address) (common.Address, error) {
	return _TokenStaking.Contract.BeneficiaryOf(&_TokenStaking.CallOpts, _operator)
}

// BeneficiaryOf is a free data retrieval call binding the contract method 0x1cdac873.
//
// Solidity: function beneficiaryOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingCallerSession) BeneficiaryOf(_operator common.Address) (common.Address, error) {
	return _TokenStaking.Contract.BeneficiaryOf(&_TokenStaking.CallOpts, _operator)
}

// MaximumLockDuration is a free data retrieval call binding the contract method 0x9ff3f125.
//
// Solidity: function maximumLockDuration() constant returns(uint256)
func (_TokenStaking *TokenStakingCaller) MaximumLockDuration(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "maximumLockDuration")
	return *ret0, err
}

// MaximumLockDuration is a free data retrieval call binding the contract method 0x9ff3f125.
//
// Solidity: function maximumLockDuration() constant returns(uint256)
func (_TokenStaking *TokenStakingSession) MaximumLockDuration() (*big.Int, error) {
	return _TokenStaking.Contract.MaximumLockDuration(&_TokenStaking.CallOpts)
}

// MaximumLockDuration is a free data retrieval call binding the contract method 0x9ff3f125.
//
// Solidity: function maximumLockDuration() constant returns(uint256)
func (_TokenStaking *TokenStakingCallerSession) MaximumLockDuration() (*big.Int, error) {
	return _TokenStaking.Contract.MaximumLockDuration(&_TokenStaking.CallOpts)
}

// MinimumStake is a free data retrieval call binding the contract method 0xec5ffac2.
//
// Solidity: function minimumStake() constant returns(uint256)
func (_TokenStaking *TokenStakingCaller) MinimumStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "minimumStake")
	return *ret0, err
}

// MinimumStake is a free data retrieval call binding the contract method 0xec5ffac2.
//
// Solidity: function minimumStake() constant returns(uint256)
func (_TokenStaking *TokenStakingSession) MinimumStake() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStake(&_TokenStaking.CallOpts)
}

// MinimumStake is a free data retrieval call binding the contract method 0xec5ffac2.
//
// Solidity: function minimumStake() constant returns(uint256)
func (_TokenStaking *TokenStakingCallerSession) MinimumStake() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStake(&_TokenStaking.CallOpts)
}

// MinimumStakeBase is a free data retrieval call binding the contract method 0x89002fed.
//
// Solidity: function minimumStakeBase() constant returns(uint256)
func (_TokenStaking *TokenStakingCaller) MinimumStakeBase(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "minimumStakeBase")
	return *ret0, err
}

// MinimumStakeBase is a free data retrieval call binding the contract method 0x89002fed.
//
// Solidity: function minimumStakeBase() constant returns(uint256)
func (_TokenStaking *TokenStakingSession) MinimumStakeBase() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStakeBase(&_TokenStaking.CallOpts)
}

// MinimumStakeBase is a free data retrieval call binding the contract method 0x89002fed.
//
// Solidity: function minimumStakeBase() constant returns(uint256)
func (_TokenStaking *TokenStakingCallerSession) MinimumStakeBase() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStakeBase(&_TokenStaking.CallOpts)
}

// MinimumStakeSchedule is a free data retrieval call binding the contract method 0x1a051082.
//
// Solidity: function minimumStakeSchedule() constant returns(uint256)
func (_TokenStaking *TokenStakingCaller) MinimumStakeSchedule(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "minimumStakeSchedule")
	return *ret0, err
}

// MinimumStakeSchedule is a free data retrieval call binding the contract method 0x1a051082.
//
// Solidity: function minimumStakeSchedule() constant returns(uint256)
func (_TokenStaking *TokenStakingSession) MinimumStakeSchedule() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStakeSchedule(&_TokenStaking.CallOpts)
}

// MinimumStakeSchedule is a free data retrieval call binding the contract method 0x1a051082.
//
// Solidity: function minimumStakeSchedule() constant returns(uint256)
func (_TokenStaking *TokenStakingCallerSession) MinimumStakeSchedule() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStakeSchedule(&_TokenStaking.CallOpts)
}

// MinimumStakeScheduleStart is a free data retrieval call binding the contract method 0xcdc3e90a.
//
// Solidity: function minimumStakeScheduleStart() constant returns(uint256)
func (_TokenStaking *TokenStakingCaller) MinimumStakeScheduleStart(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "minimumStakeScheduleStart")
	return *ret0, err
}

// MinimumStakeScheduleStart is a free data retrieval call binding the contract method 0xcdc3e90a.
//
// Solidity: function minimumStakeScheduleStart() constant returns(uint256)
func (_TokenStaking *TokenStakingSession) MinimumStakeScheduleStart() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStakeScheduleStart(&_TokenStaking.CallOpts)
}

// MinimumStakeScheduleStart is a free data retrieval call binding the contract method 0xcdc3e90a.
//
// Solidity: function minimumStakeScheduleStart() constant returns(uint256)
func (_TokenStaking *TokenStakingCallerSession) MinimumStakeScheduleStart() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStakeScheduleStart(&_TokenStaking.CallOpts)
}

// MinimumStakeSteps is a free data retrieval call binding the contract method 0x280c3846.
//
// Solidity: function minimumStakeSteps() constant returns(uint256)
func (_TokenStaking *TokenStakingCaller) MinimumStakeSteps(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "minimumStakeSteps")
	return *ret0, err
}

// MinimumStakeSteps is a free data retrieval call binding the contract method 0x280c3846.
//
// Solidity: function minimumStakeSteps() constant returns(uint256)
func (_TokenStaking *TokenStakingSession) MinimumStakeSteps() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStakeSteps(&_TokenStaking.CallOpts)
}

// MinimumStakeSteps is a free data retrieval call binding the contract method 0x280c3846.
//
// Solidity: function minimumStakeSteps() constant returns(uint256)
func (_TokenStaking *TokenStakingCallerSession) MinimumStakeSteps() (*big.Int, error) {
	return _TokenStaking.Contract.MinimumStakeSteps(&_TokenStaking.CallOpts)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) constant returns(uint256 packedParams, address owner, address beneficiary, address authorizer)
func (_TokenStaking *TokenStakingCaller) Operators(opts *bind.CallOpts, arg0 common.Address) (struct {
	PackedParams *big.Int
	Owner        common.Address
	Beneficiary  common.Address
	Authorizer   common.Address
}, error) {
	ret := new(struct {
		PackedParams *big.Int
		Owner        common.Address
		Beneficiary  common.Address
		Authorizer   common.Address
	})
	out := ret
	err := _TokenStaking.contract.Call(opts, out, "operators", arg0)
	return *ret, err
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) constant returns(uint256 packedParams, address owner, address beneficiary, address authorizer)
func (_TokenStaking *TokenStakingSession) Operators(arg0 common.Address) (struct {
	PackedParams *big.Int
	Owner        common.Address
	Beneficiary  common.Address
	Authorizer   common.Address
}, error) {
	return _TokenStaking.Contract.Operators(&_TokenStaking.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) constant returns(uint256 packedParams, address owner, address beneficiary, address authorizer)
func (_TokenStaking *TokenStakingCallerSession) Operators(arg0 common.Address) (struct {
	PackedParams *big.Int
	Owner        common.Address
	Beneficiary  common.Address
	Authorizer   common.Address
}, error) {
	return _TokenStaking.Contract.Operators(&_TokenStaking.CallOpts, arg0)
}

// OperatorsOf is a free data retrieval call binding the contract method 0xb534fbb6.
//
// Solidity: function operatorsOf(address _address) constant returns(address[])
func (_TokenStaking *TokenStakingCaller) OperatorsOf(opts *bind.CallOpts, _address common.Address) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "operatorsOf", _address)
	return *ret0, err
}

// OperatorsOf is a free data retrieval call binding the contract method 0xb534fbb6.
//
// Solidity: function operatorsOf(address _address) constant returns(address[])
func (_TokenStaking *TokenStakingSession) OperatorsOf(_address common.Address) ([]common.Address, error) {
	return _TokenStaking.Contract.OperatorsOf(&_TokenStaking.CallOpts, _address)
}

// OperatorsOf is a free data retrieval call binding the contract method 0xb534fbb6.
//
// Solidity: function operatorsOf(address _address) constant returns(address[])
func (_TokenStaking *TokenStakingCallerSession) OperatorsOf(_address common.Address) ([]common.Address, error) {
	return _TokenStaking.Contract.OperatorsOf(&_TokenStaking.CallOpts, _address)
}

// OwnerOf is a free data retrieval call binding the contract method 0x14afd79e.
//
// Solidity: function ownerOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingCaller) OwnerOf(opts *bind.CallOpts, _operator common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "ownerOf", _operator)
	return *ret0, err
}

// OwnerOf is a free data retrieval call binding the contract method 0x14afd79e.
//
// Solidity: function ownerOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingSession) OwnerOf(_operator common.Address) (common.Address, error) {
	return _TokenStaking.Contract.OwnerOf(&_TokenStaking.CallOpts, _operator)
}

// OwnerOf is a free data retrieval call binding the contract method 0x14afd79e.
//
// Solidity: function ownerOf(address _operator) constant returns(address)
func (_TokenStaking *TokenStakingCallerSession) OwnerOf(_operator common.Address) (common.Address, error) {
	return _TokenStaking.Contract.OwnerOf(&_TokenStaking.CallOpts, _operator)
}

// OwnerOperators is a free data retrieval call binding the contract method 0x4239dd8d.
//
// Solidity: function ownerOperators(address , uint256 ) constant returns(address)
func (_TokenStaking *TokenStakingCaller) OwnerOperators(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "ownerOperators", arg0, arg1)
	return *ret0, err
}

// OwnerOperators is a free data retrieval call binding the contract method 0x4239dd8d.
//
// Solidity: function ownerOperators(address , uint256 ) constant returns(address)
func (_TokenStaking *TokenStakingSession) OwnerOperators(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _TokenStaking.Contract.OwnerOperators(&_TokenStaking.CallOpts, arg0, arg1)
}

// OwnerOperators is a free data retrieval call binding the contract method 0x4239dd8d.
//
// Solidity: function ownerOperators(address , uint256 ) constant returns(address)
func (_TokenStaking *TokenStakingCallerSession) OwnerOperators(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _TokenStaking.Contract.OwnerOperators(&_TokenStaking.CallOpts, arg0, arg1)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() constant returns(address)
func (_TokenStaking *TokenStakingCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "registry")
	return *ret0, err
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() constant returns(address)
func (_TokenStaking *TokenStakingSession) Registry() (common.Address, error) {
	return _TokenStaking.Contract.Registry(&_TokenStaking.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() constant returns(address)
func (_TokenStaking *TokenStakingCallerSession) Registry() (common.Address, error) {
	return _TokenStaking.Contract.Registry(&_TokenStaking.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenStaking *TokenStakingCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenStaking *TokenStakingSession) Token() (common.Address, error) {
	return _TokenStaking.Contract.Token(&_TokenStaking.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenStaking *TokenStakingCallerSession) Token() (common.Address, error) {
	return _TokenStaking.Contract.Token(&_TokenStaking.CallOpts)
}

// UndelegationPeriod is a free data retrieval call binding the contract method 0xfdd1f986.
//
// Solidity: function undelegationPeriod() constant returns(uint256)
func (_TokenStaking *TokenStakingCaller) UndelegationPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStaking.contract.Call(opts, out, "undelegationPeriod")
	return *ret0, err
}

// UndelegationPeriod is a free data retrieval call binding the contract method 0xfdd1f986.
//
// Solidity: function undelegationPeriod() constant returns(uint256)
func (_TokenStaking *TokenStakingSession) UndelegationPeriod() (*big.Int, error) {
	return _TokenStaking.Contract.UndelegationPeriod(&_TokenStaking.CallOpts)
}

// UndelegationPeriod is a free data retrieval call binding the contract method 0xfdd1f986.
//
// Solidity: function undelegationPeriod() constant returns(uint256)
func (_TokenStaking *TokenStakingCallerSession) UndelegationPeriod() (*big.Int, error) {
	return _TokenStaking.Contract.UndelegationPeriod(&_TokenStaking.CallOpts)
}

// AuthorizeOperatorContract is a paid mutator transaction binding the contract method 0xf1654783.
//
// Solidity: function authorizeOperatorContract(address _operator, address _operatorContract) returns()
func (_TokenStaking *TokenStakingTransactor) AuthorizeOperatorContract(opts *bind.TransactOpts, _operator common.Address, _operatorContract common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "authorizeOperatorContract", _operator, _operatorContract)
}

// AuthorizeOperatorContract is a paid mutator transaction binding the contract method 0xf1654783.
//
// Solidity: function authorizeOperatorContract(address _operator, address _operatorContract) returns()
func (_TokenStaking *TokenStakingSession) AuthorizeOperatorContract(_operator common.Address, _operatorContract common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.AuthorizeOperatorContract(&_TokenStaking.TransactOpts, _operator, _operatorContract)
}

// AuthorizeOperatorContract is a paid mutator transaction binding the contract method 0xf1654783.
//
// Solidity: function authorizeOperatorContract(address _operator, address _operatorContract) returns()
func (_TokenStaking *TokenStakingTransactorSession) AuthorizeOperatorContract(_operator common.Address, _operatorContract common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.AuthorizeOperatorContract(&_TokenStaking.TransactOpts, _operator, _operatorContract)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_TokenStaking *TokenStakingTransactor) CancelStake(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "cancelStake", _operator)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_TokenStaking *TokenStakingSession) CancelStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.CancelStake(&_TokenStaking.TransactOpts, _operator)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_TokenStaking *TokenStakingTransactorSession) CancelStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.CancelStake(&_TokenStaking.TransactOpts, _operator)
}

// ClaimDelegatedAuthority is a paid mutator transaction binding the contract method 0xa590ae36.
//
// Solidity: function claimDelegatedAuthority(address delegatedAuthoritySource) returns()
func (_TokenStaking *TokenStakingTransactor) ClaimDelegatedAuthority(opts *bind.TransactOpts, delegatedAuthoritySource common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "claimDelegatedAuthority", delegatedAuthoritySource)
}

// ClaimDelegatedAuthority is a paid mutator transaction binding the contract method 0xa590ae36.
//
// Solidity: function claimDelegatedAuthority(address delegatedAuthoritySource) returns()
func (_TokenStaking *TokenStakingSession) ClaimDelegatedAuthority(delegatedAuthoritySource common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.ClaimDelegatedAuthority(&_TokenStaking.TransactOpts, delegatedAuthoritySource)
}

// ClaimDelegatedAuthority is a paid mutator transaction binding the contract method 0xa590ae36.
//
// Solidity: function claimDelegatedAuthority(address delegatedAuthoritySource) returns()
func (_TokenStaking *TokenStakingTransactorSession) ClaimDelegatedAuthority(delegatedAuthoritySource common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.ClaimDelegatedAuthority(&_TokenStaking.TransactOpts, delegatedAuthoritySource)
}

// LockStake is a paid mutator transaction binding the contract method 0x21e1625e.
//
// Solidity: function lockStake(address operator, uint256 duration) returns()
func (_TokenStaking *TokenStakingTransactor) LockStake(opts *bind.TransactOpts, operator common.Address, duration *big.Int) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "lockStake", operator, duration)
}

// LockStake is a paid mutator transaction binding the contract method 0x21e1625e.
//
// Solidity: function lockStake(address operator, uint256 duration) returns()
func (_TokenStaking *TokenStakingSession) LockStake(operator common.Address, duration *big.Int) (*types.Transaction, error) {
	return _TokenStaking.Contract.LockStake(&_TokenStaking.TransactOpts, operator, duration)
}

// LockStake is a paid mutator transaction binding the contract method 0x21e1625e.
//
// Solidity: function lockStake(address operator, uint256 duration) returns()
func (_TokenStaking *TokenStakingTransactorSession) LockStake(operator common.Address, duration *big.Int) (*types.Transaction, error) {
	return _TokenStaking.Contract.LockStake(&_TokenStaking.TransactOpts, operator, duration)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes _extraData) returns()
func (_TokenStaking *TokenStakingTransactor) ReceiveApproval(opts *bind.TransactOpts, _from common.Address, _value *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "receiveApproval", _from, _value, _token, _extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes _extraData) returns()
func (_TokenStaking *TokenStakingSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenStaking.Contract.ReceiveApproval(&_TokenStaking.TransactOpts, _from, _value, _token, _extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _value, address _token, bytes _extraData) returns()
func (_TokenStaking *TokenStakingTransactorSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenStaking.Contract.ReceiveApproval(&_TokenStaking.TransactOpts, _from, _value, _token, _extraData)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_TokenStaking *TokenStakingTransactor) RecoverStake(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "recoverStake", _operator)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_TokenStaking *TokenStakingSession) RecoverStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.RecoverStake(&_TokenStaking.TransactOpts, _operator)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_TokenStaking *TokenStakingTransactorSession) RecoverStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.RecoverStake(&_TokenStaking.TransactOpts, _operator)
}

// ReleaseExpiredLock is a paid mutator transaction binding the contract method 0xd98f233d.
//
// Solidity: function releaseExpiredLock(address operator, address operatorContract) returns()
func (_TokenStaking *TokenStakingTransactor) ReleaseExpiredLock(opts *bind.TransactOpts, operator common.Address, operatorContract common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "releaseExpiredLock", operator, operatorContract)
}

// ReleaseExpiredLock is a paid mutator transaction binding the contract method 0xd98f233d.
//
// Solidity: function releaseExpiredLock(address operator, address operatorContract) returns()
func (_TokenStaking *TokenStakingSession) ReleaseExpiredLock(operator common.Address, operatorContract common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.ReleaseExpiredLock(&_TokenStaking.TransactOpts, operator, operatorContract)
}

// ReleaseExpiredLock is a paid mutator transaction binding the contract method 0xd98f233d.
//
// Solidity: function releaseExpiredLock(address operator, address operatorContract) returns()
func (_TokenStaking *TokenStakingTransactorSession) ReleaseExpiredLock(operator common.Address, operatorContract common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.ReleaseExpiredLock(&_TokenStaking.TransactOpts, operator, operatorContract)
}

// Seize is a paid mutator transaction binding the contract method 0x09055e90.
//
// Solidity: function seize(uint256 amountToSeize, uint256 rewardMultiplier, address tattletale, address[] misbehavedOperators) returns()
func (_TokenStaking *TokenStakingTransactor) Seize(opts *bind.TransactOpts, amountToSeize *big.Int, rewardMultiplier *big.Int, tattletale common.Address, misbehavedOperators []common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "seize", amountToSeize, rewardMultiplier, tattletale, misbehavedOperators)
}

// Seize is a paid mutator transaction binding the contract method 0x09055e90.
//
// Solidity: function seize(uint256 amountToSeize, uint256 rewardMultiplier, address tattletale, address[] misbehavedOperators) returns()
func (_TokenStaking *TokenStakingSession) Seize(amountToSeize *big.Int, rewardMultiplier *big.Int, tattletale common.Address, misbehavedOperators []common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.Seize(&_TokenStaking.TransactOpts, amountToSeize, rewardMultiplier, tattletale, misbehavedOperators)
}

// Seize is a paid mutator transaction binding the contract method 0x09055e90.
//
// Solidity: function seize(uint256 amountToSeize, uint256 rewardMultiplier, address tattletale, address[] misbehavedOperators) returns()
func (_TokenStaking *TokenStakingTransactorSession) Seize(amountToSeize *big.Int, rewardMultiplier *big.Int, tattletale common.Address, misbehavedOperators []common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.Seize(&_TokenStaking.TransactOpts, amountToSeize, rewardMultiplier, tattletale, misbehavedOperators)
}

// Slash is a paid mutator transaction binding the contract method 0x8e49aa7a.
//
// Solidity: function slash(uint256 amountToSlash, address[] misbehavedOperators) returns()
func (_TokenStaking *TokenStakingTransactor) Slash(opts *bind.TransactOpts, amountToSlash *big.Int, misbehavedOperators []common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "slash", amountToSlash, misbehavedOperators)
}

// Slash is a paid mutator transaction binding the contract method 0x8e49aa7a.
//
// Solidity: function slash(uint256 amountToSlash, address[] misbehavedOperators) returns()
func (_TokenStaking *TokenStakingSession) Slash(amountToSlash *big.Int, misbehavedOperators []common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.Slash(&_TokenStaking.TransactOpts, amountToSlash, misbehavedOperators)
}

// Slash is a paid mutator transaction binding the contract method 0x8e49aa7a.
//
// Solidity: function slash(uint256 amountToSlash, address[] misbehavedOperators) returns()
func (_TokenStaking *TokenStakingTransactorSession) Slash(amountToSlash *big.Int, misbehavedOperators []common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.Slash(&_TokenStaking.TransactOpts, amountToSlash, misbehavedOperators)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_TokenStaking *TokenStakingTransactor) Undelegate(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "undelegate", _operator)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_TokenStaking *TokenStakingSession) Undelegate(_operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.Undelegate(&_TokenStaking.TransactOpts, _operator)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_TokenStaking *TokenStakingTransactorSession) Undelegate(_operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.Undelegate(&_TokenStaking.TransactOpts, _operator)
}

// UndelegateAt is a paid mutator transaction binding the contract method 0x5139a6c5.
//
// Solidity: function undelegateAt(address _operator, uint256 _undelegationTimestamp) returns()
func (_TokenStaking *TokenStakingTransactor) UndelegateAt(opts *bind.TransactOpts, _operator common.Address, _undelegationTimestamp *big.Int) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "undelegateAt", _operator, _undelegationTimestamp)
}

// UndelegateAt is a paid mutator transaction binding the contract method 0x5139a6c5.
//
// Solidity: function undelegateAt(address _operator, uint256 _undelegationTimestamp) returns()
func (_TokenStaking *TokenStakingSession) UndelegateAt(_operator common.Address, _undelegationTimestamp *big.Int) (*types.Transaction, error) {
	return _TokenStaking.Contract.UndelegateAt(&_TokenStaking.TransactOpts, _operator, _undelegationTimestamp)
}

// UndelegateAt is a paid mutator transaction binding the contract method 0x5139a6c5.
//
// Solidity: function undelegateAt(address _operator, uint256 _undelegationTimestamp) returns()
func (_TokenStaking *TokenStakingTransactorSession) UndelegateAt(_operator common.Address, _undelegationTimestamp *big.Int) (*types.Transaction, error) {
	return _TokenStaking.Contract.UndelegateAt(&_TokenStaking.TransactOpts, _operator, _undelegationTimestamp)
}

// UnlockStake is a paid mutator transaction binding the contract method 0x4a1ce599.
//
// Solidity: function unlockStake(address operator) returns()
func (_TokenStaking *TokenStakingTransactor) UnlockStake(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.contract.Transact(opts, "unlockStake", operator)
}

// UnlockStake is a paid mutator transaction binding the contract method 0x4a1ce599.
//
// Solidity: function unlockStake(address operator) returns()
func (_TokenStaking *TokenStakingSession) UnlockStake(operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.UnlockStake(&_TokenStaking.TransactOpts, operator)
}

// UnlockStake is a paid mutator transaction binding the contract method 0x4a1ce599.
//
// Solidity: function unlockStake(address operator) returns()
func (_TokenStaking *TokenStakingTransactorSession) UnlockStake(operator common.Address) (*types.Transaction, error) {
	return _TokenStaking.Contract.UnlockStake(&_TokenStaking.TransactOpts, operator)
}

// TokenStakingExpiredLockReleasedIterator is returned from FilterExpiredLockReleased and is used to iterate over the raw logs and unpacked data for ExpiredLockReleased events raised by the TokenStaking contract.
type TokenStakingExpiredLockReleasedIterator struct {
	Event *TokenStakingExpiredLockReleased // Event containing the contract specifics and raw log

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
func (it *TokenStakingExpiredLockReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingExpiredLockReleased)
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
		it.Event = new(TokenStakingExpiredLockReleased)
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
func (it *TokenStakingExpiredLockReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingExpiredLockReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingExpiredLockReleased represents a ExpiredLockReleased event raised by the TokenStaking contract.
type TokenStakingExpiredLockReleased struct {
	Operator    common.Address
	LockCreator common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterExpiredLockReleased is a free log retrieval operation binding the contract event 0x53232732101f51eae5873796ea83c72cecad5d155b851edfc11732b9dd4006f6.
//
// Solidity: event ExpiredLockReleased(address indexed operator, address lockCreator)
func (_TokenStaking *TokenStakingFilterer) FilterExpiredLockReleased(opts *bind.FilterOpts, operator []common.Address) (*TokenStakingExpiredLockReleasedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.FilterLogs(opts, "ExpiredLockReleased", operatorRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingExpiredLockReleasedIterator{contract: _TokenStaking.contract, event: "ExpiredLockReleased", logs: logs, sub: sub}, nil
}

// WatchExpiredLockReleased is a free log subscription operation binding the contract event 0x53232732101f51eae5873796ea83c72cecad5d155b851edfc11732b9dd4006f6.
//
// Solidity: event ExpiredLockReleased(address indexed operator, address lockCreator)
func (_TokenStaking *TokenStakingFilterer) WatchExpiredLockReleased(opts *bind.WatchOpts, sink chan<- *TokenStakingExpiredLockReleased, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.WatchLogs(opts, "ExpiredLockReleased", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingExpiredLockReleased)
				if err := _TokenStaking.contract.UnpackLog(event, "ExpiredLockReleased", log); err != nil {
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

// ParseExpiredLockReleased is a log parse operation binding the contract event 0x53232732101f51eae5873796ea83c72cecad5d155b851edfc11732b9dd4006f6.
//
// Solidity: event ExpiredLockReleased(address indexed operator, address lockCreator)
func (_TokenStaking *TokenStakingFilterer) ParseExpiredLockReleased(log types.Log) (*TokenStakingExpiredLockReleased, error) {
	event := new(TokenStakingExpiredLockReleased)
	if err := _TokenStaking.contract.UnpackLog(event, "ExpiredLockReleased", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingLockReleasedIterator is returned from FilterLockReleased and is used to iterate over the raw logs and unpacked data for LockReleased events raised by the TokenStaking contract.
type TokenStakingLockReleasedIterator struct {
	Event *TokenStakingLockReleased // Event containing the contract specifics and raw log

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
func (it *TokenStakingLockReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingLockReleased)
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
		it.Event = new(TokenStakingLockReleased)
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
func (it *TokenStakingLockReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingLockReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingLockReleased represents a LockReleased event raised by the TokenStaking contract.
type TokenStakingLockReleased struct {
	Operator    common.Address
	LockCreator common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLockReleased is a free log retrieval operation binding the contract event 0x9519d27283057289b75ef2605d6818602822861717fc48c918d37fe1fdc523f4.
//
// Solidity: event LockReleased(address indexed operator, address lockCreator)
func (_TokenStaking *TokenStakingFilterer) FilterLockReleased(opts *bind.FilterOpts, operator []common.Address) (*TokenStakingLockReleasedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.FilterLogs(opts, "LockReleased", operatorRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingLockReleasedIterator{contract: _TokenStaking.contract, event: "LockReleased", logs: logs, sub: sub}, nil
}

// WatchLockReleased is a free log subscription operation binding the contract event 0x9519d27283057289b75ef2605d6818602822861717fc48c918d37fe1fdc523f4.
//
// Solidity: event LockReleased(address indexed operator, address lockCreator)
func (_TokenStaking *TokenStakingFilterer) WatchLockReleased(opts *bind.WatchOpts, sink chan<- *TokenStakingLockReleased, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.WatchLogs(opts, "LockReleased", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingLockReleased)
				if err := _TokenStaking.contract.UnpackLog(event, "LockReleased", log); err != nil {
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

// ParseLockReleased is a log parse operation binding the contract event 0x9519d27283057289b75ef2605d6818602822861717fc48c918d37fe1fdc523f4.
//
// Solidity: event LockReleased(address indexed operator, address lockCreator)
func (_TokenStaking *TokenStakingFilterer) ParseLockReleased(log types.Log) (*TokenStakingLockReleased, error) {
	event := new(TokenStakingLockReleased)
	if err := _TokenStaking.contract.UnpackLog(event, "LockReleased", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingRecoveredStakeIterator is returned from FilterRecoveredStake and is used to iterate over the raw logs and unpacked data for RecoveredStake events raised by the TokenStaking contract.
type TokenStakingRecoveredStakeIterator struct {
	Event *TokenStakingRecoveredStake // Event containing the contract specifics and raw log

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
func (it *TokenStakingRecoveredStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingRecoveredStake)
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
		it.Event = new(TokenStakingRecoveredStake)
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
func (it *TokenStakingRecoveredStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingRecoveredStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingRecoveredStake represents a RecoveredStake event raised by the TokenStaking contract.
type TokenStakingRecoveredStake struct {
	Operator    common.Address
	RecoveredAt *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRecoveredStake is a free log retrieval operation binding the contract event 0x564f402e70218efee257fcbf84954e2e95e3ecc52cebe7b06c79955463752d5a.
//
// Solidity: event RecoveredStake(address operator, uint256 recoveredAt)
func (_TokenStaking *TokenStakingFilterer) FilterRecoveredStake(opts *bind.FilterOpts) (*TokenStakingRecoveredStakeIterator, error) {

	logs, sub, err := _TokenStaking.contract.FilterLogs(opts, "RecoveredStake")
	if err != nil {
		return nil, err
	}
	return &TokenStakingRecoveredStakeIterator{contract: _TokenStaking.contract, event: "RecoveredStake", logs: logs, sub: sub}, nil
}

// WatchRecoveredStake is a free log subscription operation binding the contract event 0x564f402e70218efee257fcbf84954e2e95e3ecc52cebe7b06c79955463752d5a.
//
// Solidity: event RecoveredStake(address operator, uint256 recoveredAt)
func (_TokenStaking *TokenStakingFilterer) WatchRecoveredStake(opts *bind.WatchOpts, sink chan<- *TokenStakingRecoveredStake) (event.Subscription, error) {

	logs, sub, err := _TokenStaking.contract.WatchLogs(opts, "RecoveredStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingRecoveredStake)
				if err := _TokenStaking.contract.UnpackLog(event, "RecoveredStake", log); err != nil {
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

// ParseRecoveredStake is a log parse operation binding the contract event 0x564f402e70218efee257fcbf84954e2e95e3ecc52cebe7b06c79955463752d5a.
//
// Solidity: event RecoveredStake(address operator, uint256 recoveredAt)
func (_TokenStaking *TokenStakingFilterer) ParseRecoveredStake(log types.Log) (*TokenStakingRecoveredStake, error) {
	event := new(TokenStakingRecoveredStake)
	if err := _TokenStaking.contract.UnpackLog(event, "RecoveredStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingStakeLockedIterator is returned from FilterStakeLocked and is used to iterate over the raw logs and unpacked data for StakeLocked events raised by the TokenStaking contract.
type TokenStakingStakeLockedIterator struct {
	Event *TokenStakingStakeLocked // Event containing the contract specifics and raw log

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
func (it *TokenStakingStakeLockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingStakeLocked)
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
		it.Event = new(TokenStakingStakeLocked)
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
func (it *TokenStakingStakeLockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingStakeLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingStakeLocked represents a StakeLocked event raised by the TokenStaking contract.
type TokenStakingStakeLocked struct {
	Operator    common.Address
	LockCreator common.Address
	Until       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStakeLocked is a free log retrieval operation binding the contract event 0x82358c8f3a8a41c7cae8a1196ae5106f7b58ce60eb38b7bc6fe3086d079d2a4e.
//
// Solidity: event StakeLocked(address indexed operator, address lockCreator, uint256 until)
func (_TokenStaking *TokenStakingFilterer) FilterStakeLocked(opts *bind.FilterOpts, operator []common.Address) (*TokenStakingStakeLockedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.FilterLogs(opts, "StakeLocked", operatorRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingStakeLockedIterator{contract: _TokenStaking.contract, event: "StakeLocked", logs: logs, sub: sub}, nil
}

// WatchStakeLocked is a free log subscription operation binding the contract event 0x82358c8f3a8a41c7cae8a1196ae5106f7b58ce60eb38b7bc6fe3086d079d2a4e.
//
// Solidity: event StakeLocked(address indexed operator, address lockCreator, uint256 until)
func (_TokenStaking *TokenStakingFilterer) WatchStakeLocked(opts *bind.WatchOpts, sink chan<- *TokenStakingStakeLocked, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.WatchLogs(opts, "StakeLocked", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingStakeLocked)
				if err := _TokenStaking.contract.UnpackLog(event, "StakeLocked", log); err != nil {
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

// ParseStakeLocked is a log parse operation binding the contract event 0x82358c8f3a8a41c7cae8a1196ae5106f7b58ce60eb38b7bc6fe3086d079d2a4e.
//
// Solidity: event StakeLocked(address indexed operator, address lockCreator, uint256 until)
func (_TokenStaking *TokenStakingFilterer) ParseStakeLocked(log types.Log) (*TokenStakingStakeLocked, error) {
	event := new(TokenStakingStakeLocked)
	if err := _TokenStaking.contract.UnpackLog(event, "StakeLocked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the TokenStaking contract.
type TokenStakingStakedIterator struct {
	Event *TokenStakingStaked // Event containing the contract specifics and raw log

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
func (it *TokenStakingStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingStaked)
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
		it.Event = new(TokenStakingStaked)
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
func (it *TokenStakingStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingStaked represents a Staked event raised by the TokenStaking contract.
type TokenStakingStaked struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed from, uint256 value)
func (_TokenStaking *TokenStakingFilterer) FilterStaked(opts *bind.FilterOpts, from []common.Address) (*TokenStakingStakedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TokenStaking.contract.FilterLogs(opts, "Staked", fromRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingStakedIterator{contract: _TokenStaking.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed from, uint256 value)
func (_TokenStaking *TokenStakingFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *TokenStakingStaked, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TokenStaking.contract.WatchLogs(opts, "Staked", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingStaked)
				if err := _TokenStaking.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed from, uint256 value)
func (_TokenStaking *TokenStakingFilterer) ParseStaked(log types.Log) (*TokenStakingStaked, error) {
	event := new(TokenStakingStaked)
	if err := _TokenStaking.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingTokensSeizedIterator is returned from FilterTokensSeized and is used to iterate over the raw logs and unpacked data for TokensSeized events raised by the TokenStaking contract.
type TokenStakingTokensSeizedIterator struct {
	Event *TokenStakingTokensSeized // Event containing the contract specifics and raw log

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
func (it *TokenStakingTokensSeizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingTokensSeized)
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
		it.Event = new(TokenStakingTokensSeized)
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
func (it *TokenStakingTokensSeizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingTokensSeizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingTokensSeized represents a TokensSeized event raised by the TokenStaking contract.
type TokenStakingTokensSeized struct {
	Operator common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokensSeized is a free log retrieval operation binding the contract event 0xc5c18d4a60510957a4dc6dbd6d40ecfcfa32ab9517528a6b4052729e48e183d3.
//
// Solidity: event TokensSeized(address indexed operator, uint256 amount)
func (_TokenStaking *TokenStakingFilterer) FilterTokensSeized(opts *bind.FilterOpts, operator []common.Address) (*TokenStakingTokensSeizedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.FilterLogs(opts, "TokensSeized", operatorRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingTokensSeizedIterator{contract: _TokenStaking.contract, event: "TokensSeized", logs: logs, sub: sub}, nil
}

// WatchTokensSeized is a free log subscription operation binding the contract event 0xc5c18d4a60510957a4dc6dbd6d40ecfcfa32ab9517528a6b4052729e48e183d3.
//
// Solidity: event TokensSeized(address indexed operator, uint256 amount)
func (_TokenStaking *TokenStakingFilterer) WatchTokensSeized(opts *bind.WatchOpts, sink chan<- *TokenStakingTokensSeized, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.WatchLogs(opts, "TokensSeized", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingTokensSeized)
				if err := _TokenStaking.contract.UnpackLog(event, "TokensSeized", log); err != nil {
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

// ParseTokensSeized is a log parse operation binding the contract event 0xc5c18d4a60510957a4dc6dbd6d40ecfcfa32ab9517528a6b4052729e48e183d3.
//
// Solidity: event TokensSeized(address indexed operator, uint256 amount)
func (_TokenStaking *TokenStakingFilterer) ParseTokensSeized(log types.Log) (*TokenStakingTokensSeized, error) {
	event := new(TokenStakingTokensSeized)
	if err := _TokenStaking.contract.UnpackLog(event, "TokensSeized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingTokensSlashedIterator is returned from FilterTokensSlashed and is used to iterate over the raw logs and unpacked data for TokensSlashed events raised by the TokenStaking contract.
type TokenStakingTokensSlashedIterator struct {
	Event *TokenStakingTokensSlashed // Event containing the contract specifics and raw log

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
func (it *TokenStakingTokensSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingTokensSlashed)
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
		it.Event = new(TokenStakingTokensSlashed)
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
func (it *TokenStakingTokensSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingTokensSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingTokensSlashed represents a TokensSlashed event raised by the TokenStaking contract.
type TokenStakingTokensSlashed struct {
	Operator common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokensSlashed is a free log retrieval operation binding the contract event 0xf71d4eeef63c86ffc692e306c623b01b8373e091da8ee455aa0fe26305a4981c.
//
// Solidity: event TokensSlashed(address indexed operator, uint256 amount)
func (_TokenStaking *TokenStakingFilterer) FilterTokensSlashed(opts *bind.FilterOpts, operator []common.Address) (*TokenStakingTokensSlashedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.FilterLogs(opts, "TokensSlashed", operatorRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingTokensSlashedIterator{contract: _TokenStaking.contract, event: "TokensSlashed", logs: logs, sub: sub}, nil
}

// WatchTokensSlashed is a free log subscription operation binding the contract event 0xf71d4eeef63c86ffc692e306c623b01b8373e091da8ee455aa0fe26305a4981c.
//
// Solidity: event TokensSlashed(address indexed operator, uint256 amount)
func (_TokenStaking *TokenStakingFilterer) WatchTokensSlashed(opts *bind.WatchOpts, sink chan<- *TokenStakingTokensSlashed, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.WatchLogs(opts, "TokensSlashed", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingTokensSlashed)
				if err := _TokenStaking.contract.UnpackLog(event, "TokensSlashed", log); err != nil {
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

// ParseTokensSlashed is a log parse operation binding the contract event 0xf71d4eeef63c86ffc692e306c623b01b8373e091da8ee455aa0fe26305a4981c.
//
// Solidity: event TokensSlashed(address indexed operator, uint256 amount)
func (_TokenStaking *TokenStakingFilterer) ParseTokensSlashed(log types.Log) (*TokenStakingTokensSlashed, error) {
	event := new(TokenStakingTokensSlashed)
	if err := _TokenStaking.contract.UnpackLog(event, "TokensSlashed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingUndelegatedIterator is returned from FilterUndelegated and is used to iterate over the raw logs and unpacked data for Undelegated events raised by the TokenStaking contract.
type TokenStakingUndelegatedIterator struct {
	Event *TokenStakingUndelegated // Event containing the contract specifics and raw log

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
func (it *TokenStakingUndelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingUndelegated)
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
		it.Event = new(TokenStakingUndelegated)
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
func (it *TokenStakingUndelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingUndelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingUndelegated represents a Undelegated event raised by the TokenStaking contract.
type TokenStakingUndelegated struct {
	Operator      common.Address
	UndelegatedAt *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUndelegated is a free log retrieval operation binding the contract event 0x4ae68879209bc4b489a38251122202a3653305e3d95a27baf7a5681410c90b38.
//
// Solidity: event Undelegated(address indexed operator, uint256 undelegatedAt)
func (_TokenStaking *TokenStakingFilterer) FilterUndelegated(opts *bind.FilterOpts, operator []common.Address) (*TokenStakingUndelegatedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.FilterLogs(opts, "Undelegated", operatorRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingUndelegatedIterator{contract: _TokenStaking.contract, event: "Undelegated", logs: logs, sub: sub}, nil
}

// WatchUndelegated is a free log subscription operation binding the contract event 0x4ae68879209bc4b489a38251122202a3653305e3d95a27baf7a5681410c90b38.
//
// Solidity: event Undelegated(address indexed operator, uint256 undelegatedAt)
func (_TokenStaking *TokenStakingFilterer) WatchUndelegated(opts *bind.WatchOpts, sink chan<- *TokenStakingUndelegated, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TokenStaking.contract.WatchLogs(opts, "Undelegated", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingUndelegated)
				if err := _TokenStaking.contract.UnpackLog(event, "Undelegated", log); err != nil {
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

// ParseUndelegated is a log parse operation binding the contract event 0x4ae68879209bc4b489a38251122202a3653305e3d95a27baf7a5681410c90b38.
//
// Solidity: event Undelegated(address indexed operator, uint256 undelegatedAt)
func (_TokenStaking *TokenStakingFilterer) ParseUndelegated(log types.Log) (*TokenStakingUndelegated, error) {
	event := new(TokenStakingUndelegated)
	if err := _TokenStaking.contract.UnpackLog(event, "Undelegated", log); err != nil {
		return nil, err
	}
	return event, nil
}
