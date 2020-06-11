package config

import (
	"os"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	err := os.Setenv("KEEP_ETHEREUM_PASSWORD", "not-my-password")
	if err != nil {
		t.Fatal(err)
	}

	filepath := "../test/config.toml"
	cfg, err := ReadConfig(filepath)
	if err != nil {
		t.Fatalf(
			"failed to read test config: [%v]",
			err,
		)

	}

	var configReadTests = map[string]struct {
		readValueFunc func(*Config) interface{}
		expectedValue interface{}
	}{
		"Ethereum.URL": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URL },
			expectedValue: "ws://192.168.0.158:8546",
		},
		"Ethereum.URLRPC": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URLRPC },
			expectedValue: "http://192.168.0.158:8545",
		},
		"Ethereum.Account.Address": {
			readValueFunc: func(c *Config) interface{} {
				return c.Ethereum.Account.Address
			},
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b",
		},
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"KeepRandomBeaconOperator": "0xcf64c2a367341170cb4e09cf8c0ed137d8473ceb",
			},
		},
		"Storage.DataDir": {
			readValueFunc: func(c *Config) interface{} { return c.Storage.DataDir },
			expectedValue: "/my/secure/location",
		},
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

}
