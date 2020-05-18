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

var TokenGrantCommand cli.Command

var tokenGrantDescription = `The token-grant command allows calling the TokenGrant contract on an
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
		Name:        "token-grant",
		Usage:       `Provides access to the TokenGrant contract.`,
		Description: tokenGrantDescription,
		Subcommands: []cli.Command{{
			Name:      "grant-stakes",
			Usage:     "Calls the constant method grantStakes on the TokenGrant contract.",
			ArgsUsage: "[arg0] ",
			Action:    tgGrantStakes,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "grant-indices",
			Usage:     "Calls the constant method grantIndices on the TokenGrant contract.",
			ArgsUsage: "[arg0] [arg1] ",
			Action:    tgGrantIndices,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "grantees-to-operators",
			Usage:     "Calls the constant method granteesToOperators on the TokenGrant contract.",
			ArgsUsage: "[arg0] [arg1] ",
			Action:    tgGranteesToOperators,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "grants",
			Usage:     "Calls the constant method grants on the TokenGrant contract.",
			ArgsUsage: "[arg0] ",
			Action:    tgGrants,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "balances",
			Usage:     "Calls the constant method balances on the TokenGrant contract.",
			ArgsUsage: "[arg0] ",
			Action:    tgBalances,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "withdrawable",
			Usage:     "Calls the constant method withdrawable on the TokenGrant contract.",
			ArgsUsage: "[_id] ",
			Action:    tgWithdrawable,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-grant-stake-details",
			Usage:     "Calls the constant method getGrantStakeDetails on the TokenGrant contract.",
			ArgsUsage: "[operator] ",
			Action:    tgGetGrantStakeDetails,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-grant-unlocking-schedule",
			Usage:     "Calls the constant method getGrantUnlockingSchedule on the TokenGrant contract.",
			ArgsUsage: "[_id] ",
			Action:    tgGetGrantUnlockingSchedule,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-grants",
			Usage:     "Calls the constant method getGrants on the TokenGrant contract.",
			ArgsUsage: "[_granteeOrGrantManager] ",
			Action:    tgGetGrants,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "token",
			Usage:     "Calls the constant method token on the TokenGrant contract.",
			ArgsUsage: "",
			Action:    tgToken,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "available-to-stake",
			Usage:     "Calls the constant method availableToStake on the TokenGrant contract.",
			ArgsUsage: "[_grantId] ",
			Action:    tgAvailableToStake,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "balance-of",
			Usage:     "Calls the constant method balanceOf on the TokenGrant contract.",
			ArgsUsage: "[_owner] ",
			Action:    tgBalanceOf,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-grant",
			Usage:     "Calls the constant method getGrant on the TokenGrant contract.",
			ArgsUsage: "[_id] ",
			Action:    tgGetGrant,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-grantee-operators",
			Usage:     "Calls the constant method getGranteeOperators on the TokenGrant contract.",
			ArgsUsage: "[grantee] ",
			Action:    tgGetGranteeOperators,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "num-grants",
			Usage:     "Calls the constant method numGrants on the TokenGrant contract.",
			ArgsUsage: "",
			Action:    tgNumGrants,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "stake-balance-of",
			Usage:     "Calls the constant method stakeBalanceOf on the TokenGrant contract.",
			ArgsUsage: "[_address] ",
			Action:    tgStakeBalanceOf,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "unlocked-amount",
			Usage:     "Calls the constant method unlockedAmount on the TokenGrant contract.",
			ArgsUsage: "[_id] ",
			Action:    tgUnlockedAmount,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "cancel-stake",
			Usage:     "Calls the method cancelStake on the TokenGrant contract.",
			ArgsUsage: "[_operator] ",
			Action:    tgCancelStake,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "revoke",
			Usage:     "Calls the method revoke on the TokenGrant contract.",
			ArgsUsage: "[_id] ",
			Action:    tgRevoke,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "authorize-staking-contract",
			Usage:     "Calls the method authorizeStakingContract on the TokenGrant contract.",
			ArgsUsage: "[_stakingContract] ",
			Action:    tgAuthorizeStakingContract,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "receive-approval",
			Usage:     "Calls the method receiveApproval on the TokenGrant contract.",
			ArgsUsage: "[_from] [_amount] [_token] [_extraData] ",
			Action:    tgReceiveApproval,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(4))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "stake",
			Usage:     "Calls the method stake on the TokenGrant contract.",
			ArgsUsage: "[_id] [_stakingContract] [_amount] [_extraData] ",
			Action:    tgStake,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(4))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "undelegate-revoked",
			Usage:     "Calls the method undelegateRevoked on the TokenGrant contract.",
			ArgsUsage: "[_operator] ",
			Action:    tgUndelegateRevoked,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "cancel-revoked-stake",
			Usage:     "Calls the method cancelRevokedStake on the TokenGrant contract.",
			ArgsUsage: "[_operator] ",
			Action:    tgCancelRevokedStake,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "undelegate",
			Usage:     "Calls the method undelegate on the TokenGrant contract.",
			ArgsUsage: "[_operator] ",
			Action:    tgUndelegate,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw",
			Usage:     "Calls the method withdraw on the TokenGrant contract.",
			ArgsUsage: "[_id] ",
			Action:    tgWithdraw,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "recover-stake",
			Usage:     "Calls the method recoverStake on the TokenGrant contract.",
			ArgsUsage: "[_operator] ",
			Action:    tgRecoverStake,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw-revoked",
			Usage:     "Calls the method withdrawRevoked on the TokenGrant contract.",
			ArgsUsage: "[_id] ",
			Action:    tgWithdrawRevoked,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func tgGrantStakes(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

	result, err := contract.GrantStakesAtBlock(
		arg0,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgGrantIndices(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

	result, err := contract.GrantIndicesAtBlock(
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

func tgGranteesToOperators(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

	result, err := contract.GranteesToOperatorsAtBlock(
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

func tgGrants(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

	result, err := contract.GrantsAtBlock(
		arg0,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgBalances(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

	result, err := contract.BalancesAtBlock(
		arg0,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgWithdrawable(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}
	_id, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _id, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.WithdrawableAtBlock(
		_id,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgGetGrantStakeDetails(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

	result, err := contract.GetGrantStakeDetailsAtBlock(
		operator,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgGetGrantUnlockingSchedule(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}
	_id, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _id, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetGrantUnlockingScheduleAtBlock(
		_id,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgGetGrants(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}
	_granteeOrGrantManager, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _granteeOrGrantManager, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetGrantsAtBlock(
		_granteeOrGrantManager,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgToken(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

func tgAvailableToStake(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}
	_grantId, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _grantId, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.AvailableToStakeAtBlock(
		_grantId,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgBalanceOf(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}
	_owner, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _owner, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.BalanceOfAtBlock(
		_owner,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgGetGrant(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}
	_id, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _id, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetGrantAtBlock(
		_id,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgGetGranteeOperators(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}
	grantee, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter grantee, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetGranteeOperatorsAtBlock(
		grantee,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgNumGrants(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}

	result, err := contract.NumGrantsAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgStakeBalanceOf(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

	result, err := contract.StakeBalanceOfAtBlock(
		_address,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tgUnlockedAmount(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}
	_id, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _id, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.UnlockedAmountAtBlock(
		_id,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func tgCancelStake(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

func tgRevoke(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}

	_id, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _id, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Revoke(
			_id,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRevoke(
			_id,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tgAuthorizeStakingContract(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}

	_stakingContract, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _stakingContract, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.AuthorizeStakingContract(
			_stakingContract,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallAuthorizeStakingContract(
			_stakingContract,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tgReceiveApproval(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

	_amount, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _amount, a uint256, from passed value %v",
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
			_amount,
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
			_amount,
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

func tgStake(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}

	_id, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _id, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	_stakingContract, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _stakingContract, a address, from passed value %v",
			c.Args()[1],
		)
	}

	_amount, err := hexutil.DecodeBig(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _amount, a uint256, from passed value %v",
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
		transaction, err = contract.Stake(
			_id,
			_stakingContract,
			_amount,
			_extraData,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallStake(
			_id,
			_stakingContract,
			_amount,
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

func tgUndelegateRevoked(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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
		transaction, err = contract.UndelegateRevoked(
			_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallUndelegateRevoked(
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

func tgCancelRevokedStake(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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
		transaction, err = contract.CancelRevokedStake(
			_operator,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallCancelRevokedStake(
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

func tgUndelegate(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

func tgWithdraw(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}

	_id, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _id, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Withdraw(
			_id,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallWithdraw(
			_id,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tgRecoverStake(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
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

func tgWithdrawRevoked(c *cli.Context) error {
	contract, err := initializeTokenGrant(c)
	if err != nil {
		return err
	}

	_id, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _id, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.WithdrawRevoked(
			_id,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallWithdrawRevoked(
			_id,
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

func initializeTokenGrant(c *cli.Context) (*contract.TokenGrant, error) {
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

	checkInterval := cmd.DefaultMiningCheckInterval
	maxGasPrice := cmd.DefaultMaxGasPrice
	if config.MiningCheckInterval != 0 {
		checkInterval = time.Duration(config.MiningCheckInterval) * time.Second
	}
	if config.MaxGasPrice != 0 {
		maxGasPrice = new(big.Int).SetUint64(config.MaxGasPrice)
	}

	miningWaiter := ethutil.NewMiningWaiter(client, checkInterval, maxGasPrice)

	address := common.HexToAddress(config.ContractAddresses["TokenGrant"])

	return contract.NewTokenGrant(
		address,
		key,
		client,
		ethutil.NewNonceManager(key.Address, client),
		miningWaiter,
		&sync.Mutex{},
	)
}
