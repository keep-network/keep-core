package main

import (
	"os"

	"fmt"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/logging"
	"github.com/keep-network/keep-core/build"
	"github.com/keep-network/keep-core/cmd"
	"github.com/keep-network/keep-core/config"
)

var logger = log.Logger("keep-main")

func main() {
	err := logging.Configure(os.Getenv(config.LogLevelEnvVariable))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to configure logging: [%v]\n", err)
	}

	build.LogVersion()

	rootCmd := cmd.Initialize(build.Version, build.Revision)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
		logger.Info()
	}
}
