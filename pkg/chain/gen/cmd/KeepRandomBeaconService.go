// Code generated - DO NOT EDIT.
// This file is a generated command and any manual changes will be lost.

package cmd

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/chain/gen/contract"

	"github.com/urfave/cli"
)

var KeepRandomBeaconServiceCommand cli.Command

var keepRandomBeaconServiceDescription = `The keep-random-beacon-service command allows calling the KeepRandomBeaconService contract on an
	Ethereum network. It has subcommands corresponding to each contract method,
	which respectively each take parameters based on the contract method's
	parameters.

	Subcommands will submit a non-mutating call to the network and output the
	result.

	All subcommands can be called against a specific block by passing the
	-b/--block flag.

	All subcommands can be used to investigate the result of a previous
	transaction that called that same method by passing the -t/--transaction
	flag with the transaction hash.

	Subcommands for mutating methods may be submitted as a mutating transaction
	by passing the -s/--submit flag. In this mode, this command will terminate
	successfully once the transaction has been submitted, but will not wait for
	the transaction to be included in a block. They return the transaction hash.

	Calls that require ether to be paid will get 0 ether by default, which can
	be changed by passing the -v/--value flag.`

func init() {
	AvailableCommands = append(AvailableCommands, cli.Command{
		Name:        "keep-random-beacon-service",
		Usage:       `Provides access to the KeepRandomBeaconService contract.`,
		Description: keepRandomBeaconServiceDescription,
		Subcommands: []cli.Command{{
			Name:      "dkg-contribution-margin",
			Usage:     "Calls the constant method dkgContributionMargin on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsDkgContributionMargin,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "callback-surplus-recipient",
			Usage:     "Calls the constant method callbackSurplusRecipient on the KeepRandomBeaconService contract.",
			ArgsUsage: "[requestId] ",
			Action:    krbsCallbackSurplusRecipient,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "dkg-fee-pool",
			Usage:     "Calls the constant method dkgFeePool on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsDkgFeePool,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "entry-fee-estimate",
			Usage:     "Calls the constant method entryFeeEstimate on the KeepRandomBeaconService contract.",
			ArgsUsage: "[callbackGas] ",
			Action:    krbsEntryFeeEstimate,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "initialized",
			Usage:     "Calls the constant method initialized on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsInitialized,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "version",
			Usage:     "Calls the constant method version on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsVersion,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "base-callback-gas",
			Usage:     "Calls the constant method baseCallbackGas on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsBaseCallbackGas,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "entry-fee-breakdown",
			Usage:     "Calls the constant method entryFeeBreakdown on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsEntryFeeBreakdown,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "request-subsidy-fee-pool",
			Usage:     "Calls the constant method requestSubsidyFeePool on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsRequestSubsidyFeePool,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "select-operator-contract",
			Usage:     "Calls the constant method selectOperatorContract on the KeepRandomBeaconService contract.",
			ArgsUsage: "[seed] ",
			Action:    krbsSelectOperatorContract,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "request-relay-entry0",
			Usage:     "Calls the payable method requestRelayEntry0 on the KeepRandomBeaconService contract.",
			ArgsUsage: "[callbackContract] [callbackGas] ",
			Action:    krbsRequestRelayEntry0,
			Before:    cli.BeforeFunc(cmd.PayableArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.PayableFlags,
		}, {
			Name:      "request-relay-entry",
			Usage:     "Calls the payable method requestRelayEntry on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsRequestRelayEntry,
			Before:    cli.BeforeFunc(cmd.PayableArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.PayableFlags,
		}, {
			Name:      "entry-created",
			Usage:     "Calls the method entryCreated on the KeepRandomBeaconService contract.",
			ArgsUsage: "[requestId] [entry] [submitter] ",
			Action:    krbsEntryCreated,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(3))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "execute-callback",
			Usage:     "Calls the method executeCallback on the KeepRandomBeaconService contract.",
			ArgsUsage: "[requestId] [entry] ",
			Action:    krbsExecuteCallback,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "fund-request-subsidy-fee-pool",
			Usage:     "Calls the payable method fundRequestSubsidyFeePool on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsFundRequestSubsidyFeePool,
			Before:    cli.BeforeFunc(cmd.PayableArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.PayableFlags,
		}, {
			Name:      "remove-operator-contract",
			Usage:     "Calls the method removeOperatorContract on the KeepRandomBeaconService contract.",
			ArgsUsage: "[operatorContract] ",
			Action:    krbsRemoveOperatorContract,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "add-operator-contract",
			Usage:     "Calls the method addOperatorContract on the KeepRandomBeaconService contract.",
			ArgsUsage: "[operatorContract] ",
			Action:    krbsAddOperatorContract,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "initialize",
			Usage:     "Calls the method initialize on the KeepRandomBeaconService contract.",
			ArgsUsage: "[dkgContributionMargin] [registry] ",
			Action:    krbsInitialize,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "fund-dkg-fee-pool",
			Usage:     "Calls the payable method fundDkgFeePool on the KeepRandomBeaconService contract.",
			ArgsUsage: "",
			Action:    krbsFundDkgFeePool,
			Before:    cli.BeforeFunc(cmd.PayableArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.PayableFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func krbsDkgContributionMargin(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	result, err := contract.DkgContributionMarginAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsCallbackSurplusRecipient(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}
	requestId, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter requestId, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.CallbackSurplusRecipientAtBlock(
		requestId,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsDkgFeePool(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	result, err := contract.DkgFeePoolAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsEntryFeeEstimate(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}
	callbackGas, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter callbackGas, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.EntryFeeEstimateAtBlock(
		callbackGas,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsInitialized(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	result, err := contract.InitializedAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsVersion(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	result, err := contract.VersionAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsBaseCallbackGas(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	result, err := contract.BaseCallbackGasAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsEntryFeeBreakdown(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	result, err := contract.EntryFeeBreakdownAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsRequestSubsidyFeePool(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	result, err := contract.RequestSubsidyFeePoolAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krbsSelectOperatorContract(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}
	seed, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter seed, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.SelectOperatorContractAtBlock(
		seed,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func krbsRequestRelayEntry0(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	callbackContract, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter callbackContract, a address, from passed value %v",
			c.Args()[0],
		)
	}

	callbackGas, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter callbackGas, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
		result      *big.Int
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RequestRelayEntry0(
			callbackContract,
			callbackGas,
			cmd.ValueFlagValue.Uint)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		result, err = contract.CallRequestRelayEntry0(
			callbackContract,
			callbackGas,
			cmd.ValueFlagValue.Uint, cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(result)
	}

	return nil
}

func krbsRequestRelayEntry(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
		result      *big.Int
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RequestRelayEntry(
			cmd.ValueFlagValue.Uint)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		result, err = contract.CallRequestRelayEntry(
			cmd.ValueFlagValue.Uint, cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(result)
	}

	return nil
}

func krbsEntryCreated(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	requestId, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter requestId, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	entry, err := hexutil.Decode(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter entry, a bytes, from passed value %v",
			c.Args()[1],
		)
	}

	submitter, err := ethutil.AddressFromHex(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter submitter, a address, from passed value %v",
			c.Args()[2],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.EntryCreated(
			requestId,
			entry,
			submitter,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallEntryCreated(
			requestId,
			entry,
			submitter,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krbsExecuteCallback(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	requestId, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter requestId, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	entry, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter entry, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ExecuteCallback(
			requestId,
			entry,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallExecuteCallback(
			requestId,
			entry,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krbsFundRequestSubsidyFeePool(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.FundRequestSubsidyFeePool(
			cmd.ValueFlagValue.Uint)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallFundRequestSubsidyFeePool(
			cmd.ValueFlagValue.Uint, cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krbsRemoveOperatorContract(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	operatorContract, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operatorContract, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RemoveOperatorContract(
			operatorContract,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRemoveOperatorContract(
			operatorContract,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krbsAddOperatorContract(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	operatorContract, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operatorContract, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.AddOperatorContract(
			operatorContract,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallAddOperatorContract(
			operatorContract,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krbsInitialize(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	dkgContributionMargin, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter dkgContributionMargin, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	registry, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter registry, a address, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Initialize(
			dkgContributionMargin,
			registry,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallInitialize(
			dkgContributionMargin,
			registry,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krbsFundDkgFeePool(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconService(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.FundDkgFeePool(
			cmd.ValueFlagValue.Uint)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallFundDkgFeePool(
			cmd.ValueFlagValue.Uint, cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

/// ------------------- Initialization -------------------

func initializeKeepRandomBeaconService(c *cli.Context) (*contract.KeepRandomBeaconService, error) {
	config, err := config.ReadEthereumConfig(c.GlobalString("config"))
	if err != nil {
		return nil, fmt.Errorf("error reading Ethereum config from file: [%v]", err)
	}

	client, _, _, err := ethutil.ConnectClients(config.URL, config.URLRPC)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	key, err := ethutil.DecryptKeyFile(
		config.Account.KeyFile,
		config.Account.KeyFilePassword,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read KeyFile: %s: [%v]",
			config.Account.KeyFile,
			err,
		)
	}

	checkInterval := 30 * time.Second
	maxGasPrice := big.NewInt(50000000000) // 50 Gwei
	if config.MiningCheckInterval != 0 {
		checkInterval = time.Duration(config.MiningCheckInterval) * time.Second
	}
	if config.MaxGasPrice != 0 {
		maxGasPrice = new(big.Int).SetUint64(config.MaxGasPrice)
	}

	miningWaiter := ethutil.NewMiningWaiter(client, checkInterval, maxGasPrice)

	address := common.HexToAddress(config.ContractAddresses["KeepRandomBeaconService"])

	return contract.NewKeepRandomBeaconService(
		address,
		key,
		client,
		ethutil.NewNonceManager(key.Address, client),
		miningWaiter,
		&sync.Mutex{},
	)
}
