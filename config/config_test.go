package config_test

import (
	"os"
	"testing"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/util"
)

func TestReadPeerConfig(t *testing.T) {
	setup(t)
	// This config file has a value for every setting for a non-bootstrap peer.
	cfg, err := config.ReadConfig("./testdata/config.toml")
	util.Ok(t, err)

	configReadTests := map[string]struct {
		expectedValue interface{}
		readValueFunc func(*config.Config) interface{}
	}{
		"Ethereum.URL": {
			expectedValue: "ws://192.168.0.158:8546",
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URL },
		},
		"Ethereum.URLRPC": {
			expectedValue: "http://192.168.0.158:8545",
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URLRPC },
		},
		"Ethereum.Account.Address": {
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b",
			readValueFunc: func(c *config.Config) interface{} {
				return c.Ethereum.Account.Address
			},
		},
		"Ethereum.ContractAddresses": {
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"GroupContract":    "0x139deb0dd975af8e4cc91fe9053a37e4faf37649",
			},
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.ContractAddresses },
		},
		"Node.Port": {
			expectedValue: 27001,
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Port },
		},
		"Node.Peers": {
			expectedValue: []string{
				"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
			},
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Peers },
		},
	}
	for testName, test := range configReadTests {
		t.Run(testName, func(t *testing.T) {
			util.Equals(t,
				test.expectedValue,
				test.readValueFunc(&cfg))
		})
	}
}

func TestReadBootstrapConfig(t *testing.T) {
	setup(t)
	// This config file has a value for every setting for a bootstrap peer.
	cfg, err := config.ReadConfig("./testdata/config.bootstrap.toml")
	util.Ok(t, err)

	configReadTests := map[string]struct {
		expectedValue interface{}
		readValueFunc func(*config.Config) interface{}
	}{
		"Ethereum.URL": {
			expectedValue: "ws://192.168.0.158:8546",
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URL },
		},
		"Ethereum.URLRPC": {
			expectedValue: "http://192.168.0.158:8545",
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URLRPC },
		},
		"Ethereum.Account.Address": {
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b",
			readValueFunc: func(c *config.Config) interface{} {
				return c.Ethereum.Account.Address
			},
		},
		"Ethereum.ContractAddresses": {
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"GroupContract":    "0x139deb0dd975af8e4cc91fe9053a37e4faf37649",
			},
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.ContractAddresses },
		},
		"Node.Port": {
			expectedValue: 27001,
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Port },
		},
		"Bootstrap.Seed": {
			expectedValue: 0,
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Seed },
		},
	}
	for testName, test := range configReadTests {
		t.Run(testName, func(t *testing.T) {
			util.Equals(t,
				test.expectedValue,
				test.readValueFunc(&cfg))
		})
	}
}

func TestReadIpcConfig(t *testing.T) {
	setup(t)
	// This config file has a value for every setting and uses an .ipc address for the ethereum url.
	cfg, err := config.ReadConfig("./testdata/config.ipc.toml")
	util.Ok(t, err)

	configReadTests := map[string]struct {
		expectedValue interface{}
		readValueFunc func(*config.Config) interface{}
	}{
		"Ethereum.URL": {
			expectedValue: "/tmp/geth.ipc",
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URL },
		},
		"Ethereum.URLRPC": {
			expectedValue: "http://192.168.0.158:8545",
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.URLRPC },
		},
		"Ethereum.Account.Address": {
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b",
			readValueFunc: func(c *config.Config) interface{} {
				return c.Ethereum.Account.Address
			},
		},
		"Ethereum.ContractAddresses": {
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"GroupContract":    "0x139deb0dd975af8e4cc91fe9053a37e4faf37649",
			},
			readValueFunc: func(c *config.Config) interface{} { return c.Ethereum.ContractAddresses },
		},
		"Node.Port": {
			expectedValue: 27001,
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Port },
		},
		"Node.Peers": {
			expectedValue: []string{
				"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
			},
			readValueFunc: func(c *config.Config) interface{} { return c.Node.Peers },
		},
	}
	for testName, test := range configReadTests {
		t.Run(testName, func(t *testing.T) {
			util.Equals(t,
				test.expectedValue,
				test.readValueFunc(&cfg))
		})
	}
}

func TestReadInvalidEthereumURLConfig(t *testing.T) {
	setup(t)
	cfg, err := config.NewConfig("./testdata/config.invalid.ethereum.url.toml")
	util.NotOkRead(t, err, "ethereum.URL: %v", cfg.Ethereum.URL)
}

func TestReadInvalidEthereumURLRPCConfig(t *testing.T) {
	setup(t)
	cfg, err := config.NewConfig("./testdata/config.invalid.ethereum.url.rpc.toml")
	util.NotOkRead(t, err, "ethereum.URLRPC: %v", cfg.Ethereum.URLRPC)
}

func TestReadInvalidEthereumAccountAddress(t *testing.T) {
	setup(t)
	cfg, err := config.NewConfig("./testdata/config.invalid.ethereum.account.address.toml")
	util.NotOkRead(t, err, "ethereum.account.Address: %v", cfg.Ethereum.Account.Address)
}

func TestReadInvalidEthereumAccountKeyfile(t *testing.T) {
	setup(t)
	cfg, err := config.NewConfig("./testdata/config.invalid.ethereum.account.keyfile.toml")
	util.NotOkRead(t, err, "ethereum.account.KeyFile: %v", cfg.Ethereum.Account.KeyFile)
}

func TestReadInvalidEthereumContractKeepRandomBeaconAddress(t *testing.T) {
	setup(t)
	cfg, err := config.NewConfig("./testdata/config.invalid.random.beacon.address.toml")
	util.NotOkRead(t, err, "ethereum.ContractAddresses.KeepRandomBeacon: %v", cfg.Ethereum.ContractAddresses["KeepRandomBeacon"])
}

func TestReadInvalidEthereumContractGroupContractAddress(t *testing.T) {
	setup(t)
	cfg, err := config.NewConfig("./testdata/config.invalid.group.contract.address.toml")
	util.NotOkRead(t, err, "ethereum.ContractAddresses.GroupContract: %v", cfg.Ethereum.ContractAddresses["GroupContract"])
}

func TestReadInvalidNodePortConfig(t *testing.T) {
	setup(t)
	_, err := config.NewConfig("./testdata/config.invalid.node.port.toml")
	util.NotOkRead(t, err, "unexpected error: missing value for port; see node section in config file or use --port flag")
}

func TestReadInvalidBootstrapSeedConfig(t *testing.T) {
	setup(t)
	cfg, err := config.NewConfig("./testdata/config.invalid.bootstrap.seed.toml")
	util.NotOkRead(t, err, "bootstrap.Seed: %v", cfg.Node.Seed)
}

func setup(t *testing.T) {
	err := os.Setenv("KEEP_ETHEREUM_PASSWORD", "not-my-password")
	util.Ok(t, err)
}
