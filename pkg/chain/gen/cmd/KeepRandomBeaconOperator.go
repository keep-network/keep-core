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

var KeepRandomBeaconOperatorCommand cli.Command

var keepRandomBeaconOperatorDescription = `The keep-random-beacon-operator command allows calling the KeepRandomBeaconOperator contract on an
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
		Name:        "keep-random-beacon-operator",
		Usage:       `Provides access to the KeepRandomBeaconOperator contract.`,
		Description: keepRandomBeaconOperatorDescription,
		Subcommands: []cli.Command{{
			Name:      "entry-verification-fee",
			Usage:     "Calls the constant method entryVerificationFee on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboEntryVerificationFee,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "group-threshold",
			Usage:     "Calls the constant method groupThreshold on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGroupThreshold,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "has-minimum-stake",
			Usage:     "Calls the constant method hasMinimumStake on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[staker] ",
			Action:    krboHasMinimumStake,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "current-request-previous-entry",
			Usage:     "Calls the constant method currentRequestPreviousEntry on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboCurrentRequestPreviousEntry,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "dkg-submitter-reimbursement-fee",
			Usage:     "Calls the constant method dkgSubmitterReimbursementFee on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboDkgSubmitterReimbursementFee,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-group-member-rewards",
			Usage:     "Calls the constant method getGroupMemberRewards on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[groupPubKey] ",
			Action:    krboGetGroupMemberRewards,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "selected-participants",
			Usage:     "Calls the constant method selectedParticipants on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboSelectedParticipants,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "gas-price-ceiling",
			Usage:     "Calls the constant method gasPriceCeiling on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGasPriceCeiling,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-group-public-key",
			Usage:     "Calls the constant method getGroupPublicKey on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[groupIndex] ",
			Action:    krboGetGroupPublicKey,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-stale-group",
			Usage:     "Calls the constant method isStaleGroup on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[groupPubKey] ",
			Action:    krboIsStaleGroup,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "relay-entry-timeout",
			Usage:     "Calls the constant method relayEntryTimeout on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboRelayEntryTimeout,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "submitted-tickets",
			Usage:     "Calls the constant method submittedTickets on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboSubmittedTickets,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "dkg-gas-estimate",
			Usage:     "Calls the constant method dkgGasEstimate on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboDkgGasEstimate,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-first-active-group-index",
			Usage:     "Calls the constant method getFirstActiveGroupIndex on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGetFirstActiveGroupIndex,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "has-withdrawn-rewards",
			Usage:     "Calls the constant method hasWithdrawnRewards on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[operator] [groupIndex] ",
			Action:    krboHasWithdrawnRewards,
			Before:    cmd.ArgCountChecker(2),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-group-registered",
			Usage:     "Calls the constant method isGroupRegistered on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[groupPubKey] ",
			Action:    krboIsGroupRegistered,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "group-creation-fee",
			Usage:     "Calls the constant method groupCreationFee on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGroupCreationFee,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-entry-in-progress",
			Usage:     "Calls the constant method isEntryInProgress on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboIsEntryInProgress,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-group-selection-possible",
			Usage:     "Calls the constant method isGroupSelectionPossible on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboIsGroupSelectionPossible,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "number-of-groups",
			Usage:     "Calls the constant method numberOfGroups on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboNumberOfGroups,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "current-request-group-index",
			Usage:     "Calls the constant method currentRequestGroupIndex on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboCurrentRequestGroupIndex,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "current-request-start-block",
			Usage:     "Calls the constant method currentRequestStartBlock on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboCurrentRequestStartBlock,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-group-members",
			Usage:     "Calls the constant method getGroupMembers on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[groupPubKey] ",
			Action:    krboGetGroupMembers,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "group-selection-gas-estimate",
			Usage:     "Calls the constant method groupSelectionGasEstimate on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGroupSelectionGasEstimate,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "result-publication-block-step",
			Usage:     "Calls the constant method resultPublicationBlockStep on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboResultPublicationBlockStep,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "entry-verification-gas-estimate",
			Usage:     "Calls the constant method entryVerificationGasEstimate on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboEntryVerificationGasEstimate,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "group-member-base-reward",
			Usage:     "Calls the constant method groupMemberBaseReward on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGroupMemberBaseReward,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "group-profit-fee",
			Usage:     "Calls the constant method groupProfitFee on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGroupProfitFee,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "group-size",
			Usage:     "Calls the constant method groupSize on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGroupSize,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "ticket-submission-timeout",
			Usage:     "Calls the constant method ticketSubmissionTimeout on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboTicketSubmissionTimeout,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "add-service-contract",
			Usage:     "Calls the method addServiceContract on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[serviceContract] ",
			Action:    krboAddServiceContract,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "report-unauthorized-signing",
			Usage:     "Calls the method reportUnauthorizedSigning on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[groupIndex] [signedMsgSender] ",
			Action:    krboReportUnauthorizedSigning,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "sign",
			Usage:     "Calls the payable method sign on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[requestId] [previousEntry] ",
			Action:    krboSign,
			Before:    cli.BeforeFunc(cmd.PayableArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.PayableFlags,
		}, {
			Name:      "relay-entry",
			Usage:     "Calls the method relayEntry on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[_groupSignature] ",
			Action:    krboRelayEntry,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "genesis",
			Usage:     "Calls the payable method genesis on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboGenesis,
			Before:    cli.BeforeFunc(cmd.PayableArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.PayableFlags,
		}, {
			Name:      "remove-service-contract",
			Usage:     "Calls the method removeServiceContract on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[serviceContract] ",
			Action:    krboRemoveServiceContract,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "report-relay-entry-timeout",
			Usage:     "Calls the method reportRelayEntryTimeout on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "",
			Action:    krboReportRelayEntryTimeout,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "create-group",
			Usage:     "Calls the payable method createGroup on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[_newEntry] [submitter] ",
			Action:    krboCreateGroup,
			Before:    cli.BeforeFunc(cmd.PayableArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.PayableFlags,
		}, {
			Name:      "withdraw-group-member-rewards",
			Usage:     "Calls the method withdrawGroupMemberRewards on the KeepRandomBeaconOperator contract.",
			ArgsUsage: "[operator] [groupIndex] ",
			Action:    krboWithdrawGroupMemberRewards,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(2))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func krboEntryVerificationFee(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.EntryVerificationFeeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGroupThreshold(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.GroupThresholdAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboHasMinimumStake(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
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

	result, err := contract.HasMinimumStakeAtBlock(
		staker,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboCurrentRequestPreviousEntry(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.CurrentRequestPreviousEntryAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboDkgSubmitterReimbursementFee(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.DkgSubmitterReimbursementFeeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGetGroupMemberRewards(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}
	groupPubKey, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter groupPubKey, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetGroupMemberRewardsAtBlock(
		groupPubKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboSelectedParticipants(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.SelectedParticipantsAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGasPriceCeiling(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.GasPriceCeilingAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGetGroupPublicKey(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}
	groupIndex, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter groupIndex, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetGroupPublicKeyAtBlock(
		groupIndex,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboIsStaleGroup(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}
	groupPubKey, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter groupPubKey, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.IsStaleGroupAtBlock(
		groupPubKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboRelayEntryTimeout(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.RelayEntryTimeoutAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboSubmittedTickets(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.SubmittedTicketsAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboDkgGasEstimate(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.DkgGasEstimateAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGetFirstActiveGroupIndex(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.GetFirstActiveGroupIndexAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboHasWithdrawnRewards(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
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

	groupIndex, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter groupIndex, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	result, err := contract.HasWithdrawnRewardsAtBlock(
		operator,
		groupIndex,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboIsGroupRegistered(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}
	groupPubKey, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter groupPubKey, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.IsGroupRegisteredAtBlock(
		groupPubKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGroupCreationFee(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.GroupCreationFeeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboIsEntryInProgress(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.IsEntryInProgressAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboIsGroupSelectionPossible(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.IsGroupSelectionPossibleAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboNumberOfGroups(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.NumberOfGroupsAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboCurrentRequestGroupIndex(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.CurrentRequestGroupIndexAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboCurrentRequestStartBlock(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.CurrentRequestStartBlockAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGetGroupMembers(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}
	groupPubKey, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter groupPubKey, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetGroupMembersAtBlock(
		groupPubKey,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGroupSelectionGasEstimate(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.GroupSelectionGasEstimateAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboResultPublicationBlockStep(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.ResultPublicationBlockStepAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboEntryVerificationGasEstimate(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.EntryVerificationGasEstimateAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGroupMemberBaseReward(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.GroupMemberBaseRewardAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGroupProfitFee(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.GroupProfitFeeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboGroupSize(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.GroupSizeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func krboTicketSubmissionTimeout(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	result, err := contract.TicketSubmissionTimeoutAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func krboAddServiceContract(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	serviceContract, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter serviceContract, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.AddServiceContract(
			serviceContract,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallAddServiceContract(
			serviceContract,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krboReportUnauthorizedSigning(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	groupIndex, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter groupIndex, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	signedMsgSender, err := hexutil.Decode(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter signedMsgSender, a bytes, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ReportUnauthorizedSigning(
			groupIndex,
			signedMsgSender,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallReportUnauthorizedSigning(
			groupIndex,
			signedMsgSender,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krboSign(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	requestId, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter requestId, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	previousEntry, err := hexutil.Decode(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter previousEntry, a bytes, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Sign(
			requestId,
			previousEntry,
			cmd.ValueFlagValue.Uint)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallSign(
			requestId,
			previousEntry,
			cmd.ValueFlagValue.Uint, cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krboRelayEntry(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	_groupSignature, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _groupSignature, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RelayEntry(
			_groupSignature,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRelayEntry(
			_groupSignature,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krboGenesis(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Genesis(
			cmd.ValueFlagValue.Uint)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallGenesis(
			cmd.ValueFlagValue.Uint, cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krboRemoveServiceContract(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	serviceContract, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter serviceContract, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RemoveServiceContract(
			serviceContract,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRemoveServiceContract(
			serviceContract,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krboReportRelayEntryTimeout(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ReportRelayEntryTimeout()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallReportRelayEntryTimeout(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krboCreateGroup(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
	if err != nil {
		return err
	}

	_newEntry, err := hexutil.DecodeBig(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _newEntry, a uint256, from passed value %v",
			c.Args()[0],
		)
	}

	submitter, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter submitter, a address, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.CreateGroup(
			_newEntry,
			submitter,
			cmd.ValueFlagValue.Uint)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallCreateGroup(
			_newEntry,
			submitter,
			cmd.ValueFlagValue.Uint, cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func krboWithdrawGroupMemberRewards(c *cli.Context) error {
	contract, err := initializeKeepRandomBeaconOperator(c)
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

	groupIndex, err := hexutil.DecodeBig(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter groupIndex, a uint256, from passed value %v",
			c.Args()[1],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.WithdrawGroupMemberRewards(
			operator,
			groupIndex,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallWithdrawGroupMemberRewards(
			operator,
			groupIndex,
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

func initializeKeepRandomBeaconOperator(c *cli.Context) (*contract.KeepRandomBeaconOperator, error) {
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

	checkInterval := 30 * time.Second
	maxGasPrice := big.NewInt(50000000000) // 50 Gwei
	if config.MiningCheckInterval != 0 {
		checkInterval = time.Duration(config.MiningCheckInterval) * time.Second
	}
	if config.MaxGasPrice != 0 {
		maxGasPrice = new(big.Int).SetUint64(config.MaxGasPrice)
	}

	miningWaiter := ethutil.NewMiningWaiter(client, checkInterval, maxGasPrice)

	address := common.HexToAddress(config.ContractAddresses["KeepRandomBeaconOperator"])

	return contract.NewKeepRandomBeaconOperator(
		address,
		key,
		client,
		ethutil.NewNonceManager(key.Address, client),
		miningWaiter,
		&sync.Mutex{},
	)
}
