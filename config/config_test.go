package config

import (
	"os"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	err := os.Setenv("KEEP_ETHEREUM_ACCOUNT_KEYFILEPASSWORD", "not-my-password")
	if err != nil {
		t.Fatal(err)
	}

	filepath := "../test/config.toml"
	cfg, err := ReadConfig(filepath, "", "")
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
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"KeepGroup":        "0xcf64c2a367341170cb4e09cf8c0ed137d8473ceb",
			},
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

func TestConsulReadConfig(t *testing.T) {
	// override config file settings with env variables
	err := os.Setenv(ethPasswordEnvVariable, "not-my-password")
	if err != nil {
		t.Fatal(err)
	}

	// check we can read the test config file
	filepath := "../test/config.toml"
	cfg, err := ReadConfig(filepath, "localhost:8500", "keep-client-standard-peer-0")
	if err != nil {
		t.Fatalf(
			"failed to read test config: [%v]",
			err,
		)
	}

	// test that config file variables return overrides
	var envConfigOverrideTests = map[string]struct {
		readValueFunc func(*Config) interface{}
		expectedValue interface{}
	}{
		"Ethereum.URL": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URL },
			expectedValue: "ws://192.168.0.150:8546",
		},
		"Ethereum.URLRPC": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URLRPC },
			expectedValue: "http://192.168.0.151:8545",
		},
		"Ethereum.Account.Address": {
			readValueFunc: func(c *Config) interface{} {
				return c.Ethereum.Account.Address
			},
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18c",
		},
		"Ethereum.Account.KeyFile": {
			readValueFunc: func(c *Config) interface{} {
				return c.Ethereum.Account.KeyFile
			},
			expectedValue: "/tmp/UTC--2018-03-11T01-37-33.202765887Z--c2a56884538778bacd91aa5bf343bf882c5fb18c",
		},
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37648",
				"KeepGroup":        "0xcf64c2a367341170cb4e09cf8c0ed137d8473cec",
				"StakingProxy":     "0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCD",
			},
		},
		"LibP2P.Port": {
			readValueFunc: func(c *Config) interface{} {
				return c.LibP2P.Port
			},
			expectedValue: 27002,
		},
		"LibP2P.Peers": {
			readValueFunc: func(c *Config) interface{} {
				return c.LibP2P.Peers[0] // extract the first string from the array
			},
			expectedValue: "/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
		},
		"LibP2P.Seed": {
			readValueFunc: func(c *Config) interface{} {
				return c.LibP2P.Seed
			},
			expectedValue: 0,
		},
	}

	// run tests against env variable values
	for testName, test := range envConfigOverrideTests {
		t.Run(testName, func(t *testing.T) {
			expected := test.expectedValue
			actual := test.readValueFunc(cfg)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("\nexpected: %s\nactual:   %s", expected, actual)
			}
		})
	}
}
