// Code generated - DO NOT EDIT.
// This file is a generated command and any manual changes will be lost.

package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/chain/tbtc-v2/gen/contract"

	"github.com/urfave/cli"
)

var BridgeCommand cli.Command

var bridgeDescription = `The bridge command allows calling the Bridge contract on an
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
		Name:        "bridge",
		Usage:       `Provides access to the Bridge contract.`,
		Description: bridgeDescription,
		Subcommands: []cli.Command{{
			Name:      "active-wallet-pub-key-hash",
			Usage:     "Calls the view method activeWalletPubKeyHash on the Bridge contract.",
			ArgsUsage: "",
			Action:    bActiveWalletPubKeyHash,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "contract-references",
			Usage:     "Calls the view method contractReferences on the Bridge contract.",
			ArgsUsage: "",
			Action:    bContractReferences,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "deposit-parameters",
			Usage:     "Calls the view method depositParameters on the Bridge contract.",
			ArgsUsage: "",
			Action:    bDepositParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "deposits",
			Usage:     "Calls the view method deposits on the Bridge contract.",
			ArgsUsage: "[arg_depositKey] ",
			Action:    bDeposits,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "fraud-challenges",
			Usage:     "Calls the view method fraudChallenges on the Bridge contract.",
			ArgsUsage: "[arg_challengeKey] ",
			Action:    bFraudChallenges,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "fraud-parameters",
			Usage:     "Calls the view method fraudParameters on the Bridge contract.",
			ArgsUsage: "",
			Action:    bFraudParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "governance",
			Usage:     "Calls the view method governance on the Bridge contract.",
			ArgsUsage: "",
			Action:    bGovernance,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-vault-trusted",
			Usage:     "Calls the view method isVaultTrusted on the Bridge contract.",
			ArgsUsage: "[arg_vault] ",
			Action:    bIsVaultTrusted,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "live-wallets-count",
			Usage:     "Calls the view method liveWalletsCount on the Bridge contract.",
			ArgsUsage: "",
			Action:    bLiveWalletsCount,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "moved-funds-sweep-requests",
			Usage:     "Calls the view method movedFundsSweepRequests on the Bridge contract.",
			ArgsUsage: "[arg_requestKey] ",
			Action:    bMovedFundsSweepRequests,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "moving-funds-parameters",
			Usage:     "Calls the view method movingFundsParameters on the Bridge contract.",
			ArgsUsage: "",
			Action:    bMovingFundsParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "pending-redemptions",
			Usage:     "Calls the view method pendingRedemptions on the Bridge contract.",
			ArgsUsage: "[arg_redemptionKey] ",
			Action:    bPendingRedemptions,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "redemption-parameters",
			Usage:     "Calls the view method redemptionParameters on the Bridge contract.",
			ArgsUsage: "",
			Action:    bRedemptionParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "spent-main-u-t-x-os",
			Usage:     "Calls the view method spentMainUTXOs on the Bridge contract.",
			ArgsUsage: "[arg_utxoKey] ",
			Action:    bSpentMainUTXOs,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "timed-out-redemptions",
			Usage:     "Calls the view method timedOutRedemptions on the Bridge contract.",
			ArgsUsage: "[arg_redemptionKey] ",
			Action:    bTimedOutRedemptions,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "treasury",
			Usage:     "Calls the view method treasury on the Bridge contract.",
			ArgsUsage: "",
			Action:    bTreasury,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "tx-proof-difficulty-factor",
			Usage:     "Calls the view method txProofDifficultyFactor on the Bridge contract.",
			ArgsUsage: "",
			Action:    bTxProofDifficultyFactor,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "wallet-parameters",
			Usage:     "Calls the view method walletParameters on the Bridge contract.",
			ArgsUsage: "",
			Action:    bWalletParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "defeat-fraud-challenge-with-heartbeat",
			Usage:     "Calls the nonpayable method defeatFraudChallengeWithHeartbeat on the Bridge contract.",
			ArgsUsage: "[arg_walletPublicKey] [arg_heartbeatMessage] ",
			Action:    bDefeatFraudChallengeWithHeartbeat,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "receive-balance-approval",
			Usage:     "Calls the nonpayable method receiveBalanceApproval on the Bridge contract.",
			ArgsUsage: "[arg_balanceOwner] [arg_amount] [arg_redemptionData] ",
			Action:    bReceiveBalanceApproval,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(3))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "transfer-governance",
			Usage:     "Calls the nonpayable method transferGovernance on the Bridge contract.",
			ArgsUsage: "[arg_newGovernance] ",
			Action:    bTransferGovernance,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func bActiveWalletPubKeyHash(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.ActiveWalletPubKeyHashAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bContractReferences(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.ContractReferencesAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bDepositParameters(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.DepositParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bDeposits(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_depositKey, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_depositKey, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.DepositsAtBlock(
		arg_depositKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bFraudChallenges(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_challengeKey, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_challengeKey, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.FraudChallengesAtBlock(
		arg_challengeKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bFraudParameters(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.FraudParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bGovernance(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.GovernanceAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bIsVaultTrusted(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_vault, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_vault, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.IsVaultTrustedAtBlock(
		arg_vault,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bLiveWalletsCount(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.LiveWalletsCountAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bMovedFundsSweepRequests(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_requestKey, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_requestKey, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.MovedFundsSweepRequestsAtBlock(
		arg_requestKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bMovingFundsParameters(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.MovingFundsParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bPendingRedemptions(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_redemptionKey, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_redemptionKey, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.PendingRedemptionsAtBlock(
		arg_redemptionKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bRedemptionParameters(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.RedemptionParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bSpentMainUTXOs(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_utxoKey, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_utxoKey, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.SpentMainUTXOsAtBlock(
		arg_utxoKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bTimedOutRedemptions(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}
	arg_redemptionKey, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_redemptionKey, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.TimedOutRedemptionsAtBlock(
		arg_redemptionKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bTreasury(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.TreasuryAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bTxProofDifficultyFactor(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.TxProofDifficultyFactorAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func bWalletParameters(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	result, err := contract.WalletParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func bDefeatFraudChallengeWithHeartbeat(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	arg_walletPublicKey, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_walletPublicKey, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	arg_heartbeatMessage, err := hexutil.Decode(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_heartbeatMessage, a bytes, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.DefeatFraudChallengeWithHeartbeat(
			arg_walletPublicKey,
			arg_heartbeatMessage,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallDefeatFraudChallengeWithHeartbeat(
			arg_walletPublicKey,
			arg_heartbeatMessage,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bReceiveBalanceApproval(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	arg_balanceOwner, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_balanceOwner, a address, from passed value %v",
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

	arg_redemptionData, err := hexutil.Decode(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_redemptionData, a bytes, from passed value %v",
			c.Args()[2],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ReceiveBalanceApproval(
			arg_balanceOwner,
			arg_amount,
			arg_redemptionData,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallReceiveBalanceApproval(
			arg_balanceOwner,
			arg_amount,
			arg_redemptionData,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func bTransferGovernance(c *cli.Context) error {
	contract, err := initializeBridge(c)
	if err != nil {
		return err
	}

	arg_newGovernance, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_newGovernance, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.TransferGovernance(
			arg_newGovernance,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallTransferGovernance(
			arg_newGovernance,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

/// ------------------- Initialization -------------------

func initializeBridge(c *cli.Context) (*contract.Bridge, error) {
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

	address := common.HexToAddress(config.ContractAddresses["Bridge"])

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
