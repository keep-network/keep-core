// Code generated - DO NOT EDIT.
// This file is a generated command and any manual changes will be lost.

package cmd

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/chain/gen/contract"

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
			Name:      "beneficiary-of",
			Usage:     "Calls the constant method beneficiaryOf on the TokenStaking contract.",
			ArgsUsage: "[_operator] ",
			Action:    tsBeneficiaryOf,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "minimum-stake",
			Usage:     "Calls the constant method minimumStake on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsMinimumStake,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "minimum-stake-schedule",
			Usage:     "Calls the constant method minimumStakeSchedule on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsMinimumStakeSchedule,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "owner-of",
			Usage:     "Calls the constant method ownerOf on the TokenStaking contract.",
			ArgsUsage: "[_operator] ",
			Action:    tsOwnerOf,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "balance-of",
			Usage:     "Calls the constant method balanceOf on the TokenStaking contract.",
			ArgsUsage: "[_address] ",
			Action:    tsBalanceOf,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "minimum-stake-steps",
			Usage:     "Calls the constant method minimumStakeSteps on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsMinimumStakeSteps,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "owner-operators",
			Usage:     "Calls the constant method ownerOperators on the TokenStaking contract.",
			ArgsUsage: "[arg0] [arg1] ",
			Action:    tsOwnerOperators,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-authority-source",
			Usage:     "Calls the constant method getAuthoritySource on the TokenStaking contract.",
			ArgsUsage: "[operatorContract] ",
			Action:    tsGetAuthoritySource,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "initialization-period",
			Usage:     "Calls the constant method initializationPeriod on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsInitializationPeriod,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-authorized-for-operator",
			Usage:     "Calls the constant method isAuthorizedForOperator on the TokenStaking contract.",
			ArgsUsage: "[_operator] [_operatorContract] ",
			Action:    tsIsAuthorizedForOperator,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "minimum-stake-base",
			Usage:     "Calls the constant method minimumStakeBase on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsMinimumStakeBase,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "active-stake",
			Usage:     "Calls the constant method activeStake on the TokenStaking contract.",
			ArgsUsage: "[_operator] [_operatorContract] ",
			Action:    tsActiveStake,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "maximum-lock-duration",
			Usage:     "Calls the constant method maximumLockDuration on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsMaximumLockDuration,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "eligible-stake",
			Usage:     "Calls the constant method eligibleStake on the TokenStaking contract.",
			ArgsUsage: "[_operator] [_operatorContract] ",
			Action:    tsEligibleStake,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-locks",
			Usage:     "Calls the constant method getLocks on the TokenStaking contract.",
			ArgsUsage: "[operator] ",
			Action:    tsGetLocks,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-stake-locked",
			Usage:     "Calls the constant method isStakeLocked on the TokenStaking contract.",
			ArgsUsage: "[operator] ",
			Action:    tsIsStakeLocked,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "undelegation-period",
			Usage:     "Calls the constant method undelegationPeriod on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsUndelegationPeriod,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "authorizer-of",
			Usage:     "Calls the constant method authorizerOf on the TokenStaking contract.",
			ArgsUsage: "[_operator] ",
			Action:    tsAuthorizerOf,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "registry",
			Usage:     "Calls the constant method registry on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsRegistry,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-delegation-info",
			Usage:     "Calls the constant method getDelegationInfo on the TokenStaking contract.",
			ArgsUsage: "[_operator] ",
			Action:    tsGetDelegationInfo,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "has-minimum-stake",
			Usage:     "Calls the constant method hasMinimumStake on the TokenStaking contract.",
			ArgsUsage: "[staker] [operatorContract] ",
			Action:    tsHasMinimumStake,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "operators",
			Usage:     "Calls the constant method operators on the TokenStaking contract.",
			ArgsUsage: "[arg0] ",
			Action:    tsOperators,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "operators-of",
			Usage:     "Calls the constant method operatorsOf on the TokenStaking contract.",
			ArgsUsage: "[_address] ",
			Action:    tsOperatorsOf,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "token",
			Usage:     "Calls the constant method token on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsToken,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "minimum-stake-schedule-start",
			Usage:     "Calls the constant method minimumStakeScheduleStart on the TokenStaking contract.",
			ArgsUsage: "",
			Action:    tsMinimumStakeScheduleStart,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "lock-stake",
			Usage:     "Calls the method lockStake on the TokenStaking contract.",
			ArgsUsage: "[operator] [duration] ",
			Action:    tsLockStake,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "receive-approval",
			Usage:     "Calls the method receiveApproval on the TokenStaking contract.",
			ArgsUsage: "[_from] [_value] [_token] [_extraData] ",
			Action:    tsReceiveApproval,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(4))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "undelegate",
			Usage:     "Calls the method undelegate on the TokenStaking contract.",
			ArgsUsage: "[_operator] ",
			Action:    tsUndelegate,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "authorize-operator-contract",
			Usage:     "Calls the method authorizeOperatorContract on the TokenStaking contract.",
			ArgsUsage: "[_operator] [_operatorContract] ",
			Action:    tsAuthorizeOperatorContract,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "cancel-stake",
			Usage:     "Calls the method cancelStake on the TokenStaking contract.",
			ArgsUsage: "[_operator] ",
			Action:    tsCancelStake,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "claim-delegated-authority",
			Usage:     "Calls the method claimDelegatedAuthority on the TokenStaking contract.",
			ArgsUsage: "[delegatedAuthoritySource] ",
			Action:    tsClaimDelegatedAuthority,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "recover-stake",
			Usage:     "Calls the method recoverStake on the TokenStaking contract.",
			ArgsUsage: "[_operator] ",
			Action:    tsRecoverStake,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "release-expired-lock",
			Usage:     "Calls the method releaseExpiredLock on the TokenStaking contract.",
			ArgsUsage: "[operator] [operatorContract] ",
			Action:    tsReleaseExpiredLock,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "undelegate-at",
			Usage:     "Calls the method undelegateAt on the TokenStaking contract.",
			ArgsUsage: "[_operator] [_undelegationTimestamp] ",
			Action:    tsUndelegateAt,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "unlock-stake",
			Usage:     "Calls the method unlockStake on the TokenStaking contract.",
			ArgsUsage: "[operator] ",
			Action:    tsUnlockStake,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func tsBeneficiaryOf(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.BeneficiaryOfAtBlock(
		_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsMinimumStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.MinimumStakeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsMinimumStakeSchedule(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.MinimumStakeScheduleAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsOwnerOf(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.OwnerOfAtBlock(
		_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsBalanceOf(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_address, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _address, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.BalanceOfAtBlock(
		_address,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsMinimumStakeSteps(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.MinimumStakeStepsAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsOwnerOperators(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg0, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg0, a address, from passed value %v",
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

	result, err := contract.OwnerOperatorsAtBlock(
		arg0,
		arg1,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetAuthoritySource(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
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

	result, err := contract.GetAuthoritySourceAtBlock(
		operatorContract,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsInitializationPeriod(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.InitializationPeriodAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsIsAuthorizedForOperator(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	_operatorContract, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operatorContract, a address, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.IsAuthorizedForOperatorAtBlock(
		_operator,
		_operatorContract,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsMinimumStakeBase(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.MinimumStakeBaseAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsActiveStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	_operatorContract, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operatorContract, a address, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.ActiveStakeAtBlock(
		_operator,
		_operatorContract,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsMaximumLockDuration(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.MaximumLockDurationAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsEligibleStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	_operatorContract, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operatorContract, a address, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.EligibleStakeAtBlock(
		_operator,
		_operatorContract,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetLocks(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetLocksAtBlock(
		operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsIsStakeLocked(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.IsStakeLockedAtBlock(
		operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsUndelegationPeriod(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.UndelegationPeriodAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsAuthorizerOf(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.AuthorizerOfAtBlock(
		_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsRegistry(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.RegistryAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsGetDelegationInfo(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetDelegationInfoAtBlock(
		_operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsHasMinimumStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	staker, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter staker, a address, from passed value %v",
			c.Args()[0],
		)
	}

	operatorContract, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operatorContract, a address, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.HasMinimumStakeAtBlock(
		staker,
		operatorContract,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsOperators(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	arg0, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter arg0, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.OperatorsAtBlock(
		arg0,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsOperatorsOf(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}
	_address, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _address, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.OperatorsOfAtBlock(
		_address,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsToken(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.TokenAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tsMinimumStakeScheduleStart(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	result, err := contract.MinimumStakeScheduleStartAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func tsLockStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	duration, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter duration, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.LockStake(
			operator,
			duration,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLockStake(
			operator,
			duration,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsReceiveApproval(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	_from, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _from, a address, from passed value %v",
			c.Args()[0],
		)
	}

	_value, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _value, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	_token, err := ethutil.AddressFromHex(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _token, a address, from passed value %v",
			c.Args()[2],
		)
	}

	_extraData, err := hexutil.Decode(c.Args()[3])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _extraData, a bytes, from passed value %v",
			c.Args()[3],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ReceiveApproval(
			_from,
			_value,
			_token,
			_extraData,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallReceiveApproval(
			_from,
			_value,
			_token,
			_extraData,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsUndelegate(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Undelegate(
			_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUndelegate(
			_operator,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsAuthorizeOperatorContract(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	_operatorContract, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operatorContract, a address, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.AuthorizeOperatorContract(
			_operator,
			_operatorContract,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallAuthorizeOperatorContract(
			_operator,
			_operatorContract,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsCancelStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.CancelStake(
			_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallCancelStake(
			_operator,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsClaimDelegatedAuthority(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	delegatedAuthoritySource, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter delegatedAuthoritySource, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ClaimDelegatedAuthority(
			delegatedAuthoritySource,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallClaimDelegatedAuthority(
			delegatedAuthoritySource,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsRecoverStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RecoverStake(
			_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRecoverStake(
			_operator,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsReleaseExpiredLock(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	operatorContract, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operatorContract, a address, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ReleaseExpiredLock(
			operator,
			operatorContract,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallReleaseExpiredLock(
			operator,
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

func tsUndelegateAt(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	_operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	_undelegationTimestamp, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _undelegationTimestamp, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UndelegateAt(
			_operator,
			_undelegationTimestamp,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUndelegateAt(
			_operator,
			_undelegationTimestamp,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tsUnlockStake(c *cli.Context) error {
	contract, err := initializeTokenStaking(c)
	if err != nil {
		return err
	}

	operator, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter operator, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.UnlockStake(
			operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUnlockStake(
			operator,
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

	address := common.HexToAddress(config.ContractAddresses["TokenStaking"])

	return contract.NewTokenStaking(
		address,
		key,
		client,
		&sync.Mutex{},
	)
}
