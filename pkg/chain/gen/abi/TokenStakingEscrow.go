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

// TokenStakingEscrowABI is the input ABI used to generate the binding from.
const TokenStakingEscrowABI = "[{\"inputs\":[{\"internalType\":\"contractKeepToken\",\"name\":\"_keepToken\",\"type\":\"address\"},{\"internalType\":\"contractTokenGrant\",\"name\":\"_tokenGrant\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOperator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"grantId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositRedelegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"grantee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"grantId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"grantManager\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"escrow\",\"type\":\"address\"}],\"name\":\"EscrowAuthorized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"grantManager\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RevokedDepositWithdrawn\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"anotherEscrow\",\"type\":\"address\"}],\"name\":\"authorizeEscrow\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"availableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"depositGrantId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"depositRedelegatedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"depositWithdrawnAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"depositedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"managedGrant\",\"type\":\"address\"}],\"name\":\"getManagedGrantee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"hasDeposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"keepToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receivingEscrow\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"previousOperator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"redelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenGrant\",\"outputs\":[{\"internalType\":\"contractTokenGrant\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"withdrawRevoked\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"withdrawToManagedGrantee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// TokenStakingEscrow is an auto generated Go binding around an Ethereum contract.
type TokenStakingEscrow struct {
	TokenStakingEscrowCaller     // Read-only binding to the contract
	TokenStakingEscrowTransactor // Write-only binding to the contract
	TokenStakingEscrowFilterer   // Log filterer for contract events
}

// TokenStakingEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenStakingEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStakingEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenStakingEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStakingEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenStakingEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStakingEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenStakingEscrowSession struct {
	Contract     *TokenStakingEscrow // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TokenStakingEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenStakingEscrowCallerSession struct {
	Contract *TokenStakingEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// TokenStakingEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenStakingEscrowTransactorSession struct {
	Contract     *TokenStakingEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// TokenStakingEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenStakingEscrowRaw struct {
	Contract *TokenStakingEscrow // Generic contract binding to access the raw methods on
}

// TokenStakingEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenStakingEscrowCallerRaw struct {
	Contract *TokenStakingEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// TokenStakingEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenStakingEscrowTransactorRaw struct {
	Contract *TokenStakingEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenStakingEscrow creates a new instance of TokenStakingEscrow, bound to a specific deployed contract.
func NewTokenStakingEscrow(address common.Address, backend bind.ContractBackend) (*TokenStakingEscrow, error) {
	contract, err := bindTokenStakingEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrow{TokenStakingEscrowCaller: TokenStakingEscrowCaller{contract: contract}, TokenStakingEscrowTransactor: TokenStakingEscrowTransactor{contract: contract}, TokenStakingEscrowFilterer: TokenStakingEscrowFilterer{contract: contract}}, nil
}

// NewTokenStakingEscrowCaller creates a new read-only instance of TokenStakingEscrow, bound to a specific deployed contract.
func NewTokenStakingEscrowCaller(address common.Address, caller bind.ContractCaller) (*TokenStakingEscrowCaller, error) {
	contract, err := bindTokenStakingEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowCaller{contract: contract}, nil
}

// NewTokenStakingEscrowTransactor creates a new write-only instance of TokenStakingEscrow, bound to a specific deployed contract.
func NewTokenStakingEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenStakingEscrowTransactor, error) {
	contract, err := bindTokenStakingEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowTransactor{contract: contract}, nil
}

// NewTokenStakingEscrowFilterer creates a new log filterer instance of TokenStakingEscrow, bound to a specific deployed contract.
func NewTokenStakingEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenStakingEscrowFilterer, error) {
	contract, err := bindTokenStakingEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowFilterer{contract: contract}, nil
}

// bindTokenStakingEscrow binds a generic wrapper to an already deployed contract.
func bindTokenStakingEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenStakingEscrowABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenStakingEscrow *TokenStakingEscrowRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenStakingEscrow.Contract.TokenStakingEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenStakingEscrow *TokenStakingEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.TokenStakingEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenStakingEscrow *TokenStakingEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.TokenStakingEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenStakingEscrow *TokenStakingEscrowCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenStakingEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenStakingEscrow *TokenStakingEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenStakingEscrow *TokenStakingEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.contract.Transact(opts, method, params...)
}

// AvailableAmount is a free data retrieval call binding the contract method 0x654259dd.
//
// Solidity: function availableAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) AvailableAmount(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "availableAmount", operator)
	return *ret0, err
}

// AvailableAmount is a free data retrieval call binding the contract method 0x654259dd.
//
// Solidity: function availableAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowSession) AvailableAmount(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.AvailableAmount(&_TokenStakingEscrow.CallOpts, operator)
}

// AvailableAmount is a free data retrieval call binding the contract method 0x654259dd.
//
// Solidity: function availableAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) AvailableAmount(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.AvailableAmount(&_TokenStakingEscrow.CallOpts, operator)
}

// DepositGrantId is a free data retrieval call binding the contract method 0x7589d761.
//
// Solidity: function depositGrantId(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) DepositGrantId(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "depositGrantId", operator)
	return *ret0, err
}

// DepositGrantId is a free data retrieval call binding the contract method 0x7589d761.
//
// Solidity: function depositGrantId(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowSession) DepositGrantId(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.DepositGrantId(&_TokenStakingEscrow.CallOpts, operator)
}

// DepositGrantId is a free data retrieval call binding the contract method 0x7589d761.
//
// Solidity: function depositGrantId(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) DepositGrantId(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.DepositGrantId(&_TokenStakingEscrow.CallOpts, operator)
}

// DepositRedelegatedAmount is a free data retrieval call binding the contract method 0x2c5cd1ed.
//
// Solidity: function depositRedelegatedAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) DepositRedelegatedAmount(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "depositRedelegatedAmount", operator)
	return *ret0, err
}

// DepositRedelegatedAmount is a free data retrieval call binding the contract method 0x2c5cd1ed.
//
// Solidity: function depositRedelegatedAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowSession) DepositRedelegatedAmount(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.DepositRedelegatedAmount(&_TokenStakingEscrow.CallOpts, operator)
}

// DepositRedelegatedAmount is a free data retrieval call binding the contract method 0x2c5cd1ed.
//
// Solidity: function depositRedelegatedAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) DepositRedelegatedAmount(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.DepositRedelegatedAmount(&_TokenStakingEscrow.CallOpts, operator)
}

// DepositWithdrawnAmount is a free data retrieval call binding the contract method 0xa6cedc84.
//
// Solidity: function depositWithdrawnAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) DepositWithdrawnAmount(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "depositWithdrawnAmount", operator)
	return *ret0, err
}

// DepositWithdrawnAmount is a free data retrieval call binding the contract method 0xa6cedc84.
//
// Solidity: function depositWithdrawnAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowSession) DepositWithdrawnAmount(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.DepositWithdrawnAmount(&_TokenStakingEscrow.CallOpts, operator)
}

// DepositWithdrawnAmount is a free data retrieval call binding the contract method 0xa6cedc84.
//
// Solidity: function depositWithdrawnAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) DepositWithdrawnAmount(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.DepositWithdrawnAmount(&_TokenStakingEscrow.CallOpts, operator)
}

// DepositedAmount is a free data retrieval call binding the contract method 0x4a4643f7.
//
// Solidity: function depositedAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) DepositedAmount(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "depositedAmount", operator)
	return *ret0, err
}

// DepositedAmount is a free data retrieval call binding the contract method 0x4a4643f7.
//
// Solidity: function depositedAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowSession) DepositedAmount(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.DepositedAmount(&_TokenStakingEscrow.CallOpts, operator)
}

// DepositedAmount is a free data retrieval call binding the contract method 0x4a4643f7.
//
// Solidity: function depositedAmount(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) DepositedAmount(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.DepositedAmount(&_TokenStakingEscrow.CallOpts, operator)
}

// GetManagedGrantee is a free data retrieval call binding the contract method 0x012b1973.
//
// Solidity: function getManagedGrantee(address managedGrant) constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) GetManagedGrantee(opts *bind.CallOpts, managedGrant common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "getManagedGrantee", managedGrant)
	return *ret0, err
}

// GetManagedGrantee is a free data retrieval call binding the contract method 0x012b1973.
//
// Solidity: function getManagedGrantee(address managedGrant) constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowSession) GetManagedGrantee(managedGrant common.Address) (common.Address, error) {
	return _TokenStakingEscrow.Contract.GetManagedGrantee(&_TokenStakingEscrow.CallOpts, managedGrant)
}

// GetManagedGrantee is a free data retrieval call binding the contract method 0x012b1973.
//
// Solidity: function getManagedGrantee(address managedGrant) constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) GetManagedGrantee(managedGrant common.Address) (common.Address, error) {
	return _TokenStakingEscrow.Contract.GetManagedGrantee(&_TokenStakingEscrow.CallOpts, managedGrant)
}

// HasDeposit is a free data retrieval call binding the contract method 0x1c48c7ec.
//
// Solidity: function hasDeposit(address operator) constant returns(bool)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) HasDeposit(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "hasDeposit", operator)
	return *ret0, err
}

// HasDeposit is a free data retrieval call binding the contract method 0x1c48c7ec.
//
// Solidity: function hasDeposit(address operator) constant returns(bool)
func (_TokenStakingEscrow *TokenStakingEscrowSession) HasDeposit(operator common.Address) (bool, error) {
	return _TokenStakingEscrow.Contract.HasDeposit(&_TokenStakingEscrow.CallOpts, operator)
}

// HasDeposit is a free data retrieval call binding the contract method 0x1c48c7ec.
//
// Solidity: function hasDeposit(address operator) constant returns(bool)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) HasDeposit(operator common.Address) (bool, error) {
	return _TokenStakingEscrow.Contract.HasDeposit(&_TokenStakingEscrow.CallOpts, operator)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_TokenStakingEscrow *TokenStakingEscrowSession) IsOwner() (bool, error) {
	return _TokenStakingEscrow.Contract.IsOwner(&_TokenStakingEscrow.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) IsOwner() (bool, error) {
	return _TokenStakingEscrow.Contract.IsOwner(&_TokenStakingEscrow.CallOpts)
}

// KeepToken is a free data retrieval call binding the contract method 0x467c5adf.
//
// Solidity: function keepToken() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) KeepToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "keepToken")
	return *ret0, err
}

// KeepToken is a free data retrieval call binding the contract method 0x467c5adf.
//
// Solidity: function keepToken() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowSession) KeepToken() (common.Address, error) {
	return _TokenStakingEscrow.Contract.KeepToken(&_TokenStakingEscrow.CallOpts)
}

// KeepToken is a free data retrieval call binding the contract method 0x467c5adf.
//
// Solidity: function keepToken() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) KeepToken() (common.Address, error) {
	return _TokenStakingEscrow.Contract.KeepToken(&_TokenStakingEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowSession) Owner() (common.Address, error) {
	return _TokenStakingEscrow.Contract.Owner(&_TokenStakingEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) Owner() (common.Address, error) {
	return _TokenStakingEscrow.Contract.Owner(&_TokenStakingEscrow.CallOpts)
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) TokenGrant(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "tokenGrant")
	return *ret0, err
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowSession) TokenGrant() (common.Address, error) {
	return _TokenStakingEscrow.Contract.TokenGrant(&_TokenStakingEscrow.CallOpts)
}

// TokenGrant is a free data retrieval call binding the contract method 0xd92db09f.
//
// Solidity: function tokenGrant() constant returns(address)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) TokenGrant() (common.Address, error) {
	return _TokenStakingEscrow.Contract.TokenGrant(&_TokenStakingEscrow.CallOpts)
}

// Withdrawable is a free data retrieval call binding the contract method 0xce513b6f.
//
// Solidity: function withdrawable(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCaller) Withdrawable(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenStakingEscrow.contract.Call(opts, out, "withdrawable", operator)
	return *ret0, err
}

// Withdrawable is a free data retrieval call binding the contract method 0xce513b6f.
//
// Solidity: function withdrawable(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowSession) Withdrawable(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.Withdrawable(&_TokenStakingEscrow.CallOpts, operator)
}

// Withdrawable is a free data retrieval call binding the contract method 0xce513b6f.
//
// Solidity: function withdrawable(address operator) constant returns(uint256)
func (_TokenStakingEscrow *TokenStakingEscrowCallerSession) Withdrawable(operator common.Address) (*big.Int, error) {
	return _TokenStakingEscrow.Contract.Withdrawable(&_TokenStakingEscrow.CallOpts, operator)
}

// AuthorizeEscrow is a paid mutator transaction binding the contract method 0x4d6f9646.
//
// Solidity: function authorizeEscrow(address anotherEscrow) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) AuthorizeEscrow(opts *bind.TransactOpts, anotherEscrow common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "authorizeEscrow", anotherEscrow)
}

// AuthorizeEscrow is a paid mutator transaction binding the contract method 0x4d6f9646.
//
// Solidity: function authorizeEscrow(address anotherEscrow) returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) AuthorizeEscrow(anotherEscrow common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.AuthorizeEscrow(&_TokenStakingEscrow.TransactOpts, anotherEscrow)
}

// AuthorizeEscrow is a paid mutator transaction binding the contract method 0x4d6f9646.
//
// Solidity: function authorizeEscrow(address anotherEscrow) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) AuthorizeEscrow(anotherEscrow common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.AuthorizeEscrow(&_TokenStakingEscrow.TransactOpts, anotherEscrow)
}

// Migrate is a paid mutator transaction binding the contract method 0x1068361f.
//
// Solidity: function migrate(address operator, address receivingEscrow) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) Migrate(opts *bind.TransactOpts, operator common.Address, receivingEscrow common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "migrate", operator, receivingEscrow)
}

// Migrate is a paid mutator transaction binding the contract method 0x1068361f.
//
// Solidity: function migrate(address operator, address receivingEscrow) returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) Migrate(operator common.Address, receivingEscrow common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.Migrate(&_TokenStakingEscrow.TransactOpts, operator, receivingEscrow)
}

// Migrate is a paid mutator transaction binding the contract method 0x1068361f.
//
// Solidity: function migrate(address operator, address receivingEscrow) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) Migrate(operator common.Address, receivingEscrow common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.Migrate(&_TokenStakingEscrow.TransactOpts, operator, receivingEscrow)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 value, address token, bytes extraData) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) ReceiveApproval(opts *bind.TransactOpts, from common.Address, value *big.Int, token common.Address, extraData []byte) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "receiveApproval", from, value, token, extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 value, address token, bytes extraData) returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) ReceiveApproval(from common.Address, value *big.Int, token common.Address, extraData []byte) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.ReceiveApproval(&_TokenStakingEscrow.TransactOpts, from, value, token, extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 value, address token, bytes extraData) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) ReceiveApproval(from common.Address, value *big.Int, token common.Address, extraData []byte) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.ReceiveApproval(&_TokenStakingEscrow.TransactOpts, from, value, token, extraData)
}

// Redelegate is a paid mutator transaction binding the contract method 0xaf2b5888.
//
// Solidity: function redelegate(address previousOperator, uint256 amount, bytes extraData) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) Redelegate(opts *bind.TransactOpts, previousOperator common.Address, amount *big.Int, extraData []byte) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "redelegate", previousOperator, amount, extraData)
}

// Redelegate is a paid mutator transaction binding the contract method 0xaf2b5888.
//
// Solidity: function redelegate(address previousOperator, uint256 amount, bytes extraData) returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) Redelegate(previousOperator common.Address, amount *big.Int, extraData []byte) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.Redelegate(&_TokenStakingEscrow.TransactOpts, previousOperator, amount, extraData)
}

// Redelegate is a paid mutator transaction binding the contract method 0xaf2b5888.
//
// Solidity: function redelegate(address previousOperator, uint256 amount, bytes extraData) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) Redelegate(previousOperator common.Address, amount *big.Int, extraData []byte) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.Redelegate(&_TokenStakingEscrow.TransactOpts, previousOperator, amount, extraData)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.RenounceOwnership(&_TokenStakingEscrow.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.RenounceOwnership(&_TokenStakingEscrow.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.TransferOwnership(&_TokenStakingEscrow.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.TransferOwnership(&_TokenStakingEscrow.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) Withdraw(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "withdraw", operator)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) Withdraw(operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.Withdraw(&_TokenStakingEscrow.TransactOpts, operator)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) Withdraw(operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.Withdraw(&_TokenStakingEscrow.TransactOpts, operator)
}

// WithdrawRevoked is a paid mutator transaction binding the contract method 0xe5dd7bcb.
//
// Solidity: function withdrawRevoked(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) WithdrawRevoked(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "withdrawRevoked", operator)
}

// WithdrawRevoked is a paid mutator transaction binding the contract method 0xe5dd7bcb.
//
// Solidity: function withdrawRevoked(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) WithdrawRevoked(operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.WithdrawRevoked(&_TokenStakingEscrow.TransactOpts, operator)
}

// WithdrawRevoked is a paid mutator transaction binding the contract method 0xe5dd7bcb.
//
// Solidity: function withdrawRevoked(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) WithdrawRevoked(operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.WithdrawRevoked(&_TokenStakingEscrow.TransactOpts, operator)
}

// WithdrawToManagedGrantee is a paid mutator transaction binding the contract method 0x0e6cdbdd.
//
// Solidity: function withdrawToManagedGrantee(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactor) WithdrawToManagedGrantee(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.contract.Transact(opts, "withdrawToManagedGrantee", operator)
}

// WithdrawToManagedGrantee is a paid mutator transaction binding the contract method 0x0e6cdbdd.
//
// Solidity: function withdrawToManagedGrantee(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowSession) WithdrawToManagedGrantee(operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.WithdrawToManagedGrantee(&_TokenStakingEscrow.TransactOpts, operator)
}

// WithdrawToManagedGrantee is a paid mutator transaction binding the contract method 0x0e6cdbdd.
//
// Solidity: function withdrawToManagedGrantee(address operator) returns()
func (_TokenStakingEscrow *TokenStakingEscrowTransactorSession) WithdrawToManagedGrantee(operator common.Address) (*types.Transaction, error) {
	return _TokenStakingEscrow.Contract.WithdrawToManagedGrantee(&_TokenStakingEscrow.TransactOpts, operator)
}

// TokenStakingEscrowDepositRedelegatedIterator is returned from FilterDepositRedelegated and is used to iterate over the raw logs and unpacked data for DepositRedelegated events raised by the TokenStakingEscrow contract.
type TokenStakingEscrowDepositRedelegatedIterator struct {
	Event *TokenStakingEscrowDepositRedelegated // Event containing the contract specifics and raw log

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
func (it *TokenStakingEscrowDepositRedelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingEscrowDepositRedelegated)
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
		it.Event = new(TokenStakingEscrowDepositRedelegated)
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
func (it *TokenStakingEscrowDepositRedelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingEscrowDepositRedelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingEscrowDepositRedelegated represents a DepositRedelegated event raised by the TokenStakingEscrow contract.
type TokenStakingEscrowDepositRedelegated struct {
	PreviousOperator common.Address
	NewOperator      common.Address
	GrantId          *big.Int
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDepositRedelegated is a free log retrieval operation binding the contract event 0x2b7cac54d08d3b38cfada0cdfd02770b546ef37fbd5cfe413d4986cda6e56886.
//
// Solidity: event DepositRedelegated(address indexed previousOperator, address indexed newOperator, uint256 indexed grantId, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) FilterDepositRedelegated(opts *bind.FilterOpts, previousOperator []common.Address, newOperator []common.Address, grantId []*big.Int) (*TokenStakingEscrowDepositRedelegatedIterator, error) {

	var previousOperatorRule []interface{}
	for _, previousOperatorItem := range previousOperator {
		previousOperatorRule = append(previousOperatorRule, previousOperatorItem)
	}
	var newOperatorRule []interface{}
	for _, newOperatorItem := range newOperator {
		newOperatorRule = append(newOperatorRule, newOperatorItem)
	}
	var grantIdRule []interface{}
	for _, grantIdItem := range grantId {
		grantIdRule = append(grantIdRule, grantIdItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.FilterLogs(opts, "DepositRedelegated", previousOperatorRule, newOperatorRule, grantIdRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowDepositRedelegatedIterator{contract: _TokenStakingEscrow.contract, event: "DepositRedelegated", logs: logs, sub: sub}, nil
}

// WatchDepositRedelegated is a free log subscription operation binding the contract event 0x2b7cac54d08d3b38cfada0cdfd02770b546ef37fbd5cfe413d4986cda6e56886.
//
// Solidity: event DepositRedelegated(address indexed previousOperator, address indexed newOperator, uint256 indexed grantId, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) WatchDepositRedelegated(opts *bind.WatchOpts, sink chan<- *TokenStakingEscrowDepositRedelegated, previousOperator []common.Address, newOperator []common.Address, grantId []*big.Int) (event.Subscription, error) {

	var previousOperatorRule []interface{}
	for _, previousOperatorItem := range previousOperator {
		previousOperatorRule = append(previousOperatorRule, previousOperatorItem)
	}
	var newOperatorRule []interface{}
	for _, newOperatorItem := range newOperator {
		newOperatorRule = append(newOperatorRule, newOperatorItem)
	}
	var grantIdRule []interface{}
	for _, grantIdItem := range grantId {
		grantIdRule = append(grantIdRule, grantIdItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.WatchLogs(opts, "DepositRedelegated", previousOperatorRule, newOperatorRule, grantIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingEscrowDepositRedelegated)
				if err := _TokenStakingEscrow.contract.UnpackLog(event, "DepositRedelegated", log); err != nil {
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

// ParseDepositRedelegated is a log parse operation binding the contract event 0x2b7cac54d08d3b38cfada0cdfd02770b546ef37fbd5cfe413d4986cda6e56886.
//
// Solidity: event DepositRedelegated(address indexed previousOperator, address indexed newOperator, uint256 indexed grantId, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) ParseDepositRedelegated(log types.Log) (*TokenStakingEscrowDepositRedelegated, error) {
	event := new(TokenStakingEscrowDepositRedelegated)
	if err := _TokenStakingEscrow.contract.UnpackLog(event, "DepositRedelegated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingEscrowDepositWithdrawnIterator is returned from FilterDepositWithdrawn and is used to iterate over the raw logs and unpacked data for DepositWithdrawn events raised by the TokenStakingEscrow contract.
type TokenStakingEscrowDepositWithdrawnIterator struct {
	Event *TokenStakingEscrowDepositWithdrawn // Event containing the contract specifics and raw log

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
func (it *TokenStakingEscrowDepositWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingEscrowDepositWithdrawn)
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
		it.Event = new(TokenStakingEscrowDepositWithdrawn)
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
func (it *TokenStakingEscrowDepositWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingEscrowDepositWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingEscrowDepositWithdrawn represents a DepositWithdrawn event raised by the TokenStakingEscrow contract.
type TokenStakingEscrowDepositWithdrawn struct {
	Operator common.Address
	Grantee  common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDepositWithdrawn is a free log retrieval operation binding the contract event 0x2cd6435b1b961c13f55202979edd0765a809f69a539d8a477436c94c1211e43e.
//
// Solidity: event DepositWithdrawn(address indexed operator, address indexed grantee, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) FilterDepositWithdrawn(opts *bind.FilterOpts, operator []common.Address, grantee []common.Address) (*TokenStakingEscrowDepositWithdrawnIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var granteeRule []interface{}
	for _, granteeItem := range grantee {
		granteeRule = append(granteeRule, granteeItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.FilterLogs(opts, "DepositWithdrawn", operatorRule, granteeRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowDepositWithdrawnIterator{contract: _TokenStakingEscrow.contract, event: "DepositWithdrawn", logs: logs, sub: sub}, nil
}

// WatchDepositWithdrawn is a free log subscription operation binding the contract event 0x2cd6435b1b961c13f55202979edd0765a809f69a539d8a477436c94c1211e43e.
//
// Solidity: event DepositWithdrawn(address indexed operator, address indexed grantee, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) WatchDepositWithdrawn(opts *bind.WatchOpts, sink chan<- *TokenStakingEscrowDepositWithdrawn, operator []common.Address, grantee []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var granteeRule []interface{}
	for _, granteeItem := range grantee {
		granteeRule = append(granteeRule, granteeItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.WatchLogs(opts, "DepositWithdrawn", operatorRule, granteeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingEscrowDepositWithdrawn)
				if err := _TokenStakingEscrow.contract.UnpackLog(event, "DepositWithdrawn", log); err != nil {
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

// ParseDepositWithdrawn is a log parse operation binding the contract event 0x2cd6435b1b961c13f55202979edd0765a809f69a539d8a477436c94c1211e43e.
//
// Solidity: event DepositWithdrawn(address indexed operator, address indexed grantee, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) ParseDepositWithdrawn(log types.Log) (*TokenStakingEscrowDepositWithdrawn, error) {
	event := new(TokenStakingEscrowDepositWithdrawn)
	if err := _TokenStakingEscrow.contract.UnpackLog(event, "DepositWithdrawn", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingEscrowDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the TokenStakingEscrow contract.
type TokenStakingEscrowDepositedIterator struct {
	Event *TokenStakingEscrowDeposited // Event containing the contract specifics and raw log

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
func (it *TokenStakingEscrowDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingEscrowDeposited)
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
		it.Event = new(TokenStakingEscrowDeposited)
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
func (it *TokenStakingEscrowDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingEscrowDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingEscrowDeposited represents a Deposited event raised by the TokenStakingEscrow contract.
type TokenStakingEscrowDeposited struct {
	Operator common.Address
	GrantId  *big.Int
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x73a19dd210f1a7f902193214c0ee91dd35ee5b4d920cba8d519eca65a7b488ca.
//
// Solidity: event Deposited(address indexed operator, uint256 indexed grantId, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) FilterDeposited(opts *bind.FilterOpts, operator []common.Address, grantId []*big.Int) (*TokenStakingEscrowDepositedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var grantIdRule []interface{}
	for _, grantIdItem := range grantId {
		grantIdRule = append(grantIdRule, grantIdItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.FilterLogs(opts, "Deposited", operatorRule, grantIdRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowDepositedIterator{contract: _TokenStakingEscrow.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x73a19dd210f1a7f902193214c0ee91dd35ee5b4d920cba8d519eca65a7b488ca.
//
// Solidity: event Deposited(address indexed operator, uint256 indexed grantId, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *TokenStakingEscrowDeposited, operator []common.Address, grantId []*big.Int) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var grantIdRule []interface{}
	for _, grantIdItem := range grantId {
		grantIdRule = append(grantIdRule, grantIdItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.WatchLogs(opts, "Deposited", operatorRule, grantIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingEscrowDeposited)
				if err := _TokenStakingEscrow.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x73a19dd210f1a7f902193214c0ee91dd35ee5b4d920cba8d519eca65a7b488ca.
//
// Solidity: event Deposited(address indexed operator, uint256 indexed grantId, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) ParseDeposited(log types.Log) (*TokenStakingEscrowDeposited, error) {
	event := new(TokenStakingEscrowDeposited)
	if err := _TokenStakingEscrow.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingEscrowEscrowAuthorizedIterator is returned from FilterEscrowAuthorized and is used to iterate over the raw logs and unpacked data for EscrowAuthorized events raised by the TokenStakingEscrow contract.
type TokenStakingEscrowEscrowAuthorizedIterator struct {
	Event *TokenStakingEscrowEscrowAuthorized // Event containing the contract specifics and raw log

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
func (it *TokenStakingEscrowEscrowAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingEscrowEscrowAuthorized)
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
		it.Event = new(TokenStakingEscrowEscrowAuthorized)
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
func (it *TokenStakingEscrowEscrowAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingEscrowEscrowAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingEscrowEscrowAuthorized represents a EscrowAuthorized event raised by the TokenStakingEscrow contract.
type TokenStakingEscrowEscrowAuthorized struct {
	GrantManager common.Address
	Escrow       common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterEscrowAuthorized is a free log retrieval operation binding the contract event 0xb3cf0aa974be748a2b3911557e90014a19e73153dc123266fb73f5611d7df242.
//
// Solidity: event EscrowAuthorized(address indexed grantManager, address escrow)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) FilterEscrowAuthorized(opts *bind.FilterOpts, grantManager []common.Address) (*TokenStakingEscrowEscrowAuthorizedIterator, error) {

	var grantManagerRule []interface{}
	for _, grantManagerItem := range grantManager {
		grantManagerRule = append(grantManagerRule, grantManagerItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.FilterLogs(opts, "EscrowAuthorized", grantManagerRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowEscrowAuthorizedIterator{contract: _TokenStakingEscrow.contract, event: "EscrowAuthorized", logs: logs, sub: sub}, nil
}

// WatchEscrowAuthorized is a free log subscription operation binding the contract event 0xb3cf0aa974be748a2b3911557e90014a19e73153dc123266fb73f5611d7df242.
//
// Solidity: event EscrowAuthorized(address indexed grantManager, address escrow)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) WatchEscrowAuthorized(opts *bind.WatchOpts, sink chan<- *TokenStakingEscrowEscrowAuthorized, grantManager []common.Address) (event.Subscription, error) {

	var grantManagerRule []interface{}
	for _, grantManagerItem := range grantManager {
		grantManagerRule = append(grantManagerRule, grantManagerItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.WatchLogs(opts, "EscrowAuthorized", grantManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingEscrowEscrowAuthorized)
				if err := _TokenStakingEscrow.contract.UnpackLog(event, "EscrowAuthorized", log); err != nil {
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

// ParseEscrowAuthorized is a log parse operation binding the contract event 0xb3cf0aa974be748a2b3911557e90014a19e73153dc123266fb73f5611d7df242.
//
// Solidity: event EscrowAuthorized(address indexed grantManager, address escrow)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) ParseEscrowAuthorized(log types.Log) (*TokenStakingEscrowEscrowAuthorized, error) {
	event := new(TokenStakingEscrowEscrowAuthorized)
	if err := _TokenStakingEscrow.contract.UnpackLog(event, "EscrowAuthorized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingEscrowOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenStakingEscrow contract.
type TokenStakingEscrowOwnershipTransferredIterator struct {
	Event *TokenStakingEscrowOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenStakingEscrowOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingEscrowOwnershipTransferred)
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
		it.Event = new(TokenStakingEscrowOwnershipTransferred)
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
func (it *TokenStakingEscrowOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingEscrowOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingEscrowOwnershipTransferred represents a OwnershipTransferred event raised by the TokenStakingEscrow contract.
type TokenStakingEscrowOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenStakingEscrowOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowOwnershipTransferredIterator{contract: _TokenStakingEscrow.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenStakingEscrowOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingEscrowOwnershipTransferred)
				if err := _TokenStakingEscrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) ParseOwnershipTransferred(log types.Log) (*TokenStakingEscrowOwnershipTransferred, error) {
	event := new(TokenStakingEscrowOwnershipTransferred)
	if err := _TokenStakingEscrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenStakingEscrowRevokedDepositWithdrawnIterator is returned from FilterRevokedDepositWithdrawn and is used to iterate over the raw logs and unpacked data for RevokedDepositWithdrawn events raised by the TokenStakingEscrow contract.
type TokenStakingEscrowRevokedDepositWithdrawnIterator struct {
	Event *TokenStakingEscrowRevokedDepositWithdrawn // Event containing the contract specifics and raw log

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
func (it *TokenStakingEscrowRevokedDepositWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenStakingEscrowRevokedDepositWithdrawn)
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
		it.Event = new(TokenStakingEscrowRevokedDepositWithdrawn)
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
func (it *TokenStakingEscrowRevokedDepositWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenStakingEscrowRevokedDepositWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenStakingEscrowRevokedDepositWithdrawn represents a RevokedDepositWithdrawn event raised by the TokenStakingEscrow contract.
type TokenStakingEscrowRevokedDepositWithdrawn struct {
	Operator     common.Address
	GrantManager common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRevokedDepositWithdrawn is a free log retrieval operation binding the contract event 0x85a6a53c6a902a6ef08ddf46c4ea735fea23691d9a0651e3c1516238501b87bc.
//
// Solidity: event RevokedDepositWithdrawn(address indexed operator, address indexed grantManager, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) FilterRevokedDepositWithdrawn(opts *bind.FilterOpts, operator []common.Address, grantManager []common.Address) (*TokenStakingEscrowRevokedDepositWithdrawnIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var grantManagerRule []interface{}
	for _, grantManagerItem := range grantManager {
		grantManagerRule = append(grantManagerRule, grantManagerItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.FilterLogs(opts, "RevokedDepositWithdrawn", operatorRule, grantManagerRule)
	if err != nil {
		return nil, err
	}
	return &TokenStakingEscrowRevokedDepositWithdrawnIterator{contract: _TokenStakingEscrow.contract, event: "RevokedDepositWithdrawn", logs: logs, sub: sub}, nil
}

// WatchRevokedDepositWithdrawn is a free log subscription operation binding the contract event 0x85a6a53c6a902a6ef08ddf46c4ea735fea23691d9a0651e3c1516238501b87bc.
//
// Solidity: event RevokedDepositWithdrawn(address indexed operator, address indexed grantManager, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) WatchRevokedDepositWithdrawn(opts *bind.WatchOpts, sink chan<- *TokenStakingEscrowRevokedDepositWithdrawn, operator []common.Address, grantManager []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var grantManagerRule []interface{}
	for _, grantManagerItem := range grantManager {
		grantManagerRule = append(grantManagerRule, grantManagerItem)
	}

	logs, sub, err := _TokenStakingEscrow.contract.WatchLogs(opts, "RevokedDepositWithdrawn", operatorRule, grantManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenStakingEscrowRevokedDepositWithdrawn)
				if err := _TokenStakingEscrow.contract.UnpackLog(event, "RevokedDepositWithdrawn", log); err != nil {
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

// ParseRevokedDepositWithdrawn is a log parse operation binding the contract event 0x85a6a53c6a902a6ef08ddf46c4ea735fea23691d9a0651e3c1516238501b87bc.
//
// Solidity: event RevokedDepositWithdrawn(address indexed operator, address indexed grantManager, uint256 amount)
func (_TokenStakingEscrow *TokenStakingEscrowFilterer) ParseRevokedDepositWithdrawn(log types.Log) (*TokenStakingEscrowRevokedDepositWithdrawn, error) {
	event := new(TokenStakingEscrowRevokedDepositWithdrawn)
	if err := _TokenStakingEscrow.contract.UnpackLog(event, "RevokedDepositWithdrawn", log); err != nil {
		return nil, err
	}
	return event, nil
}
