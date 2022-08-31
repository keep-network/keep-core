// Code generated - DO NOT EDIT.
// This file is a generated command and any manual changes will be lost.

package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/contract"

	"github.com/spf13/cobra"
)

var BridgeCommand *cobra.Command

var bridgeDescription = `The bridge command allows calling the Bridge contract on an
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
	BridgeCommand := &cobra.Command{
		Use:   "bridge",
		Short: `Provides access to the Bridge contract.`,
		Long:  bridgeDescription,
	}

	BridgeCommand.AddCommand(
		bActiveWalletPubKeyHashCommand(),
		bContractReferencesCommand(),
		bDepositParametersCommand(),
		bDepositsCommand(),
		bFraudChallengesCommand(),
		bFraudParametersCommand(),
		bGovernanceCommand(),
		bIsVaultTrustedCommand(),
		bLiveWalletsCountCommand(),
		bMovedFundsSweepRequestsCommand(),
		bMovingFundsParametersCommand(),
		bPendingRedemptionsCommand(),
		bRedemptionParametersCommand(),
		bSpentMainUTXOsCommand(),
		bTimedOutRedemptionsCommand(),
		bTreasuryCommand(),
		bTxProofDifficultyFactorCommand(),
		bWalletParametersCommand(),
		bDefeatFraudChallengeWithHeartbeatCommand(),
		bReceiveBalanceApprovalCommand(),
		bTransferGovernanceCommand(),
	)

	Command.AddCommand(BridgeCommand)
}

/// ------------------- Const methods -------------------

func bActiveWalletPubKeyHashCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "active-wallet-pub-key-hash",
		Short:                 "Calls the view method activeWalletPubKeyHash on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bActiveWalletPubKeyHash,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bActiveWalletPubKeyHash(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.ActiveWalletPubKeyHashAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bContractReferencesCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "contract-references",
		Short:                 "Calls the view method contractReferences on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bContractReferences,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bContractReferences(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.ContractReferencesAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bDepositParametersCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "deposit-parameters",
		Short:                 "Calls the view method depositParameters on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bDepositParameters,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bDepositParameters(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.DepositParametersAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bDepositsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "deposits [arg_depositKey]",
		Short:                 "Calls the view method deposits on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bDeposits,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bDeposits(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_depositKey, err := hexutil.DecodeBig(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_depositKey, a uint256, from passed value %v",
			args[0],
		)
	}

	result, err := contract.DepositsAtBlock(
		arg_depositKey,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bFraudChallengesCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "fraud-challenges [arg_challengeKey]",
		Short:                 "Calls the view method fraudChallenges on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bFraudChallenges,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bFraudChallenges(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_challengeKey, err := hexutil.DecodeBig(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_challengeKey, a uint256, from passed value %v",
			args[0],
		)
	}

	result, err := contract.FraudChallengesAtBlock(
		arg_challengeKey,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bFraudParametersCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "fraud-parameters",
		Short:                 "Calls the view method fraudParameters on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bFraudParameters,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bFraudParameters(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.FraudParametersAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bGovernanceCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "governance",
		Short:                 "Calls the view method governance on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bGovernance,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bGovernance(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.GovernanceAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bIsVaultTrustedCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "is-vault-trusted [arg_vault]",
		Short:                 "Calls the view method isVaultTrusted on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bIsVaultTrusted,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bIsVaultTrusted(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_vault, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_vault, a address, from passed value %v",
			args[0],
		)
	}

	result, err := contract.IsVaultTrustedAtBlock(
		arg_vault,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bLiveWalletsCountCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "live-wallets-count",
		Short:                 "Calls the view method liveWalletsCount on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bLiveWalletsCount,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bLiveWalletsCount(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.LiveWalletsCountAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bMovedFundsSweepRequestsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "moved-funds-sweep-requests [arg_requestKey]",
		Short:                 "Calls the view method movedFundsSweepRequests on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bMovedFundsSweepRequests,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bMovedFundsSweepRequests(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_requestKey, err := hexutil.DecodeBig(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_requestKey, a uint256, from passed value %v",
			args[0],
		)
	}

	result, err := contract.MovedFundsSweepRequestsAtBlock(
		arg_requestKey,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bMovingFundsParametersCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "moving-funds-parameters",
		Short:                 "Calls the view method movingFundsParameters on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bMovingFundsParameters,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bMovingFundsParameters(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.MovingFundsParametersAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bPendingRedemptionsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "pending-redemptions [arg_redemptionKey]",
		Short:                 "Calls the view method pendingRedemptions on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bPendingRedemptions,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bPendingRedemptions(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_redemptionKey, err := hexutil.DecodeBig(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_redemptionKey, a uint256, from passed value %v",
			args[0],
		)
	}

	result, err := contract.PendingRedemptionsAtBlock(
		arg_redemptionKey,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bRedemptionParametersCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "redemption-parameters",
		Short:                 "Calls the view method redemptionParameters on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bRedemptionParameters,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bRedemptionParameters(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.RedemptionParametersAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bSpentMainUTXOsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "spent-main-u-t-x-os [arg_utxoKey]",
		Short:                 "Calls the view method spentMainUTXOs on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bSpentMainUTXOs,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bSpentMainUTXOs(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_utxoKey, err := hexutil.DecodeBig(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_utxoKey, a uint256, from passed value %v",
			args[0],
		)
	}

	result, err := contract.SpentMainUTXOsAtBlock(
		arg_utxoKey,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bTimedOutRedemptionsCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "timed-out-redemptions [arg_redemptionKey]",
		Short:                 "Calls the view method timedOutRedemptions on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bTimedOutRedemptions,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bTimedOutRedemptions(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_redemptionKey, err := hexutil.DecodeBig(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_redemptionKey, a uint256, from passed value %v",
			args[0],
		)
	}

	result, err := contract.TimedOutRedemptionsAtBlock(
		arg_redemptionKey,
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bTreasuryCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "treasury",
		Short:                 "Calls the view method treasury on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bTreasury,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bTreasury(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.TreasuryAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bTxProofDifficultyFactorCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "tx-proof-difficulty-factor",
		Short:                 "Calls the view method txProofDifficultyFactor on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bTxProofDifficultyFactor,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bTxProofDifficultyFactor(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.TxProofDifficultyFactorAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bWalletParametersCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "wallet-parameters",
		Short:                 "Calls the view method walletParameters on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(0),
		RunE:                  bWalletParameters,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	cmd.InitConstFlags(c)

	return c
}

func bWalletParameters(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.WalletParametersAtBlock(
		cmd.BlockFlagValue.Int,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func bDefeatFraudChallengeWithHeartbeatCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "defeat-fraud-challenge-with-heartbeat [arg_walletPublicKey] [arg_heartbeatMessage]",
		Short:                 "Calls the nonpayable method defeatFraudChallengeWithHeartbeat on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(2),
		RunE:                  bDefeatFraudChallengeWithHeartbeat,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bDefeatFraudChallengeWithHeartbeat(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	arg_walletPublicKey, err := hexutil.Decode(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_walletPublicKey, a bytes, from passed value %v",
			args[0],
		)
	}

	arg_heartbeatMessage, err := hexutil.Decode(args[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_heartbeatMessage, a bytes, from passed value %v",
			args[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.DefeatFraudChallengeWithHeartbeat(
			arg_walletPublicKey,
			arg_heartbeatMessage,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallDefeatFraudChallengeWithHeartbeat(
			arg_walletPublicKey,
			arg_heartbeatMessage,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bReceiveBalanceApprovalCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "receive-balance-approval [arg_balanceOwner] [arg_amount] [arg_redemptionData]",
		Short:                 "Calls the nonpayable method receiveBalanceApproval on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(3),
		RunE:                  bReceiveBalanceApproval,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bReceiveBalanceApproval(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	arg_balanceOwner, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_balanceOwner, a address, from passed value %v",
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

	arg_redemptionData, err := hexutil.Decode(args[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_redemptionData, a bytes, from passed value %v",
			args[2],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ReceiveBalanceApproval(
			arg_balanceOwner,
			arg_amount,
			arg_redemptionData,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallReceiveBalanceApproval(
			arg_balanceOwner,
			arg_amount,
			arg_redemptionData,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

func bTransferGovernanceCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "transfer-governance [arg_newGovernance]",
		Short:                 "Calls the nonpayable method transferGovernance on the Bridge contract.",
		Args:                  cmd.ArgCountChecker(1),
		RunE:                  bTransferGovernance,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
	}

	c.PreRunE = cmd.NonConstArgsChecker
	cmd.InitNonConstFlags(c)

	return c
}

func bTransferGovernance(c *cobra.Command, args []string) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	arg_newGovernance, err := chainutil.AddressFromHex(args[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_newGovernance, a address, from passed value %v",
			args[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if shouldSubmit, _ := c.Flags().GetBool(cmd.SubmitFlag); shouldSubmit {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.TransferGovernance(
			arg_newGovernance,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash())
	} else {
		// Do a call.
		err = contract.CallTransferGovernance(
			arg_newGovernance,
			cmd.BlockFlagValue.Int,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput("success")
	}

	return nil
}

/// ------------------- Initialization -------------------

func initializeBridge(c *cobra.Command) (*contract.Bridge, error) {
	cfg, err := config.ReadEthereumConfig(c.Flags())
	if err != nil {
		return nil, fmt.Errorf("error reading config from file: [%v]", err)
	}

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

	address, err := cfg.ContractAddress("Bridge")
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get %s address: [%w]",
			"Bridge",
			err,
		)
	}

	return contract.NewBridge(
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
