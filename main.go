package main

import (
	"log"
	"os"

	"fmt"
	"path"
	"time"

	"github.com/dfinity/go-dfinity-crypto/bls"
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

	if err := bls.Init(bls.CurveSNARK1); err != nil {
		log.Fatal("Failed to initialize BLS.", err)
	}

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "CLI for The Keep Network"
	app.Description = "Command line interface (CLI) for running a Keep provider"
	app.Copyright = "" //TODO: Insert copyright printInfo later
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
		{
			Name: "smoke-test",
			Usage: `Simulates Distributed Key Generation (DKG) and signature verification
       --group-size 9:  Threshold relay group size; default: 10
       --threshold 9:   Minimun number of group members required to process requests; default 4
`,
			Description: "simulate Distributed Key Generation (DKG) and verify group's threshold signature",
			Action:      cmd.SmokeTest,
			Flags:       cmd.SmokeTestFlags,
		},
		{
			Name: "start",
			Usage: `Starts the Keep client in the foreground. Currently this consists of the
            threshold relay client for the Keep random beacon and the validator client
            for the Keep random beacon.
       --bootstrap             Indicates that this node is a bootstrap server
       --port                  Port this node will be listening on
       --preferred-ip-address  Preferred outbound IPv4 IP Address
       --disable-provider:     Disables the Keep provider client; default false
`,
			Description: "starts the Keep client in the foreground",
			Action:      cmd.StartNode,
			Flags:       cmd.StartFlags,
		},
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
