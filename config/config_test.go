package config

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"golang.org/x/exp/slices"
)

func TestReadConfigFromFile(t *testing.T) {
	filePaths := []string{
		"../test/config.toml",
		"../test/config.json",
		"../test/config.yaml",
	}

	var configReadTests = map[string]struct {
		readValueFunc func(*Config) interface{}
		expectedValue interface{}
	}{
		"Ethereum.URL": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URL },
			expectedValue: "ws://192.168.0.158:8546",
		},
		"Ethereum.KeyFile": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.Account.KeyFile },
			expectedValue: "/tmp/UTC--2018-03-11T01-37-33.202765887Z--c2a56884538778bacd91aa5bf343bf882c5fb18b",
		},
		"Ethereum.KeyFilePassword": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.Account.KeyFilePassword },
			expectedValue: "THIS IS TEST! Password should be defined in env variable or prompt",
		},
		"Ethereum.MiningCheckInterval": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.MiningCheckInterval },
			expectedValue: 33 * time.Second,
		},
		"Ethereum.RequestsPerSecondLimit": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.RequestsPerSecondLimit },
			expectedValue: 13,
		},
		"Ethereum.ConcurrencyLimit": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.ConcurrencyLimit },
			expectedValue: 56,
		},
		"Ethereum.MaxGasFeeCap": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.MaxGasFeeCap.Int },
			expectedValue: big.NewInt(148000000000),
		},
		"Ethereum.BalanceAlertThreshold": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.BalanceAlertThreshold.Int },
			expectedValue: big.NewInt(2300000000000000000),
		},
		"Ethereum.Developer - map": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"randombeacon":   "0xcf64c2a367341170cb4e09cf8c0ed137d8473ceb",
				"walletregistry": "0x143ba24e66fce8bca22f7d739f9a932c519b1c76",
				"tokenstaking":   "0xa363a197f1bbb8877f50350234e3f15fb4175457",
			},
		},
		"Developer - RandomBeacon": {
			readValueFunc: func(c *Config) interface{} {
				address, _ := c.Ethereum.ContractAddress(ethereum.RandomBeaconContractName)
				return address.String()
			},
			expectedValue: "0xcf64C2A367341170cb4E09cf8C0ED137d8473CEB",
		},
		"Developer.WalletRegistryAddress": {
			readValueFunc: func(c *Config) interface{} {
				address, _ := c.Ethereum.ContractAddress(ethereum.WalletRegistryContractName)
				return address.String()
			},
			expectedValue: "0x143Ba24e66FCe8bCA22F7D739F9a932c519B1C76",
		},
		"Ethereum.Developer - TokenStaking": {
			readValueFunc: func(c *Config) interface{} {
				address, _ := c.Ethereum.ContractAddress(ethereum.TokenStakingContractName)
				return address.String()
			},
			expectedValue: "0xa363a197F1BbB8877F50350234e3f15fB4175457",
		},
		"Network.Port": {
			readValueFunc: func(c *Config) interface{} { return c.LibP2P.Port },
			expectedValue: 27001,
		},
		"Network.Peers": {
			readValueFunc: func(c *Config) interface{} { return c.LibP2P.Peers },
			expectedValue: []string{
				"/ip4/127.0.0.1/tcp/3820/ipfs/16Uiu2HAmVZGi9bgF3w6C4TFVo9HjbCDyuecsQbzQnXQJvP5wBjkd",
				"/ip4/127.0.0.1/tcp/3819/ipfs/16Uiu2HAmVNfJs6t7bB3hYPTxtuKXdTdqxKYn9wKKfUiGwpCDjySM",
			},
		},
		"Network.AnnouncedAddresses": {
			readValueFunc: func(c *Config) interface{} { return c.LibP2P.AnnouncedAddresses },
			expectedValue: []string{
				"/dns4/example.com/tcp/3919",
				"/ip4/80.70.60.50/tcp/3919",
			},
		},
		"Network.DisseminationTime": {
			readValueFunc: func(c *Config) interface{} { return c.LibP2P.DisseminationTime },
			expectedValue: 76,
		},
		"Storage.DataDir": {
			readValueFunc: func(c *Config) interface{} { return c.Storage.DataDir },
			expectedValue: "/my/secure/location",
		},
		"Metrics.Port": {
			readValueFunc: func(c *Config) interface{} { return c.Metrics.Port },
			expectedValue: 3498,
		},
		"Metrics.NetworkMetricsTick": {
			readValueFunc: func(c *Config) interface{} { return c.Metrics.NetworkMetricsTick },
			expectedValue: 43 * time.Second,
		},
		"Metrics.EthereumMetricsTick": {
			readValueFunc: func(c *Config) interface{} { return c.Metrics.EthereumMetricsTick },
			expectedValue: 87 * time.Second,
		},
		"Diagnostics.Port": {
			readValueFunc: func(c *Config) interface{} { return c.Diagnostics.Port },
			expectedValue: 9723,
		},
	}

	for _, filePath := range filePaths {
		t.Run(strings.TrimPrefix(filepath.Ext(filePath), "."), func(t *testing.T) {
			cfg := &Config{}
			if err := cfg.ReadConfig(filePath, nil); err != nil {
				t.Fatalf("failed to read test config: [%v]", err)
			}

			for testName, test := range configReadTests {
				t.Run(testName, func(t *testing.T) {
					expected := test.expectedValue
					actual := test.readValueFunc(cfg)
					if !reflect.DeepEqual(expected, actual) {
						t.Errorf("\nexpected: %s\nactual:   %s", expected, actual)
					}
				})
			}
		})
	}
}

func TestReadConfig_ReadPassword(t *testing.T) {
	expectToPrompt := "expect-to-prompt"

	var configReadTests = map[string]struct {
		configFilePath      string
		envVariablePassword string
		expectedPassword    string
	}{
		"no password in file; password in environment variable": {
			configFilePath:      "../test/config_no_password.toml",
			envVariablePassword: "some-SECRET-phrase-as-ENVIRONMENT-VARIABLE",
			expectedPassword:    "some-SECRET-phrase-as-ENVIRONMENT-VARIABLE",
		},
		"password in file; password in environment variable": {
			configFilePath:      "../test/config.toml",
			envVariablePassword: "some-SECRET-phrase-as-ENVIRONMENT-VARIABLE",
			expectedPassword:    "THIS IS TEST! Password should be defined in env variable or prompt",
		},
		"password in file; no password in environment variable": {
			configFilePath:      "../test/config.toml",
			envVariablePassword: "",
			expectedPassword:    "THIS IS TEST! Password should be defined in env variable or prompt",
		},
		"no password in file; no password in environment variable": {
			configFilePath:      "../test/config_no_password.toml",
			envVariablePassword: "",
			expectedPassword:    expectToPrompt,
		},
	}

	for testName, test := range configReadTests {
		t.Run(testName, func(t *testing.T) {
			if err := os.Setenv(EthereumPasswordEnvVariable, test.envVariablePassword); err != nil {
				t.Fatal(err)
			}

			cfg := &Config{}
			err := cfg.ReadConfig(test.configFilePath, nil)

			if test.expectedPassword != expectToPrompt {
				if err != nil {
					t.Fatalf("failed to read test config: [%v]", err)
				}

				if cfg.Ethereum.Account.KeyFilePassword != test.expectedPassword {
					t.Errorf(
						"\nexpected: %s\nactual:   %s",
						test.expectedPassword,
						cfg.Ethereum.Account.KeyFilePassword,
					)
				}
			} else {
				// When password is not set neither in config file nor environment
				// variable we expect to prompt user to provide the password.
				// The prompt execution fail in unit tests so we assume that
				// if a specific error is thrown the prompt was issued.
				expectedErrors := []string{
					"unable to read password, error [operation not supported by device]",
					"unable to read password, error [inappropriate ioctl for device]",
				}
				if !slices.Contains(expectedErrors, err.Error()) {
					t.Errorf(
						"unexpected error\nexpected: %v\nactual:   %v\n",
						fmt.Sprintf(
							"\n  one of:\n  * %s",
							strings.Join(expectedErrors, "\n  * "),
						),
						err,
					)
				}
			}
		})
	}
}
