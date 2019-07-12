package main

import (
	"log"
	"os"

	"fmt"
	"path"
	"time"

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

	err := setUpLogging(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}

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
   LOG_LEVEL                 space-delimited set of log level directives; set to
                             "help" for help

`, cli.AppHelpTemplate)

	err = app.Run(os.Args)
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
