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
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"GroupContract":    "0x139deb0dd975af8e4cc91fe9053a37e4faf37649",
			},
		},
		"Bootstrap.URL": {
			readValueFunc: func(c *Config) interface{} { return c.Bootstrap.URL },
			expectedValue: []string{"/ip4/127.0.0.1/tcp/8080/ipfs/12D3KooWHHzSeKaY8xuZVzkLbKFfvNgPPeKhFBGrMbNzbm5akpqu"},
		},
		"Bootstrap.Seed": {
			readValueFunc: func(c *Config) interface{} { return c.Bootstrap.Seed },
			expectedValue: 2,
		},
	}

	for testName, test := range configReadTests {
		t.Run(testName, func(t *testing.T) {
			expected := test.expectedValue
			actual := test.readValueFunc(&cfg)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("\nexpected: %s\nactual:   %s", expected, actual)
			}
		})
	}

}
