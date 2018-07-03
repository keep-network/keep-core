package config_test

import (
	"os"
	"testing"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/util"
)

func TestMain(m *testing.M) {
	// Every test calls config.ReadConfig which requires the KEEP_ETHEREUM_PASSWORD environment variable.
	err := os.Setenv("KEEP_ETHEREUM_PASSWORD", "not-my-password")
	if err != nil {
		return
	}
	os.Exit(m.Run())
}

func TestReadPeerConfig(t *testing.T) {
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

func TestReadInvalidConfig(t *testing.T) {
	// This config file has an invalid value for every setting.
	_, err := config.ReadConfig("./testdata/config.invalid.toml")

	expectedErrMsg := `Ethereum.URL (xxx) invalid; format expected: ws://.+|\w.ipc
Ethereum.URLRPC (xxx) invalid; format expected: ^https?:\/\/(.+)\.(.+)
Ethereum.Account.Address (xxx) invalid; format expected: ([13][a-km-zA-HJ-NP-Z1-9]{25,34}|0x[a-fA-F0-9]{40}|\\w+\\.eth(\\W|$)|(?i:iban:)?XE[0-9]{2}[a-zA-Z]{16})|^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$
Ethereum.Account.KeyFile (xxx) invalid; format expected: \/.+
Ethereum.ContractAddresses[KeepRandomBeacon] (xxx) invalid; format expected: ws://.+|\w.ipc
Ethereum.ContractAddresses[GroupContract](xxx) invalid; format expected: ws://.+|\w.ipc
Node.Port (0) invalid; see node section in config file or use --port flag
non-bootstrap node should have Bootstrap.URL and a Bootstrap.Seed of 0
Node.Peers (xxx) invalid; format expected: .+\/.*`

	if err.Error() != expectedErrMsg {
		t.Fatalf("EXPECTED: \n%s\nGOT:\n%s", expectedErrMsg, err.Error())
	}
}
