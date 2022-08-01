package cmd

import (
	"github.com/keep-network/keep-core/config"
)

// Path to the configuration file. The value is set with `--config` command-line flag.
var configFilePath string

// The client configuration for command-line execution. The content is read from
// a config file and command-line flags.
var clientConfig = &config.Config{}
