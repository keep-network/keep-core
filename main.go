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

var (
	commands   []cli.Command
	configPath string

	// Version is the semantic version (added at compile time)  See scripts/version.sh
	Version string
	// Revision is the git commit id (added at compile time)
	Revision string
)

func init() {
	//TODO: Remove Version and Revision when build process auto-populates these values
	Version = "0.0.1"
	Revision = "deadbeef"
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
				Value:       cmd.DefaultConfigPath,
				Destination: &configPath,
				EnvVar:      "CONFIG_PATH",
				Usage:       "optionally, specify the `CONFIG_PATH` environment variable",
			},
		},
		Commands: cmd.Commands,
	}

	cliApp.Run(os.Args)
}
