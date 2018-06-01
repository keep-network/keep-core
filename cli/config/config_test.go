package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	os.Setenv("KEEP_ETHEREUM_PASSWORD", "not-my-password")

	cfg, err := ReadConfig("../../test/config.toml")
	if err != nil {
		t.Errorf("Error:%s\n", err)
	}
	expectedURL := "ws://192.168.0.157:8546"
	if cfg.Ethereum.URL != expectedURL {
		t.Errorf("Error: Did not correctly read in ./test/Config.toml, Expected [%s] Got [%s]\n", expectedURL, cfg.Ethereum.URL)
	}
	expectedEthereumAccountAddress := "0xc2a56884538778bacd91aa5bf343bf882c5fb18b"
	if vv := cfg.Ethereum.Account.Address; vv != expectedEthereumAccountAddress {
		t.Errorf("in test/Config.toml\ngot  Ethereum.Account.Address: %s\nwant Ethereum.Account.Address: %s\n", vv, expectedEthereumAccountAddress)
	}
	expectedAddress := "0x639deb0dd975af8e4cc91fe9053a37e4faf37649"
	if vv, ok := cfg.Ethereum.ContractAddresses["KeepRandomBeacon"]; !ok {
		t.Errorf("failed read of test/Config.toml, expected key in map [KeepRandomBeacon].  Key missing.\n")
	} else if vv != expectedAddress {
		t.Errorf("in test/Config.toml\ngot address: %s\nwant address: %s\n", vv, expectedAddress)
	}

}
