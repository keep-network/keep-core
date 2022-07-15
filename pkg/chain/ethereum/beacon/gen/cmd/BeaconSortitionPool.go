// Code generated - DO NOT EDIT.
// This file is a generated command and any manual changes will be lost.

package cmd

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen/contract"

	"github.com/urfave/cli"
)

var BeaconSortitionPoolCommand cli.Command

var beaconSortitionPoolDescription = `The beacon-sortition-pool command allows calling the BeaconSortitionPool contract on an
	ETH-like network. It has subcommands corresponding to each contract method,
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
		Name:        "beacon-sortition-pool",
		Usage:       `Provides access to the BeaconSortitionPool contract.`,
		Description: beaconSortitionPoolDescription,
		Subcommands: []cli.Command{{
			Name:      "can-restore-reward-eligibility",
			Usage:     "Calls the view method canRestoreRewardEligibility on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspCanRestoreRewardEligibility,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-available-rewards",
			Usage:     "Calls the view method getAvailableRewards on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspGetAvailableRewards,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-operator-i-d",
			Usage:     "Calls the view method getOperatorID on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspGetOperatorID,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-pool-weight",
			Usage:     "Calls the view method getPoolWeight on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspGetPoolWeight,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "ineligible-earned-rewards",
			Usage:     "Calls the view method ineligibleEarnedRewards on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspIneligibleEarnedRewards,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-eligible-for-rewards",
			Usage:     "Calls the view method isEligibleForRewards on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspIsEligibleForRewards,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-locked",
			Usage:     "Calls the view method isLocked on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspIsLocked,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-operator-in-pool",
			Usage:     "Calls the view method isOperatorInPool on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspIsOperatorInPool,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-operator-registered",
			Usage:     "Calls the view method isOperatorRegistered on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspIsOperatorRegistered,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-operator-up-to-date",
			Usage:     "Calls the view method isOperatorUpToDate on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] [arg_authorizedStake] ",
			Action:    bspIsOperatorUpToDate,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "operators-in-pool",
			Usage:     "Calls the view method operatorsInPool on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspOperatorsInPool,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "owner",
			Usage:     "Calls the view method owner on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspOwner,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "pool-weight-divisor",
			Usage:     "Calls the view method poolWeightDivisor on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspPoolWeightDivisor,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "reward-token",
			Usage:     "Calls the view method rewardToken on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspRewardToken,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "rewards-eligibility-restorable-at",
			Usage:     "Calls the view method rewardsEligibilityRestorableAt on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspRewardsEligibilityRestorableAt,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "total-weight",
			Usage:     "Calls the view method totalWeight on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspTotalWeight,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "insert-operator",
			Usage:     "Calls the nonpayable method insertOperator on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] [arg_authorizedStake] ",
			Action:    bspInsertOperator,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "lock",
			Usage:     "Calls the nonpayable method lock on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspLock,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "receive-approval",
			Usage:     "Calls the nonpayable method receiveApproval on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_sender] [arg_amount] [arg_token] [arg3] ",
			Action:    bspReceiveApproval,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(4))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "renounce-ownership",
			Usage:     "Calls the nonpayable method renounceOwnership on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspRenounceOwnership,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "restore-reward-eligibility",
			Usage:     "Calls the nonpayable method restoreRewardEligibility on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    bspRestoreRewardEligibility,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "transfer-ownership",
			Usage:     "Calls the nonpayable method transferOwnership on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_newOwner] ",
			Action:    bspTransferOwnership,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "unlock",
			Usage:     "Calls the nonpayable method unlock on the BeaconSortitionPool contract.",
			ArgsUsage: "",
			Action:    bspUnlock,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-operator-status",
			Usage:     "Calls the nonpayable method updateOperatorStatus on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] [arg_authorizedStake] ",
			Action:    bspUpdateOperatorStatus,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw-ineligible",
			Usage:     "Calls the nonpayable method withdrawIneligible on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_recipient] ",
			Action:    bspWithdrawIneligible,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw-rewards",
			Usage:     "Calls the nonpayable method withdrawRewards on the BeaconSortitionPool contract.",
			ArgsUsage: "[arg_operator] [arg_beneficiary] ",
			Action:    bspWithdrawRewards,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func bspCanRestoreRewardEligibility(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.CanRestoreRewardEligibilityAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspGetAvailableRewards(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetAvailableRewardsAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspGetOperatorID(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetOperatorIDAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspGetPoolWeight(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetPoolWeightAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIneligibleEarnedRewards(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.IneligibleEarnedRewardsAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsEligibleForRewards(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.IsEligibleForRewardsAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsLocked(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.IsLockedAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsOperatorInPool(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.IsOperatorInPoolAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsOperatorRegistered(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.IsOperatorRegisteredAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsOperatorUpToDate(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	arg_authorizedStake, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_authorizedStake, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.IsOperatorUpToDateAtBlock(
		arg_operator,
		arg_authorizedStake,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspOperatorsInPool(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.OperatorsInPoolAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspOwner(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.OwnerAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspPoolWeightDivisor(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.PoolWeightDivisorAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspRewardToken(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.RewardTokenAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspRewardsEligibilityRestorableAt(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.RewardsEligibilityRestorableAtAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspTotalWeight(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.TotalWeightAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func bspInsertOperator(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	arg_authorizedStake, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_authorizedStake, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.InsertOperator(
			arg_operator,
			arg_authorizedStake,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallInsertOperator(
			arg_operator,
			arg_authorizedStake,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspLock(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Lock()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLock(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspReceiveApproval(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_sender, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_sender, a address, from passed value %v",
			c.Args()[0],
		)
	}

	arg_amount, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_amount, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	arg_token, err := chainutil.AddressFromHex(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_token, a address, from passed value %v",
			c.Args()[2],
		)
	}

	arg3, err := hexutil.Decode(c.Args()[3])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg3, a bytes, from passed value %v",
			c.Args()[3],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ReceiveApproval(
			arg_sender,
			arg_amount,
			arg_token,
			arg3,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallReceiveApproval(
			arg_sender,
			arg_amount,
			arg_token,
			arg3,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspRenounceOwnership(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RenounceOwnership()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRenounceOwnership(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspRestoreRewardEligibility(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RestoreRewardEligibility(
			arg_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRestoreRewardEligibility(
			arg_operator,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspTransferOwnership(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_newOwner, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_newOwner, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.TransferOwnership(
			arg_newOwner,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallTransferOwnership(
			arg_newOwner,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspUnlock(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Unlock()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUnlock(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspUpdateOperatorStatus(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	arg_authorizedStake, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_authorizedStake, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateOperatorStatus(
			arg_operator,
			arg_authorizedStake,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateOperatorStatus(
			arg_operator,
			arg_authorizedStake,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspWithdrawIneligible(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_recipient, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_recipient, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.WithdrawIneligible(
			arg_recipient,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallWithdrawIneligible(
			arg_recipient,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bspWithdrawRewards(c *cli.Context) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_operator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	arg_beneficiary, err := chainutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_beneficiary, a address, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
		result      *big.Int
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.WithdrawRewards(
			arg_operator,
			arg_beneficiary,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		result, err = contract.CallWithdrawRewards(
			arg_operator,
			arg_beneficiary,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(result)
	}

	return nil
}

/// ------------------- Initialization -------------------

func initializeBeaconSortitionPool(c *cli.Context) (*contract.BeaconSortitionPool, error) {
	config, err := config.ReadEthereumConfig(c.GlobalString("config"))
	if err != nil {
		return nil, fmt.Errorf("error reading config from file: [%v]", err)
	}

	client, _, _, err := chainutil.ConnectClients(config.URL, config.URLRPC)
	if err != nil {
		return nil, fmt.Errorf("error connecting to host chain node: [%v]", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve host chain id: [%v]",
			err,
		)
	}

	key, err := chainutil.DecryptKeyFile(
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

	miningWaiter := chainutil.NewMiningWaiter(client, config)

	blockCounter, err := chainutil.NewBlockCounter(client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create block counter: [%v]",
			err,
		)
	}

	address := common.HexToAddress(config.ContractAddresses["BeaconSortitionPool"])

	return contract.NewBeaconSortitionPool(
		address,
		chainID,
		key,
		client,
		chainutil.NewNonceManager(client, key.Address),
		miningWaiter,
		blockCounter,
		&sync.Mutex{},
	)
}
