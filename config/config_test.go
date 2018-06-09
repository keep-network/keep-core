package config

import (
	"os"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	os.Setenv("KEEP_ETHEREUM_PASSWORD", "not-my-password")

	cfg, err := ReadConfig("../test/config.toml")
	if err != nil {
		t.Fatalf("unexpected error [%v]", err)
	}

	expectedURL := "ws://192.168.0.157:8546"
	if cfg.Ethereum.URL != expectedURL {
		t.Errorf(
			`did not correctly read in ./test/config.toml,
			expected [%s] got [%s]\n`,
			expectedURL, cfg.Ethereum.URL,
		)
	}

	expected := "0xc2a56884538778bacd91aa5bf343bf882c5fb18b"
	if actual := cfg.Ethereum.Account.Address; expected != actual {
		t.Errorf("\nexpected: %s\nactual:   %s", expected, actual)
	}

	expectedAddress := "0x639deb0dd975af8e4cc91fe9053a37e4faf37649"
	if vv, ok := cfg.Ethereum.ContractAddresses["KeepRandomBeacon"]; !ok {
		t.Errorf(
			`failed read of test/Config.toml, expected key in map [KeepRandomBeacon].
			 Key missing.\n`,
		)
	} else if vv != expectedAddress {
		t.Errorf(
			"in test/Config.toml\ngot address: %s\nwant address: %s\n",
			vv, expectedAddress,
		)
	}

	var configReadTests = map[string]struct {
		readValueFunc func(*Config) interface{}
		expectedValue interface{}
	}{
		"Ethereum.URL (ws)": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URL },
			expectedValue: "ws://192.168.0.157:8546",
		},
		"Ethereum.URL (ipc)": {
			readValueFunc: func(c *Config) interface{} {
				c.Ethereum.URL = "ipc:/var/folders/ts/7xznj_p13xb7_5th3w6yjmjm0000gn/T/ethereum_dev_mode/geth.ipc"
				return c.Ethereum.URL
			},
			expectedValue: "ipc:/var/folders/ts/7xznj_p13xb7_5th3w6yjmjm0000gn/T/ethereum_dev_mode/geth.ipc",
		},
		"Ethereum.Account.Address": {
			readValueFunc: func(c *Config) interface{} {
				return c.Ethereum.Account.Address
			},
			expectedValue: "0xc2a56884538778bacd91aa5bf343bf882c5fb18b",
		},
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *Config) interface{} {
				return c.Ethereum.ContractAddresses
			},
			expectedValue: map[string]string{
				"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37649",
				"GroupContract":    "0x139deb0dd975af8e4cc91fe9053a37e4faf37649",
			},
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
