package config_test

import (
	"os"
	"testing"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/util"
)

func TestReadPeerConfig(t *testing.T) {
	setup(t)

	cfg, err := config.ReadConfig("./testdata/config.toml")
	util.Ok(t, err)

	configReadTests := map[string]struct {
		readValueFunc func(*config.Config) interface{}
		expectedValue interface{}
	}{
		"Ethereum.URL": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URL },
			expectedValue: "ws://192.168.0.158:8546",
		},
		"Ethereum.URLRPC": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URLRPC },
			expectedValue: "http://192.168.0.158:8545",
		},
		"Ethereum.Account.Address": {
			readValueFunc: func(c *config.Config) interface{} {
				return c.Ethereum.Account.Address
			},
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b",
		},
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"GroupContract":    "0x139deb0dd975af8e4cc91fe9053a37e4faf37649",
			},
		},
		"Node.Port": {
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Port },
			expectedValue: 27001,
		},
		"Bootstrap.URLs": {
			readValueFunc: func(c *config.Config) interface{} { return c.Bootstrap.URLs },
			expectedValue: []string{
				"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
			},
		},
	}

	for testName, test := range configReadTests {
		t.Run(testName, func(t *testing.T) {
			expected := test.expectedValue
			actual := test.readValueFunc(&cfg)
			util.Equals(t, expected, actual)
		})
	}
}

func TestReadBootstrapConfig(t *testing.T) {
	setup(t)

	cfg, err := config.ReadConfig("./testdata/config.bootstrap.toml")
	util.Ok(t, err)

	configReadTests := map[string]struct {
		readValueFunc func(*config.Config) interface{}
		expectedValue interface{}
	}{
		"Ethereum.URL": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URL },
			expectedValue: "ws://192.168.0.158:8546",
		},
		"Ethereum.URLRPC": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URLRPC },
			expectedValue: "http://192.168.0.158:8545",
		},
		"Ethereum.Account.Address": {
			readValueFunc: func(c *config.Config) interface{} {
				return c.Ethereum.Account.Address
			},
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b",
		},
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"GroupContract":    "0x139deb0dd975af8e4cc91fe9053a37e4faf37649",
			},
		},
		"Node.Port": {
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Port },
			expectedValue: 27001,
		},
		"Bootstrap.Seed": {
			readValueFunc: func(c *config.Config) interface{} { return c.Bootstrap.Seed },
			expectedValue: 2,
		},
	}

	for testName, test := range configReadTests {
		t.Run(testName, func(t *testing.T) {
			expected := test.expectedValue
			actual := test.readValueFunc(&cfg)
			util.Equals(t, expected, actual)
		})
	}
}

func TestReadIpcConfig(t *testing.T) {
	setup(t)

	cfg, err := config.ReadConfig("./testdata/config.ipc.toml")
	util.Ok(t, err)

	configReadTests := map[string]struct {
		readValueFunc func(*config.Config) interface{}
		expectedValue interface{}
	}{
		"Ethereum.URL": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URL },
			expectedValue: "/tmp/geth.ipc",
		},
		"Ethereum.URLRPC": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URLRPC },
			expectedValue: "http://192.168.0.158:8545",
		},
		"Ethereum.Account.Address": {
			readValueFunc: func(c *config.Config) interface{} {
				return c.Ethereum.Account.Address
			},
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b",
		},
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"GroupContract":    "0x139deb0dd975af8e4cc91fe9053a37e4faf37649",
			},
		},
		"Node.Port": {
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Port },
			expectedValue: 27001,
		},
		"Bootstrap.URLs": {
			readValueFunc: func(c *config.Config) interface{} { return c.Bootstrap.URLs },
			expectedValue: []string{
				"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
			},
		},
	}

	for testName, test := range configReadTests {
		t.Run(testName, func(t *testing.T) {
			expected := test.expectedValue
			actual := test.readValueFunc(&cfg)
			util.Equals(t, expected, actual)
		})
	}
}


func setup(t *testing.T) {
	err := os.Setenv("KEEP_ETHEREUM_PASSWORD", "not-my-password")
	util.Ok(t, util.ErrWrap{ErrNo: config.InvalidErrNo, Err: err})
}
