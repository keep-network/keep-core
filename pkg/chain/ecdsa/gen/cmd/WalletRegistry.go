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
	"github.com/keep-network/keep-core/pkg/chain/ecdsa/gen/contract"

	"github.com/urfave/cli"
)

var WalletRegistryCommand cli.Command

var walletRegistryDescription = `The wallet-registry command allows calling the WalletRegistry contract on an
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
		Name:        "wallet-registry",
		Usage:       `Provides access to the WalletRegistry contract.`,
		Description: walletRegistryDescription,
		Subcommands: []cli.Command{{
			Name:      "authorization-parameters",
			Usage:     "Calls the view method authorizationParameters on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrAuthorizationParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "available-rewards",
			Usage:     "Calls the view method availableRewards on the WalletRegistry contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    wrAvailableRewards,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "dkg-parameters",
			Usage:     "Calls the view method dkgParameters on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrDkgParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "eligible-stake",
			Usage:     "Calls the view method eligibleStake on the WalletRegistry contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    wrEligibleStake,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "gas-parameters",
			Usage:     "Calls the view method gasParameters on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrGasParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-wallet-creation-state",
			Usage:     "Calls the view method getWalletCreationState on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrGetWalletCreationState,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "governance",
			Usage:     "Calls the view method governance on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrGovernance,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "has-dkg-timed-out",
			Usage:     "Calls the view method hasDkgTimedOut on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrHasDkgTimedOut,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "has-seed-timed-out",
			Usage:     "Calls the view method hasSeedTimedOut on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrHasSeedTimedOut,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-operator-in-pool",
			Usage:     "Calls the view method isOperatorInPool on the WalletRegistry contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    wrIsOperatorInPool,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-operator-up-to-date",
			Usage:     "Calls the view method isOperatorUpToDate on the WalletRegistry contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    wrIsOperatorUpToDate,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "minimum-authorization",
			Usage:     "Calls the view method minimumAuthorization on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrMinimumAuthorization,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "operator-to-staking-provider",
			Usage:     "Calls the view method operatorToStakingProvider on the WalletRegistry contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    wrOperatorToStakingProvider,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "pending-authorization-decrease",
			Usage:     "Calls the view method pendingAuthorizationDecrease on the WalletRegistry contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    wrPendingAuthorizationDecrease,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "random-beacon",
			Usage:     "Calls the view method randomBeacon on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrRandomBeacon,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "reimbursement-pool",
			Usage:     "Calls the view method reimbursementPool on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrReimbursementPool,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "remaining-authorization-decrease-delay",
			Usage:     "Calls the view method remainingAuthorizationDecreaseDelay on the WalletRegistry contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    wrRemainingAuthorizationDecreaseDelay,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "reward-parameters",
			Usage:     "Calls the view method rewardParameters on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrRewardParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "select-group",
			Usage:     "Calls the view method selectGroup on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrSelectGroup,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "slashing-parameters",
			Usage:     "Calls the view method slashingParameters on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrSlashingParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "sortition-pool",
			Usage:     "Calls the view method sortitionPool on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrSortitionPool,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "staking",
			Usage:     "Calls the view method staking on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrStaking,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "staking-provider-to-operator",
			Usage:     "Calls the view method stakingProviderToOperator on the WalletRegistry contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    wrStakingProviderToOperator,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "wallet-owner",
			Usage:     "Calls the view method walletOwner on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrWalletOwner,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "approve-authorization-decrease",
			Usage:     "Calls the nonpayable method approveAuthorizationDecrease on the WalletRegistry contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    wrApproveAuthorizationDecrease,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "beacon-callback",
			Usage:     "Calls the nonpayable method beaconCallback on the WalletRegistry contract.",
			ArgsUsage: "[arg_relayEntry] [arg1] ",
			Action:    wrBeaconCallback,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "initialize",
			Usage:     "Calls the nonpayable method initialize on the WalletRegistry contract.",
			ArgsUsage: "[arg__ecdsaDkgValidator] [arg__randomBeacon] [arg__reimbursementPool] ",
			Action:    wrInitialize,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(3))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "join-sortition-pool",
			Usage:     "Calls the nonpayable method joinSortitionPool on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrJoinSortitionPool,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-dkg-timeout",
			Usage:     "Calls the nonpayable method notifyDkgTimeout on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrNotifyDkgTimeout,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-seed-timeout",
			Usage:     "Calls the nonpayable method notifySeedTimeout on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrNotifySeedTimeout,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "register-operator",
			Usage:     "Calls the nonpayable method registerOperator on the WalletRegistry contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    wrRegisterOperator,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "request-new-wallet",
			Usage:     "Calls the nonpayable method requestNewWallet on the WalletRegistry contract.",
			ArgsUsage: "",
			Action:    wrRequestNewWallet,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "transfer-governance",
			Usage:     "Calls the nonpayable method transferGovernance on the WalletRegistry contract.",
			ArgsUsage: "[arg_newGovernance] ",
			Action:    wrTransferGovernance,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-dkg-parameters",
			Usage:     "Calls the nonpayable method updateDkgParameters on the WalletRegistry contract.",
			ArgsUsage: "[arg__seedTimeout] [arg__resultChallengePeriodLength] [arg__resultSubmissionTimeout] [arg__submitterPrecedencePeriodLength] ",
			Action:    wrUpdateDkgParameters,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(4))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-gas-parameters",
			Usage:     "Calls the nonpayable method updateGasParameters on the WalletRegistry contract.",
			ArgsUsage: "[arg_dkgResultSubmissionGas] [arg_dkgResultApprovalGasOffset] [arg_notifyOperatorInactivityGasOffset] [arg_notifySeedTimeoutGasOffset] [arg_notifyDkgTimeoutNegativeGasOffset] ",
			Action:    wrUpdateGasParameters,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(5))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-operator-status",
			Usage:     "Calls the nonpayable method updateOperatorStatus on the WalletRegistry contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    wrUpdateOperatorStatus,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-reimbursement-pool",
			Usage:     "Calls the nonpayable method updateReimbursementPool on the WalletRegistry contract.",
			ArgsUsage: "[arg__reimbursementPool] ",
			Action:    wrUpdateReimbursementPool,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-reward-parameters",
			Usage:     "Calls the nonpayable method updateRewardParameters on the WalletRegistry contract.",
			ArgsUsage: "[arg_maliciousDkgResultNotificationRewardMultiplier] [arg_sortitionPoolRewardsBanDuration] ",
			Action:    wrUpdateRewardParameters,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-wallet-owner",
			Usage:     "Calls the nonpayable method updateWalletOwner on the WalletRegistry contract.",
			ArgsUsage: "[arg__walletOwner] ",
			Action:    wrUpdateWalletOwner,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "upgrade-random-beacon",
			Usage:     "Calls the nonpayable method upgradeRandomBeacon on the WalletRegistry contract.",
			ArgsUsage: "[arg__randomBeacon] ",
			Action:    wrUpgradeRandomBeacon,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw-ineligible-rewards",
			Usage:     "Calls the nonpayable method withdrawIneligibleRewards on the WalletRegistry contract.",
			ArgsUsage: "[arg_recipient] ",
			Action:    wrWithdrawIneligibleRewards,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw-rewards",
			Usage:     "Calls the nonpayable method withdrawRewards on the WalletRegistry contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    wrWithdrawRewards,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func wrAuthorizationParameters(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.AuthorizationParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrAvailableRewards(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}
	arg_stakingProvider, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_stakingProvider, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.AvailableRewardsAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrDkgParameters(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.DkgParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrEligibleStake(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}
	arg_stakingProvider, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_stakingProvider, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.EligibleStakeAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrGasParameters(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.GasParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrGetWalletCreationState(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.GetWalletCreationStateAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrGovernance(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
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

func wrHasDkgTimedOut(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.HasDkgTimedOutAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrHasSeedTimedOut(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.HasSeedTimedOutAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrIsOperatorInPool(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
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

func wrIsOperatorUpToDate(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
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

	result, err := contract.IsOperatorUpToDateAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrMinimumAuthorization(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.MinimumAuthorizationAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrOperatorToStakingProvider(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
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

	result, err := contract.OperatorToStakingProviderAtBlock(
		arg_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrPendingAuthorizationDecrease(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}
	arg_stakingProvider, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_stakingProvider, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.PendingAuthorizationDecreaseAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrRandomBeacon(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.RandomBeaconAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrReimbursementPool(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.ReimbursementPoolAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrRemainingAuthorizationDecreaseDelay(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}
	arg_stakingProvider, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_stakingProvider, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.RemainingAuthorizationDecreaseDelayAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrRewardParameters(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.RewardParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrSelectGroup(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.SelectGroupAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrSlashingParameters(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.SlashingParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrSortitionPool(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.SortitionPoolAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrStaking(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.StakingAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrStakingProviderToOperator(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}
	arg_stakingProvider, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_stakingProvider, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.StakingProviderToOperatorAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func wrWalletOwner(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	result, err := contract.WalletOwnerAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func wrApproveAuthorizationDecrease(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg_stakingProvider, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_stakingProvider, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ApproveAuthorizationDecrease(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallApproveAuthorizationDecrease(
			arg_stakingProvider,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrBeaconCallback(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg_relayEntry, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_relayEntry, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	arg1, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg1, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.BeaconCallback(
			arg_relayEntry,
			arg1,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallBeaconCallback(
			arg_relayEntry,
			arg1,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrInitialize(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg__ecdsaDkgValidator, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__ecdsaDkgValidator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	arg__randomBeacon, err := chainutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__randomBeacon, a address, from passed value %v",
			c.Args()[1],
		)
	}

	arg__reimbursementPool, err := chainutil.AddressFromHex(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__reimbursementPool, a address, from passed value %v",
			c.Args()[2],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Initialize(
			arg__ecdsaDkgValidator,
			arg__randomBeacon,
			arg__reimbursementPool,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallInitialize(
			arg__ecdsaDkgValidator,
			arg__randomBeacon,
			arg__reimbursementPool,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrJoinSortitionPool(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.JoinSortitionPool()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallJoinSortitionPool(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrNotifyDkgTimeout(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifyDkgTimeout()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyDkgTimeout(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrNotifySeedTimeout(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifySeedTimeout()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifySeedTimeout(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrRegisterOperator(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
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
		transaction, err = contract.RegisterOperator(
			arg_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRegisterOperator(
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

func wrRequestNewWallet(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RequestNewWallet()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRequestNewWallet(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrTransferGovernance(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
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

func wrUpdateDkgParameters(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg__seedTimeout, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__seedTimeout, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	arg__resultChallengePeriodLength, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__resultChallengePeriodLength, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	arg__resultSubmissionTimeout, err := hexutil.DecodeBig(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__resultSubmissionTimeout, a uint256, from passed value %v",
			c.Args()[2],
		)
	}

	arg__submitterPrecedencePeriodLength, err := hexutil.DecodeBig(c.Args()[3])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__submitterPrecedencePeriodLength, a uint256, from passed value %v",
			c.Args()[3],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateDkgParameters(
			arg__seedTimeout,
			arg__resultChallengePeriodLength,
			arg__resultSubmissionTimeout,
			arg__submitterPrecedencePeriodLength,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateDkgParameters(
			arg__seedTimeout,
			arg__resultChallengePeriodLength,
			arg__resultSubmissionTimeout,
			arg__submitterPrecedencePeriodLength,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrUpdateGasParameters(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg_dkgResultSubmissionGas, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_dkgResultSubmissionGas, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	arg_dkgResultApprovalGasOffset, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_dkgResultApprovalGasOffset, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	arg_notifyOperatorInactivityGasOffset, err := hexutil.DecodeBig(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_notifyOperatorInactivityGasOffset, a uint256, from passed value %v",
			c.Args()[2],
		)
	}

	arg_notifySeedTimeoutGasOffset, err := hexutil.DecodeBig(c.Args()[3])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_notifySeedTimeoutGasOffset, a uint256, from passed value %v",
			c.Args()[3],
		)
	}

	arg_notifyDkgTimeoutNegativeGasOffset, err := hexutil.DecodeBig(c.Args()[4])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_notifyDkgTimeoutNegativeGasOffset, a uint256, from passed value %v",
			c.Args()[4],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateGasParameters(
			arg_dkgResultSubmissionGas,
			arg_dkgResultApprovalGasOffset,
			arg_notifyOperatorInactivityGasOffset,
			arg_notifySeedTimeoutGasOffset,
			arg_notifyDkgTimeoutNegativeGasOffset,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateGasParameters(
			arg_dkgResultSubmissionGas,
			arg_dkgResultApprovalGasOffset,
			arg_notifyOperatorInactivityGasOffset,
			arg_notifySeedTimeoutGasOffset,
			arg_notifyDkgTimeoutNegativeGasOffset,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrUpdateOperatorStatus(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
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
		transaction, err = contract.UpdateOperatorStatus(
			arg_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateOperatorStatus(
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

func wrUpdateReimbursementPool(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg__reimbursementPool, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__reimbursementPool, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateReimbursementPool(
			arg__reimbursementPool,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateReimbursementPool(
			arg__reimbursementPool,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrUpdateRewardParameters(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg_maliciousDkgResultNotificationRewardMultiplier, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_maliciousDkgResultNotificationRewardMultiplier, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	arg_sortitionPoolRewardsBanDuration, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_sortitionPoolRewardsBanDuration, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateRewardParameters(
			arg_maliciousDkgResultNotificationRewardMultiplier,
			arg_sortitionPoolRewardsBanDuration,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateRewardParameters(
			arg_maliciousDkgResultNotificationRewardMultiplier,
			arg_sortitionPoolRewardsBanDuration,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrUpdateWalletOwner(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg__walletOwner, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__walletOwner, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateWalletOwner(
			arg__walletOwner,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateWalletOwner(
			arg__walletOwner,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrUpgradeRandomBeacon(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg__randomBeacon, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg__randomBeacon, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpgradeRandomBeacon(
			arg__randomBeacon,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpgradeRandomBeacon(
			arg__randomBeacon,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func wrWithdrawIneligibleRewards(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
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
		transaction, err = contract.WithdrawIneligibleRewards(
			arg_recipient,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallWithdrawIneligibleRewards(
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

func wrWithdrawRewards(c *cli.Context) error {
	contract, err := initializeWalletRegistry(c)
	if err != nil {
		return err
	}

	arg_stakingProvider, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_stakingProvider, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.WithdrawRewards(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallWithdrawRewards(
			arg_stakingProvider,
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

func initializeWalletRegistry(c *cli.Context) (*contract.WalletRegistry, error) {
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

	address := common.HexToAddress(config.ContractAddresses["WalletRegistry"])

	return contract.NewWalletRegistry(
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
