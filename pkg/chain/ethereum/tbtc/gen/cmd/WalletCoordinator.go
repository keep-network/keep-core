// Code generated - DO NOT EDIT.
// This file is a generated command and any manual changes will be lost.

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd"
	"github.com/keep-network/keep-common/pkg/utils/decode"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/abi"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/contract"

	"github.com/spf13/cobra"
)

var WalletCoordinatorCommand *cobra.Command

var walletCoordinatorDescription = `The wallet-coordinator command allows calling the WalletCoordinator contract on an
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
	WalletCoordinatorCommand := &cobra.Command{
		Use:   "wallet-coordinator",
		Short: `Provides access to the WalletCoordinator contract.`,
		Long:  walletCoordinatorDescription,
	}

	WalletCoordinatorCommand.AddCommand(
		wcBridgeCommand(),
		wcDepositMinAgeCommand(),
		wcDepositRefundSafetyMarginCommand(),
		wcDepositSweepMaxSizeCommand(),
		wcDepositSweepProposalSubmissionGasOffsetCommand(),
		wcDepositSweepProposalValidityCommand(),
		wcHeartbeatRequestGasOffsetCommand(),
		wcHeartbeatRequestValidityCommand(),
		wcIsCoordinatorCommand(),
		wcOwnerCommand(),
		wcRedemptionMaxSizeCommand(),
		wcRedemptionProposalSubmissionGasOffsetCommand(),
		wcRedemptionProposalValidityCommand(),
		wcRedemptionRequestMinAgeCommand(),
		wcRedemptionRequestTimeoutSafetyMarginCommand(),
		wcReimbursementPoolCommand(),
		wcValidateRedemptionProposalCommand(),
		wcWalletLockCommand(),
		wcAddCoordinatorCommand(),
		wcInitializeCommand(),
		wcRemoveCoordinatorCommand(),
		wcRenounceOwnershipCommand(),
		wcRequestHeartbeatCommand(),
		wcRequestHeartbeatWithReimbursementCommand(),
		wcSubmitDepositSweepProposalCommand(),
		wcSubmitDepositSweepProposalWithReimbursementCommand(),
		wcSubmitRedemptionProposalCommand(),
		wcSubmitRedemptionProposalWithReimbursementCommand(),
		wcTransferOwnershipCommand(),
		wcUnlockWalletCommand(),
		wcUpdateDepositSweepProposalParametersCommand(),
		wcUpdateHeartbeatRequestParametersCommand(),
		wcUpdateRedemptionProposalParametersCommand(),
		wcUpdateReimbursementPoolCommand(),
	)

	ModuleCommand.AddCommand(WalletCoordinatorCommand)
}

/// ------------------- Const methods -------------------

func wcBridgeCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "bridge",
		Short:                 "Calls the view method bridge on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcBridge,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcBridge(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.BridgeAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcDepositMinAgeCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "deposit-min-age",
		Short:                 "Calls the view method depositMinAge on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcDepositMinAge,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcDepositMinAge(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.DepositMinAgeAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcDepositRefundSafetyMarginCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "deposit-refund-safety-margin",
		Short:                 "Calls the view method depositRefundSafetyMargin on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcDepositRefundSafetyMargin,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcDepositRefundSafetyMargin(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.DepositRefundSafetyMarginAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcDepositSweepMaxSizeCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "deposit-sweep-max-size",
		Short:                 "Calls the view method depositSweepMaxSize on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcDepositSweepMaxSize,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcDepositSweepMaxSize(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.DepositSweepMaxSizeAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcDepositSweepProposalSubmissionGasOffsetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "deposit-sweep-proposal-submission-gas-offset",
		Short:                 "Calls the view method depositSweepProposalSubmissionGasOffset on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcDepositSweepProposalSubmissionGasOffset,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcDepositSweepProposalSubmissionGasOffset(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.DepositSweepProposalSubmissionGasOffsetAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcDepositSweepProposalValidityCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "deposit-sweep-proposal-validity",
		Short:                 "Calls the view method depositSweepProposalValidity on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcDepositSweepProposalValidity,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcDepositSweepProposalValidity(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.DepositSweepProposalValidityAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcHeartbeatRequestGasOffsetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "heartbeat-request-gas-offset",
		Short:                 "Calls the view method heartbeatRequestGasOffset on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcHeartbeatRequestGasOffset,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcHeartbeatRequestGasOffset(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.HeartbeatRequestGasOffsetAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcHeartbeatRequestValidityCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "heartbeat-request-validity",
		Short:                 "Calls the view method heartbeatRequestValidity on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcHeartbeatRequestValidity,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcHeartbeatRequestValidity(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.HeartbeatRequestValidityAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcIsCoordinatorCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "is-coordinator [arg0]",
		Short:                 "Calls the view method isCoordinator on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcIsCoordinator,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcIsCoordinator(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg0, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg0, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.IsCoordinatorAtBlock(
		arg0,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcOwnerCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "owner",
		Short:                 "Calls the view method owner on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcOwner,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcOwner(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
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

func wcRedemptionMaxSizeCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "redemption-max-size",
		Short:                 "Calls the view method redemptionMaxSize on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcRedemptionMaxSize,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcRedemptionMaxSize(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.RedemptionMaxSizeAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcRedemptionProposalSubmissionGasOffsetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "redemption-proposal-submission-gas-offset",
		Short:                 "Calls the view method redemptionProposalSubmissionGasOffset on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcRedemptionProposalSubmissionGasOffset,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcRedemptionProposalSubmissionGasOffset(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.RedemptionProposalSubmissionGasOffsetAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcRedemptionProposalValidityCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "redemption-proposal-validity",
		Short:                 "Calls the view method redemptionProposalValidity on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcRedemptionProposalValidity,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcRedemptionProposalValidity(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.RedemptionProposalValidityAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcRedemptionRequestMinAgeCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "redemption-request-min-age",
		Short:                 "Calls the view method redemptionRequestMinAge on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcRedemptionRequestMinAge,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcRedemptionRequestMinAge(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.RedemptionRequestMinAgeAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcRedemptionRequestTimeoutSafetyMarginCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "redemption-request-timeout-safety-margin",
		Short:                 "Calls the view method redemptionRequestTimeoutSafetyMargin on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcRedemptionRequestTimeoutSafetyMargin,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcRedemptionRequestTimeoutSafetyMargin(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.RedemptionRequestTimeoutSafetyMarginAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcReimbursementPoolCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "reimbursement-pool",
		Short:                 "Calls the view method reimbursementPool on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcReimbursementPool,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcReimbursementPool(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	result, err := contract.ReimbursementPoolAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcValidateRedemptionProposalCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "validate-redemption-proposal [arg_proposal_json]",
		Short:                 "Calls the view method validateRedemptionProposal on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcValidateRedemptionProposal,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcValidateRedemptionProposal(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_proposal_json := abi.WalletCoordinatorRedemptionProposal{}
	if err := json.Unmarshal([]byte(args[0]), &arg_proposal_json); err != nil {
		return fmt.Errorf("failed to unmarshal arg_proposal_json to abi.WalletCoordinatorRedemptionProposal: %w", err)
	}

	result, err := contract.ValidateRedemptionProposalAtBlock(
		arg_proposal_json,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wcWalletLockCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "wallet-lock [arg0]",
		Short:                 "Calls the view method walletLock on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcWalletLock,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func wcWalletLock(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg0, err := decode.ParseBytes20(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg0, a bytes20, from passed value %v",
			args[0],
		)
	}

	result, err := contract.WalletLockAtBlock(
		arg0,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func wcAddCoordinatorCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "add-coordinator [arg_coordinator]",
		Short:                 "Calls the nonpayable method addCoordinator on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcAddCoordinator,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcAddCoordinator(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_coordinator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_coordinator, a address, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.AddCoordinator(
			arg_coordinator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallAddCoordinator(
			arg_coordinator,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcInitializeCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "initialize [arg__bridge]",
		Short:                 "Calls the nonpayable method initialize on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcInitialize,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcInitialize(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg__bridge, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__bridge, a address, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Initialize(
			arg__bridge,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallInitialize(
			arg__bridge,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcRemoveCoordinatorCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "remove-coordinator [arg_coordinator]",
		Short:                 "Calls the nonpayable method removeCoordinator on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcRemoveCoordinator,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcRemoveCoordinator(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_coordinator, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_coordinator, a address, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RemoveCoordinator(
			arg_coordinator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallRemoveCoordinator(
			arg_coordinator,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcRenounceOwnershipCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "renounce-ownership",
		Short:                 "Calls the nonpayable method renounceOwnership on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  wcRenounceOwnership,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcRenounceOwnership(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
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

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcRequestHeartbeatCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "request-heartbeat [arg_walletPubKeyHash] [arg_message]",
		Short:                 "Calls the nonpayable method requestHeartbeat on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(2),
		RunE:                  wcRequestHeartbeat,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcRequestHeartbeat(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_walletPubKeyHash, err := decode.ParseBytes20(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_walletPubKeyHash, a bytes20, from passed value %v",
			args[0],
		)
	}
	arg_message, err := hexutil.Decode(args[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_message, a bytes, from passed value %v",
			args[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RequestHeartbeat(
			arg_walletPubKeyHash,
			arg_message,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallRequestHeartbeat(
			arg_walletPubKeyHash,
			arg_message,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcRequestHeartbeatWithReimbursementCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "request-heartbeat-with-reimbursement [arg_walletPubKeyHash] [arg_message]",
		Short:                 "Calls the nonpayable method requestHeartbeatWithReimbursement on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(2),
		RunE:                  wcRequestHeartbeatWithReimbursement,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcRequestHeartbeatWithReimbursement(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_walletPubKeyHash, err := decode.ParseBytes20(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_walletPubKeyHash, a bytes20, from passed value %v",
			args[0],
		)
	}
	arg_message, err := hexutil.Decode(args[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_message, a bytes, from passed value %v",
			args[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RequestHeartbeatWithReimbursement(
			arg_walletPubKeyHash,
			arg_message,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallRequestHeartbeatWithReimbursement(
			arg_walletPubKeyHash,
			arg_message,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcSubmitDepositSweepProposalCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "submit-deposit-sweep-proposal [arg_proposal_json]",
		Short:                 "Calls the nonpayable method submitDepositSweepProposal on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcSubmitDepositSweepProposal,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcSubmitDepositSweepProposal(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_proposal_json := abi.WalletCoordinatorDepositSweepProposal{}
	if err := json.Unmarshal([]byte(args[0]), &arg_proposal_json); err != nil {
		return fmt.Errorf("failed to unmarshal arg_proposal_json to abi.WalletCoordinatorDepositSweepProposal: %w", err)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.SubmitDepositSweepProposal(
			arg_proposal_json,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallSubmitDepositSweepProposal(
			arg_proposal_json,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcSubmitDepositSweepProposalWithReimbursementCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "submit-deposit-sweep-proposal-with-reimbursement [arg_proposal_json]",
		Short:                 "Calls the nonpayable method submitDepositSweepProposalWithReimbursement on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcSubmitDepositSweepProposalWithReimbursement,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcSubmitDepositSweepProposalWithReimbursement(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_proposal_json := abi.WalletCoordinatorDepositSweepProposal{}
	if err := json.Unmarshal([]byte(args[0]), &arg_proposal_json); err != nil {
		return fmt.Errorf("failed to unmarshal arg_proposal_json to abi.WalletCoordinatorDepositSweepProposal: %w", err)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.SubmitDepositSweepProposalWithReimbursement(
			arg_proposal_json,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallSubmitDepositSweepProposalWithReimbursement(
			arg_proposal_json,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcSubmitRedemptionProposalCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "submit-redemption-proposal [arg_proposal_json]",
		Short:                 "Calls the nonpayable method submitRedemptionProposal on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcSubmitRedemptionProposal,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcSubmitRedemptionProposal(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_proposal_json := abi.WalletCoordinatorRedemptionProposal{}
	if err := json.Unmarshal([]byte(args[0]), &arg_proposal_json); err != nil {
		return fmt.Errorf("failed to unmarshal arg_proposal_json to abi.WalletCoordinatorRedemptionProposal: %w", err)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.SubmitRedemptionProposal(
			arg_proposal_json,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallSubmitRedemptionProposal(
			arg_proposal_json,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcSubmitRedemptionProposalWithReimbursementCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "submit-redemption-proposal-with-reimbursement [arg_proposal_json]",
		Short:                 "Calls the nonpayable method submitRedemptionProposalWithReimbursement on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcSubmitRedemptionProposalWithReimbursement,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcSubmitRedemptionProposalWithReimbursement(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_proposal_json := abi.WalletCoordinatorRedemptionProposal{}
	if err := json.Unmarshal([]byte(args[0]), &arg_proposal_json); err != nil {
		return fmt.Errorf("failed to unmarshal arg_proposal_json to abi.WalletCoordinatorRedemptionProposal: %w", err)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.SubmitRedemptionProposalWithReimbursement(
			arg_proposal_json,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallSubmitRedemptionProposalWithReimbursement(
			arg_proposal_json,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcTransferOwnershipCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "transfer-ownership [arg_newOwner]",
		Short:                 "Calls the nonpayable method transferOwnership on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcTransferOwnership,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcTransferOwnership(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
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

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcUnlockWalletCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "unlock-wallet [arg_walletPubKeyHash]",
		Short:                 "Calls the nonpayable method unlockWallet on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcUnlockWallet,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcUnlockWallet(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg_walletPubKeyHash, err := decode.ParseBytes20(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_walletPubKeyHash, a bytes20, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UnlockWallet(
			arg_walletPubKeyHash,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallUnlockWallet(
			arg_walletPubKeyHash,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcUpdateDepositSweepProposalParametersCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "update-deposit-sweep-proposal-parameters [arg__depositSweepProposalValidity] [arg__depositMinAge] [arg__depositRefundSafetyMargin] [arg__depositSweepMaxSize] [arg__depositSweepProposalSubmissionGasOffset]",
		Short:                 "Calls the nonpayable method updateDepositSweepProposalParameters on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(5),
		RunE:                  wcUpdateDepositSweepProposalParameters,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcUpdateDepositSweepProposalParameters(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg__depositSweepProposalValidity, err := decode.ParseUint[uint32](args[0], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__depositSweepProposalValidity, a uint32, from passed value %v",
			args[0],
		)
	}
	arg__depositMinAge, err := decode.ParseUint[uint32](args[1], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__depositMinAge, a uint32, from passed value %v",
			args[1],
		)
	}
	arg__depositRefundSafetyMargin, err := decode.ParseUint[uint32](args[2], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__depositRefundSafetyMargin, a uint32, from passed value %v",
			args[2],
		)
	}
	arg__depositSweepMaxSize, err := decode.ParseUint[uint16](args[3], 16)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__depositSweepMaxSize, a uint16, from passed value %v",
			args[3],
		)
	}
	arg__depositSweepProposalSubmissionGasOffset, err := decode.ParseUint[uint32](args[4], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__depositSweepProposalSubmissionGasOffset, a uint32, from passed value %v",
			args[4],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateDepositSweepProposalParameters(
			arg__depositSweepProposalValidity,
			arg__depositMinAge,
			arg__depositRefundSafetyMargin,
			arg__depositSweepMaxSize,
			arg__depositSweepProposalSubmissionGasOffset,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallUpdateDepositSweepProposalParameters(
			arg__depositSweepProposalValidity,
			arg__depositMinAge,
			arg__depositRefundSafetyMargin,
			arg__depositSweepMaxSize,
			arg__depositSweepProposalSubmissionGasOffset,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcUpdateHeartbeatRequestParametersCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "update-heartbeat-request-parameters [arg__heartbeatRequestValidity] [arg__heartbeatRequestGasOffset]",
		Short:                 "Calls the nonpayable method updateHeartbeatRequestParameters on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(2),
		RunE:                  wcUpdateHeartbeatRequestParameters,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcUpdateHeartbeatRequestParameters(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg__heartbeatRequestValidity, err := decode.ParseUint[uint32](args[0], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__heartbeatRequestValidity, a uint32, from passed value %v",
			args[0],
		)
	}
	arg__heartbeatRequestGasOffset, err := decode.ParseUint[uint32](args[1], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__heartbeatRequestGasOffset, a uint32, from passed value %v",
			args[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateHeartbeatRequestParameters(
			arg__heartbeatRequestValidity,
			arg__heartbeatRequestGasOffset,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallUpdateHeartbeatRequestParameters(
			arg__heartbeatRequestValidity,
			arg__heartbeatRequestGasOffset,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcUpdateRedemptionProposalParametersCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "update-redemption-proposal-parameters [arg__redemptionProposalValidity] [arg__redemptionRequestMinAge] [arg__redemptionRequestTimeoutSafetyMargin] [arg__redemptionMaxSize] [arg__redemptionProposalSubmissionGasOffset]",
		Short:                 "Calls the nonpayable method updateRedemptionProposalParameters on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(5),
		RunE:                  wcUpdateRedemptionProposalParameters,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcUpdateRedemptionProposalParameters(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg__redemptionProposalValidity, err := decode.ParseUint[uint32](args[0], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__redemptionProposalValidity, a uint32, from passed value %v",
			args[0],
		)
	}
	arg__redemptionRequestMinAge, err := decode.ParseUint[uint32](args[1], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__redemptionRequestMinAge, a uint32, from passed value %v",
			args[1],
		)
	}
	arg__redemptionRequestTimeoutSafetyMargin, err := decode.ParseUint[uint32](args[2], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__redemptionRequestTimeoutSafetyMargin, a uint32, from passed value %v",
			args[2],
		)
	}
	arg__redemptionMaxSize, err := decode.ParseUint[uint16](args[3], 16)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__redemptionMaxSize, a uint16, from passed value %v",
			args[3],
		)
	}
	arg__redemptionProposalSubmissionGasOffset, err := decode.ParseUint[uint32](args[4], 32)
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__redemptionProposalSubmissionGasOffset, a uint32, from passed value %v",
			args[4],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateRedemptionProposalParameters(
			arg__redemptionProposalValidity,
			arg__redemptionRequestMinAge,
			arg__redemptionRequestTimeoutSafetyMargin,
			arg__redemptionMaxSize,
			arg__redemptionProposalSubmissionGasOffset,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallUpdateRedemptionProposalParameters(
			arg__redemptionProposalValidity,
			arg__redemptionRequestMinAge,
			arg__redemptionRequestTimeoutSafetyMargin,
			arg__redemptionMaxSize,
			arg__redemptionProposalSubmissionGasOffset,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

func wcUpdateReimbursementPoolCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "update-reimbursement-pool [arg__reimbursementPool]",
		Short:                 "Calls the nonpayable method updateReimbursementPool on the WalletCoordinator contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  wcUpdateReimbursementPool,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func wcUpdateReimbursementPool(c *cobra.Command, args []string) error {
	contract, err := initializeWalletCoordinator(c)
	if err != nil {
		return err
	}

	arg__reimbursementPool, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__reimbursementPool, a address, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateReimbursementPool(
			arg__reimbursementPool,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallUpdateReimbursementPool(
			arg__reimbursementPool,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")

		cmd.PrintOutput(
			"the transaction was not submitted to the chain; " +
				"please add the `--submit` flag",
		)
	}

	return nil
}

/// ------------------- Initialization -------------------

func initializeWalletCoordinator(c *cobra.Command) (*contract.WalletCoordinator, error) {
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

	address, err := cfg.ContractAddress("WalletCoordinator")
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get %s address: [%w]",
			"WalletCoordinator",
			err,
		)
	}

	return contract.NewWalletCoordinator(
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
