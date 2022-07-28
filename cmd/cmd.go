package cmd

import (
	"github.com/keep-network/keep-core/config"
)

const defaultConfigPath = "./config.toml"

var configFilePath string

// var clientConfig *config.Config = config.NewConfig()
var clientConfig = &config.Config{}
