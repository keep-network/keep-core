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
}

var Config Conf

func GetOptions() bool {
	var configFile string

	flag.StringVar(&configFile, "config", "", "Configuration file")
	flag.StringVar(&Config.AppEnv, "app-env", "development", "Runtime environment")
	flag.Int("p2p-listen-port", 0, "Project Root directory (required). Must be absolute path")
	flag.Bool("enable-p2p-encryption", false, "Enable secure IO")
	flag.Int64("id-generation-seed", 0, "Random seed for ID generation")
	flag.BoolVar(&Config.LogTimeTrack, "log-timetrack", true, "Enable or disable logging of utils/TimeTrack() (For benchmarking/debugging)")
	flag.BoolVar(&Config.LogFullStackTrace, "log-full-stack-trace", false, "Print version information and exit")
	flag.BoolVar(&Config.LogDebugInfo, "log-debug-info", false, "Whether to log debug output to the log (set to true for debug purposes)")
	flag.BoolVar(&Config.LogDebugInfoForTests, "log-debug-info-for-tests", true, "Whether to log debug output to the log when running tests (set to true for debug purposes)")

	flag.Parse()

	if Config.P2pListenPort == 0 {
		log.Fatal("Please provide a port to bind on with --port")
	}

	if configFile != "" {
		if _, err := toml.DecodeFile(configFile, &Config); err != nil {
			HandlePanic(errors.Wrap(err, "unable to read config file"))
		}
	}
	return true
}
