package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/cmd"
	"github.com/urfave/cli"
)

const (
	defaultConfigPath string = "./config.toml"
	defaultGroupSize  int    = 10
	defaultThreshold  int    = 4
)

var (
	cmds       []cli.Command
	configPath string
	// GroupSize ...
	GroupSize int
	// Threshold ...
	Threshold int
	// Version is the semantic version (added at compile time)  See scripts/version.sh
	Version string
	// Revision is the git commit id (added at compile time)
	Revision string
)

func init() {
	//TODO: Remove Version and Revision when build process auto-populates these values
	Version = "0.0.1"
	Revision = "deadbeef"

	cmds = []cli.Command{
		{
			Name:        "validate-config",
			Usage:       fmt.Sprintf("Example: %s validate-config", path.Base(os.Args[0])),
			Description: "Validates config file",
			Action:      cmd.ValidateConfig,
		},
		{
			Name:        "smoke-test",
			Usage:       fmt.Sprintf("Usage:   %s smoke-test -g <GROUP_SIZE> -t <THRESHOLD>", path.Base(os.Args[0])),
			Description: "Simulate DKG and verify group's threshold signature",
			Action:      cmd.SmokeTest,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:        "group-size,g",
					Value:       defaultGroupSize,
					Destination: &GroupSize,
					EnvVar:      "GROUP_SIZE",
					Usage:       "optionally, specify the `GROUP_SIZE` environment variable",
				},
				&cli.IntFlag{
					Name:        "threshold,t",
					Value:       defaultThreshold,
					Destination: &Threshold,
					EnvVar:      "THRESHOLD",
					Usage:       "optionally, specify the `THRESHOLD` environment variable",
				},
			},
		},
	}
}

func main() {
	err := bls.Init(bls.CurveSNARK1)
	if err != nil {
		log.Fatal("Failed to initialize BLS.", err)
	}

	cliApp := &cli.App{
		Name:        path.Base(os.Args[0]),
		Usage:       "CLI for The Keep Network",
		Version:     fmt.Sprintf("%s (revision %s)", Version, Revision),
		Description: "Command line interface (CLI) for running a Keep provider",
		Compiled:    time.Now(),
		Authors: []cli.Author{
			cli.Author{
				Name:  "Keep Network",
				Email: "info@keep.network",
			},
		},
		Copyright: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config,c",
				Value:       defaultConfigPath,
				Destination: &configPath,
				EnvVar:      "CONFIG_PATH",
				Usage:       "optionally, specify the `CONFIG_PATH` environment variable",
			},
		},
		Commands: cmds,
	}

	cliApp.Run(os.Args)
}
