package config

import (
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestReadConfigFromFile(t *testing.T) {
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

func TestReadConfigFromEnvVars(t *testing.T) {
	expectedEthereumAccountPassword := "super-secret-PASS"
	if err := os.Setenv("KEEP_ETHEREUM_PASSWORD", expectedEthereumAccountPassword); err != nil {
		t.Fatal(err)
	}

	expectedEthereumAccountKeyFile := "some/path/to/file"
	if err := os.Setenv("KEEP_ETHEREUM_KEYFILE", expectedEthereumAccountKeyFile); err != nil {
		t.Fatal(err)
	}

	expectedEthereumUrl := "ws://infura.com:8546/api"
	if err := os.Setenv("KEEP_ETHEREUM_URL", expectedEthereumUrl); err != nil {
		t.Fatal(err)
	}

	expectedEthereumUrlRpc := "http://infura.com:8546/api"
	if err := os.Setenv("KEEP_ETHEREUM_URLRPC", expectedEthereumUrlRpc); err != nil {
		t.Fatal(err)
	}

	expectedRandomBeaconAddress := "0x0A2CA8909bF14aFC8fE93F6588AB8236CC5DA05F"
	if err := os.Setenv("KEEP_ETHEREUM_CONTRACTS_RANDOMBEACON", expectedRandomBeaconAddress); err != nil {
		t.Fatal(err)
	}

	expectedLibP2PPeers := []string{"/ip4/127.0.0.1/tcp/3919/ipfs/a1", "/ip4/127.0.0.1/tcp/3920/ipfs/a2"}
	if err := os.Setenv("KEEP_PEERS", strings.Join(expectedLibP2PPeers, ",")); err != nil {
		t.Fatal(err)
	}

	expectedLibP2PPort := 420
	if err := os.Setenv("KEEP_PORT", fmt.Sprint(expectedLibP2PPort)); err != nil {
		t.Fatal(err)
	}

	expectedLibP2PAnnouncedAddresses := []string{"/dns4/example.com/tcp/3919", "/ip4/80.70.60.50/tcp/3919"}
	if err := os.Setenv("KEEP_ANNOUNCED_ADDRESSES", strings.Join(expectedLibP2PAnnouncedAddresses, ",")); err != nil {
		t.Fatal(err)
	}

	expectedStorageDataDir := "some/path/to/storage/dir"
	if err := os.Setenv("KEEP_STORAGE_DIR", expectedStorageDataDir); err != nil {
		t.Fatal(err)
	}

	cfg, err := ReadConfig("")
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
			expectedValue: expectedEthereumUrl,
		},
		"Ethereum.URLRPC": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.URLRPC },
			expectedValue: expectedEthereumUrlRpc,
		},
		"Ethereum.Account.KeyFilePassword": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.Account.KeyFilePassword },
			expectedValue: expectedEthereumAccountPassword,
		},

		"Ethereum.Account.KeyFile": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.Account.KeyFile },
			expectedValue: expectedEthereumAccountKeyFile,
		},
		"Ethereum.ContractAddresses": {
			readValueFunc: func(c *Config) interface{} { return c.Ethereum.ContractAddresses },
			expectedValue: map[string]string{
				"randombeacon": expectedRandomBeaconAddress,
			},
		},
		"LibP2P.Peers": {
			readValueFunc: func(c *Config) interface{} { return c.LibP2P.Peers },
			expectedValue: expectedLibP2PPeers,
		},
		"LibP2P.Port": {
			readValueFunc: func(c *Config) interface{} { return c.LibP2P.Port },
			expectedValue: expectedLibP2PPort,
		},
		"LibP2P.AnnouncedAddresses": {
			readValueFunc: func(c *Config) interface{} { return c.LibP2P.AnnouncedAddresses },
			expectedValue: expectedLibP2PAnnouncedAddresses,
		},
		"Storage.DataDir": {
			readValueFunc: func(c *Config) interface{} { return c.Storage.DataDir },
			expectedValue: expectedStorageDataDir,
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
