package config

import (
	"math/big"
	"os"
	"reflect"
	"testing"
)

func TestReadConfigFromFile(t *testing.T) {
	err := os.Setenv(ethereumPasswordEnvVariable, "not-my-password")
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
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"randombeacon": "0xcf64c2a367341170cb4e09cf8c0ed137d8473ceb",
			},
		},
		"Storage.DataDir": {
			readValueFunc: func(c *Config) interface{} { return c.Storage.DataDir },
			expectedValue: "/my/secure/location",
		},
		"Ethereum.MaxGasPrice": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.MaxGasFeeCap.Int },
			expectedValue: big.NewInt(140000000000),
		},
		"Ethereum.BalanceAlertThreshold": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.BalanceAlertThreshold.Int },
			expectedValue: big.NewInt(2500000000000000000),
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
