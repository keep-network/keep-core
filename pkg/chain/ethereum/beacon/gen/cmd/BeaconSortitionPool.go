// Code generated - DO NOT EDIT.
// This file is a generated command and any manual changes will be lost.

package cmd

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen/contract"

	"github.com/spf13/cobra"
)

var BeaconSortitionPoolCommand *cobra.Command

var beaconSortitionPoolDescription = `The beacon-sortition-pool command allows calling the BeaconSortitionPool contract on an
	Ethereum network. It has subcommands corresponding to each contract method,
	which respectively each take parameters based on the contract method's
	parameters.

	Subcommands will submit a non-mutating call to the network and output the
	result.

	All subcommands can be called against a specific block by passing the
	-b/--block flag.

	Subcommands for mutating methods may be submitted as a mutating transaction
	by passing the -s/--submit flag. In this mode, this command will terminate
	successfully once the transaction has been submitted, but will not wait for
	the transaction to be included in a block. They return the transaction hash.

	Calls that require ether to be paid will get 0 ether by default, which can
	be changed by passing the -v/--value flag.`

func init() {
	BeaconSortitionPoolCommand := &cobra.Command{
		Use:   "beacon-sortition-pool",
		Short: `Provides access to the BeaconSortitionPool contract.`,
		Long:  beaconSortitionPoolDescription,
	}

	BeaconSortitionPoolCommand.AddCommand(
		bspCanRestoreRewardEligibilityCommand(),
		bspGetAvailableRewardsCommand(),
		bspGetOperatorIDCommand(),
		bspGetPoolWeightCommand(),
		bspIneligibleEarnedRewardsCommand(),
		bspIsEligibleForRewardsCommand(),
		bspIsLockedCommand(),
		bspIsOperatorInPoolCommand(),
		bspIsOperatorRegisteredCommand(),
		bspIsOperatorUpToDateCommand(),
		bspOperatorsInPoolCommand(),
		bspOwnerCommand(),
		bspPoolWeightDivisorCommand(),
		bspRewardTokenCommand(),
		bspRewardsEligibilityRestorableAtCommand(),
		bspTotalWeightCommand(),
		bspInsertOperatorCommand(),
		bspLockCommand(),
		bspReceiveApprovalCommand(),
		bspRenounceOwnershipCommand(),
		bspRestoreRewardEligibilityCommand(),
		bspTransferOwnershipCommand(),
		bspUnlockCommand(),
		bspUpdateOperatorStatusCommand(),
		bspWithdrawIneligibleCommand(),
		bspWithdrawRewardsCommand(),
	)

	ModuleCommand.AddCommand(BeaconSortitionPoolCommand)
}

/// ------------------- Const methods -------------------

func bspCanRestoreRewardEligibilityCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "can-restore-reward-eligibility [arg_operator]",
		Short:                 "Calls the view method canRestoreRewardEligibility on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspCanRestoreRewardEligibility,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspCanRestoreRewardEligibility(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.CanRestoreRewardEligibilityAtBlock(
		arg_operator,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspGetAvailableRewardsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "get-available-rewards [arg_operator]",
		Short:                 "Calls the view method getAvailableRewards on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspGetAvailableRewards,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspGetAvailableRewards(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.GetAvailableRewardsAtBlock(
		arg_operator,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspGetOperatorIDCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "get-operator-i-d [arg_operator]",
		Short:                 "Calls the view method getOperatorID on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspGetOperatorID,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspGetOperatorID(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.GetOperatorIDAtBlock(
		arg_operator,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspGetPoolWeightCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "get-pool-weight [arg_operator]",
		Short:                 "Calls the view method getPoolWeight on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspGetPoolWeight,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspGetPoolWeight(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.GetPoolWeightAtBlock(
		arg_operator,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIneligibleEarnedRewardsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "ineligible-earned-rewards",
		Short:                 "Calls the view method ineligibleEarnedRewards on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspIneligibleEarnedRewards,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspIneligibleEarnedRewards(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.IneligibleEarnedRewardsAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsEligibleForRewardsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "is-eligible-for-rewards [arg_operator]",
		Short:                 "Calls the view method isEligibleForRewards on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspIsEligibleForRewards,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspIsEligibleForRewards(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.IsEligibleForRewardsAtBlock(
		arg_operator,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsLockedCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "is-locked",
		Short:                 "Calls the view method isLocked on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspIsLocked,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspIsLocked(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.IsLockedAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsOperatorInPoolCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "is-operator-in-pool [arg_operator]",
		Short:                 "Calls the view method isOperatorInPool on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspIsOperatorInPool,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspIsOperatorInPool(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.IsOperatorInPoolAtBlock(
		arg_operator,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsOperatorRegisteredCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "is-operator-registered [arg_operator]",
		Short:                 "Calls the view method isOperatorRegistered on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspIsOperatorRegistered,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspIsOperatorRegistered(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.IsOperatorRegisteredAtBlock(
		arg_operator,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspIsOperatorUpToDateCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "is-operator-up-to-date [arg_operator] [arg_authorizedStake]",
		Short:                 "Calls the view method isOperatorUpToDate on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(2),
		RunE:                  bspIsOperatorUpToDate,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspIsOperatorUpToDate(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	arg_authorizedStake, err := hexutil.DecodeBig(args[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_authorizedStake, a uint256, from passed value %v",
			args[1],
		)
	}

	result, err := contract.IsOperatorUpToDateAtBlock(
		arg_operator,
		arg_authorizedStake,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspOperatorsInPoolCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "operators-in-pool",
		Short:                 "Calls the view method operatorsInPool on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspOperatorsInPool,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspOperatorsInPool(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.OperatorsInPoolAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspOwnerCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "owner",
		Short:                 "Calls the view method owner on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspOwner,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspOwner(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.OwnerAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspPoolWeightDivisorCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "pool-weight-divisor",
		Short:                 "Calls the view method poolWeightDivisor on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspPoolWeightDivisor,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspPoolWeightDivisor(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.PoolWeightDivisorAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspRewardTokenCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "reward-token",
		Short:                 "Calls the view method rewardToken on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspRewardToken,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspRewardToken(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.RewardTokenAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspRewardsEligibilityRestorableAtCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "rewards-eligibility-restorable-at [arg_operator]",
		Short:                 "Calls the view method rewardsEligibilityRestorableAt on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspRewardsEligibilityRestorableAt,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspRewardsEligibilityRestorableAt(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}
	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.RewardsEligibilityRestorableAtAtBlock(
		arg_operator,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bspTotalWeightCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "total-weight",
		Short:                 "Calls the view method totalWeight on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspTotalWeight,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bspTotalWeight(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	result, err := contract.TotalWeightAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func bspInsertOperatorCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "insert-operator [arg_operator] [arg_authorizedStake]",
		Short:                 "Calls the nonpayable method insertOperator on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(2),
		RunE:                  bspInsertOperator,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspInsertOperator(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	arg_authorizedStake, err := hexutil.DecodeBig(args[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_authorizedStake, a uint256, from passed value %v",
			args[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.InsertOperator(
			arg_operator,
			arg_authorizedStake,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallInsertOperator(
			arg_operator,
			arg_authorizedStake,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspLockCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "lock",
		Short:                 "Calls the nonpayable method lock on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspLock,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspLock(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Lock()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallLock(
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspReceiveApprovalCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "receive-approval [arg_sender] [arg_amount] [arg_token] [arg3]",
		Short:                 "Calls the nonpayable method receiveApproval on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(4),
		RunE:                  bspReceiveApproval,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspReceiveApproval(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_sender, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_sender, a address, from passed value %v",
			args[0],
		)
	}

	arg_amount, err := hexutil.DecodeBig(args[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_amount, a uint256, from passed value %v",
			args[1],
		)
	}

	arg_token, err := chainutil.AddressFromHex(args[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_token, a address, from passed value %v",
			args[2],
		)
	}

	arg3, err := hexutil.Decode(args[3])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg3, a bytes, from passed value %v",
			args[3],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
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

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallReceiveApproval(
			arg_sender,
			arg_amount,
			arg_token,
			arg3,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspRenounceOwnershipCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "renounce-ownership",
		Short:                 "Calls the nonpayable method renounceOwnership on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspRenounceOwnership,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspRenounceOwnership(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RenounceOwnership()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallRenounceOwnership(
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspRestoreRewardEligibilityCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "restore-reward-eligibility [arg_operator]",
		Short:                 "Calls the nonpayable method restoreRewardEligibility on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspRestoreRewardEligibility,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspRestoreRewardEligibility(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RestoreRewardEligibility(
			arg_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallRestoreRewardEligibility(
			arg_operator,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspTransferOwnershipCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "transfer-ownership [arg_newOwner]",
		Short:                 "Calls the nonpayable method transferOwnership on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspTransferOwnership,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspTransferOwnership(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_newOwner, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_newOwner, a address, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.TransferOwnership(
			arg_newOwner,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallTransferOwnership(
			arg_newOwner,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspUnlockCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "unlock",
		Short:                 "Calls the nonpayable method unlock on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bspUnlock,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspUnlock(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Unlock()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallUnlock(
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspUpdateOperatorStatusCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "update-operator-status [arg_operator] [arg_authorizedStake]",
		Short:                 "Calls the nonpayable method updateOperatorStatus on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(2),
		RunE:                  bspUpdateOperatorStatus,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspUpdateOperatorStatus(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	arg_authorizedStake, err := hexutil.DecodeBig(args[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_authorizedStake, a uint256, from passed value %v",
			args[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateOperatorStatus(
			arg_operator,
			arg_authorizedStake,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallUpdateOperatorStatus(
			arg_operator,
			arg_authorizedStake,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspWithdrawIneligibleCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "withdraw-ineligible [arg_recipient]",
		Short:                 "Calls the nonpayable method withdrawIneligible on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bspWithdrawIneligible,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspWithdrawIneligible(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_recipient, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_recipient, a address, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.WithdrawIneligible(
			arg_recipient,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallWithdrawIneligible(
			arg_recipient,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bspWithdrawRewardsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "withdraw-rewards [arg_operator] [arg_beneficiary]",
		Short:                 "Calls the nonpayable method withdrawRewards on the BeaconSortitionPool contract.",
		Args:                  cmd.ArgCountChecker(2),
		RunE:                  bspWithdrawRewards,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bspWithdrawRewards(c *cobra.Command, args []string) error {
	contract, err := initializeBeaconSortitionPool(c)
	if err != nil {
		return err
	}

	arg_operator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_operator, a address, from passed value %v",
			args[0],
		)
	}

	arg_beneficiary, err := chainutil.AddressFromHex(args[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_beneficiary, a address, from passed value %v",
			args[1],
		)
	}

	var (
		transaction *types.Transaction
		result      *big.Int
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.WithdrawRewards(
			arg_operator,
			arg_beneficiary,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		result, err = contract.CallWithdrawRewards(
			arg_operator,
			arg_beneficiary,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(result)
	}

	return nil
}

/// ------------------- Initialization -------------------

func initializeBeaconSortitionPool(c *cobra.Command) (*contract.BeaconSortitionPool, error) {
	cfg := *ModuleCommand.GetConfig()

	client, err := ethclient.Dial(cfg.URL)
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
		cfg.Account.KeyFile,
		cfg.Account.KeyFilePassword,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read KeyFile: %s: [%v]",
			cfg.Account.KeyFile,
			err,
		)
	}

	miningWaiter := chainutil.NewMiningWaiter(client, cfg)

	blockCounter, err := chainutil.NewBlockCounter(client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create block counter: [%v]",
			err,
		)
	}

	address, err := cfg.ContractAddress("BeaconSortitionPool")
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get %s address: [%w]",
			"BeaconSortitionPool",
			err,
		)
	}

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
