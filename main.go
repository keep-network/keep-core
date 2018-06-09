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

	err := newApp(version, revision).Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func newApp(version, revision string) *cli.App {

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "CLI for The Keep Network"
	app.Description = "Command line interface (CLI) for running a Keep provider"
	app.Copyright = "" //TODO: Insert copyright printInfo later
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Keep Network",
			Email: "printInfo@keep.network",
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
			Name:     "print-info",
			Usage:    "prints keep client information",
			Category: "keep client information",
			Action: func(c *cli.Context) error {
				printInfo(c)
				return nil
			},
		},
		{
			Name:        "smoke-test",
			Usage:       "Simulates DKG and signature verification",
			Description: "simulate Distributed Key Generation (DKG) and verify group's threshold signature",
			Action:      cmd.SmokeTest,
			Flags:       cmd.SmokeTestFlags,
		},
	}

	cli.AppHelpTemplate = fmt.Sprintf(`%s
ENVIRONMENT VARIABLES:
   KEEP_ETHEREUM_PASSWORD    keep client password

`, cli.AppHelpTemplate)

	return app
}

func printInfo(c *cli.Context) {
	fmt.Printf("Keep client: %s\n\n"+
		"Description: %s\n"+
		"version:     %s\n"+
		"revision:    %s\n"+
		"Config Path: %s\n",
		c.App.Name, c.App.Description, version, revision, c.GlobalString("config"))
}
