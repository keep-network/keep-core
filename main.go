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
	joinedLevelString := os.Getenv("LOG_LEVEL")

	// Nothing to do if the env var is empty.
	if len(joinedLevelString) == 0 {
		return
	}

	levelStrings := strings.Split(joinedLevelString, " ")

	// If there is only one directive and it has no = in it, treat it as a
	// global log level.
	if len(levelStrings) == 1 && !strings.Contains(levelStrings[0], "=") {
		level := levelStrings[0]
		err := logging.SetLogLevel("*", level)
		if err != nil {
			log.Fatalf("failed to parse log level [%s]: [%v]", level, err)
		}

		return
	}

	// If we're here, we want to handle subsystem=level pairs.
	for _, subsystemPair := range levelStrings {
		splitLevel := strings.Split(subsystemPair, "=")
		if len(splitLevel) != 2 {
			log.Fatalf(
				"expected string [%s] to have format subsystem=loglevel",
				splitLevel,
			)
		}

		subsystem := splitLevel[0]
		level := splitLevel[1]
		err := logging.SetLogLevel(subsystem, level)
		if err != nil {
			log.Fatalf("failed to parse log level [%s]: [%v]", level, err)
		}
	}
}
