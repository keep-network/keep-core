package cmd

import (
	"fmt"

	"github.com/keep-network/keep-core/cmd/flag"
	"github.com/urfave/cli"
)

const (
	blockFlag        string = "block"
	blockShort       string = "b"
	transactionFlag  string = "transaction"
	transactionShort string = "t"
	submitFlag       string = "submit"
	submitShort      string = "s"
	valueFlag        string = "value"
	valueShort       string = "v"
)

var (
	transactionFlagValue *flag.TransactionHash = &flag.TransactionHash{}
	blockFlagValue       *flag.Uint256         = &flag.Uint256{}
	valueFlagValue       *flag.Uint256         = &flag.Uint256{}
)

// AvailableCommands is the exported list of generated commands that can be
// installed on a CLI app. Generated contract command files set up init
// functions that add the contract's command and subcommands to this global
// variable, and any top-level command that wishes to include these commands can
// reference this variable and expect it to contain all generated contract
// commands.
var AvailableCommands []cli.Command

var (
	constFlags = []cli.Flag{
		&cli.GenericFlag{
			Name:  blockFlag + ", " + blockShort,
			Usage: "Retrieve the result of calling this method on `BLOCK`.",
			Value: blockFlagValue,
		},
		&cli.GenericFlag{
			Name:  transactionFlag + ", " + transactionShort,
			Usage: "Retrieve the already-evaluated result of this method in `TRANSACTION`",
			Value: transactionFlagValue,
		},
	}
	nonConstFlags = append(
		[]cli.Flag{
			&cli.BoolFlag{
				Name:  submitFlag + ", " + submitShort,
				Usage: "Submit this call as a gas-spending network transaction.",
			},
		},
		constFlags...,
	)
	payableFlags = append(
		[]cli.Flag{
			&cli.GenericFlag{
				Name:  valueFlag + ", " + valueShort,
				Usage: "Send `VALUE` ether with this call.",
				Value: valueFlagValue,
			},
		},
		nonConstFlags...,
	)
)

type composableArgChecker cli.BeforeFunc

func (cac composableArgChecker) andThen(nextChecker composableArgChecker) composableArgChecker {
	return func(c *cli.Context) error {
		cacErr := cac(c)
		if cacErr != nil {
			return cacErr
		}

		return nextChecker(c)
	}
}

var (
	valueArgChecker composableArgChecker = func(c *cli.Context) error {
		if !c.IsSet(valueFlag) {
			return fmt.Errorf("expected value for this payable method")
		}

		return nil
	}
	submittedArgChecker composableArgChecker = func(c *cli.Context) error {
		if c.Bool(submitFlag) {
			if c.IsSet(blockFlag) {
				return fmt.Errorf("cannot specify --block for a submitted transaction")
			}
			if c.IsSet(transactionFlag) {
				return fmt.Errorf("cannot specify --transaction for a submitted transaction")
			}
		}

		return nil
	}
)

var (
	nonConstArgsChecker = submittedArgChecker
	payableArgsChecker  = nonConstArgsChecker.andThen(valueArgChecker)
)

func argCountChecker(expectedArgCount int) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if c.NArg() != expectedArgCount {
			return fmt.Errorf(
				"Expected [%v] arguments but got [%v]",
				expectedArgCount,
				c.NArg(),
			)
		}

		return nil
	}
}
