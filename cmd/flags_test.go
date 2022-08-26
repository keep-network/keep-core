package cmd

import (
	"math/big"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	ethereumBeacon "github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen"
	ethereumEcdsa "github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen"
	ethereumThreshold "github.com/keep-network/keep-core/pkg/chain/ethereum/threshold/gen"
)

var cmdFlagsTests = map[string]struct {
	readValueFunc func(*config.Config) interface{}
	flagName      string
	flagValue     string
	// We provide arguments for flags in `flagValue` as strings, that are unmarshaled
	// to a Config specific types.
	expectedValueFromFlag interface{}
	defaultValue          interface{}
}{
	"ethereum.url": {
		readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URL },
		flagName:      "--ethereum.url",
		flagValue:     "https://eth-provider.com/mainnet",
		defaultValue:  "",
	},
	"ethereum.keyFile": {
		readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.Account.KeyFile },
		flagName:      "--ethereum.keyFile",
		flagValue:     "/tmp/UTC--2018-03-11T01-37-33.202765887Z--c2a56884538778bacd91aa5bf343bf882c5fb18b",
		defaultValue:  "",
	},
	"ethereum.miningCheckInterval": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Ethereum.MiningCheckInterval },
		flagName:              "--ethereum.miningCheckInterval",
		flagValue:             "1m45s",
		expectedValueFromFlag: 105 * time.Second,
		defaultValue:          60 * time.Second,
	},
	"ethereum.requestPerSecondLimit": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Ethereum.RequestsPerSecondLimit },
		flagName:              "--ethereum.requestPerSecondLimit",
		flagValue:             "38",
		expectedValueFromFlag: 38,
		defaultValue:          150,
	},
	"ethereum.concurrencyLimit": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Ethereum.ConcurrencyLimit },
		flagName:              "--ethereum.concurrencyLimit",
		flagValue:             "105",
		expectedValueFromFlag: 105,
		defaultValue:          30,
	},
	"ethereum.maxGasFeeCap": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Ethereum.MaxGasFeeCap.Int },
		flagName:              "--ethereum.maxGasFeeCap",
		flagValue:             "65.8 Gwei",
		expectedValueFromFlag: big.NewInt(65800000000),
		defaultValue:          big.NewInt(500000000000),
	},
	"ethereum.balanceAlertThreshold": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Ethereum.BalanceAlertThreshold.Int },
		flagName:              "--ethereum.balanceAlertThreshold",
		flagValue:             "1.25 ether",
		expectedValueFromFlag: big.NewInt(1250000000000000000),
		defaultValue:          big.NewInt(500000000000000000),
	},
	"network.port": {
		readValueFunc:         func(c *config.Config) interface{} { return c.LibP2P.Port },
		flagName:              "--network.port",
		flagValue:             "78690",
		expectedValueFromFlag: 78690,
		defaultValue:          3919,
	},
	"network.peers": {
		readValueFunc: func(c *config.Config) interface{} { return c.LibP2P.Peers },
		flagName:      "--network.peers",
		flagValue:     `"/ip4/127.0.0.1/tcp/5001/ipfs/3w6C4TFVo","/dns4/domain.local/tcp/3819/ipfs/16xtuKXdTd"`,
		expectedValueFromFlag: []string{
			"/ip4/127.0.0.1/tcp/5001/ipfs/3w6C4TFVo",
			"/dns4/domain.local/tcp/3819/ipfs/16xtuKXdTd",
		},
		defaultValue: []string{},
	},
	"network.announcedAddresses": {
		readValueFunc: func(c *config.Config) interface{} { return c.LibP2P.AnnouncedAddresses },
		flagName:      "--network.announcedAddresses",
		flagValue:     `"/dns4/boar.network/tcp/4200","/ip4/80.70.69.15/tcp/4201"`,
		expectedValueFromFlag: []string{
			"/dns4/boar.network/tcp/4200",
			"/ip4/80.70.69.15/tcp/4201",
		},
		defaultValue: []string{},
	},
	"network.disseminationTime": {
		readValueFunc:         func(c *config.Config) interface{} { return c.LibP2P.DisseminationTime },
		flagName:              "--network.disseminationTime",
		flagValue:             "486",
		expectedValueFromFlag: 486,
		defaultValue:          0,
	},
	"storage.dataDir": {
		readValueFunc: func(c *config.Config) interface{} { return c.Storage.DataDir },
		flagName:      "--storage.dataDir",
		flagValue:     "./flagged/location/dude",
		defaultValue:  "",
	},
	"metrics.port": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Metrics.Port },
		flagName:              "--metrics.port",
		flagValue:             "9870",
		expectedValueFromFlag: 9870,
		defaultValue:          8080,
	},
	"metrics.networkMetricsTick": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Metrics.NetworkMetricsTick },
		flagName:              "--metrics.networkMetricsTick",
		flagValue:             "3m9s",
		expectedValueFromFlag: 189 * time.Second,
		defaultValue:          1 * time.Minute,
	},
	"metrics.ethereumMetricsTick": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Metrics.EthereumMetricsTick },
		flagName:              "--metrics.ethereumMetricsTick",
		flagValue:             "1m16s",
		expectedValueFromFlag: 76 * time.Second,
		defaultValue:          10 * time.Minute,
	},
	"diagnostics.port": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Diagnostics.Port },
		flagName:              "--diagnostics.port",
		flagValue:             "6089",
		expectedValueFromFlag: 6089,
		defaultValue:          8081,
	},
	"tbtc.preParamsPoolSize": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Tbtc.PreParamsPoolSize },
		flagName:              "--tbtc.preParamsPoolSize",
		flagValue:             "75",
		expectedValueFromFlag: 75,
		defaultValue:          3000,
	},
	"tbtc.preParamsGenerationTimeout": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Tbtc.PreParamsGenerationTimeout },
		flagName:              "--tbtc.preParamsGenerationTimeout",
		flagValue:             "2m30s",
		expectedValueFromFlag: 150 * time.Second,
		defaultValue:          120 * time.Second,
	},
	"tbtc.preParamsGenerationDelay": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Tbtc.PreParamsGenerationDelay },
		flagName:              "--tbtc.preParamsGenerationDelay",
		flagValue:             "1m",
		expectedValueFromFlag: 60 * time.Second,
		defaultValue:          10 * time.Second,
	},
	"tbtc.preParamsGenerationConcurrency": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Tbtc.PreParamsGenerationConcurrency },
		flagName:              "--tbtc.preParamsGenerationConcurrency",
		flagValue:             "2",
		expectedValueFromFlag: 2,
		defaultValue:          1,
	},
	"tbtc.keyGenConcurrency": {
		readValueFunc:         func(c *config.Config) interface{} { return c.Tbtc.KeyGenerationConcurrency },
		flagName:              "--tbtc.keyGenerationConcurrency",
		flagValue:             "101",
		expectedValueFromFlag: 101,
		defaultValue:          runtime.GOMAXPROCS(0),
	},
	"developer.randomBeaconAddress": {
		readValueFunc: func(c *config.Config) interface{} {
			address, _ := c.Ethereum.ContractAddress(ethereum.RandomBeaconContractName)
			return address.String()
		},
		flagName:              "--developer.randomBeaconAddress",
		flagValue:             "0x3b292d36468bc7fd481987818ef2e4d28202a0ed",
		expectedValueFromFlag: "0x3B292D36468bC7fd481987818ef2E4d28202A0eD",
		defaultValue:          ethereumBeacon.RandomBeaconAddress,
	},
	"developer.walletRegistryAddress": {
		readValueFunc: func(c *config.Config) interface{} {
			address, _ := c.Ethereum.ContractAddress(ethereum.WalletRegistryContractName)
			return address.String()
		},
		flagName:              "--developer.walletRegistryAddress",
		flagValue:             "0xb76707515c3f908411b5211863a7581589a1e31f",
		expectedValueFromFlag: "0xB76707515C3f908411B5211863A7581589a1E31F",
		defaultValue:          ethereumEcdsa.WalletRegistryAddress,
	},
	"developer.tokenStakingAddress": {
		readValueFunc: func(c *config.Config) interface{} {
			address, _ := c.Ethereum.ContractAddress(ethereum.TokenStakingContractName)
			return address.String()
		},
		flagName:              "--developer.tokenStakingAddress",
		flagValue:             "0x861b021462e7864a7413edf0113030b892978617",
		expectedValueFromFlag: "0x861b021462e7864a7413edF0113030B892978617",
		defaultValue:          ethereumThreshold.TokenStakingAddress,
	},
}

func TestFlags_ReadConfigFromFlags(t *testing.T) {
	testCommand, testConfig, _ := initTestCommand()

	args := []string{}
	for _, test := range cmdFlagsTests {
		args = append(args, []string{test.flagName, test.flagValue}...)
	}
	testCommand.SetArgs(args)

	testCommand.Execute()

	for testName, test := range cmdFlagsTests {
		t.Run(testName, func(t *testing.T) {
			var expected interface{} = test.flagValue

			if test.expectedValueFromFlag != nil {
				expected = test.expectedValueFromFlag
			}

			actual := test.readValueFunc(testConfig)

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("\nexpected: %v\nactual:   %v", expected, actual)
			}
		})
	}
}

func TestFlags_ReadConfigFromFlagsWithDefaults(t *testing.T) {
	testCommand, loadedConfig, _ := initTestCommand()

	args := []string{
		cmdFlagsTests["ethereum.url"].flagName, cmdFlagsTests["ethereum.url"].flagValue,
		cmdFlagsTests["ethereum.keyFile"].flagName, cmdFlagsTests["ethereum.keyFile"].flagValue,
		cmdFlagsTests["storage.dataDir"].flagName, cmdFlagsTests["storage.dataDir"].flagValue,
	}
	testCommand.SetArgs(args)

	testCommand.Execute()

	for testName, test := range cmdFlagsTests {
		t.Run(testName, func(t *testing.T) {
			expected := test.defaultValue
			if slices.Contains(args, test.flagName) {
				expected = test.flagValue
			}

			actual := test.readValueFunc(loadedConfig)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("\nexpected: %s\nactual:   %s", expected, actual)
			}
		})
	}
}

// In this test we test a combination of properties defined in a config file and flags.
func TestFlags_Mixed(t *testing.T) {
	testCommand, testConfig, _ := initTestCommand()

	args := []string{
		"--config", "../test/config_flags.toml",
		"--ethereum.url", "https://api.url.com/123eth",
		"--ethereum.keyFile", "./keyfile-path/from/flag",
		"--network.port", "7469",
	}
	testCommand.SetArgs(args)

	testCommand.Execute()

	tests := map[string]struct {
		readValueFunc func(*config.Config) interface{}
		expectedValue interface{}
	}{
		// Properties not defined in the config file, but set with flags.
		"ethereum.url": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URL },
			expectedValue: "https://api.url.com/123eth",
		},
		// Properties provided in the config file and overwritten by the flags.
		"ethereum.keyFile": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.Account.KeyFile },
			expectedValue: "./keyfile-path/from/flag",
		},
		"network.port": {
			readValueFunc: func(c *config.Config) interface{} { return c.LibP2P.Port },
			expectedValue: 7469,
		},
		// Properties defined in the config file, not set with flags.
		"metrics.port": {
			readValueFunc: func(c *config.Config) interface{} { return c.Metrics.Port },
			expectedValue: 3097,
		},
		"storage.dataDir": {
			readValueFunc: func(c *config.Config) interface{} { return c.Storage.DataDir },
			expectedValue: "/my/secure/location",
		},
		// Properties not provided in the config file nor set with flags. Use defaults.
		"diagnostics.port": {
			readValueFunc: func(c *config.Config) interface{} { return c.Diagnostics.Port },
			expectedValue: 8081,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			expected := test.expectedValue
			actual := test.readValueFunc(testConfig)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("\nexpected: %s\nactual:   %s", expected, actual)
			}
		})
	}
}

func initTestCommand() (*cobra.Command, *config.Config, *string) {
	if err := os.Setenv(config.EthereumPasswordEnvVariable, "password from env var"); err != nil {
		panic(err)
	}

	var testConfigFilePath string
	var testConfig = &config.Config{}

	testCommand := &cobra.Command{
		Use: "Test",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := testConfig.ReadConfig(testConfigFilePath, cmd.Flags(), config.AllCategories...); err != nil {
				logger.Fatalf("error reading config: %v", err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	initFlags(testCommand, &testConfigFilePath, testConfig, config.AllCategories...)

	return testCommand, testConfig, &testConfigFilePath
}
