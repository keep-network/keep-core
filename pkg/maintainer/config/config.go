package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/mitchellh/mapstructure"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// TODO: Should each maintainer have its own config or should there be a single
//       config for all the maintainers.

// Config stores configuration for the maintainer.
type Config struct {
	Bitcoin  bitcoin.Config
	Ethereum ethereum.Config
}

func (c *Config) ReadConfig(configFilePath string) error {
	// Read configuration from a file if the config file path is set.
	if configFilePath != "" {
		if err := readConfigFile(configFilePath); err != nil {
			return fmt.Errorf(
				"unable to load config (file: [%s]): [%w]",
				configFilePath,
				err,
			)
		}
	}

	// Unmarshal config based on loaded config file and command-line flags.
	if err := unmarshalConfig(c); err != nil {
		return fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return nil
}


// unmarshalConfig unmarshals config with viper from config file and command-line
// flags into a struct.
func unmarshalConfig(config *Config) error {
	if err := viper.Unmarshal(
		config,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
				mapstructure.TextUnmarshallerHookFunc(),
			),
		),
	); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return nil
}

// readConfigFile uses viper to read configuration from a config file. The config file
// is not mandatory, if the path is
func readConfigFile(configFilePath string) error {
	// Read configuration from a file, located in `configFilePath`.
	viper.SetConfigFile(configFilePath)

	// Read configuration.
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf(
			"failed to read configuration from file [%s]: %w",
			configFilePath,
			err,
		)
	}

	return nil
}