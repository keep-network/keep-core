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

// BitcoinTxInfo2 is an auto generated low-level Go binding around an user-defined struct.
type BitcoinTxInfo2 struct {
	Version      [4]byte
	InputVector  []byte
	OutputVector []byte
	Locktime     [4]byte
}

// WalletCoordinatorDepositExtraInfo is an auto generated low-level Go binding around an user-defined struct.
type WalletCoordinatorDepositExtraInfo struct {
	FundingTx        BitcoinTxInfo2
	BlindingFactor   [8]byte
	WalletPubKeyHash [20]byte
	RefundPubKeyHash [20]byte
	RefundLocktime   [4]byte
}

// WalletCoordinatorDepositKey is an auto generated low-level Go binding around an user-defined struct.
type WalletCoordinatorDepositKey struct {
	FundingTxHash      [32]byte
	FundingOutputIndex uint32
}

// WalletCoordinatorDepositSweepProposal is an auto generated low-level Go binding around an user-defined struct.
type WalletCoordinatorDepositSweepProposal struct {
	WalletPubKeyHash [20]byte
	DepositsKeys     []WalletCoordinatorDepositKey
	SweepTxFee       *big.Int
}

// WalletCoordinatorMetaData contains all meta data concerning the WalletCoordinator contract.
var WalletCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"CoordinatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"CoordinatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"depositSweepProposalValidity\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"depositMinAge\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"depositRefundSafetyMargin\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"depositSweepMaxSize\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"depositSweepProposalSubmissionGasOffset\",\"type\":\"uint32\"}],\"name\":\"DepositSweepProposalParametersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fundingTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fundingOutputIndex\",\"type\":\"uint32\"}],\"internalType\":\"structWalletCoordinator.DepositKey[]\",\"name\":\"depositsKeys\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sweepTxFee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structWalletCoordinator.DepositSweepProposal\",\"name\":\"proposal\",\"type\":\"tuple\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"DepositSweepProposalSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"heartbeatRequestValidity\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"heartbeatRequestGasOffset\",\"type\":\"uint32\"}],\"name\":\"HeartbeatRequestParametersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"HeartbeatRequestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newReimbursementPool\",\"type\":\"address\"}],\"name\":\"ReimbursementPoolUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"}],\"name\":\"WalletManuallyUnlocked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"addCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"internalType\":\"contractBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositMinAge\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositRefundSafetyMargin\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositSweepMaxSize\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositSweepProposalSubmissionGasOffset\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositSweepProposalValidity\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"heartbeatRequestGasOffset\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"heartbeatRequestValidity\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractBridge\",\"name\":\"_bridge\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isCoordinator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reimbursementPool\",\"outputs\":[{\"internalType\":\"contractReimbursementPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"removeCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"requestHeartbeat\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"requestHeartbeatWithReimbursement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fundingTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fundingOutputIndex\",\"type\":\"uint32\"}],\"internalType\":\"structWalletCoordinator.DepositKey[]\",\"name\":\"depositsKeys\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sweepTxFee\",\"type\":\"uint256\"}],\"internalType\":\"structWalletCoordinator.DepositSweepProposal\",\"name\":\"proposal\",\"type\":\"tuple\"}],\"name\":\"submitDepositSweepProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fundingTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fundingOutputIndex\",\"type\":\"uint32\"}],\"internalType\":\"structWalletCoordinator.DepositKey[]\",\"name\":\"depositsKeys\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sweepTxFee\",\"type\":\"uint256\"}],\"internalType\":\"structWalletCoordinator.DepositSweepProposal\",\"name\":\"proposal\",\"type\":\"tuple\"}],\"name\":\"submitDepositSweepProposalWithReimbursement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"}],\"name\":\"unlockWallet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_depositSweepProposalValidity\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_depositMinAge\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_depositRefundSafetyMargin\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_depositSweepMaxSize\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_depositSweepProposalSubmissionGasOffset\",\"type\":\"uint32\"}],\"name\":\"updateDepositSweepProposalParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_heartbeatRequestValidity\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_heartbeatRequestGasOffset\",\"type\":\"uint32\"}],\"name\":\"updateHeartbeatRequestParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractReimbursementPool\",\"name\":\"_reimbursementPool\",\"type\":\"address\"}],\"name\":\"updateReimbursementPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fundingTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fundingOutputIndex\",\"type\":\"uint32\"}],\"internalType\":\"structWalletCoordinator.DepositKey[]\",\"name\":\"depositsKeys\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sweepTxFee\",\"type\":\"uint256\"}],\"internalType\":\"structWalletCoordinator.DepositSweepProposal\",\"name\":\"proposal\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes4\",\"name\":\"version\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"inputVector\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"outputVector\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"locktime\",\"type\":\"bytes4\"}],\"internalType\":\"structBitcoinTx.Info\",\"name\":\"fundingTx\",\"type\":\"tuple\"},{\"internalType\":\"bytes8\",\"name\":\"blindingFactor\",\"type\":\"bytes8\"},{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"internalType\":\"bytes20\",\"name\":\"refundPubKeyHash\",\"type\":\"bytes20\"},{\"internalType\":\"bytes4\",\"name\":\"refundLocktime\",\"type\":\"bytes4\"}],\"internalType\":\"structWalletCoordinator.DepositExtraInfo[]\",\"name\":\"depositsExtraInfo\",\"type\":\"tuple[]\"}],\"name\":\"validateDepositSweepProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes20\",\"name\":\"\",\"type\":\"bytes20\"}],\"name\":\"walletLock\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"expiresAt\",\"type\":\"uint32\"},{\"internalType\":\"enumWalletCoordinator.WalletAction\",\"name\":\"cause\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// WalletCoordinatorABI is the input ABI used to generate the binding from.
// Deprecated: Use WalletCoordinatorMetaData.ABI instead.
var WalletCoordinatorABI = WalletCoordinatorMetaData.ABI

// WalletCoordinator is an auto generated Go binding around an Ethereum contract.
type WalletCoordinator struct {
	WalletCoordinatorCaller     // Read-only binding to the contract
	WalletCoordinatorTransactor // Write-only binding to the contract
	WalletCoordinatorFilterer   // Log filterer for contract events
}

// WalletCoordinatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type WalletCoordinatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletCoordinatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WalletCoordinatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletCoordinatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WalletCoordinatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletCoordinatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WalletCoordinatorSession struct {
	Contract     *WalletCoordinator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// WalletCoordinatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WalletCoordinatorCallerSession struct {
	Contract *WalletCoordinatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// WalletCoordinatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WalletCoordinatorTransactorSession struct {
	Contract     *WalletCoordinatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// WalletCoordinatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type WalletCoordinatorRaw struct {
	Contract *WalletCoordinator // Generic contract binding to access the raw methods on
}

// WalletCoordinatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WalletCoordinatorCallerRaw struct {
	Contract *WalletCoordinatorCaller // Generic read-only contract binding to access the raw methods on
}

// WalletCoordinatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WalletCoordinatorTransactorRaw struct {
	Contract *WalletCoordinatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWalletCoordinator creates a new instance of WalletCoordinator, bound to a specific deployed contract.
func NewWalletCoordinator(address common.Address, backend bind.ContractBackend) (*WalletCoordinator, error) {
	contract, err := bindWalletCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinator{WalletCoordinatorCaller: WalletCoordinatorCaller{contract: contract}, WalletCoordinatorTransactor: WalletCoordinatorTransactor{contract: contract}, WalletCoordinatorFilterer: WalletCoordinatorFilterer{contract: contract}}, nil
}

// NewWalletCoordinatorCaller creates a new read-only instance of WalletCoordinator, bound to a specific deployed contract.
func NewWalletCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*WalletCoordinatorCaller, error) {
	contract, err := bindWalletCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorCaller{contract: contract}, nil
}

// NewWalletCoordinatorTransactor creates a new write-only instance of WalletCoordinator, bound to a specific deployed contract.
func NewWalletCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*WalletCoordinatorTransactor, error) {
	contract, err := bindWalletCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorTransactor{contract: contract}, nil
}

// NewWalletCoordinatorFilterer creates a new log filterer instance of WalletCoordinator, bound to a specific deployed contract.
func NewWalletCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*WalletCoordinatorFilterer, error) {
	contract, err := bindWalletCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorFilterer{contract: contract}, nil
}

// bindWalletCoordinator binds a generic wrapper to an already deployed contract.
func bindWalletCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WalletCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WalletCoordinator *WalletCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WalletCoordinator.Contract.WalletCoordinatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WalletCoordinator *WalletCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.WalletCoordinatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WalletCoordinator *WalletCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.WalletCoordinatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WalletCoordinator *WalletCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WalletCoordinator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WalletCoordinator *WalletCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WalletCoordinator *WalletCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.contract.Transact(opts, method, params...)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_WalletCoordinator *WalletCoordinatorSession) Bridge() (common.Address, error) {
	return _WalletCoordinator.Contract.Bridge(&_WalletCoordinator.CallOpts)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCallerSession) Bridge() (common.Address, error) {
	return _WalletCoordinator.Contract.Bridge(&_WalletCoordinator.CallOpts)
}

// DepositMinAge is a free data retrieval call binding the contract method 0xbc20a080.
//
// Solidity: function depositMinAge() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositMinAge(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositMinAge")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DepositMinAge is a free data retrieval call binding the contract method 0xbc20a080.
//
// Solidity: function depositMinAge() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) DepositMinAge() (uint32, error) {
	return _WalletCoordinator.Contract.DepositMinAge(&_WalletCoordinator.CallOpts)
}

// DepositMinAge is a free data retrieval call binding the contract method 0xbc20a080.
//
// Solidity: function depositMinAge() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositMinAge() (uint32, error) {
	return _WalletCoordinator.Contract.DepositMinAge(&_WalletCoordinator.CallOpts)
}

// DepositRefundSafetyMargin is a free data retrieval call binding the contract method 0x31d8c622.
//
// Solidity: function depositRefundSafetyMargin() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositRefundSafetyMargin(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositRefundSafetyMargin")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DepositRefundSafetyMargin is a free data retrieval call binding the contract method 0x31d8c622.
//
// Solidity: function depositRefundSafetyMargin() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) DepositRefundSafetyMargin() (uint32, error) {
	return _WalletCoordinator.Contract.DepositRefundSafetyMargin(&_WalletCoordinator.CallOpts)
}

// DepositRefundSafetyMargin is a free data retrieval call binding the contract method 0x31d8c622.
//
// Solidity: function depositRefundSafetyMargin() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositRefundSafetyMargin() (uint32, error) {
	return _WalletCoordinator.Contract.DepositRefundSafetyMargin(&_WalletCoordinator.CallOpts)
}

// DepositSweepMaxSize is a free data retrieval call binding the contract method 0xacaac705.
//
// Solidity: function depositSweepMaxSize() view returns(uint16)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositSweepMaxSize(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositSweepMaxSize")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// DepositSweepMaxSize is a free data retrieval call binding the contract method 0xacaac705.
//
// Solidity: function depositSweepMaxSize() view returns(uint16)
func (_WalletCoordinator *WalletCoordinatorSession) DepositSweepMaxSize() (uint16, error) {
	return _WalletCoordinator.Contract.DepositSweepMaxSize(&_WalletCoordinator.CallOpts)
}

// DepositSweepMaxSize is a free data retrieval call binding the contract method 0xacaac705.
//
// Solidity: function depositSweepMaxSize() view returns(uint16)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositSweepMaxSize() (uint16, error) {
	return _WalletCoordinator.Contract.DepositSweepMaxSize(&_WalletCoordinator.CallOpts)
}

// DepositSweepProposalSubmissionGasOffset is a free data retrieval call binding the contract method 0x08553292.
//
// Solidity: function depositSweepProposalSubmissionGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositSweepProposalSubmissionGasOffset(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositSweepProposalSubmissionGasOffset")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DepositSweepProposalSubmissionGasOffset is a free data retrieval call binding the contract method 0x08553292.
//
// Solidity: function depositSweepProposalSubmissionGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) DepositSweepProposalSubmissionGasOffset() (uint32, error) {
	return _WalletCoordinator.Contract.DepositSweepProposalSubmissionGasOffset(&_WalletCoordinator.CallOpts)
}

// DepositSweepProposalSubmissionGasOffset is a free data retrieval call binding the contract method 0x08553292.
//
// Solidity: function depositSweepProposalSubmissionGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositSweepProposalSubmissionGasOffset() (uint32, error) {
	return _WalletCoordinator.Contract.DepositSweepProposalSubmissionGasOffset(&_WalletCoordinator.CallOpts)
}

// DepositSweepProposalValidity is a free data retrieval call binding the contract method 0xf74b37bd.
//
// Solidity: function depositSweepProposalValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositSweepProposalValidity(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositSweepProposalValidity")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DepositSweepProposalValidity is a free data retrieval call binding the contract method 0xf74b37bd.
//
// Solidity: function depositSweepProposalValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) DepositSweepProposalValidity() (uint32, error) {
	return _WalletCoordinator.Contract.DepositSweepProposalValidity(&_WalletCoordinator.CallOpts)
}

// DepositSweepProposalValidity is a free data retrieval call binding the contract method 0xf74b37bd.
//
// Solidity: function depositSweepProposalValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositSweepProposalValidity() (uint32, error) {
	return _WalletCoordinator.Contract.DepositSweepProposalValidity(&_WalletCoordinator.CallOpts)
}

// HeartbeatRequestGasOffset is a free data retrieval call binding the contract method 0xa34fff59.
//
// Solidity: function heartbeatRequestGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) HeartbeatRequestGasOffset(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "heartbeatRequestGasOffset")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// HeartbeatRequestGasOffset is a free data retrieval call binding the contract method 0xa34fff59.
//
// Solidity: function heartbeatRequestGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) HeartbeatRequestGasOffset() (uint32, error) {
	return _WalletCoordinator.Contract.HeartbeatRequestGasOffset(&_WalletCoordinator.CallOpts)
}

// HeartbeatRequestGasOffset is a free data retrieval call binding the contract method 0xa34fff59.
//
// Solidity: function heartbeatRequestGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) HeartbeatRequestGasOffset() (uint32, error) {
	return _WalletCoordinator.Contract.HeartbeatRequestGasOffset(&_WalletCoordinator.CallOpts)
}

// HeartbeatRequestValidity is a free data retrieval call binding the contract method 0x7fbf9b39.
//
// Solidity: function heartbeatRequestValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) HeartbeatRequestValidity(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "heartbeatRequestValidity")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// HeartbeatRequestValidity is a free data retrieval call binding the contract method 0x7fbf9b39.
//
// Solidity: function heartbeatRequestValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) HeartbeatRequestValidity() (uint32, error) {
	return _WalletCoordinator.Contract.HeartbeatRequestValidity(&_WalletCoordinator.CallOpts)
}

// HeartbeatRequestValidity is a free data retrieval call binding the contract method 0x7fbf9b39.
//
// Solidity: function heartbeatRequestValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) HeartbeatRequestValidity() (uint32, error) {
	return _WalletCoordinator.Contract.HeartbeatRequestValidity(&_WalletCoordinator.CallOpts)
}

// IsCoordinator is a free data retrieval call binding the contract method 0xaec32099.
//
// Solidity: function isCoordinator(address ) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorCaller) IsCoordinator(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "isCoordinator", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsCoordinator is a free data retrieval call binding the contract method 0xaec32099.
//
// Solidity: function isCoordinator(address ) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorSession) IsCoordinator(arg0 common.Address) (bool, error) {
	return _WalletCoordinator.Contract.IsCoordinator(&_WalletCoordinator.CallOpts, arg0)
}

// IsCoordinator is a free data retrieval call binding the contract method 0xaec32099.
//
// Solidity: function isCoordinator(address ) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorCallerSession) IsCoordinator(arg0 common.Address) (bool, error) {
	return _WalletCoordinator.Contract.IsCoordinator(&_WalletCoordinator.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WalletCoordinator *WalletCoordinatorSession) Owner() (common.Address, error) {
	return _WalletCoordinator.Contract.Owner(&_WalletCoordinator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCallerSession) Owner() (common.Address, error) {
	return _WalletCoordinator.Contract.Owner(&_WalletCoordinator.CallOpts)
}

// ReimbursementPool is a free data retrieval call binding the contract method 0xc09975cd.
//
// Solidity: function reimbursementPool() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCaller) ReimbursementPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "reimbursementPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReimbursementPool is a free data retrieval call binding the contract method 0xc09975cd.
//
// Solidity: function reimbursementPool() view returns(address)
func (_WalletCoordinator *WalletCoordinatorSession) ReimbursementPool() (common.Address, error) {
	return _WalletCoordinator.Contract.ReimbursementPool(&_WalletCoordinator.CallOpts)
}

// ReimbursementPool is a free data retrieval call binding the contract method 0xc09975cd.
//
// Solidity: function reimbursementPool() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCallerSession) ReimbursementPool() (common.Address, error) {
	return _WalletCoordinator.Contract.ReimbursementPool(&_WalletCoordinator.CallOpts)
}

// ValidateDepositSweepProposal is a free data retrieval call binding the contract method 0x7698fa97.
//
// Solidity: function validateDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, ((bytes4,bytes,bytes,bytes4),bytes8,bytes20,bytes20,bytes4)[] depositsExtraInfo) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorCaller) ValidateDepositSweepProposal(opts *bind.CallOpts, proposal WalletCoordinatorDepositSweepProposal, depositsExtraInfo []WalletCoordinatorDepositExtraInfo) (bool, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "validateDepositSweepProposal", proposal, depositsExtraInfo)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidateDepositSweepProposal is a free data retrieval call binding the contract method 0x7698fa97.
//
// Solidity: function validateDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, ((bytes4,bytes,bytes,bytes4),bytes8,bytes20,bytes20,bytes4)[] depositsExtraInfo) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorSession) ValidateDepositSweepProposal(proposal WalletCoordinatorDepositSweepProposal, depositsExtraInfo []WalletCoordinatorDepositExtraInfo) (bool, error) {
	return _WalletCoordinator.Contract.ValidateDepositSweepProposal(&_WalletCoordinator.CallOpts, proposal, depositsExtraInfo)
}

// ValidateDepositSweepProposal is a free data retrieval call binding the contract method 0x7698fa97.
//
// Solidity: function validateDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, ((bytes4,bytes,bytes,bytes4),bytes8,bytes20,bytes20,bytes4)[] depositsExtraInfo) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorCallerSession) ValidateDepositSweepProposal(proposal WalletCoordinatorDepositSweepProposal, depositsExtraInfo []WalletCoordinatorDepositExtraInfo) (bool, error) {
	return _WalletCoordinator.Contract.ValidateDepositSweepProposal(&_WalletCoordinator.CallOpts, proposal, depositsExtraInfo)
}

// WalletLock is a free data retrieval call binding the contract method 0x2c259d2b.
//
// Solidity: function walletLock(bytes20 ) view returns(uint32 expiresAt, uint8 cause)
func (_WalletCoordinator *WalletCoordinatorCaller) WalletLock(opts *bind.CallOpts, arg0 [20]byte) (struct {
	ExpiresAt uint32
	Cause     uint8
}, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "walletLock", arg0)

	outstruct := new(struct {
		ExpiresAt uint32
		Cause     uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ExpiresAt = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Cause = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// WalletLock is a free data retrieval call binding the contract method 0x2c259d2b.
//
// Solidity: function walletLock(bytes20 ) view returns(uint32 expiresAt, uint8 cause)
func (_WalletCoordinator *WalletCoordinatorSession) WalletLock(arg0 [20]byte) (struct {
	ExpiresAt uint32
	Cause     uint8
}, error) {
	return _WalletCoordinator.Contract.WalletLock(&_WalletCoordinator.CallOpts, arg0)
}

// WalletLock is a free data retrieval call binding the contract method 0x2c259d2b.
//
// Solidity: function walletLock(bytes20 ) view returns(uint32 expiresAt, uint8 cause)
func (_WalletCoordinator *WalletCoordinatorCallerSession) WalletLock(arg0 [20]byte) (struct {
	ExpiresAt uint32
	Cause     uint8
}, error) {
	return _WalletCoordinator.Contract.WalletLock(&_WalletCoordinator.CallOpts, arg0)
}

// AddCoordinator is a paid mutator transaction binding the contract method 0xfbd6d77e.
//
// Solidity: function addCoordinator(address coordinator) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) AddCoordinator(opts *bind.TransactOpts, coordinator common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "addCoordinator", coordinator)
}

// AddCoordinator is a paid mutator transaction binding the contract method 0xfbd6d77e.
//
// Solidity: function addCoordinator(address coordinator) returns()
func (_WalletCoordinator *WalletCoordinatorSession) AddCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.AddCoordinator(&_WalletCoordinator.TransactOpts, coordinator)
}

// AddCoordinator is a paid mutator transaction binding the contract method 0xfbd6d77e.
//
// Solidity: function addCoordinator(address coordinator) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) AddCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.AddCoordinator(&_WalletCoordinator.TransactOpts, coordinator)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bridge) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) Initialize(opts *bind.TransactOpts, _bridge common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "initialize", _bridge)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bridge) returns()
func (_WalletCoordinator *WalletCoordinatorSession) Initialize(_bridge common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.Initialize(&_WalletCoordinator.TransactOpts, _bridge)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bridge) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) Initialize(_bridge common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.Initialize(&_WalletCoordinator.TransactOpts, _bridge)
}

// RemoveCoordinator is a paid mutator transaction binding the contract method 0x8cb5a0c0.
//
// Solidity: function removeCoordinator(address coordinator) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) RemoveCoordinator(opts *bind.TransactOpts, coordinator common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "removeCoordinator", coordinator)
}

// RemoveCoordinator is a paid mutator transaction binding the contract method 0x8cb5a0c0.
//
// Solidity: function removeCoordinator(address coordinator) returns()
func (_WalletCoordinator *WalletCoordinatorSession) RemoveCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RemoveCoordinator(&_WalletCoordinator.TransactOpts, coordinator)
}

// RemoveCoordinator is a paid mutator transaction binding the contract method 0x8cb5a0c0.
//
// Solidity: function removeCoordinator(address coordinator) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) RemoveCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RemoveCoordinator(&_WalletCoordinator.TransactOpts, coordinator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WalletCoordinator *WalletCoordinatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RenounceOwnership(&_WalletCoordinator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RenounceOwnership(&_WalletCoordinator.TransactOpts)
}

// RequestHeartbeat is a paid mutator transaction binding the contract method 0x6348c439.
//
// Solidity: function requestHeartbeat(bytes20 walletPubKeyHash, bytes message) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) RequestHeartbeat(opts *bind.TransactOpts, walletPubKeyHash [20]byte, message []byte) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "requestHeartbeat", walletPubKeyHash, message)
}

// RequestHeartbeat is a paid mutator transaction binding the contract method 0x6348c439.
//
// Solidity: function requestHeartbeat(bytes20 walletPubKeyHash, bytes message) returns()
func (_WalletCoordinator *WalletCoordinatorSession) RequestHeartbeat(walletPubKeyHash [20]byte, message []byte) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RequestHeartbeat(&_WalletCoordinator.TransactOpts, walletPubKeyHash, message)
}

// RequestHeartbeat is a paid mutator transaction binding the contract method 0x6348c439.
//
// Solidity: function requestHeartbeat(bytes20 walletPubKeyHash, bytes message) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) RequestHeartbeat(walletPubKeyHash [20]byte, message []byte) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RequestHeartbeat(&_WalletCoordinator.TransactOpts, walletPubKeyHash, message)
}

// RequestHeartbeatWithReimbursement is a paid mutator transaction binding the contract method 0x34173c2d.
//
// Solidity: function requestHeartbeatWithReimbursement(bytes20 walletPubKeyHash, bytes message) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) RequestHeartbeatWithReimbursement(opts *bind.TransactOpts, walletPubKeyHash [20]byte, message []byte) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "requestHeartbeatWithReimbursement", walletPubKeyHash, message)
}

// RequestHeartbeatWithReimbursement is a paid mutator transaction binding the contract method 0x34173c2d.
//
// Solidity: function requestHeartbeatWithReimbursement(bytes20 walletPubKeyHash, bytes message) returns()
func (_WalletCoordinator *WalletCoordinatorSession) RequestHeartbeatWithReimbursement(walletPubKeyHash [20]byte, message []byte) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RequestHeartbeatWithReimbursement(&_WalletCoordinator.TransactOpts, walletPubKeyHash, message)
}

// RequestHeartbeatWithReimbursement is a paid mutator transaction binding the contract method 0x34173c2d.
//
// Solidity: function requestHeartbeatWithReimbursement(bytes20 walletPubKeyHash, bytes message) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) RequestHeartbeatWithReimbursement(walletPubKeyHash [20]byte, message []byte) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RequestHeartbeatWithReimbursement(&_WalletCoordinator.TransactOpts, walletPubKeyHash, message)
}

// SubmitDepositSweepProposal is a paid mutator transaction binding the contract method 0x742b386b.
//
// Solidity: function submitDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) SubmitDepositSweepProposal(opts *bind.TransactOpts, proposal WalletCoordinatorDepositSweepProposal) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "submitDepositSweepProposal", proposal)
}

// SubmitDepositSweepProposal is a paid mutator transaction binding the contract method 0x742b386b.
//
// Solidity: function submitDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal) returns()
func (_WalletCoordinator *WalletCoordinatorSession) SubmitDepositSweepProposal(proposal WalletCoordinatorDepositSweepProposal) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.SubmitDepositSweepProposal(&_WalletCoordinator.TransactOpts, proposal)
}

// SubmitDepositSweepProposal is a paid mutator transaction binding the contract method 0x742b386b.
//
// Solidity: function submitDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) SubmitDepositSweepProposal(proposal WalletCoordinatorDepositSweepProposal) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.SubmitDepositSweepProposal(&_WalletCoordinator.TransactOpts, proposal)
}

// SubmitDepositSweepProposalWithReimbursement is a paid mutator transaction binding the contract method 0x6ed6405f.
//
// Solidity: function submitDepositSweepProposalWithReimbursement((bytes20,(bytes32,uint32)[],uint256) proposal) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) SubmitDepositSweepProposalWithReimbursement(opts *bind.TransactOpts, proposal WalletCoordinatorDepositSweepProposal) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "submitDepositSweepProposalWithReimbursement", proposal)
}

// SubmitDepositSweepProposalWithReimbursement is a paid mutator transaction binding the contract method 0x6ed6405f.
//
// Solidity: function submitDepositSweepProposalWithReimbursement((bytes20,(bytes32,uint32)[],uint256) proposal) returns()
func (_WalletCoordinator *WalletCoordinatorSession) SubmitDepositSweepProposalWithReimbursement(proposal WalletCoordinatorDepositSweepProposal) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.SubmitDepositSweepProposalWithReimbursement(&_WalletCoordinator.TransactOpts, proposal)
}

// SubmitDepositSweepProposalWithReimbursement is a paid mutator transaction binding the contract method 0x6ed6405f.
//
// Solidity: function submitDepositSweepProposalWithReimbursement((bytes20,(bytes32,uint32)[],uint256) proposal) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) SubmitDepositSweepProposalWithReimbursement(proposal WalletCoordinatorDepositSweepProposal) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.SubmitDepositSweepProposalWithReimbursement(&_WalletCoordinator.TransactOpts, proposal)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WalletCoordinator *WalletCoordinatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.TransferOwnership(&_WalletCoordinator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.TransferOwnership(&_WalletCoordinator.TransactOpts, newOwner)
}

// UnlockWallet is a paid mutator transaction binding the contract method 0x9039bc1c.
//
// Solidity: function unlockWallet(bytes20 walletPubKeyHash) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UnlockWallet(opts *bind.TransactOpts, walletPubKeyHash [20]byte) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "unlockWallet", walletPubKeyHash)
}

// UnlockWallet is a paid mutator transaction binding the contract method 0x9039bc1c.
//
// Solidity: function unlockWallet(bytes20 walletPubKeyHash) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UnlockWallet(walletPubKeyHash [20]byte) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UnlockWallet(&_WalletCoordinator.TransactOpts, walletPubKeyHash)
}

// UnlockWallet is a paid mutator transaction binding the contract method 0x9039bc1c.
//
// Solidity: function unlockWallet(bytes20 walletPubKeyHash) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UnlockWallet(walletPubKeyHash [20]byte) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UnlockWallet(&_WalletCoordinator.TransactOpts, walletPubKeyHash)
}

// UpdateDepositSweepProposalParameters is a paid mutator transaction binding the contract method 0x5486032e.
//
// Solidity: function updateDepositSweepProposalParameters(uint32 _depositSweepProposalValidity, uint32 _depositMinAge, uint32 _depositRefundSafetyMargin, uint16 _depositSweepMaxSize, uint32 _depositSweepProposalSubmissionGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateDepositSweepProposalParameters(opts *bind.TransactOpts, _depositSweepProposalValidity uint32, _depositMinAge uint32, _depositRefundSafetyMargin uint32, _depositSweepMaxSize uint16, _depositSweepProposalSubmissionGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateDepositSweepProposalParameters", _depositSweepProposalValidity, _depositMinAge, _depositRefundSafetyMargin, _depositSweepMaxSize, _depositSweepProposalSubmissionGasOffset)
}

// UpdateDepositSweepProposalParameters is a paid mutator transaction binding the contract method 0x5486032e.
//
// Solidity: function updateDepositSweepProposalParameters(uint32 _depositSweepProposalValidity, uint32 _depositMinAge, uint32 _depositRefundSafetyMargin, uint16 _depositSweepMaxSize, uint32 _depositSweepProposalSubmissionGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateDepositSweepProposalParameters(_depositSweepProposalValidity uint32, _depositMinAge uint32, _depositRefundSafetyMargin uint32, _depositSweepMaxSize uint16, _depositSweepProposalSubmissionGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositSweepProposalParameters(&_WalletCoordinator.TransactOpts, _depositSweepProposalValidity, _depositMinAge, _depositRefundSafetyMargin, _depositSweepMaxSize, _depositSweepProposalSubmissionGasOffset)
}

// UpdateDepositSweepProposalParameters is a paid mutator transaction binding the contract method 0x5486032e.
//
// Solidity: function updateDepositSweepProposalParameters(uint32 _depositSweepProposalValidity, uint32 _depositMinAge, uint32 _depositRefundSafetyMargin, uint16 _depositSweepMaxSize, uint32 _depositSweepProposalSubmissionGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateDepositSweepProposalParameters(_depositSweepProposalValidity uint32, _depositMinAge uint32, _depositRefundSafetyMargin uint32, _depositSweepMaxSize uint16, _depositSweepProposalSubmissionGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositSweepProposalParameters(&_WalletCoordinator.TransactOpts, _depositSweepProposalValidity, _depositMinAge, _depositRefundSafetyMargin, _depositSweepMaxSize, _depositSweepProposalSubmissionGasOffset)
}

// UpdateHeartbeatRequestParameters is a paid mutator transaction binding the contract method 0x85a0924c.
//
// Solidity: function updateHeartbeatRequestParameters(uint32 _heartbeatRequestValidity, uint32 _heartbeatRequestGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateHeartbeatRequestParameters(opts *bind.TransactOpts, _heartbeatRequestValidity uint32, _heartbeatRequestGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateHeartbeatRequestParameters", _heartbeatRequestValidity, _heartbeatRequestGasOffset)
}

// UpdateHeartbeatRequestParameters is a paid mutator transaction binding the contract method 0x85a0924c.
//
// Solidity: function updateHeartbeatRequestParameters(uint32 _heartbeatRequestValidity, uint32 _heartbeatRequestGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateHeartbeatRequestParameters(_heartbeatRequestValidity uint32, _heartbeatRequestGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateHeartbeatRequestParameters(&_WalletCoordinator.TransactOpts, _heartbeatRequestValidity, _heartbeatRequestGasOffset)
}

// UpdateHeartbeatRequestParameters is a paid mutator transaction binding the contract method 0x85a0924c.
//
// Solidity: function updateHeartbeatRequestParameters(uint32 _heartbeatRequestValidity, uint32 _heartbeatRequestGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateHeartbeatRequestParameters(_heartbeatRequestValidity uint32, _heartbeatRequestGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateHeartbeatRequestParameters(&_WalletCoordinator.TransactOpts, _heartbeatRequestValidity, _heartbeatRequestGasOffset)
}

// UpdateReimbursementPool is a paid mutator transaction binding the contract method 0x7b35b4e6.
//
// Solidity: function updateReimbursementPool(address _reimbursementPool) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateReimbursementPool(opts *bind.TransactOpts, _reimbursementPool common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateReimbursementPool", _reimbursementPool)
}

// UpdateReimbursementPool is a paid mutator transaction binding the contract method 0x7b35b4e6.
//
// Solidity: function updateReimbursementPool(address _reimbursementPool) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateReimbursementPool(_reimbursementPool common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateReimbursementPool(&_WalletCoordinator.TransactOpts, _reimbursementPool)
}

// UpdateReimbursementPool is a paid mutator transaction binding the contract method 0x7b35b4e6.
//
// Solidity: function updateReimbursementPool(address _reimbursementPool) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateReimbursementPool(_reimbursementPool common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateReimbursementPool(&_WalletCoordinator.TransactOpts, _reimbursementPool)
}

// WalletCoordinatorCoordinatorAddedIterator is returned from FilterCoordinatorAdded and is used to iterate over the raw logs and unpacked data for CoordinatorAdded events raised by the WalletCoordinator contract.
type WalletCoordinatorCoordinatorAddedIterator struct {
	Event *WalletCoordinatorCoordinatorAdded // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorCoordinatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorCoordinatorAdded)
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
		it.Event = new(WalletCoordinatorCoordinatorAdded)
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
func (it *WalletCoordinatorCoordinatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorCoordinatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorCoordinatorAdded represents a CoordinatorAdded event raised by the WalletCoordinator contract.
type WalletCoordinatorCoordinatorAdded struct {
	Coordinator common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCoordinatorAdded is a free log retrieval operation binding the contract event 0x56e6c5206d2baab7a01569b19d84b5963ca9f96e10f55076b5e5987fbf780cf5.
//
// Solidity: event CoordinatorAdded(address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterCoordinatorAdded(opts *bind.FilterOpts, coordinator []common.Address) (*WalletCoordinatorCoordinatorAddedIterator, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "CoordinatorAdded", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorCoordinatorAddedIterator{contract: _WalletCoordinator.contract, event: "CoordinatorAdded", logs: logs, sub: sub}, nil
}

// WatchCoordinatorAdded is a free log subscription operation binding the contract event 0x56e6c5206d2baab7a01569b19d84b5963ca9f96e10f55076b5e5987fbf780cf5.
//
// Solidity: event CoordinatorAdded(address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchCoordinatorAdded(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorCoordinatorAdded, coordinator []common.Address) (event.Subscription, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "CoordinatorAdded", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorCoordinatorAdded)
				if err := _WalletCoordinator.contract.UnpackLog(event, "CoordinatorAdded", log); err != nil {
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

// ParseCoordinatorAdded is a log parse operation binding the contract event 0x56e6c5206d2baab7a01569b19d84b5963ca9f96e10f55076b5e5987fbf780cf5.
//
// Solidity: event CoordinatorAdded(address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseCoordinatorAdded(log types.Log) (*WalletCoordinatorCoordinatorAdded, error) {
	event := new(WalletCoordinatorCoordinatorAdded)
	if err := _WalletCoordinator.contract.UnpackLog(event, "CoordinatorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorCoordinatorRemovedIterator is returned from FilterCoordinatorRemoved and is used to iterate over the raw logs and unpacked data for CoordinatorRemoved events raised by the WalletCoordinator contract.
type WalletCoordinatorCoordinatorRemovedIterator struct {
	Event *WalletCoordinatorCoordinatorRemoved // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorCoordinatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorCoordinatorRemoved)
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
		it.Event = new(WalletCoordinatorCoordinatorRemoved)
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
func (it *WalletCoordinatorCoordinatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorCoordinatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorCoordinatorRemoved represents a CoordinatorRemoved event raised by the WalletCoordinator contract.
type WalletCoordinatorCoordinatorRemoved struct {
	Coordinator common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCoordinatorRemoved is a free log retrieval operation binding the contract event 0xcf15831020102bbbe0547b0c3460aa49715f2646e0911760d2b069e0fb2f62a1.
//
// Solidity: event CoordinatorRemoved(address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterCoordinatorRemoved(opts *bind.FilterOpts, coordinator []common.Address) (*WalletCoordinatorCoordinatorRemovedIterator, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "CoordinatorRemoved", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorCoordinatorRemovedIterator{contract: _WalletCoordinator.contract, event: "CoordinatorRemoved", logs: logs, sub: sub}, nil
}

// WatchCoordinatorRemoved is a free log subscription operation binding the contract event 0xcf15831020102bbbe0547b0c3460aa49715f2646e0911760d2b069e0fb2f62a1.
//
// Solidity: event CoordinatorRemoved(address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchCoordinatorRemoved(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorCoordinatorRemoved, coordinator []common.Address) (event.Subscription, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "CoordinatorRemoved", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorCoordinatorRemoved)
				if err := _WalletCoordinator.contract.UnpackLog(event, "CoordinatorRemoved", log); err != nil {
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

// ParseCoordinatorRemoved is a log parse operation binding the contract event 0xcf15831020102bbbe0547b0c3460aa49715f2646e0911760d2b069e0fb2f62a1.
//
// Solidity: event CoordinatorRemoved(address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseCoordinatorRemoved(log types.Log) (*WalletCoordinatorCoordinatorRemoved, error) {
	event := new(WalletCoordinatorCoordinatorRemoved)
	if err := _WalletCoordinator.contract.UnpackLog(event, "CoordinatorRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorDepositSweepProposalParametersUpdatedIterator is returned from FilterDepositSweepProposalParametersUpdated and is used to iterate over the raw logs and unpacked data for DepositSweepProposalParametersUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalParametersUpdatedIterator struct {
	Event *WalletCoordinatorDepositSweepProposalParametersUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorDepositSweepProposalParametersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorDepositSweepProposalParametersUpdated)
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
		it.Event = new(WalletCoordinatorDepositSweepProposalParametersUpdated)
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
func (it *WalletCoordinatorDepositSweepProposalParametersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorDepositSweepProposalParametersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorDepositSweepProposalParametersUpdated represents a DepositSweepProposalParametersUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalParametersUpdated struct {
	DepositSweepProposalValidity            uint32
	DepositMinAge                           uint32
	DepositRefundSafetyMargin               uint32
	DepositSweepMaxSize                     uint16
	DepositSweepProposalSubmissionGasOffset uint32
	Raw                                     types.Log // Blockchain specific contextual infos
}

// FilterDepositSweepProposalParametersUpdated is a free log retrieval operation binding the contract event 0xc36fe6b9d74f864b5e345edc3435eb451a80669ec78c8357cc15e3a1d3897a73.
//
// Solidity: event DepositSweepProposalParametersUpdated(uint32 depositSweepProposalValidity, uint32 depositMinAge, uint32 depositRefundSafetyMargin, uint16 depositSweepMaxSize, uint32 depositSweepProposalSubmissionGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterDepositSweepProposalParametersUpdated(opts *bind.FilterOpts) (*WalletCoordinatorDepositSweepProposalParametersUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "DepositSweepProposalParametersUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorDepositSweepProposalParametersUpdatedIterator{contract: _WalletCoordinator.contract, event: "DepositSweepProposalParametersUpdated", logs: logs, sub: sub}, nil
}

// WatchDepositSweepProposalParametersUpdated is a free log subscription operation binding the contract event 0xc36fe6b9d74f864b5e345edc3435eb451a80669ec78c8357cc15e3a1d3897a73.
//
// Solidity: event DepositSweepProposalParametersUpdated(uint32 depositSweepProposalValidity, uint32 depositMinAge, uint32 depositRefundSafetyMargin, uint16 depositSweepMaxSize, uint32 depositSweepProposalSubmissionGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchDepositSweepProposalParametersUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorDepositSweepProposalParametersUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "DepositSweepProposalParametersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorDepositSweepProposalParametersUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalParametersUpdated", log); err != nil {
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

// ParseDepositSweepProposalParametersUpdated is a log parse operation binding the contract event 0xc36fe6b9d74f864b5e345edc3435eb451a80669ec78c8357cc15e3a1d3897a73.
//
// Solidity: event DepositSweepProposalParametersUpdated(uint32 depositSweepProposalValidity, uint32 depositMinAge, uint32 depositRefundSafetyMargin, uint16 depositSweepMaxSize, uint32 depositSweepProposalSubmissionGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseDepositSweepProposalParametersUpdated(log types.Log) (*WalletCoordinatorDepositSweepProposalParametersUpdated, error) {
	event := new(WalletCoordinatorDepositSweepProposalParametersUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalParametersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorDepositSweepProposalSubmittedIterator is returned from FilterDepositSweepProposalSubmitted and is used to iterate over the raw logs and unpacked data for DepositSweepProposalSubmitted events raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalSubmittedIterator struct {
	Event *WalletCoordinatorDepositSweepProposalSubmitted // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorDepositSweepProposalSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorDepositSweepProposalSubmitted)
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
		it.Event = new(WalletCoordinatorDepositSweepProposalSubmitted)
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
func (it *WalletCoordinatorDepositSweepProposalSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorDepositSweepProposalSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorDepositSweepProposalSubmitted represents a DepositSweepProposalSubmitted event raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalSubmitted struct {
	Proposal    WalletCoordinatorDepositSweepProposal
	Coordinator common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDepositSweepProposalSubmitted is a free log retrieval operation binding the contract event 0xe17b4ed4828b9f018ca549cfd5f662e5551408093b1db16d94b9cd242e6df171.
//
// Solidity: event DepositSweepProposalSubmitted((bytes20,(bytes32,uint32)[],uint256) proposal, address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterDepositSweepProposalSubmitted(opts *bind.FilterOpts, coordinator []common.Address) (*WalletCoordinatorDepositSweepProposalSubmittedIterator, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "DepositSweepProposalSubmitted", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorDepositSweepProposalSubmittedIterator{contract: _WalletCoordinator.contract, event: "DepositSweepProposalSubmitted", logs: logs, sub: sub}, nil
}

// WatchDepositSweepProposalSubmitted is a free log subscription operation binding the contract event 0xe17b4ed4828b9f018ca549cfd5f662e5551408093b1db16d94b9cd242e6df171.
//
// Solidity: event DepositSweepProposalSubmitted((bytes20,(bytes32,uint32)[],uint256) proposal, address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchDepositSweepProposalSubmitted(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorDepositSweepProposalSubmitted, coordinator []common.Address) (event.Subscription, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "DepositSweepProposalSubmitted", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorDepositSweepProposalSubmitted)
				if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalSubmitted", log); err != nil {
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

// ParseDepositSweepProposalSubmitted is a log parse operation binding the contract event 0xe17b4ed4828b9f018ca549cfd5f662e5551408093b1db16d94b9cd242e6df171.
//
// Solidity: event DepositSweepProposalSubmitted((bytes20,(bytes32,uint32)[],uint256) proposal, address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseDepositSweepProposalSubmitted(log types.Log) (*WalletCoordinatorDepositSweepProposalSubmitted, error) {
	event := new(WalletCoordinatorDepositSweepProposalSubmitted)
	if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorHeartbeatRequestParametersUpdatedIterator is returned from FilterHeartbeatRequestParametersUpdated and is used to iterate over the raw logs and unpacked data for HeartbeatRequestParametersUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorHeartbeatRequestParametersUpdatedIterator struct {
	Event *WalletCoordinatorHeartbeatRequestParametersUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorHeartbeatRequestParametersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorHeartbeatRequestParametersUpdated)
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
		it.Event = new(WalletCoordinatorHeartbeatRequestParametersUpdated)
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
func (it *WalletCoordinatorHeartbeatRequestParametersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorHeartbeatRequestParametersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorHeartbeatRequestParametersUpdated represents a HeartbeatRequestParametersUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorHeartbeatRequestParametersUpdated struct {
	HeartbeatRequestValidity  uint32
	HeartbeatRequestGasOffset uint32
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterHeartbeatRequestParametersUpdated is a free log retrieval operation binding the contract event 0xe67c0a182190bf283264290d7208533d5936483c98baee174eeac036b081360a.
//
// Solidity: event HeartbeatRequestParametersUpdated(uint32 heartbeatRequestValidity, uint32 heartbeatRequestGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterHeartbeatRequestParametersUpdated(opts *bind.FilterOpts) (*WalletCoordinatorHeartbeatRequestParametersUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "HeartbeatRequestParametersUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorHeartbeatRequestParametersUpdatedIterator{contract: _WalletCoordinator.contract, event: "HeartbeatRequestParametersUpdated", logs: logs, sub: sub}, nil
}

// WatchHeartbeatRequestParametersUpdated is a free log subscription operation binding the contract event 0xe67c0a182190bf283264290d7208533d5936483c98baee174eeac036b081360a.
//
// Solidity: event HeartbeatRequestParametersUpdated(uint32 heartbeatRequestValidity, uint32 heartbeatRequestGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchHeartbeatRequestParametersUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorHeartbeatRequestParametersUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "HeartbeatRequestParametersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorHeartbeatRequestParametersUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "HeartbeatRequestParametersUpdated", log); err != nil {
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

// ParseHeartbeatRequestParametersUpdated is a log parse operation binding the contract event 0xe67c0a182190bf283264290d7208533d5936483c98baee174eeac036b081360a.
//
// Solidity: event HeartbeatRequestParametersUpdated(uint32 heartbeatRequestValidity, uint32 heartbeatRequestGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseHeartbeatRequestParametersUpdated(log types.Log) (*WalletCoordinatorHeartbeatRequestParametersUpdated, error) {
	event := new(WalletCoordinatorHeartbeatRequestParametersUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "HeartbeatRequestParametersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorHeartbeatRequestSubmittedIterator is returned from FilterHeartbeatRequestSubmitted and is used to iterate over the raw logs and unpacked data for HeartbeatRequestSubmitted events raised by the WalletCoordinator contract.
type WalletCoordinatorHeartbeatRequestSubmittedIterator struct {
	Event *WalletCoordinatorHeartbeatRequestSubmitted // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorHeartbeatRequestSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorHeartbeatRequestSubmitted)
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
		it.Event = new(WalletCoordinatorHeartbeatRequestSubmitted)
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
func (it *WalletCoordinatorHeartbeatRequestSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorHeartbeatRequestSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorHeartbeatRequestSubmitted represents a HeartbeatRequestSubmitted event raised by the WalletCoordinator contract.
type WalletCoordinatorHeartbeatRequestSubmitted struct {
	WalletPubKeyHash [20]byte
	Message          []byte
	Coordinator      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterHeartbeatRequestSubmitted is a free log retrieval operation binding the contract event 0x46f4f1fc7a1c5c364372f2e4398c3657d5e5cc35ea286f9972459376b29b1384.
//
// Solidity: event HeartbeatRequestSubmitted(bytes20 walletPubKeyHash, bytes message, address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterHeartbeatRequestSubmitted(opts *bind.FilterOpts, coordinator []common.Address) (*WalletCoordinatorHeartbeatRequestSubmittedIterator, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "HeartbeatRequestSubmitted", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorHeartbeatRequestSubmittedIterator{contract: _WalletCoordinator.contract, event: "HeartbeatRequestSubmitted", logs: logs, sub: sub}, nil
}

// WatchHeartbeatRequestSubmitted is a free log subscription operation binding the contract event 0x46f4f1fc7a1c5c364372f2e4398c3657d5e5cc35ea286f9972459376b29b1384.
//
// Solidity: event HeartbeatRequestSubmitted(bytes20 walletPubKeyHash, bytes message, address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchHeartbeatRequestSubmitted(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorHeartbeatRequestSubmitted, coordinator []common.Address) (event.Subscription, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "HeartbeatRequestSubmitted", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorHeartbeatRequestSubmitted)
				if err := _WalletCoordinator.contract.UnpackLog(event, "HeartbeatRequestSubmitted", log); err != nil {
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

// ParseHeartbeatRequestSubmitted is a log parse operation binding the contract event 0x46f4f1fc7a1c5c364372f2e4398c3657d5e5cc35ea286f9972459376b29b1384.
//
// Solidity: event HeartbeatRequestSubmitted(bytes20 walletPubKeyHash, bytes message, address indexed coordinator)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseHeartbeatRequestSubmitted(log types.Log) (*WalletCoordinatorHeartbeatRequestSubmitted, error) {
	event := new(WalletCoordinatorHeartbeatRequestSubmitted)
	if err := _WalletCoordinator.contract.UnpackLog(event, "HeartbeatRequestSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the WalletCoordinator contract.
type WalletCoordinatorInitializedIterator struct {
	Event *WalletCoordinatorInitialized // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorInitialized)
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
		it.Event = new(WalletCoordinatorInitialized)
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
func (it *WalletCoordinatorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorInitialized represents a Initialized event raised by the WalletCoordinator contract.
type WalletCoordinatorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterInitialized(opts *bind.FilterOpts) (*WalletCoordinatorInitializedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorInitializedIterator{contract: _WalletCoordinator.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorInitialized) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorInitialized)
				if err := _WalletCoordinator.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseInitialized(log types.Log) (*WalletCoordinatorInitialized, error) {
	event := new(WalletCoordinatorInitialized)
	if err := _WalletCoordinator.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the WalletCoordinator contract.
type WalletCoordinatorOwnershipTransferredIterator struct {
	Event *WalletCoordinatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorOwnershipTransferred)
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
		it.Event = new(WalletCoordinatorOwnershipTransferred)
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
func (it *WalletCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorOwnershipTransferred represents a OwnershipTransferred event raised by the WalletCoordinator contract.
type WalletCoordinatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WalletCoordinatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorOwnershipTransferredIterator{contract: _WalletCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorOwnershipTransferred)
				if err := _WalletCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*WalletCoordinatorOwnershipTransferred, error) {
	event := new(WalletCoordinatorOwnershipTransferred)
	if err := _WalletCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorReimbursementPoolUpdatedIterator is returned from FilterReimbursementPoolUpdated and is used to iterate over the raw logs and unpacked data for ReimbursementPoolUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorReimbursementPoolUpdatedIterator struct {
	Event *WalletCoordinatorReimbursementPoolUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorReimbursementPoolUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorReimbursementPoolUpdated)
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
		it.Event = new(WalletCoordinatorReimbursementPoolUpdated)
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
func (it *WalletCoordinatorReimbursementPoolUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorReimbursementPoolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorReimbursementPoolUpdated represents a ReimbursementPoolUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorReimbursementPoolUpdated struct {
	NewReimbursementPool common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterReimbursementPoolUpdated is a free log retrieval operation binding the contract event 0x0e2d2343d31b085b7c4e56d1c8a6ec79f7ab07460386f1c9a1756239fe2533ac.
//
// Solidity: event ReimbursementPoolUpdated(address newReimbursementPool)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterReimbursementPoolUpdated(opts *bind.FilterOpts) (*WalletCoordinatorReimbursementPoolUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "ReimbursementPoolUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorReimbursementPoolUpdatedIterator{contract: _WalletCoordinator.contract, event: "ReimbursementPoolUpdated", logs: logs, sub: sub}, nil
}

// WatchReimbursementPoolUpdated is a free log subscription operation binding the contract event 0x0e2d2343d31b085b7c4e56d1c8a6ec79f7ab07460386f1c9a1756239fe2533ac.
//
// Solidity: event ReimbursementPoolUpdated(address newReimbursementPool)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchReimbursementPoolUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorReimbursementPoolUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "ReimbursementPoolUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorReimbursementPoolUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "ReimbursementPoolUpdated", log); err != nil {
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

// ParseReimbursementPoolUpdated is a log parse operation binding the contract event 0x0e2d2343d31b085b7c4e56d1c8a6ec79f7ab07460386f1c9a1756239fe2533ac.
//
// Solidity: event ReimbursementPoolUpdated(address newReimbursementPool)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseReimbursementPoolUpdated(log types.Log) (*WalletCoordinatorReimbursementPoolUpdated, error) {
	event := new(WalletCoordinatorReimbursementPoolUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "ReimbursementPoolUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorWalletManuallyUnlockedIterator is returned from FilterWalletManuallyUnlocked and is used to iterate over the raw logs and unpacked data for WalletManuallyUnlocked events raised by the WalletCoordinator contract.
type WalletCoordinatorWalletManuallyUnlockedIterator struct {
	Event *WalletCoordinatorWalletManuallyUnlocked // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorWalletManuallyUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorWalletManuallyUnlocked)
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
		it.Event = new(WalletCoordinatorWalletManuallyUnlocked)
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
func (it *WalletCoordinatorWalletManuallyUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorWalletManuallyUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorWalletManuallyUnlocked represents a WalletManuallyUnlocked event raised by the WalletCoordinator contract.
type WalletCoordinatorWalletManuallyUnlocked struct {
	WalletPubKeyHash [20]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterWalletManuallyUnlocked is a free log retrieval operation binding the contract event 0x5ea2f0397ee52fe8e2e6460e92b334c996f3580dc76c022d41b290cb985f41e6.
//
// Solidity: event WalletManuallyUnlocked(bytes20 indexed walletPubKeyHash)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterWalletManuallyUnlocked(opts *bind.FilterOpts, walletPubKeyHash [][20]byte) (*WalletCoordinatorWalletManuallyUnlockedIterator, error) {

	var walletPubKeyHashRule []interface{}
	for _, walletPubKeyHashItem := range walletPubKeyHash {
		walletPubKeyHashRule = append(walletPubKeyHashRule, walletPubKeyHashItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "WalletManuallyUnlocked", walletPubKeyHashRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorWalletManuallyUnlockedIterator{contract: _WalletCoordinator.contract, event: "WalletManuallyUnlocked", logs: logs, sub: sub}, nil
}

// WatchWalletManuallyUnlocked is a free log subscription operation binding the contract event 0x5ea2f0397ee52fe8e2e6460e92b334c996f3580dc76c022d41b290cb985f41e6.
//
// Solidity: event WalletManuallyUnlocked(bytes20 indexed walletPubKeyHash)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchWalletManuallyUnlocked(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorWalletManuallyUnlocked, walletPubKeyHash [][20]byte) (event.Subscription, error) {

	var walletPubKeyHashRule []interface{}
	for _, walletPubKeyHashItem := range walletPubKeyHash {
		walletPubKeyHashRule = append(walletPubKeyHashRule, walletPubKeyHashItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "WalletManuallyUnlocked", walletPubKeyHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorWalletManuallyUnlocked)
				if err := _WalletCoordinator.contract.UnpackLog(event, "WalletManuallyUnlocked", log); err != nil {
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

// ParseWalletManuallyUnlocked is a log parse operation binding the contract event 0x5ea2f0397ee52fe8e2e6460e92b334c996f3580dc76c022d41b290cb985f41e6.
//
// Solidity: event WalletManuallyUnlocked(bytes20 indexed walletPubKeyHash)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseWalletManuallyUnlocked(log types.Log) (*WalletCoordinatorWalletManuallyUnlocked, error) {
	event := new(WalletCoordinatorWalletManuallyUnlocked)
	if err := _WalletCoordinator.contract.UnpackLog(event, "WalletManuallyUnlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
