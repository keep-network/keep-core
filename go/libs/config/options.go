// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package config

import (
	"flag"
	"log"
	"github.com/pkg/errors"
	"github.com/BurntSushi/toml"
)

type Conf struct {
	AppEnv               string   `toml:"app_env"`
	P2pListenPort        int      `toml:"p2p_listen_port"`
	EnableP2pEncryption  bool     `toml:"enable_p2p_encryption"`
	IdGenerationSeed     int64    `toml:"id_generation_seed"`
	LogTimeTrack         bool     `toml:"log_timetrack"`
	LogFullStackTrace    bool     `toml:"log_full_stack_trace"`
	LogDebugInfo         bool     `toml:"log_debug_info"`
	LogDebugInfoForTests bool     `toml:"log_debug_info_for_tests"`
	Bootstrap            bootstrapNodes
}

type bootstrapNodes struct {
	Nodes []string
}

var Config Conf

func GetOptions() bool {
	var configFile string
	flag.StringVar(&configFile, "config", "", "Configuration file")
	flag.StringVar(&Config.AppEnv, "app-env", "development", "Runtime environment")
	flag.Int("p2p-listen-port", 7000, "p2p listen port")
	flag.Bool("enable-p2p-encryption", false, "Enable secure IO")
	flag.Int64("id-generation-seed", 0, "Random seed for ID generation")
	flag.BoolVar(&Config.LogTimeTrack, "log-timetrack", true, "Enable or disable logging of utils/TimeTrack() (For benchmarking/debugging)")
	flag.BoolVar(&Config.LogFullStackTrace, "log-full-stack-trace", false, "Print version information and exit")
	flag.BoolVar(&Config.LogDebugInfo, "log-debug-info", false, "Whether to log debug output to the log (set to true for debug purposes)")
	flag.BoolVar(&Config.LogDebugInfoForTests, "log-debug-info-for-tests", true, "Whether to log debug output to the log when running tests (set to true for debug purposes)")
	flag.Parse()
	// Load options in .toml config file into Config object
	if configFile != "" {
		if _, err := toml.DecodeFile(configFile, &Config); err != nil {
			HandlePanic(errors.Wrap(err, "Unable to read .toml config file"))
		}
	}
	// Data validations
	if Config.P2pListenPort == 7000 {
		log.Println("Using default --p2p-listen-port: 7000")
	}
	if len(Config.Bootstrap.Nodes) == 0 {
		example := `[bootstrap]
		nodes = [
			"/ip4/192.168.0.2/tcp/2701/ipfs/QmexAnfpHrhMmAC5UNQVS8iBuUUgDrMbMY17Cck2gKrqeX",
			"/ip4/192.168.0.3/tcp/2701/ipfs/Qmd3wzD2HWA95ZAs214VxnckwkwM4GHJyC6whKUCNQhNvW"
		]`
		log.Fatal("Please provide a bootstrap nodes in .toml config file like this:\n, %s", example)
	} else {
		log.Printf("Bootstrap Nodes: %+v\n", Config.Bootstrap.Nodes)
	}
	return true
}
