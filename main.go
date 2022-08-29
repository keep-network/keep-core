package main

import (
	"os"
	"path"

	"fmt"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/logging"
	"github.com/keep-network/keep-core/cmd"
	"github.com/keep-network/keep-core/config"
	"github.com/spf13/cobra"
)

var (
	version  string
	revision string

	logger = log.Logger("keep-main")
)

func main() {
	if version == "" {
		version = "unknown"
	}
	if revision == "" {
		revision = "unknown"
	}

	err := logging.Configure(os.Getenv(config.LogLevelEnvVariable))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to configure logging: [%v]\n", err)
	}

	var rootCmd = &cobra.Command{}

	rootCmd.Use = path.Base(os.Args[0])
	rootCmd.Short = "CLI for The Keep Network"
	rootCmd.Long = "Command line interface (CLI) for running a Keep provider"
	rootCmd.Version = fmt.Sprintf("%s (revision %s)", version, revision)
	rootCmd.TraverseChildren = true

	rootCmd.AddCommand(
		cmd.StartCommand,
		cmd.PingCommand,
		cmd.EthereumCommand,
	)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
