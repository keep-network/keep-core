package main

import (
	"log"
	"os"
	"strings"

	"fmt"
	"path"
	"time"

	logging "github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/cmd"
	"github.com/urfave/cli"
)

const defaultConfigPath = "./config.toml"

var (
	version  string
	revision string

	configPath string
)

func main() {
	if version == "" {
		version = "unknown"
	}
	if revision == "" {
		revision = "unknown"
	}

	readLoggingSetup()

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "CLI for The Keep Network"
	app.Description = "Command line interface (CLI) for running a Keep provider"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Keep Network",
			Email: "info@keep.network",
		},
	}
	app.Version = fmt.Sprintf("%s (revision %s)", version, revision)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,c",
			Value:       defaultConfigPath,
			Destination: &configPath,
			Usage:       "full path to the configuration file",
		},
	}
	app.Commands = []cli.Command{
		cmd.SmokeTestCommand,
		cmd.StartCommand,
		cmd.PingCommand,
		cmd.EthereumCommand,
		{
			Name:  "print-info",
			Usage: "Prints keep client information",
			Action: func(c *cli.Context) error {
				printInfo(c)
				return nil
			},
		},
	}

	cli.AppHelpTemplate = fmt.Sprintf(`%s
ENVIRONMENT VARIABLES:
   KEEP_ETHEREUM_PASSWORD    keep client password

`, cli.AppHelpTemplate)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func printInfo(c *cli.Context) {
	fmt.Printf("Keep client: %s\n\n"+
		"Description: %s\n"+
		"version:     %s\n"+
		"revision:    %s\n"+
		"Config Path: %s\n",
		c.App.Name,
		c.App.Description,
		version,
		revision,
		c.GlobalString("config"),
	)
}

func readLoggingSetup() {
	// Single string with space-delimited directives setting log level for each
	// subsystem, with = separating subsystem from log level. Example:
	//
	// "relay=debug bootstrap=info swarm2=error"
	//
	// Can also be a single log level for all subsystems:
	//
	// "info"
	//
	// If blank or unset, subsystems are left in their default initial state.
	levelDirectiveString := os.Getenv("LOG_LEVEL")

	// Default to info logs for keep.
	if len(levelDirectiveString) == 0 {
		levelDirectiveString = "keep*=info"
	}

	levelDirectives := strings.Split(levelDirectiveString, " ")
	for _, directive := range levelDirectives {
		err := evaluateLevelDirective(directive)
		if err != nil {
			log.Fatalf(
				"Failed to parse log level directive [%s]: [%v]\n"+
					"Directives can be any of:\n"+
					" - a global log level, e.g. 'debug'\n"+
					" - a subsystem=level pair, e.g. 'keep-relay=info'\n"+
					" - a subsystem*=level prefix pair, e.g. 'keep*=warn'\n",
				directive,
				err,
			)
		}
	}
}

// Takes a levelDirective that can have one of three formats:
//
//     <log-level> |
// 	   <subsystem>=<log-level> |
// 	   <subsystem-prefix>*=<log-level>
//
// In the first form, the given log-level is set on all subsystems.
//
// In the second form, the given log-level is set on the given subsystem.
//
// In the third form, the given log-level is set on any subsystem that starts
// with the given subsystem-prefix.
//
// Supported log levels are as per the ipfs/go-logging library.
func evaluateLevelDirective(levelDirective string) error {
	splitLevel := strings.Split(levelDirective, "=")

	switch len(splitLevel) {
	case 1:
		level := splitLevel[0]

		err := logging.SetLogLevel("*", level)
		if err != nil {
			return err
		}

	case 2:
		levelSubsystem := splitLevel[0]
		level := splitLevel[1]

		if strings.HasSuffix(levelSubsystem, "*") {
			subsystemPrefix := strings.TrimSuffix(levelSubsystem, "*")
			// Wildcard suffix, check for matching subsystems.
			for _, subsystem := range logging.GetSubsystems() {
				if strings.HasPrefix(subsystem, subsystemPrefix) {
					err := logging.SetLogLevel(subsystem, level)
					if err != nil {
						return err
					}
				}
			}
		} else {
			return logging.SetLogLevel(levelSubsystem, level)
		}

	default:
		return fmt.Errorf("more than two =-delimited components in directive")
	}

	return nil
}
