package main

import (
	"os"

	"fmt"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/logging"
	"github.com/keep-network/keep-core/cmd"
	"github.com/keep-network/keep-core/config"
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

	rootCmd := cmd.Initialize(version, revision)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
