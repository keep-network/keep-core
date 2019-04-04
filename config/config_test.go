package config

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

const ethPasswordEnvVariable = "KEEP_ETHEREUM_ACCOUNT_KEYFILEPASSWORD"
const ethAccountEnvVariable = "KEEP_ETHEREUM_ACCOUNT_ADDRESS"
const ethKeyfileEnvVariable = "KEEP_ETHEREUM_ACCOUNT_KEYFILE"
const ethRandomBeaconContractEnvVariable = "KEEP_ETHEREUM_CONTRACTADDRESSES_KEEPRANDOMBEACON"
const ethKeepGroupContractEnvVariable = "KEEP_ETHEREUM_CONTRACTADDRESSES_KEEPGROUP"
const ethStakingProxyContractEnvVariable = "KEEP_ETHEREUM_CONTRACTADDRESSES_STAKINGPROXY"
const ethURLEnvVariable = "KEEP_ETHEREUM_URL"
const ethURLRPCEnvVariable = "KEEP_ETHEREUM_URLRPC"
const libp2pPortEnvVariable = "KEEP_LIBP2P_PORT"
const libp2pPeersEnvVariable = "KEEP_LIBP2P_PEERS"
const libp2pSeedEnvVariable = "KEEP_LIBP2P_SEED"

func TestReadConfig(t *testing.T) {
	// check password env variable
	err := os.Setenv(ethPasswordEnvVariable, "not-my-password")
	if err != nil {
		t.Fatal(err)
	}
	if ethEnvPassword, ok := os.LookupEnv("KEEP_ETHEREUM_PASSWORD"); ok {
		if ethEnvPassword != "not-my-password" {
			t.Fatalf("Environment variable [%s] doesn't match!", "KEEP_ETHEREUM_PASSWORD")
		}
	}
	// check we can read the test config file
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
		"Ethereum.Account.KeyFile": {
			readValueFunc: func(c *Config) interface{} {
				return c.Ethereum.Account.KeyFile
			},
			expectedValue: "/tmp/UTC--2018-03-11T01-37-33.202765887Z--c2a56884538778bacd91aa5bf343bf882c5fb18b",
		},
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"KeepGroup":        "0xcf64c2a367341170cb4e09cf8c0ed137d8473ceb",
				"StakingProxy":     "0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC",
			},
		},
	}
	// run tests against config file values
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

func TestEnvReadConfig(t *testing.T) {
	// override config file settings with env variables
	err := os.Setenv(ethURLEnvVariable, "ws://192.168.0.159:8546")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(ethURLRPCEnvVariable, "http://192.168.0.159:8545")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(ethAccountEnvVariable, "0xc2a56884538778bacd91aa5bf343bf882c5fb18c")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(ethKeyfileEnvVariable, "not-my-keyfile-at-all")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(ethRandomBeaconContractEnvVariable, "0x639deb0dd975af8e4cc91fe9053a37e4faf37648")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(ethKeepGroupContractEnvVariable, "0xcf64c2a367341170cb4e09cf8c0ed137d8473cec")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(ethStakingProxyContractEnvVariable, "0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCD")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(libp2pPortEnvVariable, "27002")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(libp2pPeersEnvVariable, "[\"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA\"]")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv(libp2pSeedEnvVariable, "0")
	if err != nil {
		t.Fatal(err)
	}

	// check we can read the test config file
	filepath := "../test/config.toml"
	cfg, err := ReadConfig(filepath)
	if err != nil {
		t.Fatalf(
			"failed to read test config: [%v]",
			err,
		)
	}

	fmt.Println(cfg.Ethereum.ContractAddresses)
	fmt.Println(cfg.LibP2P.Port)

	// test that config file variables return overrides
	var envConfigOverrideTests = map[string]struct {
		readValueFunc func(*Config) interface{}
		expectedValue interface{}
	}{
		"Ethereum.URL": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URL },
			expectedValue: "ws://192.168.0.159:8546",
		},
		"Ethereum.URLRPC": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URLRPC },
			expectedValue: "http://192.168.0.159:8545",
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
			expectedValue: "not-my-keyfile-at-all",
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
