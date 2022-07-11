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
	"github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen/contract"

	"github.com/urfave/cli"
)

var RandomBeaconCommand cli.Command

var randomBeaconDescription = `The random-beacon command allows calling the RandomBeacon contract on an
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
		Name:        "random-beacon",
		Usage:       `Provides access to the RandomBeacon contract.`,
		Description: randomBeaconDescription,
		Subcommands: []cli.Command{{
			Name:      "authorization-parameters",
			Usage:     "Calls the view method authorizationParameters on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbAuthorizationParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "authorized-requesters",
			Usage:     "Calls the view method authorizedRequesters on the RandomBeacon contract.",
			ArgsUsage: "[arg0] ",
			Action:    rbAuthorizedRequesters,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "available-rewards",
			Usage:     "Calls the view method availableRewards on the RandomBeacon contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    rbAvailableRewards,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "eligible-stake",
			Usage:     "Calls the view method eligibleStake on the RandomBeacon contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    rbEligibleStake,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "gas-parameters",
			Usage:     "Calls the view method gasParameters on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbGasParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "genesis-seed",
			Usage:     "Calls the view method genesisSeed on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbGenesisSeed,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-group0",
			Usage:     "Calls the view method getGroup0 on the RandomBeacon contract.",
			ArgsUsage: "[arg_groupPubKey] ",
			Action:    rbGetGroup0,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-group-creation-state",
			Usage:     "Calls the view method getGroupCreationState on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbGetGroupCreationState,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-groups-registry",
			Usage:     "Calls the view method getGroupsRegistry on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbGetGroupsRegistry,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "governance",
			Usage:     "Calls the view method governance on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbGovernance,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "group-creation-parameters",
			Usage:     "Calls the view method groupCreationParameters on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbGroupCreationParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "has-dkg-timed-out",
			Usage:     "Calls the view method hasDkgTimedOut on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbHasDkgTimedOut,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-operator-in-pool",
			Usage:     "Calls the view method isOperatorInPool on the RandomBeacon contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    rbIsOperatorInPool,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-operator-up-to-date",
			Usage:     "Calls the view method isOperatorUpToDate on the RandomBeacon contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    rbIsOperatorUpToDate,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-relay-request-in-progress",
			Usage:     "Calls the view method isRelayRequestInProgress on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbIsRelayRequestInProgress,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "minimum-authorization",
			Usage:     "Calls the view method minimumAuthorization on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbMinimumAuthorization,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "operator-to-staking-provider",
			Usage:     "Calls the view method operatorToStakingProvider on the RandomBeacon contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    rbOperatorToStakingProvider,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "pending-authorization-decrease",
			Usage:     "Calls the view method pendingAuthorizationDecrease on the RandomBeacon contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    rbPendingAuthorizationDecrease,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "reimbursement-pool",
			Usage:     "Calls the view method reimbursementPool on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbReimbursementPool,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "relay-entry-parameters",
			Usage:     "Calls the view method relayEntryParameters on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbRelayEntryParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "remaining-authorization-decrease-delay",
			Usage:     "Calls the view method remainingAuthorizationDecreaseDelay on the RandomBeacon contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    rbRemainingAuthorizationDecreaseDelay,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "reward-parameters",
			Usage:     "Calls the view method rewardParameters on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbRewardParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "select-group",
			Usage:     "Calls the view method selectGroup on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbSelectGroup,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "slashing-parameters",
			Usage:     "Calls the view method slashingParameters on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbSlashingParameters,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "sortition-pool",
			Usage:     "Calls the view method sortitionPool on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbSortitionPool,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "staking",
			Usage:     "Calls the view method staking on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbStaking,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "staking-provider-to-operator",
			Usage:     "Calls the view method stakingProviderToOperator on the RandomBeacon contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    rbStakingProviderToOperator,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "t-token",
			Usage:     "Calls the view method tToken on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbTToken,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "approve-authorization-decrease",
			Usage:     "Calls the nonpayable method approveAuthorizationDecrease on the RandomBeacon contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    rbApproveAuthorizationDecrease,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "genesis",
			Usage:     "Calls the nonpayable method genesis on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbGenesis,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "join-sortition-pool",
			Usage:     "Calls the nonpayable method joinSortitionPool on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbJoinSortitionPool,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-dkg-timeout",
			Usage:     "Calls the nonpayable method notifyDkgTimeout on the RandomBeacon contract.",
			ArgsUsage: "",
			Action:    rbNotifyDkgTimeout,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "register-operator",
			Usage:     "Calls the nonpayable method registerOperator on the RandomBeacon contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    rbRegisterOperator,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "request-relay-entry",
			Usage:     "Calls the nonpayable method requestRelayEntry on the RandomBeacon contract.",
			ArgsUsage: "[arg_callbackContract] ",
			Action:    rbRequestRelayEntry,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "submit-relay-entry0",
			Usage:     "Calls the nonpayable method submitRelayEntry0 on the RandomBeacon contract.",
			ArgsUsage: "[arg_entry] ",
			Action:    rbSubmitRelayEntry0,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "transfer-governance",
			Usage:     "Calls the nonpayable method transferGovernance on the RandomBeacon contract.",
			ArgsUsage: "[arg_newGovernance] ",
			Action:    rbTransferGovernance,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-gas-parameters",
			Usage:     "Calls the nonpayable method updateGasParameters on the RandomBeacon contract.",
			ArgsUsage: "[arg_dkgResultSubmissionGas] [arg_dkgResultApprovalGasOffset] [arg_notifyOperatorInactivityGasOffset] [arg_relayEntrySubmissionGasOffset] ",
			Action:    rbUpdateGasParameters,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(4))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-group-creation-parameters",
			Usage:     "Calls the nonpayable method updateGroupCreationParameters on the RandomBeacon contract.",
			ArgsUsage: "[arg_groupCreationFrequency] [arg_groupLifetime] [arg_dkgResultChallengePeriodLength] [arg_dkgResultSubmissionTimeout] [arg_dkgSubmitterPrecedencePeriodLength] ",
			Action:    rbUpdateGroupCreationParameters,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(5))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-operator-status",
			Usage:     "Calls the nonpayable method updateOperatorStatus on the RandomBeacon contract.",
			ArgsUsage: "[arg_operator] ",
			Action:    rbUpdateOperatorStatus,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-reimbursement-pool",
			Usage:     "Calls the nonpayable method updateReimbursementPool on the RandomBeacon contract.",
			ArgsUsage: "[arg__reimbursementPool] ",
			Action:    rbUpdateReimbursementPool,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-relay-entry-parameters",
			Usage:     "Calls the nonpayable method updateRelayEntryParameters on the RandomBeacon contract.",
			ArgsUsage: "[arg_relayEntrySoftTimeout] [arg_relayEntryHardTimeout] [arg_callbackGasLimit] ",
			Action:    rbUpdateRelayEntryParameters,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(3))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "update-reward-parameters",
			Usage:     "Calls the nonpayable method updateRewardParameters on the RandomBeacon contract.",
			ArgsUsage: "[arg_sortitionPoolRewardsBanDuration] [arg_relayEntryTimeoutNotificationRewardMultiplier] [arg_unauthorizedSigningNotificationRewardMultiplier] [arg_dkgMaliciousResultNotificationRewardMultiplier] ",
			Action:    rbUpdateRewardParameters,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(4))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw-ineligible-rewards",
			Usage:     "Calls the nonpayable method withdrawIneligibleRewards on the RandomBeacon contract.",
			ArgsUsage: "[arg_recipient] ",
			Action:    rbWithdrawIneligibleRewards,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw-rewards",
			Usage:     "Calls the nonpayable method withdrawRewards on the RandomBeacon contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    rbWithdrawRewards,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func rbAuthorizationParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbAuthorizedRequesters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}
	arg0, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg0, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.AuthorizedRequestersAtBlock(
		arg0,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func rbAvailableRewards(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbEligibleStake(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbGasParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbGenesisSeed(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	result, err := contract.GenesisSeedAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func rbGetGroup0(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}
	arg_groupPubKey, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_groupPubKey, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetGroup0AtBlock(
		arg_groupPubKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func rbGetGroupCreationState(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	result, err := contract.GetGroupCreationStateAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func rbGetGroupsRegistry(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	result, err := contract.GetGroupsRegistryAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func rbGovernance(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbGroupCreationParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	result, err := contract.GroupCreationParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func rbHasDkgTimedOut(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbIsOperatorInPool(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbIsOperatorUpToDate(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbIsRelayRequestInProgress(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	result, err := contract.IsRelayRequestInProgressAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func rbMinimumAuthorization(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbOperatorToStakingProvider(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbPendingAuthorizationDecrease(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbReimbursementPool(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbRelayEntryParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	result, err := contract.RelayEntryParametersAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func rbRemainingAuthorizationDecreaseDelay(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbRewardParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbSelectGroup(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbSlashingParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbSortitionPool(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbStaking(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbStakingProviderToOperator(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbTToken(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	result, err := contract.TTokenAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func rbApproveAuthorizationDecrease(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbGenesis(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Genesis()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallGenesis(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func rbJoinSortitionPool(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbNotifyDkgTimeout(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbRegisterOperator(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbRequestRelayEntry(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	arg_callbackContract, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_callbackContract, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RequestRelayEntry(
			arg_callbackContract,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRequestRelayEntry(
			arg_callbackContract,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func rbSubmitRelayEntry0(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	arg_entry, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_entry, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.SubmitRelayEntry0(
			arg_entry,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallSubmitRelayEntry0(
			arg_entry,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func rbTransferGovernance(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbUpdateGasParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

	arg_relayEntrySubmissionGasOffset, err := hexutil.DecodeBig(c.Args()[3])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_relayEntrySubmissionGasOffset, a uint256, from passed value %v",
			c.Args()[3],
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
			arg_relayEntrySubmissionGasOffset,
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
			arg_relayEntrySubmissionGasOffset,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func rbUpdateGroupCreationParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	arg_groupCreationFrequency, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_groupCreationFrequency, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	arg_groupLifetime, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_groupLifetime, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	arg_dkgResultChallengePeriodLength, err := hexutil.DecodeBig(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_dkgResultChallengePeriodLength, a uint256, from passed value %v",
			c.Args()[2],
		)
	}

	arg_dkgResultSubmissionTimeout, err := hexutil.DecodeBig(c.Args()[3])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_dkgResultSubmissionTimeout, a uint256, from passed value %v",
			c.Args()[3],
		)
	}

	arg_dkgSubmitterPrecedencePeriodLength, err := hexutil.DecodeBig(c.Args()[4])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_dkgSubmitterPrecedencePeriodLength, a uint256, from passed value %v",
			c.Args()[4],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateGroupCreationParameters(
			arg_groupCreationFrequency,
			arg_groupLifetime,
			arg_dkgResultChallengePeriodLength,
			arg_dkgResultSubmissionTimeout,
			arg_dkgSubmitterPrecedencePeriodLength,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateGroupCreationParameters(
			arg_groupCreationFrequency,
			arg_groupLifetime,
			arg_dkgResultChallengePeriodLength,
			arg_dkgResultSubmissionTimeout,
			arg_dkgSubmitterPrecedencePeriodLength,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func rbUpdateOperatorStatus(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbUpdateReimbursementPool(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbUpdateRelayEntryParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	arg_relayEntrySoftTimeout, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_relayEntrySoftTimeout, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	arg_relayEntryHardTimeout, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_relayEntryHardTimeout, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	arg_callbackGasLimit, err := hexutil.DecodeBig(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_callbackGasLimit, a uint256, from passed value %v",
			c.Args()[2],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateRelayEntryParameters(
			arg_relayEntrySoftTimeout,
			arg_relayEntryHardTimeout,
			arg_callbackGasLimit,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateRelayEntryParameters(
			arg_relayEntrySoftTimeout,
			arg_relayEntryHardTimeout,
			arg_callbackGasLimit,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func rbUpdateRewardParameters(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
	if err != nil {
		return err
	}

	arg_sortitionPoolRewardsBanDuration, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_sortitionPoolRewardsBanDuration, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	arg_relayEntryTimeoutNotificationRewardMultiplier, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_relayEntryTimeoutNotificationRewardMultiplier, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	arg_unauthorizedSigningNotificationRewardMultiplier, err := hexutil.DecodeBig(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_unauthorizedSigningNotificationRewardMultiplier, a uint256, from passed value %v",
			c.Args()[2],
		)
	}

	arg_dkgMaliciousResultNotificationRewardMultiplier, err := hexutil.DecodeBig(c.Args()[3])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_dkgMaliciousResultNotificationRewardMultiplier, a uint256, from passed value %v",
			c.Args()[3],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UpdateRewardParameters(
			arg_sortitionPoolRewardsBanDuration,
			arg_relayEntryTimeoutNotificationRewardMultiplier,
			arg_unauthorizedSigningNotificationRewardMultiplier,
			arg_dkgMaliciousResultNotificationRewardMultiplier,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUpdateRewardParameters(
			arg_sortitionPoolRewardsBanDuration,
			arg_relayEntryTimeoutNotificationRewardMultiplier,
			arg_unauthorizedSigningNotificationRewardMultiplier,
			arg_dkgMaliciousResultNotificationRewardMultiplier,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func rbWithdrawIneligibleRewards(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func rbWithdrawRewards(c *cli.Context) error {
	contract, err := initializeRandomBeacon(c)
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

func initializeRandomBeacon(c *cli.Context) (*contract.RandomBeacon, error) {
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

	address := common.HexToAddress(config.ContractAddresses["RandomBeacon"])

	return contract.NewRandomBeacon(
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
