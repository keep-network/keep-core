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
	"github.com/keep-network/keep-core/pkg/chain/ethereum/threshold/gen/contract"

	"github.com/urfave/cli"
)

var TokenStakingCommand cli.Command

var tokenStakingDescription = `The token-staking command allows calling the TokenStaking contract on an
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
		Name:        "token-staking",
		Usage:       `Provides access to the TokenStaking contract.`,
		Description: tokenStakingDescription,
		Subcommands: []cli.Command{{
			Name:      "application-info",
			Usage:     "Calls the view method applicationInfo on the TokenStaking contract.",
			ArgsUsage: "[arg0] ",
			Action:    tsApplicationInfo,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "applications",
			Usage:     "Calls the view method applications on the TokenStaking contract.",
			ArgsUsage: "[arg0] ",
			Action:    tsApplications,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "authorization-ceiling",
			Usage:     "Calls the view method authorizationCeiling on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsAuthorizationCeiling,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "authorized-stake",
			Usage:     "Calls the view method authorizedStake on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] [arg_application] ",
			Action:    tsAuthorizedStake,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "delegates",
			Usage:     "Calls the view method delegates on the TokenStaking contract.",
			ArgsUsage: "[arg_account] ",
			Action:    tsDelegates,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-applications-length",
			Usage:     "Calls the view method getApplicationsLength on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsGetApplicationsLength,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-available-to-authorize",
			Usage:     "Calls the view method getAvailableToAuthorize on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] [arg_application] ",
			Action:    tsGetAvailableToAuthorize,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-past-total-supply",
			Usage:     "Calls the view method getPastTotalSupply on the TokenStaking contract.",
			ArgsUsage: "[arg_blockNumber] ",
			Action:    tsGetPastTotalSupply,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-past-votes",
			Usage:     "Calls the view method getPastVotes on the TokenStaking contract.",
			ArgsUsage: "[arg_account] [arg_blockNumber] ",
			Action:    tsGetPastVotes,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-slashing-queue-length",
			Usage:     "Calls the view method getSlashingQueueLength on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsGetSlashingQueueLength,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-start-staking-timestamp",
			Usage:     "Calls the view method getStartStakingTimestamp on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsGetStartStakingTimestamp,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-votes",
			Usage:     "Calls the view method getVotes on the TokenStaking contract.",
			ArgsUsage: "[arg_account] ",
			Action:    tsGetVotes,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "governance",
			Usage:     "Calls the view method governance on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsGovernance,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "min-t-stake-amount",
			Usage:     "Calls the view method minTStakeAmount on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsMinTStakeAmount,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "notification-reward",
			Usage:     "Calls the view method notificationReward on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsNotificationReward,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "notifiers-treasury",
			Usage:     "Calls the view method notifiersTreasury on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsNotifiersTreasury,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "num-checkpoints",
			Usage:     "Calls the view method numCheckpoints on the TokenStaking contract.",
			ArgsUsage: "[arg_account] ",
			Action:    tsNumCheckpoints,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "roles-of",
			Usage:     "Calls the view method rolesOf on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsRolesOf,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "slashing-queue",
			Usage:     "Calls the view method slashingQueue on the TokenStaking contract.",
			ArgsUsage: "[arg0] ",
			Action:    tsSlashingQueue,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "slashing-queue-index",
			Usage:     "Calls the view method slashingQueueIndex on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsSlashingQueueIndex,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "stake-discrepancy-penalty",
			Usage:     "Calls the view method stakeDiscrepancyPenalty on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsStakeDiscrepancyPenalty,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "stake-discrepancy-reward-multiplier",
			Usage:     "Calls the view method stakeDiscrepancyRewardMultiplier on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsStakeDiscrepancyRewardMultiplier,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "staked-nu",
			Usage:     "Calls the view method stakedNu on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsStakedNu,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "stakes",
			Usage:     "Calls the view method stakes on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsStakes,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "approve-application",
			Usage:     "Calls the nonpayable method approveApplication on the TokenStaking contract.",
			ArgsUsage: "[arg_application] ",
			Action:    tsApproveApplication,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "approve-authorization-decrease",
			Usage:     "Calls the nonpayable method approveAuthorizationDecrease on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsApproveAuthorizationDecrease,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "delegate-voting",
			Usage:     "Calls the nonpayable method delegateVoting on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] [arg_delegatee] ",
			Action:    tsDelegateVoting,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "disable-application",
			Usage:     "Calls the nonpayable method disableApplication on the TokenStaking contract.",
			ArgsUsage: "[arg_application] ",
			Action:    tsDisableApplication,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "force-decrease-authorization",
			Usage:     "Calls the nonpayable method forceDecreaseAuthorization on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] [arg_application] ",
			Action:    tsForceDecreaseAuthorization,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "initialize",
			Usage:     "Calls the nonpayable method initialize on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsInitialize,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-keep-stake-discrepancy",
			Usage:     "Calls the nonpayable method notifyKeepStakeDiscrepancy on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsNotifyKeepStakeDiscrepancy,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-nu-stake-discrepancy",
			Usage:     "Calls the nonpayable method notifyNuStakeDiscrepancy on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsNotifyNuStakeDiscrepancy,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "pause-application",
			Usage:     "Calls the nonpayable method pauseApplication on the TokenStaking contract.",
			ArgsUsage: "[arg_application] ",
			Action:    tsPauseApplication,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "process-slashing",
			Usage:     "Calls the nonpayable method processSlashing on the TokenStaking contract.",
			ArgsUsage: "[arg_count] ",
			Action:    tsProcessSlashing,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "refresh-keep-stake-owner",
			Usage:     "Calls the nonpayable method refreshKeepStakeOwner on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsRefreshKeepStakeOwner,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "request-authorization-decrease0",
			Usage:     "Calls the nonpayable method requestAuthorizationDecrease0 on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsRequestAuthorizationDecrease0,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "set-authorization-ceiling",
			Usage:     "Calls the nonpayable method setAuthorizationCeiling on the TokenStaking contract.",
			ArgsUsage: "[arg_ceiling] ",
			Action:    tsSetAuthorizationCeiling,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "set-panic-button",
			Usage:     "Calls the nonpayable method setPanicButton on the TokenStaking contract.",
			ArgsUsage: "[arg_application] [arg_panicButton] ",
			Action:    tsSetPanicButton,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "stake-keep",
			Usage:     "Calls the nonpayable method stakeKeep on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsStakeKeep,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "stake-nu",
			Usage:     "Calls the nonpayable method stakeNu on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] [arg_beneficiary] [arg_authorizer] ",
			Action:    tsStakeNu,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(3))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "top-up-keep",
			Usage:     "Calls the nonpayable method topUpKeep on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsTopUpKeep,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "top-up-nu",
			Usage:     "Calls the nonpayable method topUpNu on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsTopUpNu,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "transfer-governance",
			Usage:     "Calls the nonpayable method transferGovernance on the TokenStaking contract.",
			ArgsUsage: "[arg_newGuvnor] ",
			Action:    tsTransferGovernance,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "unstake-all",
			Usage:     "Calls the nonpayable method unstakeAll on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsUnstakeAll,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "unstake-keep",
			Usage:     "Calls the nonpayable method unstakeKeep on the TokenStaking contract.",
			ArgsUsage: "[arg_stakingProvider] ",
			Action:    tsUnstakeKeep,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func tsApplicationInfo(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	result, err := contract.ApplicationInfoAtBlock(
		arg0,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsApplications(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg0, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg0, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.ApplicationsAtBlock(
		arg0,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsAuthorizationCeiling(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.AuthorizationCeilingAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsAuthorizedStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	arg_application, err := chainutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_application, a address, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.AuthorizedStakeAtBlock(
		arg_stakingProvider,
		arg_application,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsDelegates(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg_account, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_account, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.DelegatesAtBlock(
		arg_account,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetApplicationsLength(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.GetApplicationsLengthAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetAvailableToAuthorize(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	arg_application, err := chainutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_application, a address, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.GetAvailableToAuthorizeAtBlock(
		arg_stakingProvider,
		arg_application,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetPastTotalSupply(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg_blockNumber, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_blockNumber, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetPastTotalSupplyAtBlock(
		arg_blockNumber,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetPastVotes(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg_account, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_account, a address, from passed value %v",
			c.Args()[0],
		)
	}

	arg_blockNumber, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_blockNumber, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.GetPastVotesAtBlock(
		arg_account,
		arg_blockNumber,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetSlashingQueueLength(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.GetSlashingQueueLengthAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetStartStakingTimestamp(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	result, err := contract.GetStartStakingTimestampAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetVotes(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg_account, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_account, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetVotesAtBlock(
		arg_account,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGovernance(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

func tsMinTStakeAmount(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.MinTStakeAmountAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsNotificationReward(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.NotificationRewardAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsNotifiersTreasury(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.NotifiersTreasuryAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsNumCheckpoints(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg_account, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_account, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.NumCheckpointsAtBlock(
		arg_account,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsRolesOf(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	result, err := contract.RolesOfAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsSlashingQueue(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg0, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg0, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.SlashingQueueAtBlock(
		arg0,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsSlashingQueueIndex(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.SlashingQueueIndexAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsStakeDiscrepancyPenalty(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.StakeDiscrepancyPenaltyAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsStakeDiscrepancyRewardMultiplier(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.StakeDiscrepancyRewardMultiplierAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsStakedNu(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	result, err := contract.StakedNuAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsStakes(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	result, err := contract.StakesAtBlock(
		arg_stakingProvider,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func tsApproveApplication(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	arg_application, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_application, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ApproveApplication(
			arg_application,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallApproveApplication(
			arg_application,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsApproveAuthorizationDecrease(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		result      *big.Int
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
		result, err = contract.CallApproveAuthorizationDecrease(
			arg_stakingProvider,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(result)
	}

	return nil
}

func tsDelegateVoting(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	arg_delegatee, err := chainutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_delegatee, a address, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.DelegateVoting(
			arg_stakingProvider,
			arg_delegatee,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallDelegateVoting(
			arg_stakingProvider,
			arg_delegatee,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsDisableApplication(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	arg_application, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_application, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.DisableApplication(
			arg_application,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallDisableApplication(
			arg_application,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsForceDecreaseAuthorization(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	arg_application, err := chainutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_application, a address, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ForceDecreaseAuthorization(
			arg_stakingProvider,
			arg_application,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallForceDecreaseAuthorization(
			arg_stakingProvider,
			arg_application,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsInitialize(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Initialize()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallInitialize(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsNotifyKeepStakeDiscrepancy(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.NotifyKeepStakeDiscrepancy(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyKeepStakeDiscrepancy(
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

func tsNotifyNuStakeDiscrepancy(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.NotifyNuStakeDiscrepancy(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyNuStakeDiscrepancy(
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

func tsPauseApplication(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	arg_application, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_application, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.PauseApplication(
			arg_application,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallPauseApplication(
			arg_application,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsProcessSlashing(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	arg_count, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_count, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ProcessSlashing(
			arg_count,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallProcessSlashing(
			arg_count,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsRefreshKeepStakeOwner(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.RefreshKeepStakeOwner(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRefreshKeepStakeOwner(
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

func tsRequestAuthorizationDecrease0(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.RequestAuthorizationDecrease0(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRequestAuthorizationDecrease0(
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

func tsSetAuthorizationCeiling(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	arg_ceiling, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_ceiling, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.SetAuthorizationCeiling(
			arg_ceiling,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallSetAuthorizationCeiling(
			arg_ceiling,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsSetPanicButton(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	arg_application, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_application, a address, from passed value %v",
			c.Args()[0],
		)
	}

	arg_panicButton, err := chainutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_panicButton, a address, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.SetPanicButton(
			arg_application,
			arg_panicButton,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallSetPanicButton(
			arg_application,
			arg_panicButton,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsStakeKeep(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.StakeKeep(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallStakeKeep(
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

func tsStakeNu(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	arg_beneficiary, err := chainutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_beneficiary, a address, from passed value %v",
			c.Args()[1],
		)
	}

	arg_authorizer, err := chainutil.AddressFromHex(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_authorizer, a address, from passed value %v",
			c.Args()[2],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.StakeNu(
			arg_stakingProvider,
			arg_beneficiary,
			arg_authorizer,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallStakeNu(
			arg_stakingProvider,
			arg_beneficiary,
			arg_authorizer,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsTopUpKeep(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.TopUpKeep(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallTopUpKeep(
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

func tsTopUpNu(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.TopUpNu(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallTopUpNu(
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

func tsTransferGovernance(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	arg_newGuvnor, err := chainutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg_newGuvnor, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.TransferGovernance(
			arg_newGuvnor,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallTransferGovernance(
			arg_newGuvnor,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsUnstakeAll(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.UnstakeAll(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUnstakeAll(
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

func tsUnstakeKeep(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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
		transaction, err = contract.UnstakeKeep(
			arg_stakingProvider,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUnstakeKeep(
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

func initializeTokenStaking(c *cli.Context) (*contract.TokenStaking, error) {
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

	address := common.HexToAddress(config.ContractAddresses["TokenStaking"])

	return contract.NewTokenStaking(
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
